package mappers

import (
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// ═══════════════════════════════════════════════════════════════
// INVOICING PROVIDER TYPE
// ═══════════════════════════════════════════════════════════════

func ProviderTypeToDomain(model *models.InvoicingProviderType) *entities.InvoicingProviderType {
	if model == nil {
		return nil
	}

	entity := &entities.InvoicingProviderType{
		ID:                 model.ID,
		CreatedAt:          model.CreatedAt,
		UpdatedAt:          model.UpdatedAt,
		Name:               model.Name,
		Code:               model.Code,
		Description:        model.Description,
		Icon:               model.Icon,
		ImageURL:           model.ImageURL,
		ApiBaseURL:         model.ApiBaseURL,
		DocumentationURL:   model.DocumentationURL,
		IsActive:           model.IsActive,
		SupportedCountries: model.SupportedCountries,
	}

	if model.DeletedAt.Valid {
		entity.DeletedAt = &model.DeletedAt.Time
	}

	return entity
}

func ProviderTypeListToDomain(models []*models.InvoicingProviderType) []*entities.InvoicingProviderType {
	entities := make([]*entities.InvoicingProviderType, 0, len(models))
	for _, model := range models {
		entities = append(entities, ProviderTypeToDomain(model))
	}
	return entities
}

// ═══════════════════════════════════════════════════════════════
// INVOICING PROVIDER
// ═══════════════════════════════════════════════════════════════

func ProviderToDomain(model *models.InvoicingProvider) *entities.InvoicingProvider {
	if model == nil {
		return nil
	}

	entity := &entities.InvoicingProvider{
		ID:             model.ID,
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
		BusinessID:     model.BusinessID,
		ProviderTypeID: model.ProviderTypeID,
		Name:           model.Name,
		Description:    model.Description,
		IsActive:       model.IsActive,
		IsDefault:      model.IsDefault,
		CreatedByID:    model.CreatedByID,
		UpdatedByID:    model.UpdatedByID,
	}

	if model.DeletedAt.Valid {
		entity.DeletedAt = &model.DeletedAt.Time
	}

	// Convertir JSONB a map
	if model.Config != nil {
		var config map[string]interface{}
		if err := json.Unmarshal(model.Config, &config); err == nil {
			entity.Config = config
		}
	}

	if model.Credentials != nil {
		var credentials map[string]interface{}
		if err := json.Unmarshal(model.Credentials, &credentials); err == nil {
			entity.Credentials = credentials
		}
	}

	return entity
}

func ProviderToModel(entity *entities.InvoicingProvider) *models.InvoicingProvider {
	if entity == nil {
		return nil
	}

	model := &models.InvoicingProvider{
		Model: gorm.Model{
			ID:        entity.ID,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},
		BusinessID:     entity.BusinessID,
		ProviderTypeID: entity.ProviderTypeID,
		Name:           entity.Name,
		Description:    entity.Description,
		IsActive:       entity.IsActive,
		IsDefault:      entity.IsDefault,
		CreatedByID:    entity.CreatedByID,
		UpdatedByID:    entity.UpdatedByID,
	}

	if entity.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *entity.DeletedAt, Valid: true}
	}

	// Convertir map a JSONB
	if entity.Config != nil {
		if data, err := json.Marshal(entity.Config); err == nil {
			model.Config = datatypes.JSON(data)
		}
	}

	if entity.Credentials != nil {
		if data, err := json.Marshal(entity.Credentials); err == nil {
			model.Credentials = datatypes.JSON(data)
		}
	}

	return model
}

func ProviderListToDomain(models []*models.InvoicingProvider) []*entities.InvoicingProvider {
	entities := make([]*entities.InvoicingProvider, 0, len(models))
	for _, model := range models {
		entities = append(entities, ProviderToDomain(model))
	}
	return entities
}

// ═══════════════════════════════════════════════════════════════
// INVOICING CONFIG
// ═══════════════════════════════════════════════════════════════

