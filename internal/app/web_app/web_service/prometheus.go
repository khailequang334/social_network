package web_service

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var latencyExporter = promauto.NewSummaryVec(
	prometheus.SummaryOpts{
		Name:       "web_request_latency",
		Help:       "Request latency in milliseconds",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"component", "status"},
)

var countExporter = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "web_request_count",
		Help: "Requests count",
	},
	[]string{"component", "type"},
)
