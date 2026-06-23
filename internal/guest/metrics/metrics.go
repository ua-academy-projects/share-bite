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

	postsCreated              prometheus.Counter
	collectionsCreated        prometheus.Counter
	postLikes                 prometheus.Counter
	followsCreated            prometheus.Counter
	collectionInvitationsSent prometheus.Counter
	postInvitationsSent       prometheus.Counter
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
		postsCreated: promauto.With(wrappedReg).NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "guest",
			Name:      "posts_created_total",
			Help:      "Total posts created",
		}),
		collectionsCreated: promauto.With(wrappedReg).NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "guest",
			Name:      "collections_created_total",
			Help:      "Total collections created",
		}),
		postLikes: promauto.With(wrappedReg).NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "guest",
			Name:      "post_likes_total",
			Help:      "Total post likes",
		}),
		followsCreated: promauto.With(wrappedReg).NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "guest",
			Name:      "follows_created_total",
			Help:      "Total follow relationships created",
		}),
		collectionInvitationsSent: promauto.With(wrappedReg).NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "guest",
			Name:      "collection_invitations_sent_total",
			Help:      "Total collection invitations sent",
		}),
		postInvitationsSent: promauto.With(wrappedReg).NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "guest",
			Name:      "post_invitations_sent_total",
			Help:      "Total post invitations sent",
		}),
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
func (m *metrics) RecordPostCreated() {
	m.postsCreated.Inc()
}

func (m *metrics) RecordCollectionCreated() {
	m.collectionsCreated.Inc()
}

func (m *metrics) RecordPostLike() {
	m.postLikes.Inc()
}

func (m *metrics) RecordFollowCreated() {
	m.followsCreated.Inc()
}

func (m *metrics) RecordCollectionInvitationSent() {
	m.collectionInvitationsSent.Inc()
}

func (m *metrics) RecordPostInvitationSent() {
	m.postInvitationsSent.Inc()
}
