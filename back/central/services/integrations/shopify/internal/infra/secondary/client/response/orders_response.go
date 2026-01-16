package response

import "time"

// OrdersResponse representa la respuesta de la API de Shopify para múltiples órdenes
type OrdersResponse struct {
	Orders []Order `json:"orders"`
}

// OrderResponse representa la respuesta de la API de Shopify para una sola orden
type OrderResponse struct {
	Order Order `json:"order"`
}

// Order representa una orden completa de Shopify
type Order struct {
	ID                       int64                 `json:"id"`
	AdminGraphQLAPIID        string                `json:"admin_graphql_api_id"`
	AppID                    *int64                `json:"app_id"`
	BrowserIP                *string               `json:"browser_ip"`
	BuyerAcceptsMarketing    bool                  `json:"buyer_accepts_marketing"`
	CancelReason             *string               `json:"cancel_reason"`
	CancelledAt              *time.Time            `json:"cancelled_at"`
	CartToken                *string               `json:"cart_token"`
	CheckoutID               *int64                `json:"checkout_id"`
	CheckoutToken            *string               `json:"checkout_token"`
	ClientDetails            *ClientDetails        `json:"client_details"`
	ClosedAt                 *time.Time            `json:"closed_at"`
	Confirmed                bool                  `json:"confirmed"`
	ContactEmail             string                `json:"contact_email"`
	CreatedAt                time.Time             `json:"created_at"`
	Currency                 string                `json:"currency"`
	CurrentSubtotalPrice     string                `json:"current_subtotal_price"`
	CurrentSubtotalPriceSet  *MoneySet             `json:"current_subtotal_price_set"`
	CurrentTotalDiscounts    string                `json:"current_total_discounts"`
	CurrentTotalDiscountsSet *MoneySet             `json:"current_total_discounts_set"`
	CurrentTotalDutiesSet    *MoneySet             `json:"current_total_duties_set"`
	CurrentTotalPrice        string                `json:"current_total_price"`
	CurrentTotalPriceSet     *MoneySet             `json:"current_total_price_set"`
	CurrentTotalTax          string                `json:"current_total_tax"`
	CurrentTotalTaxSet       *MoneySet             `json:"current_total_tax_set"`
	CustomerLocale           *string               `json:"customer_locale"`
	DeviceID                 *int64                `json:"device_id"`
	DiscountCodes            []DiscountCode        `json:"discount_codes"`
	Email                    string                `json:"email"`
	EstimatedTaxes           bool                  `json:"estimated_taxes"`
	FinancialStatus          string                `json:"financial_status"`
	FulfillmentStatus        *string               `json:"fulfillment_status"`
	Gateway                  *string               `json:"gateway"`
	LandingSite              *string               `json:"landing_site"`
	LandingSiteRef           *string               `json:"landing_site_ref"`
	LocationID               *int64                `json:"location_id"`
	MerchantOfRecordAppID    *int64                `json:"merchant_of_record_app_id"`
	Name                     string                `json:"name"`
	Note                     *string               `json:"note"`
	NoteAttributes           []NoteAttribute       `json:"note_attributes"`
	Number                   int                   `json:"number"`
	OrderNumber              int                   `json:"order_number"`
	OrderStatusURL           *string               `json:"order_status_url"`
	OriginalTotalDutiesSet   *MoneySet             `json:"original_total_duties_set"`
	PaymentGatewayNames      []string              `json:"payment_gateway_names"`
	Phone                    *string               `json:"phone"`
	PresentmentCurrency      string                `json:"presentment_currency"`
	ProcessedAt              time.Time             `json:"processed_at"`
	ProcessingMethod         *string               `json:"processing_method"`
	Reference                *string               `json:"reference"`
	ReferringSite            *string               `json:"referring_site"`
	SourceIdentifier         *string               `json:"source_identifier"`
	SourceName               string                `json:"source_name"`
	SourceURL                *string               `json:"source_url"`
	SubtotalPrice            string                `json:"subtotal_price"`
	SubtotalPriceSet         *MoneySet             `json:"subtotal_price_set"`
	Tags                     string                `json:"tags"`
	TaxLines                 []TaxLine             `json:"tax_lines"`
	TaxesIncluded            bool                  `json:"taxes_included"`
	Test                     bool                  `json:"test"`
	Token                    string                `json:"token"`
	TotalDiscounts           string                `json:"total_discounts"`
	TotalDiscountsSet        *MoneySet             `json:"total_discounts_set"`
	TotalLineItemsPrice      string                `json:"total_line_items_price"`
	TotalLineItemsPriceSet   *MoneySet             `json:"total_line_items_price_set"`
	TotalOutstanding         string                `json:"total_outstanding"`
	TotalPrice               string                `json:"total_price"`
	TotalPriceSet            *MoneySet             `json:"total_price_set"`
	TotalPriceUSD            string                `json:"total_price_usd"`
	TotalShippingPriceSet    *MoneySet             `json:"total_shipping_price_set"`
	TotalTax                 string                `json:"total_tax"`
	TotalTaxSet              *MoneySet             `json:"total_tax_set"`
	TotalTipReceived         string                `json:"total_tip_received"`
	TotalWeight              int                   `json:"total_weight"`
	UpdatedAt                time.Time             `json:"updated_at"`
	UserID                   *int64                `json:"user_id"`
	BillingAddress           *Address              `json:"billing_address"`
	Customer                 *Customer             `json:"customer"`
	DiscountApplications     []DiscountApplication `json:"discount_applications"`
	Fulfillments             []Fulfillment         `json:"fulfillments"`
	LineItems                []LineItem            `json:"line_items"`
	PaymentTerms             *PaymentTerm          `json:"payment_terms"`
	Refunds                  []Refund              `json:"refunds"`
	ShippingAddress          *Address              `json:"shipping_address"`
	ShippingLines            []ShippingLine        `json:"shipping_lines"`
}

