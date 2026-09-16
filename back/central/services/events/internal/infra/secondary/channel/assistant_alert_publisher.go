package channel

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/events/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/events/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/events/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type assistantAlertPublisher struct {
	rabbitMQ rabbitmq.IQueue
	logger   log.ILogger
}

func NewAssistantAlertPublisher(rabbitMQ rabbitmq.IQueue, logger log.ILogger) ports.IAssistantAlertPublisher {
	if rabbitMQ != nil {
		if err := rabbitMQ.DeclareQueue(rabbitmq.QueueAIAssistantAlerts, true); err != nil {
			logger.Error(context.Background()).Err(err).Str("queue", rabbitmq.QueueAIAssistantAlerts).
				Msg("No se pudo declarar la cola de alertas del asistente")
		}
	}
	return &assistantAlertPublisher{rabbitMQ: rabbitMQ, logger: logger}
}

func (p *assistantAlertPublisher) PublishToAssistant(ctx context.Context, event entities.Event) error {
	if p.rabbitMQ == nil || event.BusinessID == 0 {
		return nil
	}
	body, err := json.Marshal(dtos.EventEnvelope{
		ID:            event.ID,
		Type:          event.Type,
		Category:      event.Category,
		BusinessID:    event.BusinessID,
		IntegrationID: event.IntegrationID,
		Timestamp:     event.Timestamp,
		Data:          event.Data,
		Metadata:      event.Metadata,
	})
	if err != nil {
		return fmt.Errorf("serializar alerta del asistente: %w", err)
	}
	return p.rabbitMQ.Publish(ctx, rabbitmq.QueueAIAssistantAlerts, body)
}
