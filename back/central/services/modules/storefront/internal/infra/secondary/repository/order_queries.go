package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) ListOrdersByUserID(ctx context.Context, businessID, userID uint, page, pageSize int) ([]entities.StorefrontOrder, int64, error) {
	var orders []models.Order
	var total int64

	query := r.db.Conn(ctx).Model(&models.Order{}).
		Where("business_id = ? AND user_id = ? AND deleted_at IS NULL", businessID, userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.
		Preload("OrderItems.Product").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return mappers.OrdersToEntities(orders), total, nil
}

func (r *Repository) GetOrderByIDAndUserID(ctx context.Context, orderID string, businessID, userID uint) (*entities.StorefrontOrder, error) {
	var order models.Order
	err := r.db.Conn(ctx).
		Preload("OrderItems.Product").
		Where("id = ? AND business_id = ? AND user_id = ? AND deleted_at IS NULL", orderID, businessID, userID).
		First(&order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrOrderNotFound
		}
		return nil, err
	}
	return mappers.OrderToEntity(&order), nil
}
