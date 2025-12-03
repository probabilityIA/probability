package domain

import "time"

// OrderStatusMapping representa un mapeo de estado de orden en el dominio
type OrderStatusMapping struct {
	ID              uint
	IntegrationType string
	OriginalStatus  string
	MappedStatus    string
	IsActive        bool
	Priority        int
	Description     string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
