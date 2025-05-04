package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var BusinessOpsCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "business_operations_total",
		Help: "Count of business operations",
	},
	[]string{"operation"},
)

func init() {
	prometheus.MustRegister(BusinessOpsCounter)
}

func IncBusinessOp(operation string) {
	BusinessOpsCounter.WithLabelValues(operation).Inc()
}
