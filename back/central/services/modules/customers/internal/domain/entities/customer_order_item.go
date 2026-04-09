package entities

import "time"

type CustomerOrderItem struct {
	ID           uint
	CustomerID   uint
	BusinessID   uint
	OrderID      string
	OrderNumber  string
	ProductID    *string
	ProductName  string
	ProductSKU   string
	ProductImage *string
	Quantity     int
	UnitPrice    float64
	TotalPrice   float64
	OrderStatus  string
	OrderedAt    time.Time
}
