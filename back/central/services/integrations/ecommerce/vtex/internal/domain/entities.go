package domain

import "time"

type Integration struct {
	ID              uint
	BusinessID      *uint
	Name            string
	StoreID         string
	IntegrationType int
	Config          map[string]interface{}
}

type VTEXOrder struct {
	OrderID            string
	Sequence           string
	MarketplaceOrderID string
	Status             string
	StatusDescription  string
	Value              int
	TotalItems         int
	TotalDiscount      int
	TotalFreight       int
	CreationDate       time.Time
	LastChange         time.Time
	Currency           string
	Items              []VTEXOrderItem
	ShippingData       *VTEXShippingData
	PaymentData        *VTEXPaymentData
	ClientProfileData  *VTEXClientProfile
	RatesAndBenefits   []VTEXRateAndBenefit
	Sellers            []VTEXSeller
	Totals             []VTEXTotal
	PackageAttachment  *VTEXPackageAttachment
}

type VTEXOrderItem struct {
	UniqueID        string
	ID              string
	ProductID       string
	EANID           string
	RefID           string
	Name            string
	SKUName         string
	ImageURL        string
	DetailURL       string
	Quantity        int
	Price           int
	ListPrice       int
	SellingPrice    int
	Tax             int
	MeasurementUnit string
	UnitMultiplier  float64
}

type VTEXShippingData struct {
	Address       *VTEXAddress
	LogisticsInfo []VTEXLogisticsInfo
	SelectedSLA   string
	TrackingHints []VTEXTrackingHint
}

type VTEXAddress struct {
	AddressType    string
	ReceiverName   string
	Street         string
	Number         string
	Complement     string
	Neighborhood   string
	City           string
	State          string
	Country        string
	PostalCode     string
	Reference      string
	GeoCoordinates []float64
}

type VTEXLogisticsInfo struct {
	ItemIndex            int
	SelectedSLA          string
	LockTTL              string
	Price                int
	ListPrice            int
	SellingPrice         int
	DeliveryWindow       *VTEXDeliveryWindow
	ShippingEstimate     string
	ShippingEstimateDate *time.Time
	DeliveryCompany      string
	DeliveryIDs          []VTEXDeliveryID
}

type VTEXDeliveryWindow struct {
	StartDateUTC time.Time
	EndDateUTC   time.Time
	Price        int
}

type VTEXDeliveryID struct {
	CourierID   string
	CourierName string
	DockID      string
	Quantity    int
	WarehouseID string
}

type VTEXTrackingHint struct {
	CourierName   string
	TrackingID    string
	TrackingURL   string
	TrackingLabel string
}

type VTEXPaymentData struct {
	Transactions []VTEXTransaction
}

type VTEXTransaction struct {
	IsActive      bool
	TransactionID string
	MerchantName  string
	Payments      []VTEXPayment
}

type VTEXPayment struct {
	ID                 string
	PaymentSystem      string
	PaymentSystemName  string
	Value              int
	ReferenceValue     int
	Group              string
	ConnectorResponses map[string]string
	InstallmentCount   int
	CardHolder         string
	FirstDigits        string
	LastDigits         string
	URL                string
	TID                string
}

type VTEXClientProfile struct {
	Email         string
	FirstName     string
	LastName      string
	DocumentType  string
	Document      string
	Phone         string
	CorporateName string
	IsCorporate   bool
}

type VTEXRateAndBenefit struct {
	ID   string
	Name string
}

type VTEXSeller struct {
	ID          string
	Name        string
	SubSellerID string
}

type VTEXTotal struct {
	ID    string
	Name  string
	Value int
}

type VTEXPackageAttachment struct {
	Packages []VTEXPackage
}

type VTEXPackage struct {
	Items          []VTEXPackageItem
	CourierStatus  *VTEXCourierStatus
	TrackingNumber string
	TrackingURL    string
	InvoiceNumber  string
	InvoiceValue   int
	InvoiceURL     string
	InvoiceKey     string
	Courier        string
	Type           string
}

type VTEXPackageItem struct {
	ItemIndex   int
	Quantity    int
	Price       int
	Description string
}

type VTEXCourierStatus struct {
	Status   string
	Finished bool
	Data     []VTEXCourierEvent
}

type VTEXCourierEvent struct {
	Description string
	Date        string
	City        string
	State       string
}

type VTEXOrderListResponse struct {
	List   []VTEXOrderSummary
	Paging VTEXPaging
}

type VTEXOrderSummary struct {
	OrderID      string
	Sequence     string
	Status       string
	CreationDate time.Time
	LastChange   time.Time
	TotalValue   int
	CurrencyCode string
	Origin       string
}

type VTEXPaging struct {
	Total       int
	Pages       int
	CurrentPage int
	PerPage     int
}

type VTEXWebhookPayload struct {
	Domain        string
	OrderID       string
	State         string
	LastState     string
	LastChange    string
	CurrentChange string
	Origin        *VTEXWebhookOrigin
	IntegrationID string
}

type VTEXWebhookOrigin struct {
	Account string
	Key     string
}
