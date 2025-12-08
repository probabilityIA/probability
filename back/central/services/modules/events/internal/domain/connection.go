package domain

import (
	"fmt"
	"net/http"
)

// ───────────────────────────────────────────
//
//	SSE CONNECTION STRUCTURES
//
// ───────────────────────────────────────────

// SSEConnectionFilter define los filtros para una conexión SSE
type SSEConnectionFilter struct {
	// IntegrationID filtra eventos por integración específica (opcional)
	// Si es nil, recibe eventos de todas las integraciones del business
	IntegrationID *uint `json:"integration_id,omitempty"`

	// EventTypes filtra eventos por tipos específicos (opcional)
	// Si está vacío, recibe todos los tipos de eventos
	EventTypes []EventType `json:"event_types,omitempty"`

	// OrderIDs filtra eventos de órdenes específicas (opcional)
	// Si está vacío, recibe eventos de todas las órdenes
	OrderIDs []string `json:"order_ids,omitempty"`
}

// IsEmpty verifica si el filtro está vacío (sin filtros aplicados)
func (f *SSEConnectionFilter) IsEmpty() bool {
	return f.IntegrationID == nil && len(f.EventTypes) == 0 && len(f.OrderIDs) == 0
}

// Matches verifica si un evento coincide con los filtros
func (f *SSEConnectionFilter) Matches(event Event) bool {
	// Si el filtro está vacío, acepta todos los eventos
	if f.IsEmpty() {
		return true
	}

	// Filtrar por integration_id si está especificado
	if f.IntegrationID != nil {
		// El integration_id está en el metadata o en el evento
		if eventMetadata, ok := event.Metadata["integration_id"]; ok {
			if eventIntegrationID, ok := eventMetadata.(uint); ok {
				if eventIntegrationID != *f.IntegrationID {
					return false
				}
			} else if eventIntegrationID, ok := eventMetadata.(float64); ok {
				if uint(eventIntegrationID) != *f.IntegrationID {
					return false
				}
			} else {
				// Si no hay integration_id en metadata, verificar el IntegrationID del evento
				if uint(event.IntegrationID) != *f.IntegrationID {
					return false
				}
			}
		} else {
			// Si no hay integration_id en metadata, verificar el IntegrationID del evento
			if uint(event.IntegrationID) != *f.IntegrationID {
				return false
			}
		}
	}

	// Filtrar por event_types si están especificados
	if len(f.EventTypes) > 0 {
		matched := false
		for _, filterType := range f.EventTypes {
			if event.Type == filterType {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// Filtrar por order_ids si están especificados
	if len(f.OrderIDs) > 0 {
		// Buscar order_id en metadata o en data
		var orderID string
		if orderIDMeta, ok := event.Metadata["order_id"]; ok {
			if orderIDStr, ok := orderIDMeta.(string); ok {
				orderID = orderIDStr
			}
		}
		if orderID == "" {
			if dataMap, ok := event.Data.(map[string]interface{}); ok {
				if orderIDData, ok := dataMap["order_id"]; ok {
					if orderIDStr, ok := orderIDData.(string); ok {
						orderID = orderIDStr
					}
				}
			}
		}

		if orderID != "" {
			matched := false
			for _, filterOrderID := range f.OrderIDs {
				if orderID == filterOrderID {
					matched = true
					break
				}
			}
			if !matched {
				return false
			}
		} else {
			// Si no hay order_id en el evento, no coincide con el filtro
			return false
		}
	}

	return true
}

// SSEConnection representa una conexión SSE con sus filtros
type SSEConnection struct {
	BusinessID   uint                 // Business ID (0 = super usuario, recibe todos)
	Filter       *SSEConnectionFilter // Filtros opcionales
	Writer       http.ResponseWriter
	ConnectionID string // ID único de la conexión
}

// IsSuperUser verifica si es un super usuario (business_id = 0)
func (c *SSEConnection) IsSuperUser() bool {
	return c.BusinessID == 0
}

// MatchesBusiness verifica si un business_id coincide con esta conexión
func (c *SSEConnection) MatchesBusiness(businessID string) bool {
	// Super usuario recibe todos los businesses
	if c.IsSuperUser() {
		return true
	}

	// Convertir businessID string a uint para comparar
	// El businessID viene como string en el evento
	if businessID == "" {
		return false
	}

	// Comparar directamente con el business_id del evento
	// El businessID viene como string en el evento, necesitamos convertirlo
	// Por ahora, comparamos directamente con el string
	return businessID == fmt.Sprintf("%d", c.BusinessID)
}
