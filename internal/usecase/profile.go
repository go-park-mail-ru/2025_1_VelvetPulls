package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
)

type IUserUsecase interface {
	GetUserProfileByID(ctx context.Context, id uuid.UUID) (*model.GetUserProfile, error)
	GetUserProfileByUsername(ctx context.Context, username string) (*model.GetUserProfile, error)
	UpdateUserProfile(ctx context.Context, profile *model.UpdateUserProfile) error
}

type UserUsecase struct {
	userRepo repository.IUserRepo
}

func NewUserUsecase(userRepo repository.IUserRepo) IUserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func (uc *UserUsecase) GetUserProfileByID(ctx context.Context, id uuid.UUID) (*model.GetUserProfile, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("Fetching user profile")

	user, err := uc.userRepo.GetUserByID(ctx, id)
	if err != nil {
		logger.Error("Error fetching user profile")
		return nil, err
	}

	profile := &model.GetUserProfile{
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Username:   user.Username,
		Phone:      user.Phone,
		Email:      user.Email,
		AvatarPath: user.AvatarPath,
	}

	return profile, nil
}

func (uc *UserUsecase) GetUserProfileByUsername(ctx context.Context, username string) (*model.GetUserProfile, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("Fetching user profile")

	user, err := uc.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		logger.Error("Error fetching user profile")
		return nil, err
	}

	profile := &model.GetUserProfile{
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Username:   user.Username,
		Phone:      user.Phone,
		Email:      user.Email,
		AvatarPath: user.AvatarPath,
	}

	return profile, nil
}

func (uc *UserUsecase) UpdateUserProfile(ctx context.Context, profile *model.UpdateUserProfile) error {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("Updating user profile")

	if err := profile.Validate(); err != nil {
		logger.Error("Validation failed")
		return err
	}

	if profile.Avatar != nil {
		if !utils.IsImageFile(*profile.Avatar) {
			logger.Error("Invalid avatar file type")
			return utils.ErrNotImage
		}
	}

	avatarNewURL, avatarOldURL, err := uc.userRepo.UpdateUser(ctx, profile)
	if err != nil {
		logger.Error("Error updating user profile")
		return err
	}

	// Если есть новый аватар, сохраняем его и удаляем старый
	if avatarNewURL != "" && profile.Avatar != nil {
		if err := utils.RewritePhoto(*profile.Avatar, avatarNewURL); err != nil {
			logger.Error("Error rewriting photo")
			return err
		}
		if avatarOldURL != "" {
			go func() {
				if err := utils.RemovePhoto(avatarOldURL); err != nil {
					logger.Error("Error removing old avatar")
				}
			}()
		}
	}

	return nil
}
