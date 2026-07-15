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

type MeliOrder struct {
	ID                      int64
	Status                  string
	StatusDetail            *MeliStatusDetail
	DateCreated             time.Time
	DateClosed              *time.Time
	LastUpdated             time.Time
	TotalAmount             float64
	CurrencyID              string
	Buyer                   MeliBuyer
	Seller                  MeliSeller
	OrderItems              []MeliOrderItem
	Payments                []MeliPayment
	Shipping                *MeliShippingRef
	Tags                    []string
	PackID                  *int64
	CouponAmount            float64
	CouponID                *string
	ManufacturingEndingDate *time.Time
}

type MeliStatusDetail struct {
	Description string
	Code        string
}

type MeliBuyer struct {
	ID          int64
	Nickname    string
	FirstName   string
	LastName    string
	Email       string
	Phone       MeliPhone
	BillingInfo *MeliBillingInfo
}

type MeliPhone struct {
	AreaCode  string
	Number    string
	Extension string
}

type MeliBillingInfo struct {
	DocType   string
	DocNumber string
}

type MeliSeller struct {
	ID       int64
	Nickname string
}

type MeliOrderItem struct {
	Item          MeliItem
	Quantity      int
	UnitPrice     float64
	FullUnitPrice float64
	Currency      string
	SaleFee       float64
}

type MeliItem struct {
	ID                  string
	Title               string
	CategoryID          string
	VariationID         *int64
	SellerCustomField   *string
	SellerSKU           *string
	Condition           string
	VariationAttributes []MeliVariationAttribute
}

type MeliVariationAttribute struct {
	ID    string
	Name  string
	Value string
}

type MeliPayment struct {
	ID                 int64
	OrderID            int64
	PayerID            int64
	Status             string
	StatusDetail       string
	TransactionAmount  float64
	CurrencyID         string
	DateCreated        time.Time
	DateApproved       *time.Time
	DateLastModified   *time.Time
	PaymentMethodID    string
	PaymentType        string
	OperationType      string
	InstallmentAmount  *float64
	Installments       int
	TransactionOrderID *string
}

type MeliShippingRef struct {
	ID int64
}

type MeliShippingDetail struct {
	ID              int64
	Status          string
	SubStatus       string
	ShipmentType    string
	LogisticType    string
	LogisticMode    string
	TrackingNumber  string
	TrackingMethod  string
	DateCreated     *time.Time
	ReceiverAddress *MeliReceiverAddress
	SenderAddress   *MeliSenderAddress
	ShippingOption  *MeliShippingOption
}

type MeliReceiverAddress struct {
	ID           int64
	AddressLine  string
	StreetName   string
	StreetNumber string
	ZipCode      string
	City         MeliLocation
	State        MeliLocation
	Country      MeliLocation
	Neighborhood *MeliLocation
	Latitude     *float64
	Longitude    *float64
	Comment      string
	ReceiverName string
	ReceiverPhone string
}

type MeliSenderAddress struct {
	ID           int64
	AddressLine  string
	StreetName   string
	StreetNumber string
	ZipCode      string
	City         MeliLocation
	State        MeliLocation
	Country      MeliLocation
}

type MeliLocation struct {
	ID   string
	Name string
}

type MeliShippingOption struct {
	ID                    int64
	Name                  string
	CurrencyID            string
	Cost                  float64
	ListCost              float64
	ShippingMethodID      int64
	DeliveryType          string
	EstimatedDeliveryTime *MeliEstimatedDelivery
}

type MeliEstimatedDelivery struct {
	Type string
	Date *time.Time
}

type MeliNotification struct {
	Resource      string
	UserID        int64
	Topic         string
	ApplicationID int64
	Attempts      int
	Sent          time.Time
	Received      time.Time
}

type TokenResponse struct {
	AccessToken  string
	TokenType    string
	ExpiresIn    int
	Scope        string
	UserID       int64
	RefreshToken string
}
