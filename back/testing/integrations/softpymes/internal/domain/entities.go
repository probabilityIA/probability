package domain

import "time"

// Invoice representa una factura de SoftPymes
type Invoice struct {
	ID            string
	InvoiceNumber string
	ExternalID    string
	OrderID       string
	Comment       string
	CustomerName  string
	CustomerEmail string
	CustomerNIT   string
	Total         float64
	Subtotal      float64
	IVA           float64
	Currency      string
	CustomerPhone string
	Items         []InvoiceItem
	InvoiceURL    string
	PDFURL        string
	XMLURL        string
	CUFE          string
	IssuedAt      time.Time
	CreatedAt     time.Time
}

// InvoiceItem representa un ítem de factura
type InvoiceItem struct {
	ItemCode    string
	ItemName    string
	Description string
	Quantity    int
	UnitPrice   float64
	Tax         float64
	Discount    float64
	Total       float64
}

// CreditNote representa una nota de crédito
type CreditNote struct {
	ID               string
	CreditNoteNumber string
	ExternalID       string
	InvoiceID        string
	Amount           float64
	Reason           string
	NoteType         string // "total" o "partial"
	NoteURL          string
	PDFURL           string
	XMLURL           string
	CUFE             string
	IssuedAt         time.Time
	CreatedAt        time.Time
}

// Customer representa un cliente en Softpymes
type Customer struct {
	Identification string
	Name           string
	Email          string
	Phone          string
	Branch         string
}

// AuthToken representa un token de autenticación de SoftPymes
type AuthToken struct {
	Token     string
	ExpiresAt time.Time
	APIKey    string
	Referer   string
}
