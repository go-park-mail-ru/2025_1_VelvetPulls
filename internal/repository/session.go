package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type ISessionRepo interface {
	GetSessionBySessId(sessId string) (*model.Session, error)
	CreateSession(username string) (string, error)
	DeleteSession(sessionId string) error
}

type sessionRepo struct {
	redisClient *redis.Client
}

func NewSessionRepo(redisClient *redis.Client) ISessionRepo {
	return &sessionRepo{redisClient: redisClient}
}

func (r *sessionRepo) GetSessionBySessId(sessId string) (*model.Session, error) {
	ctx := context.Background()

	sessionData, err := r.redisClient.Get(ctx, sessId).Bytes()
	if err == redis.Nil {
		return nil, apperrors.ErrSessionNotFound
	} else if err != nil {
		return nil, err
	}

	var session model.Session
	if err := json.Unmarshal(sessionData, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *sessionRepo) CreateSession(username string) (string, error) {
	sessionId := uuid.NewString()

	session := &model.Session{
		Username: username,
		Expiry:   time.Now().Add(config.CookieDuration),
	}

	sessionData, err := json.Marshal(session)
	if err != nil {
		return "", err
	}

	err = r.redisClient.Set(context.Background(), sessionId, sessionData, config.CookieDuration).Err()
	if err != nil {
		return "", err
	}

	return sessionId, nil
}

func (r *sessionRepo) DeleteSession(sessionId string) error {
	ctx := context.Background()
	return r.redisClient.Del(ctx, sessionId).Err()
}
