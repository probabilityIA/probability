package dtos

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
)

type InvoiceRequest struct {
	Invoice      *entities.Invoice
	InvoiceItems []*entities.InvoiceItem
	Provider     *entities.InvoicingProvider
	Config       map[string]interface{}
}
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
type CreditNoteRequest struct {
	Invoice    *entities.Invoice
	CreditNote *entities.CreditNote
	Provider   *entities.InvoicingProvider
}
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
type OrderData struct {
	ID               string
	BusinessID       uint
	IntegrationID    uint
	OrderNumber      string
	TotalAmount      float64
	Subtotal         float64
	Tax              float64
	Discount         float64
	ShippingCost     float64
	CodCarrierFee    float64
	ShippingDiscount float64
	FreeShipping     bool
	Currency         string
	CustomerName     string
	CustomerEmail    string
	CustomerPhone    string
	CustomerDNI      string
	IsPaid           bool
	IsCOD            bool
	PaymentMethodID  uint
	Invoiceable      bool
	IsTest           bool
	Items            []OrderItemData
	Status           string
	OrderTypeID      uint
	OrderTypeName    string
	CustomerID       *string
	CustomerType     *string
	ShippingCity     *string
	ShippingState    *string
	ShippingCountry  *string
	CreatedAt        time.Time
}
type OrderItemData struct {
	ProductID                *string
	SKU                      string
	Name                     string
	Description              *string
	Quantity                 int
	UnitPrice                float64
	UnitPriceBase            float64
	TotalPrice               float64
	Tax                      float64
	TaxRate                  *float64
	Discount                 float64
	DiscountPercent          float64
	CategoryID               *uint
	CategoryName             *string
	UnitPricePresentment     float64
	UnitPriceBasePresentment float64
	TotalPricePresentment    float64
	DiscountPresentment      float64
	TaxPresentment           float64
}
