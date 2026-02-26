package entities

// StripePaymentResult resultado del intento de pago en Stripe
type StripePaymentResult struct {
	PaymentIntentID string
	ClientSecret    string
	Success         bool
	ErrorMessage    string
}
