package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// IUserUsecase описывает операции над профилем пользователя
type IUserUsecase interface {
	GetUserProfileByID(ctx context.Context, id uuid.UUID) (*model.GetUserProfile, error)
	GetUserProfileByUsername(ctx context.Context, username string) (*model.GetUserProfile, error)
	UpdateUserProfile(ctx context.Context, profile *model.UpdateUserProfile) error
}

// UserUsecase реализует логику работы с профилем пользователя
type UserUsecase struct {
	userRepo repository.IUserRepo
}

// NewUserUsecase создаёт экземпляр UserUsecase
func NewUserUsecase(userRepo repository.IUserRepo) IUserUsecase {
	return &UserUsecase{userRepo: userRepo}
}

// fetchProfile возвращает унифицированную модель GetUserProfile
func (uc *UserUsecase) fetchProfile(ctx context.Context, user *model.User) *model.GetUserProfile {
	return &model.GetUserProfile{
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Username:   user.Username,
		Phone:      user.Phone,
		Email:      user.Email,
		AvatarPath: user.AvatarPath,
	}
}

// GetUserProfileByID возвращает профиль по ID
func (uc *UserUsecase) GetUserProfileByID(ctx context.Context, id uuid.UUID) (*model.GetUserProfile, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("GetUserProfileByID start", zap.String("userID", id.String()))

	user, err := uc.userRepo.GetUserByID(ctx, id)
	if err != nil {
		logger.Error("GetUserByID failed", zap.Error(err))
		return nil, err
	}
	profile := uc.fetchProfile(ctx, user)
	return profile, nil
}

// GetUserProfileByUsername возвращает профиль по username
func (uc *UserUsecase) GetUserProfileByUsername(ctx context.Context, username string) (*model.GetUserProfile, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("GetUserProfileByUsername start", zap.String("username", username))

	user, err := uc.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		logger.Error("GetUserByUsername failed", zap.Error(err))
		return nil, err
	}
	profile := uc.fetchProfile(ctx, user)
	return profile, nil
}

// UpdateUserProfile обновляет данные пользователя и аватар
func (uc *UserUsecase) UpdateUserProfile(ctx context.Context, req *model.UpdateUserProfile) error {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("UpdateUserProfile start", zap.String("userID", req.ID.String()))

	// валидация входных данных
	if err := req.Validate(); err != nil {
		logger.Error("Validation failed", zap.Error(err))
		return err
	}

	// проверка формата аватара
	if req.Avatar != nil {
		if !utils.IsImageFile(*req.Avatar) {
			logger.Error("Invalid avatar file type")
			return utils.ErrNotImage
		}
	}

	// обновление в репозитории
	newURL, oldURL, err := uc.userRepo.UpdateUser(ctx, req)
	if err != nil {
		logger.Error("UpdateUser failed", zap.Error(err))
		return err
	}

	// если появился новый аватар — сохраняем и удаляем старый
	if req.Avatar != nil && newURL != "" {
		if err := utils.RewritePhoto(*req.Avatar, newURL); err != nil {
			logger.Error("RewritePhoto failed", zap.Error(err))
			return err
		}
		uc.handleAvatarCleanup(oldURL)
	}

	logger.Info("UpdateUserProfile done", zap.String("userID", req.ID.String()))
	return nil
}

// --- Private Helpers ---

func (uc *UserUsecase) handleAvatarCleanup(oldURL string) {
	if oldURL != "" {
		go func(url string) {
			if err := utils.RemovePhoto(url); err != nil {
				zap.L().Warn("Old avatar remove failed", zap.String("url", url), zap.Error(err))
			}
		}(oldURL)
	}
}
