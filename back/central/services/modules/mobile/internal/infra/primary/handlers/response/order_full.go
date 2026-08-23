package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/mobile/internal/domain/entities"
)

type OrderSummary struct {
	ID              string    `json:"id"`
	OrderNumber     string    `json:"order_number"`
	InternalNumber  string    `json:"internal_number"`
	IntegrationID   uint      `json:"integration_id"`
	IntegrationType string    `json:"integration_type"`
	Platform        string    `json:"platform"`
	Status          string    `json:"status"`
	StatusID        *uint     `json:"status_id"`
	IsPaid          bool      `json:"is_paid"`
	IsCod           bool      `json:"is_cod"`
	CodTotal        *float64  `json:"cod_total"`
	Subtotal        float64   `json:"subtotal"`
	Tax             float64   `json:"tax"`
	Discount        float64   `json:"discount"`
	ShippingCost    float64   `json:"shipping_cost"`
	TotalAmount     float64   `json:"total_amount"`
	Currency        string    `json:"currency"`
	CustomerName    string    `json:"customer_name"`
	CustomerEmail   string    `json:"customer_email"`
	CustomerPhone   string    `json:"customer_phone"`
	ShippingStreet  string    `json:"shipping_street"`
	ShippingCity    string    `json:"shipping_city"`
	WarehouseName   string    `json:"warehouse_name"`
	UserName        string    `json:"user_name"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type OrderLine struct {
	ID         uint    `json:"id"`
	ProductID  *string `json:"product_id"`
	SKU        string  `json:"sku"`
	Name       string  `json:"name"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
}

type ShipmentSummary struct {
	ID                   uint      `json:"id"`
	TrackingNumber       *string   `json:"tracking_number"`
	TrackingURL          *string   `json:"tracking_url"`
	Carrier              *string   `json:"carrier"`
	GuideURL             *string   `json:"guide_url"`
	Status               string    `json:"status"`
	CarrierStatus        *string   `json:"carrier_status"`
	DestinationCity      string    `json:"destination_city"`
	InsuranceCost        *float64  `json:"insurance_cost"`
	TotalCost            *float64  `json:"total_cost"`
	CarrierCost          *float64  `json:"carrier_cost"`
	AppliedMargin        *float64  `json:"applied_margin"`
	CodCarrierFee        *float64  `json:"cod_carrier_fee"`
	CodProbabilityMargin *float64  `json:"cod_probability_margin"`
	CreatedAt            time.Time `json:"created_at"`
}

type InvoiceSummary struct {
	ID            uint       `json:"id"`
	InvoiceNumber string     `json:"invoice_number"`
	Status        string     `json:"status"`
	TotalAmount   float64    `json:"total_amount"`
	Currency      string     `json:"currency"`
	InvoiceURL    *string    `json:"invoice_url"`
	CUFE          *string    `json:"cufe"`
	IssuedAt      *time.Time `json:"issued_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

type OrderFull struct {
	Order    OrderSummary     `json:"order"`
	Items    []OrderLine      `json:"items"`
	Shipment *ShipmentSummary `json:"shipment"`
	Invoice  *InvoiceSummary  `json:"invoice"`
}

func FromOrderFull(source *entities.OrderFull) OrderFull {
	if source == nil {
		return OrderFull{}
	}

	items := make([]OrderLine, 0, len(source.Items))
	for _, item := range source.Items {
		items = append(items, OrderLine{
			ID:         item.ID,
			ProductID:  item.ProductID,
			SKU:        item.SKU,
			Name:       item.Name,
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice,
			TotalPrice: item.TotalPrice,
		})
	}

	result := OrderFull{
		Order: OrderSummary(source.Order),
		Items: items,
	}

	if source.Shipment != nil {
		shipment := ShipmentSummary(*source.Shipment)
		result.Shipment = &shipment
	}

	if source.Invoice != nil {
		invoice := InvoiceSummary(*source.Invoice)
		result.Invoice = &invoice
	}

	return result
}
