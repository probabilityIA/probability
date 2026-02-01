package repository

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// CreateIntegration crea una nueva integración
func (r *Repository) CreateIntegration(ctx context.Context, integration *domain.Integration) error {
	// Encriptar credenciales antes de guardar
	if len(integration.Credentials) > 0 {
		// datatypes.JSON es []byte, convertir a map para encriptar
		var credentialsMap map[string]interface{}
		credentialsBytes := []byte(integration.Credentials)
		if err := json.Unmarshal(credentialsBytes, &credentialsMap); err == nil {
			// VALIDACIÓN: Rechazar credenciales que parezcan ser el wrapper encriptado
			// Esto previene el bug de doble-encriptación
			if _, hasEncrypted := credentialsMap["encrypted"]; hasEncrypted && len(credentialsMap) == 1 {
				r.log.Error(ctx).Msg("Las credenciales parecen ser el wrapper encriptado, no las credenciales reales")
				return fmt.Errorf("credenciales inválidas: no envíe el wrapper encriptado como credenciales")
			}

			encrypted, err := r.encryptionService.EncryptCredentials(ctx, credentialsMap)
			if err != nil {
				r.log.Error(ctx).Err(err).Msg("Error al encriptar credenciales")
				return fmt.Errorf("error al encriptar credenciales: %w", err)
			}
			// Codificar en base64 para guardar en JSONB (que requiere UTF-8)
			encoded := base64.StdEncoding.EncodeToString(encrypted)
			// Crear un JSON con el valor codificado
			encodedJSON, err := json.Marshal(map[string]string{"encrypted": encoded})
			if err != nil {
				r.log.Error(ctx).Err(err).Msg("Error al codificar credenciales en JSON")
				return fmt.Errorf("error al codificar credenciales: %w", err)
			}
			integration.Credentials = encodedJSON
		}
	}

	// Convertir domain.Integration a models.Integration
	model := r.toModel(integration)

	// Asegurar que el ID sea 0 para crear un nuevo registro (GORM generará el ID automáticamente)
	// También resetear timestamps para que se generen automáticamente
	model.ID = 0
	model.CreatedAt = time.Time{}
	model.UpdatedAt = time.Time{}

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

// UpdateIntegration actualiza una integración existente
func (r *Repository) UpdateIntegration(ctx context.Context, id uint, integration *domain.Integration) error {
	// Encriptar credenciales si se están actualizando
	if len(integration.Credentials) > 0 {
		var credentialsMap map[string]interface{}
		credentialsBytes := []byte(integration.Credentials)
		if err := json.Unmarshal(credentialsBytes, &credentialsMap); err == nil {
			// VALIDACIÓN: Rechazar credenciales que parezcan ser el wrapper encriptado
			// Esto previene el bug de doble-encriptación
			if _, hasEncrypted := credentialsMap["encrypted"]; hasEncrypted && len(credentialsMap) == 1 {
				r.log.Error(ctx).Uint("id", id).Msg("Las credenciales parecen ser el wrapper encriptado, no las credenciales reales")
				return fmt.Errorf("credenciales inválidas: no envíe el wrapper encriptado como credenciales")
			}

			encrypted, err := r.encryptionService.EncryptCredentials(ctx, credentialsMap)
			if err != nil {
				r.log.Error(ctx).Err(err).Msg("Error al encriptar credenciales")
				return fmt.Errorf("error al encriptar credenciales: %w", err)
			}
			// Codificar en base64 para guardar en JSONB (que requiere UTF-8)
			encoded := base64.StdEncoding.EncodeToString(encrypted)
			// Crear un JSON con el valor codificado
			encodedJSON, err := json.Marshal(map[string]string{"encrypted": encoded})
			if err != nil {
				r.log.Error(ctx).Err(err).Msg("Error al codificar credenciales en JSON")
				return fmt.Errorf("error al codificar credenciales: %w", err)
			}
			integration.Credentials = encodedJSON
		}
	}

	model := r.toModel(integration)
	model.ID = id

	// Preparar campos para actualizar, excluyendo UpdatedByID si es 0 o nil
	updateFields := map[string]interface{}{
		"name":                model.Name,
		"code":                model.Code,
		"integration_type_id": model.IntegrationTypeID,
		"business_id":         model.BusinessID,
		"store_id":            model.StoreID,
		"is_active":           model.IsActive,
		"is_default":          model.IsDefault,
		"config":              model.Config,
		"credentials":         model.Credentials,
		"description":         model.Description,
		"updated_at":          model.UpdatedAt,
	}

	// Solo incluir updated_by_id si tiene un valor válido (mayor que 0)
	if model.UpdatedByID != nil && *model.UpdatedByID > 0 {
		updateFields["updated_by_id"] = *model.UpdatedByID
	}

	if err := r.db.Conn(ctx).Model(&models.Integration{}).Where("id = ?", id).Updates(updateFields).Error; err != nil {
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

// GetIntegrationByID obtiene una integración por su ID
func (r *Repository) GetIntegrationByID(ctx context.Context, id uint) (*domain.Integration, error) {
	var model models.Integration
	if err := r.db.Conn(ctx).
		Preload("IntegrationType").
		Preload("IntegrationType.Category").
		First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("integración con ID %d no encontrada", id)
		}
		r.log.Error(ctx).Err(err).Uint("id", id).Msg("Error al obtener integración")
		return nil, fmt.Errorf("error al obtener integración: %w", err)
	}

	return r.toDomain(&model), nil
}

// DeleteIntegration elimina una integración
func (r *Repository) DeleteIntegration(ctx context.Context, id uint) error {
	if err := r.db.Conn(ctx).Delete(&models.Integration{}, id).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("id", id).Msg("Error al eliminar integración")
		return fmt.Errorf("error al eliminar integración: %w", err)
	}
	return nil
}

