package request

// UpdateConfig es el request para actualizar una configuración de facturación
type UpdateConfig struct {
	InvoicingProviderID *uint                   `json:"invoicing_provider_id,omitempty"`
	Enabled             *bool                   `json:"enabled,omitempty"`
	AutoInvoice         *bool                   `json:"auto_invoice,omitempty"`
	Filters             *map[string]interface{} `json:"filters,omitempty"`
}
