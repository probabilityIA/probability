package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/auth/users/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
	var user domain.UserAuthInfo
	if err := r.database.Conn(ctx).
		Model(&models.User{}).
		Select("id, name, email, password, phone, avatar_url, is_active, last_login_at, created_at, updated_at, deleted_at").
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

func (r *Repository) GetUserByID(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
	var user domain.UserAuthInfo
	if err := r.database.Conn(ctx).
		Model(&models.User{}).
		Select("id, name, email, password, phone, avatar_url, is_active, last_login_at, created_at, updated_at, deleted_at").
		Where("id = ?", userID).
		First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error().Uint("user_id", userID).Err(err).Msg("Error al obtener usuario por ID")
		return nil, err
	}
	return &user, nil
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

func (r *Repository) AssignBusinessStaffRelationships(ctx context.Context, userID uint, assignments []domain.BusinessRoleAssignment) error {
	db := r.database.Conn(ctx)

	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		r.logger.Error().Uint("user_id", userID).Err(err).Msg("Error al encontrar usuario")
		return err
	}

	if len(assignments) > 0 {
		businessSet := make(map[uint]struct{})
		roleSet := make(map[uint]struct{})
		for _, assignment := range assignments {
			businessSet[assignment.BusinessID] = struct{}{}
			roleSet[assignment.RoleID] = struct{}{}
		}
		uniqueBusinessIDs := make([]uint, 0, len(businessSet))
		for id := range businessSet {
			uniqueBusinessIDs = append(uniqueBusinessIDs, id)
		}
		uniqueRoleIDs := make([]uint, 0, len(roleSet))
		for id := range roleSet {
			uniqueRoleIDs = append(uniqueRoleIDs, id)
		}

		var businessCount int64
		if err := db.Model(&models.Business{}).Where("id IN ?", uniqueBusinessIDs).Count(&businessCount).Error; err != nil {
			r.logger.Error().Uint("user_id", userID).Err(err).Msg("Error al verificar businesses")
			return err
		}
		if businessCount != int64(len(uniqueBusinessIDs)) {
			return fmt.Errorf("algunos businesses no existen")
		}

		var roleCount int64
		if err := db.Model(&models.Role{}).Where("id IN ?", uniqueRoleIDs).Count(&roleCount).Error; err != nil {
			r.logger.Error().Uint("user_id", userID).Err(err).Msg("Error al verificar roles")
			return err
		}
		if roleCount != int64(len(uniqueRoleIDs)) {
			return fmt.Errorf("algunos roles no existen")
		}
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("business_staff").Where("user_id = ?", userID).Delete(nil).Error; err != nil {
			return err
		}
		if err := tx.Table("user_businesses").Where("user_id = ?", userID).Delete(nil).Error; err != nil {
			return err
		}
		if err := tx.Table("user_roles").Where("user_id = ?", userID).Delete(nil).Error; err != nil {
			return err
		}

		if len(assignments) > 0 {
			businessStaffRecords := make([]models.BusinessStaff, len(assignments))
			businessSet := make(map[uint]bool)
			roleSet := make(map[uint]bool)

			for i, assignment := range assignments {
				rid := assignment.RoleID
				bid := assignment.BusinessID
				businessStaffRecords[i] = models.BusinessStaff{
					UserID:     userID,
					BusinessID: &bid,
					RoleID:     &rid,
				}
				businessSet[assignment.BusinessID] = true
				roleSet[assignment.RoleID] = true
			}

			if err := tx.CreateInBatches(businessStaffRecords, 100).Error; err != nil {
				return err
			}

			for businessID := range businessSet {
				if err := tx.Table("user_businesses").Create(map[string]interface{}{
					"user_id":     userID,
					"business_id": businessID,
				}).Error; err != nil {
					r.logger.Warn().Err(err).Msg("Error al insertar en user_businesses")
				}
			}

			for roleID := range roleSet {
				if err := tx.Table("user_roles").Create(map[string]interface{}{
					"user_id": userID,
					"role_id": roleID,
				}).Error; err != nil {
					r.logger.Warn().Err(err).Msg("Error al insertar en user_roles")
				}
			}
		}

		return nil
	})
}
func (r *Repository) GetBusinessStaffRelationships(ctx context.Context, userID uint) ([]domain.BusinessRoleAssignmentDetailed, error) {
	var businessStaffList []models.BusinessStaff

	err := r.database.Conn(ctx).
		Preload("Business.BusinessType").
		Preload("Role").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Find(&businessStaffList).Error

	if err != nil {
		r.logger.Error().Uint("user_id", userID).Err(err).Msg("Error al obtener relaciones business_staff")
		return nil, err
	}

	r.logger.Info().
		Uint("user_id", userID).
		Int("business_staff_count", len(businessStaffList)).
		Msg("Relaciones business_staff obtenidas")

	assignments := make([]domain.BusinessRoleAssignmentDetailed, 0, len(businessStaffList))
	for _, bs := range businessStaffList {
		businessName := ""
		var businessID uint
		if bs.BusinessID != nil {
			businessID = *bs.BusinessID
			if bs.Business.ID != 0 {
				businessName = bs.Business.Name
			}
		}

		var roleID uint
		roleName := ""
		if bs.RoleID != nil {
			roleID = *bs.RoleID
			r.logger.Debug().
				Uint("user_id", userID).
				Uint("business_id", businessID).
				Uint("role_id", roleID).
				Bool("role_loaded", bs.Role.ID != 0).
				Uint("role_model_id", bs.Role.ID).
				Msg("Verificando carga del rol")

			// Verificar si el Role fue cargado correctamente desde preload
			if bs.Role.ID != 0 && bs.Role.ID == roleID {
				roleName = bs.Role.Name
				r.logger.Debug().
					Uint("user_id", userID).
					Uint("business_id", businessID).
					Uint("role_id", roleID).
					Str("role_name", roleName).
					Msg("Rol cargado correctamente desde preload")
			} else {
				// Si el preload no funcionó, consultar el rol manualmente
				r.logger.Warn().
					Uint("user_id", userID).
					Uint("business_id", businessID).
					Uint("role_id", roleID).
					Uint("role_model_id", bs.Role.ID).
					Msg("El rol no fue cargado desde preload, consultando manualmente")

				var role models.Role
				if err := r.database.Conn(ctx).Where("id = ?", roleID).First(&role).Error; err == nil && role.ID != 0 {
					roleName = role.Name
					r.logger.Info().
						Uint("user_id", userID).
						Uint("business_id", businessID).
						Uint("role_id", roleID).
						Str("role_name", roleName).
						Msg("Rol obtenido manualmente")
				} else {
					r.logger.Error().
						Uint("user_id", userID).
						Uint("business_id", businessID).
						Uint("role_id", roleID).
						Err(err).
						Msg("Error al consultar rol manualmente")
				}
			}
		} else {
			r.logger.Debug().
				Uint("user_id", userID).
				Uint("business_id", businessID).
				Msg("No hay role_id asignado en business_staff")
		}

		assignments = append(assignments, domain.BusinessRoleAssignmentDetailed{
			BusinessID:   businessID,
			BusinessName: businessName,
			RoleID:       roleID,
			RoleName:     roleName,
		})
	}

	r.logger.Info().
		Uint("user_id", userID).
		Int("assignments_count", len(assignments)).
		Msg("Assignments construidos desde business_staff")

	return assignments, nil
}

