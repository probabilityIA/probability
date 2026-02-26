package entities

// PayUPaymentResult resultado del intento de pago en PayU
type PayUPaymentResult struct {
	TransactionID string
	RedirectURL   string
	Success       bool
	ErrorMessage  string
}
