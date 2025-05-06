package metrics

import (
	"context"
	"net/http"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func ErrorCounterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := NewResponseWriter(w)

		next.ServeHTTP(rw, r)

		status := rw.statusCode
		ErrorCounter.WithLabelValues(r.URL.Path, r.Method, strconv.Itoa(status)).Inc()
	})
}

func GrpcErrorCounterInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		st, _ := status.FromError(err)

		GrpcErrorCounter.WithLabelValues(
			info.FullMethod,
			strconv.Itoa(int(st.Code())),
		).Inc()

		if err != nil {
			ErrorCounter.WithLabelValues(
				info.FullMethod,
				"grpc",
				strconv.Itoa(int(st.Code())),
			).Inc()
		}

		return resp, err
	}
}
