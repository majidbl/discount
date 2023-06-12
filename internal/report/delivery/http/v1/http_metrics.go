package v1

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	successRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_report_success_incoming_messages_total",
		Help: "The total number of success incoming report HTTP requests",
	})
	errorRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_report_error_incoming_message_total",
		Help: "The total number of error incoming report HTTP requests",
	})
	getByIdRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_report_get_by_id_incoming_requests_total",
		Help: "The total number of incoming get by id report HTTP requests",
	})
)
