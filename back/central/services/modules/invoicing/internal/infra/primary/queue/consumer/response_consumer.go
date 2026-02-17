package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/constants"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// ResponseConsumer consume responses de proveedores de facturación
type ResponseConsumer struct {
	queue        rabbitmq.IQueue
	repo         ports.IRepository
	ssePublisher ports.IInvoiceSSEPublisher
	eventPub     ports.IEventPublisher
	log          log.ILogger
}

// NewResponseConsumer crea un nuevo consumer de responses
func NewResponseConsumer(
	queue rabbitmq.IQueue,
	repo ports.IRepository,
	ssePublisher ports.IInvoiceSSEPublisher,
	eventPub ports.IEventPublisher,
	logger log.ILogger,
) *ResponseConsumer {
	return &ResponseConsumer{
		queue:        queue,
		repo:         repo,
		ssePublisher: ssePublisher,
		eventPub:     eventPub,
		log:          logger.WithModule("invoicing.response_consumer"),
	}
}

const (
	QueueInvoiceResponses = "invoicing.responses"
)

// Start inicia el consumo de responses de proveedores
func (c *ResponseConsumer) Start(ctx context.Context) error {
	// Declarar la cola si no existe
	if err := c.queue.DeclareQueue(QueueInvoiceResponses, true); err != nil {
		c.log.Error(ctx).Err(err).Msg("Error al declarar cola de responses")
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	c.log.Info(ctx).
		Str("queue", QueueInvoiceResponses).
		Msg("📥 Starting invoice response consumer")

	// Iniciar consumo
	if err := c.queue.Consume(ctx, QueueInvoiceResponses, c.handleResponse); err != nil {
		c.log.Error(ctx).Err(err).Msg("Error al iniciar consumer de responses")
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	return nil
}

// handleResponse procesa una respuesta del proveedor
func (c *ResponseConsumer) handleResponse(message []byte) error {
	ctx := context.Background()

	// Deserializar response
	var response dtos.InvoiceResponseMessage
	if err := json.Unmarshal(message, &response); err != nil {
		c.log.Error(ctx).Err(err).Msg("Error al deserializar response")
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	c.log.Info(ctx).
		Uint("invoice_id", response.InvoiceID).
		Str("provider", response.Provider).
		Str("status", response.Status).
		Str("correlation_id", response.CorrelationID).
		Msg("📨 Processing invoice response")

	// Obtener factura
	invoice, err := c.repo.GetInvoiceByID(ctx, response.InvoiceID)
	if err != nil {
		c.log.Error(ctx).
			Err(err).
			Uint("invoice_id", response.InvoiceID).
			Msg("Failed to get invoice")
		return nil // No requeue - invoice no existe
	}

	// Obtener sync log actual (el más reciente en processing)
	syncLogs, err := c.repo.GetSyncLogsByInvoiceID(ctx, response.InvoiceID)
	if err != nil || len(syncLogs) == 0 {
		c.log.Warn(ctx).
			Uint("invoice_id", response.InvoiceID).
			Msg("No sync logs found")
		// Continuar sin sync log (no es crítico)
	}

	var syncLog *entities.InvoiceSyncLog
	if len(syncLogs) > 0 {
		// Usar el sync log más reciente directamente
		syncLog = syncLogs[0]
	}

	// Procesar según status
	if response.Status == dtos.ResponseStatusSuccess {
		c.handleSuccess(ctx, invoice, syncLog, &response)
	} else {
		c.handleError(ctx, invoice, syncLog, &response)
	}

	return nil
}

// handleSuccess procesa una respuesta exitosa
func (c *ResponseConsumer) handleSuccess(
	ctx context.Context,
	invoice *entities.Invoice,
	syncLog *entities.InvoiceSyncLog,
	response *dtos.InvoiceResponseMessage,
) {
	c.log.Info(ctx).
		Uint("invoice_id", invoice.ID).
		Str("invoice_number", response.InvoiceNumber).
		Msg("✅ Invoice created successfully by provider")

	// Actualizar invoice con datos del proveedor
	invoice.InvoiceNumber = response.InvoiceNumber
	if response.ExternalID != "" {
		invoice.ExternalID = &response.ExternalID
	}
	if response.InvoiceURL != "" {
		invoice.InvoiceURL = &response.InvoiceURL
	}
	if response.PDFURL != "" {
		invoice.PDFURL = &response.PDFURL
	}
	if response.XMLURL != "" {
		invoice.XMLURL = &response.XMLURL
	}
	if response.CUFE != "" {
		invoice.CUFE = &response.CUFE
	}
	if response.IssuedAt != nil {
		invoice.IssuedAt = response.IssuedAt
	}

	invoice.Status = constants.InvoiceStatusIssued

	// Guardar documento completo si existe
	if response.DocumentJSON != nil {
		invoice.ProviderResponse = response.DocumentJSON
	}

	// Actualizar en BD
	if err := c.repo.UpdateInvoice(ctx, invoice); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to update invoice")
		return
	}

	// Actualizar sync log como success
	if syncLog != nil {
		completedAt := time.Now()
		duration := int(completedAt.Sub(syncLog.StartedAt).Milliseconds())
		syncLog.Status = constants.SyncStatusSuccess
		syncLog.CompletedAt = &completedAt
		syncLog.Duration = &duration

		// Guardar response completa
		syncLog.ResponseBody = response.DocumentJSON

		// Poblar audit data del request/response HTTP al proveedor
		c.populateSyncLogAudit(syncLog, response)

		if err := c.repo.UpdateInvoiceSyncLog(ctx, syncLog); err != nil {
			c.log.Error(ctx).Err(err).Msg("Failed to update sync log")
		}
	}

	// Actualizar información de factura en la orden
	invoiceURL := ""
	if invoice.InvoiceURL != nil {
		invoiceURL = *invoice.InvoiceURL
	}
	if err := c.repo.UpdateOrderInvoiceInfo(ctx, invoice.OrderID, invoice.InvoiceNumber, invoiceURL); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to update order invoice info")
	}

	// Publicar evento de factura creada (RabbitMQ)
	if err := c.eventPub.PublishInvoiceCreated(ctx, invoice); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to publish invoice created event")
	}

	// Publicar evento SSE
	if err := c.ssePublisher.PublishInvoiceCreated(ctx, invoice); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to publish SSE event")
	}

	// Actualizar contadores de bulk job si la factura pertenece a uno
	c.updateBulkJobOnResult(ctx, invoice.ID, true)

	c.log.Info(ctx).
		Uint("invoice_id", invoice.ID).
		Str("invoice_number", invoice.InvoiceNumber).
		Int64("processing_time_ms", response.ProcessingTime).
		Msg("✅ Invoice response processed successfully")
}

