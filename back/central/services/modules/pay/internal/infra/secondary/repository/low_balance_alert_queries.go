package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
)

func (r *Repository) ListLowBalanceAlertCandidates(ctx context.Context, threshold float64) ([]entities.LowBalanceAlertCandidate, error) {
	var rows []struct {
		BusinessID     uint
		BusinessName   string
		Balance        float64
		WarehousePhone string
	}

	err := r.db.Conn(ctx).
		Table("wallet w").
		Select(`w.business_id, b.name AS business_name, w.balance,
			COALESCE(wh.phone, '') AS warehouse_phone`).
		Joins("INNER JOIN business b ON b.id = w.business_id AND b.deleted_at IS NULL").
		Joins(`LEFT JOIN warehouses wh ON wh.business_id = w.business_id AND wh.is_default = true AND wh.deleted_at IS NULL`).
		Where("w.balance <= ?", threshold).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]entities.LowBalanceAlertCandidate, len(rows))
	for i, row := range rows {
		result[i] = entities.LowBalanceAlertCandidate{
			BusinessID:     row.BusinessID,
			BusinessName:   row.BusinessName,
			Balance:        row.Balance,
			WarehousePhone: row.WarehousePhone,
		}
	}
	return result, nil
}

func (r *Repository) GetLowBalanceAlertCandidate(ctx context.Context, businessID uint, threshold float64) (*entities.LowBalanceAlertCandidate, error) {
	var row struct {
		BusinessID     uint
		BusinessName   string
		Balance        float64
		WarehousePhone string
	}

	err := r.db.Conn(ctx).
		Table("wallet w").
		Select(`w.business_id, b.name AS business_name, w.balance,
			COALESCE(wh.phone, '') AS warehouse_phone`).
		Joins("INNER JOIN business b ON b.id = w.business_id AND b.deleted_at IS NULL").
		Joins(`LEFT JOIN warehouses wh ON wh.business_id = w.business_id AND wh.is_default = true AND wh.deleted_at IS NULL`).
		Where("w.business_id = ? AND w.balance <= ?", businessID, threshold).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.BusinessID == 0 {
		return nil, nil
	}

	return &entities.LowBalanceAlertCandidate{
		BusinessID:     row.BusinessID,
		BusinessName:   row.BusinessName,
		Balance:        row.Balance,
		WarehousePhone: row.WarehousePhone,
	}, nil
}

func (r *Repository) HasRecentSystemAlert(ctx context.Context, phoneNumber string) (bool, error) {
	var count int64
	err := r.db.Conn(ctx).
		Table("whatsapp_conversations").
		Where("phone_number = ? AND conversation_type = 'system_alert' AND expires_at > now()", phoneNumber).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
