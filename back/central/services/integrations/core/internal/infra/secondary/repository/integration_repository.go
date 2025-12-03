package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

type integrationRepository struct {
	db                db.IDatabase
	log               log.ILogger
	encryptionService domain.IEncryptionService
}

// Create crea una nueva integración
func (r *integrationRepository) Create(ctx context.Context, integration *domain.Integration) error {
	// Encriptar credenciales antes de guardar
	if len(integration.Credentials) > 0 {
		// datatypes.JSON es []byte, convertir a map para encriptar
		var credentialsMap map[string]interface{}
		credentialsBytes := []byte(integration.Credentials)
		if err := json.Unmarshal(credentialsBytes, &credentialsMap); err == nil {
			encrypted, err := r.encryptionService.EncryptCredentials(ctx, credentialsMap)
			if err != nil {
				r.log.Error(ctx).Err(err).Msg("Error al encriptar credenciales")
				return fmt.Errorf("error al encriptar credenciales: %w", err)
			}
			integration.Credentials = encrypted
		}
	}

	// Convertir domain.Integration a models.Integration
	model := r.toModel(integration)

	if err := r.db.Conn(ctx).Create(&model).Error; err != nil {
		r.log.Error(ctx).Err(err).Msg("Error al crear integración")
		return fmt.Errorf("error al crear integración: %w", err)
	}

	// Actualizar el ID en el dominio
	integration.ID = model.ID
	integration.CreatedAt = model.CreatedAt
	integration.UpdatedAt = model.UpdatedAt

	return nil
}

// Update actualiza una integración existente
func (r *integrationRepository) Update(ctx context.Context, id uint, integration *domain.Integration) error {
	// Encriptar credenciales si se están actualizando
	if len(integration.Credentials) > 0 {
		var credentialsMap map[string]interface{}
		credentialsBytes := []byte(integration.Credentials)
		if err := json.Unmarshal(credentialsBytes, &credentialsMap); err == nil {
			encrypted, err := r.encryptionService.EncryptCredentials(ctx, credentialsMap)
			if err != nil {
				r.log.Error(ctx).Err(err).Msg("Error al encriptar credenciales")
				return fmt.Errorf("error al encriptar credenciales: %w", err)
			}
			integration.Credentials = encrypted
		}
	}

	model := r.toModel(integration)
	model.ID = id

	if err := r.db.Conn(ctx).Model(&models.Integration{}).Where("id = ?", id).Updates(model).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("id", id).Msg("Error al actualizar integración")
		return fmt.Errorf("error al actualizar integración: %w", err)
	}

	// Obtener la integración actualizada para obtener timestamps
	var updated models.Integration
	if err := r.db.Conn(ctx).First(&updated, id).Error; err != nil {
		return fmt.Errorf("error al obtener integración actualizada: %w", err)
	}

	integration.UpdatedAt = updated.UpdatedAt
	return nil
}

// GetByID obtiene una integración por su ID
func (r *integrationRepository) GetByID(ctx context.Context, id uint) (*domain.Integration, error) {
	var model models.Integration
	if err := r.db.Conn(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("integración con ID %d no encontrada", id)
		}
		r.log.Error(ctx).Err(err).Uint("id", id).Msg("Error al obtener integración")
		return nil, fmt.Errorf("error al obtener integración: %w", err)
	}

	return r.toDomain(&model), nil
}

// Delete elimina una integración
func (r *integrationRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.Conn(ctx).Delete(&models.Integration{}, id).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("id", id).Msg("Error al eliminar integración")
		return fmt.Errorf("error al eliminar integración: %w", err)
	}
	return nil
}

// List lista integraciones con filtros
func (r *integrationRepository) List(ctx context.Context, filters domain.IntegrationFilters) ([]*domain.Integration, int64, error) {
	var integrationModels []models.Integration
	var total int64

	query := r.db.Conn(ctx).Model(&models.Integration{})

	// Aplicar filtros
	if filters.Type != nil {
		query = query.Where("type = ?", *filters.Type)
	}
	if filters.Category != nil {
		query = query.Where("category = ?", *filters.Category)
	}
	if filters.BusinessID != nil {
		query = query.Where("business_id = ?", *filters.BusinessID)
	}
	// Si BusinessID es nil, NO filtrar (mostrar todas: globales y por business)
	if filters.IsActive != nil {
		query = query.Where("is_active = ?", *filters.IsActive)
	}
	if filters.Search != nil && *filters.Search != "" {
		search := "%" + *filters.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ?", search, search)
	}

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		r.log.Error(ctx).Err(err).Msg("Error al contar integraciones")
		return nil, 0, fmt.Errorf("error al contar integraciones: %w", err)
	}

	// Aplicar paginación
	page := filters.Page
	if page < 1 {
		page = 1
	}
	pageSize := filters.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	// Obtener resultados
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&integrationModels).Error; err != nil {
		r.log.Error(ctx).Err(err).Msg("Error al listar integraciones")
		return nil, 0, fmt.Errorf("error al listar integraciones: %w", err)
	}

	// Convertir a dominio
	integrations := make([]*domain.Integration, len(integrationModels))
	for i, model := range integrationModels {
		integrations[i] = r.toDomain(&model)
	}

	return integrations, total, nil
}

// GetByType obtiene una integración por tipo y business_id
func (r *integrationRepository) GetByType(ctx context.Context, integrationType string, businessID *uint) (*domain.Integration, error) {
	var model models.Integration
	query := r.db.Conn(ctx).Where("type = ?", integrationType)

	if businessID != nil {
		query = query.Where("business_id = ?", *businessID)
	} else {
		query = query.Where("business_id IS NULL")
	}

	if err := query.First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("integración de tipo '%s' no encontrada", integrationType)
		}
		r.log.Error(ctx).Err(err).Str("type", integrationType).Msg("Error al obtener integración por tipo")
		return nil, fmt.Errorf("error al obtener integración por tipo: %w", err)
	}

	return r.toDomain(&model), nil
}

