package repository

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) ListBusinessesForAdmin(ctx context.Context, page, pageSize int, search, statusFilter string) ([]entities.AdminBusinessRow, int64, error) {
	applyFilters := func(tx *gorm.DB) *gorm.DB {
		tx = tx.Model(&models.Business{}).Where("deleted_at IS NULL")

		if search != "" {
			like := "%" + search + "%"
			tx = tx.Where("name ILIKE ? OR code ILIKE ?", like, like)
		}

		now := time.Now()
		switch statusFilter {
		case entities.AdminStatusFilterActive:
			tx = tx.Where("subscription_status = ?", entities.BusinessStatusActive)
		case entities.AdminStatusFilterExpiringSoon:
			tx = tx.Where("subscription_status = ? AND subscription_end_date BETWEEN ? AND ?",
				entities.BusinessStatusActive, now, now.AddDate(0, 0, 5))
		case entities.AdminStatusFilterExpired:
			tx = tx.Where("subscription_status = ?", entities.BusinessStatusExpired)
		case entities.AdminStatusFilterCancelled:
			tx = tx.Where("subscription_status = ?", entities.BusinessStatusCancelled)
		case entities.AdminStatusFilterNoPlan:
			tx = tx.Where("subscription_type_id IS NULL")
		}
		return tx
	}

	var total int64
	if err := applyFilters(r.db.Conn(ctx)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var businesses []models.Business
	offset := (page - 1) * pageSize
	if err := applyFilters(r.db.Conn(ctx)).
		Preload("SubscriptionType").
		Order("created_at DESC").
		Limit(pageSize).Offset(offset).
		Find(&businesses).Error; err != nil {
		return nil, 0, err
	}

	if len(businesses) == 0 {
		return []entities.AdminBusinessRow{}, total, nil
	}

	ids := make([]uint, len(businesses))
	for i, b := range businesses {
		ids[i] = b.ID
	}

	var subs []models.BusinessSubscription
	if err := r.db.Conn(ctx).
		Where("business_id IN ? AND status = ?", ids, entities.SubscriptionStatusPaid).
		Order("created_at DESC").
		Find(&subs).Error; err != nil {
		return nil, 0, err
	}

	latestByBusiness := make(map[uint]models.BusinessSubscription, len(ids))
	for _, s := range subs {
		if _, ok := latestByBusiness[s.BusinessID]; !ok {
			latestByBusiness[s.BusinessID] = s
		}
	}

	rows := make([]entities.AdminBusinessRow, len(businesses))
	for i, b := range businesses {
		row := entities.AdminBusinessRow{
			ID:     b.ID,
			Name:   b.Name,
			Code:   b.Code,
			Status: b.SubscriptionStatus,
		}
		if b.SubscriptionType != nil {
			row.PlanName = b.SubscriptionType.Name
		}
		if latest, ok := latestByBusiness[b.ID]; ok {
			start := latest.StartDate
			end := latest.EndDate
			row.CycleStartDate = &start
			row.CycleEndDate = &end
			amount := latest.Amount
			created := latest.CreatedAt
			row.LastPaymentAmount = &amount
			row.LastPaymentDate = &created

			if b.SubscriptionType != nil {
				forecast, err := r.forecastNextPayment(ctx, b.ID, *b.SubscriptionType, start, end)
				if err != nil {
					return nil, 0, err
				}
				row.ForecastedPayment = &forecast
			}
		}
		rows[i] = row
	}

	return rows, total, nil
}

func (r *Repository) forecastNextPayment(ctx context.Context, businessID uint, plan models.SubscriptionType, cycleStart, cycleEnd time.Time) (float64, error) {
	forecast := plan.Price

	if plan.IncludedShipments != nil && plan.ShipmentOveragePrice != nil {
		count, err := r.countShipmentsInRange(ctx, businessID, cycleStart, cycleEnd)
		if err != nil {
			return 0, err
		}
		if extra := count - int64(*plan.IncludedShipments); extra > 0 {
			forecast += float64(extra) * *plan.ShipmentOveragePrice
		}
	}

	if plan.IncludedInvoices != nil && plan.InvoiceOveragePrice != nil {
		count, err := r.countInvoicesInRange(ctx, businessID, cycleStart, cycleEnd)
		if err != nil {
			return 0, err
		}
		if extra := count - int64(*plan.IncludedInvoices); extra > 0 {
			forecast += float64(extra) * *plan.InvoiceOveragePrice
		}
	}

	if plan.IncludedOrders != nil && plan.OrderOveragePrice != nil {
		count, err := r.countOrdersInRange(ctx, businessID, cycleStart, cycleEnd)
		if err != nil {
			return 0, err
		}
		if extra := count - int64(*plan.IncludedOrders); extra > 0 {
			forecast += float64(extra) * *plan.OrderOveragePrice
		}
	}

	return forecast, nil
}

func (r *Repository) countShipmentsInRange(ctx context.Context, businessID uint, start, end time.Time) (int64, error) {
	var count int64
	err := r.db.Conn(ctx).Table("shipments").
		Joins("JOIN orders ON orders.id = shipments.order_id").
		Where("orders.business_id = ? AND shipments.deleted_at IS NULL AND shipments.created_at BETWEEN ? AND ?", businessID, start, end).
		Count(&count).Error
	return count, err
}

func (r *Repository) countInvoicesInRange(ctx context.Context, businessID uint, start, end time.Time) (int64, error) {
	var count int64
	err := r.db.Conn(ctx).Table("invoices").
		Where("business_id = ? AND deleted_at IS NULL AND created_at BETWEEN ? AND ?", businessID, start, end).
		Count(&count).Error
	return count, err
}

func (r *Repository) countOrdersInRange(ctx context.Context, businessID uint, start, end time.Time) (int64, error) {
	var count int64
	err := r.db.Conn(ctx).Table("orders").
		Where("business_id = ? AND deleted_at IS NULL AND created_at BETWEEN ? AND ?", businessID, start, end).
		Count(&count).Error
	return count, err
}

func (r *Repository) GetAdminKPIs(ctx context.Context) (entities.AdminKPIs, error) {
	var kpis entities.AdminKPIs

	var activeBusinesses []models.Business
	if err := r.db.Conn(ctx).Model(&models.Business{}).
		Where("deleted_at IS NULL AND subscription_status = ?", entities.BusinessStatusActive).
		Preload("SubscriptionType").
		Find(&activeBusinesses).Error; err != nil {
		return kpis, err
	}
	kpis.ActiveCount = int64(len(activeBusinesses))
	for _, b := range activeBusinesses {
		if b.SubscriptionType != nil {
			kpis.MRR += b.SubscriptionType.Price
		}
	}

	now := time.Now()
	if err := r.db.Conn(ctx).Model(&models.Business{}).
		Where("deleted_at IS NULL AND subscription_status = ? AND subscription_end_date BETWEEN ? AND ?",
			entities.BusinessStatusActive, now, now.AddDate(0, 0, 5)).
		Count(&kpis.ExpiringSoonCount).Error; err != nil {
		return kpis, err
	}

	if err := r.db.Conn(ctx).Model(&models.Business{}).
		Where("deleted_at IS NULL AND subscription_status IN ?", []string{entities.BusinessStatusExpired, entities.BusinessStatusCancelled}).
		Count(&kpis.ExpiredOrSuspendedCount).Error; err != nil {
		return kpis, err
	}

	return kpis, nil
}
