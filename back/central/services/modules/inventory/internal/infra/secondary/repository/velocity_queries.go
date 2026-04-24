package repository

import (
	"context"
	"sort"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) ComputeVelocities(ctx context.Context, businessID, warehouseID uint, period string) error {
	db := r.db.Conn(ctx)

	type row struct {
		ProductID  string
		UnitsMoved int
	}
	var results []row

	days := periodToDays(period)
	if days <= 0 {
		days = 30
	}
	since := time.Now().AddDate(0, 0, -days)

	err := db.Table("stock_movements").
		Select("product_id, SUM(ABS(quantity)) AS units_moved").
		Where("business_id = ? AND warehouse_id = ? AND created_at >= ? AND deleted_at IS NULL", businessID, warehouseID, since).
		Group("product_id").
		Order("units_moved DESC").
		Scan(&results).Error
	if err != nil {
		return err
	}
	if len(results) == 0 {
		return nil
	}

	sort.SliceStable(results, func(i, j int) bool {
		return results[i].UnitsMoved > results[j].UnitsMoved
	})

	total := 0
	for _, r := range results {
		total += r.UnitsMoved
	}
	accum := 0
	computedAt := time.Now()

	for _, row := range results {
		accum += row.UnitsMoved
		pct := float64(accum) / float64(total) * 100
		rank := "C"
		if pct <= 80 {
			rank = "A"
		} else if pct <= 95 {
			rank = "B"
		}

		var existing models.ProductVelocity
		q := db.Where("business_id = ? AND product_id = ? AND warehouse_id = ? AND period = ?", businessID, row.ProductID, warehouseID, period)
		if err := q.First(&existing).Error; err == nil {
			existing.UnitsMoved = row.UnitsMoved
			existing.Rank = rank
			existing.ComputedAt = computedAt
			if err := db.Save(&existing).Error; err != nil {
				return err
			}
			continue
		}

		m := models.ProductVelocity{
			BusinessID:  businessID,
			ProductID:   row.ProductID,
			WarehouseID: warehouseID,
			Period:      period,
			UnitsMoved:  row.UnitsMoved,
			Rank:        rank,
			ComputedAt:  computedAt,
		}
		if err := db.Create(&m).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) ListVelocities(ctx context.Context, params dtos.ListVelocityParams) ([]entities.ProductVelocity, error) {
	var ml []models.ProductVelocity
	q := r.db.Conn(ctx).Model(&models.ProductVelocity{}).
		Where("business_id = ? AND warehouse_id = ?", params.BusinessID, params.WarehouseID)
	if params.Period != "" {
		q = q.Where("period = ?", params.Period)
	}
	if params.Rank != "" {
		q = q.Where("rank = ?", params.Rank)
	}
	limit := params.Limit
	if limit <= 0 {
		limit = 100
	}
	if err := q.Order("units_moved DESC").Limit(limit).Find(&ml).Error; err != nil {
		return nil, err
	}
	out := make([]entities.ProductVelocity, len(ml))
	for i := range ml {
		out[i] = *mappers.ProductVelocityModelToEntity(&ml[i])
	}
	return out, nil
}

func periodToDays(period string) int {
	switch period {
	case "7d":
		return 7
	case "30d":
		return 30
	case "90d":
		return 90
	case "180d":
		return 180
	case "365d":
		return 365
	default:
		return 30
	}
}
