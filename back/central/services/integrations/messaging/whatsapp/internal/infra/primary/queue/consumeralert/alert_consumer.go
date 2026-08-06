package consumeralert

import (
	"context"
	"encoding/json"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
	whaErrors "github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/errors"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const (
	queueName    = rabbitmq.QueueMonitoringAlerts
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

	// Obtener credenciales desde cache Redis (calentadas por integrations/core al startup)
	config, err := c.credentialsCache.GetWhatsAppDefaultConfig(context.Background())
	if err != nil {
		c.log.Error().
			Err(err).
			Msg("[AlertConsumer] Error obteniendo credenciales de WhatsApp desde cache - verifica warmup de integrations/core")
		// Retornar nil para hacer ack del mensaje y evitar loop infinito de reintentos.
		// Este error es de configuración, no del mensaje — reintentar no lo resuelve.
		return nil
	}

	// Construir variables del template alerta_servidor
	variables := map[string]string{
		"1": event.AlertType, // {{1}} tipo de alerta: RAM / CPU / Disco
		"2": event.Summary,   // {{2}} descripción: "87.3% - supera umbral de 85%"
	}

	enviados := 0
	var falloTransitorio error

	for _, phone := range c.phones {
		msg := buildAlertTemplateMessage(phone, variables)

		messageID, err := c.wa.SendMessage(context.Background(), config.PhoneNumberID, msg, config.AccessToken)
		if err != nil {
			if whaErrors.IsNonRetryable(err) {
				c.log.Warn().
					Err(err).
					Str("alert_type", event.AlertType).
					Str("phone", phone).
					Msg("[AlertConsumer] Destinatario descartado - error no reintentable")
				continue
			}
			falloTransitorio = err
			c.log.Error().
				Err(err).
				Str("alert_type", event.AlertType).
				Str("phone", phone).
				Msg("[AlertConsumer] Error enviando alerta por WhatsApp")
			continue
		}

		enviados++
		c.log.Info().
			Str("message_id", messageID).
			Str("alert_type", event.AlertType).
			Str("phone", phone).
			Msg("[AlertConsumer] Alerta enviada por WhatsApp exitosamente")
	}

	if falloTransitorio != nil {
		c.log.Error().
			Int("enviados", enviados).
			Int("destinatarios", len(c.phones)).
			Msg("[AlertConsumer] Reintentando alerta por fallo transitorio - los destinatarios ya enviados recibiran duplicado")
		return falloTransitorio
	}

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
