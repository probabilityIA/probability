package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/auth/users/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (r *Repository) GetUserByEmailForAuth(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
	var userAuth domain.UserAuthInfo
	if err := r.database.Conn(ctx).
		Model(&models.User{}).
		Select("id, name, email, password, phone, avatar_url, is_active, last_login_at, created_at, updated_at, deleted_at").
		Where("email = ?", email).
		First(&userAuth).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error().Str("email", email).Err(err).Msg("Error al obtener usuario por email")
		return nil, err
	}
	return &userAuth, nil
}

func (r *Repository) GetUserRoles(ctx context.Context, userID uint) ([]domain.Role, error) {
	var userRoles []models.BusinessStaff
	var roles []domain.Role

	err := r.database.Conn(ctx).
		Model(&models.BusinessStaff{}).
		Preload("Role.Scope").
		Where("user_id = ?", userID).
		Find(&userRoles).Error

	if err != nil {
		r.logger.Error().Uint("user_id", userID).Msg("Error al obtener roles del usuario")
		return nil, err
	}

	for _, role := range userRoles {
		roles = append(roles, domain.Role{
			ID:          role.Role.ID,
			Name:        role.Role.Name,
			Description: role.Role.Description,
			Level:       role.Role.Level,
			IsSystem:    role.Role.IsSystem,
			ScopeID:     role.Role.ScopeID,
			ScopeName:   role.Role.Scope.Name,
			ScopeCode:   role.Role.Scope.Code,
			CreatedAt:   role.Role.CreatedAt,
			UpdatedAt:   role.Role.UpdatedAt,
		})
	}

	return roles, nil
}

// GetUserRoleByBusiness obtiene el rol de un usuario para un business específico desde user_roles
// Valida que el rol coincida con el tipo de business
func (r *Repository) GetUserRoleByBusiness(ctx context.Context, userID uint, businessID uint) (*domain.Role, error) {
	// Obtener el business para conocer su tipo
	var business models.Business
	if err := r.database.Conn(ctx).
		Preload("BusinessType").
		Where("id = ?", businessID).
		First(&business).Error; err != nil {
		r.logger.Error().
			Uint("business_id", businessID).
			Err(err).
			Msg("Error al obtener business para validar tipo")
		return nil, err
	}

	// Buscar roles del usuario que coincidan con el tipo de business usando business_staff
	var staffEntries []models.BusinessStaff
	if err := r.database.Conn(ctx).
		Preload("Role.Scope").
		Preload("Role.BusinessType").
		Where("user_id = ? AND business_id = ?", userID, businessID).
		Find(&staffEntries).Error; err != nil {
		r.logger.Error().
			Uint("user_id", userID).
			Uint("business_id", businessID).
			Err(err).
			Msg("Error al obtener roles del usuario desde business_staff")
		return nil, err
	}

	if len(staffEntries) == 0 {
		return nil, nil // No tiene asociación en business_staff
	}

	var selectedRole *models.Role
	for _, entry := range staffEntries {
		role := entry.Role
		if role.ID == 0 {
			continue
		}
		if role.BusinessTypeID != nil && *role.BusinessTypeID == business.BusinessTypeID {
			selectedRole = &role
			break
		}
		if selectedRole == nil {
			selectedRole = &role
		}
	}

	if selectedRole == nil {
		return nil, nil
	}

	businessTypeID := uint(0)
	businessTypeName := ""
	if selectedRole.BusinessTypeID != nil {
		businessTypeID = *selectedRole.BusinessTypeID
	}
	if selectedRole.BusinessType != nil {
		businessTypeName = selectedRole.BusinessType.Name
	}

	role := &domain.Role{
		ID:               selectedRole.ID,
		Name:             selectedRole.Name,
		Description:      selectedRole.Description,
		Level:            selectedRole.Level,
		IsSystem:         selectedRole.IsSystem,
		ScopeID:          selectedRole.ScopeID,
		ScopeName:        selectedRole.Scope.Name,
		ScopeCode:        selectedRole.Scope.Code,
		BusinessTypeID:   businessTypeID,
		BusinessTypeName: businessTypeName,
		CreatedAt:        selectedRole.CreatedAt,
		UpdatedAt:        selectedRole.UpdatedAt,
	}

	return role, nil
}

