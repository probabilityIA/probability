package entities

// MeliPagoPaymentResult resultado del intento de pago en MercadoPago
type MeliPagoPaymentResult struct {
	PreferenceID string
	CheckoutURL  string
	Success      bool
	ErrorMessage string
}
