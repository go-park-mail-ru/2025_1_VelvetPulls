package auth

import (
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/utils"
)

type AuthUsecaseInterface interface {
	RegisterUser(values model.RegisterCredentials) (string, error)
	LoginUser(values model.LoginCredentials) (string, error)
}

// authUsecase реализует интерфейс AuthUsecase.
type AuthUsecase struct {
	userRepo    repository.UserRepoInterface
	sessionRepo repository.SessionRepoInterface
}

// NewAuthUsecase создает новый экземпляр AuthUsecase.
func NewAuthUsecase(userRepo repository.UserRepoInterface, sessionRepo repository.SessionRepoInterface) AuthUsecaseInterface {
	return &AuthUsecase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func parseRegisterCredentials(values model.RegisterCredentials) error {
	if values.Password != values.ConfirmPassword {
		return apperrors.ErrPasswordsDoNotMatch
	}

	if !utils.ValidateRegistrationPassword(values.Password) {
		return apperrors.ErrInvalidPassword
	}

	if !utils.ValidateRegistrationPhone(values.Phone) {
		return apperrors.ErrInvalidPhoneFormat
	}

	if !utils.ValidateRegistrationUsername(values.Username) {
		return apperrors.ErrInvalidUsername
	}

	return nil
}

// RegisterUser регистрирует нового пользователя.
func (uc *AuthUsecase) RegisterUser(values model.RegisterCredentials) (string, error) {
	if err := parseRegisterCredentials(values); err != nil {
		return "", err
	}

	if _, err := uc.userRepo.GetUserByUsername(values.Username); err == nil {
		return "", apperrors.ErrUsernameTaken
	} else if err != apperrors.ErrUserNotFound {
		return "", err
	}

	if _, err := uc.userRepo.GetUserByPhone(values.Phone); err == nil {
		return "", apperrors.ErrPhoneTaken
	} else if err != apperrors.ErrUserNotFound {
		return "", err
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

	if err := uc.userRepo.CreateUser(user); err != nil {
		return "", err
	}

	sessionID, err := uc.sessionRepo.CreateSession(user.Username)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

// LoginUser аутентифицирует пользователя и создает сессию.
func (uc *AuthUsecase) LoginUser(values model.LoginCredentials) (string, error) {
	user, err := uc.userRepo.GetUserByUsername(values.Username)
	if err != nil {
		return "", apperrors.ErrUserNotFound
	}

	if !utils.ValidatePassword(user.Password, values.Password) {
		return "", apperrors.ErrInvalidCredentials
	}

	sessionID, err := uc.sessionRepo.CreateSession(user.Username)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}
