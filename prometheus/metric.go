package prometheus

import "github.com/prometheus/client_golang/prometheus"

var (
	requestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_requests_total",
		Help: "Count number of requests.",
	}, []string{"type", "service", "method", "status"})
	requestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "grpc_request_duration_seconds",
		Help:    "Hisogram for the runtime of request.",
		Buckets: prometheus.DefBuckets,
	}, []string{"type", "service", "method"})
)

func init() {
	registerMetrics()
}

func registerMetrics() {
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(requestDuration)
}
