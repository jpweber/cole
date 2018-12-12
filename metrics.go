package main

import "github.com/prometheus/client_golang/prometheus"

var (
	timerCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "cole",
		Subsystem: "server",
		Name:      "timer_count",
		Help:      "Count of active timers in the system",
	})

	httpDurations = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace:  "cole",
			Name:       "http_durations_seconds",
			Help:       "HTTP latency distributions.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
	)

	dmAlertsRecieved = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "cole",
			Name:      "deadman_alerts_recvd",
			Help:      "Number of times a DeadManSwitch has sent cole and alert.",
		},
	)
)

func init() {
	prometheus.MustRegister(timerCount)
	prometheus.MustRegister(httpDurations)
	prometheus.MustRegister(dmAlertsRecieved)
}
