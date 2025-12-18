package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/auth/login/internal/domain"
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

// AssignBusinessStaffRelationships asigna relaciones usuario-negocio-rol usando la tabla business_staff

// GetBusinessStaffRelationships obtiene todas las relaciones business_staff de un usuario con información completa

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
	var user domain.UserAuthInfo
	if err := r.database.Conn(ctx).
		Model(&models.User{}).
		Where("email = ?", email).
		First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error().Str("email", email).Err(err).Msg("Error al obtener usuario por email")
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id uint) (*domain.UserAuthInfo, error) {
	var user domain.UserAuthInfo
	if err := r.database.Conn(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error().Uint("id", id).Err(err).Msg("Error al obtener usuario por ID")
		return nil, err
	}
	return &user, nil
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

// RemovePermissionFromRole elimina un permiso específico de un rol

// GetRolePermissionsIDs obtiene los IDs de los permisos asignados a un rol

// GetUserRoleIDFromBusinessStaff retorna el role_id de business_staff para un usuario y business dado
// Si businessID es nil, busca filas con business_id NULL (caso super)

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

// AssignRoleToUserBusiness asigna o actualiza roles de un usuario en múltiples businesses
// Valida que el usuario esté asociado a cada business y que cada rol sea del mismo tipo de business

func (r *Repository) GetBusinessByID(ctx context.Context, businessID uint) (*domain.BusinessInfo, error) {
	var business models.Business
	if err := r.database.Conn(ctx).
		Preload("BusinessType").
		Where("id = ?", businessID).
		First(&business).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error().Uint("business_id", businessID).Err(err).Msg("Error al obtener business por ID")
		return nil, err
	}

	return &domain.BusinessInfo{
		ID:             business.ID,
		Name:           business.Name,
		Code:           business.Code,
		BusinessTypeID: business.BusinessTypeID,
		BusinessType: domain.BusinessTypeInfo{
			ID:          business.BusinessType.ID,
			Name:        business.BusinessType.Name,
			Code:        business.BusinessType.Code,
			Description: business.BusinessType.Description,
			Icon:        business.BusinessType.Icon,
		},
		Timezone:        business.Timezone,
		Address:         business.Address,
		Description:     business.Description,
		LogoURL:         business.LogoURL,
		PrimaryColor:    business.PrimaryColor,
		SecondaryColor:  business.SecondaryColor,
		TertiaryColor:   business.TertiaryColor,
		QuaternaryColor: business.QuaternaryColor,
		NavbarImageURL:  business.NavbarImageURL,
		CustomDomain: func() string {
			if business.CustomDomain != nil {
				return *business.CustomDomain
			}
			return ""
		}(),
		IsActive:           business.IsActive,
		EnableDelivery:     business.EnableDelivery,
		EnablePickup:       business.EnablePickup,
		EnableReservations: business.EnableReservations,
	}, nil
}

func (r *Repository) GetRoleByID(ctx context.Context, id uint) (*domain.Role, error) {
	var roleModel models.Role
	if err := r.database.Conn(ctx).
		Preload("Scope").
		Preload("BusinessType").
		Where("id = ?", id).
		First(&roleModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error().Uint("role_id", id).Err(err).Msg("Error al obtener rol por ID")
		return nil, err
	}

	businessTypeID := uint(0)
	businessTypeName := ""
	if roleModel.BusinessTypeID != nil {
		businessTypeID = *roleModel.BusinessTypeID
	}
	if roleModel.BusinessType != nil {
		businessTypeName = roleModel.BusinessType.Name
	}

	return &domain.Role{
		ID:               roleModel.ID,
		Name:             roleModel.Name,
		Description:      roleModel.Description,
		Level:            roleModel.Level,
		IsSystem:         roleModel.IsSystem,
		ScopeID:          roleModel.ScopeID,
		ScopeName:        roleModel.Scope.Name,
		ScopeCode:        roleModel.Scope.Code,
		BusinessTypeID:   businessTypeID,
		BusinessTypeName: businessTypeName,
		CreatedAt:        roleModel.CreatedAt,
		UpdatedAt:        roleModel.UpdatedAt,
	}, nil
}
