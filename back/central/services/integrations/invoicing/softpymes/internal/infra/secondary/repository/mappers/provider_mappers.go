package mappers

import (
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/repository/models"
)

// ═══════════════════════════════════════════════════════════════
// PROVIDER TYPE - Tipo de proveedor
// ═══════════════════════════════════════════════════════════════

// ProviderTypeToDomain convierte modelo GORM a entidad de dominio
func ProviderTypeToDomain(model *models.ProviderType) *entities.ProviderType {
	if model == nil {
		return nil
	}

	entity := &entities.ProviderType{
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

// ProviderTypeListToDomain convierte lista de modelos a entidades
func ProviderTypeListToDomain(models []*models.ProviderType) []*entities.ProviderType {
	entities := make([]*entities.ProviderType, 0, len(models))
	for _, model := range models {
		entities = append(entities, ProviderTypeToDomain(model))
	}
	return entities
}

// ═══════════════════════════════════════════════════════════════
// PROVIDER - Proveedor de facturación
// ═══════════════════════════════════════════════════════════════

// ProviderToDomain convierte modelo GORM a entidad de dominio
func ProviderToDomain(model *models.Provider) *entities.Provider {
	if model == nil {
		return nil
	}

	entity := &entities.Provider{
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

// ProviderToModel convierte entidad de dominio a modelo GORM
func ProviderToModel(entity *entities.Provider) *models.Provider {
	if entity == nil {
		return nil
	}

	model := &models.Provider{
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

// ProviderListToDomain convierte lista de modelos a entidades
func ProviderListToDomain(models []*models.Provider) []*entities.Provider {
	entities := make([]*entities.Provider, 0, len(models))
	for _, model := range models {
		entities = append(entities, ProviderToDomain(model))
	}
	return entities
}
