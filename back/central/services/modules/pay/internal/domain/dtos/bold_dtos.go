package dtos

// BoldSignatureResponse contiene la firma de integridad para Bold.co
type BoldSignatureResponse struct {
	OrderID            string  `json:"order_id"`
	IntegritySignature string  `json:"integrity_signature"`
	Amount             float64 `json:"amount"`
	Currency           string  `json:"currency"`
	IdentityKey        string  `json:"identity_key"`
}

// BoldStatusResponse representa la respuesta de estado desde Bold.co
type BoldStatusResponse struct {
	BoldOrderID   string  `json:"bold_order_id"`
	Status        string  `json:"status"` // e.g., 'APPROVED', 'REJECTED', 'PENDING'
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	PaymentMethod string  `json:"payment_method"`
}
