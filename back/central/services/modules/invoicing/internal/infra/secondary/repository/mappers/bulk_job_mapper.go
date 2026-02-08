package mappers

import (
	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// JobToDomain convierte un modelo GORM de BulkInvoiceJob a entidad de dominio
func JobToDomain(model *models.BulkInvoiceJob) *entities.BulkInvoiceJob {
	if model == nil {
		return nil
	}

	return &entities.BulkInvoiceJob{
		ID:              model.ID.String(),
		BusinessID:      model.BusinessID,
		CreatedByUserID: model.CreatedByUserID,
		TotalOrders:     model.TotalOrders,
		Processed:       model.Processed,
		Successful:      model.Successful,
		Failed:          model.Failed,
		Status:          model.Status,
		StartedAt:       model.StartedAt,
		CompletedAt:     model.CompletedAt,
		ErrorMessage:    model.ErrorMessage,
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
	}
}

// JobToModel convierte una entidad de dominio BulkInvoiceJob a modelo GORM
func JobToModel(entity *entities.BulkInvoiceJob) models.BulkInvoiceJob {
	model := models.BulkInvoiceJob{
		BusinessID:      entity.BusinessID,
		CreatedByUserID: entity.CreatedByUserID,
		TotalOrders:     entity.TotalOrders,
		Processed:       entity.Processed,
		Successful:      entity.Successful,
		Failed:          entity.Failed,
		Status:          entity.Status,
		StartedAt:       entity.StartedAt,
		CompletedAt:     entity.CompletedAt,
		ErrorMessage:    entity.ErrorMessage,
		CreatedAt:       entity.CreatedAt,
		UpdatedAt:       entity.UpdatedAt,
	}

	// Parsear UUID si existe
	if entity.ID != "" {
		if id, err := uuid.Parse(entity.ID); err == nil {
			model.ID = id
		}
	}

	return model
}

// JobItemToDomain convierte un modelo GORM de BulkInvoiceJobItem a entidad de dominio
func JobItemToDomain(model *models.BulkInvoiceJobItem) *entities.BulkInvoiceJobItem {
	if model == nil {
		return nil
	}

	return &entities.BulkInvoiceJobItem{
		ID:           model.ID,
		JobID:        model.JobID.String(),
		OrderID:      model.OrderID,
		InvoiceID:    model.InvoiceID,
		Status:       model.Status,
		ErrorMessage: model.ErrorMessage,
		ProcessedAt:  model.ProcessedAt,
		CreatedAt:    model.CreatedAt,
	}
}

// JobItemToModel convierte una entidad de dominio BulkInvoiceJobItem a modelo GORM
func JobItemToModel(entity *entities.BulkInvoiceJobItem) models.BulkInvoiceJobItem {
	jobUUID, _ := uuid.Parse(entity.JobID)

	return models.BulkInvoiceJobItem{
		ID:           entity.ID,
		JobID:        jobUUID,
		OrderID:      entity.OrderID,
		InvoiceID:    entity.InvoiceID,
		Status:       entity.Status,
		ErrorMessage: entity.ErrorMessage,
		ProcessedAt:  entity.ProcessedAt,
		CreatedAt:    entity.CreatedAt,
	}
}
