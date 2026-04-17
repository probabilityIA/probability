package domain

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func normalize(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, err := transform.String(t, s)
	if err != nil {
		return strings.ToLower(strings.TrimSpace(s))
	}
	return strings.ToLower(strings.TrimSpace(result))
}

func ApiStatusToStep(status, statusDetail string) string {
	combined := normalize(status + " " + statusDetail)
	switch {
	case strings.Contains(combined, "entregado"), strings.Contains(combined, "entregada"):
		return "Entregado"
	case strings.Contains(combined, "no entregad"):
		return "No Entregado"
	case strings.Contains(combined, "devuelto"), strings.Contains(combined, "regresado"):
		return "Devuelto"
	case strings.Contains(combined, "cancelad"):
		return "Cancelado"
	case strings.Contains(combined, "pendiente"):
		return "Pendiente"
	case strings.Contains(combined, "novedad"),
		strings.Contains(combined, "incidencia"),
		strings.Contains(combined, "ausente"):
		return "Novedad"
	case strings.Contains(combined, "distribucion"), strings.Contains(combined, "reparto"):
		return "En Distribucion"
	case strings.Contains(combined, "recolec"):
		return "Envio Recolectado"
	case strings.Contains(combined, "transit"):
		return "En Transito"
	}
	return status
}

func MapStatusStepToProbability(statusStep string, incidence bool) (status ProbabilityShipmentStatus, unknown bool) {
	normalized := normalize(statusStep)

	if incidence {
		if normalized == "entregado" || normalized == "entregada" {
			return StatusDelivered, false
		}
		return StatusOnHold, false
	}

	switch normalized {
	case "pendiente", "pendiente de recoleccion":
		return StatusPending, false
	case "envio recolectado", "recolectado":
		return StatusPickedUp, false
	case "en transito":
		return StatusInTransit, false
	case "en distribucion":
		return StatusOutForDelivery, false
	case "entregado", "entregada":
		return StatusDelivered, false
	case "novedad", "incidencia":
		return StatusOnHold, false
	case "no entregado", "no entregada":
		return StatusFailed, false
	case "devuelto", "devuelta", "regresado", "regresada":
		return StatusReturned, false
	case "cancelado", "cancelada":
		return StatusCancelled, false
	}

	return StatusInTransit, true
}

func (p *WebhookPayload) ToNormalizedUpdate() *NormalizedWebhookUpdate {
	event := p.LatestEvent()
	if event == nil {
		return nil
	}

	status, unknown := MapStatusStepToProbability(event.StatusStep, event.Incidence)

	update := &NormalizedWebhookUpdate{
		TrackingNumber:      p.TrackingCode,
		MyShipmentReference: p.MyShipmentReference,
		ProbabilityStatus:   status,
		RawStatusStep:       event.StatusStep,
		HasIncidence:        event.Incidence,
		IsUnknownStatus:     unknown,
		EventDescription:    event.Description,
		EventTimestamp:      event.Timestamp,
	}

	if p.RealDeliveryDate != "" && status == StatusDelivered {
		update.DeliveredAt = &p.RealDeliveryDate
	}

	if p.RealPickupDate != "" {
		switch status {
		case StatusPickedUp, StatusInTransit, StatusOutForDelivery, StatusDelivered:
			update.ShippedAt = &p.RealPickupDate
		}
	}

	return update
}
