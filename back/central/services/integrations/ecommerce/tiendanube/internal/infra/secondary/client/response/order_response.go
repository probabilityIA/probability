package response

import "encoding/json"

type OrderAddress struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Number    string `json:"number"`
	Floor     string `json:"floor"`
	Locality  string `json:"locality"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	Zipcode   string `json:"zipcode"`
}

type OrderProduct struct {
	ID        json.Number `json:"id"`
	ProductID json.Number `json:"product_id"`
	VariantID json.Number `json:"variant_id"`
	Name      string      `json:"name"`
	SKU       string      `json:"sku"`
	Price     string      `json:"price"`
	Quantity  json.Number `json:"quantity"`
	Weight    string      `json:"weight"`
	ImageURL  string      `json:"image_url"`
}

type OrderCustomer struct {
	ID             json.Number `json:"id"`
	Name           string      `json:"name"`
	Email          string      `json:"email"`
	Phone          string      `json:"phone"`
	Identification string      `json:"identification"`
}

type Order struct {
	ID       json.Number `json:"id"`
	Number   json.Number `json:"number"`
	Token    string      `json:"token"`
	StoreID  json.Number `json:"store_id"`
	Currency string      `json:"currency"`

	Status         string `json:"status"`
	PaymentStatus  string `json:"payment_status"`
	ShippingStatus string `json:"shipping_status"`

	Subtotal             string `json:"subtotal"`
	Discount             string `json:"discount"`
	Total                string `json:"total"`
	ShippingCostCustomer string `json:"shipping_cost_customer"`
	ShippingCostOwner    string `json:"shipping_cost_owner"`
	Weight               string `json:"weight"`

	ContactName           string `json:"contact_name"`
	ContactEmail          string `json:"contact_email"`
	ContactPhone          string `json:"contact_phone"`
	ContactIdentification string `json:"contact_identification"`

	Gateway     string          `json:"gateway"`
	GatewayName string          `json:"gateway_name"`
	Note        string          `json:"note"`
	Coupon      json.RawMessage `json:"coupon"`

	ShippingOption         string `json:"shipping_option"`
	ShippingTrackingNumber string `json:"shipping_tracking_number"`
	ShippingTrackingURL    string `json:"shipping_tracking_url"`

	BillingName     string `json:"billing_name"`
	BillingPhone    string `json:"billing_phone"`
	BillingAddress  string `json:"billing_address"`
	BillingNumber   string `json:"billing_number"`
	BillingFloor    string `json:"billing_floor"`
	BillingLocality string `json:"billing_locality"`
	BillingCity     string `json:"billing_city"`
	BillingProvince string `json:"billing_province"`
	BillingCountry  string `json:"billing_country"`
	BillingZipcode  string `json:"billing_zipcode"`

	ShippingAddress OrderAddress   `json:"shipping_address"`
	Customer        OrderCustomer  `json:"customer"`
	Products        []OrderProduct `json:"products"`

	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	PaidAt      string `json:"paid_at"`
	CancelledAt string `json:"cancelled_at"`
	ClosedAt    string `json:"closed_at"`
}
