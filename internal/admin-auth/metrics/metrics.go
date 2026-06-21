package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metrics struct {
	requestCounter        *prometheus.CounterVec
	histogramResponseTime *prometheus.HistogramVec
}

func New(namespace, appName string, reg prometheus.Registerer) *metrics {
	wrappedReg := prometheus.WrapRegistererWith(prometheus.Labels{"app": appName}, reg)

	wrappedReg.MustRegister(collectors.NewGoCollector())
	wrappedReg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	return &metrics{
		requestCounter: promauto.With(wrappedReg).NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Total number of HTTP requests",
		}, []string{"path", "method", "status"}),
		histogramResponseTime: promauto.With(wrappedReg).NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "http",
			Name:      "response_time_seconds",
			Help:      "Response time in seconds",
			Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
		}, []string{"path", "method", "status"}),
	}
}

func (m *metrics) IncRequestCounter(path, method, status string) {
	m.requestCounter.WithLabelValues(path, method, status).Inc()
}

func (m *metrics) HistogramResponseTimeObserve(path, method, status string, time float64) {
	m.histogramResponseTime.WithLabelValues(path, method, status).Observe(time)
}