// ListIntegrations lista integraciones con filtros
func (r *Repository) ListIntegrations(ctx context.Context, filters domain.IntegrationFilters) ([]*domain.Integration, int64, error) {
	var integrationModels []models.Integration
	var total int64

	query := r.db.Conn(ctx).Model(&models.Integration{})

	// Aplicar filtros
	if filters.IntegrationTypeID != nil {
		query = query.Where("integration_type_id = ?", *filters.IntegrationTypeID)
	} else if filters.IntegrationTypeCode != nil {
		// Si se filtra por código, necesitamos hacer un JOIN con integration_type
		query = query.Joins("JOIN integration_type ON integration.integration_type_id = integration_type.id").
			Where("integration_type.code = ?", *filters.IntegrationTypeCode)
	}
	if filters.Category != nil {
		// Filtrar por categoría a través de integration_types -> integration_categories
		query = query.Joins("JOIN integration_types it ON integrations.integration_type_id = it.id").
			Joins("JOIN integration_categories ic ON it.category_id = ic.id").
			Where("ic.code = ?", *filters.Category)
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
	if filters.StoreID != nil && *filters.StoreID != "" {
		query = query.Where("store_id = ?", *filters.StoreID)
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

	// Obtener resultados con la relación IntegrationType y su categoría cargadas
	if err := query.
		Preload("IntegrationType").
		Preload("IntegrationType.Category").
		Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&integrationModels).Error; err != nil {
		r.log.Error(ctx).Err(err).Msg("Error al listar integraciones")
		return nil, 0, fmt.Errorf("error al listar integraciones: %w", err)
	}

	// DEBUG: Log para verificar qué se cargó de la BD
	if len(integrationModels) > 0 {
		firstModel := integrationModels[0]
		r.log.Info(ctx).
			Uint("integration_id", firstModel.ID).
			Str("integration_name", firstModel.Name).
			Uint("integration_type_id", firstModel.IntegrationTypeID).
			Bool("integration_type_loaded", firstModel.IntegrationType.ID != 0).
			Interface("category_id_ptr", firstModel.IntegrationType.CategoryID).
			Bool("category_loaded", firstModel.IntegrationType.Category != nil).
			Msg("[DEBUG] ListIntegrations - First result loaded from DB")

		if firstModel.IntegrationType.Category != nil {
			r.log.Info(ctx).
				Uint("category_id", firstModel.IntegrationType.Category.ID).
				Str("category_code", firstModel.IntegrationType.Category.Code).
				Str("category_name", firstModel.IntegrationType.Category.Name).
				Msg("[DEBUG] ListIntegrations - Category data found")
		} else {
			r.log.Warn(ctx).
				Uint("integration_type_id", firstModel.IntegrationType.ID).
				Interface("category_id", firstModel.IntegrationType.CategoryID).
				Msg("[DEBUG] ListIntegrations - Category is NIL despite CategoryID being set")
		}
	}

	// Convertir a dominio
	integrations := make([]*domain.Integration, len(integrationModels))
	for i, model := range integrationModels {
		integrations[i] = r.toDomain(&model)
	}

	return integrations, total, nil
}

// GetIntegrationByIntegrationTypeID obtiene una integración por tipo y business_id
func (r *Repository) GetIntegrationByIntegrationTypeID(ctx context.Context, integrationTypeID uint, businessID *uint) (*domain.Integration, error) {
	var model models.Integration
	query := r.db.Conn(ctx).
		Preload("IntegrationType").
		Preload("IntegrationType.Category").
		Where("integration_type_id = ?", integrationTypeID)

	if businessID != nil {
		query = query.Where("business_id = ?", *businessID)
	} else {
		query = query.Where("business_id IS NULL")
	}

	if err := query.First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("integración de tipo con ID %d no encontrada", integrationTypeID)
		}
		r.log.Error(ctx).Err(err).Uint("integration_type_id", integrationTypeID).Msg("Error al obtener integración por tipo")
		return nil, fmt.Errorf("error al obtener integración por tipo: %w", err)
	}

	return r.toDomain(&model), nil
}

