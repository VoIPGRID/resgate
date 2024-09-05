package metrics

import (
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	uUIDRegex     = regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")
	iDRegex       = regexp.MustCompile("[0-9]+")
	stringIDRegex = regexp.MustCompile(`^(dashboard\.signed\.).*$`)
)

type MetricSet struct {
	// WebSocket connectionws
	WSConnections prometheus.Gauge
	// WebSocket requests
	WSRequests *prometheus.CounterVec
	// Cache
	CacheResources     prometheus.Gauge
	CacheSubscriptions *prometheus.GaugeVec
	// HTTP requests
	HTTPRequests *prometheus.CounterVec
}

func (m *MetricSet) Register(version string, protocolVersion string) {
	// Resgate info
	prometheus.MustRegister(prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "resgate",
		Name:      "info",
		Help:      "Information about resgate.",
		ConstLabels: prometheus.Labels{
			"version":  version,
			"protocol": protocolVersion,
		},
	}))

	// WebSocket connections
	m.WSConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "resgate",
		Subsystem: "ws",
		Name:      "current_connections",
		Help:      "Current established WebSocket connections.",
	})
	prometheus.MustRegister(m.WSConnections)

	// WebSocket requests
	m.WSRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "resgate",
		Subsystem: "ws",
		Name:      "requests",
		Help:      "Total WebSocket client requests.",
	}, []string{"method"})
	prometheus.MustRegister(m.WSRequests)

	// Cache
	m.CacheResources = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "resgate",
		Subsystem: "cache",
		Name:      "resources",
		Help:      "Current number of resources stored in the cache.",
	})
	prometheus.MustRegister(m.CacheResources)
	m.CacheSubscriptions = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "resgate",
		Subsystem: "cache",
		Name:      "subscriptions",
		Help:      "Current number of subscriptions on cached resources.",
	}, []string{"name"})
	prometheus.MustRegister(m.CacheSubscriptions)

	// HTTP requests
	m.HTTPRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "resgate",
		Subsystem: "http",
		Name:      "requests",
		Help:      "Total HTTP client requests.",
	}, []string{"method"})
	prometheus.MustRegister(m.HTTPRequests)
}

func SanitizedString(s string) string {
	s = stringIDRegex.ReplaceAllString(s, "$1{stringid}")
	s = uUIDRegex.ReplaceAllString(s, "{uuid}")
	s = iDRegex.ReplaceAllString(s, "{id}")
	return s
}
