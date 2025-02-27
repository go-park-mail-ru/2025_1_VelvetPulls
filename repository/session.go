package repository

import (
	"time"

	errors "github.com/go-park-mail-ru/2025_1_VelvetPulls/errors"
	model "github.com/go-park-mail-ru/2025_1_VelvetPulls/model"
)

var sessions = make(map[int64]model.Session)

func GetSessionByID(id string) (model.Session, error) {
	for _, session := range sessions {
		if session.ID == id {
			return session, nil
		}
	}
	return model.Session{}, errors.ErrSessionNotFound
}

func CreateSession(sessionId string, session model.Session) error {
	if _, exists := GetSessionByID(sessionId); exists == nil {
		return errors.ErrSessionAlreadyExists
	}

	session.ID = sessionId
	session.Expiry = time.Now() //TODO поменять на другое время
	return nil
}