// GetActiveByType obtiene una integración activa por tipo y business_id
func (r *integrationRepository) GetActiveByType(ctx context.Context, integrationType string, businessID *uint) (*domain.Integration, error) {
	var model models.Integration
	query := r.db.Conn(ctx).Where("type = ? AND is_active = ?", integrationType, true)

	if businessID != nil {
		query = query.Where("business_id = ?", *businessID)
	} else {
		query = query.Where("business_id IS NULL")
	}

	if err := query.First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("integración activa de tipo '%s' no encontrada", integrationType)
		}
		r.log.Error(ctx).Err(err).Str("type", integrationType).Msg("Error al obtener integración activa por tipo")
		return nil, fmt.Errorf("error al obtener integración activa por tipo: %w", err)
	}

	return r.toDomain(&model), nil
}

// ListByBusiness lista integraciones de un business
func (r *integrationRepository) ListByBusiness(ctx context.Context, businessID uint) ([]*domain.Integration, error) {
	var integrationModels []models.Integration
	if err := r.db.Conn(ctx).Where("business_id = ?", businessID).Find(&integrationModels).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("business_id", businessID).Msg("Error al listar integraciones por business")
		return nil, fmt.Errorf("error al listar integraciones por business: %w", err)
	}

	integrations := make([]*domain.Integration, len(integrationModels))
	for i, model := range integrationModels {
		integrations[i] = r.toDomain(&model)
	}

	return integrations, nil
}

// ListByType lista integraciones por tipo
func (r *integrationRepository) ListByType(ctx context.Context, integrationType string) ([]*domain.Integration, error) {
	var integrationModels []models.Integration
	if err := r.db.Conn(ctx).Where("type = ?", integrationType).Find(&integrationModels).Error; err != nil {
		r.log.Error(ctx).Err(err).Str("type", integrationType).Msg("Error al listar integraciones por tipo")
		return nil, fmt.Errorf("error al listar integraciones por tipo: %w", err)
	}

	integrations := make([]*domain.Integration, len(integrationModels))
	for i, model := range integrationModels {
		integrations[i] = r.toDomain(&model)
	}

	return integrations, nil
}

// SetAsDefault marca una integración como default
func (r *integrationRepository) SetAsDefault(ctx context.Context, id uint) error {
	// Primero obtener la integración para saber su tipo y business_id
	var integration models.Integration
	if err := r.db.Conn(ctx).First(&integration, id).Error; err != nil {
		return fmt.Errorf("integración no encontrada: %w", err)
	}

	// Desmarcar todas las demás del mismo tipo y business como no default
	query := r.db.Conn(ctx).Model(&models.Integration{}).
		Where("type = ? AND id != ?", integration.Type, id)

	if integration.BusinessID != nil {
		query = query.Where("business_id = ?", *integration.BusinessID)
	} else {
		query = query.Where("business_id IS NULL")
	}

	if err := query.Update("is_default", false).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("id", id).Msg("Error al desmarcar otras integraciones como default")
		return fmt.Errorf("error al desmarcar otras integraciones: %w", err)
	}

	// Marcar esta como default
	if err := r.db.Conn(ctx).Model(&models.Integration{}).Where("id = ?", id).Update("is_default", true).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("id", id).Msg("Error al marcar integración como default")
		return fmt.Errorf("error al marcar integración como default: %w", err)
	}

	return nil
}

// ExistsByCode verifica si existe una integración con el código dado
func (r *integrationRepository) ExistsByCode(ctx context.Context, code string, businessID *uint) (bool, error) {
	var count int64
	query := r.db.Conn(ctx).Model(&models.Integration{}).Where("code = ?", code)

	if businessID != nil {
		query = query.Where("business_id = ?", *businessID)
	} else {
		query = query.Where("business_id IS NULL")
	}

	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("error al verificar existencia de código: %w", err)
	}

	return count > 0, nil
}

// toModel convierte domain.Integration a models.Integration
func (r *integrationRepository) toModel(integration *domain.Integration) *models.Integration {
	model := &models.Integration{
		Model: gorm.Model{
			ID:        integration.ID,
			CreatedAt: integration.CreatedAt,
			UpdatedAt: integration.UpdatedAt,
		},
		Name:        integration.Name,
		Code:        integration.Code,
		Type:        integration.Type,
		Category:    integration.Category,
		BusinessID:  integration.BusinessID,
		IsActive:    integration.IsActive,
		IsDefault:   integration.IsDefault,
		Config:      integration.Config,
		Credentials: integration.Credentials,
		Description: integration.Description,
		CreatedByID: integration.CreatedByID,
	}
	if integration.UpdatedByID != nil {
		model.UpdatedByID = integration.UpdatedByID
	}
	return model
}

// toDomain convierte models.Integration a domain.Integration
func (r *integrationRepository) toDomain(model *models.Integration) *domain.Integration {
	businessID := model.BusinessID
	var updatedByID *uint
	if model.UpdatedByID != nil {
		updatedByID = model.UpdatedByID
	}

	return &domain.Integration{
		ID:          model.ID,
		Name:        model.Name,
		Code:        model.Code,
		Type:        model.Type,
		Category:    model.Category,
		BusinessID:  businessID,
		IsActive:    model.IsActive,
		IsDefault:   model.IsDefault,
		Config:      model.Config,
		Credentials: model.Credentials, // Mantener encriptado
		Description: model.Description,
		CreatedByID: model.CreatedByID,
		UpdatedByID: updatedByID,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}
