package dtos

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

// OrderSummary representa un resumen de la orden para listados
// ✅ DTO PURO - SIN TAGS
type OrderSummary struct {
	ID                     string
	CreatedAt              time.Time
	BusinessID             uint
	IntegrationID          uint
	IntegrationType        string
	IntegrationLogoURL     *string
	Platform               string
	ExternalID             string
	OrderNumber            string
	TotalAmount            float64
	Currency               string
	TotalAmountPresentment float64
	CurrencyPresentment    string
	CustomerName           string
	CustomerFirstName      string
	CustomerLastName       string
	CustomerEmail          string
	CustomerPhone          string
	ShippingStreet         string
	ShippingCity           string
	ShippingState          string
	Weight                 *float64
	Height                 *float64
	Width                  *float64
	Length                 *float64
	Status                 string
	ItemsCount             int
	DeliveryProbability    *float64
	NegativeFactors        []string
	OrderStatus            *entities.OrderStatusInfo
	PaymentStatus          *entities.PaymentStatusInfo
	FulfillmentStatus      *entities.FulfillmentStatusInfo
	GuideLink              *string
	IsPaid                 bool
	IsConfirmed            *bool
	Novelty                *string
}

// OrderRawResponse representa la respuesta con los datos crudos
type OrderRawResponse struct {
	OrderID       string
	ChannelSource string
	RawData       []byte
}

// OrdersListResponse representa la respuesta paginada de órdenes
type OrdersListResponse struct {
	Data       []OrderSummary
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}
