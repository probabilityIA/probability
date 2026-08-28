package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/constants"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type compareItemDetail struct {
	ItemCode string `json:"item_code"`
	ItemName string `json:"item_name"`
	Quantity string `json:"quantity"`
	Value    string `json:"value"`
	IVA      string `json:"iva"`
}

type ResponseConsumer struct {
	queue        rabbitmq.IQueue
	repo         ports.IRepository
	ssePublisher ports.IInvoiceSSEPublisher
	eventPub     ports.IEventPublisher
	compareCache ports.ICompareCache
	log          log.ILogger
}

func NewResponseConsumer(
	queue rabbitmq.IQueue,
	repo ports.IRepository,
	ssePublisher ports.IInvoiceSSEPublisher,
	eventPub ports.IEventPublisher,
	compareCache ports.ICompareCache,
	logger log.ILogger,
) *ResponseConsumer {
	return &ResponseConsumer{
		queue:        queue,
		repo:         repo,
		ssePublisher: ssePublisher,
		eventPub:     eventPub,
		compareCache: compareCache,
		log:          logger.WithModule("invoicing.response_consumer"),
	}
}

const (
	QueueInvoiceResponses = rabbitmq.QueueInvoicingResponses
)

type compareProviderDocument struct {
	DocumentNumber     string              `json:"document_number"`
	DocumentDate       string              `json:"document_date"`
	DocumentName       string              `json:"document_name"`
	Total              string              `json:"total"`
	CustomerNit        string              `json:"customer_nit"`
	CustomerName       string              `json:"customer_name"`
	Comment            string              `json:"comment"`
	Prefix             string              `json:"prefix"`
	Annuled            bool                `json:"annuled"`
	ElectronicDocument bool                `json:"electronic_document"`
	Details            []compareItemDetail `json:"details,omitempty"`
}

type compareResponseMessage struct {
	Operation         string                    `json:"operation"`
	Mode              string                    `json:"mode,omitempty"`
	CorrelationID     string                    `json:"correlation_id"`
	BusinessID        uint                      `json:"business_id"`
	DateFrom          string                    `json:"date_from"`
	DateTo            string                    `json:"date_to"`
	ProviderDocuments []compareProviderDocument `json:"provider_documents"`
	Error             string                    `json:"error,omitempty"`
	Timestamp         time.Time                 `json:"timestamp"`
}

type responseDiscriminator struct {
	Operation string `json:"operation,omitempty"`
	InvoiceID uint   `json:"invoice_id"`
}

func (c *ResponseConsumer) Start(ctx context.Context) error {

	if err := c.queue.DeclareQueue(QueueInvoiceResponses, true); err != nil {
		c.log.Error(ctx).Err(err).Msg("Error al declarar cola de responses")
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	c.log.Info(ctx).
		Str("queue", QueueInvoiceResponses).
		Msg("Starting invoice response consumer")

	if err := c.queue.Consume(ctx, QueueInvoiceResponses, c.handleResponse); err != nil {
		c.log.Error(ctx).Err(err).Msg("Error al iniciar consumer de responses")
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	return nil
}

func (c *ResponseConsumer) handleResponse(message []byte) error {
	ctx := context.Background()

	var disc responseDiscriminator
	if err := json.Unmarshal(message, &disc); err != nil {
		c.log.Error(ctx).Err(err).Msg("Error al deserializar discriminator")
		return fmt.Errorf("failed to unmarshal discriminator: %w", err)
	}

	if disc.Operation == dtos.OperationCompare {
		return c.handleCompareResponse(ctx, message)
	}

	if disc.Operation == dtos.OperationListItems {
		return c.handleListItemsResponse(ctx, message)
	}

	if disc.Operation == dtos.OperationListBankAccounts {
		return c.handleListBankAccountsResponse(ctx, message)
	}

	if disc.Operation == dtos.OperationListSiigoWarehouses {
		return c.handleListSiigoWarehousesResponse(ctx, message)
	}

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
		Msg("Processing invoice response")

	invoice, err := c.repo.GetInvoiceByID(ctx, response.InvoiceID)
	if err != nil {
		c.log.Error(ctx).
			Err(err).
			Uint("invoice_id", response.InvoiceID).
			Msg("Failed to get invoice")
		return nil
	}

	syncLogs, err := c.repo.GetSyncLogsByInvoiceID(ctx, response.InvoiceID)
	if err != nil || len(syncLogs) == 0 {
		c.log.Warn(ctx).
			Uint("invoice_id", response.InvoiceID).
			Msg("No sync logs found")

	}

	var syncLog *entities.InvoiceSyncLog
	if len(syncLogs) > 0 {

		syncLog = syncLogs[0]
	}

	if response.Operation == dtos.OperationCancel {
		if response.Status == dtos.ResponseStatusSuccess {
			c.handleCancelSuccess(ctx, invoice, syncLog, &response)
		} else {
			c.handleCancelError(ctx, invoice, syncLog, &response)
		}
	} else if response.Operation == dtos.OperationCashReceipt {
		c.handleCashReceiptResponse(ctx, invoice, syncLog, &response)
	} else if response.Operation == dtos.OperationCreditNote {
		c.handleCreditNoteResponse(ctx, invoice, syncLog, &response)
	} else {
		switch response.Status {
		case dtos.ResponseStatusSuccess:
			c.handleSuccess(ctx, invoice, syncLog, &response)
		case dtos.ResponseStatusPendingValidation:
			c.handlePendingValidation(ctx, invoice, syncLog, &response)
		case dtos.ResponseStatusRetryRequeued:
			c.handleRetryRequeued(ctx, invoice, syncLog, &response)
		default:
			c.handleError(ctx, invoice, syncLog, &response)
		}
	}

	return nil
}

func (c *ResponseConsumer) handleSuccess(
	ctx context.Context,
	invoice *entities.Invoice,
	syncLog *entities.InvoiceSyncLog,
	response *dtos.InvoiceResponseMessage,
) {
	c.log.Info(ctx).
		Uint("invoice_id", invoice.ID).
		Str("invoice_number", response.InvoiceNumber).
		Msg("Invoice created successfully by provider")

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

	if response.DocumentJSON != nil {
		invoice.ProviderResponse = response.DocumentJSON
	}

	if err := c.repo.UpdateInvoice(ctx, invoice); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to update invoice")
		return
	}

	if syncLog != nil {
		completedAt := time.Now()
		duration := int(completedAt.Sub(syncLog.StartedAt).Milliseconds())
		syncLog.Status = constants.SyncStatusSuccess
		syncLog.CompletedAt = &completedAt
		syncLog.Duration = &duration

		syncLog.ResponseBody = response.DocumentJSON

		c.populateSyncLogAudit(syncLog, response)

		if err := c.repo.UpdateInvoiceSyncLog(ctx, syncLog); err != nil {
			c.log.Error(ctx).Err(err).Msg("Failed to update sync log")
		}
	}

	if response.Operation != dtos.OperationCreateJournal {
		invoiceURL := ""
		if invoice.InvoiceURL != nil {
			invoiceURL = *invoice.InvoiceURL
		}
		if err := c.repo.UpdateOrderInvoiceInfo(ctx, invoice.OrderID, invoice.InvoiceNumber, invoiceURL); err != nil {
			c.log.Error(ctx).Err(err).Msg("Failed to update order invoice info")
		}
	}

	if err := c.eventPub.PublishInvoiceCreated(ctx, invoice); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to publish invoice created event")
	}

	if err := c.ssePublisher.PublishInvoiceCreated(ctx, invoice); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to publish SSE event")
	}

	c.updateBulkJobOnResult(ctx, invoice.ID, true)

	c.log.Info(ctx).
		Uint("invoice_id", invoice.ID).
		Str("invoice_number", invoice.InvoiceNumber).
		Int64("processing_time_ms", response.ProcessingTime).
		Msg("Invoice response processed successfully")
}

