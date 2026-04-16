package dtos

// UpdatePaymentMethod representa la solicitud para actualizar un método de pago (PURO - sin tags)
type UpdatePaymentMethod struct {
	Name        string
	Description string
	Category    string // card, digital_wallet, bank_transfer, cash
	Provider    string
	Icon        string
	Color       string
}