func (r *Repository) AssignRoleToUserBusiness(ctx context.Context, userID uint, assignments []domain.BusinessRoleAssignment) error {
	db := r.database.Conn(ctx)

	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("usuario no encontrado")
	}

	if len(assignments) == 0 {
		return fmt.Errorf("no se proporcionaron asignaciones")
	}

	businessSet := make(map[uint]struct{})
	roleSet := make(map[uint]struct{})
	for _, assignment := range assignments {
		businessSet[assignment.BusinessID] = struct{}{}
		roleSet[assignment.RoleID] = struct{}{}
	}

	businessIDs := make([]uint, 0, len(businessSet))
	for id := range businessSet {
		businessIDs = append(businessIDs, id)
	}

	roleIDs := make([]uint, 0, len(roleSet))
	for id := range roleSet {
		roleIDs = append(roleIDs, id)
	}

	var businesses []models.Business
	if err := db.Preload("BusinessType").Where("id IN ?", businessIDs).Find(&businesses).Error; err != nil {
		return fmt.Errorf("error al verificar businesses")
	}

	if len(businesses) != len(businessIDs) {
		return fmt.Errorf("algunos businesses no existen")
	}

	businessMap := make(map[uint]models.Business)
	for _, b := range businesses {
		businessMap[b.ID] = b
	}

	var roles []models.Role
	if err := db.Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
		return fmt.Errorf("error al verificar roles")
	}

	if len(roles) != len(roleIDs) {
		return fmt.Errorf("algunos roles no existen")
	}

	roleMap := make(map[uint]models.Role)
	for _, role := range roles {
		roleMap[role.ID] = role
	}

	for _, assignment := range assignments {
		business, businessExists := businessMap[assignment.BusinessID]
		if !businessExists {
			return fmt.Errorf("business %d no encontrado", assignment.BusinessID)
		}

		role, roleExists := roleMap[assignment.RoleID]
		if !roleExists {
			return fmt.Errorf("rol %d no encontrado", assignment.RoleID)
		}

		if role.BusinessTypeID == nil || *role.BusinessTypeID != business.BusinessTypeID {
			return fmt.Errorf("el rol %d no corresponde al tipo de business del business %d", assignment.RoleID, assignment.BusinessID)
		}
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for _, assignment := range assignments {
			var existingBS models.BusinessStaff
			if err := tx.Where("user_id = ? AND business_id = ?", userID, assignment.BusinessID).First(&existingBS).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return fmt.Errorf("el usuario no está asociado al business %d", assignment.BusinessID)
				}
				return fmt.Errorf("error al verificar relación usuario-business")
			}

			if err := tx.Model(&existingBS).Update("role_id", assignment.RoleID).Error; err != nil {
				return fmt.Errorf("error al asignar rol %d al business %d", assignment.RoleID, assignment.BusinessID)
			}
		}
		return nil
	})
}

