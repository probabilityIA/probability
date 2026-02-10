package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// CreateBulkInvoiceJob crea un nuevo job de facturación masiva en la base de datos
func (r *Repository) CreateJob(ctx context.Context, job *entities.BulkInvoiceJob) error {
	model := mappers.JobToModel(job)

	if err := r.db.Conn(ctx).Create(&model).Error; err != nil {
		r.log.Error(ctx).Err(err).Str("job_id", job.ID).Msg("Failed to create bulk invoice job")
		return err
	}

	// Actualizar el ID generado
	job.ID = model.ID.String()

	r.log.Info(ctx).
		Str("job_id", job.ID).
		Int("total_orders", job.TotalOrders).
		Msg("Bulk invoice job created")

	return nil
}

// CreateBulkInvoiceJobItems crea múltiples items de job en la base de datos (batch insert)
func (r *Repository) CreateJobItems(ctx context.Context, items []*entities.BulkInvoiceJobItem) error {
	if len(items) == 0 {
		return nil
	}

	itemModels := make([]models.BulkInvoiceJobItem, len(items))
	for i, item := range items {
		itemModels[i] = mappers.JobItemToModel(item)
	}

	if err := r.db.Conn(ctx).Create(&itemModels).Error; err != nil {
		r.log.Error(ctx).Err(err).Int("items_count", len(items)).Msg("Failed to create bulk job items")
		return err
	}

	r.log.Info(ctx).Int("items_count", len(items)).Msg("Bulk job items created")

	return nil
}

// GetBulkInvoiceJobByID busca un job de facturación masiva por su UUID en la base de datos
func (r *Repository) GetJobByID(ctx context.Context, jobID string) (*entities.BulkInvoiceJob, error) {
	jobUUID, err := uuid.Parse(jobID)
	if err != nil {
		return nil, err
	}

	var jobModel models.BulkInvoiceJob
	if err := r.db.Conn(ctx).Where("id = ?", jobUUID).First(&jobModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.log.Error(ctx).Err(err).Str("job_id", jobID).Msg("Failed to get bulk job")
		return nil, err
	}

	return mappers.JobToDomain(&jobModel), nil
}

// GetBulkInvoiceJobItemsByJobID busca todos los items de un job específico en la base de datos
func (r *Repository) GetJobItems(ctx context.Context, jobID string) ([]*entities.BulkInvoiceJobItem, error) {
	jobUUID, err := uuid.Parse(jobID)
	if err != nil {
		return nil, err
	}

	var itemModels []models.BulkInvoiceJobItem
	if err := r.db.Conn(ctx).
		Where("job_id = ?", jobUUID).
		Order("created_at ASC").
		Find(&itemModels).Error; err != nil {
		r.log.Error(ctx).Err(err).Str("job_id", jobID).Msg("Failed to get bulk job items")
		return nil, err
	}

	items := make([]*entities.BulkInvoiceJobItem, len(itemModels))
	for i, model := range itemModels {
		items[i] = mappers.JobItemToDomain(&model)
	}

	return items, nil
}

// UpdateBulkInvoiceJob actualiza un job existente en la base de datos (todos los campos)
func (r *Repository) UpdateJob(ctx context.Context, job *entities.BulkInvoiceJob) error {
	jobUUID, err := uuid.Parse(job.ID)
	if err != nil {
		return err
	}

	jobModel := mappers.JobToModel(job)
	jobModel.ID = jobUUID

	if err := r.db.Conn(ctx).Save(&jobModel).Error; err != nil {
		r.log.Error(ctx).Err(err).Str("job_id", job.ID).Msg("Failed to update bulk job")
		return err
	}

	return nil
}

// UpdateBulkInvoiceJobItem actualiza un item específico de un job en la base de datos
func (r *Repository) UpdateJobItem(ctx context.Context, item *entities.BulkInvoiceJobItem) error {
	itemModel := mappers.JobItemToModel(item)

	if err := r.db.Conn(ctx).Save(&itemModel).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("item_id", item.ID).Msg("Failed to update bulk job item")
		return err
	}

	return nil
}

// ListBulkInvoiceJobsByBusinessID obtiene jobs paginados de un negocio desde la base de datos
func (r *Repository) ListJobs(ctx context.Context, businessID uint, page, pageSize int) ([]*entities.BulkInvoiceJob, int64, error) {
	var jobModels []models.BulkInvoiceJob
	var total int64

	query := r.db.Conn(ctx).Model(&models.BulkInvoiceJob{}).Where("business_id = ?", businessID)

	// Contar total de registros
	if err := query.Count(&total).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("business_id", businessID).Msg("Failed to count bulk jobs")
		return nil, 0, err
	}

	// Obtener página específica
	offset := (page - 1) * pageSize
	if err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&jobModels).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("business_id", businessID).Msg("Failed to list bulk jobs")
		return nil, 0, err
	}

	jobs := make([]*entities.BulkInvoiceJob, len(jobModels))
	for i, model := range jobModels {
		jobs[i] = mappers.JobToDomain(&model)
	}

	return jobs, total, nil
}

// IncrementBulkInvoiceJobCounters actualiza contadores de un job de forma atómica en la BD (evita race conditions)
func (r *Repository) IncrementJobCounters(ctx context.Context, jobID string, processed, successful, failed int) error {
	jobUUID, err := uuid.Parse(jobID)
	if err != nil {
		return err
	}

	result := r.db.Conn(ctx).
		Model(&models.BulkInvoiceJob{}).
		Where("id = ?", jobUUID).
		UpdateColumns(map[string]interface{}{
			"processed":  gorm.Expr("processed + ?", processed),
			"successful": gorm.Expr("successful + ?", successful),
			"failed":     gorm.Expr("failed + ?", failed),
		})

	if result.Error != nil {
		r.log.Error(ctx).
			Err(result.Error).
			Str("job_id", jobID).
			Msg("Failed to increment job counters")
		return result.Error
	}

	r.log.Debug(ctx).
		Str("job_id", jobID).
		Int("processed", processed).
		Int("successful", successful).
		Int("failed", failed).
		Msg("Job counters incremented")

	return nil
}
