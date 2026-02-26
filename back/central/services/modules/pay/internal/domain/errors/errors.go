package errors

import "errors"

var (
	ErrPaymentNotFound         = errors.New("payment transaction not found")
	ErrInvalidAmount           = errors.New("amount must be greater than 0")
	ErrInvalidGateway          = errors.New("unsupported payment gateway")
	ErrPaymentAlreadyProcessed = errors.New("payment transaction already processed")
	ErrMaxRetriesReached       = errors.New("maximum retry attempts reached")
)
