package queue

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/monitoring/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/monitoring/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const queueName = "monitoring.alerts"

type publisher struct {
	queue rabbitmq.IQueue
	log   log.ILogger
}

// New crea una nueva instancia del publisher de alertas
func New(queue rabbitmq.IQueue, logger log.ILogger) ports.IAlertPublisher {
	return &publisher{
		queue: queue,
		log:   logger,
	}
}

// Publish publica un evento de alerta en la cola monitoring.alerts
func (p *publisher) Publish(ctx context.Context, event entities.AlertEvent) error {
	if err := p.queue.DeclareQueue(queueName, true); err != nil {
		p.log.Error(ctx).
			Err(err).
			Str("queue", queueName).
			Msg("[Monitoring] Error declarando cola")
		return err
	}

	body, err := json.Marshal(event)
	if err != nil {
		p.log.Error(ctx).Err(err).Msg("[Monitoring] Error serializando evento de alerta")
		return err
	}

	if err := p.queue.Publish(ctx, queueName, body); err != nil {
		p.log.Error(ctx).
			Err(err).
			Str("queue", queueName).
			Msg("[Monitoring] Error publicando evento de alerta")
		return err
	}

	p.log.Info(ctx).
		Str("queue", queueName).
		Str("alert_type", event.AlertType).
		Msg("[Monitoring] Evento de alerta publicado")

	return nil
}
