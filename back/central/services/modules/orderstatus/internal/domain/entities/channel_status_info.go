package entities

// ChannelStatusInfo contiene información de un estado nativo de un canal de integración
// (ej. "paid" de Shopify, "pending" de MercadoLibre)
// PURO - Sin tags JSON
type ChannelStatusInfo struct {
	ID                uint
	IntegrationTypeID uint
	Code              string
	Name              string
	Description       string
	IsActive          bool
	DisplayOrder      int

	// Relación opcional (disponible en GET y LIST con preload)
	IntegrationType *IntegrationTypeInfo
}
