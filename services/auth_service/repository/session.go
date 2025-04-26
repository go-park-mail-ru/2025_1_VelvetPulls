package repository

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type ISessionRepo interface {
	GetUserIDByToken(ctx context.Context, sessionID string) (string, error)
	CreateSession(ctx context.Context, userID uuid.UUID) (string, error)
	DeleteSession(ctx context.Context, sessionID string) error
}

type sessionRepo struct {
	redisClient *redis.Client
}

func NewSessionRepo(redisClient *redis.Client) ISessionRepo {
	return &sessionRepo{redisClient: redisClient}
}

func (r *sessionRepo) GetUserIDByToken(ctx context.Context, sessionID string) (string, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("Getting user ID by session token")

	userID, err := r.redisClient.Get(ctx, sessionID).Result()
	if err == redis.Nil {
		logger.Error("Session not found")
		return "", ErrSessionNotFound
	} else if err != nil {
		logger.Error("Error during Redis operation")
		return "", ErrDatabaseOperation
	}
	return userID, nil
}
func (r *sessionRepo) CreateSession(ctx context.Context, userID uuid.UUID) (string, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("Creating new session")

	sessionId := uuid.NewString()

	err := r.redisClient.Set(ctx, sessionId, userID.String(), config.CookieDuration).Err()
	if err != nil {
		logger.Error("Error creating session in Redis")
		return "", ErrDatabaseOperation
	}

	return sessionId, nil
}

func (r *sessionRepo) DeleteSession(ctx context.Context, sessionID string) error {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("Deleting session")

	exists, err := r.redisClient.Exists(ctx, sessionID).Result()
	if err != nil {
		logger.Error("Error checking if session exists")
		return ErrDatabaseOperation
	}

	if exists == 0 {
		logger.Error("Session not found")
		return ErrSessionNotFound
	}

	err = r.redisClient.Del(ctx, sessionID).Err()
	if err != nil {
		logger.Error("Error deleting session in Redis")
		return ErrDatabaseOperation
	}

	return nil
}
