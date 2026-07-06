package consumerauthotp

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const (
	queueName    = rabbitmq.QueueAuthPasswordResetOTP
	templateName = "recuperacion_codigo"
)

type OTPEvent struct {
	Phone          string `json:"phone"`
	Code           string `json:"code"`
	UserName       string `json:"user_name"`
	ExpiresMinutes int    `json:"expires_minutes"`
}

func (c *consumerAuthOTP) Start(ctx context.Context) error {
	if err := c.queue.DeclareQueue(queueName, true); err != nil {
		c.log.Error().
			Err(err).
			Str("queue", queueName).
			Msg("[AuthOTPConsumer] Error declarando cola")
		return err
	}

	go func() {
		if err := c.queue.Consume(ctx, queueName, c.handleMessage); err != nil {
			c.log.Error().Err(err).Msg("[AuthOTPConsumer] Error consumiendo cola de OTP")
		}
	}()

	c.log.Info().
		Str("queue", queueName).
		Msg("[AuthOTPConsumer] Consumer de OTP de recuperacion iniciado")

	return nil
}

func (c *consumerAuthOTP) handleMessage(body []byte) error {
	var event OTPEvent
	if err := json.Unmarshal(body, &event); err != nil {
		c.log.Error().Err(err).Msg("[AuthOTPConsumer] Error deserializando evento OTP")
		return err
	}

	if event.Phone == "" || event.Code == "" {
		c.log.Error().Msg("[AuthOTPConsumer] Evento OTP sin telefono o codigo - descartando")
		return nil
	}

	config, err := c.credentialsCache.GetWhatsAppDefaultConfig(context.Background())
	if err != nil {
		c.log.Error().
			Err(err).
			Msg("[AuthOTPConsumer] Error obteniendo credenciales de WhatsApp desde cache")
		return nil
	}

	msg := buildOTPTemplateMessage(event.Phone, event.Code)

	messageID, err := c.wa.SendMessage(context.Background(), config.PhoneNumberID, msg, config.AccessToken)
	if err != nil {
		c.log.Error().
			Err(err).
			Str("phone", event.Phone).
			Msg("[AuthOTPConsumer] Error enviando codigo OTP por WhatsApp - se descarta el mensaje (el usuario puede reenviar)")
		return nil
	}

	c.log.Info().
		Str("message_id", messageID).
		Str("phone", event.Phone).
		Msg("[AuthOTPConsumer] Codigo OTP enviado por WhatsApp exitosamente")

	return nil
}

func buildOTPTemplateMessage(phoneNumber, code string) entities.TemplateMessage {
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
					Type: "body",
					Parameters: []entities.TemplateParameter{
						{Type: "text", Text: code},
					},
				},
				{
					Type:    "button",
					SubType: "url",
					Index:   "0",
					Parameters: []entities.TemplateParameter{
						{Type: "text", Text: code},
					},
				},
			},
		},
	}
}
