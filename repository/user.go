package repository

import (
	"time"

	err "github.com/go-park-mail-ru/2025_1_VelvetPulls/errors"
	model "github.com/go-park-mail-ru/2025_1_VelvetPulls/model"
)

var users = make(map[int64]model.User)

func GetUserByID(ID int64) (model.User, error) {
	user, exists := users[ID]
	if !exists {
		return model.User{}, err.ErrUserNotFound
	}
	return user, nil
}

func CreateUser(user model.User) error {
	if _, exists := users[user.ID]; exists {
		return err.ErrUserAlreadyExists
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	users[user.ID] = user
	return nil
}

// TODO: подумать как реализовать изменение и удаление пользователя.
