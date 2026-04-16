package dtos

import "time"

type ListCustomerAddressesParams struct {
	CustomerID uint
	BusinessID uint
	Page       int
	PageSize   int
}

type ListCustomerProductsParams struct {
	CustomerID uint
	BusinessID uint
	Page       int
	PageSize   int
	SortBy     string
}

type ListCustomerOrderItemsParams struct {
	CustomerID uint
	BusinessID uint
	Page       int
	PageSize   int
}

type OrderEventItemDTO struct {
	ProductID   *string
	ProductName string
	ProductSKU  string
	ProductImage *string
	Quantity    int
	UnitPrice   float64
	TotalPrice  float64
}

type OrderEventDTO struct {
	EventType           string
	OrderID             string
	BusinessID          uint
	CustomerID          *uint
	CustomerName        string
	CustomerEmail       string
	CustomerPhone       string
	CustomerDNI         string
	TotalAmount         float64
	Currency            string
	Platform            string
	Status              string
	IsPaid              bool
	DeliveryProbability float64
	ShippingStreet      string
	ShippingCity        string
	ShippingState       string
	ShippingCountry     string
	ShippingPostalCode  string
	ShippingLat         *float64
	ShippingLng         *float64
	OrderNumber         string
	OrderedAt           time.Time
	Items               []OrderEventItemDTO
	PreviousStatus      string
	CurrentStatus       string
}
