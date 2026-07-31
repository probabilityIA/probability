package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/errors"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) GetClientByUserID(ctx context.Context, businessID, userID uint) (*entities.StorefrontClient, error) {
	var client models.Client
	err := r.db.Conn(ctx).
		Where("business_id = ? AND user_id = ? AND deleted_at IS NULL", businessID, userID).
		First(&client).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrClientNotFound
		}
		return nil, err
	}
	return &entities.StorefrontClient{
		ID:         client.ID,
		BusinessID: client.BusinessID,
		UserID:     client.UserID,
		Name:       client.Name,
		Email:      client.Email,
		Phone:      client.Phone,
		Dni:        client.Dni,
	}, nil
}

func (r *Repository) GetClientByBusinessAndEmail(ctx context.Context, businessID uint, email string) (*entities.StorefrontClient, error) {
	var client models.Client
	err := r.db.Conn(ctx).
		Where("business_id = ? AND email = ? AND deleted_at IS NULL", businessID, email).
		First(&client).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &entities.StorefrontClient{
		ID:         client.ID,
		BusinessID: client.BusinessID,
		UserID:     client.UserID,
		Name:       client.Name,
		Email:      client.Email,
		Phone:      client.Phone,
		Dni:        client.Dni,
	}, nil
}

func (r *Repository) LinkClientUser(ctx context.Context, clientID, userID uint) error {
	return r.db.Conn(ctx).Model(&models.Client{}).
		Where("id = ?", clientID).
		Update("user_id", userID).Error
}

func (r *Repository) CreateClient(ctx context.Context, client *entities.StorefrontClient) error {
	model := &models.Client{
		BusinessID: client.BusinessID,
		UserID:     client.UserID,
		Name:       client.Name,
		Email:      client.Email,
		Phone:      client.Phone,
		Dni:        client.Dni,
	}
	return r.db.Conn(ctx).Create(model).Error
}
