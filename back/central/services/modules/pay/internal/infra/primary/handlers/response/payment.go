package response

import "time"

// PaymentTransactionResponse es la respuesta HTTP de una transacción de pago
type PaymentTransactionResponse struct {
	ID            uint      `json:"id"`
	BusinessID    uint      `json:"business_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	Status        string    `json:"status"`
	GatewayCode   string    `json:"gateway_code"`
	ExternalID    *string   `json:"external_id,omitempty"`
	Reference     string    `json:"reference"`
	PaymentMethod string    `json:"payment_method"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// PaginatedPaymentsResponse respuesta paginada de transacciones
type PaginatedPaymentsResponse struct {
	Data       []*PaymentTransactionResponse `json:"data"`
	Total      int64                         `json:"total"`
	Page       int                           `json:"page"`
	PageSize   int                           `json:"page_size"`
	TotalPages int                           `json:"total_pages"`
}
