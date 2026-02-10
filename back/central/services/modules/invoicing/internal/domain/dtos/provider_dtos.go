package dtos

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
)

// ═══════════════════════════════════════════════════════════════
// DTOs PARA COMUNICACIÓN CON PROVEEDORES DE FACTURACIÓN
// ═══════════════════════════════════════════════════════════════

// InvoiceRequest representa los datos necesarios para crear una factura en el proveedor
type InvoiceRequest struct {
	Invoice      *entities.Invoice
	InvoiceItems []*entities.InvoiceItem
	Provider     *entities.InvoicingProvider
	Config       map[string]interface{}
}

// InvoiceResponse representa la respuesta del proveedor al crear una factura
type InvoiceResponse struct {
	InvoiceNumber string
	ExternalID    string
	InvoiceURL    *string
	PDFURL        *string
	XMLURL        *string
	CUFE          *string
	IssuedAt      string
	RawResponse   map[string]interface{}
}

// CreditNoteRequest representa los datos necesarios para crear una nota de crédito
type CreditNoteRequest struct {
	Invoice    *entities.Invoice
	CreditNote *entities.CreditNote
	Provider   *entities.InvoicingProvider
}

// CreditNoteResponse representa la respuesta del proveedor al crear una nota de crédito
type CreditNoteResponse struct {
	CreditNoteNumber string
	ExternalID       string
	NoteURL          *string
	PDFURL           *string
	XMLURL           *string
	CUFE             *string
	IssuedAt         string
	RawResponse      map[string]interface{}
}

// ═══════════════════════════════════════════════════════════════
// DTOs PARA COMUNICACIÓN CON MÓDULO DE ÓRDENES
// ═══════════════════════════════════════════════════════════════

// OrderData representa los datos mínimos necesarios de una orden para facturación
type OrderData struct {
	// Campos básicos
	ID               string
	BusinessID       uint
	IntegrationID    uint
	OrderNumber      string
	TotalAmount      float64
	Subtotal         float64
	Tax              float64
	Discount         float64
	ShippingCost     float64
	Currency         string
	CustomerName     string
	CustomerEmail    string
	CustomerPhone    string
	CustomerDNI      string
	IsPaid           bool
	PaymentMethodID  uint
	Invoiceable      bool
	Items            []OrderItemData

	// Campos para filtros avanzados
	Status          string
	OrderTypeID     uint
	OrderTypeName   string
	CustomerID      *string
	CustomerType    *string
	ShippingCity    *string
	ShippingState   *string
	ShippingCountry *string
	CreatedAt       time.Time
}

// OrderItemData representa un item de orden
type OrderItemData struct {
	ProductID    *string
	SKU          string
	Name         string
	Description  *string
	Quantity     int
	UnitPrice    float64
	TotalPrice   float64
	Tax          float64
	TaxRate      *float64
	Discount     float64
	CategoryID   *uint
	CategoryName *string
}
