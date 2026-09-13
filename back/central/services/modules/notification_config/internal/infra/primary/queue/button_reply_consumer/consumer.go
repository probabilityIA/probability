package button_reply_consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app/templates"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type Consumer struct {
	rabbit  rabbitmq.IQueue
	useCase templates.IUseCase
	logger  log.ILogger
}

func New(rabbit rabbitmq.IQueue, useCase templates.IUseCase, logger log.ILogger) *Consumer {
	return &Consumer{
		rabbit:  rabbit,
		useCase: useCase,
		logger:  logger.WithModule("whatsapp_button_reply_consumer"),
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	if err := c.rabbit.DeclareQueue(rabbitmq.QueueWhatsAppButtonReplies, true); err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", rabbitmq.QueueWhatsAppButtonReplies, err)
	}

	return c.rabbit.Consume(ctx, rabbitmq.QueueWhatsAppButtonReplies, func(body []byte) error {
		var event dtos.ButtonReplyEvent

		if err := json.Unmarshal(body, &event); err != nil {
			c.logger.Warn(ctx).Err(err).Msg("Respuesta de boton ilegible - se descarta (ACK)")
			return nil
		}

		if err := c.useCase.HandleButtonReply(ctx, event); err != nil {
			c.logger.Error(ctx).Err(err).
				Uint("business_id", event.BusinessID).
				Str("button_text", event.ButtonText).
				Msg("Error resolviendo la respuesta encadenada - se reintenta")
			return err
		}

		return nil
	})
}
