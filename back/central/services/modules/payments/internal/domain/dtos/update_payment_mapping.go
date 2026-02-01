package dtos

// UpdatePaymentMapping representa la solicitud para actualizar un mapeo (PURO - sin tags)
type UpdatePaymentMapping struct {
	OriginalMethod  string
	PaymentMethodID uint
	Priority        int
}
