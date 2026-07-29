package domain

import "sync"

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
	Complement  string `json:"complement"`
	CommuneID   uint   `json:"commune_id"`
	CommuneName string `json:"commune_name"`
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Kind        string `json:"kind"`
}

type CourierSpec struct {
	ID        uint   `json:"id"`
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
	Algorithm     int     `json:"algorithm"`
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

type StoredShipment struct {
	ID             int64
	Reference      string
	TrackingNumber string
	Courier        string
	Status         string
	Sandbox        bool
	Destiny        Destiny
	Statuses       []TrackingStatus
	CreatedAt      string
	UpdatedAt      string
}

type Repository struct {
	mu        sync.RWMutex
	shipments map[string]*StoredShipment
	nextID    int64
}

func NewRepository() *Repository {
	return &Repository{
		shipments: make(map[string]*StoredShipment),
		nextID:    5000,
	}
}

func (r *Repository) NextID() int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nextID++
	return r.nextID
}

func (r *Repository) Save(s *StoredShipment) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.shipments[s.TrackingNumber] = s
}

func (r *Repository) GetByTracking(trackingNumber string) *StoredShipment {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.shipments[trackingNumber]
}

func (r *Repository) GetAll() []*StoredShipment {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*StoredShipment, 0, len(r.shipments))
	for _, s := range r.shipments {
		result = append(result, s)
	}
	return result
}
