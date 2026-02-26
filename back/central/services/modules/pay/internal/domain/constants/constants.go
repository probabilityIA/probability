package constants

// Estados de transacción
const (
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
	StatusCancelled  = "cancelled"
)

// Gateways disponibles
const (
	GatewayNequi = "nequi"
)

// Métodos de pago
const (
	PaymentMethodQRCode      = "qr_code"
	PaymentMethodPaymentLink = "payment_link"
)

// Colas RabbitMQ
const (
	QueuePayRequests  = "pay.requests"
	QueuePayResponses = "pay.responses"
)

// Configuración de reintentos
const (
	MaxRetries = 3
)
