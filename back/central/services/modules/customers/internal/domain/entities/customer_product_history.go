package entities

import "time"

type CustomerProductHistory struct {
	ID             uint
	CustomerID     uint
	BusinessID     uint
	ProductID      string
	ProductName    string
	ProductSKU     string
	ProductImage   *string
	TimesOrdered   int
	TotalQuantity  int
	TotalSpent     float64
	FirstOrderedAt time.Time
	LastOrderedAt  time.Time
}
