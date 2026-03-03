package consumer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	queueMappers "github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/queue/mappers"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/queue/messages"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const QueueBulkInvoiceJobs = rabbitmq.QueueInvoicingBulkCreate

// BulkInvoiceConsumer procesa mensajes de facturación masiva
type BulkInvoiceConsumer struct {
	queue        rabbitmq.IQueue
	useCase      ports.IUseCase
	repo         ports.IRepository
	ssePublisher ports.IInvoiceSSEPublisher
	log          log.ILogger
}

// NewBulkInvoiceConsumer crea un nuevo consumer de facturación masiva
func NewBulkInvoiceConsumer(
	queue rabbitmq.IQueue,
	useCase ports.IUseCase,
	repo ports.IRepository,
	ssePublisher ports.IInvoiceSSEPublisher,
	logger log.ILogger,
) *BulkInvoiceConsumer {
	return &BulkInvoiceConsumer{
		queue:        queue,
		useCase:      useCase,
		repo:         repo,
		ssePublisher: ssePublisher,
		log:          logger.WithModule("invoicing.bulk_consumer"),
	}
}

// Start inicia el consumer
func (c *BulkInvoiceConsumer) Start(ctx context.Context) error {
	// Declarar cola (idempotente)
	c.queue.DeclareQueue(QueueBulkInvoiceJobs, true)

	// Consumir mensajes
	return c.queue.Consume(ctx, QueueBulkInvoiceJobs, c.handleMessage)
}

// handleMessage procesa un mensaje individual
func (c *BulkInvoiceConsumer) handleMessage(message []byte) error {
	ctx := context.Background()

	// 1. Deserializar mensaje de queue
	var queueMsg messages.BulkInvoiceJobMessage
	if err := json.Unmarshal(message, &queueMsg); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to unmarshal bulk invoice job message")
		return nil // No requeue - mensaje corrupto
	}

	// 2. Convertir mensaje de queue a DTO de dominio
	msg := queueMappers.BulkJobMessageToDTO(&queueMsg)

	c.log.Info(ctx).
		Str("job_id", msg.JobID).
		Str("order_id", msg.OrderID).
		Int("attempt", msg.AttemptNumber).
		Msg("Processing bulk invoice message")

	// 3. Marcar item como processing
	if err := c.updateItemStatus(ctx, msg.JobID, msg.OrderID, entities.JobItemStatusProcessing, nil); err != nil {
		c.log.Warn(ctx).Err(err).Str("job_id", msg.JobID).Str("order_id", msg.OrderID).Msg("Failed to update item to processing")
		// Continuar de todos modos
	}

	// 4. Inyectar business_id en contexto (necesario para CreateInvoice)
	ctx = context.WithValue(ctx, "business_id", msg.BusinessID)
	if msg.CreatedBy != nil {
		ctx = context.WithValue(ctx, "user_id", *msg.CreatedBy)
	}

	// 5. Crear factura usando el use case existente
	invoice, err := c.useCase.CreateInvoice(ctx, &dtos.CreateInvoiceDTO{
		OrderID:  msg.OrderID,
		IsManual: msg.IsManual,
	})

	// 6. Actualizar según resultado
	if err != nil {
		// Factura falló
		c.handleInvoiceError(ctx, msg.JobID, msg.OrderID, err)
		// NO retornar error - no hacer requeue (RetryConsumer maneja reintentos de provider)
		return nil
	}

	// 7. Solicitud de factura enviada al proveedor (aún pendiente de confirmación)
	c.handleInvoiceSuccess(ctx, msg.JobID, msg.OrderID, invoice.ID)

	// 8. NO verificar job completion aquí - el response_consumer lo hará
	// cuando reciba la respuesta de Softpymes (success o failure)

	return nil
}

// handleInvoiceError maneja cuando la creación de factura falla
func (c *BulkInvoiceConsumer) handleInvoiceError(ctx context.Context, jobID, orderID string, err error) {
	errMsg := err.Error()

	// Buscar invoice fallida en BD (puede existir con status "failed")
	invoice, getErr := c.repo.GetInvoiceByOrderID(ctx, orderID)
	if getErr == nil && invoice != nil {
		// Publicar evento individual de factura fallida
		if pubErr := c.ssePublisher.PublishInvoiceFailed(ctx, invoice, errMsg); pubErr != nil {
			c.log.Error(ctx).Err(pubErr).Msg("Failed to publish invoice.failed SSE")
		}
	} else {
		c.log.Warn(ctx).
			Err(getErr).
			Str("order_id", orderID).
			Msg("Invoice not found in DB after failure - skipping individual SSE event")
	}

	// Actualizar item como failed
	if updateErr := c.updateItemStatus(ctx, jobID, orderID, entities.JobItemStatusFailed, &errMsg); updateErr != nil {
		c.log.Error(ctx).Err(updateErr).Str("job_id", jobID).Str("order_id", orderID).Msg("Failed to update item to failed")
	}

	// Incrementar contadores del job (processed +1, failed +1)
	if incrementErr := c.repo.IncrementJobCounters(ctx, jobID, 1, 0, 1); incrementErr != nil {
		c.log.Error(ctx).Err(incrementErr).Str("job_id", jobID).Msg("Failed to increment job counters for failed invoice")
	}

	// Publicar progreso SSE
	c.publishJobProgress(ctx, jobID)

	// Para errores de pre-procesamiento (validación, BD, cola), verificar completion
	// ya que estas facturas NO pasarán por el response_consumer
	c.checkJobCompletion(ctx, jobID)

	c.log.Warn(ctx).
		Err(err).
		Str("job_id", jobID).
		Str("order_id", orderID).
		Msg("Invoice creation failed in bulk job")
}

