package utils

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type contextKey string

const (
	LOGGER_ID_KEY  = contextKey("loggerID")
	REQUEST_ID_KEY = contextKey("requestID")
	USER_ID_KEY    = contextKey("userID")
)

func GetRequestIDFromCtx(ctx context.Context) string {
	return ctx.Value(REQUEST_ID_KEY).(string)
}

func GetLoggerFromCtx(ctx context.Context) *zap.Logger {
	return ctx.Value(LOGGER_ID_KEY).(*zap.Logger)
}

func GetUserIDFromCtx(ctx context.Context) uuid.UUID {
	return ctx.Value(USER_ID_KEY).(uuid.UUID)
}
