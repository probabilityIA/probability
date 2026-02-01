package entities

import "time"

// PaymentStatus representa un estado de pago en el dominio (PURO - sin tags)
type PaymentStatus struct {
	ID          uint
	Code        string
	Name        string
	Description string
	Category    string
	IsActive    bool
	Icon        string
	Color       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
