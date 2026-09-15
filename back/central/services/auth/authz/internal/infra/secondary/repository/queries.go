package repository

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/auth/authz/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) GetStaffRole(ctx context.Context, userID uint, businessID *uint) (*entities.StaffRole, error) {
	var staff models.BusinessStaff
	q := r.db.Conn(ctx).Where("user_id = ?", userID)
	if businessID == nil {
		q = q.Where("business_id IS NULL")
	} else {
		q = q.Where("business_id = ?", *businessID)
	}
	if err := q.Preload("Business").Preload("Role.Scope").First(&staff).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	result := &entities.StaffRole{
		BusinessID: staff.BusinessID,
		RoleID:     staff.RoleID,
	}
	if staff.BusinessID != nil {
		result.BusinessName = staff.Business.Name
		result.SubscriptionStatus = staff.Business.SubscriptionStatus
	}
	if staff.RoleID != nil {
		result.RoleName = staff.Role.Name
		result.RoleScopeCode = staff.Role.Scope.Code
	}
	return result, nil
}

func (r *Repository) GetRolePermissions(ctx context.Context, roleID uint) ([]entities.RolePermission, error) {
	var rows []struct {
		ResourceName string
		ActionName   string
	}
	err := r.db.Conn(ctx).
		Model(&models.Permission{}).
		Select("resource.name AS resource_name, action.name AS action_name").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permission.id").
		Joins("JOIN resource ON resource.id = permission.resource_id AND resource.deleted_at IS NULL").
		Joins("JOIN action ON action.id = permission.action_id AND action.deleted_at IS NULL").
		Where("role_permissions.role_id = ?", roleID).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]entities.RolePermission, len(rows))
	for i, row := range rows {
		result[i] = entities.RolePermission{ResourceName: row.ResourceName, ActionName: row.ActionName}
	}
	return result, nil
}

func (r *Repository) GetBusinessResourceConfig(ctx context.Context, businessID uint) ([]entities.ResourceConfig, error) {
	var rows []struct {
		ResourceName string
		Active       bool
	}
	err := r.db.Conn(ctx).
		Model(&models.BusinessResourceConfigured{}).
		Select("resource.name AS resource_name, business_resource_configured.active AS active").
		Joins("JOIN resource ON resource.id = business_resource_configured.resource_id AND resource.deleted_at IS NULL").
		Where("business_resource_configured.business_id = ?", businessID).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]entities.ResourceConfig, len(rows))
	for i, row := range rows {
		result[i] = entities.ResourceConfig{ResourceName: row.ResourceName, Active: row.Active}
	}
	return result, nil
}

func (r *Repository) GetBusinessName(ctx context.Context, businessID uint) (string, error) {
	var business models.Business
	err := r.db.Conn(ctx).Select("id", "name").Where("id = ?", businessID).First(&business).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	return business.Name, nil
}
