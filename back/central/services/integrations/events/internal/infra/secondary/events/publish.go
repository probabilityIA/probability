package events

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/events/internal/domain"
)

// PublishEvent publica un evento al sistema
func (m *IntegrationEventManager) PublishEvent(event domain.IntegrationEvent) {
	// Agregar timestamp si no está presente
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Generar ID si no está presente
	if event.ID == "" {
		event.ID = fmt.Sprintf("%d-%s", time.Now().Unix(), generateRandomID())
	}

	// #region agent log
	if f, err := os.OpenFile("/home/cam/Desktop/probability/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		logData, _ := json.Marshal(map[string]interface{}{
			"location": "publish.go:PublishEvent",
			"message":  "Publishing integration event to channel",
			"data": map[string]interface{}{
				"event_id":       event.ID,
				"event_type":     string(event.Type),
				"integration_id": event.IntegrationID,
				"business_id":    event.BusinessID,
			},
			"timestamp":    time.Now().UnixMilli(),
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "C",
		})
		f.WriteString(string(logData) + "\n")
		f.Close()
	}
	// #endregion

	// Enviar al canal de eventos (no bloqueante)
	select {
	case m.eventChan <- event:
		// #region agent log
		if f, err := os.OpenFile("/home/cam/Desktop/probability/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			logData, _ := json.Marshal(map[string]interface{}{
				"location": "publish.go:PublishEvent",
				"message":  "Event sent to channel successfully",
				"data": map[string]interface{}{
					"event_id":   event.ID,
					"event_type": string(event.Type),
				},
				"timestamp":    time.Now().UnixMilli(),
				"sessionId":    "debug-session",
				"runId":        "run1",
				"hypothesisId": "C",
			})
			f.WriteString(string(logData) + "\n")
			f.Close()
		}
		// #endregion
	default:
		// Si el canal está lleno, loggear warning pero no bloquear
		// #region agent log
		if f, err := os.OpenFile("/home/cam/Desktop/probability/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			logData, _ := json.Marshal(map[string]interface{}{
				"location": "publish.go:PublishEvent",
				"message":  "Event channel full, dropping event",
				"data": map[string]interface{}{
					"event_id":   event.ID,
					"event_type": string(event.Type),
				},
				"timestamp":    time.Now().UnixMilli(),
				"sessionId":    "debug-session",
				"runId":        "run1",
				"hypothesisId": "C",
			})
			f.WriteString(string(logData) + "\n")
			f.Close()
		}
		// #endregion
		if m.logger != nil {
			m.logger.Warn().
				Str("event_id", event.ID).
				Str("event_type", string(event.Type)).
				Msg("Event channel full, dropping event")
		}
	}
}

// generateRandomID genera un ID aleatorio corto
func generateRandomID() string {
	return fmt.Sprintf("%x", time.Now().UnixNano()%1000000)
}
