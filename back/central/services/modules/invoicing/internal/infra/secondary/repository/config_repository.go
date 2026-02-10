package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// ═══════════════════════════════════════════
// INVOICING CONFIGS - Métodos del Repository
// ═══════════════════════════════════════════

// CreateInvoicingConfig crea una nueva configuración de facturación en la base de datos
func (r *Repository) CreateInvoicingConfig(ctx context.Context, config *entities.InvoicingConfig) error {
	model := mappers.ConfigToModel(config)

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		r.log.Error(ctx).Err(err).Msg("Failed to create invoicing config")
		return fmt.Errorf("failed to create config: %w", err)
	}

	config.ID = model.ID
	r.log.Info(ctx).Uint("config_id", config.ID).Msg("Invoicing config created")

	// Guardar en caché en background
	go r.configCache.Set(context.Background(), config)

	return nil
}

// GetInvoicingConfigByID obtiene una configuración de facturación por su ID desde la base de datos
func (r *Repository) GetInvoicingConfigByID(ctx context.Context, id uint) (*entities.InvoicingConfig, error) {
	var model models.InvoicingConfig

	if err := r.db.Conn(ctx).First(&model, id).Error; err != nil {
		return nil, fmt.Errorf("config not found: %w", err)
	}

	return mappers.ConfigToDomain(&model), nil
}

// GetConfigByIntegration obtiene una configuración de facturación por ID de integración
// Implementa read-through cache: primero intenta desde Redis, luego desde BD
func (r *Repository) GetConfigByIntegration(ctx context.Context, integrationID uint) (*entities.InvoicingConfig, error) {
	// 1. Intentar desde caché (cache HIT → retornar inmediatamente)
	cachedConfig, err := r.configCache.Get(ctx, integrationID)
	if err != nil {
		// Error al leer caché (no debería pasar, pero loggear por si acaso)
		r.log.Warn(ctx).Err(err).Uint("integration_id", integrationID).Msg("Error al leer caché, consultando BD")
	}
	if cachedConfig != nil {
		return cachedConfig, nil
	}

	// 2. Cache MISS - consultar base de datos
	var model models.InvoicingConfig

	if err := r.db.Conn(ctx).Where("integration_id = ?", integrationID).First(&model).Error; err != nil {
		return nil, fmt.Errorf("config not found for integration: %w", err)
	}

	config := mappers.ConfigToDomain(&model)

	// 3. Actualizar caché en background (no bloquear el retorno)
	go r.configCache.Set(context.Background(), config)

	return config, nil
}

// ListInvoicingConfigs lista todas las configuraciones de facturación de un negocio desde la base de datos
func (r *Repository) ListInvoicingConfigs(ctx context.Context, businessID uint) ([]*entities.InvoicingConfig, error) {
	var configModels []*models.InvoicingConfig

	if err := r.db.Conn(ctx).
		Preload("Integration").                            // Cargar integración de e-commerce
		Preload("InvoicingIntegration").                   // Cargar integración de facturación (Softpymes)
		Preload("InvoicingIntegration.IntegrationType").  // Cargar tipo para obtener el logo
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

	// Usar Updates() con Omit() para no actualizar campos de auditoría
	// Save() actualiza TODOS los campos, incluyendo created_by_id con valor 0 (violación FK)
	if err := r.db.Conn(ctx).Model(&models.InvoicingConfig{}).
		Where("id = ?", model.ID).
		Omit("created_by_id", "created_at").
		Updates(model).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("config_id", config.ID).Msg("Failed to update config")
		return fmt.Errorf("failed to update config: %w", err)
	}

	r.log.Info(ctx).Uint("config_id", config.ID).Msg("Invoicing config updated")

	// Invalidar caché para forzar recarga en próxima consulta
	go r.configCache.Invalidate(context.Background(), config.IntegrationID)

	// Actualizar caché con nuevos valores
	go r.configCache.Set(context.Background(), config)

	return nil
}

// DeleteInvoicingConfig elimina PERMANENTEMENTE una configuración de facturación (hard delete)
func (r *Repository) DeleteInvoicingConfig(ctx context.Context, id uint) error {
	// Primero obtener la config para saber el integration_id (necesario para invalidar caché)
	// Usar Unscoped() para buscar incluso si ya está soft-deleted
	var model models.InvoicingConfig
	if err := r.db.Conn(ctx).Unscoped().First(&model, id).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("config_id", id).Msg("Config not found for deletion")
		return fmt.Errorf("config not found: %w", err)
	}

	// HARD DELETE - eliminar permanentemente usando Unscoped()
	if err := r.db.Conn(ctx).Unscoped().Delete(&model).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("config_id", id).Msg("Failed to delete config")
		return fmt.Errorf("failed to delete config: %w", err)
	}

	r.log.Info(ctx).Uint("config_id", id).Msg("Invoicing config permanently deleted (hard delete)")

	// Invalidar caché en background
	go r.configCache.Invalidate(context.Background(), model.IntegrationID)

	return nil
}

// ConfigExistsForIntegration verifica si existe una configuración de facturación para una integración específica
func (r *Repository) ConfigExistsForIntegration(ctx context.Context, integrationID uint) (bool, error) {
	var count int64

	if err := r.db.Conn(ctx).Model(&models.InvoicingConfig{}).
		Where("integration_id = ?", integrationID).
		Where("deleted_at IS NULL").
		Count(&count).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("integration_id", integrationID).Msg("Failed to check config existence")
		return false, fmt.Errorf("failed to check config existence: %w", err)
	}

	return count > 0, nil
}
