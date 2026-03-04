package request

// UpdateConfig es el request para actualizar una configuración de facturación
type UpdateConfig struct {
	InvoicingProviderID    *uint                   `json:"invoicing_provider_id,omitempty"`    // DEPRECATED: usar invoicing_integration_id
	InvoicingIntegrationID *uint                   `json:"invoicing_integration_id,omitempty"` // Campo actual (FK a integrations)
	IntegrationIDs         *[]uint                 `json:"integration_ids,omitempty"`          // IDs de integraciones de e-commerce (reemplaza todas las actuales)
	Enabled                *bool                   `json:"enabled,omitempty"`
	AutoInvoice            *bool                   `json:"auto_invoice,omitempty"`
	Filters                *map[string]interface{} `json:"filters,omitempty"`
}
