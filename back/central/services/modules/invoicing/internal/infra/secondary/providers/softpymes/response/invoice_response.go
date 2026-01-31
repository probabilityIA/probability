package response

import "time"

// InvoiceResponse representa la respuesta de creación de factura de Softpymes
type InvoiceResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Data    *InvoiceData `json:"data,omitempty"`
}

// InvoiceData contiene los datos de la factura creada
type InvoiceData struct {
	InvoiceID     string    `json:"invoice_id"`
	InvoiceNumber string    `json:"invoice_number"`
	CUFE          string    `json:"cufe"`           // Código Único de Factura Electrónica
	PDFURL        string    `json:"pdf_url"`
	XMLURL        string    `json:"xml_url"`
	InvoiceURL    string    `json:"invoice_url"`
	IssuedAt      time.Time `json:"issued_at"`
	Status        string    `json:"status"`
	QRCode        string    `json:"qr_code,omitempty"`
}

// ErrorDetail representa detalles de error de la API
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}
