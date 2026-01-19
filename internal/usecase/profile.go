package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config/metrics"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IUserModule interface {
	GetUserProfileByID(ctx context.Context, id uuid.UUID) (*model.GetUserProfile, error)
	GetUserProfileByUsername(ctx context.Context, username string) (*model.GetUserProfile, error)
	UpdateUserProfile(ctx context.Context, profile *model.UpdateUserProfile) error
}

type UserModule struct {
	userMapper repository.IUserMapper
}

func NewUserModule(userMapper repository.IUserMapper) IUserModule {
	return &UserModule{userMapper: userMapper}
}

func (uc *UserModule) fetchProfile(ctx context.Context, user *model.User) *model.GetUserProfile {
	profile := &model.GetUserProfile{
		ID:         user.ID,
		AvatarPath: user.AvatarPath,
		Name:       user.Name,
		Username:   user.Username,
		BirthDate:  user.BirthDate,
	}

	return profile
}

func (uc *UserModule) GetUserProfileByID(ctx context.Context, id uuid.UUID) (*model.GetUserProfile, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("GetUserProfileByID start", zap.String("userID", id.String()))

	user, err := uc.userMapper.GetUserByID(ctx, id)
	if err != nil {
		logger.Error("GetUserByID failed", zap.Error(err))
		return nil, err
	}
	profile := uc.fetchProfile(ctx, user)
	metrics.IncBusinessOp("get_self_profile")
	return profile, nil
}

func (uc *UserModule) GetUserProfileByUsername(ctx context.Context, username string) (*model.GetUserProfile, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("GetUserProfileByUsername start", zap.String("username", username))

	user, err := uc.userMapper.GetUserByUsername(ctx, username)
	if err != nil {
		logger.Error("GetUserByUsername failed", zap.Error(err))
		return nil, err
	}
	profile := uc.fetchProfile(ctx, user)
	metrics.IncBusinessOp("get_profile")
	return profile, nil
}

func (uc *UserModule) UpdateUserProfile(ctx context.Context, req *model.UpdateUserProfile) error {
	logger := utils.GetLoggerFromCtx(ctx)
	if req == nil {
		logger.Error("UpdateUserProfile received nil request")
		return repository.ErrEmptyField
	}

	logger.Info("UpdateUserProfile start", zap.String("userID", req.ID.String()))

	if req.ID == uuid.Nil {
		logger.Error("Invalid UUID")
		return repository.ErrInvalidUUID
	}

	if err := req.Validate(); err != nil {
		logger.Error("Validation failed", zap.Error(err))
		return repository.ErrInvalidInput
	}

	if req.Avatar != nil {
		if !utils.IsImageFile(*req.Avatar) {
			logger.Error("Invalid avatar file type")
			return repository.ErrInvalidInput
		}
	}

	newURL, oldURL, err := uc.userMapper.UpdateUser(ctx, req)
	if err != nil {
		logger.Error("UpdateUser failed", zap.Error(err))
		return err
	}

	if req.Avatar != nil && newURL != "" {
		if err := utils.RewritePhoto(*req.Avatar, newURL); err != nil {
			logger.Error("RewritePhoto failed", zap.Error(err))
			return repository.ErrDatabaseOperation
		}
		uc.handleAvatarCleanup(oldURL)
	}

	logger.Info("UpdateUserProfile done", zap.String("userID", req.ID.String()))
	metrics.IncBusinessOp("update_profile")
	return nil
}

// --- Private Helpers ---

func (uc *UserModule) handleAvatarCleanup(oldURL string) {
	if oldURL != "" {
		go func(url string) {
			if err := utils.RemovePhoto(url); err != nil {
				zap.L().Warn("Old avatar remove failed", zap.String("url", url), zap.Error(err))
			}
		}(oldURL)
	}
}
