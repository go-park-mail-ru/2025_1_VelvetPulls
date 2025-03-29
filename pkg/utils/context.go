package utils

import (
	"context"

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

func GetContextLogger(ctx context.Context) *zap.SugaredLogger {
	return ctx.Value(LOGGER_ID_KEY).(*zap.SugaredLogger)
}

func GetUserIDFromCtx(ctx context.Context) string {
	return ctx.Value(USER_ID_KEY).(string)
}
