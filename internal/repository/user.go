package repository

import (
	"sync"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
)

type UserRepoInterface interface {
	GetUserByUsername(username string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	GetUserByPhone(phone string) (*model.User, error)
	CreateUser(user *model.User) error
}

type userRepo struct {
	users []*model.User
	mu    sync.RWMutex // Используем RWMutex для безопасного чтения и записи
}

func NewUserRepo() UserRepoInterface {
	return &userRepo{
		users: make([]*model.User, 0),
	}
}

// Получение пользователя по имени пользователя (с безопасностью для чтения)
func (r *userRepo) GetUserByUsername(username string) (*model.User, error) {
	r.mu.RLock() // Чтение - блокируем только для чтения
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, apperrors.ErrUserNotFound
}

// Получение пользователя по email (с безопасностью для чтения)
func (r *userRepo) GetUserByEmail(email string) (*model.User, error) {
	r.mu.RLock() // Чтение - блокируем только для чтения
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, apperrors.ErrUserNotFound
}

// Получение пользователя по телефону (с безопасностью для чтения)
func (r *userRepo) GetUserByPhone(phone string) (*model.User, error) {
	r.mu.RLock() // Чтение - блокируем только для чтения
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Phone == phone {
			return user, nil
		}
	}
	return nil, apperrors.ErrUserNotFound
}

// Создание нового пользователя (с безопасностью для записи)
func (r *userRepo) CreateUser(user *model.User) error {
	r.mu.Lock() // Запись - блокируем для записи
	defer r.mu.Unlock()

	for _, u := range r.users {
		if u.Username == user.Username {
			return apperrors.ErrUsernameTaken
		}
		if u.Phone == user.Phone {
			return apperrors.ErrPhoneTaken
		}
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	r.users = append(r.users, user)
	return nil
}
