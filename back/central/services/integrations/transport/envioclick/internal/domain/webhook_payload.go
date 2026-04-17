package domain

type WebhookPayload struct {
	IDCarrier           int64           `json:"idCarrier"`
	Carrier             string          `json:"carrier"`
	IDOrder             int64           `json:"idOrder"`
	TrackingCode        string          `json:"trackingCode"`
	RealPickupDate      string          `json:"realPickupDate"`
	ArrivalDate         string          `json:"arrivalDate"`
	RealDeliveryDate    string          `json:"realDeliveryDate"`
	MyShipmentReference string          `json:"myShipmentReference"`
	Events              []WebhookEvent  `json:"events"`
}

type WebhookEvent struct {
	Timestamp     string  `json:"timestamp"`
	Status        string  `json:"status"`
	StatusDetail  string  `json:"statusDetail"`
	StatusStep    string  `json:"statusStep"`
	Incidence     bool    `json:"incidence"`
	IncidenceType *string `json:"incidenceType"`
	Description   string  `json:"description"`
	ReceivedBy    *string `json:"receivedBy"`
}

func (p *WebhookPayload) LatestEvent() *WebhookEvent {
	if len(p.Events) == 0 {
		return nil
	}
	return &p.Events[len(p.Events)-1]
}

type NormalizedWebhookUpdate struct {
	TrackingNumber      string
	MyShipmentReference string
	ProbabilityStatus   ProbabilityShipmentStatus
	RawStatusStep       string
	HasIncidence        bool
	IsUnknownStatus     bool
	EventDescription    string
	EventTimestamp      string
	ShippedAt           *string
	DeliveredAt         *string
}