func (r *Repository) GetRolePermissions(ctx context.Context, roleID uint) ([]domain.Permission, error) {
	var role models.Role
	var permissions []domain.Permission

	err := r.database.Conn(ctx).
		Model(&models.Role{}).
		Preload("Permissions.Scope").
		Preload("Permissions.Resource").
		Preload("Permissions.Action").
		Where("id = ?", roleID).
		First(&role).Error

	if err != nil {
		r.logger.Error().Uint("role_id", roleID).Msg("Error al obtener permisos del rol")
		return nil, err
	}

	for _, permission := range role.Permissions {
		businessTypeID := uint(0)
		if permission.BusinessTypeID != nil {
			businessTypeID = *permission.BusinessTypeID
		}

		permissions = append(permissions, domain.Permission{
			ID:               permission.Model.ID,
			Name:             permission.Name,
			Description:      permission.Description,
			Resource:         permission.Resource.Name,
			Action:           permission.Action.Name,
			ResourceID:       permission.ResourceID,
			ActionID:         permission.ActionID,
			ScopeID:          permission.ScopeID,
			BusinessTypeID:   businessTypeID,
			BusinessTypeName: "", // Se puede agregar si se necesita
		})
	}

	return permissions, nil
}

func (r *Repository) UpdateLastLogin(ctx context.Context, userID uint) error {
	now := time.Now()
	if err := r.database.Conn(ctx).Table("user").Where("id = ?", userID).Update("last_login_at", now).Error; err != nil {
		r.logger.Error().Uint("user_id", userID).Msg("Error al actualizar último login")
		return err
	}
	return nil
}

func (r *Repository) GetUserByIDForAuth(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
	var userAuth domain.UserAuthInfo
	if err := r.database.Conn(ctx).
		Model(&models.User{}).
		Select("id, name, email, password, phone, avatar_url, is_active, last_login_at, created_at, updated_at, deleted_at").
		Where("id = ?", userID).
		First(&userAuth).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error().Uint("user_id", userID).Err(err).Msg("Error al obtener usuario por ID")
		return nil, err
	}
	return &userAuth, nil
}

func (r *Repository) ChangePassword(ctx context.Context, userID uint, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		r.logger.Error().Err(err).Msg("Error al hashear nueva contraseña")
		return fmt.Errorf("error al procesar contraseña")
	}

	if err := r.database.Conn(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("password", string(hashedPassword)).Error; err != nil {
		r.logger.Error().Uint("user_id", userID).Err(err).Msg("Error al actualizar contraseña")
		return err
	}

	return nil
}

