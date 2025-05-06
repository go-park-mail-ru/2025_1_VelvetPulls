package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config/metrics"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/internal/repository"
	"go.uber.org/zap"
)

type IAuthUsecase interface {
	RegisterUser(ctx context.Context, values model.RegisterCredentials) (string, error)
	LoginUser(ctx context.Context, values model.LoginCredentials) (string, error)
	LogoutUser(ctx context.Context, sessionId string) error
}

type AuthUsecase struct {
	authRepo    repository.IAuthRepo
	sessionRepo repository.ISessionRepo
}

func NewAuthUsecase(authRepo repository.IAuthRepo, sessionRepo repository.ISessionRepo) IAuthUsecase {
	return &AuthUsecase{
		authRepo:    authRepo,
		sessionRepo: sessionRepo,
	}
}

func (uc *AuthUsecase) RegisterUser(ctx context.Context, values model.RegisterCredentials) (string, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("Registering new user")

	if err := values.Validate(); err != nil {
		logger.Error("Validation failed", zap.Error(err))
		return "", err
	}

	if _, err := uc.authRepo.GetUserByUsername(ctx, values.Username); err == nil {
		return "", ErrUsernameIsTaken
	}

	if _, err := uc.authRepo.GetUserByPhone(ctx, values.Phone); err == nil {
		return "", ErrPhoneIsTaken
	}

	hashedPassword, err := utils.HashAndSalt(values.Password)
	if err != nil {
		logger.Error("Error hashing password")
		return "", ErrHashPassword
	}

	user := &model.User{
		Username: values.Username,
		Password: hashedPassword,
		Phone:    values.Phone,
	}

	userID, err := uc.authRepo.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}
	sessionID, err := uc.sessionRepo.CreateSession(ctx, userID)
	if err != nil {
		return "", err
	}
	metrics.IncBusinessOp("registration")
	return sessionID, nil
}

func (uc *AuthUsecase) LoginUser(ctx context.Context, values model.LoginCredentials) (string, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("User login attempt")
	if err := values.Validate(); err != nil {
		logger.Error("Validation failed", zap.Error(err))
		return "", err
	}

	user, err := uc.authRepo.GetUserByUsername(ctx, values.Username)
	if err != nil {
		return "", ErrInvalidUsername
	}

	if !utils.CheckPassword(user.Password, values.Password) {
		logger.Error("Invalid password")
		return "", ErrInvalidPassword
	}

	sessionID, err := uc.sessionRepo.CreateSession(ctx, user.ID)
	if err != nil {
		return "", err
	}

	metrics.IncBusinessOp("login")
	return sessionID, nil
}

func (uc *AuthUsecase) LogoutUser(ctx context.Context, sessionId string) error {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("User logout")

	metrics.IncBusinessOp("logout")
	return uc.sessionRepo.DeleteSession(ctx, sessionId)
}
