package nats

import "time"

const (
	ackWait     = 60 * time.Second
	durableName = "Discount-dur"
	maxInflight = 25

	createDiscountWorkers = 6
	chargeDiscountWorkers = 6

	createDiscountSubject = "discount:create"
	chargeDiscountSubject = "discount:charge"
	DiscountGroupName     = "discount_service"

	deadLetterQueueSubject = "Discount:errors"
	maxRedeliveryCount     = 3
)
