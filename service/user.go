package service

import (
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/utils"
)

func RegisterUser(values model.RegisterCredentials) (UserResponse, error) {
	if values.Password != values.ConfirmPassword {
		return UserResponse{
			StatusCode: 400,
			Body:       apperrors.ErrPasswordsDoNotMatch,
		}, apperrors.ErrPasswordsDoNotMatch
	}
	_, err := repository.GetUserByUsername(values.Username)
	if err == nil {
		return UserResponse{
			StatusCode: 400,
			Body:       apperrors.ErrUserAlreadyExists,
		}, apperrors.ErrUserAlreadyExists
	} else if err != apperrors.ErrUserNotFound {
		return UserResponse{
			StatusCode: 500,
			Body:       err,
		}, err
	}

	_, err = repository.GetUserByPhone(values.Phone)
	if err == nil {
		return UserResponse{
			StatusCode: 400,
			Body:       apperrors.ErrPhoneTaken,
		}, apperrors.ErrPhoneTaken
	} else if err != apperrors.ErrUserNotFound {
		return UserResponse{
			StatusCode: 500,
			Body:       err,
		}, err
	}

	hashedPassword, err := utils.HashAndSalt(values.Password)
	if err != nil {
		return UserResponse{
			StatusCode: 500,
			Body:       err,
		}, err
	}

	user := model.User{
		Username: values.Username,
		Password: hashedPassword,
		Phone:    values.Phone,
	}

	err = repository.CreateUser(user)
	if err != nil {
		return UserResponse{
			StatusCode: 500,
			Body:       apperrors.ErrUserCreation,
		}, apperrors.ErrUserCreation
	}

	sessionId, err := repository.CreateSession(user.Username)
	if err != nil {
		return UserResponse{
			StatusCode: 500,
			Body:       err,
		}, err
	}

	return UserResponse{
		StatusCode: 201,
		Body:       sessionId,
	}, nil
}

func LoginUser(values model.LoginCredentials) (UserResponse, error) {
	user, err := repository.GetUserByUsername(values.Username)
	if err != nil {
		return UserResponse{
			StatusCode: 400,
			Body:       apperrors.ErrUsernameTaken,
		}, apperrors.ErrUserNotFound
	}

	if !utils.ValidatePassword(user.Password, values.Password) {
		return UserResponse{
			StatusCode: 400,
			Body:       apperrors.ErrInvalidCredentials,
		}, apperrors.ErrInvalidCredentials
	}

	sessionId, err := repository.CreateSession(values.Username)
	if err != nil {
		return UserResponse{
			StatusCode: 500,
			Body:       err,
		}, err
	}

	return UserResponse{
		StatusCode: 201,
		Body:       sessionId,
	}, nil
}
