package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type alertEnvelope struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	BusinessID uint           `json:"business_id"`
	Timestamp  time.Time      `json:"timestamp"`
	Data       map[string]any `json:"data"`
}

func (c *Consumer) StartAlerts(ctx context.Context) error {
	queueName := rabbitmq.QueueAIAssistantAlerts
	if err := c.queue.DeclareQueue(queueName, true); err != nil {
		return fmt.Errorf("declarar cola %s: %w", queueName, err)
	}
	if err := c.queue.Consume(ctx, queueName, c.handleAlert); err != nil {
		return fmt.Errorf("consumir cola %s: %w", queueName, err)
	}
	c.log.Info(ctx).Str("queue", queueName).Msg("[ai.assistant] consumidor de alertas iniciado")
	return nil
}

func (c *Consumer) handleAlert(body []byte) error {
	ctx := context.Background()

	var envelope alertEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		c.log.Warn(ctx).Err(err).Msg("[ai.assistant] alerta malformada (ACK)")
		return nil
	}
	if envelope.ID == "" || envelope.BusinessID == 0 {
		c.log.Warn(ctx).Str("event_type", envelope.Type).Msg("[ai.assistant] alerta sin id o sin negocio (ACK)")
		return nil
	}

	err := c.uc.IngestAlertEvent(ctx, dtos.AlertEvent{
		ID:         envelope.ID,
		Type:       envelope.Type,
		BusinessID: envelope.BusinessID,
		Timestamp:  envelope.Timestamp,
		Data:       envelope.Data,
	})
	if err != nil {
		if isPermanent(err) {
			c.log.Warn(ctx).Err(err).Str("event_id", envelope.ID).Msg("[ai.assistant] alerta descartada por error permanente (ACK)")
			return nil
		}
		c.log.Error(ctx).Err(err).Str("event_id", envelope.ID).Msg("[ai.assistant] no se pudo guardar la alerta, se reintenta")
		return err
	}
	return nil
}
