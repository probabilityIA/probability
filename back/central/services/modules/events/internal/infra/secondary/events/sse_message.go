package events

import (
	"net/http"

	"github.com/secamc93/probability/back/central/services/modules/events/internal/domain"
)

// sendSSEMessage envía un evento como mensaje SSE
func (m *EventManager) sendSSEMessage(w http.ResponseWriter, event domain.Event) error {
	seqVal, ok := event.Metadata["sse_seq"]
	var seq string
	if ok {
		seq = m.toString(seqVal)
	}

	message := ""
	if ok && seq != "" && seq != "0" {
		message += "id: " + seq + "\n"
	}
	message += "event: " + string(event.Type) + "\n"
	message += "data: " + m.eventToJSON(event) + "\n\n"

	if _, err := w.Write([]byte(message)); err != nil {
		return err
	}
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
	return nil
}
