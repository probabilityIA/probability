package app

import (
	"fmt"
	"strings"
)

// buildSubject genera el asunto del email basado en el tipo de evento
func buildSubject(eventType string) string {
	return fmt.Sprintf("Notificación: %s", eventType)
}

// buildHTML genera el contenido HTML del email basado en el tipo de evento y sus datos
func buildHTML(eventType string, eventData map[string]interface{}) string {
	var sb strings.Builder

	sb.WriteString(`<!DOCTYPE html><html><head><meta charset="UTF-8"></head><body style="font-family:Arial,sans-serif;max-width:600px;margin:0 auto;padding:20px;">`)
	sb.WriteString(fmt.Sprintf(`<h2 style="color:#333;">Evento: %s</h2>`, eventType))
	sb.WriteString(`<p>Se ha producido un evento en tu cuenta.</p>`)

	// Incluir datos relevantes si están disponibles
	if len(eventData) > 0 {
		sb.WriteString(`<table style="border-collapse:collapse;width:100%;margin-top:16px;">`)
		for key, value := range eventData {
			sb.WriteString(fmt.Sprintf(
				`<tr><td style="padding:8px;border:1px solid #ddd;font-weight:bold;">%s</td><td style="padding:8px;border:1px solid #ddd;">%v</td></tr>`,
				key, value,
			))
		}
		sb.WriteString(`</table>`)
	}

	sb.WriteString(`<hr style="margin-top:24px;"><p style="color:#999;font-size:12px;">Este es un mensaje automático, no responder.</p>`)
	sb.WriteString(`</body></html>`)

	return sb.String()
}
