package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/auth/roles/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// CreateRole crea un nuevo rol en la base de datos
func (r *Repository) CreateRole(ctx context.Context, roleDTO domain.CreateRoleDTO) (*domain.Role, error) {
	// Crear el modelo de GORM
	role := models.Role{
		Name:           roleDTO.Name,
		Description:    roleDTO.Description,
		Level:          roleDTO.Level,
		IsSystem:       roleDTO.IsSystem,
		ScopeID:        roleDTO.ScopeID,
		BusinessTypeID: &roleDTO.BusinessTypeID, // Convertir a puntero
	}

	// Insertar en la base de datos
	err := r.database.Conn(ctx).Create(&role).Error
	if err != nil {
		r.logger.Error().
			Err(err).
			Str("name", roleDTO.Name).
			Msg("Error al crear rol en la base de datos")
		return nil, err
	}

	// Convertir a entidad de dominio
	domainRole := &domain.Role{
		ID:             role.ID,
		Name:           role.Name,
		Description:    role.Description,
		Level:          role.Level,
		IsSystem:       role.IsSystem,
		ScopeID:        role.ScopeID,
		BusinessTypeID: *role.BusinessTypeID, // Convertir de puntero a valor
		CreatedAt:      role.CreatedAt,
		UpdatedAt:      role.UpdatedAt,
	}

	r.logger.Info().
		Uint("role_id", role.ID).
		Str("name", role.Name).
		Msg("Rol creado exitosamente en la base de datos")

	return domainRole, nil
}

// GetRoleByID obtiene un rol por su ID
func (r *Repository) GetRoleByID(ctx context.Context, roleID uint) (*domain.Role, error) {
	var role models.Role

	err := r.database.Conn(ctx).
		Preload("Scope").
		Preload("BusinessType").
		Where("id = ?", roleID).
		First(&role).Error

	if err != nil {
		r.logger.Error().Uint("role_id", roleID).Err(err).Msg("Error al obtener rol por ID")
		return nil, err
	}

	businessTypeID := uint(0)
	businessTypeName := ""
	scopeName := ""
	scopeCode := ""
	if role.BusinessTypeID != nil {
		businessTypeID = *role.BusinessTypeID
	}
	if role.BusinessType != nil {
		businessTypeName = role.BusinessType.Name
	}
	if role.ScopeID > 0 && role.Scope.Name != "" {
		scopeName = role.Scope.Name
		scopeCode = role.Scope.Code
	}

	domainRole := &domain.Role{
		ID:               role.ID,
		Name:             role.Name,
		Description:      role.Description,
		Level:            role.Level,
		IsSystem:         role.IsSystem,
		ScopeID:          role.ScopeID,
		ScopeName:        scopeName,
		ScopeCode:        scopeCode,
		BusinessTypeID:   businessTypeID,
		BusinessTypeName: businessTypeName,
		CreatedAt:        role.CreatedAt,
		UpdatedAt:        role.UpdatedAt,
	}

	return domainRole, nil
}

// GetRoles obtiene todos los roles
func (r *Repository) GetRoles(ctx context.Context) ([]domain.Role, error) {
	var roles []models.Role

	err := r.database.Conn(ctx).
		Preload("Scope").
		Preload("BusinessType").
		Find(&roles).Error

	if err != nil {
		r.logger.Error().Err(err).Msg("Error al obtener roles")
		return nil, err
	}

	domainRoles := make([]domain.Role, len(roles))
	for i, role := range roles {
		businessTypeID := uint(0)
		businessTypeName := ""
		scopeName := ""
		scopeCode := ""

		if role.BusinessTypeID != nil {
			businessTypeID = *role.BusinessTypeID
			if role.BusinessType != nil {
				businessTypeName = role.BusinessType.Name
			}
		}

		if role.ScopeID > 0 && role.Scope.Name != "" {
			scopeName = role.Scope.Name
			scopeCode = role.Scope.Code
		}

		domainRoles[i] = domain.Role{
			ID:               role.ID,
			Name:             role.Name,
			Description:      role.Description,
			Level:            role.Level,
			IsSystem:         role.IsSystem,
			ScopeID:          role.ScopeID,
			ScopeName:        scopeName,
			ScopeCode:        scopeCode,
			BusinessTypeID:   businessTypeID,
			BusinessTypeName: businessTypeName,
			CreatedAt:        role.CreatedAt,
			UpdatedAt:        role.UpdatedAt,
		}
	}

	return domainRoles, nil
}

