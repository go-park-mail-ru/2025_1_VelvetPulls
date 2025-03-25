package utils

import (
	"context"
)

type contextKey string

const (
	USER_ID_KEY = contextKey("userID")
)

func GetUserIDFromCtx(ctx context.Context) string {
	return ctx.Value(USER_ID_KEY).(string)
}
