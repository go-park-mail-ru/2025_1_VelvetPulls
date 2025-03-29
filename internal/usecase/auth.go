package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
)

type IAuthUsecase interface {
	RegisterUser(ctx context.Context, values model.RegisterCredentials) (string, error)
	LoginUser(ctx context.Context, values model.LoginCredentials) (string, error)
	LogoutUser(ctx context.Context, sessionId string) error
}

type AuthUsecase struct {
	userRepo    repository.IUserRepo
	sessionRepo repository.ISessionRepo
}

func NewAuthUsecase(userRepo repository.IUserRepo, sessionRepo repository.ISessionRepo) IAuthUsecase {
	return &AuthUsecase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (uc *AuthUsecase) RegisterUser(ctx context.Context, values model.RegisterCredentials) (string, error) {
	if err := values.Validate(); err != nil {
		return "", err
	}

	if _, err := uc.userRepo.GetUserByUsername(ctx, values.Username); err == nil {
		return "", apperrors.ErrUsernameTaken
	}

	if _, err := uc.userRepo.GetUserByPhone(ctx, values.Phone); err == nil {
		return "", apperrors.ErrPhoneTaken
	}

	hashedPassword, err := utils.HashAndSalt(values.Password)
	if err != nil {
		return "", err
	}

	user := &model.User{
		Username: values.Username,
		Password: hashedPassword,
		Phone:    values.Phone,
	}

	userID, err := uc.userRepo.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}

	sessionID, err := uc.sessionRepo.CreateSession(ctx, userID)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (uc *AuthUsecase) LoginUser(ctx context.Context, values model.LoginCredentials) (string, error) {
	if err := values.Validate(); err != nil {
		return "", err
	}

	user, err := uc.userRepo.GetUserByUsername(ctx, values.Username)
	if err != nil {
		return "", apperrors.ErrUserNotFound
	}

	if !utils.CheckPassword(user.Password, values.Password) {
		return "", apperrors.ErrInvalidCredentials
	}

	sessionID, err := uc.sessionRepo.CreateSession(ctx, user.ID.String())
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (uc *AuthUsecase) LogoutUser(ctx context.Context, sessionId string) error {
	return uc.sessionRepo.DeleteSession(ctx, sessionId)
}
