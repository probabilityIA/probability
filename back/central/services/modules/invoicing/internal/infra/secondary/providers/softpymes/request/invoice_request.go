package request

import "time"

// InvoiceRequest representa la solicitud de creación de factura a Softpymes
type InvoiceRequest struct {
	// Información general
	Referer    string    `json:"referer"`     // NIT del facturador
	BranchCode string    `json:"branch_code"` // Código de sucursal
	Date       time.Time `json:"date"`        // Fecha de factura
	DueDate    *time.Time `json:"due_date,omitempty"` // Fecha de vencimiento

	// Información del cliente
	Customer CustomerData `json:"customer"`

	// Items de la factura
	Items []InvoiceItem `json:"items"`

	// Información financiera
	Currency     string  `json:"currency"`      // Código de moneda (COP, USD)
	Notes        *string `json:"notes,omitempty"` // Notas adicionales
	ExchangeRate float64 `json:"exchange_rate,omitempty"` // Tasa de cambio si no es COP

	// Información de pago
	PaymentMethod string `json:"payment_method,omitempty"` // Método de pago
}

// CustomerData representa los datos del cliente en la factura
type CustomerData struct {
	IdentificationType string `json:"identification_type"` // "CC", "NIT", "CE", "Pasaporte"
	IdentificationNumber string `json:"identification_number"`
	Name               string `json:"name"`
	Email              string `json:"email,omitempty"`
	Phone              string `json:"phone,omitempty"`
	Address            string `json:"address,omitempty"`
	City               string `json:"city,omitempty"`
	Country            string `json:"country,omitempty"`
}

// InvoiceItem representa un item de la factura
type InvoiceItem struct {
	Code        string  `json:"code,omitempty"`        // Código/SKU del producto
	Description string  `json:"description"`           // Descripción del producto
	Quantity    float64 `json:"quantity"`              // Cantidad
	UnitPrice   float64 `json:"unit_price"`            // Precio unitario
	TotalPrice  float64 `json:"total_price"`           // Precio total
	Tax         float64 `json:"tax,omitempty"`         // Impuesto
	TaxRate     float64 `json:"tax_rate,omitempty"`    // Tasa de impuesto (0.19 = 19%)
	Discount    float64 `json:"discount,omitempty"`    // Descuento
}
