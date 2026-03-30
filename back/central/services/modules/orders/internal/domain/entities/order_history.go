package entities

import "time"

// OrderHistory representa un registro de cambio de estado de una orden
type OrderHistory struct {
	ID             uint
	CreatedAt      time.Time
	OrderID        string
	PreviousStatus string
	NewStatus      string
	ChangedBy      *uint
	ChangedByName  string
	Reason         *string
	Metadata       []byte
}
