package repository

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) CreateSubscriptionType(ctx context.Context, subType *entities.SubscriptionType) error {
	typeDB := &models.SubscriptionType{
		Name:          subType.Name,
		Code:          subType.Code,
		Description:   subType.Description,
		Price:         subType.Price,
		BillingPeriod: subType.BillingPeriod,
		Active:        subType.Active,
		Features:      marshalModuleCodes(subType.ModuleCodes),
	}

	if err := r.db.Conn(ctx).Create(typeDB).Error; err != nil {
		return err
	}

	subType.ID = typeDB.ID
	subType.CreatedAt = typeDB.CreatedAt
	subType.UpdatedAt = typeDB.UpdatedAt
	return nil
}

func (r *Repository) UpdateSubscriptionType(ctx context.Context, subType *entities.SubscriptionType) error {
	updates := map[string]interface{}{
		"name":           subType.Name,
		"description":    subType.Description,
		"price":          subType.Price,
		"billing_period": subType.BillingPeriod,
		"active":         subType.Active,
		"features":       marshalModuleCodes(subType.ModuleCodes),
	}
	return r.db.Conn(ctx).Model(&models.SubscriptionType{}).Where("id = ?", subType.ID).Updates(updates).Error
}

func (r *Repository) DeleteSubscriptionType(ctx context.Context, id uint) error {
	return r.db.Conn(ctx).Delete(&models.SubscriptionType{}, id).Error
}

func (r *Repository) GetSubscriptionType(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
	var typeDB models.SubscriptionType
	err := r.db.Conn(ctx).First(&typeDB, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return subscriptionTypeToEntity(&typeDB), nil
}

func (r *Repository) ListSubscriptionTypes(ctx context.Context, activeOnly bool) ([]entities.SubscriptionType, error) {
	query := r.db.Conn(ctx).Order("price ASC")
	if activeOnly {
		query = query.Where("active = ?", true)
	}

	var typesDB []models.SubscriptionType
	if err := query.Find(&typesDB).Error; err != nil {
		return nil, err
	}

	types := make([]entities.SubscriptionType, len(typesDB))
	for i, t := range typesDB {
		types[i] = *subscriptionTypeToEntity(&t)
	}
	return types, nil
}
