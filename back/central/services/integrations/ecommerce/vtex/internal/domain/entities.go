package domain

import "time"

// Integration representa los datos de una integración de VTEX
// tal como se obtienen del core de integraciones.
type Integration struct {
	ID              uint
	BusinessID      *uint
	Name            string
	StoreID         string
	IntegrationType int
	Config          map[string]interface{}
}

// VTEXOrder representa una orden de VTEX (OMS API).
// Los valores monetarios están en centavos (dividir por 100 para obtener el valor real).
type VTEXOrder struct {
	OrderID            string
	Sequence           string
	MarketplaceOrderID string
	Status             string
	StatusDescription  string
	Value              int // centavos
	TotalItems         int // centavos
	TotalDiscount      int // centavos
	TotalFreight       int // centavos
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

// VTEXOrderItem representa un item de la orden de VTEX.
type VTEXOrderItem struct {
	UniqueID       string
	ID             string
	ProductID      string
	EANID          string
	RefID          string
	Name           string
	SKUName        string
	ImageURL       string
	DetailURL      string
	Quantity       int
	Price          int // centavos
	ListPrice      int // centavos
	SellingPrice   int // centavos
	Tax            int // centavos
	MeasurementUnit string
	UnitMultiplier float64
}

// VTEXShippingData datos de envío de la orden.
type VTEXShippingData struct {
	Address          *VTEXAddress
	LogisticsInfo    []VTEXLogisticsInfo
	SelectedSLA      string
	TrackingHints    []VTEXTrackingHint
}

// VTEXAddress dirección de envío.
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

// VTEXLogisticsInfo información logística por item.
type VTEXLogisticsInfo struct {
	ItemIndex            int
	SelectedSLA          string
	LockTTL              string
	Price                int // centavos
	ListPrice            int // centavos
	SellingPrice         int // centavos
	DeliveryWindow       *VTEXDeliveryWindow
	ShippingEstimate     string
	ShippingEstimateDate *time.Time
	DeliveryCompany      string
	DeliveryIDs          []VTEXDeliveryID
}

// VTEXDeliveryWindow ventana de entrega.
type VTEXDeliveryWindow struct {
	StartDateUTC time.Time
	EndDateUTC   time.Time
	Price        int // centavos
}

// VTEXDeliveryID identificador de entrega.
type VTEXDeliveryID struct {
	CourierID   string
	CourierName string
	DockID      string
	Quantity    int
	WarehouseID string
}

// VTEXTrackingHint información de rastreo.
type VTEXTrackingHint struct {
	CourierName    string
	TrackingID     string
	TrackingURL    string
	TrackingLabel  string
}

// VTEXPaymentData datos de pago de la orden.
type VTEXPaymentData struct {
	Transactions []VTEXTransaction
}

// VTEXTransaction representa una transacción de pago.
type VTEXTransaction struct {
	IsActive       bool
	TransactionID  string
	MerchantName   string
	Payments       []VTEXPayment
}

// VTEXPayment representa un pago individual.
type VTEXPayment struct {
	ID                  string
	PaymentSystem       string
	PaymentSystemName   string
	Value               int // centavos
	ReferenceValue      int // centavos
	Group               string
	ConnectorResponses  map[string]string
	InstallmentCount    int
	CardHolder          string
	FirstDigits         string
	LastDigits          string
	URL                 string
	TID                 string
}

// VTEXClientProfile datos del cliente que realizó el pedido.
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

// VTEXRateAndBenefit promociones y descuentos aplicados.
type VTEXRateAndBenefit struct {
	ID   string
	Name string
}

// VTEXSeller vendedor asociado a la orden.
type VTEXSeller struct {
	ID         string
	Name       string
	SubSellerID string
}

// VTEXTotal totalizador de la orden.
type VTEXTotal struct {
	ID    string
	Name  string
	Value int // centavos
}

// VTEXPackageAttachment paquetes de envío.
type VTEXPackageAttachment struct {
	Packages []VTEXPackage
}

// VTEXPackage paquete individual con datos de rastreo.
type VTEXPackage struct {
	Items            []VTEXPackageItem
	CourierStatus    *VTEXCourierStatus
	TrackingNumber   string
	TrackingURL      string
	InvoiceNumber    string
	InvoiceValue     int // centavos
	InvoiceURL       string
	InvoiceKey       string
	Courier          string
	Type             string
}

// VTEXPackageItem item dentro de un paquete.
type VTEXPackageItem struct {
	ItemIndex   int
	Quantity    int
	Price       int // centavos
	Description string
}

// VTEXCourierStatus estado del courier.
type VTEXCourierStatus struct {
	Status     string
	Finished   bool
	Data       []VTEXCourierEvent
}

// VTEXCourierEvent evento de rastreo del courier.
type VTEXCourierEvent struct {
	Description string
	Date        string
	City        string
	State       string
}

// VTEXOrderListResponse respuesta de GET /api/oms/pvt/orders.
type VTEXOrderListResponse struct {
	List    []VTEXOrderSummary
	Paging  VTEXPaging
}

// VTEXOrderSummary resumen de orden en la lista (no tiene detalle completo).
type VTEXOrderSummary struct {
	OrderID       string
	Sequence      string
	Status        string
	CreationDate  time.Time
	LastChange    time.Time
	TotalValue    int // centavos
	CurrencyCode  string
	Origin        string
}

// VTEXPaging datos de paginación de VTEX.
type VTEXPaging struct {
	Total   int
	Pages   int
	CurrentPage int
	PerPage int
}

// VTEXWebhookPayload es el payload que VTEX envía al webhook (Hook v1).
type VTEXWebhookPayload struct {
	Domain        string
	OrderID       string
	State         string
	LastState     string
	LastChange    string
	CurrentChange string
	Origin        *VTEXWebhookOrigin
}

// VTEXWebhookOrigin origen del webhook.
type VTEXWebhookOrigin struct {
	Account string
	Key     string
}
