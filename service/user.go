package service

import (
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/errors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/repository"
)

func RegisterUser(user model.User) (UserResponse, error) {
	// Проверяем, существует ли уже такой пользователь
	existingUser, err := repository.GetUserByID(user.ID)
	if err == nil && existingUser.ID != 0 {
		return UserResponse{
			StatusCode: 400, // Bad Request
			Body:       err,
		}, err
	}

	// Создаем нового пользователя
	err = repository.CreateUser(user)
	if err != nil {
		return UserResponse{
			StatusCode: 500, // Internal Server Error
			Body:       errors.ErrUserCreation,
		}, errors.ErrUserCreation
	}

	// Если все прошло успешно, возвращаем успешный ответ
	return UserResponse{
		StatusCode: 201, // Created
		Body:       user,
	}, nil
}
