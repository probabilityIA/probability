package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// GetPermissions obtiene todos los permisos con filtros opcionales
func (r *Repository) GetPermissions(ctx context.Context, businessTypeID *uint, name *string, scopeID *uint, resource *string) ([]domain.Permission, error) {
	var permissions []models.Permission
	query := r.database.Conn(ctx).
		Model(&models.Permission{}).
		Preload("Scope").
		Preload("Resource").
		Preload("Action").
		Preload("BusinessType")

	// Filtrar por business_type_id si se proporciona
	// Incluye permisos genéricos (NULL) o del tipo especificado
	if businessTypeID != nil {
		query = query.Where("business_type_id = ? OR business_type_id IS NULL", *businessTypeID)
	}

	// JOIN a resource solo si alguno de los filtros lo requiere
	needsResourceJoin := (name != nil && *name != "") || (resource != nil && *resource != "")
	if needsResourceJoin {
		query = query.Joins("JOIN resource ON permission.resource_id = resource.id")
	}

	// Filtrar por name (búsqueda parcial en permission.name)
	if name != nil && *name != "" {
		query = query.Where("permission.name ILIKE ?", "%"+*name+"%")
	}

	// Filtrar por resource (búsqueda parcial en resource.name)
	if resource != nil && *resource != "" {
		query = query.Where("resource.name ILIKE ?", "%"+*resource+"%")
	}

	// Filtrar por scope_id
	if scopeID != nil {
		query = query.Where("scope_id = ?", *scopeID)
	}

	if err := query.Find(&permissions).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al obtener permisos")
		return nil, err
	}
	return mappers.ToPermissionEntitySlice(permissions), nil
}

// GetPermissionByID obtiene un permiso por su ID
func (r *Repository) GetPermissionByID(ctx context.Context, id uint) (*domain.Permission, error) {
	var permission models.Permission
	if err := r.database.Conn(ctx).
		Model(&models.Permission{}).
		Preload("Scope").
		Preload("Resource").
		Preload("Action").
		Where("id = ?", id).
		First(&permission).Error; err != nil {
		r.logger.Error().Uint("id", id).Err(err).Msg("Error al obtener permiso por ID")
		return nil, err
	}

	entity := mappers.ToPermissionEntity(permission)
	return &entity, nil
}

// GetPermissionsByScopeID obtiene permisos por scope ID
func (r *Repository) GetPermissionsByScopeID(ctx context.Context, scopeID uint) ([]domain.Permission, error) {
	var permissions []models.Permission
	if err := r.database.Conn(ctx).
		Model(&models.Permission{}).
		Preload("Scope").
		Preload("Resource").
		Preload("Action").
		Where("scope_id = ?", scopeID).
		Find(&permissions).Error; err != nil {
		r.logger.Error().Uint("scope_id", scopeID).Err(err).Msg("Error al obtener permisos por scope ID")
		return nil, err
	}
	return mappers.ToPermissionEntitySlice(permissions), nil
}

// GetPermissionsByResource obtiene permisos por recurso
func (r *Repository) GetPermissionsByResource(ctx context.Context, resource string) ([]domain.Permission, error) {
	var permissions []models.Permission
	if err := r.database.Conn(ctx).
		Model(&models.Permission{}).
		Preload("Scope").
		Preload("Resource").
		Preload("Action").
		Joins("JOIN resource ON permission.resource_id = resource.id").
		Where("resource.name = ?", resource).
		Find(&permissions).Error; err != nil {
		r.logger.Error().Str("resource", resource).Err(err).Msg("Error al obtener permisos por recurso")
		return nil, err
	}
	return mappers.ToPermissionEntitySlice(permissions), nil
}

// PermissionExistsByName verifica si existe un permiso con el nombre especificado
func (r *Repository) PermissionExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := r.database.Conn(ctx).
		Model(&models.Permission{}).
		Where("name = ?", name).
		Count(&count).Error; err != nil {
		r.logger.Error().Err(err).Str("name", name).Msg("Error verificando existencia de permiso por nombre")
		return false, fmt.Errorf("error verificando existencia de permiso por nombre: %w", err)
	}
	return count > 0, nil
}

