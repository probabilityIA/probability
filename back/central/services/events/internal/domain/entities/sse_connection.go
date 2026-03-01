package entities

import (
	"fmt"
	"net/http"
)

// SSEConnectionFilter define los filtros para una conexión SSE
type SSEConnectionFilter struct {
	IntegrationID *uint
	EventTypes    []string
	OrderIDs      []string
}

// IsEmpty verifica si el filtro está vacío (sin filtros aplicados)
func (f *SSEConnectionFilter) IsEmpty() bool {
	return f.IntegrationID == nil && len(f.EventTypes) == 0 && len(f.OrderIDs) == 0
}

// Matches verifica si un evento coincide con los filtros
func (f *SSEConnectionFilter) Matches(event Event) bool {
	if f.IsEmpty() {
		return true
	}

	// Filtrar por integration_id
	if f.IntegrationID != nil {
		if eventMetadata, ok := event.Metadata["integration_id"]; ok {
			if eventIntegrationID, ok := eventMetadata.(uint); ok {
				if eventIntegrationID != *f.IntegrationID {
					return false
				}
			} else if eventIntegrationID, ok := eventMetadata.(float64); ok {
				if uint(eventIntegrationID) != *f.IntegrationID {
					return false
				}
			} else if event.IntegrationID != *f.IntegrationID {
				return false
			}
		} else if event.IntegrationID != *f.IntegrationID {
			return false
		}
	}

	// Filtrar por event_types
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

	// Filtrar por order_ids
	if len(f.OrderIDs) > 0 {
		var orderID string
		if orderIDMeta, ok := event.Metadata["order_id"]; ok {
			if orderIDStr, ok := orderIDMeta.(string); ok {
				orderID = orderIDStr
			}
		}
		if orderID == "" {
			if dataMap := event.Data; dataMap != nil {
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
			return false
		}
	}

	return true
}

// SSEConnection representa una conexión SSE con sus filtros
type SSEConnection struct {
	BusinessID   uint
	Filter       *SSEConnectionFilter
	Writer       http.ResponseWriter
	ConnectionID string
}

// IsSuperUser verifica si es un super usuario (business_id = 0)
func (c *SSEConnection) IsSuperUser() bool {
	return c.BusinessID == 0
}

// MatchesBusiness verifica si un business_id coincide con esta conexión
func (c *SSEConnection) MatchesBusiness(businessID uint) bool {
	if c.IsSuperUser() {
		return true
	}
	return c.BusinessID == businessID
}

// MatchesBusinessStr verifica si un business_id string coincide
func (c *SSEConnection) MatchesBusinessStr(businessID string) bool {
	if c.IsSuperUser() {
		return true
	}
	if businessID == "" {
		return false
	}
	return businessID == fmt.Sprintf("%d", c.BusinessID)
}
