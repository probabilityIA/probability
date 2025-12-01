package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/auth/roles/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// AssignPermissionsToRole asigna permisos a un rol
func (r *Repository) AssignPermissionsToRole(ctx context.Context, roleID uint, permissionIDs []uint) error {
	// Verificar que el rol existe
	var role models.Role
	if err := r.database.Conn(ctx).First(&role, roleID).Error; err != nil {
		r.logger.Error().Err(err).Uint("role_id", roleID).Msg("Error al buscar rol para asignar permisos")
		return err
	}

	// Buscar los permisos
	var permissions []models.Permission
	if err := r.database.Conn(ctx).Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al buscar permisos para asignar")
		return err
	}

	// Reemplazar los permisos del rol
	if err := r.database.Conn(ctx).Model(&role).Association("Permissions").Replace(permissions); err != nil {
		r.logger.Error().Err(err).Uint("role_id", roleID).Msg("Error al asignar permisos al rol")
		return err
	}

	r.logger.Info().
		Uint("role_id", roleID).
		Int("permissions_count", len(permissions)).
		Msg("Permisos asignados al rol exitosamente")

	return nil
}

// GetRolePermissions obtiene los permisos asignados a un rol
func (r *Repository) GetRolePermissions(ctx context.Context, roleID uint) ([]domain.Permission, error) {
	var role models.Role

	// Cargar el rol con sus permisos y las relaciones de los permisos
	err := r.database.Conn(ctx).
		Preload("Permissions").
		Preload("Permissions.Resource").
		Preload("Permissions.Action").
		Preload("Permissions.Scope").
		Preload("Permissions.BusinessType").
		First(&role, roleID).Error

	if err != nil {
		r.logger.Error().Err(err).Uint("role_id", roleID).Msg("Error al obtener permisos del rol")
		return nil, err
	}

	domainPermissions := make([]domain.Permission, len(role.Permissions))
	for i, p := range role.Permissions {
		businessTypeID := uint(0)
		businessTypeName := ""
		scopeName := ""
		scopeCode := ""

		if p.BusinessTypeID != nil {
			businessTypeID = *p.BusinessTypeID
			if p.BusinessType != nil {
				businessTypeName = p.BusinessType.Name
			}
		}

		if p.Scope.Name != "" {
			scopeName = p.Scope.Name
			scopeCode = p.Scope.Code
		}

		domainPermissions[i] = domain.Permission{
			ID:               p.ID,
			Name:             p.Name,
			Description:      p.Description,
			Resource:         p.Resource.Name,
			Action:           p.Action.Name,
			ResourceID:       p.ResourceID,
			ActionID:         p.ActionID,
			ScopeID:          p.ScopeID,
			ScopeName:        scopeName,
			ScopeCode:        scopeCode,
			BusinessTypeID:   businessTypeID,
			BusinessTypeName: businessTypeName,
		}
	}

	return domainPermissions, nil
}

// RemovePermissionFromRole elimina un permiso de un rol
func (r *Repository) RemovePermissionFromRole(ctx context.Context, roleID uint, permissionID uint) error {
	// Verificar que el rol existe
	var role models.Role
	if err := r.database.Conn(ctx).First(&role, roleID).Error; err != nil {
		r.logger.Error().Err(err).Uint("role_id", roleID).Msg("Error al buscar rol para eliminar permiso")
		return err
	}

	// Buscar el permiso
	var permission models.Permission
	if err := r.database.Conn(ctx).First(&permission, permissionID).Error; err != nil {
		r.logger.Error().Err(err).Uint("permission_id", permissionID).Msg("Error al buscar permiso para eliminar")
		return err
	}

	// Eliminar la asociación
	if err := r.database.Conn(ctx).Model(&role).Association("Permissions").Delete(&permission); err != nil {
		r.logger.Error().Err(err).Uint("role_id", roleID).Uint("permission_id", permissionID).Msg("Error al eliminar permiso del rol")
		return err
	}

	r.logger.Info().
		Uint("role_id", roleID).
		Uint("permission_id", permissionID).
		Msg("Permiso eliminado del rol exitosamente")

	return nil
}
