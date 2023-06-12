package nats

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	totalSubscribeMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "nats_discount_incoming_messages_total",
		Help: "The total number of incoming report NATS messages",
	})
	successSubscribeMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "nats_discount_success_incoming_messages_total",
		Help: "The total number of success report NATS messages",
	})
	errorSubscribeMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "nats_discount_error_incoming_messages_total",
		Help: "The total number of error report NATS messages",
	})
)
