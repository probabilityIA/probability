package domain

import "time"

type ShopifyOrdersResponse struct {
	Orders []ShopifyAPIOrder
}

type ShopifyAPIOrder struct {
	ID                  int64
	Name                string
	OrderNumber         int
	Email               string
	Phone               string
	CreatedAt           string
	UpdatedAt           string
	ProcessedAt         string
	Currency            string
	TotalPrice          string
	SubtotalPrice       string
	TotalTax            string
	TotalDiscounts      string
	FinancialStatus     string
	FulfillmentStatus   *string
	SourceName          string
	PaymentGatewayNames []string
	Customer            *ShopifyAPICustomer
	ShippingAddress     *ShopifyAPIAddress
	BillingAddress      *ShopifyAPIAddress
	LineItems           []ShopifyLineItem
	ShippingLines       []ShopifyShippingLine
	Fulfillments        []ShopifyFulfillment
	LocationID          *int64
	Note                *string
	Tags                string
	TotalWeight         int
	RawData             map[string]interface{}
}

type ShopifyAPICustomer struct {
	ID            int64
	Email         string
	FirstName     string
	LastName      string
	Phone         *string
	VerifiedEmail bool
	OrdersCount   int
	State         string
	TotalSpent    string
	CreatedAt     string
	UpdatedAt     string
}

type ShopifyAPIAddress struct {
	FirstName    string
	LastName     string
	Company      *string
	Address1     string
	Address2     *string
	City         string
	Province     string
	ProvinceCode string
	Country      string
	CountryCode  string
	Zip          string
	Phone        *string
	Latitude     *float64
	Longitude    *float64
}

type ShopifyLineItem struct {
	ID                int64
	VariantID         *int64
	ProductID         *int64
	Title             string
	VariantTitle      *string
	SKU               string
	Quantity          int
	Price             string
	Grams             int
	TotalDiscount     string
	FulfillmentStatus *string
	Name              string
}

type ShopifyShippingLine struct {
	ID                            int64
	Title                         string
	Price                         string
	Code                          string
	Source                        string
	Phone                         *string
	RequestedFulfillmentServiceID *string
	DeliveryCategory              *string
	CarrierIdentifier             *string
}

type ShopifyFulfillment struct {
	ID              int64
	OrderID         int64
	Status          string
	CreatedAt       string
	UpdatedAt       string
	TrackingCompany *string
	TrackingNumber  *string
	TrackingNumbers []string
	TrackingURL     *string
	TrackingURLs    []string
	ShipmentStatus  *string
}

type FetchOrdersParams struct {
	Status       string
	Limit        int
	CreatedAtMin *time.Time
	CreatedAtMax *time.Time
	UpdatedAtMin *time.Time
	Fields       []string
	SinceID      *int64
}
