package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// CreateIntegrationType crea un nuevo tipo de integración
func (r *Repository) CreateIntegrationType(ctx context.Context, integrationType *domain.IntegrationType) error {
	model := toIntegrationTypeModel(integrationType)
	if err := r.db.Conn(ctx).Create(&model).Error; err != nil {
		r.log.Error(ctx).Err(err).Msg("Error al crear tipo de integración")
		return fmt.Errorf("error al crear tipo de integración: %w", err)
	}
	*integrationType = toIntegrationTypeDomain(model)
	return nil
}

// UpdateIntegrationType actualiza un tipo de integración
func (r *Repository) UpdateIntegrationType(ctx context.Context, id uint, integrationType *domain.IntegrationType) error {
	model := toIntegrationTypeModel(integrationType)
	if err := r.db.Conn(ctx).Model(&models.IntegrationType{}).Where("id = ?", id).Updates(&model).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("id", id).Msg("Error al actualizar tipo de integración")
		return fmt.Errorf("error al actualizar tipo de integración: %w", err)
	}
	*integrationType = toIntegrationTypeDomain(model)
	return nil
}

// GetIntegrationTypeByID obtiene un tipo de integración por ID
func (r *Repository) GetIntegrationTypeByID(ctx context.Context, id uint) (*domain.IntegrationType, error) {
	var model models.IntegrationType
	if err := r.db.Conn(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("tipo de integración con ID %d no encontrado", id)
		}
		r.log.Error(ctx).Err(err).Uint("id", id).Msg("Error al obtener tipo de integración por ID")
		return nil, fmt.Errorf("error al obtener tipo de integración: %w", err)
	}
	result := toIntegrationTypeDomain(model)
	return &result, nil
}

// GetIntegrationTypeByCode obtiene un tipo de integración por código
func (r *Repository) GetIntegrationTypeByCode(ctx context.Context, code string) (*domain.IntegrationType, error) {
	var model models.IntegrationType
	if err := r.db.Conn(ctx).Where("code = ?", code).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("tipo de integración con código '%s' no encontrado", code)
		}
		r.log.Error(ctx).Err(err).Str("code", code).Msg("Error al obtener tipo de integración por código")
		return nil, fmt.Errorf("error al obtener tipo de integración: %w", err)
	}
	result := toIntegrationTypeDomain(model)
	return &result, nil
}

// GetIntegrationTypeByName obtiene un tipo de integración por nombre
func (r *Repository) GetIntegrationTypeByName(ctx context.Context, name string) (*domain.IntegrationType, error) {
	var model models.IntegrationType
	if err := r.db.Conn(ctx).Where("name = ?", name).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("tipo de integración con nombre '%s' no encontrado", name)
		}
		r.log.Error(ctx).Err(err).Str("name", name).Msg("Error al obtener tipo de integración por nombre")
		return nil, fmt.Errorf("error al obtener tipo de integración: %w", err)
	}
	result := toIntegrationTypeDomain(model)
	return &result, nil
}

// DeleteIntegrationType elimina un tipo de integración
func (r *Repository) DeleteIntegrationType(ctx context.Context, id uint) error {
	if err := r.db.Conn(ctx).Delete(&models.IntegrationType{}, id).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("id", id).Msg("Error al eliminar tipo de integración")
		return fmt.Errorf("error al eliminar tipo de integración: %w", err)
	}
	return nil
}

// ListIntegrationTypes obtiene todos los tipos de integración
func (r *Repository) ListIntegrationTypes(ctx context.Context) ([]*domain.IntegrationType, error) {
	var models []models.IntegrationType
	if err := r.db.Conn(ctx).Order("name ASC").Find(&models).Error; err != nil {
		r.log.Error(ctx).Err(err).Msg("Error al listar tipos de integración")
		return nil, fmt.Errorf("error al listar tipos de integración: %w", err)
	}
	result := make([]*domain.IntegrationType, len(models))
	for i := range models {
		domainType := toIntegrationTypeDomain(models[i])
		result[i] = &domainType
	}
	return result, nil
}

// ListActiveIntegrationTypes obtiene solo los tipos de integración activos
func (r *Repository) ListActiveIntegrationTypes(ctx context.Context) ([]*domain.IntegrationType, error) {
	var models []models.IntegrationType
	if err := r.db.Conn(ctx).Where("is_active = ?", true).Order("name ASC").Find(&models).Error; err != nil {
		r.log.Error(ctx).Err(err).Msg("Error al listar tipos de integración activos")
		return nil, fmt.Errorf("error al listar tipos de integración activos: %w", err)
	}
	result := make([]*domain.IntegrationType, len(models))
	for i := range models {
		domainType := toIntegrationTypeDomain(models[i])
		result[i] = &domainType
	}
	return result, nil
}

// toIntegrationTypeModel convierte domain.IntegrationType a models.IntegrationType
func toIntegrationTypeModel(d *domain.IntegrationType) models.IntegrationType {
	return models.IntegrationType{
		Model: gorm.Model{
			ID:        d.ID,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		},
		Name:              d.Name,
		Code:              d.Code,
		Description:       d.Description,
		Icon:              d.Icon,
		ImageURL:          d.ImageURL,
		Category:          d.Category,
		IsActive:          d.IsActive,
		ConfigSchema:      d.ConfigSchema,
		CredentialsSchema: d.CredentialsSchema,
		SetupInstructions: d.SetupInstructions,
	}
}

// toIntegrationTypeDomain convierte models.IntegrationType a domain.IntegrationType
func toIntegrationTypeDomain(m models.IntegrationType) domain.IntegrationType {
	return domain.IntegrationType{
		ID:                m.ID,
		Name:              m.Name,
		Code:              m.Code,
		Description:       m.Description,
		Icon:              m.Icon,
		ImageURL:          m.ImageURL,
		Category:          m.Category,
		IsActive:          m.IsActive,
		ConfigSchema:      m.ConfigSchema,
		CredentialsSchema: m.CredentialsSchema,
		SetupInstructions: m.SetupInstructions,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}
