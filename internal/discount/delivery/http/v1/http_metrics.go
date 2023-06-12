package v1

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	successRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_discount_success_incoming_messages_total",
		Help: "The total number of success incoming report HTTP requests",
	})
	errorRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_discount_error_incoming_message_total",
		Help: "The total number of error incoming report HTTP requests",
	})
	createRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_discount_create_incoming_requests_total",
		Help: "The total number of incoming create report HTTP requests",
	})
)
