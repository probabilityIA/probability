package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/errors"
	"github.com/secamc93/probability/back/migration/shared/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (r *Repository) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.Conn(ctx).Model(&models.User{}).
		Where("email = ? AND deleted_at IS NULL", email).
		Count(&count).Error
	return count > 0, err
}

func (r *Repository) CreateUser(ctx context.Context, user *entities.NewUser) (uint, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("error hasheando password: %w", err)
	}

	// Business scope ID = 2
	scopeID := uint(2)
	model := &models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: string(hashedPassword),
		Phone:    user.Phone,
		IsActive: true,
		ScopeID:  &scopeID,
	}

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return 0, err
	}
	return model.ID, nil
}

func (r *Repository) CreateBusinessStaff(ctx context.Context, userID, businessID, roleID uint) error {
	model := &models.BusinessStaff{
		UserID:     userID,
		BusinessID: &businessID,
		RoleID:     &roleID,
	}
	return r.db.Conn(ctx).Create(model).Error
}

func (r *Repository) GetBusinessByCode(ctx context.Context, code string) (*entities.StorefrontBusiness, error) {
	var business models.Business
	err := r.db.Conn(ctx).
		Where("code = ? AND is_active = true AND deleted_at IS NULL", code).
		First(&business).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrBusinessNotFound
		}
		return nil, err
	}
	return &entities.StorefrontBusiness{
		ID:   business.ID,
		Name: business.Name,
		Code: business.Code,
	}, nil
}

func (r *Repository) GetPlatformIntegrationID(ctx context.Context, businessID uint) (uint, error) {
	var result struct {
		ID uint
	}
	err := r.db.Conn(ctx).
		Table("integrations").
		Select("id").
		Where("business_id = ? AND integration_type_id = ? AND is_active = true AND deleted_at IS NULL", businessID, 6).
		Limit(1).
		Scan(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, domainerrors.ErrIntegrationNotFound
		}
		return 0, err
	}
	if result.ID == 0 {
		return 0, domainerrors.ErrIntegrationNotFound
	}
	return result.ID, nil
}
