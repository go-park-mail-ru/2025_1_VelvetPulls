package middleware

import (
	"bufio"
	"context"
	"fmt"
	"net"
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

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(p []byte) (int, error) {
	rw.body += string(p)
	return rw.ResponseWriter.Write(p)
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("wrapped ResponseWriter does not implement http.Hijacker")
	}
	return hj.Hijack()
}

func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrappedWriter := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		requestID := utils.GetRequestIDFromCtx(r.Context())
		contextLogger := utils.Logger.With(
			zap.String("request_id", requestID),
			zap.String("method", r.Method),
			zap.String("url", r.URL.Path),
		)

		ctx := context.WithValue(r.Context(), utils.LOGGER_ID_KEY, contextLogger)
		next.ServeHTTP(wrappedWriter, r.WithContext(ctx))

		contentType := w.Header().Get("Content-Type")
		if contentType == "application/json" {
			contextLogger.Info("HTTP Request completed",
				zap.Int("status", wrappedWriter.status),
				zap.Duration("execution_time", time.Since(start)),
				zap.String("body", wrappedWriter.body),
			)
		} else {
			contextLogger.Info("HTTP Request completed",
				zap.Int("status", wrappedWriter.status),
				zap.Duration("execution_time", time.Since(start)),
			)
		}
	})
}
