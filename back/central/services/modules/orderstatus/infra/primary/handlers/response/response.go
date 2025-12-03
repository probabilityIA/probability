package response

import "time"

// OrderStatusMappingResponse representa la respuesta de un mapeo de estado
type OrderStatusMappingResponse struct {
	ID              uint      `json:"id"`
	IntegrationType string    `json:"integration_type"`
	OriginalStatus  string    `json:"original_status"`
	MappedStatus    string    `json:"mapped_status"`
	IsActive        bool      `json:"is_active"`
	Priority        int       `json:"priority"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// OrderStatusMappingsListResponse representa la respuesta de lista de mapeos
type OrderStatusMappingsListResponse struct {
	Data  []OrderStatusMappingResponse `json:"data"`
	Total int64                        `json:"total"`
}
