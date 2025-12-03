package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/auth/resources/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"

	"gorm.io/gorm"
)

// GetResources obtiene todos los recursos con filtros y paginación
func (r *Repository) GetResources(ctx context.Context, filters domain.ResourceFilters) ([]domain.Resource, int64, error) {
	r.logger.Info().Interface("filters", filters).Msg("Iniciando búsqueda de recursos")

	// Configurar paginación por defecto
	if filters.PageSize <= 0 {
		filters.PageSize = 10
	}
	if filters.Page <= 0 {
		filters.Page = 1
	}

	offset := (filters.Page - 1) * filters.PageSize

	// Construir query base
	query := r.database.Conn(ctx).Model(&models.Resource{})

	// Aplicar filtros
	if filters.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filters.Name+"%")
	}
	if filters.Description != "" {
		query = query.Where("description ILIKE ?", "%"+filters.Description+"%")
	}

	// Filtrar por business_type_id
	// Si es nil, solo muestra recursos genéricos (business_type_id IS NULL)
	// Si tiene valor, muestra recursos de ese tipo o genéricos
	if filters.BusinessTypeID != nil {
		query = query.Where("business_type_id = ? OR business_type_id IS NULL", *filters.BusinessTypeID)
	} else {
		// Si no se especifica business_type_id, mostrar todos los recursos
		// (tanto genéricos como específicos por tipo)
		// No aplicar ningún filtro adicional
	}

	// Contar total antes de aplicar paginación
	var total int64
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al contar recursos")
		return nil, 0, err
	}

	// Aplicar ordenamiento
	orderBy := "created_at DESC" // Por defecto
	if filters.SortBy != "" {
		direction := "ASC"
		if strings.ToUpper(filters.SortOrder) == "DESC" {
			direction = "DESC"
		}

		// Validar campos de ordenamiento permitidos
		allowedSortFields := map[string]bool{
			"name":       true,
			"created_at": true,
			"updated_at": true,
		}

		if allowedSortFields[filters.SortBy] {
			orderBy = fmt.Sprintf("%s %s", filters.SortBy, direction)
		}
	}

	// Aplicar paginación y ordenamiento con preload de BusinessType
	var resources []models.Resource
	if err := query.Preload("BusinessType").Order(orderBy).Offset(offset).Limit(filters.PageSize).Find(&resources).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al obtener recursos")
		return nil, 0, err
	}

	// Convertir a entidades de dominio
	var domainResources []domain.Resource
	for _, resource := range resources {
		var deletedAt *time.Time
		if resource.DeletedAt.Valid {
			deletedAt = &resource.DeletedAt.Time
		}

		businessTypeID := uint(0)
		businessTypeName := ""
		if resource.BusinessTypeID != nil {
			businessTypeID = *resource.BusinessTypeID
			if resource.BusinessType != nil {
				businessTypeName = resource.BusinessType.Name
			}
		}

		domainResources = append(domainResources, domain.Resource{
			ID:               resource.ID,
			Name:             resource.Name,
			Description:      resource.Description,
			BusinessTypeID:   businessTypeID,
			BusinessTypeName: businessTypeName,
			CreatedAt:        resource.CreatedAt,
			UpdatedAt:        resource.UpdatedAt,
			DeletedAt:        deletedAt,
		})
	}

	r.logger.Info().
		Int64("total", total).
		Int("returned", len(domainResources)).
		Int("page", filters.Page).
		Int("page_size", filters.PageSize).
		Msg("Recursos obtenidos exitosamente")

	return domainResources, total, nil
}

// GetResourceByID obtiene un recurso por su ID
func (r *Repository) GetResourceByID(ctx context.Context, id uint) (*domain.Resource, error) {
	r.logger.Info().Uint("resource_id", id).Msg("Obteniendo recurso por ID")

	var resource models.Resource
	if err := r.database.Conn(ctx).Where("id = ?", id).First(&resource).Error; err != nil {
		r.logger.Error().Err(err).Uint("resource_id", id).Msg("Error al obtener recurso por ID")
		return nil, err
	}

	var deletedAt *time.Time
	if resource.DeletedAt.Valid {
		deletedAt = &resource.DeletedAt.Time
	}

	businessTypeID := uint(0)
	businessTypeName := ""
	if resource.BusinessTypeID != nil {
		businessTypeID = *resource.BusinessTypeID
		if resource.BusinessType != nil {
			businessTypeName = resource.BusinessType.Name
		}
	}

	domainResource := &domain.Resource{
		ID:               resource.ID,
		Name:             resource.Name,
		Description:      resource.Description,
		BusinessTypeID:   businessTypeID,
		BusinessTypeName: businessTypeName,
		CreatedAt:        resource.CreatedAt,
		UpdatedAt:        resource.UpdatedAt,
		DeletedAt:        deletedAt,
	}

	r.logger.Info().Uint("resource_id", id).Str("name", resource.Name).Msg("Recurso obtenido exitosamente")
	return domainResource, nil
}

