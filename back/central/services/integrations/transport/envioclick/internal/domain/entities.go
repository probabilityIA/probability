package domain

// EnvioClick API Models (moved from modules/shipments/internal/domain/envioclick_models.go)

type QuoteRequest struct {
	IDRate              int64     `json:"idRate"`
	MyShipmentReference string    `json:"myShipmentReference"`
	ExternalOrderID     string    `json:"external_order_id"`
	OrderUUID           string    `json:"order_uuid,omitempty"`
	RequestPickup       bool      `json:"requestPickup"`
	PickupDate          string    `json:"pickupDate"`
	Insurance           bool      `json:"insurance"`
	Description         string    `json:"description"`
	ContentValue        float64   `json:"contentValue"`
	CODValue            float64   `json:"codValue,omitempty"`
	IncludeGuideCost    bool      `json:"includeGuideCost,omitempty"`
	CODPaymentMethod    string    `json:"codPaymentMethod,omitempty"`
	Packages            []Package `json:"packages"`
	Origin              Address   `json:"origin"`
	Destination         Address   `json:"destination"`
}

type Package struct {
	Weight float64 `json:"weight"`
	Height float64 `json:"height"`
	Width  float64 `json:"width"`
	Length float64 `json:"length"`
}

type Address struct {
	Company     string `json:"company"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	Suburb      string `json:"suburb"`
	CrossStreet string `json:"crossStreet"`
	Reference   string `json:"reference"`
	DaneCode    string `json:"daneCode"`
}

type QuoteResponse struct {
	Status string    `json:"status"`
	Data   QuoteData `json:"data"`
}

type QuoteData struct {
	Rates []Rate `json:"rates"`
}

type Rate struct {
	IDRate           int64   `json:"idRate"`
	IDProduct        int64   `json:"idProduct"`
	Product          string  `json:"product"`
	IDCarrier        int64   `json:"idCarrier"`
	Carrier          string  `json:"carrier"`
	Flete            float64 `json:"flete"`
	MinimumInsurance float64 `json:"minimumInsurance"`
	ExtraInsurance   float64 `json:"extraInsurance"`
	DeliveryDays     int     `json:"deliveryDays"`
	QuotationType    string  `json:"quotationType"`
}

type GenerateResponse struct {
	Status string       `json:"status"`
	Data   GenerateData `json:"data"`
}

type GenerateData struct {
	TrackingNumber   string `json:"tracker"`
	LabelURL         string `json:"url"`
	MyGuideReference string `json:"myGuideReference"`
	Carrier          string `json:"carrier"`
	IDOrder          int64  `json:"idOrder"`
}

type TrackingResponse struct {
	Status         string         `json:"status"`
	StatusCodes    []int          `json:"status_codes"`
	StatusMessages []StatusMsg    `json:"status_messages"`
	Data           TrackingData   `json:"data"`
}

type StatusMsg struct {
	Request string `json:"request"`
}

type TrackingData struct {
	TrackingNumber string         `json:"trackingNumber"`
	Carrier        string         `json:"carrier"`
	Status         string         `json:"status"`
	StatusDetail   string         `json:"statusDetail"`
	ArrivalDate    *string        `json:"arrivalDate"`
	Events         []TrackHistory `json:"history"`
}

type TrackHistory struct {
	Date        string `json:"date"`
	Status      string `json:"status"`
	Description string `json:"description"`
	Location    string `json:"location"`
}

type CancelResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type CancelBatchRequest struct {
	IDOrders []int64 `json:"idOrders"`
}

type CancelBatchResponse struct {
	Status         string          `json:"status"`
	StatusCodes    []int           `json:"status_codes"`
	StatusMessages []StatusMsg     `json:"status_messages"`
	Data           CancelBatchData `json:"data,omitempty"`
}

type CancelBatchData struct {
	NotValidOrders   []int64 `json:"not_valid_orders"`
	OnlyCancelOrders []int64 `json:"only_cancel_orders"`
	ToRefundOrders   []int64 `json:"to_refund_orders"`
}
