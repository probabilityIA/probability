package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
)

func (r *Repository) GetKardex(ctx context.Context, params dtos.KardexQueryParams) ([]entities.KardexEntry, error) {
	type row struct {
		MovementID       uint
		CreatedAt        string
		MovementTypeCode string
		MovementTypeName string
		Quantity         int
		PreviousQty      int
		NewQty           int
		Reason           string
		ReferenceType    *string
		ReferenceID      *string
		LocationID       *uint
		LotID            *uint
	}
	var rows []row

	q := r.db.Conn(ctx).
		Table("stock_movements sm").
		Select(`sm.id AS movement_id,
			sm.created_at,
			smt.code AS movement_type_code,
			smt.name AS movement_type_name,
			sm.quantity,
			sm.previous_qty,
			sm.new_qty,
			sm.reason,
			sm.reference_type,
			sm.reference_id,
			sm.location_id,
			sm.lot_id`).
		Joins("INNER JOIN stock_movement_types smt ON smt.id = sm.movement_type_id").
		Where("sm.business_id = ? AND sm.product_id = ? AND sm.warehouse_id = ? AND sm.deleted_at IS NULL",
			params.BusinessID, params.ProductID, params.WarehouseID)

	if params.From != nil {
		q = q.Where("sm.created_at >= ?", *params.From)
	}
	if params.To != nil {
		q = q.Where("sm.created_at <= ?", *params.To)
	}

	if err := q.Order("sm.created_at ASC, sm.id ASC").Scan(&rows).Error; err != nil {
		return nil, err
	}

	entries := make([]entities.KardexEntry, len(rows))
	balance := 0
	for i, row := range rows {
		balance += row.Quantity
		entries[i] = entities.KardexEntry{
			MovementID:       row.MovementID,
			MovementTypeCode: row.MovementTypeCode,
			MovementTypeName: row.MovementTypeName,
			Quantity:         row.Quantity,
			PreviousQty:      row.PreviousQty,
			NewQty:           row.NewQty,
			RunningBalance:   balance,
			Reason:           row.Reason,
			ReferenceType:    row.ReferenceType,
			ReferenceID:      row.ReferenceID,
			LocationID:       row.LocationID,
			LotID:            row.LotID,
		}
	}
	return entries, nil
}
