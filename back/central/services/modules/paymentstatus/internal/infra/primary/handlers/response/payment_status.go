package response

// PaymentStatusResponse DTO de respuesta HTTP (con tags JSON)
type PaymentStatusResponse struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Color       string `json:"color"`
}

// PaymentStatusListResponse respuesta de lista de estados de pago
type PaymentStatusListResponse struct {
	Success bool                    `json:"success"`
	Message string                  `json:"message"`
	Data    []PaymentStatusResponse `json:"data"`
}
