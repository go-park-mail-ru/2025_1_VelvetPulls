package service

import (
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/errors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/model"
)

type UserRepository interface {
	GetUserByID(ID int64) (model.User, error)
	CreateUser(user model.User) error
	// TODO: расширить функционал
}

type SessionRepository interface {
	// TODO: добавить методы сессий
}

type UserService struct {
	userRepo    UserRepository
	sessionRepo SessionRepository
}

func NewUserService(userRepo UserRepository, sessionRepo SessionRepository) *UserService {
	return &UserService{userRepo: userRepo, sessionRepo: sessionRepo}
}

func (s *UserService) RegisterUser(user model.User) (UserResponse, error) {
	// Проверяем, существует ли уже такой пользователь
	existingUser, err := s.userRepo.GetUserByID(user.ID)
	if err == nil && existingUser.ID != 0 {
		return UserResponse{
			StatusCode: 400, // Bad Request
			Body:       err,
		}, err
	}

	// Создаем нового пользователя
	err = s.userRepo.CreateUser(user)
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
