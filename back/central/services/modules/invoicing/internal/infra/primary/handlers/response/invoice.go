package response

import "time"

// Invoice es la respuesta de una factura
type Invoice struct {
	ID                  uint                   `json:"id"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
	OrderID             string                 `json:"order_id"`
	BusinessID          uint                   `json:"business_id"`
	InvoicingProviderID uint                   `json:"invoicing_provider_id"`
	InvoiceNumber       string                 `json:"invoice_number"`
	InternalNumber      string                 `json:"internal_number"`
	ExternalID          *string                `json:"external_id,omitempty"`
	Status              string                 `json:"status"`
	TotalAmount         float64                `json:"total_amount"`
	Subtotal            float64                `json:"subtotal"`
	Tax                 float64                `json:"tax"`
	Discount            float64                `json:"discount"`
	Currency            string                 `json:"currency"`
	IssuedAt            *time.Time             `json:"issued_at,omitempty"`
	CancelledAt         *time.Time             `json:"cancelled_at,omitempty"`
	CUFE                *string                `json:"cufe,omitempty"`
	PDFURL              *string                `json:"pdf_url,omitempty"`
	XMLURL              *string                `json:"xml_url,omitempty"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
	Items               []InvoiceItem          `json:"items,omitempty"`
}

// InvoiceItem es un item de factura
type InvoiceItem struct {
	ID          uint    `json:"id"`
	ProductSKU  *string `json:"product_sku,omitempty"`
	ProductName string  `json:"product_name"`
	Description *string `json:"description,omitempty"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalPrice  float64 `json:"total_price"`
	Tax         float64 `json:"tax"`
	TaxRate     *float64 `json:"tax_rate,omitempty"`
	Discount    float64 `json:"discount"`
}

// InvoiceList es la respuesta de listado de facturas
type InvoiceList struct {
	Items      []Invoice `json:"items"`
	TotalCount int64     `json:"total_count"`
	Page       int       `json:"page"`
	PageSize   int       `json:"page_size"`
}