func (r *Repository) AssignBusinessesToUser(ctx context.Context, userID uint, businessIDs []uint) error {
	db := r.database.Conn(ctx)

	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return err
	}

	if len(businessIDs) > 0 {
		var count int64
		if err := db.Model(&models.Business{}).Where("id IN ?", businessIDs).Count(&count).Error; err != nil {
			return err
		}
		if count != int64(len(businessIDs)) {
			return fmt.Errorf("algunos businesses no existen")
		}
	}

	return db.Transaction(func(tx *gorm.DB) error {
		var existingRecords []models.BusinessStaff
		if err := tx.Unscoped().Where("user_id = ? AND business_id IS NOT NULL", userID).Find(&existingRecords).Error; err != nil {
			return err
		}

		existingBusinessMap := make(map[uint]bool)
		existingRecordMap := make(map[uint]*models.BusinessStaff)
		for i := range existingRecords {
			if existingRecords[i].BusinessID != nil {
				bid := *existingRecords[i].BusinessID
				existingBusinessMap[bid] = true
				existingRecordMap[bid] = &existingRecords[i]
			}
		}

		newBusinessMap := make(map[uint]bool)
		for _, bid := range businessIDs {
			newBusinessMap[bid] = true
		}

		for bid := range existingBusinessMap {
			if !newBusinessMap[bid] {
				if err := tx.Unscoped().Where("user_id = ? AND business_id = ?", userID, bid).Delete(&models.BusinessStaff{}).Error; err != nil {
					return err
				}
			}
		}

		if len(businessIDs) > 0 {
			for _, bid := range businessIDs {
				if existingRecord, exists := existingRecordMap[bid]; exists {
					if err := tx.Unscoped().Model(existingRecord).
						Updates(map[string]interface{}{
							"deleted_at": nil,
							"role_id":    nil,
						}).Error; err != nil {
						return err
					}
				} else {
					b := bid
					newRecord := models.BusinessStaff{
						UserID:     userID,
						BusinessID: &b,
						RoleID:     nil,
					}
					if err := tx.Create(&newRecord).Error; err != nil {
						return err
					}
				}
			}
		}

		return nil
	})
}