// handleError procesa una respuesta de error
func (c *ResponseConsumer) handleError(
	ctx context.Context,
	invoice *entities.Invoice,
	syncLog *entities.InvoiceSyncLog,
	response *dtos.InvoiceResponseMessage,
) {
	c.log.Error(ctx).
		Uint("invoice_id", invoice.ID).
		Str("error", response.Error).
		Str("error_code", response.ErrorCode).
		Msg("❌ Provider returned error")

	// Marcar invoice como failed
	invoice.Status = constants.InvoiceStatusFailed

	if err := c.repo.UpdateInvoice(ctx, invoice); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to update invoice status to failed")
		return
	}

	// Actualizar sync log como failed
	if syncLog != nil {
		completedAt := time.Now()
		duration := int(completedAt.Sub(syncLog.StartedAt).Milliseconds())
		syncLog.Status = constants.SyncStatusFailed
		syncLog.CompletedAt = &completedAt
		syncLog.Duration = &duration
		syncLog.ErrorMessage = &response.Error
		syncLog.ErrorCode = &response.ErrorCode

		// Guardar detalles del error
		if response.ErrorDetails != nil {
			syncLog.ErrorDetails = response.ErrorDetails
		}

		// Poblar audit data del request/response HTTP al proveedor
		c.populateSyncLogAudit(syncLog, response)

		// Calcular próximo reintento si no se excedió el límite
		if syncLog.RetryCount < syncLog.MaxRetries {
			nextRetry := c.calculateNextRetry(syncLog.RetryCount)
			syncLog.NextRetryAt = &nextRetry
		}

		if err := c.repo.UpdateInvoiceSyncLog(ctx, syncLog); err != nil {
			c.log.Error(ctx).Err(err).Msg("Failed to update sync log")
		}
	}

	// Publicar evento de factura fallida
	if err := c.eventPub.PublishInvoiceFailed(ctx, invoice, response.Error); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to publish invoice failed event")
	}

	// Publicar evento SSE
	if err := c.ssePublisher.PublishInvoiceFailed(ctx, invoice, response.Error); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to publish SSE failed event")
	}

	// Actualizar contadores de bulk job si la factura pertenece a uno
	c.updateBulkJobOnResult(ctx, invoice.ID, false)
}

