package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// ═══════════════════════════════════════════
// INVOICING CONFIGS - Métodos del Repository
// ═══════════════════════════════════════════

// CreateInvoicingConfig crea una nueva configuración de facturación en la base de datos
// y sus entradas en la tabla join invoicing_config_integrations.
func (r *Repository) CreateInvoicingConfig(ctx context.Context, config *entities.InvoicingConfig) error {
	model := mappers.ConfigToModel(config)

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		if isUniqueConstraintError(err) {
			r.log.Warn(ctx).
				Interface("integration_ids", config.IntegrationIDs).
				Msg("Duplicate config detected at DB level")
			return fmt.Errorf("config already exists: %w", err)
		}
		r.log.Error(ctx).Err(err).Msg("Failed to create invoicing config")
		return fmt.Errorf("failed to create config: %w", err)
	}

	config.ID = model.ID

	// Insertar entradas en la tabla join para cada integration ID
	for _, integrationID := range config.IntegrationIDs {
		joinEntry := &models.InvoicingConfigIntegration{
			ConfigID:      config.ID,
			IntegrationID: integrationID,
		}
		if err := r.db.Conn(ctx).Create(joinEntry).Error; err != nil {
			r.log.Error(ctx).Err(err).
				Uint("config_id", config.ID).
				Uint("integration_id", integrationID).
				Msg("Failed to create config-integration join entry")
			return fmt.Errorf("failed to link integration to config: %w", err)
		}
	}

	r.log.Info(ctx).
		Uint("config_id", config.ID).
		Interface("integration_ids", config.IntegrationIDs).
		Msg("Invoicing config created")

	// Guardar en caché por cada integration_id
	for _, integrationID := range config.IntegrationIDs {
		id := integrationID
		go r.configCache.Set(context.Background(), id, config)
	}

	return nil
}

// GetInvoicingConfigByID obtiene una configuración por su ID con todos los preloads necesarios
func (r *Repository) GetInvoicingConfigByID(ctx context.Context, id uint) (*entities.InvoicingConfig, error) {
	var model models.InvoicingConfig

	if err := r.db.Conn(ctx).
		Preload("ConfigIntegrations.Integration").
		Preload("InvoicingIntegration").
		Preload("InvoicingIntegration.IntegrationType").
		First(&model, id).Error; err != nil {
		return nil, fmt.Errorf("config not found: %w", err)
	}

	return mappers.ConfigToDomain(&model), nil
}

// GetConfigByIntegration obtiene una configuración por ID de integración de e-commerce
// Implementa read-through cache: primero intenta desde Redis, luego desde BD via join table
func (r *Repository) GetConfigByIntegration(ctx context.Context, integrationID uint) (*entities.InvoicingConfig, error) {
	// 1. Intentar desde caché (cache HIT → retornar inmediatamente)
	cachedConfig, err := r.configCache.Get(ctx, integrationID)
	if err != nil {
		r.log.Warn(ctx).Err(err).Uint("integration_id", integrationID).Msg("Error al leer caché, consultando BD")
	}
	if cachedConfig != nil {
		return cachedConfig, nil
	}

	// 2. Cache MISS - consultar base de datos via tabla join
	var model models.InvoicingConfig

	if err := r.db.Conn(ctx).
		Preload("ConfigIntegrations.Integration").
		Preload("InvoicingIntegration").
		Preload("InvoicingIntegration.IntegrationType").
		Joins("JOIN invoicing_config_integrations ici ON ici.config_id = invoicing_configs.id AND ici.deleted_at IS NULL").
		Where("ici.integration_id = ?", integrationID).
		Where("invoicing_configs.deleted_at IS NULL").
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying invoicing config: %w", err)
	}

	config := mappers.ConfigToDomain(&model)

	// 3. Actualizar caché en background
	go r.configCache.Set(context.Background(), integrationID, config)

	return config, nil
}

// GetConfigByInvoicingIntegration retorna una config por (business_id, invoicing_integration_id).
// Usa el nuevo constraint único parcial.
func (r *Repository) GetConfigByInvoicingIntegration(ctx context.Context, businessID uint, invoicingIntegrationID uint) (*entities.InvoicingConfig, error) {
	var model models.InvoicingConfig

	err := r.db.Conn(ctx).
		Preload("ConfigIntegrations.Integration").
		Preload("InvoicingIntegration").
		Preload("InvoicingIntegration.IntegrationType").
		Where("business_id = ?", businessID).
		Where("invoicing_integration_id = ?", invoicingIntegrationID).
		Where("deleted_at IS NULL").
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying invoicing config: %w", err)
	}

	return mappers.ConfigToDomain(&model), nil
}