// GetRolesByLevel obtiene roles por nivel
func (r *Repository) GetRolesByLevel(ctx context.Context, level int) ([]domain.Role, error) {
	var roles []models.Role

	err := r.database.Conn(ctx).
		Preload("Scope").
		Preload("BusinessType").
		Where("level = ?", level).
		Find(&roles).Error

	if err != nil {
		r.logger.Error().Err(err).Int("level", level).Msg("Error al obtener roles por nivel")
		return nil, err
	}

	domainRoles := make([]domain.Role, len(roles))
	for i, role := range roles {
		businessTypeID := uint(0)
		businessTypeName := ""
		scopeName := ""
		scopeCode := ""
		if role.BusinessTypeID != nil {
			businessTypeID = *role.BusinessTypeID
		}
		if role.BusinessType != nil {
			businessTypeName = role.BusinessType.Name
		}
		if role.ScopeID > 0 && role.Scope.Name != "" {
			scopeName = role.Scope.Name
			scopeCode = role.Scope.Code
		}

		domainRoles[i] = domain.Role{
			ID:               role.ID,
			Name:             role.Name,
			Description:      role.Description,
			Level:            role.Level,
			IsSystem:         role.IsSystem,
			ScopeID:          role.ScopeID,
			ScopeName:        scopeName,
			ScopeCode:        scopeCode,
			BusinessTypeID:   businessTypeID,
			BusinessTypeName: businessTypeName,
			CreatedAt:        role.CreatedAt,
			UpdatedAt:        role.UpdatedAt,
		}
	}

	return domainRoles, nil
}

// GetRolesByScopeID obtiene roles por scope
func (r *Repository) GetRolesByScopeID(ctx context.Context, scopeID uint) ([]domain.Role, error) {
	var roles []models.Role

	err := r.database.Conn(ctx).
		Preload("Scope").
		Preload("BusinessType").
		Where("scope_id = ?", scopeID).
		Find(&roles).Error

	if err != nil {
		r.logger.Error().Err(err).Uint("scope_id", scopeID).Msg("Error al obtener roles por scope")
		return nil, err
	}

	domainRoles := make([]domain.Role, len(roles))
	for i, role := range roles {
		businessTypeID := uint(0)
		businessTypeName := ""
		scopeName := ""
		scopeCode := ""
		if role.BusinessTypeID != nil {
			businessTypeID = *role.BusinessTypeID
		}
		if role.BusinessType != nil {
			businessTypeName = role.BusinessType.Name
		}
		if role.ScopeID > 0 && role.Scope.Name != "" {
			scopeName = role.Scope.Name
			scopeCode = role.Scope.Code
		}

		domainRoles[i] = domain.Role{
			ID:               role.ID,
			Name:             role.Name,
			Description:      role.Description,
			Level:            role.Level,
			IsSystem:         role.IsSystem,
			ScopeID:          role.ScopeID,
			ScopeName:        scopeName,
			ScopeCode:        scopeCode,
			BusinessTypeID:   businessTypeID,
			BusinessTypeName: businessTypeName,
			CreatedAt:        role.CreatedAt,
			UpdatedAt:        role.UpdatedAt,
		}
	}

	return domainRoles, nil
}

// GetSystemRoles obtiene roles del sistema
func (r *Repository) GetSystemRoles(ctx context.Context) ([]domain.Role, error) {
	var roles []models.Role

	err := r.database.Conn(ctx).
		Preload("Scope").
		Preload("BusinessType").
		Where("is_system = ?", true).
		Find(&roles).Error

	if err != nil {
		r.logger.Error().Err(err).Msg("Error al obtener roles del sistema")
		return nil, err
	}

	domainRoles := make([]domain.Role, len(roles))
	for i, role := range roles {
		businessTypeID := uint(0)
		businessTypeName := ""
		scopeName := ""
		scopeCode := ""
		if role.BusinessTypeID != nil {
			businessTypeID = *role.BusinessTypeID
		}
		if role.BusinessType != nil {
			businessTypeName = role.BusinessType.Name
		}
		if role.ScopeID > 0 && role.Scope.Name != "" {
			scopeName = role.Scope.Name
			scopeCode = role.Scope.Code
		}

		domainRoles[i] = domain.Role{
			ID:               role.ID,
			Name:             role.Name,
			Description:      role.Description,
			Level:            role.Level,
			IsSystem:         role.IsSystem,
			ScopeID:          role.ScopeID,
			ScopeName:        scopeName,
			ScopeCode:        scopeCode,
			BusinessTypeID:   businessTypeID,
			BusinessTypeName: businessTypeName,
			CreatedAt:        role.CreatedAt,
			UpdatedAt:        role.UpdatedAt,
		}
	}

	return domainRoles, nil
}

