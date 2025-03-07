package repository

import (
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/models"
)

var sessions = make(map[int64]models.Session)

func GetSessionByID(id string) (models.Session, error) {
	for _, session := range sessions {
		if session.ID == id {
			return session, nil
		}
	}
	return models.Session{}, apperrors.ErrSessionNotFound
}

func CreateSession(sessionId string, session models.Session) error {
	if _, exists := GetSessionByID(sessionId); exists == nil {
		return apperrors.ErrSessionAlreadyExists
	}

	session.ID = sessionId
	session.Expiry = time.Now() //TODO поменять на другое время
	return nil
}
