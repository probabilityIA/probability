package mappers

import (
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// ═══════════════════════════════════════════════════════════════
// INVOICE ITEM
// ═══════════════════════════════════════════════════════════════

func InvoiceItemToDomain(model *models.InvoiceItem) *entities.InvoiceItem {
	if model == nil {
		return nil
	}

	entity := &entities.InvoiceItem{
		ID:             model.ID,
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
		InvoiceID:      model.InvoiceID,
		ProductID:      model.ProductID,
		SKU:            model.SKU,
		Name:           model.Name,
		Description:    model.Description,
		Quantity:       model.Quantity,
		UnitPrice:      model.UnitPrice,
		TotalPrice:     model.TotalPrice,
		Currency:       model.Currency,
		Tax:            model.Tax,
		TaxRate:        model.TaxRate,
		Discount:       model.Discount,
		ProviderItemID: model.ProviderItemID,
	}

	if model.DeletedAt.Valid {
		entity.DeletedAt = &model.DeletedAt.Time
	}

	if model.Metadata != nil {
		var metadata map[string]interface{}
		if err := json.Unmarshal(model.Metadata, &metadata); err == nil {
			entity.Metadata = metadata
		}
	}

	return entity
}

func InvoiceItemToModel(entity *entities.InvoiceItem) *models.InvoiceItem {
	if entity == nil {
		return nil
	}

	model := &models.InvoiceItem{
		Model: gorm.Model{
			ID:        entity.ID,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},
		InvoiceID:      entity.InvoiceID,
		ProductID:      entity.ProductID,
		SKU:            entity.SKU,
		Name:           entity.Name,
		Description:    entity.Description,
		Quantity:       entity.Quantity,
		UnitPrice:      entity.UnitPrice,
		TotalPrice:     entity.TotalPrice,
		Currency:       entity.Currency,
		Tax:            entity.Tax,
		TaxRate:        entity.TaxRate,
		Discount:       entity.Discount,
		ProviderItemID: entity.ProviderItemID,
	}

	if entity.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *entity.DeletedAt, Valid: true}
	}

	if entity.Metadata != nil {
		if data, err := json.Marshal(entity.Metadata); err == nil {
			model.Metadata = datatypes.JSON(data)
		}
	}

	return model
}

func InvoiceItemListToDomain(models []*models.InvoiceItem) []*entities.InvoiceItem {
	entities := make([]*entities.InvoiceItem, 0, len(models))
	for _, model := range models {
		entities = append(entities, InvoiceItemToDomain(model))
	}
	return entities
}

// ═══════════════════════════════════════════════════════════════
// INVOICE SYNC LOG
// ═══════════════════════════════════════════════════════════════

func SyncLogToDomain(model *models.InvoiceSyncLog) *entities.InvoiceSyncLog {
	if model == nil {
		return nil
	}

	entity := &entities.InvoiceSyncLog{
		ID:            model.ID,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
		InvoiceID:     model.InvoiceID,
		OperationType: model.OperationType,
		Status:        model.Status,
		RequestURL:    model.RequestURL,
		ResponseStatus: model.ResponseStatus,
		ErrorMessage:  model.ErrorMessage,
		ErrorCode:     model.ErrorCode,
		RetryCount:    model.RetryCount,
		NextRetryAt:   model.NextRetryAt,
		MaxRetries:    model.MaxRetries,
		RetriedAt:     model.RetriedAt,
		StartedAt:     model.StartedAt,
		CompletedAt:   model.CompletedAt,
		Duration:      model.Duration,
		TriggeredBy:   model.TriggeredBy,
		UserID:        model.UserID,
	}

	if model.DeletedAt.Valid {
		entity.DeletedAt = &model.DeletedAt.Time
	}

	// Convertir JSONB a map
	if model.RequestPayload != nil {
		var payload map[string]interface{}
		if err := json.Unmarshal(model.RequestPayload, &payload); err == nil {
			entity.RequestPayload = payload
		}
	}

	if model.RequestHeaders != nil {
		var headers map[string]interface{}
		if err := json.Unmarshal(model.RequestHeaders, &headers); err == nil {
			entity.RequestHeaders = headers
		}
	}

	if model.ResponseBody != nil {
		var body map[string]interface{}
		if err := json.Unmarshal(model.ResponseBody, &body); err == nil {
			entity.ResponseBody = body
		}
	}

	if model.ResponseHeaders != nil {
		var headers map[string]interface{}
		if err := json.Unmarshal(model.ResponseHeaders, &headers); err == nil {
			entity.ResponseHeaders = headers
		}
	}

	if model.ErrorDetails != nil {
		var details map[string]interface{}
		if err := json.Unmarshal(model.ErrorDetails, &details); err == nil {
			entity.ErrorDetails = details
		}
	}

	return entity
}

