package repository

import (
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/model"
)

var users = []model.User{
	{
		FirstName: "Ruslan",
		LastName:  "Novikov",
		Username:  "ruslantus228",
		Phone:     "+79128234765",
		Email:     "rumail@mail.ru",
		Password:  "qwerty",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		FirstName: "Ilya",
		LastName:  "Zeonov",
		Username:  "ilyaaaaaaaaz",
		Phone:     "+79476781543",
		Email:     "zeonzeonych@mail.ru",
		Password:  "qwerty",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		FirstName: "Aleksey",
		LastName:  "Lupenkov",
		Username:  "lumpaumpenkov",
		Phone:     "+77777777777",
		Email:     "seniorjunior@mail.ru",
		Password:  "qwerty",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}

func GetUserByUsername(username string) (model.User, error) {
	for _, user := range users {
		if user.Username == username {
			return user, nil
		}
	}
	return model.User{}, apperrors.ErrUserNotFound
}

func GetUserByEmail(email string) (model.User, error) {
	for _, user := range users {
		if user.Email == email {
			return user, nil
		}
	}
	return model.User{}, apperrors.ErrUserNotFound
}

func GetUserByPhone(phone string) (model.User, error) {
	for _, user := range users {
		if user.Phone == phone {
			return user, nil
		}
	}
	return model.User{}, apperrors.ErrUserNotFound
}

func CreateUser(user model.User) error {
	for _, u := range users {
		if u.Username == user.Username {
			return apperrors.ErrUserAlreadyExists
		}
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	users = append(users, user)
	return nil
}
