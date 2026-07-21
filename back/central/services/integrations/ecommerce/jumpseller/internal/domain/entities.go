package domain

import "time"

type Integration struct {
	ID              uint
	BusinessID      *uint
	Name            string
	StoreID         string
	IntegrationType int
	Config          map[string]interface{}
	IsTesting       bool
	BaseURL         string
	BaseURLTest     string
}

type TokenResponse struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int
	CreatedAt    int64
}

type StoreInfo struct {
	Code             string
	Name             string
	URL              string
	Country          string
	Currency         string
	HooksToken       string
	WeightUnit       string
	SubscriptionPlan string
}

type LocationsInfo struct {
	Locations        []Location
	SubscriptionPlan string
	MultiLocation    bool
	StockOriginName  string
}

type Location struct {
	ID            int64
	Name          string
	Main          bool
	IsStockOrigin bool
	PickupPoint   bool
	City          string
	Country       string
}

type WebhookItem struct {
	ID        string
	Address   string
	Topic     string
	Format    string
	CreatedAt string
}

type JumpsellerOrder struct {
	ID                int64
	CreatedAt         time.Time
	Status            string
	Currency          string
	Subtotal          float64
	Tax               float64
	ShippingTax       float64
	Shipping          float64
	ShippingRequired  bool
	Total             float64
	Discount          float64
	ShippingDiscount  float64
	FulfillmentStatus string
	ShipmentStatus    string
	ShippingMethodID  int64
	ShippingMethod    string
	PaymentMethodName string
	PaymentMethodType string
	PaymentInfo       string
	AdditionalInfo    string
	TrackingNumber    string
	TrackingCompany   string
	TrackingURL       string
	ShippingOption    string
	Customer          OrderCustomer
	ShippingAddress   Address
	BillingAddress    Address
	Products          []OrderProduct
	AdditionalFields  []OrderAdditionalField
}

type OrderCustomer struct {
	ID    string
	Name  string
	Email string
	Phone string
	IP    string
}

type Address struct {
	Name         string
	Surname      string
	TaxID        string
	Address      string
	StreetNumber string
	City         string
	Postal       string
	Region       string
	Country      string
	CountryCode  string
	RegionCode   string
	Latitude     *float64
	Longitude    *float64
}

type OrderProduct struct {
	ID        int64
	VariantID int64
	SKU       string
	Name      string
	Qty       int
	Price     float64
	Tax       float64
	Discount  float64
	Weight    float64
}

type OrderAdditionalField struct {
	ID    int64
	Label string
	Value string
	Area  string
}

type JumpsellerProduct struct {
	ID             int64
	Name           string
	SKU            string
	Barcode        string
	Description    string
	Price          float64
	Stock          int
	StockUnlimited bool
	Status         string
	Weight         float64
	Height         float64
	Width          float64
	Length         float64
	Diameter       float64
	PackageFormat  string
	Variants       []ProductVariant
}

type ProductVariant struct {
	ID             int64
	SKU            string
	Price          float64
	Stock          int
	StockUnlimited bool
}

type StockTarget struct {
	ProductID int64
	VariantID int64
	Found     bool
}

type UpdateOrderFields struct {
	Status          string
	ShipmentStatus  string
	TrackingNumber  string
	TrackingCompany string
	AdditionalInfo  string
}