func (c *ResponseConsumer) handlePendingValidation(
	ctx context.Context,
	invoice *entities.Invoice,
	syncLog *entities.InvoiceSyncLog,
	response *dtos.InvoiceResponseMessage,
) {
	c.log.Info(ctx).
		Uint("invoice_id", invoice.ID).
		Str("provider_message", response.Error).
		Msg("Invoice accepted by provider, pending DIAN validation")

	invoice.Status = constants.InvoiceStatusPending
	if err := c.repo.UpdateInvoice(ctx, invoice); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to update invoice status to pending")
		return
	}

	if syncLog != nil {
		completedAt := time.Now()
		duration := int(completedAt.Sub(syncLog.StartedAt).Milliseconds())
		syncLog.Status = constants.SyncStatusPending
		syncLog.CompletedAt = &completedAt
		syncLog.Duration = &duration
		providerMsg := response.Error
		syncLog.ErrorMessage = &providerMsg

		c.populateSyncLogAudit(syncLog, response)

		nextCheck := c.calculateCheckBackoff(syncLog.RetryCount)
		syncLog.NextRetryAt = &nextCheck
		c.log.Info(ctx).
			Uint("invoice_id", invoice.ID).
			Time("next_check_at", nextCheck).
			Int("check_count", syncLog.RetryCount).
			Msg("Scheduled next check_status for pending DIAN validation")

		if err := c.repo.UpdateInvoiceSyncLog(ctx, syncLog); err != nil {
			c.log.Error(ctx).Err(err).Msg("Failed to update sync log for pending validation")
		}
	}

	_ = rabbitmq.PublishEvent(ctx, c.queue, rabbitmq.EventEnvelope{
		Type:       "invoice.pending_validation",
		Category:   "invoice",
		BusinessID: invoice.BusinessID,
		Data: map[string]interface{}{
			"invoice_id":       invoice.ID,
			"order_id":         invoice.OrderID,
			"provider_message": response.Error,
		},
	})

	c.updateBulkJobOnResult(ctx, invoice.ID, true)

	c.log.Info(ctx).
		Uint("invoice_id", invoice.ID).
		Msg("Invoice kept as pending - awaiting DIAN validation")
}

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
		Msg("Provider returned error")

	invoice.Status = constants.InvoiceStatusFailed

	if err := c.repo.UpdateInvoice(ctx, invoice); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to update invoice status to failed")
		return
	}

	if syncLog != nil {
		completedAt := time.Now()
		duration := int(completedAt.Sub(syncLog.StartedAt).Milliseconds())
		syncLog.Status = constants.SyncStatusFailed
		syncLog.CompletedAt = &completedAt
		syncLog.Duration = &duration
		syncLog.ErrorMessage = &response.Error
		syncLog.ErrorCode = &response.ErrorCode

		if response.ErrorDetails != nil {
			syncLog.ErrorDetails = response.ErrorDetails
		}

		c.populateSyncLogAudit(syncLog, response)

		if esErrorDeConfiguracion(response.ErrorCode) {
			syncLog.NextRetryAt = nil
			c.log.Warn(ctx).
				Uint("invoice_id", invoice.ID).
				Str("error_code", response.ErrorCode).
				Msg("Error de configuracion en el proveedor: no se reintenta, hay que corregirlo antes de volver a facturar")
		} else if isProviderUnavailableError(response.Error) {
			if syncLog.RetryCount > 0 {
				syncLog.RetryCount--
			}
			elapsed := time.Since(invoice.CreatedAt)
			if elapsed >= providerUnavailableMaxWindow {
				syncLog.Status = constants.SyncStatusCancelled
				syncLog.NextRetryAt = nil
				c.log.Warn(ctx).
					Uint("invoice_id", invoice.ID).
					Dur("elapsed", elapsed).
					Dur("max_window", providerUnavailableMaxWindow).
					Msg("Provider unavailable window exhausted - automatic retries stopped, invoice needs manual review")
			} else {
				nextRetry := time.Now().Add(providerUnavailableRetryDelay(elapsed))
				syncLog.NextRetryAt = &nextRetry
				c.log.Warn(ctx).
					Uint("invoice_id", invoice.ID).
					Time("next_retry_at", nextRetry).
					Dur("elapsed", elapsed).
					Int("retry_count", syncLog.RetryCount).
					Msg("Provider unavailable - retry rescheduled without consuming retry budget")
			}
		} else if syncLog.RetryCount < syncLog.MaxRetries {
			nextRetry := c.calculateNextRetry(syncLog.RetryCount)
			syncLog.NextRetryAt = &nextRetry
		}

		if err := c.repo.UpdateInvoiceSyncLog(ctx, syncLog); err != nil {
			c.log.Error(ctx).Err(err).Msg("Failed to update sync log")
		}
	}

	if err := c.eventPub.PublishInvoiceFailed(ctx, invoice, response.Error); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to publish invoice failed event")
	}

	if err := c.ssePublisher.PublishInvoiceFailed(ctx, invoice, response.Error); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to publish SSE failed event")
	}

	c.updateBulkJobOnResult(ctx, invoice.ID, false)
}

