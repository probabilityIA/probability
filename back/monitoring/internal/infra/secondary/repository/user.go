package repository

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/monitoring/internal/domain/entities"
	domainErrors "github.com/secamc93/probability/back/monitoring/internal/domain/errors"
	"gorm.io/gorm"
)

type userModel struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string
	Email    string
	Password string
	IsActive bool
	ScopeID  *uint
}

func (userModel) TableName() string {
	return "user"
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entities.MonitoringUser, string, error) {
	var user userModel

	err := r.db.WithContext(ctx).
		Where("LOWER(email) = ? AND deleted_at IS NULL", strings.ToLower(email)).
		First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, "", domainErrors.ErrInvalidCredentials
		}
		return nil, "", err
	}

	return &entities.MonitoringUser{
		ID:       user.ID,
		Email:    user.Email,
		Name:     user.Name,
		ScopeID:  user.ScopeID,
		IsActive: user.IsActive,
	}, user.Password, nil
}