// GetActiveIntegrationByIntegrationTypeID obtiene una integración activa por tipo y business_id
func (r *Repository) GetActiveIntegrationByIntegrationTypeID(ctx context.Context, integrationTypeID uint, businessID *uint) (*domain.Integration, error) {
	var model models.Integration
	query := r.db.Conn(ctx).
		Preload("IntegrationType").
		Preload("IntegrationType.Category").
		Where("integration_type_id = ? AND is_active = ?", integrationTypeID, true)

	if businessID != nil {
		query = query.Where("business_id = ?", *businessID)
	} else {
		query = query.Where("business_id IS NULL")
	}

	if err := query.First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("integración activa de tipo con ID %d no encontrada", integrationTypeID)
		}
		r.log.Error(ctx).Err(err).Uint("integration_type_id", integrationTypeID).Msg("Error al obtener integración activa por tipo")
		return nil, fmt.Errorf("error al obtener integración activa por tipo: %w", err)
	}

	return r.toDomain(&model), nil
}

// ListIntegrationsByBusiness lista integraciones de un business
func (r *Repository) ListIntegrationsByBusiness(ctx context.Context, businessID uint) ([]*domain.Integration, error) {
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

// ListIntegrationsByIntegrationTypeID lista integraciones por tipo de integración
func (r *Repository) ListIntegrationsByIntegrationTypeID(ctx context.Context, integrationTypeID uint) ([]*domain.Integration, error) {
	var integrationModels []models.Integration
	if err := r.db.Conn(ctx).Where("integration_type_id = ?", integrationTypeID).Find(&integrationModels).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("integration_type_id", integrationTypeID).Msg("Error al listar integraciones por tipo")
		return nil, fmt.Errorf("error al listar integraciones por tipo: %w", err)
	}

	integrations := make([]*domain.Integration, len(integrationModels))
	for i, model := range integrationModels {
		integrations[i] = r.toDomain(&model)
	}

	return integrations, nil
}

