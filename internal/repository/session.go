package repository

import (
	"sync"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/google/uuid"
)

type SessionRepoInterface interface {
	GetSessionBySessId(sessId string) (*model.Session, error)
	CreateSession(username string) (string, error)
	DeleteSession(sessionId string) error
}

type sessionRepo struct {
	sessions map[string]*model.Session
	mu       sync.RWMutex // Мьютекс для безопасного чтения и записи
}

func NewSessionRepo() SessionRepoInterface {
	return &sessionRepo{
		sessions: make(map[string]*model.Session),
	}
}

// Получение сессии по ее ID
func (r *sessionRepo) GetSessionBySessId(sessId string) (*model.Session, error) {
	r.mu.RLock() // Блокировка для чтения
	defer r.mu.RUnlock()

	session, exists := r.sessions[sessId]
	if !exists {
		return nil, apperrors.ErrSessionNotFound
	}

	// Возвращаем копию сессии, чтобы избежать гонок данных
	sessionCopy := *session
	return &sessionCopy, nil
}

// Создание новой сессии
func (r *sessionRepo) CreateSession(username string) (string, error) {
	r.mu.Lock() // Блокировка для записи
	defer r.mu.Unlock()

	sessionId := uuid.NewString()

	// Проверяем, существует ли уже сессия с таким ID (хотя это маловероятно для UUID)
	if _, exists := r.sessions[sessionId]; exists {
		return "", apperrors.ErrSessionAlreadyExists
	}

	r.sessions[sessionId] = &model.Session{
		Username: username,
		Expiry:   time.Now().Add(config.CookieDuration),
	}
	return sessionId, nil
}

// Удаление сессии
func (r *sessionRepo) DeleteSession(sessionId string) error {
	r.mu.Lock() // Блокировка для записи
	defer r.mu.Unlock()

	if _, exists := r.sessions[sessionId]; !exists {
		return apperrors.ErrSessionNotFound
	}
	delete(r.sessions, sessionId)
	return nil
}
