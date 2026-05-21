package repository

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) CarrierConfigs(ctx context.Context, businessID uint) ([]entities.CarrierConfig, error) {
	var rows []models.CarrierCodConfig
	err := r.db.Conn(ctx).
		Where("business_id = ?", businessID).
		Order("carrier_name ASC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]entities.CarrierConfig, len(rows))
	for i := range rows {
		out[i] = entities.CarrierConfig{
			ID:                 rows[i].ID,
			BusinessID:         rows[i].BusinessID,
			CarrierName:        rows[i].CarrierName,
			DiscountPercentage: rows[i].DiscountPercentage,
			IsActive:           rows[i].IsActive,
		}
	}
	return out, nil
}

func (r *Repository) DiscoveredCarriers(ctx context.Context, businessID uint) ([]string, error) {
	var rows []string
	query := `
SELECT DISTINCT UPPER(TRIM(COALESCE(NULLIF(sh.carrier,''),'SIN TRANSPORTADORA'))) AS carrier
FROM shipments sh
INNER JOIN orders o ON o.id = sh.order_id
WHERE sh.deleted_at IS NULL AND o.deleted_at IS NULL
	AND o.cod_total > 0 AND o.business_id = ?
ORDER BY 1`
	if err := r.db.Conn(ctx).Raw(query, businessID).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *Repository) SaveCarrierConfig(ctx context.Context, d dtos.SaveCarrierConfigDTO) (*entities.CarrierConfig, error) {
	name := strings.ToUpper(strings.TrimSpace(d.CarrierName))
	var existing models.CarrierCodConfig
	err := r.db.Conn(ctx).
		Where("business_id = ? AND carrier_name = ?", d.BusinessID, name).
		First(&existing).Error
	if err == nil {
		existing.DiscountPercentage = d.DiscountPercentage
		existing.IsActive = d.IsActive
		if err := r.db.Conn(ctx).Save(&existing).Error; err != nil {
			return nil, err
		}
		return &entities.CarrierConfig{
			ID:                 existing.ID,
			BusinessID:         existing.BusinessID,
			CarrierName:        existing.CarrierName,
			DiscountPercentage: existing.DiscountPercentage,
			IsActive:           existing.IsActive,
		}, nil
	}

	row := models.CarrierCodConfig{
		BusinessID:         d.BusinessID,
		CarrierName:        name,
		DiscountPercentage: d.DiscountPercentage,
		IsActive:           d.IsActive,
	}
	if err := r.db.Conn(ctx).Create(&row).Error; err != nil {
		return nil, err
	}
	return &entities.CarrierConfig{
		ID:                 row.ID,
		BusinessID:         row.BusinessID,
		CarrierName:        row.CarrierName,
		DiscountPercentage: row.DiscountPercentage,
		IsActive:           row.IsActive,
	}, nil
}

func (r *Repository) UserName(ctx context.Context, userID uint) string {
	if userID == 0 {
		return ""
	}
	var name string
	err := r.db.Conn(ctx).
		Table("\"user\"").
		Select("name").
		Where("id = ?", userID).
		Limit(1).
		Scan(&name).Error
	if err != nil {
		return ""
	}
	return name
}
