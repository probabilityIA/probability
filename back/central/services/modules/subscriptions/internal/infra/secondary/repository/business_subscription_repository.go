package repository

import (
	"context"
	"errors"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) CreateBusinessSubscription(ctx context.Context, subscription *entities.BusinessSubscription) error {
	subTypeID := subscription.SubscriptionTypeID
	months := subscription.Months

	subDB := &models.BusinessSubscription{
		BusinessID:         subscription.BusinessID,
		SubscriptionTypeID: &subTypeID,
		Months:             &months,
		Amount:             subscription.Amount,
		StartDate:          subscription.StartDate,
		EndDate:            subscription.EndDate,
		Status:             subscription.Status,
		PaymentReference:   subscription.PaymentReference,
		Notes:              subscription.Notes,
	}

	if err := r.db.Conn(ctx).Create(subDB).Error; err != nil {
		return err
	}

	subscription.ID = subDB.ID
	subscription.CreatedAt = subDB.CreatedAt
	subscription.UpdatedAt = subDB.UpdatedAt
	return nil
}

func (r *Repository) GetLatestByBusinessID(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error) {
	var subDB models.BusinessSubscription
	err := r.db.Conn(ctx).
		Preload("SubscriptionType").
		Where("business_id = ?", businessID).
		Order("created_at desc").
		First(&subDB).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return subscriptionToEntity(&subDB), nil
}

func (r *Repository) ListByBusinessID(ctx context.Context, businessID uint) ([]entities.BusinessSubscription, error) {
	var subsDB []models.BusinessSubscription
	err := r.db.Conn(ctx).
		Preload("SubscriptionType").
		Where("business_id = ?", businessID).
		Order("created_at desc").
		Find(&subsDB).Error
	if err != nil {
		return nil, err
	}

	subs := make([]entities.BusinessSubscription, len(subsDB))
	for i, s := range subsDB {
		subs[i] = *subscriptionToEntity(&s)
	}
	return subs, nil
}

func (r *Repository) UpdateBusinessCurrentSubscriptionType(ctx context.Context, businessID uint, subscriptionTypeID uint, status string, endDate time.Time) error {
	updates := map[string]interface{}{
		"subscription_type_id":  subscriptionTypeID,
		"subscription_status":   status,
		"subscription_end_date": endDate,
	}
	return r.db.Conn(ctx).Model(&models.Business{}).Where("id = ? AND deleted_at IS NULL", businessID).Updates(updates).Error
}

func (r *Repository) UpdateBusinessSubscriptionStatus(ctx context.Context, businessID uint, status string, endDate *time.Time) error {
	updates := map[string]interface{}{"subscription_status": status}
	if endDate != nil {
		updates["subscription_end_date"] = *endDate
	}
	return r.db.Conn(ctx).Model(&models.Business{}).Where("id = ? AND deleted_at IS NULL", businessID).Updates(updates).Error
}

func (r *Repository) GetBusinessCurrentSubscriptionTypeID(ctx context.Context, businessID uint) (*uint, error) {
	var business models.Business
	err := r.db.Conn(ctx).Select("subscription_type_id").Where("id = ? AND deleted_at IS NULL", businessID).First(&business).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return business.SubscriptionTypeID, nil
}

func (r *Repository) ListBusinessesExpiringBetween(ctx context.Context, from, to time.Time) ([]uint, error) {
	var businesses []models.Business
	err := r.db.Conn(ctx).Select("id").
		Where("deleted_at IS NULL AND subscription_status = ? AND subscription_end_date BETWEEN ? AND ?", entities.BusinessStatusActive, from, to).
		Find(&businesses).Error
	if err != nil {
		return nil, err
	}

	ids := make([]uint, len(businesses))
	for i, b := range businesses {
		ids[i] = b.ID
	}
	return ids, nil
}

func (r *Repository) ListBusinessesJustExpired(ctx context.Context, before time.Time) ([]uint, error) {
	var businesses []models.Business
	err := r.db.Conn(ctx).Select("id").
		Where("deleted_at IS NULL AND subscription_status = ? AND subscription_end_date < ?", entities.BusinessStatusActive, before).
		Find(&businesses).Error
	if err != nil {
		return nil, err
	}

	ids := make([]uint, len(businesses))
	for i, b := range businesses {
		ids[i] = b.ID
	}
	return ids, nil
}
