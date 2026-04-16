package entities

import "time"

// StorefrontOrder represents an order summary for the storefront customer
type StorefrontOrder struct {
	ID          string
	OrderNumber string
	Status      string
	TotalAmount float64
	Currency    string
	CreatedAt   time.Time
	Items       []StorefrontOrderItem
}

// StorefrontOrderItem represents an item within a storefront order
type StorefrontOrderItem struct {
	ProductName string
	Quantity    int
	UnitPrice   float64
	TotalPrice  float64
	ImageURL    *string
}
