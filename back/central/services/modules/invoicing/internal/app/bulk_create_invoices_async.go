package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
)

// BulkCreateInvoicesAsync crea un job de facturación masiva y publica mensajes a RabbitMQ
func (uc *useCase) BulkCreateInvoicesAsync(ctx context.Context, dto *dtos.BulkCreateInvoicesDTO) (string, error) {
	// 1. Validar DTO
	if len(dto.OrderIDs) == 0 {
		return "", fmt.Errorf("order_ids cannot be empty")
	}

	if len(dto.OrderIDs) > 500 {
		return "", fmt.Errorf("maximum 500 orders per batch (received %d)", len(dto.OrderIDs))
	}

	// 2. Determinar business_id:
	//    - Super admin: viene en el DTO (seleccionado en frontend)
	//    - Usuario normal: viene en el contexto JWT
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

	// 3. Generar jobID
	jobID := uuid.New().String()

	// 4. Crear job en BD
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

	// Agregar user_id si está en contexto
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

	// 5. Crear job items (una fila por orden)
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

	// 6. Publicar mensajes a RabbitMQ (uno por orden)
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

		// Publicar mensaje - no fallar todo el job si uno falla
		if err := uc.eventPublisher.PublishBulkInvoiceJob(ctx, message); err != nil {
			uc.log.Warn(ctx).
				Err(err).
				Str("job_id", jobID).
				Str("order_id", orderID).
				Msg("Failed to publish bulk invoice job message (will continue with others)")

			// Marcar el job item como failed inmediatamente para que el usuario sepa cuáles no se publicaron
			errMsg := fmt.Sprintf("failed to publish to queue: %s", err.Error())
			now := time.Now()
			failedItem := &entities.BulkInvoiceJobItem{
				JobID:        jobID,
				OrderID:      orderID,
				Status:       entities.JobItemStatusFailed,
				ErrorMessage: &errMsg,
				ProcessedAt:  &now,
			}
			// Buscar el item existente y actualizarlo
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

	// 7. Actualizar job status a processing si al menos un mensaje se publicó
	if publishedCount > 0 {
		now := time.Now()
		job.Status = entities.JobStatusProcessing
		job.StartedAt = &now
		uc.repo.UpdateJob(ctx, job)
	}

	return jobID, nil
}

// GetBulkJobStatus obtiene el estado de un job de facturación masiva
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

// GetBulkJobItems obtiene los items de un job de facturación masiva
func (uc *useCase) GetBulkJobItems(ctx context.Context, jobID string) ([]*entities.BulkInvoiceJobItem, error) {
	items, err := uc.repo.GetJobItems(ctx, jobID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Str("job_id", jobID).Msg("Failed to get bulk job items")
		return nil, err
	}

	return items, nil
}

// ListBulkJobs lista los jobs de facturación masiva de un negocio
func (uc *useCase) ListBulkJobs(ctx context.Context, businessID uint, page, pageSize int) ([]*entities.BulkInvoiceJob, int64, error) {
	// Validar paginación
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
