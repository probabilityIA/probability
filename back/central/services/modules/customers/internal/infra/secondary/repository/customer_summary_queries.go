package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) GetCustomerSummary(ctx context.Context, businessID, customerID uint) (*entities.CustomerSummary, error) {
	var model models.CustomerSummary
	err := r.db.Conn(ctx).
		Where("customer_id = ? AND business_id = ?", customerID, businessID).
		First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return mapCustomerSummaryToEntity(&model), nil
}
