package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
<<<<<<< HEAD
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
=======
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
)

<<<<<<< HEAD
type InvoicingConfigRepository struct {
	*Repository
}

func NewInvoicingConfigRepository(repo *Repository) ports.IInvoicingConfigRepository {
	return &InvoicingConfigRepository{Repository: repo}
}

func (r *InvoicingConfigRepository) Create(ctx context.Context, config *entities.InvoicingConfig) error {
	model := mappers.ConfigToModel(config)

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
=======
// ═══════════════════════════════════════════
// INVOICING CONFIGS - Métodos del Repository
// ═══════════════════════════════════════════

// CreateInvoicingConfig crea una nueva configuración de facturación en la base de datos
func (r *Repository) CreateInvoicingConfig(ctx context.Context, config *entities.InvoicingConfig) error {
	model := mappers.ConfigToModel(config)

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		r.log.Error(ctx).Err(err).Msg("Failed to create invoicing config")
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
		return fmt.Errorf("failed to create config: %w", err)
	}

	config.ID = model.ID
<<<<<<< HEAD
	return nil
}

func (r *InvoicingConfigRepository) GetByID(ctx context.Context, id uint) (*entities.InvoicingConfig, error) {
=======
	r.log.Info(ctx).Uint("config_id", config.ID).Msg("Invoicing config created")
	return nil
}

// GetInvoicingConfigByID obtiene una configuración de facturación por su ID desde la base de datos
func (r *Repository) GetInvoicingConfigByID(ctx context.Context, id uint) (*entities.InvoicingConfig, error) {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	var model models.InvoicingConfig

	if err := r.db.Conn(ctx).First(&model, id).Error; err != nil {
		return nil, fmt.Errorf("config not found: %w", err)
	}

	return mappers.ConfigToDomain(&model), nil
}

<<<<<<< HEAD
func (r *InvoicingConfigRepository) GetByIntegration(ctx context.Context, integrationID uint) (*entities.InvoicingConfig, error) {
=======
// GetConfigByIntegration obtiene una configuración de facturación por ID de integración desde la base de datos
func (r *Repository) GetConfigByIntegration(ctx context.Context, integrationID uint) (*entities.InvoicingConfig, error) {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	var model models.InvoicingConfig

	if err := r.db.Conn(ctx).Where("integration_id = ?", integrationID).First(&model).Error; err != nil {
		return nil, fmt.Errorf("config not found for integration: %w", err)
	}

	return mappers.ConfigToDomain(&model), nil
}

<<<<<<< HEAD
func (r *InvoicingConfigRepository) List(ctx context.Context, businessID uint) ([]*entities.InvoicingConfig, error) {
	var models []*models.InvoicingConfig

	if err := r.db.Conn(ctx).
		Preload("Integration").             // Cargar integración de e-commerce
		Preload("InvoicingIntegration").   // Cargar integración de facturación (Softpymes)
		Where("business_id = ?", businessID).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list configs: %w", err)
	}

	return mappers.ConfigListToDomain(models), nil
}

func (r *InvoicingConfigRepository) Update(ctx context.Context, config *entities.InvoicingConfig) error {
	model := mappers.ConfigToModel(config)

	if err := r.db.Conn(ctx).Save(model).Error; err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	return nil
}

func (r *InvoicingConfigRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.Conn(ctx).Delete(&models.InvoicingConfig{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete config: %w", err)
	}

	return nil
}

func (r *InvoicingConfigRepository) ExistsForIntegration(ctx context.Context, integrationID uint) (bool, error) {
=======
// ListInvoicingConfigs lista todas las configuraciones de facturación de un negocio desde la base de datos
func (r *Repository) ListInvoicingConfigs(ctx context.Context, businessID uint) ([]*entities.InvoicingConfig, error) {
	var configModels []*models.InvoicingConfig

	if err := r.db.Conn(ctx).
		Preload("Integration").            // Cargar integración de e-commerce
		Preload("InvoicingIntegration").   // Cargar integración de facturación (Softpymes)
		Where("business_id = ?", businessID).
		Order("created_at DESC").
		Find(&configModels).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("business_id", businessID).Msg("Failed to list configs")
		return nil, fmt.Errorf("failed to list configs: %w", err)
	}

	return mappers.ConfigListToDomain(configModels), nil
}

// UpdateInvoicingConfig actualiza una configuración de facturación existente en la base de datos
func (r *Repository) UpdateInvoicingConfig(ctx context.Context, config *entities.InvoicingConfig) error {
	model := mappers.ConfigToModel(config)

	if err := r.db.Conn(ctx).Save(model).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("config_id", config.ID).Msg("Failed to update config")
		return fmt.Errorf("failed to update config: %w", err)
	}

	r.log.Info(ctx).Uint("config_id", config.ID).Msg("Invoicing config updated")
	return nil
}

// DeleteInvoicingConfig elimina (soft delete) una configuración de facturación de la base de datos
func (r *Repository) DeleteInvoicingConfig(ctx context.Context, id uint) error {
	if err := r.db.Conn(ctx).Delete(&models.InvoicingConfig{}, id).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("config_id", id).Msg("Failed to delete config")
		return fmt.Errorf("failed to delete config: %w", err)
	}

	r.log.Info(ctx).Uint("config_id", id).Msg("Invoicing config deleted")
	return nil
}

// ConfigExistsForIntegration verifica si existe una configuración de facturación para una integración específica
func (r *Repository) ConfigExistsForIntegration(ctx context.Context, integrationID uint) (bool, error) {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	var count int64

	if err := r.db.Conn(ctx).Model(&models.InvoicingConfig{}).
		Where("integration_id = ?", integrationID).
<<<<<<< HEAD
		Count(&count).Error; err != nil {
=======
		Where("deleted_at IS NULL").
		Count(&count).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("integration_id", integrationID).Msg("Failed to check config existence")
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
		return false, fmt.Errorf("failed to check config existence: %w", err)
	}

	return count > 0, nil
}
