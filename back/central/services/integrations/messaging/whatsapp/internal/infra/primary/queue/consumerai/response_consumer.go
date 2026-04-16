package consumerai

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// aiResponseDTO coincide con el DTO publicado por ai_sales/infra/secondary/queue/response_publisher.go
type aiResponseDTO struct {
	PhoneNumber  string `json:"PhoneNumber"`
	ResponseText string `json:"ResponseText"`
	BusinessID   uint   `json:"BusinessID"`
	SessionID    string `json:"SessionID"`
	Timestamp    int64  `json:"Timestamp"`
}

// Start inicia el consumer de respuestas AI
func (c *consumer) Start(ctx context.Context) error {
	queueName := rabbitmq.QueueWhatsAppAIResponse
	if err := c.queue.DeclareQueue(queueName, true); err != nil {
		c.log.Error().
			Err(err).
			Str("queue", queueName).
			Msg("[AIResponseConsumer] Error declarando cola")
		return err
	}

	go func() {
		if err := c.queue.Consume(ctx, queueName, c.handleMessage); err != nil {
			c.log.Error().Err(err).Msg("[AIResponseConsumer] Error consumiendo cola de respuestas AI")
		}
	}()

	c.log.Info().
		Str("queue", queueName).
		Msg("[AIResponseConsumer] Consumer de respuestas AI iniciado")

	return nil
}

// handleMessage procesa cada respuesta del agente AI y la envía por WhatsApp
func (c *consumer) handleMessage(body []byte) error {
	var dto aiResponseDTO
	if err := json.Unmarshal(body, &dto); err != nil {
		c.log.Error().Err(err).Msg("[AIResponseConsumer] Mensaje malformado - descartando (ACK)")
		return nil
	}

	c.log.Info().
		Str("phone", dto.PhoneNumber).
		Uint("business_id", dto.BusinessID).
		Msg("[AIResponseConsumer] Procesando respuesta AI")

	// Obtener credenciales WhatsApp desde cache Redis
	config, err := c.credentialsCache.GetWhatsAppDefaultConfig(context.Background())
	if err != nil {
		c.log.Error().
			Err(err).
			Msg("[AIResponseConsumer] Credenciales no disponibles - descartando (ACK)")
		// Error de configuración: no tiene sentido reintentar
		return nil
	}

	// Construir mensaje de texto libre
	msg := entities.TemplateMessage{
		MessagingProduct: "whatsapp",
		To:               dto.PhoneNumber,
		Type:             "text",
		TextBody:         dto.ResponseText,
	}

	// Enviar mensaje vía WhatsApp API
	messageID, err := c.wa.SendMessage(context.Background(), config.PhoneNumberID, msg, config.AccessToken)
	if err != nil {
		c.log.Error().
			Err(err).
			Str("phone", dto.PhoneNumber).
			Msg("[AIResponseConsumer] Error enviando respuesta AI por WhatsApp")
		return err // Retryable: puede ser error temporal de red
	}

	c.log.Info().
		Str("message_id", messageID).
		Str("phone", dto.PhoneNumber).
		Msg("[AIResponseConsumer] Respuesta AI enviada exitosamente")

	return nil
}
