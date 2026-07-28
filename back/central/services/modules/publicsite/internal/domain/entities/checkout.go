package entities

import "time"

type CheckoutItem struct {
	ProductID string  `json:"product_id"`
	SKU       string  `json:"sku"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	ImageURL  string  `json:"image_url"`
}

type CheckoutAddress struct {
	Street       string `json:"street"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
	PostalCode   string `json:"postal_code"`
	Instructions string `json:"instructions"`
}

type PublicCheckout struct {
	ID              string
	BusinessID      uint
	IntegrationID   uint
	Slug            string
	Reference       string
	Status          string
	Amount          float64
	Currency        string
	Items           []CheckoutItem
	CustomerName    string
	CustomerEmail   string
	CustomerPhone   string
	CustomerDni     string
	ShippingAddress *CheckoutAddress
	BoldOrderID     string
	CreatedAt       time.Time
	PaidAt          *time.Time
}

const (
	CheckoutStatusPending = "pending"
	CheckoutStatusPaid    = "paid"
	CheckoutStatusFailed  = "failed"
	CheckoutStatusExpired = "expired"
)
