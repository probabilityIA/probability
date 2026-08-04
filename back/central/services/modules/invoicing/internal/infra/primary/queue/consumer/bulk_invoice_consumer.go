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

type BulkInvoiceConsumer struct {
	queue        rabbitmq.IQueue
	useCase      ports.IUseCase
	repo         ports.IRepository
	ssePublisher ports.IInvoiceSSEPublisher
	log          log.ILogger
}

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

func (c *BulkInvoiceConsumer) Start(ctx context.Context) error {
	c.queue.DeclareQueue(QueueBulkInvoiceJobs, true)

	return c.queue.Consume(ctx, QueueBulkInvoiceJobs, c.handleMessage)
}

func (c *BulkInvoiceConsumer) handleMessage(message []byte) error {
	ctx := context.Background()

	var queueMsg messages.BulkInvoiceJobMessage
	if err := json.Unmarshal(message, &queueMsg); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to unmarshal bulk invoice job message")
		return nil
	}

	msg := queueMappers.BulkJobMessageToDTO(&queueMsg)

	c.log.Info(ctx).
		Str("job_id", msg.JobID).
		Str("order_id", msg.OrderID).
		Int("attempt", msg.AttemptNumber).
		Msg("Processing bulk invoice message")

	if err := c.updateItemStatus(ctx, msg.JobID, msg.OrderID, entities.JobItemStatusProcessing, nil); err != nil {
		c.log.Warn(ctx).Err(err).Str("job_id", msg.JobID).Str("order_id", msg.OrderID).Msg("Failed to update item to processing")
	}

	ctx = context.WithValue(ctx, "business_id", msg.BusinessID)
	if msg.CreatedBy != nil {
		ctx = context.WithValue(ctx, "user_id", *msg.CreatedBy)
	}

	invoice, err := c.useCase.CreateInvoice(ctx, &dtos.CreateInvoiceDTO{
		OrderID:  msg.OrderID,
		IsManual: msg.IsManual,
	})

	if err != nil {
		c.handleInvoiceError(ctx, msg.JobID, msg.OrderID, msg.BusinessID, err)
		return nil
	}

	c.handleInvoiceSuccess(ctx, msg.JobID, msg.OrderID, invoice.ID)

	return nil
}

func (c *BulkInvoiceConsumer) handleInvoiceError(ctx context.Context, jobID, orderID string, businessID uint, err error) {
	errMsg := err.Error()

	invoice, getErr := c.repo.GetInvoiceByOrderID(ctx, orderID)

	invoiceID := uint(0)
	if getErr == nil && invoice != nil {
		invoiceID = invoice.ID
	} else {
		c.log.Warn(ctx).
			Err(getErr).
			Str("order_id", orderID).
			Msg("Invoice not found in DB after failure - publishing error event without invoice_id")
	}

	_ = rabbitmq.PublishEvent(ctx, c.queue, rabbitmq.EventEnvelope{
		Type:       "invoice.failed",
		Category:   "invoice",
		BusinessID: businessID,
		Data: map[string]interface{}{
			"invoice_id":    invoiceID,
			"order_id":      orderID,
			"error_message": errMsg,
		},
	})

	if updateErr := c.updateItemStatus(ctx, jobID, orderID, entities.JobItemStatusFailed, &errMsg); updateErr != nil {
		c.log.Error(ctx).Err(updateErr).Str("job_id", jobID).Str("order_id", orderID).Msg("Failed to update item to failed")
	}

	if incrementErr := c.repo.IncrementJobCounters(ctx, jobID, 1, 0, 1); incrementErr != nil {
		c.log.Error(ctx).Err(incrementErr).Str("job_id", jobID).Msg("Failed to increment job counters for failed invoice")
	}

	c.publishJobProgress(ctx, jobID)

	c.checkJobCompletion(ctx, jobID)

	c.log.Warn(ctx).
		Err(err).
		Str("job_id", jobID).
		Str("order_id", orderID).
		Msg("Invoice creation failed in bulk job")
}

func (c *BulkInvoiceConsumer) handleInvoiceSuccess(ctx context.Context, jobID, orderID string, invoiceID uint) {

	now := time.Now()
	item, err := c.repo.GetJobItemByJobAndOrder(ctx, jobID, orderID)
	if err != nil {
		c.log.Error(ctx).Err(err).Str("job_id", jobID).Str("order_id", orderID).Msg("Failed to get job item")
		return
	}

	if item != nil {
		item.Status = entities.JobItemStatusProcessing
		item.InvoiceID = &invoiceID
		item.ProcessedAt = &now
		item.ErrorMessage = nil

		if updateErr := c.repo.UpdateJobItem(ctx, item); updateErr != nil {
			c.log.Error(ctx).Err(updateErr).Str("job_id", jobID).Str("order_id", orderID).Msg("Failed to update item to processing")
		}
	}

	if incrementErr := c.repo.IncrementJobCounters(ctx, jobID, 1, 0, 0); incrementErr != nil {
		c.log.Error(ctx).Err(incrementErr).Str("job_id", jobID).Msg("Failed to increment job counters")
	}

	c.publishJobProgress(ctx, jobID)

	c.log.Info(ctx).
		Str("job_id", jobID).
		Str("order_id", orderID).
		Uint("invoice_id", invoiceID).
		Msg("Invoice created successfully in bulk job")
}

func (c *BulkInvoiceConsumer) updateItemStatus(ctx context.Context, jobID, orderID, status string, errorMsg *string) error {
	item, err := c.repo.GetJobItemByJobAndOrder(ctx, jobID, orderID)
	if err != nil {
		return err
	}

	if item == nil {
		c.log.Warn(ctx).
			Str("job_id", jobID).
			Str("order_id", orderID).
			Msg("Job item not found for update")
		return nil
	}

	item.Status = status
	item.ErrorMessage = errorMsg

	if status == entities.JobItemStatusProcessing {
		now := time.Now()
		item.ProcessedAt = &now
	}

	return c.repo.UpdateJobItem(ctx, item)
}

func (c *BulkInvoiceConsumer) checkJobCompletion(ctx context.Context, jobID string) {
	job, err := c.repo.GetJobByID(ctx, jobID)
	if err != nil {
		c.log.Error(ctx).Err(err).Str("job_id", jobID).Msg("Failed to get job for completion check")
		return
	}

	if job == nil {
		return
	}

	if job.Processed >= job.TotalOrders {
		now := time.Now()
		job.Status = entities.JobStatusCompleted
		job.CompletedAt = &now

		if updateErr := c.repo.UpdateJob(ctx, job); updateErr != nil {
			c.log.Error(ctx).Err(updateErr).Str("job_id", jobID).Msg("Failed to mark job as completed")
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
			Str("job_id", jobID).
			Int("successful", job.Successful).
			Int("failed", job.Failed).
			Int("total", job.TotalOrders).
			Msg("Bulk invoice job completed")
	}
}

func (c *BulkInvoiceConsumer) publishJobProgress(ctx context.Context, jobID string) {
	job, err := c.repo.GetJobByID(ctx, jobID)
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
}