func (r *Repository) GetUsers(ctx context.Context, filters domain.UserFilters) ([]domain.UserQueryDTO, int64, error) {
	var total int64

	// Query base usando Model para que GORM determine el nombre de tabla automáticamente
	query := r.database.Conn(ctx).
		Model(&models.User{}).
		Select(`"user".id, "user".name, "user".email, "user".phone, "user".avatar_url, 
			"user".is_active, "user".last_login_at, "user".scope_id, 
			"user".created_at, "user".updated_at, "user".deleted_at,
			scope.code as scope_code, scope.name as scope_name`).
		Joins("LEFT JOIN scope ON scope.id = \"user\".scope_id")

	// FILTRO DE SEGURIDAD POR SCOPE
	// Si el usuario que solicita es de scope "business", solo puede ver usuarios "business"
	if filters.RequesterScope == "business" {
		r.logger.Info().Str("requester_scope", filters.RequesterScope).Msg("Filtrando usuarios: solo scope business")
		// Obtener el ID del scope "business"
		var businessScope models.Scope
		if err := r.database.Conn(ctx).Where("code = ?", "business").First(&businessScope).Error; err == nil {
			query = query.Where("scope_id = ? OR scope_id IS NULL", businessScope.ID)
		} else {
			// Si no existe el scope "business", filtrar por scope_id NULL (comportamiento por defecto)
			query = query.Where("scope_id IS NULL")
		}

		// Además, si tiene un business_id, solo mostrar usuarios de sus mismos negocios
		if filters.BusinessID != nil {
			subquery := r.database.Conn(ctx).
				Table("user_businesses").
				Select("user_id").
				Where("business_id = ?", *filters.BusinessID)
			query = query.Where("\"user\".id IN (?)", subquery)
		}
	}

	// Filtros opcionales
	if filters.Email != "" {
		query = query.Where("email LIKE ?", "%"+filters.Email+"%")
	}
	if filters.Name != "" {
		query = query.Where("name LIKE ?", "%"+filters.Name+"%")
	}
	if filters.Phone != "" {
		query = query.Where("phone LIKE ?", "%"+filters.Phone+"%")
	}
	if len(filters.UserIDs) > 0 {
		query = query.Where("\"user\".id IN ?", filters.UserIDs)
	}
	if filters.IsActive != nil {
		query = query.Where("is_active = ?", *filters.IsActive)
	}
	if filters.RoleID != nil {
		subquery := r.database.Conn(ctx).
			Table("user_roles").
			Select("user_id").
			Where("role_id = ?", *filters.RoleID)
		query = query.Where("\"user\".id IN (?)", subquery)
	}
	if filters.BusinessID != nil && filters.RequesterScope != "business" {
		// Solo aplicar si no es usuario business (ya se aplicó arriba)
		subquery := r.database.Conn(ctx).
			Table("user_businesses").
			Select("user_id").
			Where("business_id = ?", *filters.BusinessID)
		query = query.Where("\"user\".id IN (?)", subquery)
	}
	// Filtro por scope específico
	if filters.ScopeID != nil {
		query = query.Where("scope_id = ?", *filters.ScopeID)
	}
	if filters.ScopeCode != "" {
		query = query.Where("scope.code = ?", filters.ScopeCode)
	}

	// Ordenamiento
	if filters.SortBy != "" && filters.SortOrder != "" {
		query = query.Order(filters.SortBy + " " + filters.SortOrder)
	} else {
		query = query.Order("created_at desc")
	}

	// Contar total
	countQuery := query.Session(&gorm.Session{})
	if err := countQuery.Count(&total).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al contar usuarios")
		return nil, 0, err
	}

	// Paginación
	offset := (filters.Page - 1) * filters.PageSize
	query = query.Offset(offset).Limit(filters.PageSize)

	// Ejecutar query
	rows, err := query.Rows()
	if err != nil {
		r.logger.Error().Err(err).Msg("Error al obtener usuarios")
		return nil, 0, err
	}
	defer rows.Close()

	var users []domain.UserQueryDTO
	for rows.Next() {
		var user domain.UserQueryDTO
		if err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.Phone, &user.AvatarURL,
			&user.IsActive, &user.LastLoginAt, &user.ScopeID,
			&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
			&user.ScopeCode, &user.ScopeName,
		); err != nil {
			r.logger.Error().Err(err).Msg("Error al escanear usuario")
			continue
		}
		users = append(users, user)
	}

	return users, total, nil
}

func (r *Repository) GetUserBusinesses(ctx context.Context, userID uint) ([]domain.BusinessInfoEntity, error) {
	// Obtener relaciones business_staff con preload de Business y Role
	var businessStaffList []models.BusinessStaff

	if err := r.database.Conn(ctx).
		Preload("Business.BusinessType").
		Preload("Role").
		Where("user_id = ? AND business_id IS NOT NULL", userID).
		Find(&businessStaffList).Error; err != nil {
		r.logger.Error().Uint("user_id", userID).Err(err).Msg("Error al obtener negocios del usuario desde business_staff")
		return nil, err
	}

	result := make([]domain.BusinessInfoEntity, 0, len(businessStaffList))
	for _, bs := range businessStaffList {
		if bs.BusinessID == nil || bs.Business.ID == 0 {
			continue // Saltar si no hay business
		}

		businessInfo := domain.BusinessInfoEntity{
			ID:              bs.Business.ID,
			Name:            bs.Business.Name,
			Code:            bs.Business.Code,
			BusinessTypeID:  bs.Business.BusinessTypeID,
			Timezone:        bs.Business.Timezone,
			Address:         bs.Business.Address,
			Description:     bs.Business.Description,
			LogoURL:         bs.Business.LogoURL,
			PrimaryColor:    bs.Business.PrimaryColor,
			SecondaryColor:  bs.Business.SecondaryColor,
			TertiaryColor:   bs.Business.TertiaryColor,
			QuaternaryColor: bs.Business.QuaternaryColor,
			NavbarImageURL:  bs.Business.NavbarImageURL,
			CustomDomain: func() string {
				if bs.Business.CustomDomain != nil {
					return *bs.Business.CustomDomain
				}
				return ""
			}(),
			IsActive:           bs.Business.IsActive,
			EnableDelivery:     bs.Business.EnableDelivery,
			EnablePickup:       bs.Business.EnablePickup,
			EnableReservations: bs.Business.EnableReservations,
		}

		if bs.Business.BusinessType.ID != 0 {
			businessInfo.BusinessTypeName = bs.Business.BusinessType.Name
			businessInfo.BusinessTypeCode = bs.Business.BusinessType.Code
		}

		result = append(result, businessInfo)
	}

	return result, nil
}

