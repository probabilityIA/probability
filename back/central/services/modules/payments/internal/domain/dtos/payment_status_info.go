package dtos

// PaymentStatusInfo DTO para información de estado de pago (sin tags)
type PaymentStatusInfo struct {
	ID          uint
	Code        string
	Name        string
	Description string
	Category    string
	Color       string
	IsActive    bool
}
