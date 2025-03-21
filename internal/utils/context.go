package utils

import (
	"context"

	"go.uber.org/zap"
)

type contextKey string

const (
	LOGGER_ID_KEY  = contextKey("loggerID")
	REQUEST_ID_KEY = contextKey("requestID")
	SESSION_ID_KEY = contextKey("sessionID")
)

func GetRequestIDFromCtx(ctx context.Context) string {
	return ctx.Value(REQUEST_ID_KEY).(string)
}

func GetSessionIDFromCtx(ctx context.Context) string {
	return ctx.Value(SESSION_ID_KEY).(string)
}

func GetContextLogger(ctx context.Context) *zap.SugaredLogger {
	return ctx.Value(LOGGER_ID_KEY).(*zap.SugaredLogger)
}
