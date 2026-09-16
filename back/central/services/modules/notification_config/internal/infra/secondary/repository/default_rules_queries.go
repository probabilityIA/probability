package repository

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"gorm.io/gorm"
)

type defaultRulesQuerier struct {
	db db.IDatabase
}

func NewDefaultRulesQuerier(database db.IDatabase) ports.IDefaultRulesQuerier {
	return &defaultRulesQuerier{db: database}
}

func (q *defaultRulesQuerier) PlatformIntegrationID(ctx context.Context, businessID uint) (uint, error) {
	var row struct{ ID uint }
	err := q.db.Conn(ctx).Table("integrations").
		Select("integrations.id").
		Joins("JOIN integration_types ON integration_types.id = integrations.integration_type_id").
		Where("integrations.business_id = ? AND integration_types.code = ? AND integrations.deleted_at IS NULL", businessID, "platform").
		Order("integrations.id").
		Limit(1).
		Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	return row.ID, err
}

func (q *defaultRulesQuerier) OrderStatusIDsByCodes(ctx context.Context, codes []string) ([]uint, error) {
	var ids []uint
	err := q.db.Conn(ctx).Table("order_statuses").
		Where("code IN ? AND deleted_at IS NULL", codes).
		Order("id").
		Pluck("id", &ids).Error
	return ids, err
}