func (c *ResponseConsumer) handleRetryRequeued(
	ctx context.Context,
	invoice *entities.Invoice,
	syncLog *entities.InvoiceSyncLog,
	response *dtos.InvoiceResponseMessage,
) {
	invoice.Status = constants.InvoiceStatusFailed
	if err := c.repo.UpdateInvoice(ctx, invoice); err != nil {
		c.log.Error(ctx).Err(err).Uint("invoice_id", invoice.ID).Msg("Failed to update invoice for requeue")
		return
	}

	if syncLog == nil {
		c.log.Warn(ctx).Uint("invoice_id", invoice.ID).Msg("No sync log to requeue")
		return
	}

	syncLog.Status = constants.SyncStatusFailed
	if syncLog.RetryCount >= syncLog.MaxRetries {
		syncLog.MaxRetries = syncLog.RetryCount + constants.MaxRetries
	}
	nextRetry := time.Now()
	syncLog.NextRetryAt = &nextRetry
	errMsg := response.Error
	if errMsg == "" {
		errMsg = "re-encolada por reintento masivo de fallidas"
	}
	syncLog.ErrorMessage = &errMsg

	if err := c.repo.UpdateInvoiceSyncLog(ctx, syncLog); err != nil {
		c.log.Error(ctx).Err(err).Uint("invoice_id", invoice.ID).Msg("Failed to requeue sync log")
		return
	}

	c.log.Info(ctx).
		Uint("invoice_id", invoice.ID).
		Int("retry_count", syncLog.RetryCount).
		Int("max_retries", syncLog.MaxRetries).
		Msg("Invoice requeued for creation retry after reconcile")
}

func (c *ResponseConsumer) handleCancelSuccess(
	ctx context.Context,
	invoice *entities.Invoice,
	syncLog *entities.InvoiceSyncLog,
	response *dtos.InvoiceResponseMessage,
) {
	c.log.Info(ctx).Uint("invoice_id", invoice.ID).Msg("Invoice cancelled successfully by provider")

	now := time.Now()
	invoice.Status = constants.InvoiceStatusCancelled
	invoice.CancelledAt = &now

	if err := c.repo.UpdateInvoice(ctx, invoice); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to update invoice status to cancelled")
		return
	}

	if syncLog != nil {
		completedAt := time.Now()
		duration := int(completedAt.Sub(syncLog.StartedAt).Milliseconds())
		syncLog.Status = constants.SyncStatusSuccess
		syncLog.CompletedAt = &completedAt
		syncLog.Duration = &duration
		if err := c.repo.UpdateInvoiceSyncLog(ctx, syncLog); err != nil {
			c.log.Error(ctx).Err(err).Msg("Failed to update cancel sync log")
		}
	}

	_ = rabbitmq.PublishEvent(ctx, c.queue, rabbitmq.EventEnvelope{
		Type:       "invoice.cancelled",
		Category:   "invoice",
		BusinessID: invoice.BusinessID,
		Data: map[string]interface{}{
			"invoice_id": invoice.ID,
			"order_id":   invoice.OrderID,
		},
	})

	c.log.Info(ctx).Uint("invoice_id", invoice.ID).Msg("Invoice cancellation processed successfully")
}

func (c *ResponseConsumer) handleCancelError(
	ctx context.Context,
	invoice *entities.Invoice,
	syncLog *entities.InvoiceSyncLog,
	response *dtos.InvoiceResponseMessage,
) {
	c.log.Error(ctx).
		Uint("invoice_id", invoice.ID).
		Str("error", response.Error).
		Msg("Provider returned error on cancellation")

	if syncLog != nil {
		completedAt := time.Now()
		duration := int(completedAt.Sub(syncLog.StartedAt).Milliseconds())
		syncLog.Status = constants.SyncStatusFailed
		syncLog.CompletedAt = &completedAt
		syncLog.Duration = &duration
		syncLog.ErrorMessage = &response.Error
		syncLog.ErrorCode = &response.ErrorCode
		if err := c.repo.UpdateInvoiceSyncLog(ctx, syncLog); err != nil {
			c.log.Error(ctx).Err(err).Msg("Failed to update cancel sync log")
		}
	}

	_ = rabbitmq.PublishEvent(ctx, c.queue, rabbitmq.EventEnvelope{
		Type:       "invoice.cancel_failed",
		Category:   "invoice",
		BusinessID: invoice.BusinessID,
		Data: map[string]interface{}{
			"invoice_id":    invoice.ID,
			"order_id":      invoice.OrderID,
			"error_message": response.Error,
		},
	})
}