func (r *Repository) CreateUser(ctx context.Context, user domain.UsersEntity) (uint, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		r.logger.Error().Err(err).Msg("Error al hashear contraseña")
		return 0, fmt.Errorf("error al procesar contraseña")
	}
	user.Password = string(hashedPassword)

	// Usar el mapper para convertir entities.User a models.User
	userModel := mappers.CreateUserModel(user)

	if err := r.database.Conn(ctx).Create(&userModel).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al crear usuario")
		return 0, err
	}

	return userModel.Model.ID, nil
}

func (r *Repository) UpdateUser(ctx context.Context, id uint, user domain.UsersEntity) (string, error) {
	if err := r.database.Conn(ctx).Model(&models.User{}).Where("id = ?", id).Updates(&user).Error; err != nil {
		r.logger.Error().Uint("id", id).Err(err).Msg("Error al actualizar usuario")
		return "", err
	}

	return fmt.Sprintf("Usuario actualizado con ID: %d", id), nil
}

func (r *Repository) DeleteUser(ctx context.Context, id uint) (string, error) {
	if err := r.database.Conn(ctx).Delete(&models.User{}, id).Error; err != nil {
		r.logger.Error().Uint("id", id).Err(err).Msg("Error al eliminar usuario")
		return "", err
	}

	return fmt.Sprintf("Usuario eliminado con ID: %d", id), nil
}

func (r *Repository) AssignRolesToUser(ctx context.Context, userID uint, roleIDs []uint) error {
	db := r.database.Conn(ctx)

	// Verificar que el usuario existe
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		r.logger.Error().Uint("user_id", userID).Err(err).Msg("Error al encontrar usuario")
		return err
	}

	// Verificar que los roles existen
	if len(roleIDs) > 0 {
		var count int64
		if err := db.Model(&models.Role{}).Where("id IN ?", roleIDs).Count(&count).Error; err != nil {
			r.logger.Error().Uint("user_id", userID).Err(err).Msg("Error al verificar roles")
			return err
		}
		if count != int64(len(roleIDs)) {
			r.logger.Error().
				Uint("user_id", userID).
				Int64("expected", int64(len(roleIDs))).
				Int64("found", count).
				Msg("Algunos roles no existen")
			return fmt.Errorf("algunos roles no existen")
		}
	}

	// Iniciar transacción
	return db.Transaction(func(tx *gorm.DB) error {
		// Eliminar todos los roles existentes
		if err := tx.Table("user_roles").Where("user_id = ?", userID).Delete(nil).Error; err != nil {
			r.logger.Error().Uint("user_id", userID).Err(err).Msg("Error al eliminar roles existentes")
			return err
		}

		// Insertar los nuevos roles si hay alguno
		if len(roleIDs) > 0 {
			values := make([]map[string]interface{}, len(roleIDs))
			for i, roleID := range roleIDs {
				values[i] = map[string]interface{}{
					"user_id": userID,
					"role_id": roleID,
				}
			}

			if err := tx.Table("user_roles").CreateInBatches(values, 100).Error; err != nil {
				r.logger.Error().
					Uint("user_id", userID).
					Err(err).
					Msg("Error al insertar roles")
				return err
			}
		}

		r.logger.Info().
			Uint("user_id", userID).
			Int("role_count", len(roleIDs)).
			Msg("Roles asignados exitosamente")

		return nil
	})
}

func (r *Repository) CreateAPIKey(ctx context.Context, apiKey domain.APIKey, keyHash string) (uint, error) {
	dbAPIKey := mappers.CreateAPIKeyModel(apiKey, keyHash)

	if err := r.database.Conn(ctx).Model(&models.APIKey{}).Create(&dbAPIKey).Error; err != nil {
		r.logger.Error().Err(err).
			Uint("user_id", apiKey.UserID).
			Uint("business_id", apiKey.BusinessID).
			Msg("Error al crear API Key")
		return 0, err
	}

	return dbAPIKey.Model.ID, nil
}

