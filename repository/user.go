package repository

import (
	"sync"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/utils"
)

var userPassword, _ = utils.HashAndSalt("qwerty")

var (
	users = []*model.User{
		{
			FirstName: "Ruslan",
			LastName:  "Novikov",
			Username:  "ruslantus228",
			Phone:     "+79128234765",
			Email:     "rumail@mail.ru",
			Password:  userPassword,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			FirstName: "Ilya",
			LastName:  "Zeonov",
			Username:  "ilyaaaaaaaaz",
			Phone:     "+79476781543",
			Email:     "zeonzeonych@mail.ru",
			Password:  userPassword,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			FirstName: "Aleksey",
			LastName:  "Lupenkov",
			Username:  "lumpaumpenkov",
			Phone:     "+77777777777",
			Email:     "seniorjunior@mail.ru",
			Password:  userPassword,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	muUser sync.RWMutex // Используем RWMutex для безопасного чтения и записи
)

// Получение пользователя по имени пользователя (с безопасностью для чтения)
func GetUserByUsername(username string) (*model.User, error) {
	muUser.RLock() // Чтение - блокируем только для чтения
	defer muUser.RUnlock()

	for _, user := range users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, apperrors.ErrUserNotFound
}

// Получение пользователя по email (с безопасностью для чтения)
func GetUserByEmail(email string) (*model.User, error) {
	muUser.RLock() // Чтение - блокируем только для чтения
	defer muUser.RUnlock()

	for _, user := range users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, apperrors.ErrUserNotFound
}

// Получение пользователя по телефону (с безопасностью для чтения)
func GetUserByPhone(phone string) (*model.User, error) {
	muUser.RLock() // Чтение - блокируем только для чтения
	defer muUser.RUnlock()

	for _, user := range users {
		if user.Phone == phone {
			return user, nil
		}
	}
	return nil, apperrors.ErrUserNotFound
}

// Создание нового пользователя (с безопасностью для записи)
func CreateUser(user *model.User) error {
	muUser.Lock() // Запись - блокируем для записи
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
