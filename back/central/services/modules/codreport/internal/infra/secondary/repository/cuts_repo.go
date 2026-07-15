package repository

import (
	"context"
	"encoding/json"

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

