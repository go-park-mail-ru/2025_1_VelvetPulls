package service

import (
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/errors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/repository"
)

func RegisterUser(user model.User) (UserResponse, error) {
	// Проверяем, существует ли уже такой пользователь по Email
	_, err := repository.GetUserByEmail(user.Email)
	if err == nil { // Если пользователь найден, значит Email уже занят
		return UserResponse{
			StatusCode: 400,
			Body:       errors.ErrEmailTaken,
		}, errors.ErrEmailTaken
	} else if err != errors.ErrUserNotFound { // Если другая ошибка (например, БД)
		return UserResponse{
			StatusCode: 500,
			Body:       err,
		}, err
	}

	// Проверяем, существует ли уже такой пользователь по Phone
	_, err = repository.GetUserByPhone(user.Phone)
	if err == nil {
		return UserResponse{
			StatusCode: 400, // Bad Request
			Body:       errors.ErrPhoneTaken,
		}, errors.ErrPhoneTaken
	} else if err != errors.ErrUserNotFound {
		return UserResponse{
			StatusCode: 500,
			Body:       err,
		}, err
	}

	// Создаем нового пользователя
	err = repository.CreateUser(user)
	if err != nil {
		return UserResponse{
			StatusCode: 500,
			Body:       errors.ErrUserCreation,
		}, errors.ErrUserCreation
	}

	// Если все прошло успешно, возвращаем успешный ответ
	return UserResponse{
		StatusCode: 201,
		Body:       user,
	}, nil
}
