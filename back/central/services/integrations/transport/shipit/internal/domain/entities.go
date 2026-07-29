package domain

type Credentials struct {
	Email       string
	AccessToken string
}

type GuideRequest struct {
	MyShipmentReference string    `json:"myShipmentReference"`
	ExternalOrderID     string    `json:"external_order_id,omitempty"`
	OrderUUID           string    `json:"order_uuid,omitempty"`
	Description         string    `json:"description"`
	ContentValue        float64   `json:"contentValue"`
	CODValue            float64   `json:"codValue,omitempty"`
	Carrier             string    `json:"carrier,omitempty"`
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
	Company     string `json:"company,omitempty"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	Suburb      string `json:"suburb,omitempty"`
	City        string `json:"city,omitempty"`
	State       string `json:"state,omitempty"`
	Reference   string `json:"reference,omitempty"`
	CommuneID   uint   `json:"communeId,omitempty"`
	CommuneName string `json:"communeName,omitempty"`
}

type ShipmentRequest struct {
	Shipment ShipmentBody `json:"shipment"`
}

type ShipmentBody struct {
	Kind      int         `json:"kind"`
	Platform  int         `json:"platform"`
	Reference string      `json:"reference"`
	Items     int         `json:"items"`
	Sandbox   bool        `json:"sandbox"`
	Sizes     Sizes       `json:"sizes"`
	Destiny   Destiny     `json:"destiny"`
	Courier   CourierSpec `json:"courier"`
	Insurance Insurance   `json:"insurance"`
}

type Sizes struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Length float64 `json:"length"`
	Weight float64 `json:"weight"`
}

type Destiny struct {
	Street      string `json:"street"`
	Number      string `json:"number"`
	Complement  string `json:"complement,omitempty"`
	CommuneID   uint   `json:"commune_id,omitempty"`
	CommuneName string `json:"commune_name"`
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Kind        string `json:"kind"`
}

type CourierSpec struct {
	ID        uint   `json:"id,omitempty"`
	Client    string `json:"client"`
	Selected  bool   `json:"selected"`
	Payable   bool   `json:"payable"`
	Algorithm int    `json:"algorithm"`
}

type Insurance struct {
	TicketNumber string  `json:"ticket_number"`
	TicketAmount float64 `json:"ticket_amount"`
	Detail       string  `json:"detail"`
	Extra        bool    `json:"extra"`
}

type ShipmentResponse struct {
	ID               int64   `json:"id"`
	Reference        string  `json:"reference"`
	Status           string  `json:"status"`
	SubStatus        string  `json:"sub_status"`
	TrackingNumber   *string `json:"tracking_number"`
	ShipitCode       string  `json:"shipit_code"`
	TrackingLink     string  `json:"shipit_tracking_link"`
	CourierURL       string  `json:"courier_url"`
	CourierForClient string  `json:"courier_for_client"`
	TicketPDFURL     *string `json:"ticket_shipit_pdf_url"`
	TicketURL        *string `json:"ticket_url"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

type RatesRequest struct {
	Parcel Parcel `json:"parcel"`
}

type Parcel struct {
	Length        float64 `json:"length"`
	Width         float64 `json:"width"`
	Height        float64 `json:"height"`
	Weight        float64 `json:"weight"`
	OriginID      uint    `json:"origin_id"`
	DestinyID     uint    `json:"destiny_id"`
	TypeOfDestiny string  `json:"type_of_destiny"`
	Algorithm     int     `json:"algorithm,omitempty"`
}

type RatesResponse struct {
	Algorithm string  `json:"algorithm"`
	Prices    []Price `json:"prices"`
}

type Price struct {
	Courier             PriceCourier `json:"courier"`
	Name                string       `json:"name"`
	Price               float64      `json:"price"`
	Days                int          `json:"days"`
	AvailableToShipping bool         `json:"available_to_shipping"`
}

type PriceCourier struct {
	Name string `json:"name"`
}

type TrackingResponse struct {
	Data TrackingData `json:"data"`
}

type TrackingData struct {
	ID          int64            `json:"id"`
	Courier     string           `json:"courier"`
	Number      string           `json:"number"`
	PackageID   int64            `json:"package_id"`
	IsDelivered bool             `json:"is_delivered"`
	CreatedAt   string           `json:"created_at"`
	Statuses    []TrackingStatus `json:"statuses"`
}

type TrackingStatus struct {
	ID               int64  `json:"id"`
	TrackingID       int64  `json:"tracking_id"`
	GenericStatus    int    `json:"generic_status"`
	CourierStatus    string `json:"courier_status"`
	CourierUpdatedAt string `json:"courier_updated_at"`
	CreatedAt        string `json:"created_at"`
}

type Commune struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type QuoteResponse struct {
	Status string    `json:"status"`
	Data   QuoteData `json:"data"`
}

type QuoteData struct {
	Rates []Rate `json:"rates"`
}

type Rate struct {
	IDRate       int64   `json:"idRate"`
	Product      string  `json:"product"`
	Carrier      string  `json:"carrier"`
	Flete        float64 `json:"flete"`
	DeliveryDays int     `json:"deliveryDays"`
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

type TrackOutput struct {
	Status         string       `json:"status"`
	Carrier        string       `json:"carrier"`
	TrackingNumber string       `json:"trackingNumber"`
	StatusDetail   string       `json:"statusDetail"`
	IsDelivered    bool         `json:"isDelivered"`
	History        []TrackEvent `json:"history"`
}

type TrackEvent struct {
	Date        string `json:"date"`
	Status      string `json:"status"`
	Description string `json:"description"`
	Location    string `json:"location"`
}