func SyncLogToModel(entity *entities.InvoiceSyncLog) *models.InvoiceSyncLog {
	if entity == nil {
		return nil
	}

	model := &models.InvoiceSyncLog{
		Model: gorm.Model{
			ID:        entity.ID,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},
		InvoiceID:      entity.InvoiceID,
		OperationType:  entity.OperationType,
		Status:         entity.Status,
		RequestURL:     entity.RequestURL,
		ResponseStatus: entity.ResponseStatus,
		ErrorMessage:   entity.ErrorMessage,
		ErrorCode:      entity.ErrorCode,
		RetryCount:     entity.RetryCount,
		NextRetryAt:    entity.NextRetryAt,
		MaxRetries:     entity.MaxRetries,
		RetriedAt:      entity.RetriedAt,
		StartedAt:      entity.StartedAt,
		CompletedAt:    entity.CompletedAt,
		Duration:       entity.Duration,
		TriggeredBy:    entity.TriggeredBy,
		UserID:         entity.UserID,
	}

	if entity.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *entity.DeletedAt, Valid: true}
	}

	// Convertir map a JSONB
	if entity.RequestPayload != nil {
		if data, err := json.Marshal(entity.RequestPayload); err == nil {
			model.RequestPayload = datatypes.JSON(data)
		}
	}

	if entity.RequestHeaders != nil {
		if data, err := json.Marshal(entity.RequestHeaders); err == nil {
			model.RequestHeaders = datatypes.JSON(data)
		}
	}

	if entity.ResponseBody != nil {
		if data, err := json.Marshal(entity.ResponseBody); err == nil {
			model.ResponseBody = datatypes.JSON(data)
		}
	}

	if entity.ResponseHeaders != nil {
		if data, err := json.Marshal(entity.ResponseHeaders); err == nil {
			model.ResponseHeaders = datatypes.JSON(data)
		}
	}

	if entity.ErrorDetails != nil {
		if data, err := json.Marshal(entity.ErrorDetails); err == nil {
			model.ErrorDetails = datatypes.JSON(data)
		}
	}

	return model
}

func SyncLogListToDomain(models []*models.InvoiceSyncLog) []*entities.InvoiceSyncLog {
	entities := make([]*entities.InvoiceSyncLog, 0, len(models))
	for _, model := range models {
		entities = append(entities, SyncLogToDomain(model))
	}
	return entities
}

// ═══════════════════════════════════════════════════════════════
// CREDIT NOTE
// ═══════════════════════════════════════════════════════════════

func CreditNoteToDomain(model *models.CreditNote) *entities.CreditNote {
	if model == nil {
		return nil
	}

	entity := &entities.CreditNote{
		ID:               model.ID,
		CreatedAt:        model.CreatedAt,
		UpdatedAt:        model.UpdatedAt,
		InvoiceID:        model.InvoiceID,
		BusinessID:       model.BusinessID,
		CreditNoteNumber: model.CreditNoteNumber,
		ExternalID:       model.ExternalID,
		InternalNumber:   model.InternalNumber,
		NoteType:         model.NoteType,
		Amount:           model.Amount,
		Currency:         model.Currency,
		Reason:           model.Reason,
		Description:      model.Description,
		Status:           model.Status,
		IssuedAt:         model.IssuedAt,
		CancelledAt:      model.CancelledAt,
		NoteURL:          model.NoteURL,
		PDFURL:           model.PDFURL,
		XMLURL:           model.XMLURL,
		CUFE:             model.CUFE,
		CreatedByID:      model.CreatedByID,
	}

	if model.DeletedAt.Valid {
		entity.DeletedAt = &model.DeletedAt.Time
	}

	// Convertir JSONB a map
	if model.Metadata != nil {
		var metadata map[string]interface{}
		if err := json.Unmarshal(model.Metadata, &metadata); err == nil {
			entity.Metadata = metadata
		}
	}

	if model.ProviderResponse != nil {
		var response map[string]interface{}
		if err := json.Unmarshal(model.ProviderResponse, &response); err == nil {
			entity.ProviderResponse = response
		}
	}

	return entity
}

func CreditNoteToModel(entity *entities.CreditNote) *models.CreditNote {
	if entity == nil {
		return nil
	}

	model := &models.CreditNote{
		Model: gorm.Model{
			ID:        entity.ID,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},
		InvoiceID:        entity.InvoiceID,
		BusinessID:       entity.BusinessID,
		CreditNoteNumber: entity.CreditNoteNumber,
		ExternalID:       entity.ExternalID,
		InternalNumber:   entity.InternalNumber,
		NoteType:         entity.NoteType,
		Amount:           entity.Amount,
		Currency:         entity.Currency,
		Reason:           entity.Reason,
		Description:      entity.Description,
		Status:           entity.Status,
		IssuedAt:         entity.IssuedAt,
		CancelledAt:      entity.CancelledAt,
		NoteURL:          entity.NoteURL,
		PDFURL:           entity.PDFURL,
		XMLURL:           entity.XMLURL,
		CUFE:             entity.CUFE,
		CreatedByID:      entity.CreatedByID,
	}

	if entity.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *entity.DeletedAt, Valid: true}
	}

	// Convertir map a JSONB
	if entity.Metadata != nil {
		if data, err := json.Marshal(entity.Metadata); err == nil {
			model.Metadata = datatypes.JSON(data)
		}
	}

	if entity.ProviderResponse != nil {
		if data, err := json.Marshal(entity.ProviderResponse); err == nil {
			model.ProviderResponse = datatypes.JSON(data)
		}
	}

	return model
}

func CreditNoteListToDomain(models []*models.CreditNote) []*entities.CreditNote {
	entities := make([]*entities.CreditNote, 0, len(models))
	for _, model := range models {
		entities = append(entities, CreditNoteToDomain(model))
	}
	return entities
}
