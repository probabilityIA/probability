package domain

// EnvioClick Models

type EnvioClickQuoteRequest struct {
	IDRate              int64               `json:"idRate"`
	MyShipmentReference string              `json:"myShipmentReference"`
	ExternalOrderID     string              `json:"external_order_id"`
	RequestPickup       bool                `json:"requestPickup"`
	PickupDate          string              `json:"pickupDate"`
	Insurance           bool                `json:"insurance"`
	Description         string              `json:"description"`
	ContentValue        float64             `json:"contentValue"`
	CODValue            float64             `json:"codValue"`
	IncludeGuideCost    bool                `json:"includeGuideCost"`
	CODPaymentMethod    string              `json:"codPaymentMethod"` // "cash" or "data_phone"
	Packages            []EnvioClickPackage `json:"packages"`
	Origin              EnvioClickAddress   `json:"origin"`
	Destination         EnvioClickAddress   `json:"destination"`
}

type EnvioClickPackage struct {
	Weight float64 `json:"weight"`
	Height float64 `json:"height"`
	Width  float64 `json:"width"`
	Length float64 `json:"length"`
}

type EnvioClickAddress struct {
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

type EnvioClickQuoteResponse struct {
	Status string              `json:"status"`
	Data   EnvioClickQuoteData `json:"data"`
}

type EnvioClickQuoteData struct {
	Rates []EnvioClickRate `json:"rates"`
}

type EnvioClickRate struct {
	IDRate        int64   `json:"idRate"`
	IDProduct     int64   `json:"idProduct"`
	Product       string  `json:"product"`
	IDCarrier     int64   `json:"idCarrier"`
	Carrier       string  `json:"carrier"`
	Flete         float64 `json:"flete"`
	DeliveryDays  int     `json:"deliveryDays"`
	QuotationType string  `json:"quotationType"`
}

type EnvioClickGenerateResponse struct {
	Status string                 `json:"status"`
	Data   EnvioClickGenerateData `json:"data"`
}

type EnvioClickGenerateData struct {
	TrackingNumber   string `json:"tracker"`
	LabelURL         string `json:"url"` // Matches API response key "url"
	MyGuideReference string `json:"myGuideReference"`
}
