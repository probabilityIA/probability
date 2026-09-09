package consumer

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/push/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/push/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/push/internal/infra/primary/queue/consumer/request"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func (c *consumer) Start(ctx context.Context) error {
	queueName := rabbitmq.QueuePushNotificationRequests

	if err := c.queue.DeclareQueue(queueName, true); err != nil {
		c.log.Error().Err(err).Str("queue", queueName).Msg("no se pudo declarar la cola de push")
		return err
	}

	go func() {
		if err := c.queue.Consume(ctx, queueName, c.handleMessage); err != nil {
			c.log.Error().Err(err).Str("queue", queueName).Msg("error consumiendo la cola de push")
		}
	}()

	return nil
}

func (c *consumer) handleMessage(messageBody []byte) error {
	ctx := context.Background()

	var event request.PushEvent
	if err := json.Unmarshal(messageBody, &event); err != nil {
		c.log.Warn().Err(err).Msg("mensaje de push malformado, se descarta")
		return nil
	}

	if event.BusinessID == 0 {
		c.log.Warn().Str("event_type", event.EventType).Msg("evento de push sin business, se descarta")
		return nil
	}

	dto := dtos.PushEventDTO{
		EventType:     event.EventType,
		EventCategory: event.EventCategory,
		BusinessID:    event.BusinessID,
		OrderNumber:   event.OrderNumber,
		OrderID:       event.OrderID,
		Carrier:       event.Carrier,
		TrackingNo:    event.TrackingNo,
		StatusCode:    event.StatusCode,
		StatusName:    event.StatusName,
		CustomerName:  event.CustomerName,
		CodTotal:      event.CodTotal,
		Balance:       event.Balance,
	}

	if err := c.useCase.NotifyBusiness(ctx, dto); err != nil {
		if domainerrors.IsNonRetryable(err) {
			c.log.Warn().
				Err(err).
				Uint("business_id", event.BusinessID).
				Str("event_type", event.EventType).
				Msg("push descartado por error permanente")
			return nil
		}
		c.log.Error().
			Err(err).
			Uint("business_id", event.BusinessID).
			Msg("error enviando push, se reintenta")
		return err
	}

	return nil
}
