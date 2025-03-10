package repository

import (
	"sync"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/model"
	"github.com/google/uuid"
)

var (
	sessions  = make(map[string]model.Session)
	muSession sync.Mutex
)

func GetSessionBySessId(sessId string) (model.Session, error) {
	muSession.Lock()
	session, exists := sessions[sessId]
	muSession.Unlock()

	if !exists {
		return model.Session{}, apperrors.ErrSessionNotFound
	}
	return session, nil
}

func CreateSession(username string) (string, error) {
	sessionId := uuid.NewString()
	muSession.Lock()
	sessions[sessionId] = model.Session{
		Username: username,
		Expiry:   time.Now().Add(3 * time.Hour),
	}
	muSession.Unlock()
	return sessionId, nil
}

func DeleteSession(sessionId string) error {
	muSession.Lock()
	defer muSession.Unlock()

	if _, exists := sessions[sessionId]; !exists {
		return apperrors.ErrSessionNotFound
	}
	delete(sessions, sessionId)
	return nil
}
