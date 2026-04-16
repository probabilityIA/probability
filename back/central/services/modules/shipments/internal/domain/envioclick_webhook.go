package domain

type EnvioClickWebhookPayload struct {
	IDCarrier           int64                    `json:"idCarrier"`
	Carrier             string                   `json:"carrier"`
	IDOrder             int64                    `json:"idOrder"`
	TrackingCode        string                   `json:"trackingCode"`
	RealPickupDate      string                   `json:"realPickupDate"`
	ArrivalDate         string                   `json:"arrivalDate"`
	RealDeliveryDate    string                   `json:"realDeliveryDate"`
	MyShipmentReference string                   `json:"myShipmentReference"`
	Events              []EnvioClickWebhookEvent `json:"events"`
}

type EnvioClickWebhookEvent struct {
	Timestamp     string  `json:"timestamp"`
	Status        string  `json:"status"`
	StatusDetail  string  `json:"statusDetail"`
	StatusStep    string  `json:"statusStep"`
	Incidence     bool    `json:"incidence"`
	IncidenceType *string `json:"incidenceType"`
	Description   string  `json:"description"`
	ReceivedBy    *string `json:"receivedBy"`
}

func MapEnvioClickStatus(statusStep string, incidence bool) string {
	if incidence {
		return "failed"
	}
	switch statusStep {
	case "Pendiente de Recolección", "Pendiente":
		return "pending"
	case "En tránsito", "En Tránsito", "En Transito", "En transito",
		"Envío Recolectado", "Envio Recolectado", "En Distribución", "En distribucion":
		return "in_transit"
	case "Entregado", "Entregada":
		return "delivered"
	case "Incidencia", "Novedad", "No entregado", "No Entregado":
		return "failed"
	default:
		return "in_transit"
	}
}
