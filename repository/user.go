package repository

import (
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/models"
)

var users = map[int64]models.User{
	1: {
		ID:        1,
		FirstName: "Ruslan",
		LastName:  "Novikov",
		Username:  "ruslantus228",
		Phone:     "+79128234765",
		Email:     "rumail@mail.ru",
		Password:  "qwerty",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	2: {
		ID:        2,
		FirstName: "Ilya",
		LastName:  "Zeonov",
		Username:  "ilyaaaaaaaaz",
		Phone:     "+79476781543",
		Email:     "zeonzeonych@mail.ru",
		Password:  "qwerty",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	3: {
		ID:        3,
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

func GetUserByUsername(username string) (models.User, error) {
	for _, user := range users {
		if user.Username == username {
			return user, nil
		}
	}
	return models.User{}, apperrors.ErrUserNotFound
}

func GetUserByEmail(email string) (models.User, error) {
	for _, user := range users {
		if user.Email == email {
			return user, nil
		}
	}
	return models.User{}, apperrors.ErrUserNotFound
}

func GetUserByPhone(phone string) (models.User, error) {
	for _, user := range users {
		if user.Phone == phone {
			return user, nil
		}
	}
	return models.User{}, apperrors.ErrUserNotFound
}

func CreateUser(user models.User) error {
	if _, exists := users[user.ID]; exists {
		return apperrors.ErrUserAlreadyExists
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	users[user.ID] = user
	return nil
}

// TODO: подумать как реализовать изменение и удаление пользователя.