func (c *ResponseConsumer) handleCashReceiptResponse(
	ctx context.Context,
	invoice *entities.Invoice,
	syncLog *entities.InvoiceSyncLog,
	response *dtos.InvoiceResponseMessage,
) {
	isSuccess := response.Status == dtos.ResponseStatusSuccess

	if isSuccess {
		c.log.Info(ctx).
			Uint("invoice_id", invoice.ID).
			Msg("Cash receipt generated successfully")

		if response.DocumentJSON != nil {
			invoice.ProviderResponse = response.DocumentJSON
			if err := c.repo.UpdateInvoice(ctx, invoice); err != nil {
				c.log.Error(ctx).Err(err).Msg("Failed to update invoice document_json after cash receipt")
			}
		}
	} else {
		c.log.Error(ctx).
			Uint("invoice_id", invoice.ID).
			Str("error", response.Error).
			Msg("Cash receipt generation failed")
	}

	if syncLog != nil {
		completedAt := time.Now()
		duration := int(completedAt.Sub(syncLog.StartedAt).Milliseconds())
		syncLog.CompletedAt = &completedAt
		syncLog.Duration = &duration

		if isSuccess {
			syncLog.Status = constants.SyncStatusSuccess
			syncLog.ResponseBody = response.DocumentJSON
		} else {
			syncLog.Status = constants.SyncStatusFailed
			syncLog.ErrorMessage = &response.Error
			syncLog.ErrorCode = &response.ErrorCode
		}

		c.populateSyncLogAudit(syncLog, response)

		if err := c.repo.UpdateInvoiceSyncLog(ctx, syncLog); err != nil {
			c.log.Error(ctx).Err(err).Msg("Failed to update cash receipt sync log")
		}
	}

	if isSuccess {
		_ = c.ssePublisher.PublishInvoiceCreated(ctx, invoice)
	} else {
		_ = c.ssePublisher.PublishInvoiceFailed(ctx, invoice, response.Error)
	}
}

func (c *ResponseConsumer) updateBulkJobOnResult(ctx context.Context, invoiceID uint, success bool) {

	jobItem, err := c.repo.GetJobItemByInvoiceID(ctx, invoiceID)
	if err != nil {
		c.log.Warn(ctx).Err(err).Uint("invoice_id", invoiceID).Msg("Error checking bulk job item")
		return
	}
	if jobItem == nil {
		return
	}

	if success {
		jobItem.Status = "success"
	} else {
		jobItem.Status = "failed"
	}
	if updateErr := c.repo.UpdateJobItem(ctx, jobItem); updateErr != nil {
		c.log.Error(ctx).Err(updateErr).Msg("Failed to update bulk job item status")
	}

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

	job, err := c.repo.GetJobByID(ctx, jobItem.JobID)
	if err != nil || job == nil {
		return
	}

	_ = rabbitmq.PublishEvent(ctx, c.queue, rabbitmq.EventEnvelope{
		Type:       "bulk_job.progress",
		Category:   "invoice",
		BusinessID: job.BusinessID,
		Data: map[string]interface{}{
			"job_id":       job.ID,
			"total_orders": job.TotalOrders,
			"processed":    job.Processed,
			"successful":   job.Successful,
			"failed":       job.Failed,
			"progress":     job.GetProgress(),
			"status":       job.Status,
		},
	})

	if job.Successful+job.Failed >= job.TotalOrders {
		c.completeBulkJob(ctx, job)
	}
}

func (c *ResponseConsumer) completeBulkJob(ctx context.Context, job *entities.BulkInvoiceJob) {
	now := time.Now()
	job.Status = "completed"
	job.CompletedAt = &now

	if updateErr := c.repo.UpdateJob(ctx, job); updateErr != nil {
		c.log.Error(ctx).Err(updateErr).Str("job_id", job.ID).Msg("Failed to mark bulk job as completed")
		return
	}

	_ = rabbitmq.PublishEvent(ctx, c.queue, rabbitmq.EventEnvelope{
		Type:       "bulk_job.completed",
		Category:   "invoice",
		BusinessID: job.BusinessID,
		Data: map[string]interface{}{
			"job_id":       job.ID,
			"total_orders": job.TotalOrders,
			"processed":    job.Processed,
			"successful":   job.Successful,
			"failed":       job.Failed,
			"progress":     100,
			"status":       job.Status,
		},
	})

	c.log.Info(ctx).
		Str("job_id", job.ID).
		Int("successful", job.Successful).
		Int("failed", job.Failed).
		Int("total", job.TotalOrders).
		Msg("Bulk invoice job completed (from response consumer)")
}

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

	if response.CashReceiptRequestURL != "" {
		syncLog.CashReceiptRequestURL = response.CashReceiptRequestURL
	}
	if response.CashReceiptRequestPayload != nil {
		syncLog.CashReceiptRequestPayload = response.CashReceiptRequestPayload
	}
	if response.CashReceiptResponseStatus != 0 {
		syncLog.CashReceiptResponseStatus = response.CashReceiptResponseStatus
	}
	if response.CashReceiptResponseBody != "" {
		var bodyMap map[string]interface{}
		if json.Unmarshal([]byte(response.CashReceiptResponseBody), &bodyMap) == nil {
			syncLog.CashReceiptResponseBody = bodyMap
		} else {

			syncLog.CashReceiptResponseBody = map[string]interface{}{"raw": response.CashReceiptResponseBody}
		}
	}
}

