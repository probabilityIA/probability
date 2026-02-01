package request

// UpdatePaymentMapping representa la solicitud HTTP para actualizar un mapeo
type UpdatePaymentMapping struct {
	OriginalMethod  string `json:"original_method" binding:"required,max=128"`
	PaymentMethodID uint   `json:"payment_method_id" binding:"required"`
	Priority        int    `json:"priority"`
}
