package domain

import "time"

// Shopify API Response Types
// Based on Shopify Admin API 2024-10

type ShopifyOrdersResponse struct {
	Orders []ShopifyOrder `json:"orders"`
}

type ShopifyOrder struct {
	ID                  int64                  `json:"id"`
	Name                string                 `json:"name"`
	OrderNumber         int                    `json:"order_number"`
	Email               string                 `json:"email"`
	Phone               string                 `json:"phone"`
	CreatedAt           string                 `json:"created_at"`
	UpdatedAt           string                 `json:"updated_at"`
	ProcessedAt         string                 `json:"processed_at"`
	Currency            string                 `json:"currency"`
	TotalPrice          string                 `json:"total_price"`
	SubtotalPrice       string                 `json:"subtotal_price"`
	TotalTax            string                 `json:"total_tax"`
	TotalDiscounts      string                 `json:"total_discounts"`
	FinancialStatus     string                 `json:"financial_status"`
	FulfillmentStatus   *string                `json:"fulfillment_status"`
	SourceName          string                 `json:"source_name"`
	PaymentGatewayNames []string               `json:"payment_gateway_names"`
	Customer            *ShopifyCustomer       `json:"customer"`
	ShippingAddress     *ShopifyAddress        `json:"shipping_address"`
	BillingAddress      *ShopifyAddress        `json:"billing_address"`
	LineItems           []ShopifyLineItem      `json:"line_items"`
	ShippingLines       []ShopifyShippingLine  `json:"shipping_lines"`
	Fulfillments        []ShopifyFulfillment   `json:"fulfillments"`
	LocationID          *int64                 `json:"location_id"`
	Note                *string                `json:"note"`
	Tags                string                 `json:"tags"`
	TotalWeight         int                    `json:"total_weight"`
	RawData             map[string]interface{} `json:"-"` // For storing complete response
}

type ShopifyCustomer struct {
	ID            int64   `json:"id"`
	Email         string  `json:"email"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	Phone         *string `json:"phone"`
	VerifiedEmail bool    `json:"verified_email"`
	OrdersCount   int     `json:"orders_count"`
	State         string  `json:"state"`
	TotalSpent    string  `json:"total_spent"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type ShopifyAddress struct {
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name"`
	Company      *string  `json:"company"`
	Address1     string   `json:"address1"`
	Address2     *string  `json:"address2"`
	City         string   `json:"city"`
	Province     string   `json:"province"`
	ProvinceCode string   `json:"province_code"`
	Country      string   `json:"country"`
	CountryCode  string   `json:"country_code"`
	Zip          string   `json:"zip"`
	Phone        *string  `json:"phone"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
}

type ShopifyLineItem struct {
	ID                int64   `json:"id"`
	VariantID         *int64  `json:"variant_id"`
	ProductID         *int64  `json:"product_id"`
	Title             string  `json:"title"`
	VariantTitle      *string `json:"variant_title"`
	SKU               string  `json:"sku"`
	Quantity          int     `json:"quantity"`
	Price             string  `json:"price"`
	Grams             int     `json:"grams"`
	TotalDiscount     string  `json:"total_discount"`
	FulfillmentStatus *string `json:"fulfillment_status"`
	Name              string  `json:"name"`
}

type ShopifyShippingLine struct {
	ID                            int64   `json:"id"`
	Title                         string  `json:"title"`
	Price                         string  `json:"price"`
	Code                          string  `json:"code"`
	Source                        string  `json:"source"`
	Phone                         *string `json:"phone"`
	RequestedFulfillmentServiceID *string `json:"requested_fulfillment_service_id"`
	DeliveryCategory              *string `json:"delivery_category"`
	CarrierIdentifier             *string `json:"carrier_identifier"`
}

type ShopifyFulfillment struct {
	ID              int64    `json:"id"`
	OrderID         int64    `json:"order_id"`
	Status          string   `json:"status"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
	TrackingCompany *string  `json:"tracking_company"`
	TrackingNumber  *string  `json:"tracking_number"`
	TrackingNumbers []string `json:"tracking_numbers"`
	TrackingURL     *string  `json:"tracking_url"`
	TrackingURLs    []string `json:"tracking_urls"`
	ShipmentStatus  *string  `json:"shipment_status"`
}

// FetchOrdersParams contains parameters for fetching orders from Shopify
type FetchOrdersParams struct {
	Status       string     // "any", "open", "closed", "cancelled"
	Limit        int        // Max 250
	CreatedAtMin *time.Time // Filter orders created after this date
	CreatedAtMax *time.Time // Filter orders created before this date
	UpdatedAtMin *time.Time // Filter orders updated after this date
	Fields       []string   // Specific fields to retrieve
	SinceID      *int64     // For pagination
}
