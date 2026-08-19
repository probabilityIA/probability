package repository

import (
	"context"
	"errors"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	errs "github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/errors"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) CreateSubscriptionAndActivate(ctx context.Context, subscription *entities.BusinessSubscription, subscriptionTypeID uint, endDate time.Time) error {
	return r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		subTypeID := subscription.SubscriptionTypeID
		months := subscription.Months
		paymentMethod := subscription.PaymentMethod
		if paymentMethod == "" {
			paymentMethod = entities.PaymentMethodWallet
		}

		subDB := &models.BusinessSubscription{
			BusinessID:         subscription.BusinessID,
			SubscriptionTypeID: &subTypeID,
			Months:             &months,
			Amount:             subscription.Amount,
			StartDate:          subscription.StartDate,
			EndDate:            subscription.EndDate,
			Status:             subscription.Status,
			PaymentMethod:      paymentMethod,
			PaymentReference:   subscription.PaymentReference,
			Notes:              subscription.Notes,
		}
		if err := tx.Create(subDB).Error; err != nil {
			return err
		}
		subscription.ID = subDB.ID
		subscription.CreatedAt = subDB.CreatedAt
		subscription.UpdatedAt = subDB.UpdatedAt

		updates := map[string]interface{}{
			"subscription_type_id":  subscriptionTypeID,
			"subscription_status":   entities.BusinessStatusActive,
			"subscription_end_date": endDate,
		}
		return tx.Model(&models.Business{}).
			Where("id = ? AND deleted_at IS NULL", subscription.BusinessID).
			Updates(updates).Error
	})
}

func (r *Repository) GetLatestByBusinessID(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error) {
	var subDB models.BusinessSubscription
	err := r.db.Conn(ctx).
		Preload("SubscriptionType").
		Where("business_id = ? AND status = ?", businessID, entities.SubscriptionStatusPaid).
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

func (r *Repository) GetSubscriptionByID(ctx context.Context, id uint) (*entities.BusinessSubscription, error) {
	var subDB models.BusinessSubscription
	err := r.db.Conn(ctx).Preload("SubscriptionType").First(&subDB, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return subscriptionToEntity(&subDB), nil
}

func (r *Repository) RevertSubscriptionAndRecalculate(ctx context.Context, subscriptionID uint) (*entities.BusinessSubscription, error) {
	var reverted *entities.BusinessSubscription

	err := r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		var subDB models.BusinessSubscription
		if err := tx.First(&subDB, subscriptionID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errs.ErrSubscriptionNotFound
			}
			return err
		}
		if subDB.Status == entities.SubscriptionStatusReverted {
			return errs.ErrSubscriptionAlreadyReverted
		}
		if subDB.Status != entities.SubscriptionStatusPaid {
			return errs.ErrSubscriptionNotPaid
		}

		if err := tx.Model(&subDB).Update("status", entities.SubscriptionStatusReverted).Error; err != nil {
			return err
		}

		var next models.BusinessSubscription
		result := tx.Preload("SubscriptionType").
			Where("business_id = ? AND status = ?", subDB.BusinessID, entities.SubscriptionStatusPaid).
			Order("created_at desc").
			First(&next)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}

		updates := map[string]interface{}{}
		if result.RowsAffected == 0 {
			updates["subscription_type_id"] = nil
			updates["subscription_status"] = entities.BusinessStatusExpired
			updates["subscription_end_date"] = nil
		} else {
			status := entities.BusinessStatusActive
			if next.EndDate.Before(time.Now()) {
				status = entities.BusinessStatusExpired
			}
			updates["subscription_type_id"] = next.SubscriptionTypeID
			updates["subscription_status"] = status
			updates["subscription_end_date"] = next.EndDate
		}
		if err := tx.Model(&models.Business{}).
			Where("id = ? AND deleted_at IS NULL", subDB.BusinessID).
			Updates(updates).Error; err != nil {
			return err
		}

		subDB.Status = entities.SubscriptionStatusReverted
		reverted = subscriptionToEntity(&subDB)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return reverted, nil
}

func (r *Repository) UpdateBusinessSubscriptionStatus(ctx context.Context, businessID uint, status string, endDate *time.Time) error {
	updates := map[string]interface{}{"subscription_status": status}
	if endDate != nil {
		updates["subscription_end_date"] = *endDate
	}
	return r.db.Conn(ctx).Model(&models.Business{}).Where("id = ? AND deleted_at IS NULL", businessID).Updates(updates).Error
}

func (r *Repository) MarkExpiredIfStillActive(ctx context.Context, businessID uint, before time.Time) error {
	updates := map[string]interface{}{"subscription_status": entities.BusinessStatusExpired}
	return r.db.Conn(ctx).Model(&models.Business{}).
		Where("id = ? AND deleted_at IS NULL AND subscription_status = ? AND subscription_end_date < ?",
			businessID, entities.BusinessStatusActive, before).
		Updates(updates).Error
}

func (r *Repository) UpdateSubscriptionDates(ctx context.Context, id uint, startDate, endDate time.Time) error {
	return r.db.Conn(ctx).Model(&models.BusinessSubscription{}).Where("id = ?", id).
		Updates(map[string]interface{}{"start_date": startDate, "end_date": endDate}).Error
}

func (r *Repository) UpdateBusinessSubscriptionEndDate(ctx context.Context, businessID uint, endDate time.Time) error {
	return r.db.Conn(ctx).Model(&models.Business{}).Where("id = ? AND deleted_at IS NULL", businessID).
		Update("subscription_end_date", endDate).Error
}

func (r *Repository) GetBusinessCurrentSubscriptionTypeID(ctx context.Context, businessID uint) (*uint, error) {
	var business models.Business
	err := r.db.Conn(ctx).Select("subscription_type_id", "subscription_status").Where("id = ? AND deleted_at IS NULL", businessID).First(&business).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if business.SubscriptionStatus != entities.BusinessStatusActive {
		return nil, nil
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
