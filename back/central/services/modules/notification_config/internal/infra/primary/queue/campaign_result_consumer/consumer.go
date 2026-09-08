package campaign_result_consumer

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app/campaigns"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type Consumer struct {
	rabbit  rabbitmq.IQueue
	useCase campaigns.IUseCase
	logger  log.ILogger
}

func New(rabbit rabbitmq.IQueue, useCase campaigns.IUseCase, logger log.ILogger) *Consumer {
	return &Consumer{
		rabbit:  rabbit,
		useCase: useCase,
		logger:  logger.WithModule("campaign_result_consumer"),
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	if err := c.rabbit.DeclareQueue(rabbitmq.QueueCampaignResults, true); err != nil {
		c.logger.Error(ctx).Err(err).Str("queue", rabbitmq.QueueCampaignResults).
			Msg("Error declarando la cola")
		return err
	}

	return c.rabbit.Consume(ctx, rabbitmq.QueueCampaignResults, func(body []byte) error {
		var result dtos.CampaignSendResult

		if err := json.Unmarshal(body, &result); err != nil {
			c.logger.Warn(ctx).Err(err).
				Msg("Resultado de envio de campana ilegible - se descarta (ACK)")
			return nil
		}

		if result.SendID == 0 {
			c.logger.Warn(ctx).Msg("Resultado de campana sin send_id - se descarta (ACK)")
			return nil
		}

		if err := c.useCase.MarkSendResult(ctx, result); err != nil {
			c.logger.Error(ctx).Err(err).Uint("send_id", result.SendID).
				Msg("Error guardando el resultado del envio de campana - se reintenta")
			return err
		}

		return nil
	})
}
