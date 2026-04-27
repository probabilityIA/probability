package errors

import "errors"

var (
	ErrPaymentNotFound         = errors.New("payment transaction not found")
	ErrInvalidAmount           = errors.New("amount must be greater than 0")
	ErrInvalidGateway          = errors.New("unsupported payment gateway")
	ErrPaymentAlreadyProcessed = errors.New("payment transaction already processed")
	ErrMaxRetriesReached       = errors.New("maximum retry attempts reached")

	ErrWalletNotFound         = errors.New("wallet not found")
	ErrTransactionNotFound    = errors.New("wallet transaction not found")
	ErrTransactionNotPending  = errors.New("transaction is not in pending status")
	ErrMinimumRechargeAmount  = errors.New("amount is below the minimum recharge amount")
	ErrInsufficientBalance    = errors.New("insufficient wallet balance")

	ErrBoldConfigNotFound      = errors.New("bold integration type not configured")
	ErrBoldCredentialsMissing  = errors.New("bold credentials are not configured")
	ErrBoldOrderNotFound       = errors.New("bold order not found")
	ErrBoldUnauthorized        = errors.New("bold credentials are invalid or unauthorized")
	ErrBoldUpstreamUnavailable = errors.New("bold api unavailable")
)