// updateBulkJobOnResult actualiza los contadores del bulk job cuando se recibe el resultado del proveedor
func (c *ResponseConsumer) updateBulkJobOnResult(ctx context.Context, invoiceID uint, success bool) {
	// Buscar si esta factura pertenece a un bulk job
	jobItem, err := c.repo.GetJobItemByInvoiceID(ctx, invoiceID)
	if err != nil {
		c.log.Warn(ctx).Err(err).Uint("invoice_id", invoiceID).Msg("Error checking bulk job item")
		return
	}
	if jobItem == nil {
		return // No pertenece a un bulk job
	}

	// Actualizar estado del item
	if success {
		jobItem.Status = "success"
	} else {
		jobItem.Status = "failed"
	}
	if updateErr := c.repo.UpdateJobItem(ctx, jobItem); updateErr != nil {
		c.log.Error(ctx).Err(updateErr).Msg("Failed to update bulk job item status")
	}

	// Incrementar contadores del job
	successful, failed := 0, 0
	if success {
		successful = 1
	} else {
		failed = 1
	}
	if incrementErr := c.repo.IncrementJobCounters(ctx, jobItem.JobID, 0, successful, failed); incrementErr != nil {
		c.log.Error(ctx).Err(incrementErr).Str("job_id", jobItem.JobID).Msg("Failed to increment bulk job counters")
		return
	}

	// Publicar progreso SSE del job
	job, err := c.repo.GetJobByID(ctx, jobItem.JobID)
	if err != nil || job == nil {
		return
	}

	if pubErr := c.ssePublisher.PublishBulkJobProgress(ctx, job); pubErr != nil {
		c.log.Error(ctx).Err(pubErr).Str("job_id", jobItem.JobID).Msg("Failed to publish bulk job progress SSE")
	}

	// Verificar si el job completó (successful + failed = total)
	if job.Successful+job.Failed >= job.TotalOrders {
		c.completeBulkJob(ctx, job)
	}
}

// completeBulkJob marca un bulk job como completado
func (c *ResponseConsumer) completeBulkJob(ctx context.Context, job *entities.BulkInvoiceJob) {
	now := time.Now()
	job.Status = "completed"
	job.CompletedAt = &now

	if updateErr := c.repo.UpdateJob(ctx, job); updateErr != nil {
		c.log.Error(ctx).Err(updateErr).Str("job_id", job.ID).Msg("Failed to mark bulk job as completed")
		return
	}

	if pubErr := c.ssePublisher.PublishBulkJobCompleted(ctx, job); pubErr != nil {
		c.log.Error(ctx).Err(pubErr).Str("job_id", job.ID).Msg("Failed to publish bulk job completed SSE")
	}

	c.log.Info(ctx).
		Str("job_id", job.ID).
		Int("successful", job.Successful).
		Int("failed", job.Failed).
		Int("total", job.TotalOrders).
		Msg("Bulk invoice job completed (from response consumer)")
}

// populateSyncLogAudit extrae audit data del response message y la almacena en el sync log
func (c *ResponseConsumer) populateSyncLogAudit(syncLog *entities.InvoiceSyncLog, response *dtos.InvoiceResponseMessage) {
	if response.AuditRequestURL != "" {
		syncLog.RequestURL = response.AuditRequestURL
	}
	if response.AuditRequestPayload != nil {
		syncLog.RequestPayload = response.AuditRequestPayload
	}
	if response.AuditResponseStatus != 0 {
		syncLog.ResponseStatus = response.AuditResponseStatus
	}
	if response.AuditResponseBody != "" {
		var bodyMap map[string]interface{}
		if json.Unmarshal([]byte(response.AuditResponseBody), &bodyMap) == nil {
			syncLog.ResponseBody = bodyMap
		}
	}
}

// calculateNextRetry calcula el próximo intento (exponential backoff)
func (c *ResponseConsumer) calculateNextRetry(retryCount int) time.Time {
	// Backoff exponencial: 5min, 15min, 30min
	delays := []time.Duration{
		5 * time.Minute,
		15 * time.Minute,
		30 * time.Minute,
	}

	delayIndex := retryCount
	if delayIndex >= len(delays) {
		delayIndex = len(delays) - 1
	}

	return time.Now().Add(delays[delayIndex])
}
