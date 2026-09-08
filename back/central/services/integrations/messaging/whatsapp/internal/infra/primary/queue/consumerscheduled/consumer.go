package consumerscheduled

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/app/usecasemessaging"
	whaErrors "github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/errors"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type scheduledSend struct {
	SendID       uint     `json:"send_id"`
	RuleID       uint     `json:"rule_id"`
	BusinessID   uint     `json:"business_id"`
	ClientID     uint     `json:"client_id"`
	Phone        string   `json:"phone"`
	TemplateName string   `json:"template_name"`
	Language     string   `json:"language"`
	Parameters   []string `json:"parameters"`
}

type scheduledResult struct {
	SendID       uint   `json:"send_id"`
	Status       string `json:"status"`
	MessageID    string `json:"message_id"`
	ErrorMessage string `json:"error_message"`
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
		logger:  logger.WithModule("whatsapp-scheduled-send-consumer"),
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	if err := c.rabbit.DeclareQueue(rabbitmq.QueueScheduledNotificationSends, true); err != nil {
		c.logger.Error(ctx).Err(err).Str("queue", rabbitmq.QueueScheduledNotificationSends).
			Msg("Error declarando la cola")
		return err
	}

	return c.rabbit.Consume(ctx, rabbitmq.QueueScheduledNotificationSends, func(body []byte) error {
		var send scheduledSend

		if err := json.Unmarshal(body, &send); err != nil {
			c.logger.Warn(ctx).Err(err).Msg("Envio programado ilegible - se descarta (ACK)")
			return nil
		}

		if send.SendID == 0 || send.Phone == "" || send.TemplateName == "" {
			c.logger.Warn(ctx).Uint("send_id", send.SendID).
				Msg("Envio programado incompleto - se descarta (ACK)")
			return c.publishResult(ctx, scheduledResult{
				SendID:       send.SendID,
				Status:       "failed",
				ErrorMessage: "mensaje incompleto",
			})
		}

		messageID, err := c.useCase.SendCustomTemplate(
			ctx,
			send.TemplateName,
			send.Language,
			send.Phone,
			send.Parameters,
			send.BusinessID,
		)

		if err != nil {
			if whaErrors.IsNonRetryable(err) {
				c.logger.Warn(ctx).Err(err).
					Uint("send_id", send.SendID).
					Msg("Envio programado descartado - error permanente (ACK)")
				return c.publishResult(ctx, scheduledResult{
					SendID:       send.SendID,
					Status:       "failed",
					ErrorMessage: err.Error(),
				})
			}

			c.logger.Error(ctx).Err(err).
				Uint("send_id", send.SendID).
				Msg("Error transitorio enviando el mensaje programado - se reintenta")
			return err
		}

		return c.publishResult(ctx, scheduledResult{
			SendID:    send.SendID,
			Status:    "sent",
			MessageID: messageID,
		})
	})
}

func (c *Consumer) publishResult(ctx context.Context, result scheduledResult) error {
	payload, err := json.Marshal(result)
	if err != nil {
		c.logger.Error(ctx).Err(err).Uint("send_id", result.SendID).
			Msg("Error serializando el resultado del envio programado")
		return nil
	}

	if err := c.rabbit.Publish(ctx, rabbitmq.QueueScheduledNotificationResults, payload); err != nil {
		c.logger.Error(ctx).Err(err).Uint("send_id", result.SendID).
			Msg("Error publicando el resultado del envio programado")
		return err
	}

	return nil
}