func (r *Repository) ValidateAPIKey(ctx context.Context, apiKey string) (*domain.APIKey, error) {
	var dbAPIKeys []models.APIKey

	err := r.database.Conn(ctx).
		Model(&models.APIKey{}).
		Where("revoked = ?", false).
		Find(&dbAPIKeys).Error

	if err != nil {
		r.logger.Error().Err(err).Msg("Error al buscar API Keys")
		return nil, err
	}

	for _, dbAPIKey := range dbAPIKeys {
		if apiKey == dbAPIKey.KeyHash {
			if err := r.UpdateAPIKeyLastUsed(ctx, dbAPIKey.Model.ID); err != nil {
				r.logger.Warn().Uint("api_key_id", dbAPIKey.Model.ID).Err(err).Msg("Error al actualizar último uso")
			}

			entity := mappers.ToAPIKeyEntity(dbAPIKey)
			return &entity, nil
		}
	}

	return nil, nil
}

func (r *Repository) UpdateAPIKeyLastUsed(ctx context.Context, apiKeyID uint) error {
	now := time.Now()
	if err := r.database.Conn(ctx).
		Model(&models.APIKey{}).
		Where("id = ?", apiKeyID).
		Update("last_used_at", now).Error; err != nil {
		r.logger.Error().Uint("api_key_id", apiKeyID).Err(err).Msg("Error al actualizar último uso de API Key")
		return err
	}

	return nil
}

func (r *Repository) GetAPIKeysByUser(ctx context.Context, userID uint) ([]domain.APIKeyInfo, error) {
	var dbAPIKeys []models.APIKey

	err := r.database.Conn(ctx).
		Model(&models.APIKey{}).
		Where("user_id = ?", userID).
		Find(&dbAPIKeys).Error

	if err != nil {
		r.logger.Error().Uint("user_id", userID).Err(err).Msg("Error al obtener API Keys del usuario")
		return nil, err
	}

	apiKeys := mappers.ToAPIKeyInfoEntitySlice(dbAPIKeys)
	return apiKeys, nil
}

func (r *Repository) RevokeAPIKey(ctx context.Context, apiKeyID uint) error {
	now := time.Now()
	if err := r.database.Conn(ctx).
		Model(&models.APIKey{}).
		Where("id = ?", apiKeyID).
		Updates(map[string]interface{}{
			"revoked":    true,
			"revoked_at": now,
			"updated_at": now,
		}).Error; err != nil {
		return err
	}
	return nil
}

// GetBusinessConfiguredResourcesIDs obtiene los IDs de recursos ACTIVOS configurados para un business específico
func (r *Repository) GetBusinessConfiguredResourcesIDs(ctx context.Context, businessID uint) ([]uint, error) {
	var resourcesIDs []uint

	// Obtener solo los resource_ids que están activos (active = true) en business_resource_configured
	err := r.database.Conn(ctx).
		Model(&models.BusinessResourceConfigured{}).
		Where("business_id = ? AND active = ? AND deleted_at IS NULL", businessID, true).
		Pluck("resource_id", &resourcesIDs).Error

	if err != nil {
		r.logger.Error().Err(err).Uint("business_id", businessID).Msg("Error al obtener recursos activos configurados del business")
		return nil, err
	}

	r.logger.Info().Uint("business_id", businessID).Int("resources_count", len(resourcesIDs)).Msg("Recursos activos configurados del business obtenidos exitosamente")

	return resourcesIDs, nil
}

