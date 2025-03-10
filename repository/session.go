package repository

import (
	"sync"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/model"
	"github.com/google/uuid"
)

var (
	sessions  = make(map[string]*model.Session) // Используем указатели на модель сессии
	muSession sync.Mutex
)

// Получение сессии по ее ID
func GetSessionBySessId(sessId string) (*model.Session, error) {
	muSession.Lock()
	session, exists := sessions[sessId]
	muSession.Unlock()

	if !exists {
		return nil, apperrors.ErrSessionNotFound
	}
	return session, nil
}

// Создание новой сессии
func CreateSession(username string) (string, error) {
	sessionId := uuid.NewString()
	muSession.Lock()
	sessions[sessionId] = &model.Session{ // Сохраняем указатель на сессию
		Username: username,
		Expiry:   time.Now().Add(config.CookieDuration),
	}
	muSession.Unlock()
	return sessionId, nil
}

// Удаление сессии
func DeleteSession(sessionId string) error {
	muSession.Lock()
	defer muSession.Unlock()

	if _, exists := sessions[sessionId]; !exists {
		return apperrors.ErrSessionNotFound
	}
	delete(sessions, sessionId)
	return nil
}
