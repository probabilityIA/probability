package domain

// EnvioClickWebhookPayload representa el payload recibido desde el webhook de EnvioClick
// cuando ocurre un evento de tracking de un envío.
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

// EnvioClickWebhookEvent representa cada evento de tracking dentro del webhook.
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

// MapEnvioClickStatus convierte un estado de EnvioClick al estado interno de Probability.
// Referencia de estados EnvioClick:
//   - "Pendiente de Recolección"  -> pending
//   - "En tránsito" / "En Tránsito" -> in_transit
//   - "Entregado"                 -> delivered
//   - "Incidencia" / incidence=true -> failed
//   - Cualquier otro              -> in_transit (por defecto seguro)
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
