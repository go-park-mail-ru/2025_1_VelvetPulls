package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	servererrors "github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/server_errors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
)

type IUserUsecase interface {
	GetUserProfile(ctx context.Context, id uuid.UUID) (*model.GetUserProfile, error)
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

func (uc *UserUsecase) GetUserProfile(ctx context.Context, id uuid.UUID) (*model.GetUserProfile, error) {
	user, err := uc.userRepo.GetUserByID(ctx, id)
	if err != nil {
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
	if err := profile.Validate(); err != nil {
		return servererrors.ErrValidation
	}

	if profile.Avatar != nil {
		if !utils.IsImageFile(*profile.Avatar) {
			return utils.ErrNotImage
		}
	}
	avatarNewURL, avatarOldURL, err := uc.userRepo.UpdateUser(ctx, profile)
	if err != nil {
		return err
	}

	if avatarNewURL != nil && profile.Avatar != nil {
		if err := utils.RewritePhoto(*profile.Avatar, *avatarNewURL); err != nil {
			return err
		}
		if avatarOldURL != nil {
			go func() {
				if err := utils.RemovePhoto(*avatarOldURL); err != nil {
					// TODO: log
				}
			}()
		}
	}

	return nil
}
