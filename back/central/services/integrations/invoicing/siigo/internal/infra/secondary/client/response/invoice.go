package response

// CreateInvoiceResponse representa la respuesta de Siigo al crear una factura
type CreateInvoiceResponse struct {
	ID           string               `json:"id"`
	Document     InvoiceDocumentRef   `json:"document"`
	Date         string               `json:"date"`
	Number       int                  `json:"number"`
	Name         string               `json:"name"`         // Número completo: "FV-123"
	Customer     InvoiceCustomerInfo  `json:"customer"`
	Total        float64              `json:"total"`
	TotalTax     float64              `json:"total_tax"`
	Balance      float64              `json:"balance"`
	ErrorCode    string               `json:"error_code,omitempty"`
	Errors       []SiigoError         `json:"Errors,omitempty"`
	PublicURL    string               `json:"public_url,omitempty"`
	Metadata     InvoiceMetadata      `json:"metadata,omitempty"`
}

// InvoiceDocumentRef referencia al tipo de documento
type InvoiceDocumentRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// InvoiceCustomerInfo información del cliente en la factura
type InvoiceCustomerInfo struct {
	Identification string `json:"identification"`
	Name           string `json:"name"`
}

// InvoiceMetadata metadatos DIAN
type InvoiceMetadata struct {
	CUFE string `json:"cufe,omitempty"`
	QR   string `json:"qr,omitempty"`
}

// SiigoError error retornado por Siigo
type SiigoError struct {
	Code    string `json:"Code,omitempty"`
	Message string `json:"Message,omitempty"`
	Params  string `json:"Params,omitempty"`
}
