package entities

// BoldPaymentResult representa el resultado de la creación de un link de pago Bold
type BoldPaymentResult struct {
	PaymentLinkID string
	CheckoutURL   string
	Status        string
	Success       bool
	ErrorMessage  string
}
