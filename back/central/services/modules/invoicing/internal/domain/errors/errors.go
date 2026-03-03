package errors

import "errors"

// Errores de validación
var (
	ErrInvalidInvoiceData    = errors.New("datos de factura inválidos")
	ErrInvalidProviderConfig = errors.New("configuración del proveedor inválida")
	ErrInvalidCredentials    = errors.New("credenciales inválidas")
	ErrInvalidAmount         = errors.New("monto inválido")
	ErrInvalidCurrency       = errors.New("moneda inválida")
	ErrInvalidCustomerData   = errors.New("datos del cliente inválidos")
	ErrInvalidOrderData      = errors.New("datos de la orden inválidos")
	ErrMissingRequiredField  = errors.New("campo requerido faltante")
)

// Errores de estado
var (
	ErrInvoiceAlreadyExists     = errors.New("ya existe una factura para esta orden")
	ErrInvoiceNotFound          = errors.New("factura no encontrada")
	ErrInvoiceAlreadyIssued     = errors.New("la factura ya fue emitida")
	ErrInvoiceAlreadyCancelled  = errors.New("la factura ya fue cancelada")
	ErrInvoiceCannotBeCancelled = errors.New("la factura no puede ser cancelada")
	ErrCancelNotImplemented     = errors.New("la cancelación de facturas aún no está implementada")
	ErrOrderNotInvoiceable      = errors.New("la orden no es facturable")
	ErrOrderAlreadyInvoiced     = errors.New("la orden ya tiene una factura")
)

// Errores de proveedor
var (
	ErrProviderNotFound          = errors.New("proveedor de facturación no encontrado")
	ErrProviderNotActive         = errors.New("el proveedor de facturación no está activo")
	ErrProviderTypeNotFound      = errors.New("tipo de proveedor de facturación no encontrado")
	ErrProviderNotConfigured     = errors.New("proveedor de facturación no configurado para esta integración")
	ErrProviderAPIError          = errors.New("error en la API del proveedor")
	ErrProviderTimeout           = errors.New("tiempo de espera del proveedor agotado")
	ErrProviderUnauthorized      = errors.New("no autorizado por el proveedor")
	ErrProviderRateLimitExceeded = errors.New("límite de solicitudes del proveedor excedido")
)

// Errores de configuración
var (
	ErrConfigNotFound              = errors.New("configuración de facturación no encontrada")
	ErrConfigNotEnabled            = errors.New("la configuración de facturación no está habilitada")
	ErrConfigAlreadyExists         = errors.New("ya existe una configuración de facturación para esta integración")
	ErrActiveInvoicingConfigExists = errors.New("ya existe una configuración de facturación activa para este negocio")
	ErrAutoInvoiceNotEnabled       = errors.New("la facturación automática no está habilitada")
)

// Errores de filtros
var (
	// Monto
	ErrOrderBelowMinAmount = errors.New("el monto de la orden está por debajo del mínimo requerido")
	ErrOrderAboveMaxAmount = errors.New("el monto de la orden supera el máximo permitido")

	// Pago
	ErrOrderNotPaid            = errors.New("la orden no está pagada")
	ErrPaymentMethodNotAllowed = errors.New("el método de pago no está permitido")

	// Orden
	ErrOrderTypeNotAllowed = errors.New("el tipo de orden no está permitido")
	ErrOrderStatusExcluded = errors.New("el estado de la orden está excluido de la facturación")

	// Productos
	ErrProductExcluded   = errors.New("la orden contiene productos excluidos")
	ErrProductNotAllowed = errors.New("la orden contiene productos que no están en la lista permitida")
	ErrMinItemsNotMet    = errors.New("la orden no cumple con el mínimo de artículos requerido")
	ErrMaxItemsExceeded  = errors.New("la orden supera el máximo de artículos permitido")

	// Cliente
	ErrCustomerTypeNotAllowed = errors.New("el tipo de cliente no está permitido")
	ErrCustomerExcluded       = errors.New("el cliente está excluido de la facturación")

	// Ubicación
	ErrShippingRegionNotAllowed = errors.New("la región de envío no está permitida")

	// Fecha
	ErrOrderOutsideDateRange = errors.New("la orden está fuera del rango de fechas permitido")

	// Config
	ErrInvalidFilterConfig = errors.New("configuración de filtros inválida")
)

// Errores de comparación
var (
	ErrCompareDateRangeTooLarge = errors.New("el rango de fechas para comparación no puede superar 7 días")
)

// Errores de sincronización
var (
	ErrSyncFailed         = errors.New("la sincronización falló")
	ErrMaxRetriesExceeded = errors.New("número máximo de reintentos excedido")
	ErrRetryNotAllowed    = errors.New("reintento no permitido")
	ErrSyncInProgress     = errors.New("sincronización ya en progreso")
	ErrSyncLogNotFound    = errors.New("registro de sincronización no encontrado")
	ErrNoRetriesToCancel  = errors.New("no se encontraron reintentos pendientes para cancelar")
)

// Errores de notas de crédito
var (
	ErrCreditNoteNotFound      = errors.New("nota de crédito no encontrada")
	ErrCreditNoteAlreadyIssued = errors.New("la nota de crédito ya fue emitida")
	ErrCreditNoteAmountExceeds = errors.New("el monto de la nota de crédito supera el total de la factura")
	ErrInvoiceNotIssued        = errors.New("la factura debe estar emitida antes de crear una nota de crédito")
)

// Errores de encriptación
var (
	ErrEncryptionFailed     = errors.New("error al encriptar")
	ErrDecryptionFailed     = errors.New("error al desencriptar")
	ErrInvalidEncryptionKey = errors.New("clave de encriptación inválida")
)

// Errores de autenticación con proveedor
var (
	ErrAuthenticationFailed = errors.New("autenticación fallida con el proveedor")
	ErrTokenExpired         = errors.New("token del proveedor expirado")
	ErrTokenRefreshFailed   = errors.New("renovación del token fallida")
)
