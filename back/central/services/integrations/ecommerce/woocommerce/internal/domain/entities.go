package domain

import "time"

// Integration representa los datos de una integración de WooCommerce
// tal como se obtienen del core de integraciones.
type Integration struct {
	ID              uint
	BusinessID      *uint
	Name            string
	StoreID         string
	IntegrationType int
	Config          map[string]interface{}
}

// WooCommerceOrder representa una orden de WooCommerce (API v3).
// Estructura de dominio pura — sin tags JSON.
type WooCommerceOrder struct {
	ID                  int64
	ParentID            int64
	Number              string
	OrderKey            string
	CreatedVia          string
	Version             string
	Status              string
	Currency            string
	DateCreated         time.Time
	DateModified        time.Time
	DiscountTotal       string
	DiscountTax         string
	ShippingTotal       string
	ShippingTax         string
	CartTax             string
	Total               string
	TotalTax            string
	PricesIncludeTax    bool
	CustomerID          int64
	CustomerNote        string
	Billing             WooCommerceBilling
	Shipping            WooCommerceShipping
	PaymentMethod       string
	PaymentMethodTitle  string
	TransactionID       string
	DatePaid            *time.Time
	DateCompleted       *time.Time
	LineItems           []WooCommerceLineItem
	ShippingLines       []WooCommerceShippingLine
	FeeLines            []WooCommerceFeeLine
	CouponLines         []WooCommerceCouponLine
	MetaData            []WooCommerceMetaData
}

// WooCommerceBilling representa la dirección de facturación.
type WooCommerceBilling struct {
	FirstName string
	LastName  string
	Company   string
	Address1  string
	Address2  string
	City      string
	State     string
	Postcode  string
	Country   string
	Email     string
	Phone     string
}

// WooCommerceShipping representa la dirección de envío.
type WooCommerceShipping struct {
	FirstName string
	LastName  string
	Company   string
	Address1  string
	Address2  string
	City      string
	State     string
	Postcode  string
	Country   string
	Phone     string
}

// WooCommerceLineItem representa un producto de la orden.
type WooCommerceLineItem struct {
	ID          int64
	Name        string
	ProductID   int64
	VariationID int64
	Quantity    int
	TaxClass    string
	Subtotal    string
	SubtotalTax string
	Total       string
	TotalTax    string
	SKU         string
	Price       float64
	ImageURL    string
	MetaData    []WooCommerceMetaData
}

// WooCommerceShippingLine representa una línea de envío.
type WooCommerceShippingLine struct {
	ID          int64
	MethodTitle string
	MethodID    string
	Total       string
	TotalTax    string
	MetaData    []WooCommerceMetaData
}

// WooCommerceFeeLine representa una línea de tarifa adicional.
type WooCommerceFeeLine struct {
	ID        int64
	Name      string
	TaxClass  string
	TaxStatus string
	Total     string
	TotalTax  string
}

// WooCommerceCouponLine representa un cupón aplicado a la orden.
type WooCommerceCouponLine struct {
	ID          int64
	Code        string
	Discount    string
	DiscountTax string
}

// WooCommerceMetaData representa metadata genérica de WooCommerce.
type WooCommerceMetaData struct {
	ID    int64
	Key   string
	Value interface{}
}
