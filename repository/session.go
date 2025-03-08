package repository

import (
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/model"
	"github.com/google/uuid"
)

var sessions = make(map[string]model.Session)

func GetSessionBySessId(sessId string) (model.Session, error) {
	session, exists := sessions[sessId]
	if !exists {
		return model.Session{}, apperrors.ErrSessionNotFound
	}

	return session, nil
}

func CreateSession(username string) (string, error) {
	sessionId := uuid.NewString()
	sessions[uuid.NewString()] = model.Session{
		Username: username,
		Expiry:   time.Now().Add(3 * time.Hour), // Сессия истекает через 3 часа
	}
	// Может быть ошибка, если читать из redis
	return sessionId, nil
}

func DeleteSession(sessionID string) error {
	if _, exists := sessions[sessionID]; !exists {
		return apperrors.ErrSessionNotFound
	}

	delete(sessions, sessionID)
	return nil
}
