package consumeralert

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
)

const (
	queueName    = "monitoring.alerts"
	adminPhone   = "573023406789"
	templateName = "alerta_servidor"
)

// AlertEvent coincide con el struct publicado por el módulo monitoring
type AlertEvent struct {
	AlertType string    `json:"alert_type"`
	Summary   string    `json:"summary"`
	Status    string    `json:"status"`
	FiredAt   time.Time `json:"fired_at"`
}

// Start inicia el consumer de alertas de monitoreo
func (c *consumerAlert) Start(ctx context.Context) error {
	if err := c.queue.DeclareQueue(queueName, true); err != nil {
		c.log.Error().
			Err(err).
			Str("queue", queueName).
			Msg("[AlertConsumer] Error declarando cola")
		return err
	}

	go func() {
		if err := c.queue.Consume(ctx, queueName, c.handleMessage); err != nil {
			c.log.Error().Err(err).Msg("[AlertConsumer] Error consumiendo cola de alertas")
		}
	}()

	c.log.Info().
		Str("queue", queueName).
		Msg("[AlertConsumer] Consumer de alertas de monitoreo iniciado")

	return nil
}

// handleMessage procesa cada mensaje de alerta
func (c *consumerAlert) handleMessage(body []byte) error {
	var event AlertEvent
	if err := json.Unmarshal(body, &event); err != nil {
		c.log.Error().Err(err).Msg("[AlertConsumer] Error deserializando evento de alerta")
		return err
	}

	c.log.Info().
		Str("alert_type", event.AlertType).
		Str("status", event.Status).
		Str("summary", event.Summary).
		Msg("[AlertConsumer] Procesando alerta de monitoreo")

	// Ignorar alertas que no estén en estado "firing"
	if event.Status != "firing" {
		c.log.Info().
			Str("status", event.Status).
			Msg("[AlertConsumer] Ignorando alerta no activa")
		return nil
	}

	// Obtener credenciales desde env vars (no desde DB - es alerta de infra)
	phoneNumberIDStr := c.env.Get("WHATSAPP_PHONE_NUMBER_ID")
	phoneNumberID, err := strconv.ParseUint(phoneNumberIDStr, 10, 64)
	if err != nil {
		c.log.Error().
			Err(err).
			Str("phone_number_id", phoneNumberIDStr).
			Msg("[AlertConsumer] WHATSAPP_PHONE_NUMBER_ID inválido")
		return err
	}

	token := c.env.Get("WHATSAPP_TOKEN")
	if token == "" {
		c.log.Error().Msg("[AlertConsumer] WHATSAPP_TOKEN no configurado")
		return nil
	}

	// Construir variables del template alerta_servidor
	variables := map[string]string{
		"1": event.AlertType, // {{1}} tipo de alerta: RAM / CPU / Disco
		"2": event.Summary,   // {{2}} descripción: "87.3% - supera umbral de 85%"
	}

	// Construir el mensaje de plantilla directamente (sin pasar por el use case que requiere DB)
	msg := buildAlertTemplateMessage(adminPhone, variables)

	// Enviar mensaje de WhatsApp
	messageID, err := c.wa.SendMessage(context.Background(), uint(phoneNumberID), msg, token)
	if err != nil {
		c.log.Error().
			Err(err).
			Str("alert_type", event.AlertType).
			Str("phone", adminPhone).
			Msg("[AlertConsumer] Error enviando alerta por WhatsApp")
		return err
	}

	c.log.Info().
		Str("message_id", messageID).
		Str("alert_type", event.AlertType).
		Str("phone", adminPhone).
		Msg("[AlertConsumer] Alerta enviada por WhatsApp exitosamente")

	return nil
}

// buildAlertTemplateMessage construye el TemplateMessage para la plantilla alerta_servidor
func buildAlertTemplateMessage(phoneNumber string, variables map[string]string) entities.TemplateMessage {
	bodyParams := []entities.TemplateParameter{
		{Type: "text", Text: variables["1"]},
		{Type: "text", Text: variables["2"]},
	}

	return entities.TemplateMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               phoneNumber,
		Type:             "template",
		Template: entities.TemplateData{
			Name:     templateName,
			Language: entities.TemplateLanguage{Code: "es"},
			Components: []entities.TemplateComponent{
				{
					Type:       "body",
					Parameters: bodyParams,
				},
			},
		},
	}
}
