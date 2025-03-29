package middleware

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
)

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		ctx := context.WithValue(r.Context(), utils.REQUEST_ID_KEY, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
