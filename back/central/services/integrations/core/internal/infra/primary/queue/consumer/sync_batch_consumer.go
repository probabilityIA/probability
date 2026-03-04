package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const defaultBatchDelay = 5 * time.Second

// ProviderLookupFn resuelve un provider por integration_type_id.
type ProviderLookupFn func(integrationTypeID int) (domain.IIntegrationContract, bool)

// SyncBatchConsumer procesa lotes de sincronización desde la cola integration.sync.batches.
type SyncBatchConsumer struct {
	queue      rabbitmq.IQueue
	providerFn ProviderLookupFn
	log        log.ILogger
	batchDelay time.Duration
}

// NewSyncBatchConsumer crea un nuevo consumer de lotes de sincronización.
func NewSyncBatchConsumer(queue rabbitmq.IQueue, providerFn ProviderLookupFn, logger log.ILogger) *SyncBatchConsumer {
	return &SyncBatchConsumer{
		queue:      queue,
		providerFn: providerFn,
		log:        logger.WithModule("sync_batch_consumer"),
		batchDelay: defaultBatchDelay,
	}
}

// Start declara la cola y registra el consumer.
func (c *SyncBatchConsumer) Start(ctx context.Context) error {
	if err := c.queue.DeclareQueue(rabbitmq.QueueSyncBatches, true); err != nil {
		return fmt.Errorf("error al declarar cola %s: %w", rabbitmq.QueueSyncBatches, err)
	}

	c.log.Info(ctx).
		Str("queue", rabbitmq.QueueSyncBatches).
		Msg("🚀 SyncBatchConsumer iniciado")

	return c.queue.Consume(ctx, rabbitmq.QueueSyncBatches, c.handleMessage)
}

func (c *SyncBatchConsumer) handleMessage(body []byte) error {
	ctx := context.Background()
	ctx = log.WithFunctionCtx(ctx, "SyncBatchConsumer.handleMessage")

	var msg domain.SyncBatchMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		c.log.Error(ctx).Err(err).Msg("Error al deserializar mensaje de lote — se omite (no se reencola)")
		// Retornar nil para ACK — mensajes corruptos no se pueden reintentar
		return nil
	}

	c.log.Info(ctx).
		Str("job_id", msg.JobID).
		Str("integration_id", msg.IntegrationID).
		Int("batch_index", msg.BatchIndex).
		Int("total_batches", msg.TotalBatches).
		Time("date_from", msg.CreatedAtMin).
		Time("date_to", msg.CreatedAtMax).
		Msg("📦 Procesando lote de sincronización")

	// Resolver IDs para los eventos SSE (necesarios antes del evento batch.processing)
	var businessID uint
	if msg.BusinessID != nil {
		businessID = *msg.BusinessID
	}
	integrationIDUint, _ := strconv.ParseUint(msg.IntegrationID, 10, 64)

	// Resolver provider
	provider, ok := c.providerFn(msg.IntegrationTypeID)
	if !ok {
		c.log.Error(ctx).
			Int("integration_type_id", msg.IntegrationTypeID).
			Msg("No hay provider registrado para este tipo de integración — se omite lote")
		return nil
	}

	// Publicar evento SSE: este lote está siendo procesado.
	// El frontend usa esto para completar el lote anterior y asignar órdenes al lote correcto.
	rabbitmq.PublishEvent(ctx, c.queue, rabbitmq.EventEnvelope{ //nolint:errcheck
		Type:          "integration.sync.batch.processing",
		Category:      "integration",
		BusinessID:    businessID,
		IntegrationID: uint(integrationIDUint),
		Data: map[string]interface{}{
			"job_id":        msg.JobID,
			"batch_index":   msg.BatchIndex,
			"total_batches": msg.TotalBatches,
		},
	})

	// Construir params como map[string]interface{} para el provider
	params := map[string]interface{}{
		"created_at_min": msg.CreatedAtMin,
		"created_at_max": msg.CreatedAtMax,
	}
	if msg.Status != "" {
		params["status"] = msg.Status
	}
	if msg.FinancialStatus != "" {
		params["financial_status"] = msg.FinancialStatus
	}
	if msg.FulfillmentStatus != "" {
		params["fulfillment_status"] = msg.FulfillmentStatus
	}

	startTime := time.Now()

	// Ejecutar sincronización del lote
	err := provider.SyncOrdersByIntegrationIDWithParams(ctx, msg.IntegrationID, params)
	duration := time.Since(startTime)

	if err != nil {
		c.log.Error(ctx).
			Err(err).
			Str("job_id", msg.JobID).
			Int("batch_index", msg.BatchIndex).
			Dur("duration", duration).
			Msg("❌ Lote de sincronización falló")

		// Publicar evento SSE de fallo
		rabbitmq.PublishEvent(ctx, c.queue, rabbitmq.EventEnvelope{ //nolint:errcheck
			Type:          "integration.sync.batch.failed",
			Category:      "integration",
			BusinessID:    businessID,
			IntegrationID: uint(integrationIDUint),
			Data: map[string]interface{}{
				"job_id":         msg.JobID,
				"batch_index":    msg.BatchIndex,
				"total_batches":  msg.TotalBatches,
				"error":          err.Error(),
				"duration":       duration.String(),
				"created_at_min": msg.CreatedAtMin.Format(time.RFC3339),
				"created_at_max": msg.CreatedAtMax.Format(time.RFC3339),
			},
		})

		// Retornar nil para ACK — no reencolar lotes individuales (el provider ya reintentó internamente)
		return nil
	}

	c.log.Info(ctx).
		Str("job_id", msg.JobID).
		Int("batch_index", msg.BatchIndex).
		Int("total_batches", msg.TotalBatches).
		Dur("duration", duration).
		Msg("✅ Lote de sincronización completado")

	// Publicar evento SSE de completado
	rabbitmq.PublishEvent(ctx, c.queue, rabbitmq.EventEnvelope{ //nolint:errcheck
		Type:          "integration.sync.batch.completed",
		Category:      "integration",
		BusinessID:    businessID,
		IntegrationID: uint(integrationIDUint),
		Data: map[string]interface{}{
			"job_id":         msg.JobID,
			"batch_index":    msg.BatchIndex,
			"total_batches":  msg.TotalBatches,
			"duration":       duration.String(),
			"created_at_min": msg.CreatedAtMin.Format(time.RFC3339),
			"created_at_max": msg.CreatedAtMax.Format(time.RFC3339),
		},
	})

	// Backpressure: esperar antes de procesar el siguiente lote
	time.Sleep(c.batchDelay)

	return nil
}
