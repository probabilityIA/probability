package sse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/secamc93/probability/back/central/services/events/internal/domain/entities"
)

// sendSSEMessage envía un evento como mensaje SSE
func (m *EventManager) sendSSEMessage(w http.ResponseWriter, event entities.Event) error {
	seqVal, ok := event.Metadata["sse_seq"]
	var seq string
	if ok {
		seq = m.toString(seqVal)
	}

	message := ""
	if ok && seq != "" && seq != "0" {
		message += "id: " + seq + "\n"
	}
	message += "event: " + event.Type + "\n"
	message += "data: " + m.eventToJSON(event) + "\n\n"

	if _, err := w.Write([]byte(message)); err != nil {
		return err
	}
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
	return nil
}

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
func (m *EventManager) eventToJSON(event entities.Event) string {
	eventMap := map[string]interface{}{
		"id":          event.ID,
		"type":        event.Type,
		"business_id": event.BusinessID,
		"timestamp":   event.Timestamp.Format(time.RFC3339),
	}

	if event.Data != nil {
		eventMap["data"] = event.Data

		// Compatibilidad con campos legacy en raíz
		if dataMap, ok := event.Data["sku"]; ok {
			eventMap["sku"] = dataMap
		}
		if dataMap, ok := event.Data["quantity"]; ok {
			eventMap["quantity"] = dataMap
		}
		if dataMap, ok := event.Data["error"]; ok {
			eventMap["error"] = dataMap
		}
		if dataMap, ok := event.Data["summary"]; ok {
			eventMap["summary"] = dataMap
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
