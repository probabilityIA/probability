package request

// ListOrdersFilters representa los filtros para listar órdenes
// ✅ DTO HTTP - CON TAGS (json)
type ListOrdersFilters struct {
	// Paginación
	Page     int `json:"page" form:"page" binding:"omitempty,min=1"`
	PageSize int `json:"page_size" form:"page_size" binding:"omitempty,min=1,max=100"`

	// Filtros
	BusinessID      *uint   `json:"business_id" form:"business_id"`
	IntegrationID   *uint   `json:"integration_id" form:"integration_id"`
	IntegrationType *string `json:"integration_type" form:"integration_type"`
	Platform        *string `json:"platform" form:"platform"`
	Status          *string `json:"status" form:"status"`
	CustomerEmail   *string `json:"customer_email" form:"customer_email"`
	OrderNumber     *string `json:"order_number" form:"order_number"`
	ExternalID      *string `json:"external_id" form:"external_id"`

	// Filtros de fecha
	CreatedFrom *string `json:"created_from" form:"created_from"` // RFC3339 format
	CreatedTo   *string `json:"created_to" form:"created_to"`     // RFC3339 format

	// Ordenamiento
	SortBy    *string `json:"sort_by" form:"sort_by"`       // campo por el que ordenar
	SortOrder *string `json:"sort_order" form:"sort_order"` // asc o desc
}
