package domain

import (
	"github.com/secamc93/probability/back/central/shared/productmatch"
)

type Integration struct {
	ID                uint
	BusinessID        *uint
	Name              string
	StoreID           string
	IntegrationType   int
	Config            map[string]interface{}
	IsTesting         bool
	BaseURL           string
	BaseURLTest       string
	ProductMatchRules []productmatch.Rule
}

type Credential struct {
	AccessToken string
	StoreID     string
	BaseURL     string
	UserAgent   string
}

type StoreInfo struct {
	ID       int64
	Name     string
	URL      string
	Country  string
	Currency string
	Language string
}

type TiendanubeVariant struct {
	ID              int64
	ProductID       int64
	SKU             string
	Barcode         string
	Price           float64
	Stock           int
	StockManagement bool
	Weight          float64
	Depth           float64
	Width           float64
	Height          float64
}

type TiendanubeProduct struct {
	ID          int64
	Name        string
	Description string
	ImageURL    string
	Published   bool
	Variants    []TiendanubeVariant
}

type StockTarget struct {
	ProductID int64
	VariantID int64
	Found     bool
}

type WebhookItem struct {
	ID        string
	Address   string
	Topic     string
	Format    string
	CreatedAt string
	UpdatedAt string
}

type CreateWebhooksResult struct {
	WebhookURL       string
	ExistingWebhooks []WebhookItem
	CreatedWebhooks  []string
	FailedWebhooks   []string
}

type TiendanubeAddress struct {
	FirstName string
	LastName  string
	Name      string
	Phone     string
	Street    string
	Number    string
	Floor     string
	Locality  string
	City      string
	Province  string
	Country   string
	Zipcode   string
}

type TiendanubeOrderItem struct {
	ID        string
	ProductID string
	VariantID string
	Name      string
	SKU       string
	Price     float64
	Quantity  int
	Weight    float64
	ImageURL  string
}

type TiendanubeOrder struct {
	ID                    string
	Number                string
	Currency              string
	Status                string
	PaymentStatus         string
	ShippingStatus        string
	Subtotal              float64
	Discount              float64
	Total                 float64
	ShippingCost          float64
	Weight                float64
	ContactName           string
	ContactEmail          string
	ContactPhone          string
	ContactIdentification string
	Gateway               string
	GatewayName           string
	Note                  string
	ShippingOption        string
	TrackingNumber        string
	TrackingURL           string
	BillingAddress        TiendanubeAddress
	ShippingAddress       TiendanubeAddress
	Items                 []TiendanubeOrderItem
	CreatedAt             string
	UpdatedAt             string
	PaidAt                string
	CancelledAt           string
}

type OrderFilters struct {
	CreatedAtMin  string
	CreatedAtMax  string
	Status        string
	PaymentStatus string
	Limit         int
}
