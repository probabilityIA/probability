package dtos

import "time"

// OrderHistoryResponse representa un registro de cambio de estado
type OrderHistoryResponse struct {
	ID             uint
	CreatedAt      time.Time
	OrderID        string
	PreviousStatus string
	NewStatus      string
	ChangedBy      *uint
	ChangedByName  string
	Reason         *string
}
