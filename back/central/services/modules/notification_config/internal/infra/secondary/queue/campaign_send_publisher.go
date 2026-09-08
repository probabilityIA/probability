package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app/campaigns"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type campaignSendPublisher struct {
	rabbit rabbitmq.IQueue
	logger log.ILogger
}

func NewCampaignSendPublisher(rabbit rabbitmq.IQueue, logger log.ILogger) campaigns.ICampaignPublisher {
	return &campaignSendPublisher{
		rabbit: rabbit,
		logger: logger.WithModule("campaign_send_publisher"),
	}
}

func (p *campaignSendPublisher) PublishCampaignSend(ctx context.Context, message dtos.CampaignSendMessage) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error serializando el envio de campana: %w", err)
	}

	return p.rabbit.Publish(ctx, rabbitmq.QueueCampaignSends, payload)
}
