package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"go.uber.org/zap"
)

type responseWriter struct {
	http.ResponseWriter
	status int
	body   string
}

func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrappedWriter := &responseWriter{ResponseWriter: w}

		requestID := utils.GetRequestIDFromCtx(r.Context())

		contextLogger := utils.Logger.With(
			zap.String("request_id", requestID),
			zap.String("method", r.Method),
			zap.String("url", r.URL.Path),
		)

		ctx := context.WithValue(r.Context(), utils.LOGGER_ID_KEY, contextLogger)
		next.ServeHTTP(wrappedWriter, r.WithContext(ctx))

		contextLogger.Info("HTTP Request",
			zap.Int("status", wrappedWriter.status),
			zap.Duration("execution_time", time.Since(start)),
		)

		if wrappedWriter.body != "" {
			contextLogger.Error(wrappedWriter.body)
		} else {
			contextLogger.Info("HTTP Request completed successfully")
		}
	})
}
