package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Repository) SelectableCutOrders(ctx context.Context, f dtos.SelectableOrdersFilter) ([]entities.CodOrder, error) {
	sql := fmt.Sprintf(`
SELECT o.id AS order_id, o.order_number, o.customer_name, o.cod_total, o.currency, o.created_at,
	s.id AS shipment_id,
	UPPER(TRIM(COALESCE(NULLIF(s.carrier,''),'SIN TRANSPORTADORA'))) AS carrier,
	COALESCE(s.shipping_cost,0) AS shipping_cost,
	COALESCE(s.cod_carrier_fee,0) AS cod_carrier_fee,
	s.status, s.delivered_at,
	true AS collected,
	false AS paid,
	`+hasGuideExpr+` AS has_guide
FROM orders o %s
WHERE o.deleted_at IS NULL AND o.cod_total > 0 AND o.business_id = ?
	AND s.status = 'delivered'
	AND COALESCE(s.delivered_at, s.updated_at) BETWEEN ? AND ?
	AND NOT %s
ORDER BY COALESCE(s.delivered_at, o.created_at) DESC`, latestShipmentJoin, paidExpr)

	var rows []codOrderRow
	if err := r.db.Conn(ctx).Raw(sql, f.BusinessID, f.PeriodStart, f.PeriodEnd).Scan(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]entities.CodOrder, len(rows))
	for i := range rows {
		out[i] = entities.CodOrder{
			OrderID:       rows[i].OrderID,
			OrderNumber:   rows[i].OrderNumber,
			ShipmentID:    rows[i].ShipmentID,
			HasGuide:      rows[i].HasGuide,
			CustomerName:  rows[i].CustomerName,
			Carrier:       rows[i].Carrier,
			CodTotal:      rows[i].CodTotal,
			CodCarrierFee: rows[i].CodCarrierFee,
			ShippingCost:  rows[i].ShippingCost,
			Currency:      rows[i].Currency,
			Status:        rows[i].Status,
			DeliveredAt:   rows[i].DeliveredAt,
		}
	}
	return out, nil
}

type payoutRow struct {
	OrderID   string
	Carrier   string
	CodAmount float64
}

func (r *Repository) PayoutOrders(ctx context.Context, businessID uint, orderIDs []string) ([]entities.PayoutOrder, error) {
	if len(orderIDs) == 0 {
		return []entities.PayoutOrder{}, nil
	}

	sql := fmt.Sprintf(`
SELECT o.id AS order_id,
	UPPER(TRIM(COALESCE(NULLIF(s.carrier,''),'SIN TRANSPORTADORA'))) AS carrier,
	o.cod_total AS cod_amount
FROM orders o %s
WHERE o.deleted_at IS NULL AND o.cod_total > 0 AND o.business_id = ?
	AND s.status = 'delivered'
	AND o.id IN ?
	AND NOT %s`, latestShipmentJoin, paidExpr)

	var rows []payoutRow
	if err := r.db.Conn(ctx).Raw(sql, businessID, orderIDs).Scan(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]entities.PayoutOrder, len(rows))
	for i := range rows {
		out[i] = entities.PayoutOrder{
			OrderID:   rows[i].OrderID,
			Carrier:   rows[i].Carrier,
			CodAmount: rows[i].CodAmount,
		}
	}
	return out, nil
}

func (r *Repository) UpsertCutOrders(ctx context.Context, cut entities.PaymentCut, orders []entities.PayoutOrder, userID uint, userName string) (uint, error) {
	now := time.Now().UTC()
	var cutID uint

	err := r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		row := models.CodPaymentCut{
			BusinessID:      cut.BusinessID,
			PeriodStart:     cut.PeriodStart,
			PeriodEnd:       cut.PeriodEnd,
			Status:          "confirmed",
			ConfirmedBy:     userID,
			ConfirmedByName: userName,
			ConfirmedAt:     &now,
		}

		var existing models.CodPaymentCut
		e := tx.Where("business_id = ? AND period_start = ? AND period_end = ?", cut.BusinessID, cut.PeriodStart, cut.PeriodEnd).First(&existing).Error
		if e == nil {
			row.ID = existing.ID
			row.CreatedAt = existing.CreatedAt
			if err := tx.Save(&row).Error; err != nil {
				return err
			}
		} else if err := tx.Create(&row).Error; err != nil {
			return err
		}

		for i := range orders {
			link := models.CodPaymentCutOrder{
				CodPaymentCutID: row.ID,
				BusinessID:      cut.BusinessID,
				OrderID:         orders[i].OrderID,
				Carrier:         orders[i].Carrier,
				CodAmount:       orders[i].CodAmount,
				PaidAt:          now,
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "order_id"}},
				DoNothing: true,
			}).Create(&link).Error; err != nil {
				return err
			}
		}

		cutID = row.ID
		return nil
	})
	if err != nil {
		return 0, err
	}
	return cutID, nil
}

func (r *Repository) PaidAggregatesForCut(ctx context.Context, cutID uint) ([]entities.CarrierAggregate, error) {
	var rows []carrierAggRow
	sql := `
SELECT carrier, COUNT(*) AS orders_count, COALESCE(SUM(cod_amount),0) AS total_collected
FROM cod_payment_cut_orders
WHERE cod_payment_cut_id = ? AND deleted_at IS NULL
GROUP BY carrier
ORDER BY total_collected DESC`
	if err := r.db.Conn(ctx).Raw(sql, cutID).Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]entities.CarrierAggregate, len(rows))
	for i := range rows {
		out[i] = entities.CarrierAggregate{
			Carrier:        rows[i].Carrier,
			OrdersCount:    rows[i].OrdersCount,
			TotalCollected: rows[i].TotalCollected,
		}
	}
	return out, nil
}

func (r *Repository) UpdateCutTotals(ctx context.Context, cut entities.PaymentCut) error {
	breakdown, _ := json.Marshal(cut.ByCarrier)
	return r.db.Conn(ctx).Model(&models.CodPaymentCut{}).
		Where("id = ?", cut.ID).
		Updates(map[string]any{
			"orders_count":      cut.OrdersCount,
			"total_collected":   cut.TotalCollected,
			"total_discount":    cut.TotalDiscount,
			"total_net":         cut.TotalNet,
			"carrier_breakdown": string(breakdown),
		}).Error
}
