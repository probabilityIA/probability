package response

import "time"

// OrderStatusInfo representa la información del estado de orden
type OrderStatusInfo struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Color       string `json:"color"`
	Priority    int    `json:"priority"`
}

// IntegrationTypeInfo representa la información del tipo de integración
type IntegrationTypeInfo struct {
	ID       uint   `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

// OrderStatusMappingResponse representa la respuesta de un mapeo de estado
type OrderStatusMappingResponse struct {
	ID                uint                 `json:"id"`
	IntegrationTypeID uint                 `json:"integration_type_id"`
	IntegrationType   *IntegrationTypeInfo `json:"integration_type,omitempty"`
	OriginalStatus    string               `json:"original_status"`
	OrderStatusID     uint                 `json:"order_status_id"`
	OrderStatus       *OrderStatusInfo     `json:"order_status,omitempty"`
	IsActive          bool                 `json:"is_active"`
	Description       string               `json:"description"`
	CreatedAt         time.Time            `json:"created_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
}

// OrderStatusMappingsListResponse representa la respuesta de lista de mapeos
type OrderStatusMappingsListResponse struct {
	Data       []OrderStatusMappingResponse `json:"data"`
	Total      int64                        `json:"total"`
	Page       int                          `json:"page"`
	PageSize   int                          `json:"page_size"`
	TotalPages int                          `json:"total_pages"`
}