func (c *ResponseConsumer) handleCompareResponse(ctx context.Context, message []byte) error {
	var msg compareResponseMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to unmarshal compare response")
		return fmt.Errorf("failed to unmarshal compare response: %w", err)
	}

	c.log.Info(ctx).
		Str("correlation_id", msg.CorrelationID).
		Uint("business_id", msg.BusinessID).
		Str("date_from", msg.DateFrom).
		Str("date_to", msg.DateTo).
		Int("provider_docs", len(msg.ProviderDocuments)).
		Msg("Processing compare response")

	if msg.Error != "" {
		c.log.Warn(ctx).Str("error", msg.Error).Msg("Provider returned error in compare response")
		data := &dtos.CompareResponseData{
			CorrelationID: msg.CorrelationID,
			BusinessID:    msg.BusinessID,
			DateFrom:      msg.DateFrom,
			DateTo:        msg.DateTo,
			Results: []dtos.CompareResult{
				{Status: dtos.CompareStatusProviderOnly, Comment: "Error del proveedor: " + msg.Error},
			},
			Summary: dtos.CompareSummary{},
		}

		if err := c.compareCache.StoreCompareResult(ctx, msg.CorrelationID, data); err != nil {
			c.log.Warn(ctx).Err(err).Str("correlation_id", msg.CorrelationID).Msg("Failed to store compare error result in Redis (non-fatal)")
		}
		go c.publishCompareEvent(ctx, data)
		return c.ssePublisher.PublishCompareReady(ctx, data)
	}

	systemInvoices, err := c.repo.GetIssuedInvoicesByDateRange(ctx, msg.BusinessID, msg.DateFrom, msg.DateTo)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to get system invoices for comparison")
		return fmt.Errorf("failed to get system invoices: %w", err)
	}

	systemMap := make(map[string]*entities.Invoice, len(systemInvoices))
	for _, inv := range systemInvoices {
		if inv.InvoiceNumber != "" {
			systemMap[inv.InvoiceNumber] = inv
		}
	}

	orderIDs := make([]string, 0, len(systemInvoices))
	for _, inv := range systemInvoices {
		if inv.OrderID != "" {
			orderIDs = append(orderIDs, inv.OrderID)
		}
	}
	orderDates, _ := c.repo.GetOrderCreatedAtsByIDs(ctx, orderIDs)

	providerMap := make(map[string]compareProviderDocument, len(msg.ProviderDocuments))
	for _, doc := range msg.ProviderDocuments {
		if doc.DocumentNumber != "" {
			providerMap[doc.DocumentNumber] = doc
		}
	}

	isSync := msg.Mode == "sync"
	results := make([]dtos.CompareResult, 0)
	matched, systemOnly, providerOnly, annulledInProvider, released := 0, 0, 0, 0, 0

	c.log.Info(ctx).
		Int("total_provider_docs", len(providerMap)).
		Msg("Starting comparison of provider documents")

	for docNum, doc := range providerMap {
		c.log.Debug(ctx).
			Str("doc_number", docNum).
			Str("document_name", doc.DocumentName).
			Bool("annuled", doc.Annuled).
			Msg("Processing document from provider")
		if sysInv, found := systemMap[docNum]; found {
			total := sysInv.TotalAmount
			orderID := sysInv.OrderID
			orderCreatedAt := formatOrderDate(orderDates, sysInv.OrderID)
			status := dtos.CompareStatusMatched
			wasReleased := false
			isCreditNote := strings.Contains(strings.ToUpper(doc.DocumentName), "NOTA") && strings.Contains(strings.ToUpper(doc.DocumentName), "CREDITO")
			if (doc.Annuled || isCreditNote) && sysInv.Status != constants.InvoiceStatusCancelled {
				status = dtos.CompareStatusAnnulledInProvider
				annulledInProvider++
				if isSync {
					ok, _, relErr := c.repo.CancelInvoiceAndReleaseOrder(ctx, sysInv.ID)
					if relErr != nil {
						c.log.Error(ctx).Err(relErr).Uint("invoice_id", sysInv.ID).Msg("Failed to release order for annulled invoice")
					} else if ok {
						wasReleased = true
						released++
					}
				}
			} else {
				matched++
				continue
			}
			results = append(results, dtos.CompareResult{
				Status:          status,
				InvoiceNumber:   docNum,
				Prefix:          doc.Prefix,
				DocumentDate:    doc.DocumentDate,
				ProviderTotal:   doc.Total,
				ProviderAnnuled: doc.Annuled,
				Released:        wasReleased,
				SystemInvoiceID: &sysInv.ID,
				SystemOrderID:   &orderID,
				SystemTotal:     &total,
				SystemStatus:    sysInv.Status,
				CustomerNit:     doc.CustomerNit,
				CustomerName:    doc.CustomerName,
				Comment:         doc.Comment,
				OrderCreatedAt:  orderCreatedAt,
				ProviderDetails: mapProviderDetailsToCompareDetails(doc.Details),
				SystemItems:     mapInvoiceItemsToCompareDetails(sysInv.Items),
			})
		} else {

			results = append(results, dtos.CompareResult{
				Status:          dtos.CompareStatusProviderOnly,
				InvoiceNumber:   docNum,
				Prefix:          doc.Prefix,
				DocumentDate:    doc.DocumentDate,
				ProviderTotal:   doc.Total,
				CustomerNit:     doc.CustomerNit,
				CustomerName:    doc.CustomerName,
				Comment:         doc.Comment,
				ProviderDetails: mapProviderDetailsToCompareDetails(doc.Details),
			})
			providerOnly++
		}
	}

	for invNum, sysInv := range systemMap {
		if _, found := providerMap[invNum]; !found {
			total := sysInv.TotalAmount
			orderID := sysInv.OrderID

			customerNit := sysInv.CustomerDNI
			orderCreatedAt := formatOrderDate(orderDates, sysInv.OrderID)
			results = append(results, dtos.CompareResult{
				Status:          dtos.CompareStatusSystemOnly,
				InvoiceNumber:   invNum,
				SystemInvoiceID: &sysInv.ID,
				SystemOrderID:   &orderID,
				SystemTotal:     &total,
				CustomerNit:     customerNit,
				OrderCreatedAt:  orderCreatedAt,
				SystemItems:     mapInvoiceItemsToCompareDetails(sysInv.Items),
			})
			systemOnly++
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].DocumentDate != results[j].DocumentDate {
			return results[i].DocumentDate > results[j].DocumentDate
		}
		return results[i].InvoiceNumber > results[j].InvoiceNumber
	})

	summary := dtos.CompareSummary{
		Matched:            matched,
		SystemOnly:         systemOnly,
		ProviderOnly:       providerOnly,
		AnnulledInProvider: annulledInProvider,
		Released:           released,
	}

	if isSync {
		filtered := make([]dtos.CompareResult, 0, annulledInProvider)
		for _, r := range results {
			if r.Status == dtos.CompareStatusAnnulledInProvider {
				filtered = append(filtered, r)
			}
		}
		results = filtered
	}

	responseData := &dtos.CompareResponseData{
		CorrelationID: msg.CorrelationID,
		BusinessID:    msg.BusinessID,
		DateFrom:      msg.DateFrom,
		DateTo:        msg.DateTo,
		Results:       results,
		Summary:       summary,
	}

	c.log.Info(ctx).
		Str("correlation_id", msg.CorrelationID).
		Int("matched", matched).
		Int("system_only", systemOnly).
		Int("provider_only", providerOnly).
		Msg("Comparison complete, storing in Redis + publishing SSE")

	if err := c.compareCache.StoreCompareResult(ctx, msg.CorrelationID, responseData); err != nil {
		c.log.Warn(ctx).Err(err).Str("correlation_id", msg.CorrelationID).Msg("Failed to store compare result in Redis (non-fatal)")
	}

	go c.publishCompareEvent(ctx, responseData)
	return c.ssePublisher.PublishCompareReady(ctx, responseData)
}

