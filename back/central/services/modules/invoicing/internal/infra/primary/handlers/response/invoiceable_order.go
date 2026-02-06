package response

import "time"

// InvoiceableOrder representa una orden facturable en la respuesta HTTP
type InvoiceableOrder struct {
	ID           string    `json:"id"`
	BusinessID   uint      `json:"business_id"`   // Para que super admin vea de qué business es
	OrderNumber  string    `json:"order_number"`
	CustomerName string    `json:"customer_name"`
	TotalAmount  float64   `json:"total_amount"`
	Currency     string    `json:"currency"`
	CreatedAt    time.Time `json:"created_at"`
}

// PaginatedInvoiceableOrders representa la respuesta paginada de órdenes facturables
type PaginatedInvoiceableOrders struct {
	Data     []InvoiceableOrder `json:"data"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}
