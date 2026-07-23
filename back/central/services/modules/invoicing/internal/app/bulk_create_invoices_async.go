package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
)

func (uc *useCase) BulkCreateInvoicesAsync(ctx context.Context, dto *dtos.BulkCreateInvoicesDTO) (string, error) {
	if len(dto.OrderIDs) == 0 {
		return "", fmt.Errorf("order_ids cannot be empty")
	}

	if len(dto.OrderIDs) > dtos.MaxBulkInvoiceOrders {
		return "", fmt.Errorf("maximum %d orders per batch (received %d)", dtos.MaxBulkInvoiceOrders, len(dto.OrderIDs))
	}

	var businessID uint
	if dto.BusinessID != nil && *dto.BusinessID > 0 {
		businessID = *dto.BusinessID
	} else {
		businessID, _ = ctx.Value("business_id").(uint)
	}

	if businessID == 0 {
		uc.log.Error(ctx).Msg("business_id is required")
		return "", fmt.Errorf("business_id is required: super admin must select a business")
	}

	jobID := uuid.New().String()

	job := &entities.BulkInvoiceJob{
		ID:          jobID,
		BusinessID:  businessID,
		TotalOrders: len(dto.OrderIDs),
		Processed:   0,
		Successful:  0,
		Failed:      0,
		Status:      entities.JobStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if userID, ok := ctx.Value("user_id").(uint); ok && userID > 0 {
		job.CreatedByUserID = &userID
	}

	if err := uc.repo.CreateJob(ctx, job); err != nil {
		uc.log.Error(ctx).
			Err(err).
			Str("job_id", jobID).
			Msg("Failed to create bulk invoice job")
		return "", fmt.Errorf("failed to create bulk invoice job: %w", err)
	}

	uc.log.Info(ctx).
		Str("job_id", jobID).
		Int("total_orders", len(dto.OrderIDs)).
		Uint("business_id", businessID).
		Msg("Bulk invoice job created")

	items := make([]*entities.BulkInvoiceJobItem, len(dto.OrderIDs))
	for i, orderID := range dto.OrderIDs {
		items[i] = &entities.BulkInvoiceJobItem{
			JobID:     jobID,
			OrderID:   orderID,
			Status:    entities.JobItemStatusPending,
			CreatedAt: time.Now(),
		}
	}

	if err := uc.repo.CreateJobItems(ctx, items); err != nil {
		uc.log.Error(ctx).
			Err(err).
			Str("job_id", jobID).
			Msg("Failed to create bulk job items")
		return "", fmt.Errorf("failed to create bulk job items: %w", err)
	}

	publishedCount := 0
	for _, orderID := range dto.OrderIDs {
		message := &dtos.BulkInvoiceJobMessage{
			JobID:         jobID,
			OrderID:       orderID,
			BusinessID:    businessID,
			IsManual:      true,
			AttemptNumber: 1,
		}

		if job.CreatedByUserID != nil {
			message.CreatedBy = job.CreatedByUserID
		}

		if err := uc.eventPublisher.PublishBulkInvoiceJob(ctx, message); err != nil {
			uc.log.Warn(ctx).
				Err(err).
				Str("job_id", jobID).
				Str("order_id", orderID).
				Msg("Failed to publish bulk invoice job message (will continue with others)")

			errMsg := fmt.Sprintf("failed to publish to queue: %s", err.Error())
			now := time.Now()
			failedItem := &entities.BulkInvoiceJobItem{
				JobID:        jobID,
				OrderID:      orderID,
				Status:       entities.JobItemStatusFailed,
				ErrorMessage: &errMsg,
				ProcessedAt:  &now,
			}
			for _, item := range items {
				if item.OrderID == orderID {
					failedItem.ID = item.ID
					break
				}
			}
			if failedItem.ID > 0 {
				if updateErr := uc.repo.UpdateJobItem(ctx, failedItem); updateErr != nil {
					uc.log.Error(ctx).Err(updateErr).Str("order_id", orderID).Msg("Failed to mark job item as failed")
				}
			}
			job.Failed++
			continue
		}

		publishedCount++
	}

	uc.log.Info(ctx).
		Str("job_id", jobID).
		Int("published", publishedCount).
		Int("total", len(dto.OrderIDs)).
		Msg("Bulk invoice job messages published to RabbitMQ")

	if publishedCount > 0 {
		now := time.Now()
		job.Status = entities.JobStatusProcessing
		job.StartedAt = &now
		uc.repo.UpdateJob(ctx, job)
	}

	return jobID, nil
}

func (uc *useCase) GetBulkJobStatus(ctx context.Context, jobID string) (*entities.BulkInvoiceJob, error) {
	job, err := uc.repo.GetJobByID(ctx, jobID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Str("job_id", jobID).Msg("Failed to get bulk job status")
		return nil, err
	}

	if job == nil {
		return nil, fmt.Errorf("job not found: %s", jobID)
	}

	return job, nil
}

func (uc *useCase) GetBulkJobItems(ctx context.Context, jobID string) ([]*entities.BulkInvoiceJobItem, error) {
	items, err := uc.repo.GetJobItems(ctx, jobID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Str("job_id", jobID).Msg("Failed to get bulk job items")
		return nil, err
	}

	return items, nil
}

func (uc *useCase) ListBulkJobs(ctx context.Context, businessID uint, page, pageSize int) ([]*entities.BulkInvoiceJob, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	jobs, total, err := uc.repo.ListJobs(ctx, businessID, page, pageSize)
	if err != nil {
		uc.log.Error(ctx).
			Err(err).
			Uint("business_id", businessID).
			Msg("Failed to list bulk jobs")
		return nil, 0, err
	}

	return jobs, total, nil
}
