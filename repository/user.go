package repository

import (
	"time"

	err "github.com/go-park-mail-ru/2025_1_VelvetPulls/errors"
	model "github.com/go-park-mail-ru/2025_1_VelvetPulls/model"
)

type UserRepository struct {
	users map[int64]model.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[int64]model.User),
	}
}

func (r *UserRepository) GetUserByID(ID int64) (model.User, error) {
	user, exists := r.users[ID]
	if !exists {
		return model.User{}, err.ErrUserNotFound
	}
	return user, nil
}

func (r *UserRepository) CreateUser(user model.User) error {
	if _, exists := r.users[user.ID]; exists {
		return err.ErrUserAlreadyExists
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	r.users[user.ID] = user
	return nil
}

// TODO: подумать как реализовать изменение и удаление пользователя.