// CreatePermission crea un nuevo permiso
func (r *Repository) CreatePermission(ctx context.Context, permission domain.Permission) (string, error) {
	// Verificar que el Resource existe
	var resource models.Resource
	if err := r.database.Conn(ctx).Where("id = ?", permission.ResourceID).First(&resource).Error; err != nil {
		r.logger.Error().Uint("resource_id", permission.ResourceID).Err(err).Msg("Error al buscar resource")
		return "", fmt.Errorf("resource no encontrado con ID: %d", permission.ResourceID)
	}

	// Verificar que el Action existe
	var action models.Action
	if err := r.database.Conn(ctx).Where("id = ?", permission.ActionID).First(&action).Error; err != nil {
		r.logger.Error().Uint("action_id", permission.ActionID).Err(err).Msg("Error al buscar action")
		return "", fmt.Errorf("action no encontrada con ID: %d", permission.ActionID)
	}

	// Crear el modelo Permission
	permissionModel := models.Permission{
		Name:        permission.Name,
		Description: permission.Description,
		ResourceID:  permission.ResourceID,
		ActionID:    permission.ActionID,
		ScopeID:     permission.ScopeID,
	}

	// Agregar business_type_id si está presente
	if permission.BusinessTypeID > 0 {
		btID := permission.BusinessTypeID
		permissionModel.BusinessTypeID = &btID
	}

	if err := r.database.Conn(ctx).Create(&permissionModel).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al crear permiso")
		return "", err
	}

	r.logger.Info().
		Uint("permission_id", permissionModel.ID).
		Uint("resource_id", permission.ResourceID).
		Uint("action_id", permission.ActionID).
		Msg("Permiso creado exitosamente")

	return fmt.Sprintf("Permiso creado con ID: %d", permissionModel.ID), nil
}

// UpdatePermission actualiza un permiso existente
func (r *Repository) UpdatePermission(ctx context.Context, id uint, permission domain.Permission) (string, error) {
	// Verificar que el permiso existe
	var existingPermission models.Permission
	if err := r.database.Conn(ctx).Where("id = ?", id).First(&existingPermission).Error; err != nil {
		r.logger.Error().Uint("id", id).Err(err).Msg("Error al buscar permiso para actualizar")
		return "", fmt.Errorf("permiso no encontrado con ID: %d", id)
	}

	// Verificar que el Resource existe
	var resource models.Resource
	if err := r.database.Conn(ctx).Where("id = ?", permission.ResourceID).First(&resource).Error; err != nil {
		r.logger.Error().Uint("resource_id", permission.ResourceID).Err(err).Msg("Error al buscar resource")
		return "", fmt.Errorf("resource no encontrado con ID: %d", permission.ResourceID)
	}

	// Verificar que el Action existe
	var action models.Action
	if err := r.database.Conn(ctx).Where("id = ?", permission.ActionID).First(&action).Error; err != nil {
		r.logger.Error().Uint("action_id", permission.ActionID).Err(err).Msg("Error al buscar action")
		return "", fmt.Errorf("action no encontrada con ID: %d", permission.ActionID)
	}

	// Preparar los campos a actualizar
	updates := map[string]interface{}{
		"name":        permission.Name,
		"description": permission.Description,
		"resource_id": permission.ResourceID,
		"action_id":   permission.ActionID,
		"scope_id":    permission.ScopeID,
	}

	// Agregar business_type_id si está presente o si se necesita limpiar
	if permission.BusinessTypeID > 0 {
		btID := permission.BusinessTypeID
		updates["business_type_id"] = &btID
	} else if permission.BusinessTypeID == 0 && existingPermission.BusinessTypeID != nil {
		// Si viene 0 pero el existente tiene valor, limpiar
		updates["business_type_id"] = nil
	}

	// Actualizar el permiso
	if err := r.database.Conn(ctx).
		Model(&models.Permission{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		r.logger.Error().Err(err).Uint("id", id).Msg("Error al actualizar permiso")
		return "", err
	}

	r.logger.Info().
		Uint("permission_id", id).
		Str("name", permission.Name).
		Msg("Permiso actualizado exitosamente")

	return fmt.Sprintf("Permiso actualizado con ID: %d", id), nil
}

// DeletePermission elimina un permiso permanentemente
func (r *Repository) DeletePermission(ctx context.Context, id uint) (string, error) {
	r.logger.Info().Uint("id", id).Msg("Eliminando permiso permanentemente")

	// Usar Unscoped().Delete() para eliminación física (no soft delete)
	// Esto activará la eliminación en cascada de las relaciones definidas en el modelo
	result := r.database.Conn(ctx).Unscoped().Delete(&models.Permission{}, id)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Uint("id", id).Msg("Error al eliminar permiso")
		return "", result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.Warn().Uint("id", id).Msg("Permiso no encontrado para eliminar")
		return "", fmt.Errorf("permiso con ID %d no encontrado", id)
	}

	r.logger.Info().Uint("id", id).Msg("Permiso eliminado permanentemente")
	return fmt.Sprintf("Permiso eliminado con ID: %d", id), nil
}