// handleInvoiceSuccess maneja cuando la solicitud de factura fue enviada al proveedor
// NOTA: En este punto la factura está en estado "pending" - aún NO fue confirmada por Softpymes.
// El response_consumer publicará el SSE correcto (invoice.created o invoice.failed) cuando
// Softpymes responda, y también actualizará los contadores del job (successful/failed).
func (c *BulkInvoiceConsumer) handleInvoiceSuccess(ctx context.Context, jobID, orderID string, invoiceID uint) {
	// NO publicar invoice.created SSE aquí - la factura aún está pendiente de confirmación.
	// El response_consumer se encarga del SSE cuando Softpymes confirme o falle.

	// Actualizar item como processing (enviado al proveedor, esperando respuesta)
	now := time.Now()
	items, err := c.repo.GetJobItems(ctx, jobID)
	if err != nil {
		c.log.Error(ctx).Err(err).Str("job_id", jobID).Msg("Failed to get job items")
		return
	}

	// Encontrar el item específico y vincularlo con la factura
	for _, item := range items {
		if item.OrderID == orderID {
			item.Status = entities.JobItemStatusProcessing
			item.InvoiceID = &invoiceID
			item.ProcessedAt = &now
			item.ErrorMessage = nil

			if updateErr := c.repo.UpdateJobItem(ctx, item); updateErr != nil {
				c.log.Error(ctx).Err(updateErr).Str("job_id", jobID).Str("order_id", orderID).Msg("Failed to update item to processing")
			}
			break
		}
	}

	// Solo incrementar "processed" (NO "successful" - eso lo hará el response_consumer)
	if incrementErr := c.repo.IncrementJobCounters(ctx, jobID, 1, 0, 0); incrementErr != nil {
		c.log.Error(ctx).Err(incrementErr).Str("job_id", jobID).Msg("Failed to increment job counters")
	}

	// Publicar progreso SSE
	c.publishJobProgress(ctx, jobID)

	c.log.Info(ctx).
		Str("job_id", jobID).
		Str("order_id", orderID).
		Uint("invoice_id", invoiceID).
		Msg("Invoice created successfully in bulk job")
}

// updateItemStatus actualiza el estado de un item
func (c *BulkInvoiceConsumer) updateItemStatus(ctx context.Context, jobID, orderID, status string, errorMsg *string) error {
	items, err := c.repo.GetJobItems(ctx, jobID)
	if err != nil {
		return err
	}

	// Encontrar el item
	for _, item := range items {
		if item.OrderID == orderID {
			item.Status = status
			item.ErrorMessage = errorMsg

			if status == entities.JobItemStatusProcessing {
				now := time.Now()
				item.ProcessedAt = &now
			}

			return c.repo.UpdateJobItem(ctx, item)
		}
	}

	c.log.Warn(ctx).
		Str("job_id", jobID).
		Str("order_id", orderID).
		Msg("Job item not found for update")

	return nil
}

// checkJobCompletion verifica si el job completó y actualiza su estado
func (c *BulkInvoiceConsumer) checkJobCompletion(ctx context.Context, jobID string) {
	job, err := c.repo.GetJobByID(ctx, jobID)
	if err != nil {
		c.log.Error(ctx).Err(err).Str("job_id", jobID).Msg("Failed to get job for completion check")
		return
	}

	if job == nil {
		return
	}

	// Verificar si todos los items fueron procesados
	if job.Processed >= job.TotalOrders {
		now := time.Now()
		job.Status = entities.JobStatusCompleted
		job.CompletedAt = &now

		if updateErr := c.repo.UpdateJob(ctx, job); updateErr != nil {
			c.log.Error(ctx).Err(updateErr).Str("job_id", jobID).Msg("Failed to mark job as completed")
			return
		}

		// Publicar evento SSE de job completado
		if pubErr := c.ssePublisher.PublishBulkJobCompleted(ctx, job); pubErr != nil {
			c.log.Error(ctx).Err(pubErr).Str("job_id", jobID).Msg("Failed to publish bulk job completed SSE event")
		}

		c.log.Info(ctx).
			Str("job_id", jobID).
			Int("successful", job.Successful).
			Int("failed", job.Failed).
			Int("total", job.TotalOrders).
			Msg("Bulk invoice job completed")
	}
}

// publishJobProgress obtiene el estado actual del job y publica progreso SSE
func (c *BulkInvoiceConsumer) publishJobProgress(ctx context.Context, jobID string) {
	job, err := c.repo.GetJobByID(ctx, jobID)
	if err != nil || job == nil {
		return
	}
	if pubErr := c.ssePublisher.PublishBulkJobProgress(ctx, job); pubErr != nil {
		c.log.Error(ctx).Err(pubErr).Str("job_id", jobID).Msg("Failed to publish bulk job progress SSE event")
	}
}