// SetIntegrationAsDefault marca una integración como default
func (r *Repository) SetIntegrationAsDefault(ctx context.Context, id uint) error {
	// Primero obtener la integración para saber su tipo y business_id
	var integration models.Integration
	if err := r.db.Conn(ctx).First(&integration, id).Error; err != nil {
		return fmt.Errorf("integración no encontrada: %w", err)
	}

	// Desmarcar todas las demás del mismo tipo y business como no default
	query := r.db.Conn(ctx).Model(&models.Integration{}).
		Where("integration_type_id = ? AND id != ?", integration.IntegrationTypeID, id)

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

// ExistsIntegrationByCode verifica si existe una integración con el código dado
func (r *Repository) ExistsIntegrationByCode(ctx context.Context, code string, businessID *uint) (bool, error) {
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
func (r *Repository) toModel(integration *domain.Integration) *models.Integration {
	model := &models.Integration{
		Model: gorm.Model{
			ID:        integration.ID,
			CreatedAt: integration.CreatedAt,
			UpdatedAt: integration.UpdatedAt,
		},
		Name:              integration.Name,
		Code:              integration.Code,
		IntegrationTypeID: integration.IntegrationTypeID,
		BusinessID:        integration.BusinessID,
		StoreID:           integration.StoreID,
		IsActive:          integration.IsActive,
		IsDefault:         integration.IsDefault,
		Config:            integration.Config,
		Credentials:       integration.Credentials,
		Description:       integration.Description,
		CreatedByID:       integration.CreatedByID,
	}
	if integration.UpdatedByID != nil {
		model.UpdatedByID = integration.UpdatedByID
	}
	return model
}

// toDomain convierte models.Integration a domain.Integration
func (r *Repository) toDomain(model *models.Integration) *domain.Integration {
	businessID := model.BusinessID
	var updatedByID *uint
	if model.UpdatedByID != nil {
		updatedByID = model.UpdatedByID
	}

	integration := &domain.Integration{
		ID:                model.ID,
		Name:              model.Name,
		Code:              model.Code,
		IntegrationTypeID: model.IntegrationTypeID,
		BusinessID:        businessID,
		StoreID:           model.StoreID,
		IsActive:          model.IsActive,
		IsDefault:         model.IsDefault,
		Config:            model.Config,
		Credentials:       model.Credentials, // Mantener encriptado
		Description:       model.Description,
		CreatedByID:       model.CreatedByID,
		UpdatedByID:       updatedByID,
		CreatedAt:         model.CreatedAt,
		UpdatedAt:         model.UpdatedAt,
	}

	// Cargar IntegrationType si está disponible en el modelo
	if model.IntegrationType != nil && model.IntegrationType.ID != 0 {
		// DEBUG: Log para verificar si Category se está cargando
		r.log.Info(context.Background()).
			Uint("integration_type_id", model.IntegrationType.ID).
			Str("integration_type_code", model.IntegrationType.Code).
			Interface("category_id", model.IntegrationType.CategoryID).
			Bool("category_is_nil", model.IntegrationType.Category == nil).
			Msg("[DEBUG] toDomain - IntegrationType loaded")

		var category *domain.IntegrationCategory
		if model.IntegrationType.Category != nil {
			r.log.Info(context.Background()).
				Uint("category_id", model.IntegrationType.Category.ID).
				Str("category_code", model.IntegrationType.Category.Code).
				Str("category_name", model.IntegrationType.Category.Name).
				Msg("[DEBUG] toDomain - Category found and mapping")

			category = &domain.IntegrationCategory{
				ID:               model.IntegrationType.Category.ID,
				Code:             model.IntegrationType.Category.Code,
				Name:             model.IntegrationType.Category.Name,
				Description:      model.IntegrationType.Category.Description,
				Icon:             model.IntegrationType.Category.Icon,
				Color:            model.IntegrationType.Category.Color,
				DisplayOrder:     model.IntegrationType.Category.DisplayOrder,
				ParentCategoryID: model.IntegrationType.Category.ParentCategoryID,
				IsActive:         model.IntegrationType.Category.IsActive,
				IsVisible:        model.IntegrationType.Category.IsVisible,
				CreatedAt:        model.IntegrationType.Category.CreatedAt,
				UpdatedAt:        model.IntegrationType.Category.UpdatedAt,
			}
		}

		categoryID := uint(0)
		if model.IntegrationType.CategoryID != nil {
			categoryID = *model.IntegrationType.CategoryID
		}

		integrationType := domain.IntegrationType{
			ID:                model.IntegrationType.ID,
			Name:              model.IntegrationType.Name,
			Code:              model.IntegrationType.Code,
			Description:       model.IntegrationType.Description,
			Icon:              model.IntegrationType.Icon,
			ImageURL:          model.IntegrationType.ImageURL,
			CategoryID:        categoryID,
			Category:          category,
			IsActive:          model.IntegrationType.IsActive,
			ConfigSchema:      model.IntegrationType.ConfigSchema,
			CredentialsSchema: model.IntegrationType.CredentialsSchema,
			SetupInstructions: model.IntegrationType.SetupInstructions,
			CreatedAt:         model.IntegrationType.CreatedAt,
			UpdatedAt:         model.IntegrationType.UpdatedAt,
		}
		integration.IntegrationType = &integrationType
	}

	return integration
}
