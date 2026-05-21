package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) ConfirmedCuts(ctx context.Context, businessID uint) ([]entities.PaymentCut, error) {
	var rows []models.CodPaymentCut
	err := r.db.Conn(ctx).
		Where("business_id = ?", businessID).
		Order("period_start DESC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]entities.PaymentCut, len(rows))
	for i := range rows {
		var byCarrier []entities.CarrierAggregate
		if rows[i].CarrierBreakdown != "" {
			_ = json.Unmarshal([]byte(rows[i].CarrierBreakdown), &byCarrier)
		}
		out[i] = entities.PaymentCut{
			ID:              rows[i].ID,
			BusinessID:      rows[i].BusinessID,
			PeriodStart:     rows[i].PeriodStart,
			PeriodEnd:       rows[i].PeriodEnd,
			Status:          rows[i].Status,
			OrdersCount:     rows[i].OrdersCount,
			TotalCollected:  rows[i].TotalCollected,
			TotalDiscount:   rows[i].TotalDiscount,
			TotalNet:        rows[i].TotalNet,
			ByCarrier:       byCarrier,
			ConfirmedBy:     rows[i].ConfirmedBy,
			ConfirmedByName: rows[i].ConfirmedByName,
			ConfirmedAt:     rows[i].ConfirmedAt,
		}
	}
	return out, nil
}

func (r *Repository) SaveConfirmedCut(ctx context.Context, cut entities.PaymentCut, userID uint, userName string) (*entities.PaymentCut, error) {
	breakdown, _ := json.Marshal(cut.ByCarrier)
	now := time.Now().UTC()
	row := models.CodPaymentCut{
		BusinessID:       cut.BusinessID,
		PeriodStart:      cut.PeriodStart,
		PeriodEnd:        cut.PeriodEnd,
		Status:           "confirmed",
		OrdersCount:      cut.OrdersCount,
		TotalCollected:   cut.TotalCollected,
		TotalDiscount:    cut.TotalDiscount,
		TotalNet:         cut.TotalNet,
		CarrierBreakdown: string(breakdown),
		ConfirmedBy:      userID,
		ConfirmedByName:  userName,
		ConfirmedAt:      &now,
	}

	var existing models.CodPaymentCut
	err := r.db.Conn(ctx).
		Where("business_id = ? AND period_start = ? AND period_end = ?", cut.BusinessID, cut.PeriodStart, cut.PeriodEnd).
		First(&existing).Error
	if err == nil {
		row.ID = existing.ID
		row.CreatedAt = existing.CreatedAt
		if err := r.db.Conn(ctx).Save(&row).Error; err != nil {
			return nil, err
		}
	} else if err := r.db.Conn(ctx).Create(&row).Error; err != nil {
		return nil, err
	}

	cut.ID = row.ID
	cut.Status = "confirmed"
	cut.ConfirmedBy = userID
	cut.ConfirmedByName = userName
	cut.ConfirmedAt = &now
	return &cut, nil
}
