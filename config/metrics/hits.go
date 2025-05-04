package metrics

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

var HitCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "requests_amount_count",
		Help: "Accumulates incoming requests",
	},
	[]string{"path", "method"},
)

func HitCounterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HitCounter.WithLabelValues(r.URL.Path, r.Method).Inc()

		next.ServeHTTP(w, r)
	})
}

var GrpcHitCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "grpc_requests_amount_count",
		Help: "Accumulates incoming gRPC requests",
	},
	[]string{"method"}, // или другие подходящие метки
)

func GrpcHitCounterInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		GrpcHitCounter.WithLabelValues(info.FullMethod).Inc()
		return handler(ctx, req)
	}
}
