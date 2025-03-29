package repository

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type ISessionRepo interface {
	GetUserIDByToken(ctx context.Context, sessId string) (string, error)
	CreateSession(ctx context.Context, userID string) (string, error)
	DeleteSession(ctx context.Context, sessionId string) error
}

type sessionRepo struct {
	redisClient *redis.Client
}

func NewSessionRepo(redisClient *redis.Client) ISessionRepo {
	return &sessionRepo{redisClient: redisClient}
}

func (r *sessionRepo) GetUserIDByToken(ctx context.Context, sessId string) (string, error) {
	userID, err := r.redisClient.Get(ctx, sessId).Result()
	if err == redis.Nil {
		return "", apperrors.ErrSessionNotFound
	} else if err != nil {
		return "", err
	}
	return userID, nil
}

func (r *sessionRepo) CreateSession(ctx context.Context, userID string) (string, error) {
	sessionId := uuid.NewString()

	err := r.redisClient.Set(ctx, sessionId, userID, config.CookieDuration).Err()
	if err != nil {
		return "", err
	}

	return sessionId, nil
}

func (r *sessionRepo) DeleteSession(ctx context.Context, sessionId string) error {
	exists, err := r.redisClient.Exists(ctx, sessionId).Result()
	if err != nil {
		return err
	}
	if exists == 0 {
		return errors.New("session not found")
	}
	return r.redisClient.Del(ctx, sessionId).Err()
}