// ListInvoicingConfigs lista todas las configuraciones de facturación de un negocio.
// Si businessID == 0 (super admin sin filtro), retorna todas las configs de todos los negocios.
func (r *Repository) ListInvoicingConfigs(ctx context.Context, businessID uint) ([]*entities.InvoicingConfig, error) {
	var configModels []*models.InvoicingConfig

	query := r.db.Conn(ctx).
		Preload("ConfigIntegrations.Integration").
		Preload("InvoicingIntegration").
		Preload("InvoicingIntegration.IntegrationType").
		Order("created_at DESC")

	if businessID > 0 {
		query = query.Where("business_id = ?", businessID)
	}

	if err := query.Find(&configModels).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("business_id", businessID).Msg("Failed to list configs")
		return nil, fmt.Errorf("failed to list configs: %w", err)
	}

	return mappers.ConfigListToDomain(configModels), nil
}

// UpdateInvoicingConfig actualiza una configuración existente y sincroniza la tabla join.
func (r *Repository) UpdateInvoicingConfig(ctx context.Context, config *entities.InvoicingConfig) error {
	// Obtener los integration IDs actuales para invalidar caché luego
	var oldJoinEntries []models.InvoicingConfigIntegration
	r.db.Conn(ctx).Where("config_id = ? AND deleted_at IS NULL", config.ID).Find(&oldJoinEntries)

	var oldIntegrationIDs []uint
	for _, entry := range oldJoinEntries {
		oldIntegrationIDs = append(oldIntegrationIDs, entry.IntegrationID)
	}

	// Actualizar columnas principales del config
	updates := map[string]interface{}{
		"enabled":      config.Enabled,
		"auto_invoice": config.AutoInvoice,
		"description":  config.Description,
	}

	if config.InvoicingIntegrationID != nil && *config.InvoicingIntegrationID > 0 {
		updates["invoicing_integration_id"] = *config.InvoicingIntegrationID
	}

	if config.Filters != nil {
		if data, err := json.Marshal(config.Filters); err == nil {
			updates["filters"] = datatypes.JSON(data)
		}
	}

	if config.InvoiceConfig != nil {
		if data, err := json.Marshal(config.InvoiceConfig); err == nil {
			updates["invoice_config"] = datatypes.JSON(data)
		}
	}

	if err := r.db.Conn(ctx).Model(&models.InvoicingConfig{}).
		Where("id = ? AND deleted_at IS NULL", config.ID).
		Updates(updates).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("config_id", config.ID).Msg("Failed to update config")
		return fmt.Errorf("failed to update config: %w", err)
	}

	// Sincronizar tabla join si se proveen integration IDs
	if len(config.IntegrationIDs) > 0 {
		// Hard delete de entradas anteriores
		if err := r.db.Conn(ctx).Unscoped().
			Where("config_id = ?", config.ID).
			Delete(&models.InvoicingConfigIntegration{}).Error; err != nil {
			r.log.Error(ctx).Err(err).Uint("config_id", config.ID).Msg("Failed to clear config integrations")
			return fmt.Errorf("failed to clear config integrations: %w", err)
		}

		// Insertar nuevas entradas
		for _, integrationID := range config.IntegrationIDs {
			joinEntry := &models.InvoicingConfigIntegration{
				ConfigID:      config.ID,
				IntegrationID: integrationID,
			}
			if err := r.db.Conn(ctx).Create(joinEntry).Error; err != nil {
				r.log.Error(ctx).Err(err).
					Uint("config_id", config.ID).
					Uint("integration_id", integrationID).
					Msg("Failed to create config-integration join entry on update")
				return fmt.Errorf("failed to link integration to config: %w", err)
			}
		}
	}

	r.log.Info(ctx).Uint("config_id", config.ID).Msg("Invoicing config updated")

	// Invalidar caché para IDs anteriores
	for _, integrationID := range oldIntegrationIDs {
		id := integrationID
		go r.configCache.Invalidate(context.Background(), id)
	}

	// Actualizar caché con nuevos IDs
	newIDs := config.IntegrationIDs
	if len(newIDs) == 0 {
		newIDs = oldIntegrationIDs // Si no cambiaron, re-cachear con los mismos
	}
	for _, integrationID := range newIDs {
		id := integrationID
		go r.configCache.Set(context.Background(), id, config)
	}

	return nil
}