func ConfigToDomain(model *models.InvoicingConfig) *entities.InvoicingConfig {
	if model == nil {
		return nil
	}

	entity := &entities.InvoicingConfig{
		ID:                     model.ID,
		CreatedAt:              model.CreatedAt,
		UpdatedAt:              model.UpdatedAt,
		BusinessID:             model.BusinessID,
		IntegrationID:          model.IntegrationID,
		InvoicingProviderID:    model.InvoicingProviderID,    // Campo deprecado (legacy)
		InvoicingIntegrationID: model.InvoicingIntegrationID, // Campo nuevo (actual)
		Enabled:                model.Enabled,
		AutoInvoice:            model.AutoInvoice,
		Description:            model.Description,
		CreatedByID:            model.CreatedByID,
		UpdatedByID:            model.UpdatedByID,
	}

	if model.DeletedAt.Valid {
		entity.DeletedAt = &model.DeletedAt.Time
	}

	// Extraer nombre de la integración de e-commerce si está preloaded
	if model.Integration.ID > 0 && model.Integration.Name != "" {
		entity.IntegrationName = &model.Integration.Name
	}

	// Extraer nombre y logo de la integración de facturación (Softpymes) si está preloaded
	if model.InvoicingIntegration.ID > 0 && model.InvoicingIntegration.Name != "" {
		entity.ProviderName = &model.InvoicingIntegration.Name

		// Extraer IsTesting desde la integración de facturación
		entity.IsTesting = model.InvoicingIntegration.IsTesting

		// Extraer logo, BaseURL y BaseURLTest del IntegrationType si está preloaded
		if model.InvoicingIntegration.IntegrationType != nil && model.InvoicingIntegration.IntegrationType.ID > 0 {
			if model.InvoicingIntegration.IntegrationType.ImageURL != "" {
				entity.ProviderImageURL = &model.InvoicingIntegration.IntegrationType.ImageURL
			}
			entity.BaseURL = model.InvoicingIntegration.IntegrationType.BaseURL
			entity.BaseURLTest = model.InvoicingIntegration.IntegrationType.BaseURLTest
		}
	}

	// Convertir JSONB a map
	if model.Filters != nil {
		var filters map[string]interface{}
		if err := json.Unmarshal(model.Filters, &filters); err == nil {
			entity.Filters = filters
		}
	}

	if model.InvoiceConfig != nil {
		var config map[string]interface{}
		if err := json.Unmarshal(model.InvoiceConfig, &config); err == nil {
			entity.InvoiceConfig = config
		}
	}

	return entity
}

func ConfigToModel(entity *entities.InvoicingConfig) *models.InvoicingConfig {
	if entity == nil {
		return nil
	}

	model := &models.InvoicingConfig{
		Model: gorm.Model{
			ID:        entity.ID,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},
		BusinessID:             entity.BusinessID,
		IntegrationID:          entity.IntegrationID,
		InvoicingProviderID:    entity.InvoicingProviderID,    // Campo deprecado (legacy)
		InvoicingIntegrationID: entity.InvoicingIntegrationID, // Campo nuevo (actual)
		Enabled:                entity.Enabled,
		AutoInvoice:            entity.AutoInvoice,
		Description:            entity.Description,
		CreatedByID:            entity.CreatedByID,
		UpdatedByID:            entity.UpdatedByID,
	}

	if entity.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *entity.DeletedAt, Valid: true}
	}

	// Convertir map a JSONB
	if entity.Filters != nil {
		if data, err := json.Marshal(entity.Filters); err == nil {
			model.Filters = datatypes.JSON(data)
		}
	}

	if entity.InvoiceConfig != nil {
		if data, err := json.Marshal(entity.InvoiceConfig); err == nil {
			model.InvoiceConfig = datatypes.JSON(data)
		}
	}

	return model
}

func ConfigListToDomain(models []*models.InvoicingConfig) []*entities.InvoicingConfig {
	entities := make([]*entities.InvoicingConfig, 0, len(models))
	for _, model := range models {
		entities = append(entities, ConfigToDomain(model))
	}
	return entities
}
