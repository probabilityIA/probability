package consumerflow

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/app/usecasemessaging"
	whaErrors "github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/errors"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type flowSend struct {
	BusinessID     uint   `json:"business_id"`
	Phone          string `json:"phone"`
	TemplateName   string `json:"template_name"`
	Language       string `json:"language"`
	HeaderImageURL string `json:"header_image_url"`
}

type Consumer struct {
	rabbit  rabbitmq.IQueue
	useCase usecasemessaging.IUseCase
	logger  log.ILogger
}

func New(rabbit rabbitmq.IQueue, useCase usecasemessaging.IUseCase, logger log.ILogger) *Consumer {
	return &Consumer{
		rabbit:  rabbit,
		useCase: useCase,
		logger:  logger.WithModule("whatsapp-flow-send-consumer"),
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	if err := c.rabbit.DeclareQueue(rabbitmq.QueueWhatsAppFlowSends, true); err != nil {
		c.logger.Error(ctx).Err(err).Str("queue", rabbitmq.QueueWhatsAppFlowSends).
			Msg("Error declarando la cola")
		return err
	}

	return c.rabbit.Consume(ctx, rabbitmq.QueueWhatsAppFlowSends, func(body []byte) error {
		var send flowSend

		if err := json.Unmarshal(body, &send); err != nil {
			c.logger.Warn(ctx).Err(err).Msg("Respuesta de flujo ilegible - se descarta (ACK)")
			return nil
		}

		if send.Phone == "" || send.TemplateName == "" {
			c.logger.Warn(ctx).Uint("business_id", send.BusinessID).
				Msg("Respuesta de flujo incompleta - se descarta (ACK)")
			return nil
		}

		messageID, err := c.useCase.SendCustomTemplate(
			ctx,
			send.TemplateName,
			send.Language,
			send.Phone,
			nil,
			send.HeaderImageURL,
			send.BusinessID,
		)

		if err != nil {
			if whaErrors.IsNonRetryable(err) {
				c.logger.Warn(ctx).Err(err).
					Str("template", send.TemplateName).
					Msg("Respuesta de flujo descartada - error permanente (ACK)")
				return nil
			}

			c.logger.Error(ctx).Err(err).
				Str("template", send.TemplateName).
				Msg("Error transitorio enviando la respuesta del flujo - se reintenta")
			return err
		}

		c.logger.Info(ctx).
			Str("template", send.TemplateName).
			Str("message_id", messageID).
			Msg("Respuesta del flujo enviada")

		return nil
	})
}
