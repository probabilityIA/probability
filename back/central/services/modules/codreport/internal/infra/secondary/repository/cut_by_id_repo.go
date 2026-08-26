package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) CutByID(ctx context.Context, businessID uint, cutID uint) (*entities.PaymentCut, error) {
	var row models.CodPaymentCut
	err := r.db.Conn(ctx).
		Where("id = ? AND business_id = ? AND deleted_at IS NULL", cutID, businessID).
		Limit(1).
		Find(&row).Error
	if err != nil {
		return nil, err
	}
	if row.ID == 0 {
		return nil, fmt.Errorf("corte de pago no encontrado")
	}
	var byCarrier []entities.CarrierAggregate
	if row.CarrierBreakdown != "" {
		_ = json.Unmarshal([]byte(row.CarrierBreakdown), &byCarrier)
	}
	return &entities.PaymentCut{
		ID:              row.ID,
		BusinessID:      row.BusinessID,
		PeriodStart:     row.PeriodStart,
		PeriodEnd:       row.PeriodEnd,
		Status:          row.Status,
		OrdersCount:     row.OrdersCount,
		TotalCollected:  row.TotalCollected,
		TotalDiscount:   row.TotalDiscount,
		TotalNet:        row.TotalNet,
		ByCarrier:       byCarrier,
		ConfirmedBy:     row.ConfirmedBy,
		ConfirmedByName: row.ConfirmedByName,
		ConfirmedAt:     row.ConfirmedAt,
	}, nil
}

func (r *Repository) BusinessName(ctx context.Context, businessID uint) string {
	var name string
	err := r.db.Conn(ctx).
		Table("business").
		Select("name").
		Where("id = ? AND deleted_at IS NULL", businessID).
		Limit(1).
		Scan(&name).Error
	if err != nil {
		return ""
	}
	return name
}
