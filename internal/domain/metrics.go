package domain

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Status string

const (
	Success   Status = "success"
	Failed    Status = "failed"
	Skipped   Status = "skipped"
	Processed Status = "skipped"
)

type Type string

const (
	Command  Type = "command"
	Text     Type = "text"
	Callback Type = "callback"
	Skip     Type = "skip"
)

var requestMetrics = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Namespace:  "bot",
	Subsystem:  "tg",
	Name:       "request",
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
}, []string{"type", "text", "status"})

func observeRequest(d time.Duration, requestType Type, text string, status Status) {
	requestMetrics.WithLabelValues(string(requestType), string(text), string(status)).Observe(d.Seconds())
}
