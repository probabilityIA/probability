package events

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/events/internal/domain"
)

// toString convierte cualquier valor a string
func (m *EventManager) toString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case int:
		return fmt.Sprintf("%d", val)
	case int64:
		return fmt.Sprintf("%d", val)
	case int32:
		return fmt.Sprintf("%d", val)
	case float64:
		return fmt.Sprintf("%.2f", val)
	case float32:
		return fmt.Sprintf("%.2f", val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", v)
	}
}

// eventToJSON convierte un evento a JSON string usando json.Marshal para seguridad
func (m *EventManager) eventToJSON(event domain.Event) string {
	eventMap := map[string]interface{}{
		"id":          event.ID,
		"type":        string(event.Type),
		"business_id": event.BusinessID,
		"timestamp":   event.Timestamp.Format(time.RFC3339),
	}

	if event.Data != nil {
		eventMap["data"] = event.Data

		// Mantener compatibilidad con campos legacy en raíz
		if dataMap, ok := event.Data.(map[string]interface{}); ok {
			if sku, ok := dataMap["sku"]; ok {
				eventMap["sku"] = sku
			}
			if quantity, ok := dataMap["quantity"]; ok {
				eventMap["quantity"] = quantity
			}
			if errorMsg, ok := dataMap["error"]; ok {
				eventMap["error"] = errorMsg
			}
			if summary, ok := dataMap["summary"]; ok {
				eventMap["summary"] = summary
			}
		}
	}

	if len(event.Metadata) > 0 {
		eventMap["metadata"] = event.Metadata
	}

	jsonBytes, err := json.Marshal(eventMap)
	if err != nil {
		return "{}"
	}
	return string(jsonBytes)
}
