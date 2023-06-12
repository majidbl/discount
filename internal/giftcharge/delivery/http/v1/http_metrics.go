package v1

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	successRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_giftcharge_success_incoming_messages_total",
		Help: "The total number of success incoming giftcharge HTTP requests",
	})
	errorRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_giftcharge_error_incoming_message_total",
		Help: "The total number of error incoming giftcharge HTTP requests",
	})
	createRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_giftcharge_create_incoming_requests_total",
		Help: "The total number of incoming create giftcharge HTTP requests",
	})
	getByIdRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_giftcharge_get_by_id_incoming_requests_total",
		Help: "The total number of incoming get by id giftcharge HTTP requests",
	})
)
