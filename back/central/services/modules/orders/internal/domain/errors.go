package domain

import "errors"

var (
	// ErrOrderAlreadyExists indicates that an order with the same external ID already exists for the integration
	ErrOrderAlreadyExists = errors.New("order with this external_id already exists for this integration")
)