// AssignPermissionsToRole asigna permisos a un rol
func (r *Repository) AssignPermissionsToRole(ctx context.Context, roleID uint, permissionIDs []uint) error {
	db := r.database.Conn(ctx)

	// Verificar que el rol existe y obtener su business_type_id
	var role models.Role
	err := db.Where("id = ?", roleID).First(&role).Error
	if err != nil {
		r.logger.Error().Uint("role_id", roleID).Err(err).Msg("Error al encontrar rol")
		return fmt.Errorf("rol no encontrado")
	}

	// Si el rol tiene business_type_id, validar que los permisos pertenezcan al mismo business_type
	if role.BusinessTypeID != nil && len(permissionIDs) > 0 {
		var permissions []models.Permission
		err := db.Where("id IN ?", permissionIDs).Find(&permissions).Error
		if err != nil {
			r.logger.Error().Uint("role_id", roleID).Err(err).Msg("Error al verificar permisos")
			return fmt.Errorf("error al verificar permisos")
		}

		// Validar que todos los permisos pertenezcan al mismo business_type o sean genéricos
		for _, permission := range permissions {
			// Si el permiso tiene business_type_id, debe coincidir con el del rol
			if permission.BusinessTypeID != nil {
				if *permission.BusinessTypeID != *role.BusinessTypeID {
					roleBTID := uint(0)
					if role.BusinessTypeID != nil {
						roleBTID = *role.BusinessTypeID
					}

					r.logger.Error().
						Uint("role_id", roleID).
						Uint("permission_id", permission.ID).
						Uint("role_business_type", roleBTID).
						Uint("permission_business_type", *permission.BusinessTypeID).
						Msg("El permiso no pertenece al mismo business_type que el rol")
					return fmt.Errorf("el permiso con ID %d no pertenece al mismo business_type que el rol", permission.ID)
				}
			}
			// Si es NULL, es genérico y se puede asignar a cualquier tipo
		}
	}

	// Iniciar transacción
	return db.Transaction(func(tx *gorm.DB) error {
		// Eliminar todos los permisos existentes del rol
		if err := tx.Table("role_permissions").Where("role_id = ?", roleID).Delete(nil).Error; err != nil {
			r.logger.Error().Uint("role_id", roleID).Err(err).Msg("Error al eliminar permisos existentes")
			return err
		}

		// Insertar los nuevos permisos si hay alguno
		if len(permissionIDs) > 0 {
			values := make([]map[string]interface{}, len(permissionIDs))
			for i, permissionID := range permissionIDs {
				values[i] = map[string]interface{}{
					"role_id":       roleID,
					"permission_id": permissionID,
				}
			}

			if err := tx.Table("role_permissions").CreateInBatches(values, 100).Error; err != nil {
				r.logger.Error().
					Uint("role_id", roleID).
					Err(err).
					Msg("Error al insertar permisos")
				return err
			}
		}

		r.logger.Info().
			Uint("role_id", roleID).
			Int("permission_count", len(permissionIDs)).
			Msg("Permisos asignados exitosamente al rol")

		return nil
	})
}

// RemovePermissionFromRole elimina un permiso específico de un rol
func (r *Repository) RemovePermissionFromRole(ctx context.Context, roleID uint, permissionID uint) error {
	db := r.database.Conn(ctx)

	// Verificar que el rol existe
	var role models.Role
	err := db.Where("id = ?", roleID).First(&role).Error
	if err != nil {
		r.logger.Error().Uint("role_id", roleID).Err(err).Msg("Error al encontrar rol")
		return fmt.Errorf("rol no encontrado")
	}

	// Verificar que el permiso existe
	var permission models.Permission
	err = db.Where("id = ?", permissionID).First(&permission).Error
	if err != nil {
		r.logger.Error().Uint("permission_id", permissionID).Err(err).Msg("Error al encontrar permiso")
		return fmt.Errorf("permiso no encontrado")
	}

	// Eliminar la relación
	err = db.Table("role_permissions").
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete(nil).Error

	if err != nil {
		r.logger.Error().
			Uint("role_id", roleID).
			Uint("permission_id", permissionID).
			Err(err).
			Msg("Error al eliminar permiso del rol")
		return err
	}

	r.logger.Info().
		Uint("role_id", roleID).
		Uint("permission_id", permissionID).
		Msg("Permiso eliminado exitosamente del rol")

	return nil
}

// GetRolePermissionsIDs obtiene los IDs de los permisos asignados a un rol
func (r *Repository) GetRolePermissionsIDs(ctx context.Context, roleID uint) ([]uint, error) {
	db := r.database.Conn(ctx)

	// Verificar que el rol existe
	var role models.Role
	err := db.Where("id = ?", roleID).First(&role).Error
	if err != nil {
		r.logger.Error().Uint("role_id", roleID).Err(err).Msg("Error al encontrar rol")
		return nil, fmt.Errorf("rol no encontrado")
	}

	var permissionIDs []uint
	err = db.Table("role_permissions").
		Where("role_id = ?", roleID).
		Pluck("permission_id", &permissionIDs).Error

	if err != nil {
		r.logger.Error().Uint("role_id", roleID).Err(err).Msg("Error al obtener permisos del rol")
		return nil, err
	}

	r.logger.Info().
		Uint("role_id", roleID).
		Int("permission_count", len(permissionIDs)).
		Msg("IDs de permisos del rol obtenidos exitosamente")

	return permissionIDs, nil
}

