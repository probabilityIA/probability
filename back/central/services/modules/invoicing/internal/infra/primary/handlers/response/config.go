package response

import "time"

// Config es la respuesta de una configuración de facturación
type Config struct {
	ID                  uint                   `json:"id"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
	BusinessID          uint                   `json:"business_id"`
	IntegrationID       uint                   `json:"integration_id"`
	InvoicingProviderID uint                   `json:"invoicing_provider_id"`
	Enabled             bool                   `json:"enabled"`
	AutoInvoice         bool                   `json:"auto_invoice"`
	Filters             map[string]interface{} `json:"filters,omitempty"`

	// Nombres de relaciones (para frontend)
	IntegrationName  *string `json:"integration_name,omitempty"`
	ProviderName     *string `json:"provider_name,omitempty"`
	ProviderImageURL *string `json:"provider_image_url,omitempty"`
	Description      *string `json:"description,omitempty"`
}

// ConfigList es la respuesta de listado de configuraciones
type ConfigList struct {
	Items      []Config `json:"items"`
	TotalCount int64    `json:"total_count"`
	Page       int      `json:"page"`
	PageSize   int      `json:"page_size"`
}
