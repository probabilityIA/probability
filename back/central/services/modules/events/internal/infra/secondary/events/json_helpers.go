package events

import (
	"fmt"
	"strings"
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

// eventToJSON convierte un evento a JSON string
func (m *EventManager) eventToJSON(event domain.Event) string {
	result := "{"
	result += `"id":"` + event.ID + `",`
	result += `"type":"` + string(event.Type) + `",`
	result += `"business_id":"` + event.BusinessID + `",`
	result += `"timestamp":"` + event.Timestamp.Format(time.RFC3339) + `"`

	// Incluir campos específicos del Data si existen
	if event.Data != nil {
		// Serializar todo el objeto Data en el campo "data"
		result += `,"data":` + m.dataToJSON(event.Data)

		// Mantener compatibilidad con campos específicos si es necesario (opcional)
		if dataMap, ok := event.Data.(map[string]interface{}); ok {
			if sku, ok := dataMap["sku"]; ok {
				result += `,"sku":"` + m.escapeJSONString(m.toString(sku)) + `"`
			}
			if quantity, ok := dataMap["quantity"]; ok {
				result += `,"quantity":` + m.toString(quantity)
			}
			if errorMsg, ok := dataMap["error"]; ok {
				result += `,"error":"` + m.escapeJSONString(m.toString(errorMsg)) + `"`
			}
			if summary, ok := dataMap["summary"]; ok {
				result += `,"summary":` + m.dataToJSON(summary)
			}
		}
	}

	if len(event.Metadata) > 0 {
		result += `,"metadata":` + m.dataToJSON(event.Metadata)
	}
	result += "}"
	return result
}

// dataToJSON convierte datos a JSON string
func (m *EventManager) dataToJSON(data interface{}) string {
	switch v := data.(type) {
	case string:
		return `"` + v + `"`
	case int64:
		return string(rune(v))
	case time.Time:
		return `"` + v.Format(time.RFC3339) + `"`
	case map[string]interface{}:
		result := "{"
		for key, value := range v {
			if result != "{" {
				result += ","
			}
			result += `"` + key + `":` + m.dataToJSON(value)
		}
		result += "}"
		return result
	default:
		return `"` + m.toString(v) + `"`
	}
}

// escapeJSONString escapa caracteres especiales en una cadena para JSON
func (m *EventManager) escapeJSONString(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\r", `\r`)
	s = strings.ReplaceAll(s, "\t", `\t`)
	return s
}
