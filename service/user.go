package service

import (
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/models"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/repository"
)

func RegisterUser(user models.User) (UserResponse, error) {
	// Проверяем, существует ли уже такой пользователь по Email
	_, err := repository.GetUserByEmail(user.Email)
	if err == nil { // Если пользователь найден, значит Email уже занят
		return UserResponse{
			StatusCode: 400,
			Body:       apperrors.ErrEmailTaken,
		}, apperrors.ErrEmailTaken
	} else if err != apperrors.ErrUserNotFound { // Если другая ошибка (например, БД)
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
			Body:       apperrors.ErrPhoneTaken,
		}, apperrors.ErrPhoneTaken
	} else if err != apperrors.ErrUserNotFound {
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
			Body:       apperrors.ErrUserCreation,
		}, apperrors.ErrUserCreation
	}

	// Если все прошло успешно, возвращаем успешный ответ
	return UserResponse{
		StatusCode: 201,
		Body:       user,
	}, nil
}
func AuthenticateUser(values models.AuthCredentials, session models.Session) (UserResponse, error) {
	user, err := repository.GetUserByUsername(values.Username)
	if err != nil {
		return UserResponse{
			StatusCode: 400,
			Body:       apperrors.ErrUsernameTaken,
		}, apperrors.ErrUserNotFound
	}
	//TODO делать хеширование пароля и сверять хеееееееееееееш
	if user.Password != values.Password {
		return UserResponse{
			StatusCode: 400,
			Body:       apperrors.ErrInvalidCredentials,
		}, apperrors.ErrInvalidCredentials
	}
	err = repository.CreateSession(values.Username, session)
	if err != nil {
		return UserResponse{
			StatusCode: 500,
			Body:       err,
		}, err
	}

	return UserResponse{
		StatusCode: 201,
		Body:       session,
	}, nil
}
