package response

import "time"

// OrderHistoryResponse representa la respuesta HTTP de un registro de historial
type OrderHistoryResponse struct {
	ID             uint      `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	OrderID        string    `json:"order_id"`
	PreviousStatus string    `json:"previous_status"`
	NewStatus      string    `json:"new_status"`
	ChangedBy      *uint     `json:"changed_by,omitempty"`
	ChangedByName  string    `json:"changed_by_name"`
	Reason         *string   `json:"reason,omitempty"`
}
