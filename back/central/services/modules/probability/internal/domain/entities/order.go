package entities

// ScoreOrder contiene los datos necesarios para calcular el score de entrega
type ScoreOrder struct {
	ID                 string
	BusinessID         *uint
	IntegrationID      uint
	CustomerID         *uint
	CustomerEmail      string
	CustomerName       string
	Platform           string
	CustomerPhone      string
	ShippingStreet     string
	Address2           string
	CustomerOrderCount int
	OrderNumber        string
	CodTotal           *float64
	Metadata           []byte
	PaymentDetails     []byte
	Payments           []ScorePayment
	Addresses          []ScoreAddress
	ChannelMetadata    []ScoreChannelMetadata
	TotalAmount        float64
	IsPaid             bool
	IsConfirmed        *bool
	Coupon             *string
	Weight             *float64
	PaymentMethodID    uint
	OrderItemCount     int
	CustomerHistory    *CustomerHistory
	DeliveryHistory    *DeliveryHistory
}

// ScorePayment contiene los datos de pago relevantes para scoring
type ScorePayment struct {
	Gateway *string
}

// ScoreAddress contiene los datos de direccion relevantes para scoring
type ScoreAddress struct {
	Type    string
	Street2 string
}

// ScoreChannelMetadata contiene los metadatos del canal relevantes para scoring
type ScoreChannelMetadata struct {
	RawData []byte
}
