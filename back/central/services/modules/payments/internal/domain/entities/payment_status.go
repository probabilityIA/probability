package entities

import "time"

// PaymentStatus representa un estado de pago en el sistema
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
