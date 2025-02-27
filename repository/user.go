package repository

import (
	"time"

	errors "github.com/go-park-mail-ru/2025_1_VelvetPulls/errors"
	model "github.com/go-park-mail-ru/2025_1_VelvetPulls/model"
)

var users = make(map[int64]model.User)

func GetUserByUsername(username string) (model.User, error) {
	for _, user := range users {
		if user.Username == username {
			return user, nil
		}
	}
	return model.User{}, errors.ErrUserNotFound
}
func GetUserByEmail(email string) (model.User, error) {
	for _, user := range users {
		if user.Email == email {
			return user, nil
		}
	}
	return model.User{}, errors.ErrUserNotFound
}

func GetUserByPhone(phone string) (model.User, error) {
	for _, user := range users {
		if user.Phone == phone {
			return user, nil
		}
	}
	return model.User{}, errors.ErrUserNotFound
}

func CreateUser(user model.User) error {
	if _, exists := users[user.ID]; exists {
		return errors.ErrUserAlreadyExists
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	users[user.ID] = user
	return nil
}

// TODO: подумать как реализовать изменение и удаление пользователя.
