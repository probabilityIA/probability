package domain

import "errors"

var (
	ErrIntegrationNotFound          = errors.New("integracion de Jumpseller no encontrada")
	ErrInvalidCredentials           = errors.New("credenciales invalidas: verifica el Login y el Auth Token de tu tienda Jumpseller")
	ErrMissingAPIKey                = errors.New("falta el Login (api_key) en las credenciales")
	ErrMissingAPISecret             = errors.New("falta el Auth Token (api_secret) en las credenciales")
	ErrMissingBaseURL               = errors.New("el tipo de integracion Jumpseller no tiene base_url configurada en base de datos")
	ErrMissingBaseURLTest           = errors.New("la integracion esta en modo pruebas pero el tipo de integracion Jumpseller no tiene base_url_test configurada en base de datos")
	ErrWebhookInvalidSignature      = errors.New("la firma del webhook de Jumpseller no es valida")
	ErrWebhookCreationFailed        = errors.New("Jumpseller rechazo el registro de los webhooks. La URL de entrega debe ser publica y accesible desde internet: Jumpseller la consulta antes de registrar el hook")
	ErrWebhookMissingToken          = errors.New("falta el hooks_token de la tienda para validar el webhook")
	ErrPerLocationStockNotSupported = errors.New("emparejaste bodegas hacia varias bodegas de Jumpseller, pero todavia no podemos escribir el stock por bodega. Deja un solo destino o usa el modo de una bodega")
	ErrRateLimited                  = errors.New("Jumpseller rechazo la peticion por exceso de solicitudes (rate limit)")
	ErrNoOrdersFound                = errors.New("no se encontraron ordenes en Jumpseller")
	ErrProductNotFound              = errors.New("no se encontro el producto en Jumpseller para el SKU indicado")
	ErrStatusNotMapped              = errors.New("el estado no tiene homologacion con Jumpseller")
)
