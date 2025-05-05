package metrics

func IncBusinessOp(operation string) {
	BusinessOpsCounter.WithLabelValues(operation).Inc()
}