// GetResourceByName obtiene un recurso por su nombre
func (r *Repository) GetResourceByName(ctx context.Context, name string) (*domain.Resource, error) {
	r.logger.Info().Str("name", name).Msg("Obteniendo recurso por nombre")

	var resource models.Resource
	if err := r.database.Conn(ctx).Where("name = ?", name).First(&resource).Error; err != nil {
		r.logger.Error().Err(err).Str("name", name).Msg("Error al obtener recurso por nombre")
		return nil, err
	}

	var deletedAt *time.Time
	if resource.DeletedAt.Valid {
		deletedAt = &resource.DeletedAt.Time
	}

	businessTypeID := uint(0)
	businessTypeName := ""
	if resource.BusinessTypeID != nil {
		businessTypeID = *resource.BusinessTypeID
		if resource.BusinessType != nil {
			businessTypeName = resource.BusinessType.Name
		}
	}

	domainResource := &domain.Resource{
		ID:               resource.ID,
		Name:             resource.Name,
		Description:      resource.Description,
		BusinessTypeID:   businessTypeID,
		BusinessTypeName: businessTypeName,
		CreatedAt:        resource.CreatedAt,
		UpdatedAt:        resource.UpdatedAt,
		DeletedAt:        deletedAt,
	}

	r.logger.Info().Uint("resource_id", resource.ID).Str("name", name).Msg("Recurso obtenido exitosamente por nombre")
	return domainResource, nil
}

// CreateResource crea un nuevo recurso
func (r *Repository) CreateResource(ctx context.Context, resource domain.Resource) (uint, error) {
	r.logger.Info().Str("name", resource.Name).Msg("Creando nuevo recurso")

	modelResource := models.Resource{
		Name:        resource.Name,
		Description: resource.Description,
	}

	// Agregar business_type_id si está presente
	if resource.BusinessTypeID > 0 {
		btID := resource.BusinessTypeID
		modelResource.BusinessTypeID = &btID
	}

	// Iniciar transacción para crear recurso y sus relaciones
	tx := r.database.Conn(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&modelResource).Error; err != nil {
		tx.Rollback()
		r.logger.Error().Err(err).Str("name", resource.Name).Msg("Error al crear recurso")
		return 0, err
	}

	// Si el recurso tiene business_type_id, crear relaciones con todos los businesses de ese tipo
	if resource.BusinessTypeID > 0 {
		if err := r.createBusinessRelationsForResource(tx, modelResource.ID, resource.BusinessTypeID); err != nil {
			tx.Rollback()
			r.logger.Error().Err(err).Uint("resource_id", modelResource.ID).Uint("business_type_id", resource.BusinessTypeID).Msg("Error al crear relaciones con businesses")
			return 0, err
		}
	}

	// Confirmar transacción
	if err := tx.Commit().Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al confirmar transacción")
		return 0, err
	}

	r.logger.Info().
		Uint("resource_id", modelResource.ID).
		Str("name", resource.Name).
		Msg("Recurso creado exitosamente")

	return modelResource.ID, nil
}

// createBusinessRelationsForResource crea relaciones entre el recurso y todos los businesses del tipo especificado
func (r *Repository) createBusinessRelationsForResource(tx *gorm.DB, resourceID uint, businessTypeID uint) error {
	// Obtener todos los businesses del tipo especificado
	var businesses []models.Business
	if err := tx.Where("business_type_id = ?", businessTypeID).Find(&businesses).Error; err != nil {
		r.logger.Error().Err(err).Uint("business_type_id", businessTypeID).Msg("Error al obtener businesses del tipo")
		return err
	}

	if len(businesses) == 0 {
		r.logger.Info().Uint("business_type_id", businessTypeID).Msg("No hay businesses de este tipo, no se crean relaciones")
		return nil
	}

	// Crear relaciones en BusinessResourceConfigured con Active = false
	for _, business := range businesses {
		businessResource := models.BusinessResourceConfigured{
			BusinessID: business.ID,
			ResourceID: resourceID,
			Active:     false, // Nuevo recurso inactivo por defecto
		}

		if err := tx.Create(&businessResource).Error; err != nil {
			r.logger.Error().Err(err).
				Uint("business_id", business.ID).
				Uint("resource_id", resourceID).
				Msg("Error al crear relación business-resource")
			return err
		}
	}

	r.logger.Info().
		Uint("resource_id", resourceID).
		Uint("business_type_id", businessTypeID).
		Int("businesses_count", len(businesses)).
		Msg("Relaciones con businesses creadas exitosamente")

	return nil
}

