package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metrics struct {
	requestCounter        *prometheus.CounterVec
	histogramResponseTime *prometheus.HistogramVec
	activeRequests        prometheus.Gauge

	registrationsTotal   *prometheus.CounterVec
	loginsTotal          *prometheus.CounterVec
	businessReviewsTotal *prometheus.CounterVec
}

func New(namespace, appName string, reg prometheus.Registerer) *metrics {
	wrappedReg := prometheus.WrapRegistererWith(prometheus.Labels{"app": appName}, reg)

	wrappedReg.MustRegister(collectors.NewGoCollector())
	wrappedReg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	m := &metrics{
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
		activeRequests: promauto.With(wrappedReg).NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "http",
			Name:      "active_requests",
			Help:      "Number of requests currently being processed",
		}),
		registrationsTotal: promauto.With(wrappedReg).NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "auth",
			Name:      "registrations_total",
			Help:      "Total user registrations",
		}, []string{"provider"}),
		loginsTotal: promauto.With(wrappedReg).NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "auth",
			Name:      "logins_total",
			Help:      "Total user logins",
		}, []string{"provider"}),
		businessReviewsTotal: promauto.With(wrappedReg).NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "admin",
			Name:      "business_reviews_total",
			Help:      "Total business reviews processed",
		}, []string{"status"}),
	}

	return m
}

// middleware
func (m *metrics) IncRequestCounter(path, method, status string) {
	m.requestCounter.WithLabelValues(path, method, status).Inc()
}

func (m *metrics) HistogramResponseTimeObserve(path, method, status string, time float64) {
	m.histogramResponseTime.WithLabelValues(path, method, status).Observe(time)
}

func (m *metrics) IncActiveRequests() {
	m.activeRequests.Inc()
}

func (m *metrics) DecActiveRequests() {
	m.activeRequests.Dec()
}

// handler
func (m *metrics) RecordRegistration(provider string) {
	m.registrationsTotal.WithLabelValues(provider).Inc()
}

func (m *metrics) RecordLogin(provider string) {
	m.loginsTotal.WithLabelValues(provider).Inc()
}

func (m *metrics) RecordBusinessReview(status string) {
	m.businessReviewsTotal.WithLabelValues(status).Inc()
}
