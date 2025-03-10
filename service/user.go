package service

import (
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/utils"
)

// RegisterUser регистрирует нового пользователя.
func RegisterUser(values model.RegisterCredentials) (string, error) {
	if values.Password != values.ConfirmPassword {
		return "", apperrors.ErrPasswordsDoNotMatch
	}

	if !utils.ValidateRegistrationPassword(values.Password) {
		return "", apperrors.ErrInvalidPassword
	}

	if !utils.ValidateRegistrationPhone(values.Phone) {
		return "", apperrors.ErrInvalidPhoneFormat
	}

	if !utils.ValidateRegistrationUsername(values.Username) {
		return "", apperrors.ErrInvalidUsername
	}

	if _, err := repository.GetUserByUsername(values.Username); err == nil {
		return "", apperrors.ErrUserAlreadyExists
	} else if err != apperrors.ErrUserNotFound {
		return "", err
	}

	if _, err := repository.GetUserByPhone(values.Phone); err == nil {
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

	if err := repository.CreateUser(user); err != nil {
		return "", apperrors.ErrUserCreation
	}

	sessionID, err := repository.CreateSession(user.Username)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

// LoginUser аутентифицирует пользователя и создает сессию.
func LoginUser(values model.LoginCredentials) (string, error) {
	user, err := repository.GetUserByUsername(values.Username)
	if err != nil {
		return "", apperrors.ErrUserNotFound
	}

	if !utils.ValidatePassword(user.Password, values.Password) {
		return "", apperrors.ErrInvalidCredentials
	}

	sessionID, err := repository.CreateSession(user.Username)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}
