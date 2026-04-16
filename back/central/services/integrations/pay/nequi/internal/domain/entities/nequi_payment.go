package entities

// NequiPaymentResult resultado del intento de pago en Nequi
type NequiPaymentResult struct {
	QRValue       string
	TransactionID string
	Success       bool
	ErrorMessage  string
}
