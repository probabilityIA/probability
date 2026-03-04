package domain

import "errors"

var (
	ErrOrderPayloadNil     = errors.New("order payload is nil")
	ErrBusinessIDMissing   = errors.New("integration has no business_id assigned")
	ErrIntegrationNotFound = errors.New("integration not found")
	ErrPublishFailed       = errors.New("failed to publish order to queue")
	ErrInvalidCredentials  = errors.New("invalid shopify credentials")
	ErrInvalidIntegrationID = errors.New("invalid integration_id")
)
