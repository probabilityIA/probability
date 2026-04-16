package entities

// WompiPaymentResult resultado del intento de pago en Wompi
type WompiPaymentResult struct {
	TransactionID string
	RedirectURL   string
	Success       bool
	ErrorMessage  string
}
