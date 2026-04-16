package domain

import "time"

// Integration representa los datos de una integración de MercadoLibre
// tal como se obtienen del core de integraciones.
type Integration struct {
	ID              uint
	BusinessID      *uint
	Name            string
	StoreID         string
	IntegrationType int
	Config          map[string]interface{}
}

// MeliOrder representa una orden de MercadoLibre (API v2 orders).
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

// MeliStatusDetail detalle del estado de la orden.
type MeliStatusDetail struct {
	Description string
	Code        string
}

// MeliBuyer representa al comprador en una orden.
type MeliBuyer struct {
	ID          int64
	Nickname    string
	FirstName   string
	LastName    string
	Email       string
	Phone       MeliPhone
	BillingInfo *MeliBillingInfo
}

// MeliPhone teléfono del comprador.
type MeliPhone struct {
	AreaCode  string
	Number    string
	Extension string
}

// MeliBillingInfo información de facturación del comprador.
type MeliBillingInfo struct {
	DocType   string
	DocNumber string
}

// MeliSeller representa al vendedor.
type MeliSeller struct {
	ID       int64
	Nickname string
}

// MeliOrderItem representa un item de la orden.
type MeliOrderItem struct {
	Item          MeliItem
	Quantity      int
	UnitPrice     float64
	FullUnitPrice float64
	Currency      string
	SaleFee       float64
}

// MeliItem representa un producto/publicación de MeLi.
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

// MeliVariationAttribute atributo de variación (color, talla, etc).
type MeliVariationAttribute struct {
	ID    string
	Name  string
	Value string
}

// MeliPayment representa un pago de la orden.
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

// MeliShippingRef referencia al envío dentro de la orden (solo ID).
// Para obtener detalles completos se debe llamar a GET /shipments/{id}.
type MeliShippingRef struct {
	ID int64
}

// MeliShippingDetail datos completos del envío obtenidos de GET /shipments/{id}.
type MeliShippingDetail struct {
	ID              int64
	Status          string
	SubStatus       string
	ShipmentType    string
	DateCreated     *time.Time
	ReceiverAddress *MeliReceiverAddress
	SenderAddress   *MeliSenderAddress
	ShippingOption  *MeliShippingOption
}

// MeliReceiverAddress dirección del comprador.
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
}

// MeliSenderAddress dirección del vendedor.
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

// MeliLocation representa una ubicación con ID y nombre.
type MeliLocation struct {
	ID   string
	Name string
}

// MeliShippingOption opción de envío seleccionada.
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

// MeliEstimatedDelivery tiempo estimado de entrega.
type MeliEstimatedDelivery struct {
	Type string
	Date *time.Time
}

// MeliNotification representa la notificación IPN que MercadoLibre envía.
type MeliNotification struct {
	Resource      string // "/orders/123456789"
	UserID        int64
	Topic         string // "orders_v2", "payments", "items"
	ApplicationID int64
	Attempts      int
	Sent          time.Time
	Received      time.Time
}

// TokenResponse representa la respuesta del endpoint OAuth /oauth/token.
type TokenResponse struct {
	AccessToken  string
	TokenType    string
	ExpiresIn    int
	Scope        string
	UserID       int64
	RefreshToken string
}
