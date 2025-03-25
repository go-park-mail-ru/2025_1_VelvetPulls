package utils

import (
	"context"
)

type contextKey string

const (
	SESSION_ID_KEY = contextKey("sessionID")
	USER_ID_KEY    = contextKey("userID")
)

func GetSessionIDFromCtx(ctx context.Context) string {
	return ctx.Value(SESSION_ID_KEY).(string)
}

func GetUserIDFromCtx(ctx context.Context) string {
	return ctx.Value(USER_ID_KEY).(string)
}