// DeleteInvoicingConfig elimina PERMANENTEMENTE una configuración (hard delete).
// ON DELETE CASCADE en invoicing_config_integrations limpia la join table automáticamente.
func (r *Repository) DeleteInvoicingConfig(ctx context.Context, id uint) error {
	var model models.InvoicingConfig
	if err := r.db.Conn(ctx).Unscoped().First(&model, id).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("config_id", id).Msg("Config not found for deletion")
		return fmt.Errorf("config not found: %w", err)
	}

	// Obtener integration IDs antes de eliminar (para invalidar caché)
	var joinEntries []models.InvoicingConfigIntegration
	r.db.Conn(ctx).Unscoped().Where("config_id = ?", id).Find(&joinEntries)

	var integrationIDs []uint
	for _, entry := range joinEntries {
		integrationIDs = append(integrationIDs, entry.IntegrationID)
	}

	// HARD DELETE - ON DELETE CASCADE limpia invoicing_config_integrations
	if err := r.db.Conn(ctx).Unscoped().Delete(&model).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("config_id", id).Msg("Failed to delete config")
		return fmt.Errorf("failed to delete config: %w", err)
	}

	r.log.Info(ctx).Uint("config_id", id).Msg("Invoicing config permanently deleted (hard delete)")

	// Invalidar caché para todos los integration IDs
	for _, integrationID := range integrationIDs {
		integID := integrationID
		go r.configCache.Invalidate(context.Background(), integID)
	}

	return nil
}

// ConfigExistsForIntegration verifica si existe una configuración para un integration ID de e-commerce
// usando la tabla join invoicing_config_integrations.
func (r *Repository) ConfigExistsForIntegration(ctx context.Context, integrationID uint) (bool, error) {
	var count int64

	if err := r.db.Conn(ctx).
		Table("invoicing_config_integrations ici").
		Joins("JOIN invoicing_configs ic ON ic.id = ici.config_id").
		Where("ici.integration_id = ?", integrationID).
		Where("ici.deleted_at IS NULL").
		Where("ic.deleted_at IS NULL").
		Count(&count).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("integration_id", integrationID).Msg("Failed to check config existence")
		return false, fmt.Errorf("failed to check config existence: %w", err)
	}

	return count > 0, nil
}

// GetEnabledConfigByBusiness retorna la primera configuración activa (enabled=true) de un negocio.
// Retorna nil (sin error) si no existe ninguna activa.
func (r *Repository) GetEnabledConfigByBusiness(ctx context.Context, businessID uint) (*entities.InvoicingConfig, error) {
	var model models.InvoicingConfig

	err := r.db.Conn(ctx).
		Where("business_id = ?", businessID).
		Where("enabled = ?", true).
		Where("deleted_at IS NULL").
		Limit(1).
		First(&model).Error

	if err != nil {
		return nil, nil
	}

	return mappers.ConfigToDomain(&model), nil
}

// GetAnyConfigByBusiness retorna la primera configuración de un negocio sin filtrar por enabled.
// Usado para operaciones de auditoría donde solo se necesitan credenciales.
func (r *Repository) GetAnyConfigByBusiness(ctx context.Context, businessID uint) (*entities.InvoicingConfig, error) {
	var model models.InvoicingConfig

	err := r.db.Conn(ctx).
		Preload("InvoicingIntegration").
		Preload("InvoicingIntegration.IntegrationType").
		Where("business_id = ?", businessID).
		Where("deleted_at IS NULL").
		Order("enabled DESC, created_at DESC").
		Limit(1).
		First(&model).Error

	if err != nil {
		return nil, nil
	}

	return mappers.ConfigToDomain(&model), nil
}

// ListAllActiveConfigs lista todas las configuraciones activas de todos los negocios.
// Usado para cache warming al iniciar el servidor.
func (r *Repository) ListAllActiveConfigs(ctx context.Context) ([]*entities.InvoicingConfig, error) {
	var configModels []*models.InvoicingConfig

	if err := r.db.Conn(ctx).
		Preload("ConfigIntegrations").
		Preload("InvoicingIntegration").
		Preload("InvoicingIntegration.IntegrationType").
		Where("enabled = ?", true).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&configModels).Error; err != nil {
		r.log.Error(ctx).Err(err).Msg("Failed to list all active configs")
		return nil, fmt.Errorf("failed to list all active configs: %w", err)
	}

	return mappers.ConfigListToDomain(configModels), nil
}

// isUniqueConstraintError detecta errores de clave duplicada de PostgreSQL
func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "idx_business_invoicing_integration") ||
		strings.Contains(msg, "idx_config_integration") ||
		strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "unique constraint")
}
