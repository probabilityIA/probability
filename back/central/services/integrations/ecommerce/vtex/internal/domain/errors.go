package domain

import "errors"

var (
	ErrIntegrationNotFound  = errors.New("vtex: integration not found")
	ErrIntegrationNotOwned  = errors.New("vtex: la integracion no pertenece al negocio")
	ErrInvalidCredentials   = errors.New("vtex: invalid credentials")
	ErrMissingAppKey        = errors.New("vtex: missing app_key in credentials")
	ErrMissingAppToken      = errors.New("vtex: missing app_token in credentials")
	ErrMissingAccountName   = errors.New("vtex: missing account_name in config")
	ErrOrderNotFound        = errors.New("vtex: order not found")
	ErrRateLimited          = errors.New("vtex: rate limited, too many requests")
	ErrSKUNotFound          = errors.New("vtex: sku not found")
	ErrNoWarehousesMapped   = errors.New("vtex: no warehouses mapped for inventory sync")
	ErrInventorySyncDisabled = errors.New("vtex: inventory sync is disabled for this integration")
	ErrProductNotFound      = errors.New("vtex: product not found")
	ErrForeignHookExists    = errors.New("vtex: la cuenta ya tiene un webhook de otra herramienta registrado")
)