func (c *ResponseConsumer) publishCompareEvent(ctx context.Context, data *dtos.CompareResponseData) {
	results := make([]map[string]interface{}, 0, len(data.Results))
	for _, r := range data.Results {
		row := map[string]interface{}{
			"status":           r.Status,
			"invoice_number":   r.InvoiceNumber,
			"prefix":           r.Prefix,
			"document_date":    r.DocumentDate,
			"provider_total":   r.ProviderTotal,
			"customer_nit":     r.CustomerNit,
			"customer_name":    r.CustomerName,
			"comment":          r.Comment,
			"order_created_at": r.OrderCreatedAt,
		}
		if r.SystemInvoiceID != nil {
			row["system_invoice_id"] = *r.SystemInvoiceID
		}
		if r.SystemOrderID != nil {
			row["system_order_id"] = *r.SystemOrderID
		}
		if r.SystemTotal != nil {
			row["system_total"] = *r.SystemTotal
		}
		if len(r.ProviderDetails) > 0 {
			details := make([]map[string]interface{}, 0, len(r.ProviderDetails))
			for _, d := range r.ProviderDetails {
				details = append(details, map[string]interface{}{
					"item_code":  d.ItemCode,
					"item_name":  d.ItemName,
					"quantity":   d.Quantity,
					"unit_value": d.UnitValue,
					"iva":        d.IVA,
				})
			}
			row["provider_details"] = details
		}
		if len(r.SystemItems) > 0 {
			items := make([]map[string]interface{}, 0, len(r.SystemItems))
			for _, d := range r.SystemItems {
				items = append(items, map[string]interface{}{
					"item_code":  d.ItemCode,
					"item_name":  d.ItemName,
					"quantity":   d.Quantity,
					"unit_value": d.UnitValue,
					"iva":        d.IVA,
				})
			}
			row["system_items"] = items
		}
		results = append(results, row)
	}

	_ = rabbitmq.PublishEvent(ctx, c.queue, rabbitmq.EventEnvelope{
		Type:       "invoice.compare_ready",
		Category:   "invoice",
		BusinessID: data.BusinessID,
		Data: map[string]interface{}{
			"correlation_id": data.CorrelationID,
			"date_from":      data.DateFrom,
			"date_to":        data.DateTo,
			"results":        results,
			"summary": map[string]interface{}{
				"matched":       data.Summary.Matched,
				"system_only":   data.Summary.SystemOnly,
				"provider_only": data.Summary.ProviderOnly,
			},
		},
	})
}

func mapInvoiceItemsToCompareDetails(items []entities.InvoiceItem) []dtos.CompareItemDetail {
	result := make([]dtos.CompareItemDetail, 0, len(items))
	for _, it := range items {
		iva := "0"
		if it.TaxRate != nil {
			iva = fmt.Sprintf("%.0f", *it.TaxRate*100)
		}
		result = append(result, dtos.CompareItemDetail{
			ItemCode:  it.SKU,
			ItemName:  it.Name,
			Quantity:  fmt.Sprintf("%d", it.Quantity),
			UnitValue: fmt.Sprintf("%.2f", it.UnitPrice),
			IVA:       iva,
		})
	}
	return result
}

func mapProviderDetailsToCompareDetails(details []compareItemDetail) []dtos.CompareItemDetail {
	result := make([]dtos.CompareItemDetail, 0, len(details))
	for _, d := range details {
		result = append(result, dtos.CompareItemDetail{
			ItemCode:  d.ItemCode,
			ItemName:  d.ItemName,
			Quantity:  d.Quantity,
			UnitValue: d.Value,
			IVA:       d.IVA,
		})
	}
	return result
}

func formatOrderDate(orderDates map[string]*time.Time, orderID string) *string {
	if t, ok := orderDates[orderID]; ok && t != nil {
		s := t.Format("2006-01-02")
		return &s
	}
	return nil
}

const providerUnavailableMaxWindow = 6 * time.Hour

func providerUnavailableRetryDelay(elapsed time.Duration) time.Duration {
	switch {
	case elapsed < time.Hour:
		return 15 * time.Minute
	case elapsed < 2*time.Hour:
		return 30 * time.Minute
	case elapsed < 4*time.Hour:
		return time.Hour
	default:
		return 2 * time.Hour
	}
}

var providerUnavailableErrorMarkers = []string{
	"context deadline exceeded",
	"client.timeout",
	"connection refused",
	"connection reset",
	"no such host",
	"i/o timeout",
	"tls handshake timeout",
	"authentication request failed",
	"reintento abortado",
	"no fue posible verificar",
	"response has no info",
	"internal server error",
	"status 429",
	"status 502",
	"status 503",
	"status 504",
}

var codigosDeConfiguracion = map[string]bool{
	"document_inactive":               true,
	"document_not_electronic":         true,
	"document_settings":               true,
	"product_inactive":                true,
	"customer_inactive":               true,
	"payment_method_inactive":         true,
	"seller_inactive":                 true,
	"parameter_inactive":              true,
	"invalid_reference":               true,
	"invalid_product_code":            true,
	"invalid_credentials":             true,
	"incomplete_siigo_config":         true,
	"missing_customer_identification": true,
}

func esErrorDeConfiguracion(codigo string) bool {
	return codigosDeConfiguracion[codigo]
}