// RoleExistsByName verifica si existe un rol con el nombre especificado
// excludeID permite excluir un rol específico (útil para actualizaciones)
func (r *Repository) RoleExistsByName(ctx context.Context, name string, excludeID *uint) (bool, error) {
	var count int64
	query := r.database.Conn(ctx).Model(&models.Role{}).Where("name = ?", name)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		r.logger.Error().Err(err).Str("name", name).Msg("Error verificando existencia de rol por nombre")
		return false, fmt.Errorf("error verificando existencia de rol por nombre: %w", err)
	}
	return count > 0, nil
}

// UpdateRole actualiza un rol existente en la base de datos
func (r *Repository) UpdateRole(ctx context.Context, id uint, roleDTO domain.UpdateRoleDTO) (*domain.Role, error) {
	// Obtener el rol existente
	var existingRole models.Role
	err := r.database.Conn(ctx).Where("id = ?", id).First(&existingRole).Error
	if err != nil {
		r.logger.Error().
			Err(err).
			Uint("role_id", id).
			Msg("Error al obtener rol existente para actualizar")
		return nil, err
	}

	// Actualizar solo los campos que no son nil
	updates := make(map[string]interface{})

	if roleDTO.Name != nil {
		updates["name"] = *roleDTO.Name
	}
	if roleDTO.Description != nil {
		updates["description"] = *roleDTO.Description
	}
	if roleDTO.Level != nil {
		updates["level"] = *roleDTO.Level
	}
	if roleDTO.IsSystem != nil {
		updates["is_system"] = *roleDTO.IsSystem
	}
	if roleDTO.ScopeID != nil {
		updates["scope_id"] = *roleDTO.ScopeID
	}
	if roleDTO.BusinessTypeID != nil {
		updates["business_type_id"] = *roleDTO.BusinessTypeID
	}

	// Si no hay campos para actualizar, devolver el rol existente
	if len(updates) == 0 {
		r.logger.Warn().
			Uint("role_id", id).
			Msg("No hay campos para actualizar")

		// Convertir a entidad de dominio
		businessTypeID := uint(0)
		scopeName := ""
		scopeCode := ""
		if existingRole.BusinessTypeID != nil {
			businessTypeID = *existingRole.BusinessTypeID
		}

		return &domain.Role{
			ID:             existingRole.ID,
			Name:           existingRole.Name,
			Description:    existingRole.Description,
			Level:          existingRole.Level,
			IsSystem:       existingRole.IsSystem,
			ScopeID:        existingRole.ScopeID,
			ScopeName:      scopeName,
			ScopeCode:      scopeCode,
			BusinessTypeID: businessTypeID,
			CreatedAt:      existingRole.CreatedAt,
			UpdatedAt:      existingRole.UpdatedAt,
		}, nil
	}

	// Actualizar en la base de datos
	err = r.database.Conn(ctx).Model(&existingRole).Updates(updates).Error
	if err != nil {
		r.logger.Error().
			Err(err).
			Uint("role_id", id).
			Msg("Error al actualizar rol en la base de datos")
		return nil, err
	}

	// Obtener el rol actualizado con las relaciones
	var updatedRole models.Role
	err = r.database.Conn(ctx).
		Preload("Scope").
		Preload("BusinessType").
		Where("id = ?", id).
		First(&updatedRole).Error
	if err != nil {
		r.logger.Error().
			Err(err).
			Uint("role_id", id).
			Msg("Error al obtener rol actualizado")
		return nil, err
	}

	// Convertir a entidad de dominio
	businessTypeID := uint(0)
	businessTypeName := ""
	scopeName := ""
	scopeCode := ""
	if updatedRole.BusinessTypeID != nil {
		businessTypeID = *updatedRole.BusinessTypeID
	}
	if updatedRole.BusinessType != nil {
		businessTypeName = updatedRole.BusinessType.Name
	}
	if updatedRole.ScopeID > 0 && updatedRole.Scope.Name != "" {
		scopeName = updatedRole.Scope.Name
		scopeCode = updatedRole.Scope.Code
	}

	domainRole := &domain.Role{
		ID:               updatedRole.ID,
		Name:             updatedRole.Name,
		Description:      updatedRole.Description,
		Level:            updatedRole.Level,
		IsSystem:         updatedRole.IsSystem,
		ScopeID:          updatedRole.ScopeID,
		ScopeName:        scopeName,
		ScopeCode:        scopeCode,
		BusinessTypeID:   businessTypeID,
		BusinessTypeName: businessTypeName,
		CreatedAt:        updatedRole.CreatedAt,
		UpdatedAt:        updatedRole.UpdatedAt,
	}

	r.logger.Info().
		Uint("role_id", updatedRole.ID).
		Str("name", updatedRole.Name).
		Msg("Rol actualizado exitosamente en la base de datos")

	return domainRole, nil
}