// MoneySet representa un conjunto de valores monetarios
type MoneySet struct {
	ShopMoney        Money `json:"shop_money"`
	PresentmentMoney Money `json:"presentment_money"`
}

// Money representa un valor monetario
type Money struct {
	Amount       string `json:"amount"`
	CurrencyCode string `json:"currency_code"`
}

// ClientDetails representa los detalles del cliente
type ClientDetails struct {
	AcceptLanguage *string `json:"accept_language"`
	BrowserHeight  *int    `json:"browser_height"`
	BrowserIP      *string `json:"browser_ip"`
	BrowserWidth   *int    `json:"browser_width"`
	SessionHash    *string `json:"session_hash"`
	UserAgent      *string `json:"user_agent"`
}

// DiscountCode representa un código de descuento
type DiscountCode struct {
	Code   string `json:"code"`
	Amount string `json:"amount"`
	Type   string `json:"type"`
}

// NoteAttribute representa un atributo de nota
type NoteAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// TaxLine representa una línea de impuesto
type TaxLine struct {
	Price         string    `json:"price"`
	Rate          float64   `json:"rate"`
	Title         string    `json:"title"`
	PriceSet      *MoneySet `json:"price_set"`
	ChannelLiable bool      `json:"channel_liable"`
}

// Address representa una dirección
type Address struct {
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name"`
	Company      *string  `json:"company"`
	Address1     string   `json:"address1"`
	Address2     *string  `json:"address2"`
	City         string   `json:"city"`
	Province     string   `json:"province"`
	Country      string   `json:"country"`
	Zip          string   `json:"zip"`
	Phone        *string  `json:"phone"`
	Name         string   `json:"name"`
	ProvinceCode *string  `json:"province_code"`
	CountryCode  string   `json:"country_code"`
	CountryName  *string  `json:"country_name"`
	Default      *bool    `json:"default"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
}

// Customer representa un cliente
type Customer struct {
	ID                        int64      `json:"id"`
	Email                     string     `json:"email"`
	AcceptsMarketing          bool       `json:"accepts_marketing"`
	CreatedAt                 time.Time  `json:"created_at"`
	UpdatedAt                 time.Time  `json:"updated_at"`
	FirstName                 string     `json:"first_name"`
	LastName                  string     `json:"last_name"`
	State                     string     `json:"state"`
	Note                      *string    `json:"note"`
	VerifiedEmail             bool       `json:"verified_email"`
	MultipassIdentifier       *string    `json:"multipass_identifier"`
	TaxExempt                 bool       `json:"tax_exempt"`
	Phone                     *string    `json:"phone"`
	Tags                      string     `json:"tags"`
	Currency                  string     `json:"currency"`
	AcceptsMarketingUpdatedAt *time.Time `json:"accepts_marketing_updated_at"`
	MarketingOptInLevel       *string    `json:"marketing_opt_in_level"`
	AdminGraphQLAPIID         string     `json:"admin_graphql_api_id"`
	DefaultAddress            *Address   `json:"default_address"`
	OrdersCount               int        `json:"orders_count"`
	TotalSpent                string     `json:"total_spent"`
}

// DiscountApplication representa una aplicación de descuento
type DiscountApplication struct {
	Type             string `json:"type"`
	Value            string `json:"value"`
	ValueType        string `json:"value_type"`
	AllocationMethod string `json:"allocation_method"`
	TargetSelection  string `json:"target_selection"`
	TargetType       string `json:"target_type"`
	Title            string `json:"title"`
	Description      string `json:"description"`
}

// Fulfillment representa un cumplimiento de orden
type Fulfillment struct {
	ID                int64      `json:"id"`
	OrderID           int64      `json:"order_id"`
	Status            string     `json:"status"`
	CreatedAt         time.Time  `json:"created_at"`
	Service           *string    `json:"service"`
	UpdatedAt         time.Time  `json:"updated_at"`
	TrackingCompany   *string    `json:"tracking_company"`
	ShipmentStatus    *string    `json:"shipment_status"`
	LocationID        *int64     `json:"location_id"`
	OriginAddress     *Address   `json:"origin_address"`
	Receipt           *Receipt   `json:"receipt"`
	Name              string     `json:"name"`
	AdminGraphQLAPIID string     `json:"admin_graphql_api_id"`
	TrackingNumbers   []string   `json:"tracking_numbers"`
	TrackingUrls      []string   `json:"tracking_urls"`
	TrackingNumber    *string    `json:"tracking_number"`
	TrackingURL       *string    `json:"tracking_url"`
	UpdatedAtCustom   *time.Time `json:"updated_at_custom"`
	LineItems         []LineItem `json:"line_items"`
}

// Receipt representa un recibo
type Receipt struct {
	Testcase      bool   `json:"testcase"`
	Authorization string `json:"authorization"`
}

// LineItem representa un item de línea
type LineItem struct {
	ID                         int64                `json:"id"`
	AdminGraphQLAPIID          string               `json:"admin_graphql_api_id"`
	FulfillableQuantity        int                  `json:"fulfillable_quantity"`
	FulfillmentService         *string              `json:"fulfillment_service"`
	FulfillmentStatus          *string              `json:"fulfillment_status"`
	GiftCard                   bool                 `json:"gift_card"`
	Grams                      int                  `json:"grams"`
	Name                       string               `json:"name"`
	Price                      string               `json:"price"`
	PriceSet                   *MoneySet            `json:"price_set"`
	ProductExists              bool                 `json:"product_exists"`
	ProductID                  int64                `json:"product_id"`
	Properties                 []Property           `json:"properties"`
	Quantity                   int                  `json:"quantity"`
	RequiresShipping           bool                 `json:"requires_shipping"`
	SKU                        string               `json:"sku"`
	Taxable                    bool                 `json:"taxable"`
	Title                      string               `json:"title"`
	TotalDiscount              string               `json:"total_discount"`
	TotalDiscountSet           *MoneySet            `json:"total_discount_set"`
	VariantID                  int64                `json:"variant_id"`
	VariantInventoryManagement *string              `json:"variant_inventory_management"`
	VariantTitle               *string              `json:"variant_title"`
	Vendor                     *string              `json:"vendor"`
	TaxLines                   []TaxLine            `json:"tax_lines"`
	Duties                     []Duty               `json:"duties"`
	DiscountAllocations        []DiscountAllocation `json:"discount_allocations"`
}

// Property representa una propiedad de un item
type Property struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Duty representa un deber/impuesto
type Duty struct {
	ID                   *int64    `json:"id"`
	HarmonizedSystemCode *string   `json:"harmonized_system_code"`
	CountryCodeOfOrigin  *string   `json:"country_code_of_origin"`
	AdminGraphQLAPIID    string    `json:"admin_graphql_api_id"`
	Name                 string    `json:"name"`
	Price                string    `json:"price"`
	PriceSet             *MoneySet `json:"price_set"`
	TaxLines             []TaxLine `json:"tax_lines"`
}

// DiscountAllocation representa una asignación de descuento
type DiscountAllocation struct {
	Amount                   string    `json:"amount"`
	AmountSet                *MoneySet `json:"amount_set"`
	DiscountApplicationIndex int       `json:"discount_application_index"`
}

// PaymentTerm representa términos de pago
type PaymentTerm struct {
	ID               *int64  `json:"id"`
	DueInDays        *int    `json:"due_in_days"`
	PaymentTermsName *string `json:"payment_terms_name"`
	PaymentTermsType *string `json:"payment_terms_type"`
}

// Refund representa un reembolso
type Refund struct {
	ID                int64             `json:"id"`
	OrderID           int64             `json:"order_id"`
	CreatedAt         time.Time         `json:"created_at"`
	Note              *string           `json:"note"`
	UserID            *int64            `json:"user_id"`
	ProcessedAt       time.Time         `json:"processed_at"`
	Restock           bool              `json:"restock"`
	AdminGraphQLAPIID string            `json:"admin_graphql_api_id"`
	RefundLineItems   []RefundLineItem  `json:"refund_line_items"`
	Transactions      []Transaction     `json:"transactions"`
	OrderAdjustments  []OrderAdjustment `json:"order_adjustments"`
}

// RefundLineItem representa un item de reembolso
type RefundLineItem struct {
	ID          int64     `json:"id"`
	LineItemID  int64     `json:"line_item_id"`
	LineItem    LineItem  `json:"line_item"`
	Quantity    int       `json:"quantity"`
	Subtotal    float64   `json:"subtotal"`
	SubtotalSet *MoneySet `json:"subtotal_set"`
	TotalTax    float64   `json:"total_tax"`
	TotalTaxSet *MoneySet `json:"total_tax_set"`
	RestockType string    `json:"restock_type"`
}

// Transaction representa una transacción
type Transaction struct {
	ID                int64     `json:"id"`
	OrderID           int64     `json:"order_id"`
	Kind              string    `json:"kind"`
	Gateway           string    `json:"gateway"`
	Status            string    `json:"status"`
	Message           *string   `json:"message"`
	CreatedAt         time.Time `json:"created_at"`
	Test              bool      `json:"test"`
	Authorization     *string   `json:"authorization"`
	LocationID        *int64    `json:"location_id"`
	UserID            *int64    `json:"user_id"`
	ParentID          *int64    `json:"parent_id"`
	ProcessedAt       time.Time `json:"processed_at"`
	DeviceID          *int64    `json:"device_id"`
	ErrorCode         *string   `json:"error_code"`
	SourceName        string    `json:"source_name"`
	Amount            string    `json:"amount"`
	Currency          string    `json:"currency"`
	AdminGraphQLAPIID string    `json:"admin_graphql_api_id"`
}

// OrderAdjustment representa un ajuste de orden
type OrderAdjustment struct {
	ID           int64     `json:"id"`
	OrderID      int64     `json:"order_id"`
	RefundID     int64     `json:"refund_id"`
	Amount       string    `json:"amount"`
	TaxAmount    string    `json:"tax_amount"`
	Kind         string    `json:"kind"`
	Reason       string    `json:"reason"`
	AmountSet    *MoneySet `json:"amount_set"`
	TaxAmountSet *MoneySet `json:"tax_amount_set"`
}

// ShippingLine representa una línea de envío
type ShippingLine struct {
	ID                            int64                `json:"id"`
	CarrierIdentifier             *string              `json:"carrier_identifier"`
	Code                          *string              `json:"code"`
	DeliveryCategory              *string              `json:"delivery_category"`
	DiscountedPrice               string               `json:"discounted_price"`
	DiscountedPriceSet            *MoneySet            `json:"discounted_price_set"`
	Phone                         *string              `json:"phone"`
	Price                         string               `json:"price"`
	PriceSet                      *MoneySet            `json:"price_set"`
	RequestedFulfillmentServiceID *int64               `json:"requested_fulfillment_service_id"`
	Source                        *string              `json:"source"`
	Title                         string               `json:"title"`
	TaxLines                      []TaxLine            `json:"tax_lines"`
	DiscountAllocations           []DiscountAllocation `json:"discount_allocations"`
}