func isProviderUnavailableError(errMsg string) bool {
	lower := strings.ToLower(errMsg)
	for _, marker := range providerUnavailableErrorMarkers {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

func (c *ResponseConsumer) calculateNextRetry(retryCount int) time.Time {

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

func (c *ResponseConsumer) calculateCheckBackoff(checkCount int) time.Time {
	delays := []time.Duration{
		15 * time.Minute,
		30 * time.Minute,
		1 * time.Hour,
		2 * time.Hour,
		4 * time.Hour,
	}

	delayIndex := checkCount
	if delayIndex >= len(delays) {
		delayIndex = len(delays) - 1
	}

	return time.Now().Add(delays[delayIndex])
}

type listItemsProviderItem struct {
	ItemCode      string  `json:"item_code"`
	ItemName      string  `json:"item_name"`
	ItemPrice     float64 `json:"item_price"`
	UnitCost      float64 `json:"unit_cost"`
	Description   string  `json:"description"`
	MinimumStock  string  `json:"minimum_stock"`
	OrderQuantity string  `json:"order_quantity"`
}

type listItemsResponseMessage struct {
	Operation     string                  `json:"operation"`
	CorrelationID string                  `json:"correlation_id"`
	BusinessID    uint                    `json:"business_id"`
	Items         []listItemsProviderItem `json:"items"`
	Error         string                  `json:"error,omitempty"`
	Timestamp     time.Time               `json:"timestamp"`
}

func (c *ResponseConsumer) handleListItemsResponse(ctx context.Context, message []byte) error {
	var msg listItemsResponseMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to unmarshal list_items response")
		return fmt.Errorf("failed to unmarshal list_items response: %w", err)
	}

	c.log.Info(ctx).
		Str("correlation_id", msg.CorrelationID).
		Uint("business_id", msg.BusinessID).
		Int("provider_items", len(msg.Items)).
		Msg("Processing list_items response")

	if msg.Error != "" {
		c.log.Warn(ctx).Str("error", msg.Error).Msg("Provider returned error in list_items response")
		data := &dtos.ItemCompareResponseData{
			CorrelationID: msg.CorrelationID,
			BusinessID:    msg.BusinessID,
			Results:       []dtos.ItemCompareResult{},
			Summary:       dtos.ItemCompareSummary{},
		}

		if err := c.compareCache.StoreItemCompareResult(ctx, msg.CorrelationID, data); err != nil {
			c.log.Warn(ctx).Err(err).Str("correlation_id", msg.CorrelationID).Msg("Failed to store item compare error result in Redis (non-fatal)")
		}
		go c.publishListItemsEvent(ctx, data)
		return c.ssePublisher.PublishListItemsReady(ctx, data)
	}

	systemProducts, err := c.repo.ListProductsByBusinessID(ctx, msg.BusinessID)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to get system products for items comparison")
		return fmt.Errorf("failed to get system products: %w", err)
	}

	providerMap := make(map[string]listItemsProviderItem, len(msg.Items))
	for _, item := range msg.Items {
		if item.ItemCode != "" {
			providerMap[item.ItemCode] = item
		}
	}

	systemMap := make(map[string]dtos.SystemProduct, len(systemProducts))
	for _, prod := range systemProducts {
		if prod.SKU != "" {
			systemMap[prod.SKU] = prod
		}
	}

	results := make([]dtos.ItemCompareResult, 0)
	matched, providerOnly, systemOnly := 0, 0, 0

	for code, pItem := range providerMap {
		if sProd, found := systemMap[code]; found {

			results = append(results, dtos.ItemCompareResult{
				Status:        dtos.CompareStatusMatched,
				ItemCode:      code,
				ProviderName:  pItem.ItemName,
				SystemName:    sProd.Name,
				ProviderPrice: pItem.ItemPrice,
				SystemPrice:   sProd.Price,
				PriceDiff:     pItem.ItemPrice - sProd.Price,
				UnitCost:      pItem.UnitCost,
				Description:   pItem.Description,
			})
			matched++
		} else {

			results = append(results, dtos.ItemCompareResult{
				Status:        dtos.CompareStatusProviderOnly,
				ItemCode:      code,
				ProviderName:  pItem.ItemName,
				ProviderPrice: pItem.ItemPrice,
				UnitCost:      pItem.UnitCost,
				Description:   pItem.Description,
			})
			providerOnly++
		}
	}

	for sku, sProd := range systemMap {
		if _, found := providerMap[sku]; !found {
			results = append(results, dtos.ItemCompareResult{
				Status:      dtos.CompareStatusSystemOnly,
				ItemCode:    sku,
				SystemName:  sProd.Name,
				SystemPrice: sProd.Price,
			})
			systemOnly++
		}
	}

	sort.Slice(results, func(i, j int) bool {
		statusOrder := map[string]int{
			dtos.CompareStatusMatched:      0,
			dtos.CompareStatusProviderOnly: 1,
			dtos.CompareStatusSystemOnly:   2,
		}
		if statusOrder[results[i].Status] != statusOrder[results[j].Status] {
			return statusOrder[results[i].Status] < statusOrder[results[j].Status]
		}
		return results[i].ItemCode < results[j].ItemCode
	})

	summary := dtos.ItemCompareSummary{
		Matched:       matched,
		ProviderOnly:  providerOnly,
		SystemOnly:    systemOnly,
		TotalProvider: len(msg.Items),
		TotalSystem:   len(systemProducts),
	}

	responseData := &dtos.ItemCompareResponseData{
		CorrelationID: msg.CorrelationID,
		BusinessID:    msg.BusinessID,
		Results:       results,
		Summary:       summary,
	}

	c.log.Info(ctx).
		Str("correlation_id", msg.CorrelationID).
		Int("matched", matched).
		Int("provider_only", providerOnly).
		Int("system_only", systemOnly).
		Msg("Items comparison complete, storing in Redis + publishing SSE")

	if err := c.compareCache.StoreItemCompareResult(ctx, msg.CorrelationID, responseData); err != nil {
		c.log.Warn(ctx).Err(err).Str("correlation_id", msg.CorrelationID).Msg("Failed to store item compare result in Redis (non-fatal)")
	}

	go c.publishListItemsEvent(ctx, responseData)
	return c.ssePublisher.PublishListItemsReady(ctx, responseData)
}

type listBankAccountsResponseMessage struct {
	Operation     string `json:"operation"`
	CorrelationID string `json:"correlation_id"`
	BusinessID    uint   `json:"business_id"`
	Items         []struct {
		AccountNumber string `json:"account_number"`
		Name          string `json:"name"`
		NameType      string `json:"name_type"`
	} `json:"items"`
	Error     string    `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func (c *ResponseConsumer) handleListBankAccountsResponse(ctx context.Context, message []byte) error {
	var msg listBankAccountsResponseMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to unmarshal list_bank_accounts response")
		return fmt.Errorf("failed to unmarshal list_bank_accounts response: %w", err)
	}

	c.log.Info(ctx).
		Str("correlation_id", msg.CorrelationID).
		Uint("business_id", msg.BusinessID).
		Int("accounts", len(msg.Items)).
		Msg("Processing list_bank_accounts response")

	if msg.Error != "" {
		c.log.Warn(ctx).Str("error", msg.Error).Msg("Provider returned error in list_bank_accounts response")
		data := &dtos.BankAccountsResponseData{
			CorrelationID: msg.CorrelationID,
			BusinessID:    msg.BusinessID,
			Results:       []dtos.BankAccountResult{},
		}
		if err := c.compareCache.StoreBankAccountsResult(ctx, msg.CorrelationID, data); err != nil {
			c.log.Warn(ctx).Err(err).Str("correlation_id", msg.CorrelationID).Msg("Failed to store bank accounts error result in Redis (non-fatal)")
		}
		_ = rabbitmq.PublishEvent(ctx, c.queue, rabbitmq.EventEnvelope{
			Type:       "invoice.list_bank_accounts_ready",
			Category:   "invoice",
			BusinessID: msg.BusinessID,
			Data: map[string]interface{}{
				"correlation_id": msg.CorrelationID,
				"results":        []map[string]interface{}{},
				"error":          msg.Error,
			},
		})
		return nil
	}

	results := make([]dtos.BankAccountResult, 0, len(msg.Items))
	for _, item := range msg.Items {
		results = append(results, dtos.BankAccountResult{
			AccountNumber: item.AccountNumber,
			Name:          item.Name,
			NameType:      item.NameType,
		})
	}

	responseData := &dtos.BankAccountsResponseData{
		CorrelationID: msg.CorrelationID,
		BusinessID:    msg.BusinessID,
		Results:       results,
	}

	if err := c.compareCache.StoreBankAccountsResult(ctx, msg.CorrelationID, responseData); err != nil {
		c.log.Warn(ctx).Err(err).Str("correlation_id", msg.CorrelationID).Msg("Failed to store bank accounts result in Redis (non-fatal)")
	}

	resultsData := make([]map[string]interface{}, 0, len(results))
	for _, r := range results {
		resultsData = append(resultsData, map[string]interface{}{
			"account_number": r.AccountNumber,
			"name":           r.Name,
			"name_type":      r.NameType,
		})
	}

	_ = rabbitmq.PublishEvent(ctx, c.queue, rabbitmq.EventEnvelope{
		Type:       "invoice.list_bank_accounts_ready",
		Category:   "invoice",
		BusinessID: msg.BusinessID,
		Data: map[string]interface{}{
			"correlation_id": msg.CorrelationID,
			"results":        resultsData,
		},
	})

	c.log.Info(ctx).
		Str("correlation_id", msg.CorrelationID).
		Int("accounts", len(results)).
		Msg("Bank accounts response processed successfully")

	return nil
}

type listSiigoWarehousesResponseMessage struct {
	Operation     string `json:"operation"`
	CorrelationID string `json:"correlation_id"`
	BusinessID    uint   `json:"business_id"`
	Items         []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"items"`
	Error     string    `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func (c *ResponseConsumer) handleListSiigoWarehousesResponse(ctx context.Context, message []byte) error {
	var msg listSiigoWarehousesResponseMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to unmarshal list_siigo_warehouses response")
		return fmt.Errorf("failed to unmarshal list_siigo_warehouses response: %w", err)
	}

	results := make([]map[string]interface{}, 0, len(msg.Items))
	for _, item := range msg.Items {
		results = append(results, map[string]interface{}{
			"id":   item.ID,
			"name": item.Name,
		})
	}

	data := map[string]interface{}{
		"correlation_id": msg.CorrelationID,
		"results":        results,
	}
	if msg.Error != "" {
		data["error"] = msg.Error
	}

	_ = rabbitmq.PublishEvent(ctx, c.queue, rabbitmq.EventEnvelope{
		Type:       "invoice.siigo_warehouses_ready",
		Category:   "invoice",
		BusinessID: msg.BusinessID,
		Data:       data,
	})

	c.log.Info(ctx).
		Str("correlation_id", msg.CorrelationID).
		Int("warehouses", len(results)).
		Msg("Siigo warehouses response processed")

	return nil
}

func (c *ResponseConsumer) publishListItemsEvent(ctx context.Context, data *dtos.ItemCompareResponseData) {
	results := make([]map[string]interface{}, 0, len(data.Results))
	for _, r := range data.Results {
		results = append(results, map[string]interface{}{
			"status":         r.Status,
			"item_code":      r.ItemCode,
			"provider_name":  r.ProviderName,
			"system_name":    r.SystemName,
			"provider_price": r.ProviderPrice,
			"system_price":   r.SystemPrice,
			"price_diff":     r.PriceDiff,
			"unit_cost":      r.UnitCost,
			"description":    r.Description,
		})
	}

	_ = rabbitmq.PublishEvent(ctx, c.queue, rabbitmq.EventEnvelope{
		Type:       "invoice.list_items_ready",
		Category:   "invoice",
		BusinessID: data.BusinessID,
		Data: map[string]interface{}{
			"correlation_id": data.CorrelationID,
			"results":        results,
			"summary": map[string]interface{}{
				"matched":        data.Summary.Matched,
				"provider_only":  data.Summary.ProviderOnly,
				"system_only":    data.Summary.SystemOnly,
				"total_provider": data.Summary.TotalProvider,
				"total_system":   data.Summary.TotalSystem,
			},
		},
	})
}
