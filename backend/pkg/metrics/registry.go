package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration, httpResponseSize)
}

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests processed, partitioned by status code and method.",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of latencies for HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	httpResponseSize = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_response_size_bytes",
			Help: "Size of HTTP responses.",
		},
		[]string{"method", "path"},
	)

	carRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "car_requests",
			Help: "Total number of requests for car reservations",
		},
		[]string{"car_id"},
	)

	businessRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "business_requests",
			Help: "Total number of requests for business reservations",
		},
		[]string{"business_id"},
	)
)

func IncHttpRequestsTotal(method, path, status string) {
	httpRequestsTotal.WithLabelValues(method, path, status).Inc()
}

func SetHttpRequestDuration(method, path string, duration float64) {
	httpRequestDuration.WithLabelValues(method, path).Observe(duration)
}

func SetHttpResponseSize(method, path string, size float64) {
	httpResponseSize.WithLabelValues(method, path).Observe(size)
}

func IncCarRequestsTotal(car_id string) {
	carRequests.WithLabelValues(car_id).Inc()
}

func IncBusinessRequestsTotal(business_id string) {
	businessRequests.WithLabelValues(business_id).Inc()
}