// GetUserRoleIDFromBusinessStaff retorna el role_id de business_staff para un usuario y business dado
// Si businessID es nil, busca filas con business_id NULL (caso super)
func (r *Repository) GetUserRoleIDFromBusinessStaff(ctx context.Context, userID uint, businessID *uint) (*uint, error) {
	db := r.database.Conn(ctx)
	var bs models.BusinessStaff
	q := db.Where("user_id = ?", userID)
	if businessID == nil {
		q = q.Where("business_id IS NULL")
	} else {
		q = q.Where("business_id = ?", *businessID)
	}
	if err := q.First(&bs).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error().Uint("user_id", userID).Err(err).Msg("Error consultando business_staff")
		return nil, err
	}
	return bs.RoleID, nil
}

// GetBusinessStaffRelation obtiene la relación completa user-business-role desde business_staff
func (r *Repository) GetBusinessStaffRelation(ctx context.Context, userID uint, businessID *uint) (*domain.BusinessStaffRelation, error) {
	db := r.database.Conn(ctx)
	var bs models.BusinessStaff
	q := db.Where("user_id = ?", userID)
	if businessID == nil {
		q = q.Where("business_id IS NULL")
		r.logger.Info().Uint("user_id", userID).Msg("Buscando relación con business_id NULL")
	} else {
		q = q.Where("business_id = ?", *businessID)
		r.logger.Info().Uint("user_id", userID).Uint("business_id", *businessID).Msg("Buscando relación con business_id específico")
	}

	// Preload de Business (si existe) y Role (si existe)
	if err := q.Preload("Business.BusinessType").Preload("Role").First(&bs).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Warn().Uint("user_id", userID).Any("business_id", businessID).Msg("Relación no encontrada en business_staff")
			return nil, nil
		}
		r.logger.Error().Uint("user_id", userID).Err(err).Msg("Error consultando business_staff")
		return nil, err
	}

	r.logger.Info().
		Uint("user_id", bs.UserID).
		Any("business_id", bs.BusinessID).
		Any("role_id", bs.RoleID).
		Bool("has_business", bs.Business.ID != 0).
		Bool("has_role", bs.Role.ID != 0).
		Msg("Relación business_staff encontrada")

	rel := &domain.BusinessStaffRelation{
		UserID:     bs.UserID,
		BusinessID: bs.BusinessID,
		RoleID:     bs.RoleID,
	}

	// Si hay business, mapear la información
	if bs.BusinessID != nil && bs.Business.ID != 0 {
		businessInfo := domain.BusinessInfoEntity{
			ID:              bs.Business.ID,
			Name:            bs.Business.Name,
			Code:            bs.Business.Code,
			BusinessTypeID:  bs.Business.BusinessTypeID,
			Timezone:        bs.Business.Timezone,
			Address:         bs.Business.Address,
			Description:     bs.Business.Description,
			LogoURL:         bs.Business.LogoURL,
			PrimaryColor:    bs.Business.PrimaryColor,
			SecondaryColor:  bs.Business.SecondaryColor,
			TertiaryColor:   bs.Business.TertiaryColor,
			QuaternaryColor: bs.Business.QuaternaryColor,
			NavbarImageURL:  bs.Business.NavbarImageURL,
			CustomDomain: func() string {
				if bs.Business.CustomDomain != nil {
					return *bs.Business.CustomDomain
				}
				return ""
			}(),
			IsActive:           bs.Business.IsActive,
			EnableDelivery:     bs.Business.EnableDelivery,
			EnablePickup:       bs.Business.EnablePickup,
			EnableReservations: bs.Business.EnableReservations,
		}
		if bs.Business.BusinessType.ID != 0 {
			businessInfo.BusinessTypeName = bs.Business.BusinessType.Name
			businessInfo.BusinessTypeCode = bs.Business.BusinessType.Code
		}
		rel.Business = &businessInfo
	}

	return rel, nil
}
