package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app/templates"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type flowSendPublisher struct {
	rabbit rabbitmq.IQueue
	logger log.ILogger
}

func NewFlowSendPublisher(rabbit rabbitmq.IQueue, logger log.ILogger) templates.IFlowPublisher {
	return &flowSendPublisher{
		rabbit: rabbit,
		logger: logger.WithModule("flow_send_publisher"),
	}
}

func (p *flowSendPublisher) PublishFlowSend(ctx context.Context, message dtos.FlowSendMessage) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error serializando la respuesta del flujo: %w", err)
	}

	return p.rabbit.Publish(ctx, rabbitmq.QueueWhatsAppFlowSends, payload)
}
