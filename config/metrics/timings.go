package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
)

var HttpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path", "method"})

func TimingHistogramMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := prometheus.NewTimer(HttpDuration.WithLabelValues(r.URL.Path, r.Method))

		next.ServeHTTP(w, r)

		timer.ObserveDuration()
	})
}

var GrpcDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "grpc_response_time_seconds",
	Help:    "Duration of gRPC requests.",
	Buckets: prometheus.DefBuckets,
}, []string{"method"})

func GrpcTimingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)

		GrpcDuration.WithLabelValues(info.FullMethod).Observe(time.Since(start).Seconds())

		return resp, err
	}
}
