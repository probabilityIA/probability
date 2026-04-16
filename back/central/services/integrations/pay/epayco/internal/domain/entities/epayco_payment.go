package entities

// EPaycoPaymentResult resultado del intento de pago en ePayco
type EPaycoPaymentResult struct {
	CheckoutID   string
	RedirectURL  string
	Success      bool
	ErrorMessage string
}
