package dtos

// CreatePaymentMethod representa la solicitud para crear un método de pago (PURO - sin tags)
type CreatePaymentMethod struct {
	Code        string
	Name        string
	Description string
	Category    string // card, digital_wallet, bank_transfer, cash
	Provider    string
	Icon        string
	Color       string
}
