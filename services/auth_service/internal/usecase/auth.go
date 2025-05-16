package usecase

import (
	"context"
	"fmt"

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

	// Валидируем введенные данные
	if err := values.Validate(); err != nil {
		logger.Error("Validation failed", zap.Error(err))
		return "", err
	}

	// Проверяем, не занят ли логин
	if _, err := uc.authRepo.GetUserByUsername(ctx, values.Username); err == nil {
		return "", ErrUsernameIsTaken
	}

	// Хешируем пароль
	hashedPassword, err := utils.HashAndSalt(values.Password)
	if err != nil {
		logger.Error("Error hashing password")
		return "", ErrHashPassword
	}

	// Создаем пользователя
	user := &model.User{
		Username: values.Username,
		Password: hashedPassword,
		Name:     values.Name,
	}

	// Создаем пользователя в репозитории
	userID, err := uc.authRepo.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}

	// Создаем сессию
	sessionID, err := uc.sessionRepo.CreateSession(ctx, userID)
	if err != nil {
		return "", err
	}

	// Увеличиваем счетчик операций
	metrics.IncBusinessOp("registration")
	return sessionID, nil
}

func (uc *AuthUsecase) LoginUser(ctx context.Context, values model.LoginCredentials) (string, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("User login attempt")

	// Валидируем данные для логина
	if err := values.Validate(); err != nil {
		logger.Error("Validation failed", zap.Error(err))
		return "", err
	}

	fmt.Println(values.Username)
	// Находим пользователя по логину
	user, err := uc.authRepo.GetUserByUsername(ctx, values.Username)
	if err != nil {
		return "", ErrInvalidUsername
	}
	fmt.Println(user.Password, values.Password)

	// Проверяем пароль
	if !utils.CheckPassword(user.Password, values.Password) {
		logger.Error("Invalid password")
		return "", ErrInvalidPassword
	}

	// Создаем сессию
	sessionID, err := uc.sessionRepo.CreateSession(ctx, user.ID)
	if err != nil {
		return "", err
	}

	// Увеличиваем счетчик операций
	metrics.IncBusinessOp("login")
	return sessionID, nil
}

func (uc *AuthUsecase) LogoutUser(ctx context.Context, sessionId string) error {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("User logout")

	// Удаляем сессию
	metrics.IncBusinessOp("logout")
	return uc.sessionRepo.DeleteSession(ctx, sessionId)
}
