package repository

import (
	"sync"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/model"
)

var (
	users = []*model.User{
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
	muUser sync.Mutex
)

func GetUserByUsername(username string) (*model.User, error) {
	muUser.Lock()
	defer muUser.Unlock()

	for _, user := range users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, apperrors.ErrUserNotFound
}

func GetUserByEmail(email string) (*model.User, error) {
	muUser.Lock()
	defer muUser.Unlock()

	for _, user := range users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, apperrors.ErrUserNotFound
}

func GetUserByPhone(phone string) (*model.User, error) {
	muUser.Lock()
	defer muUser.Unlock()

	for _, user := range users {
		if user.Phone == phone {
			return user, nil
		}
	}
	return nil, apperrors.ErrUserNotFound
}

func CreateUser(user *model.User) error {
	muUser.Lock()
	defer muUser.Unlock()

	for _, u := range users {
		if u.Username == user.Username {
			return apperrors.ErrUsernameTaken
		}
		if u.Phone == user.Phone {
			return apperrors.ErrPhoneTaken
		}
		if u.Email == user.Email {
			return apperrors.ErrEmailTaken
		}
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	users = append(users, user)
	return nil
}
