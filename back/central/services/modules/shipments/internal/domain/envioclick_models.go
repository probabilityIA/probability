package domain

// EnvioClick Models

type EnvioClickQuoteRequest struct {
	Packages         []EnvioClickPackage `json:"packages"`
	Description      string              `json:"description"`
	ContentValue     float64             `json:"contentValue"`
	CODValue         float64             `json:"codValue,omitempty"`
	IncludeGuideCost bool                `json:"includeGuideCost"`
	CODPaymentMethod string              `json:"codPaymentMethod,omitempty"` // "cash" or "data_phone"
	Origin           EnvioClickAddress   `json:"origin"`
	Destination      EnvioClickAddress   `json:"destination"`
}

type EnvioClickPackage struct {
	Weight float64 `json:"weight"`
	Height float64 `json:"height"`
	Width  float64 `json:"width"`
	Length float64 `json:"length"`
}

type EnvioClickAddress struct {
	DaneCode string `json:"daneCode"`
	Address  string `json:"address"`
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
	TrackingNumber   string `json:"trackingNumber"`
	LabelURL         string `json:"labelUrl"` // Adjust based on actual API response key
	MyGuideReference string `json:"myGuideReference"`
}
