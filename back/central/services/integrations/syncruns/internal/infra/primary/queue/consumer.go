package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const queueName = "integrations.sync_runs"

var inventoryRoutingKeys = []string{
	"shopify.inventory.sync.completed",
	"meli.inventory.sync.completed",
	"woo.inventory.sync.completed",
	"jumpseller.inventory.sync.completed",
	"vtex.inventory.sync.completed",
}

type envelope struct {
	Type          string                 `json:"type"`
	BusinessID    uint                   `json:"business_id"`
	IntegrationID uint                   `json:"integration_id"`
	Timestamp     time.Time              `json:"timestamp"`
	Data          map[string]interface{} `json:"data"`
}

type Consumer struct {
	queue   rabbitmq.IQueue
	useCase domain.IUseCase
	logger  log.ILogger
}

func New(queue rabbitmq.IQueue, useCase domain.IUseCase, logger log.ILogger) *Consumer {
	return &Consumer{
		queue:   queue,
		useCase: useCase,
		logger:  logger.WithModule("syncruns"),
	}
}

func (c *Consumer) Start(ctx context.Context) {
	if c.queue == nil {
		return
	}

	if err := c.queue.DeclareQueue(queueName, true); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Error al declarar la cola de resultados de sincronizacion")
		return
	}

	for _, key := range inventoryRoutingKeys {
		if err := c.queue.BindQueue(queueName, rabbitmq.ExchangeEvents, key); err != nil {
			c.logger.Error(ctx).Err(err).Str("routing_key", key).
				Msg("Error al enlazar la cola de resultados de sincronizacion")
			return
		}
	}

	go func() {
		if err := c.queue.Consume(ctx, queueName, func(body []byte) error {
			c.handle(ctx, body)
			return nil
		}); err != nil {
			c.logger.Error(ctx).Err(err).Msg("Error al consumir la cola de resultados de sincronizacion")
		}
	}()

	c.logger.Info(ctx).Msg("Consumer de resultados de sincronizacion iniciado")
}

func (c *Consumer) handle(ctx context.Context, body []byte) {
	var msg envelope
	if err := json.Unmarshal(body, &msg); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Evento de sincronizacion invalido")
		return
	}

	if msg.BusinessID == 0 || msg.IntegrationID == 0 {
		return
	}

	finished := msg.Timestamp
	if finished.IsZero() {
		finished = time.Now()
	}

	updated := intOf(msg.Data, "updated")
	if updated == 0 {
		updated = intOf(msg.Data, "synced")
	}
	failed := intOf(msg.Data, "failed")

	run := domain.SyncRun{
		BusinessID:    msg.BusinessID,
		IntegrationID: msg.IntegrationID,
		Kind:          domain.KindInventory,
		CorrelationID: stringOf(msg.Data, "correlation_id"),
		StartedAt:     finished,
		FinishedAt:    &finished,
		Total:         intOf(msg.Data, "total"),
		Updated:       updated,
		Unchanged:     intOf(msg.Data, "unchanged"),
		Skipped:       intOf(msg.Data, "skipped"),
		Failed:        failed,
		Status:        domain.StatusCompleted,
		Detail:        failedDetail(msg.Data),
	}

	if err := c.useCase.Record(ctx, &run); err != nil {
		c.logger.Warn(ctx).Err(err).
			Uint("integration_id", run.IntegrationID).
			Msg("No se guardo el resultado de la sincronizacion de inventario")
	}
}

func failedDetail(data map[string]interface{}) []domain.DetailItem {
	raw, ok := data["failed_skus"].([]interface{})
	if !ok {
		return nil
	}
	detail := make([]domain.DetailItem, 0, len(raw))
	for _, item := range raw {
		switch value := item.(type) {
		case string:
			detail = append(detail, domain.DetailItem{SKU: value, Label: "fallo al actualizar", Tone: "error"})
		case map[string]interface{}:
			sku, _ := value["sku"].(string)
			message, _ := value["error"].(string)
			if message == "" {
				message = "fallo al actualizar"
			}
			detail = append(detail, domain.DetailItem{SKU: sku, Label: message, Tone: "error"})
		}
	}
	return detail
}

func intOf(data map[string]interface{}, key string) int {
	switch value := data[key].(type) {
	case float64:
		return int(value)
	case int:
		return value
	}
	return 0
}

func stringOf(data map[string]interface{}, key string) string {
	value, _ := data[key].(string)
	return value
}
