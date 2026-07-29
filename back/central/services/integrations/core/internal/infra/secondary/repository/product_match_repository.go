package repository

import (
	"context"

	"gorm.io/datatypes"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) UpdateProductMatchRules(ctx context.Context, id uint, rules datatypes.JSON) error {
	return r.db.Conn(ctx).
		Model(&models.Integration{}).
		Where("id = ?", id).
		Update("product_match_rules", rules).Error
}
