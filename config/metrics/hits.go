package metrics

import (
	"context"
	"net/http"

	"google.golang.org/grpc"
)

func HitCounterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HitCounter.WithLabelValues(r.URL.Path, r.Method).Inc()

		next.ServeHTTP(w, r)
	})
}

func GrpcHitCounterInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		GrpcHitCounter.WithLabelValues(info.FullMethod).Inc()
		return handler(ctx, req)
	}
}
