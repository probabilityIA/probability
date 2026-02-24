package domain

import "time"

// ───────────────────────────────────────────
//
//	SHIPMENT EVENT TYPES
//
// ───────────────────────────────────────────

// ShipmentEventType define los tipos de eventos de envíos
type ShipmentEventType string

const (
	ShipmentEventTypeQuoteReceived   ShipmentEventType = "shipment.quote_received"
	ShipmentEventTypeQuoteFailed     ShipmentEventType = "shipment.quote_failed"
	ShipmentEventTypeGuideGenerated  ShipmentEventType = "shipment.guide_generated"
	ShipmentEventTypeGuideFailed     ShipmentEventType = "shipment.guide_failed"
	ShipmentEventTypeTrackingUpdated ShipmentEventType = "shipment.tracking_updated"
	ShipmentEventTypeTrackingFailed  ShipmentEventType = "shipment.tracking_failed"
	ShipmentEventTypeCancelled       ShipmentEventType = "shipment.cancelled"
	ShipmentEventTypeCancelFailed    ShipmentEventType = "shipment.cancel_failed"
)

// IsValid verifica si el tipo de evento es válido
func (t ShipmentEventType) IsValid() bool {
	switch t {
	case ShipmentEventTypeQuoteReceived, ShipmentEventTypeQuoteFailed,
		ShipmentEventTypeGuideGenerated, ShipmentEventTypeGuideFailed,
		ShipmentEventTypeTrackingUpdated, ShipmentEventTypeTrackingFailed,
		ShipmentEventTypeCancelled, ShipmentEventTypeCancelFailed:
		return true
	}
	return false
}

// ───────────────────────────────────────────
//
//	SHIPMENT EVENT STRUCTURES
//
// ───────────────────────────────────────────

// ShipmentEvent representa un evento de envíos recibido desde Redis
type ShipmentEvent struct {
	ID         string                 `json:"id"`
	Type       ShipmentEventType      `json:"event_type"`
	BusinessID uint                   `json:"business_id"`
	Timestamp  time.Time              `json:"timestamp"`
	Data       map[string]interface{} `json:"data"`
}
