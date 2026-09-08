package consumercampaign

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/app/usecasemessaging"
	whaErrors "github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/errors"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type campaignSend struct {
	SendID       uint     `json:"send_id"`
	CampaignID   uint     `json:"campaign_id"`
	BusinessID   uint     `json:"business_id"`
	ClientID     uint     `json:"client_id"`
	Phone        string   `json:"phone"`
	TemplateName string   `json:"template_name"`
	Language     string   `json:"language"`
	Parameters   []string `json:"parameters"`
}

type campaignResult struct {
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
		logger:  logger.WithModule("whatsapp-campaign-send-consumer"),
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	if err := c.rabbit.DeclareQueue(rabbitmq.QueueCampaignSends, true); err != nil {
		c.logger.Error(ctx).Err(err).Str("queue", rabbitmq.QueueCampaignSends).
			Msg("Error declarando la cola")
		return err
	}

	return c.rabbit.Consume(ctx, rabbitmq.QueueCampaignSends, func(body []byte) error {
		var send campaignSend

		if err := json.Unmarshal(body, &send); err != nil {
			c.logger.Warn(ctx).Err(err).Msg("Envio de campana ilegible - se descarta (ACK)")
			return nil
		}

		if send.SendID == 0 || send.Phone == "" || send.TemplateName == "" {
			c.logger.Warn(ctx).Uint("send_id", send.SendID).
				Msg("Envio de campana incompleto - se descarta (ACK)")
			return c.publishResult(ctx, campaignResult{
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
					Msg("Envio de campana descartado - error permanente (ACK)")
				return c.publishResult(ctx, campaignResult{
					SendID:       send.SendID,
					Status:       "failed",
					ErrorMessage: err.Error(),
				})
			}

			c.logger.Error(ctx).Err(err).
				Uint("send_id", send.SendID).
				Msg("Error transitorio enviando el mensaje de campana - se reintenta")
			return err
		}

		return c.publishResult(ctx, campaignResult{
			SendID:    send.SendID,
			Status:    "sent",
			MessageID: messageID,
		})
	})
}

func (c *Consumer) publishResult(ctx context.Context, result campaignResult) error {
	payload, err := json.Marshal(result)
	if err != nil {
		c.logger.Error(ctx).Err(err).Uint("send_id", result.SendID).
			Msg("Error serializando el resultado del envio de campana")
		return nil
	}

	if err := c.rabbit.Publish(ctx, rabbitmq.QueueCampaignResults, payload); err != nil {
		c.logger.Error(ctx).Err(err).Uint("send_id", result.SendID).
			Msg("Error publicando el resultado del envio de campana")
		return err
	}

	return nil
}
