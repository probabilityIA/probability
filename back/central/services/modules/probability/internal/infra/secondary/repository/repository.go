package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/probability/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/probability/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) GetOrderForScoring(ctx context.Context, orderID string) (*entities.ScoreOrder, error) {
	var order models.Order
	err := r.db.Conn(ctx).
		Preload("Payments").
		Preload("Addresses").
		Preload("ChannelMetadata").
		First(&order, "id = ?", orderID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("order %s not found", orderID)
		}
		return nil, err
	}
	return mappers.OrderToScoreOrder(&order), nil
}

func (r *Repository) CountOrdersByClientID(ctx context.Context, clientID uint) (int64, error) {
	var count int64
	err := r.db.Conn(ctx).Model(&models.Order{}).
		Where("customer_id = ?", clientID).
		Where("deleted_at IS NULL").
		Count(&count).Error
	return count, err
}

func (r *Repository) UpdateOrderScore(ctx context.Context, orderID string, score float64, factors []byte) error {
	return r.db.Conn(ctx).Model(&models.Order{}).
		Where("id = ?", orderID).
		Updates(map[string]interface{}{
			"delivery_probability": score,
			"negative_factors":     factors,
		}).Error
}
