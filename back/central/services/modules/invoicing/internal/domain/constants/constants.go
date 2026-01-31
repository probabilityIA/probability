package constants

// Estados de facturas
const (
	InvoiceStatusDraft     = "draft"
	InvoiceStatusPending   = "pending"
	InvoiceStatusIssued    = "issued"
	InvoiceStatusCancelled = "cancelled"
	InvoiceStatusFailed    = "failed"
)

// Estados de notas de crédito
const (
	CreditNoteStatusDraft     = "draft"
	CreditNoteStatusPending   = "pending"
	CreditNoteStatusIssued    = "issued"
	CreditNoteStatusCancelled = "cancelled"
	CreditNoteStatusFailed    = "failed"
)

// Tipos de notas de crédito
const (
	CreditNoteTypeFullRefund    = "full_refund"
	CreditNoteTypePartialRefund = "partial_refund"
	CreditNoteTypeCancellation  = "cancellation"
	CreditNoteTypeCorrection    = "correction"
)

// Estados de sincronización
const (
	SyncStatusPending    = "pending"
	SyncStatusProcessing = "processing"
	SyncStatusSuccess    = "success"
	SyncStatusFailed     = "failed"
)

// Tipos de operación de sincronización
const (
	OperationTypeCreate     = "create"
	OperationTypeCancel     = "cancel"
	OperationTypeCreditNote = "credit_note"
	OperationTypeQuery      = "query"
)

// Triggers de sincronización
const (
	TriggerAuto   = "auto"
	TriggerManual = "manual"
	TriggerRetry  = "retry_job"
)

// Códigos de proveedores
const (
	ProviderCodeSoftpymes = "softpymes"
	ProviderCodeSiigo     = "siigo"
	ProviderCodeFacturama = "facturama"
	ProviderCodeAlegra    = "alegra"
	ProviderCodeNubeFact  = "nubefact"
)

// Reintentos
const (
	MaxRetries              = 3
	DefaultRetryIntervalMin = 5 // minutos
)

// Monedas
const (
	CurrencyCOP = "COP" // Peso colombiano
	CurrencyUSD = "USD" // Dólar estadounidense
	CurrencyMXN = "MXN" // Peso mexicano
	CurrencyPEN = "PEN" // Sol peruano
	CurrencyCLP = "CLP" // Peso chileno
)
