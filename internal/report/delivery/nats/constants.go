package nats

import "time"

const (
	ackWait     = 60 * time.Second
	durableName = "report-dur"
	maxInflight = 25

	createReportWorkers = 6
	chargeReportWorkers = 6

	createReportSubject = "report:create"
	reportGroupName     = "report_service"

	deadLetterQueueSubject = "report:errors"
	maxRedeliveryCount     = 3
)
