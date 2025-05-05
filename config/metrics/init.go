package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	initOnce sync.Once

	BusinessOpsCounter *prometheus.CounterVec

	ErrorCounter     *prometheus.CounterVec
	GrpcErrorCounter *prometheus.CounterVec

	HitCounter     *prometheus.CounterVec
	GrpcHitCounter *prometheus.CounterVec

	HttpDuration *prometheus.HistogramVec
	GrpcDuration *prometheus.HistogramVec
)

func init() {
	initOnce.Do(func() {
		// Бизнес-метрики
		BusinessOpsCounter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "business_operations_total",
				Help: "Count of business operations",
			},
			[]string{"operation"},
		)

		// Метрики ошибок
		ErrorCounter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "errors_amount_count",
				Help: "Accumulates outgoing errors",
			},
			[]string{"path", "method", "status"},
		)

		GrpcErrorCounter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "grpc_errors_amount_count",
				Help: "Accumulates outgoing gRPC errors",
			},
			[]string{"method", "status"},
		)

		// Метрики запросов
		HitCounter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "requests_amount_count",
				Help: "Accumulates incoming requests",
			},
			[]string{"path", "method"},
		)

		GrpcHitCounter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "grpc_requests_amount_count",
				Help: "Accumulates incoming gRPC requests",
			},
			[]string{"method"},
		)

		// Метрики времени выполнения
		HttpDuration = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_response_time_seconds",
				Help:    "Duration of HTTP requests.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"path", "method"},
		)

		GrpcDuration = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "grpc_response_time_seconds",
				Help:    "Duration of gRPC requests.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method"},
		)

		// Регистрация всех метрик
		prometheus.MustRegister(
			BusinessOpsCounter,
			ErrorCounter,
			GrpcErrorCounter,
			HitCounter,
			GrpcHitCounter,
			HttpDuration,
			GrpcDuration,
		)
	})
}
