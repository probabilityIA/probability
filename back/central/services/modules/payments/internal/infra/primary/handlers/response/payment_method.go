package response

import "time"

// PaymentMethod representa la respuesta HTTP de un método de pago
type PaymentMethod struct {
	ID          uint      `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Provider    string    `json:"provider"`
	IsActive    bool      `json:"is_active"`
	Icon        string    `json:"icon"`
	Color       string    `json:"color"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PaymentMethodsList representa la respuesta paginada HTTP de métodos de pago
type PaymentMethodsList struct {
	Data       []PaymentMethod `json:"data"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}