// UpdateResource actualiza un recurso existente
func (r *Repository) UpdateResource(ctx context.Context, id uint, resource domain.Resource) (string, error) {
	r.logger.Info().Uint("resource_id", id).Str("name", resource.Name).Msg("Actualizando recurso")

	// Obtener el recurso actual para comparar business_type_id
	var currentResource models.Resource
	if err := r.database.Conn(ctx).Where("id = ?", id).First(&currentResource).Error; err != nil {
		r.logger.Error().Err(err).Uint("resource_id", id).Msg("Error al obtener recurso actual")
		return "", fmt.Errorf("recurso con ID %d no encontrado", id)
	}

	// Iniciar transacción
	tx := r.database.Conn(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	updateData := map[string]interface{}{
		"name":        resource.Name,
		"description": resource.Description,
	}

	var oldBusinessTypeID *uint
	var newBusinessTypeID *uint

	// Extraer business_type_id actual y nuevo
	if currentResource.BusinessTypeID != nil {
		oldBusinessTypeID = currentResource.BusinessTypeID
	}

	if resource.BusinessTypeID > 0 {
		btID := resource.BusinessTypeID
		updateData["business_type_id"] = &btID
		newBusinessTypeID = &btID
	}

	// Actualizar el recurso
	result := tx.Model(&models.Resource{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		tx.Rollback()
		r.logger.Error().Err(result.Error).Uint("resource_id", id).Msg("Error al actualizar recurso")
		return "", result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		r.logger.Warn().Uint("resource_id", id).Msg("Recurso no encontrado para actualizar")
		return "", fmt.Errorf("recurso con ID %d no encontrado", id)
	}

	// Verificar si cambió el business_type_id
	businessTypeChanged := false
	if oldBusinessTypeID != nil && newBusinessTypeID != nil {
		businessTypeChanged = *oldBusinessTypeID != *newBusinessTypeID
	} else if oldBusinessTypeID == nil && newBusinessTypeID != nil {
		businessTypeChanged = true // de NULL a un valor
	} else if oldBusinessTypeID != nil && newBusinessTypeID == nil {
		businessTypeChanged = true // de un valor a NULL
	}

	// Si cambió el business_type_id, actualizar las relaciones
	if businessTypeChanged {
		r.logger.Info().
			Uint("resource_id", id).
			Interface("old_business_type_id", oldBusinessTypeID).
			Interface("new_business_type_id", newBusinessTypeID).
			Msg("Cambió business_type_id, actualizando relaciones")

		// Eliminar todas las relaciones existentes
		if err := tx.Unscoped().Where("resource_id = ?", id).Delete(&models.BusinessResourceConfigured{}).Error; err != nil {
			tx.Rollback()
			r.logger.Error().Err(err).Uint("resource_id", id).Msg("Error al eliminar relaciones antiguas")
			return "", err
		}

		// Crear nuevas relaciones si el nuevo business_type_id no es NULL
		if newBusinessTypeID != nil {
			if err := r.createBusinessRelationsForResource(tx, id, *newBusinessTypeID); err != nil {
				tx.Rollback()
				r.logger.Error().Err(err).Uint("resource_id", id).Uint("business_type_id", *newBusinessTypeID).Msg("Error al crear relaciones nuevas")
				return "", err
			}
		} else {
			// Si se cambia a NULL (recurso genérico), no se crean relaciones
			r.logger.Info().Uint("resource_id", id).Msg("Recurso cambiado a genérico (NULL), no se crean relaciones")
		}
	}

	// Confirmar transacción
	if err := tx.Commit().Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al confirmar transacción")
		return "", err
	}

	message := fmt.Sprintf("Recurso actualizado con ID: %d", id)
	r.logger.Info().Uint("resource_id", id).Msg("Recurso actualizado exitosamente")

	return message, nil
}

// DeleteResource elimina un recurso permanentemente con eliminación en cascada
func (r *Repository) DeleteResource(ctx context.Context, id uint) (string, error) {
	r.logger.Info().Uint("resource_id", id).Msg("Eliminando recurso permanentemente")

	// Usar Unscoped().Delete() para eliminación física (no soft delete)
	// Esto activará la eliminación en cascada de las relaciones definidas en el modelo
	result := r.database.Conn(ctx).Unscoped().Delete(&models.Resource{}, id)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Uint("resource_id", id).Msg("Error al eliminar recurso")
		return "", result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.Warn().Uint("resource_id", id).Msg("Recurso no encontrado para eliminar")
		return "", fmt.Errorf("recurso con ID %d no encontrado", id)
	}

	message := fmt.Sprintf("Recurso eliminado permanentemente con ID: %d", id)
	r.logger.Info().
		Uint("resource_id", id).
		Int64("rows_affected", result.RowsAffected).
		Msg("Recurso eliminado exitosamente (eliminación en cascada aplicada)")

	return message, nil
}
