package inventorycompare

import (
	"context"
	"strings"
	"time"

	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LoadOptions struct {
	Page     int
	PageSize int
	OnlyDiff bool
	Search   string
	SKUs     []string
}

func Save(
	ctx context.Context,
	conn *gorm.DB,
	businessID, integrationID uint,
	rows []Row,
	checkedAt time.Time,
) error {
	if len(rows) == 0 {
		return nil
	}

	vistos := make(map[string]bool, len(rows))
	filas := make([]models.InventoryCompareSnapshot, 0, len(rows))
	for _, fila := range rows {
		if strings.TrimSpace(fila.ProductID) == "" {
			continue
		}
		clave := fila.ProductID + "|" + strings.TrimSpace(fila.ExternalVariantID)
		if vistos[clave] {
			continue
		}
		vistos[clave] = true
		filas = append(filas, models.InventoryCompareSnapshot{
			BusinessID:        businessID,
			IntegrationID:     integrationID,
			ProductID:         fila.ProductID,
			SKU:               fila.SKU,
			Name:              fila.Name,
			ImageURL:          fila.ImageURL,
			ExternalItemID:    fila.ExternalItemID,
			ExternalVariantID: fila.ExternalVariantID,
			ProbabilityQty:    fila.ProbabilityQty,
			ChannelQty:        fila.ChannelQty,
			Delta:             fila.Delta,
			Action:            fila.Action,
			Reason:            fila.Reason,
			CheckedAt:         checkedAt,
		})
	}
	if len(filas) == 0 {
		return nil
	}

	return conn.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "integration_id"}, {Name: "product_id"}, {Name: "external_variant_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"business_id", "sku", "name", "image_url", "external_item_id",
			"probability_qty", "channel_qty", "delta", "action", "reason",
			"checked_at", "updated_at",
		}),
	}).CreateInBatches(&filas, 200).Error
}

func Load(
	ctx context.Context,
	conn *gorm.DB,
	businessID, integrationID uint,
	opts LoadOptions,
) (*Page, error) {
	page, pageSize := NormalizePaging(opts.Page, opts.PageSize)

	resumen := struct {
		Total     int64
		ToUpdate  int64
		Unchanged int64
		Skipped   int64
		OldestAt  *time.Time
		NewestAt  *time.Time
	}{}

	err := conn.WithContext(ctx).
		Model(&models.InventoryCompareSnapshot{}).
		Select(`COUNT(*) AS total,
			COUNT(*) FILTER (WHERE action = 'update') AS to_update,
			COUNT(*) FILTER (WHERE action = 'unchanged') AS unchanged,
			COUNT(*) FILTER (WHERE action = 'skip') AS skipped,
			MIN(checked_at) AS oldest_at,
			MAX(checked_at) AS newest_at`).
		Where("business_id = ? AND integration_id = ?", businessID, integrationID).
		Scan(&resumen).Error
	if err != nil {
		return nil, err
	}

	resultado := &Page{
		Rows:      []Row{},
		Page:      page,
		PageSize:  pageSize,
		FromCache: true,
		Totals: Totals{
			Total:     int(resumen.Total),
			ToUpdate:  int(resumen.ToUpdate),
			Unchanged: int(resumen.Unchanged),
			Skipped:   int(resumen.Skipped),
		},
	}
	if resumen.OldestAt != nil {
		resultado.CheckedAt = *resumen.OldestAt
	}
	if resumen.NewestAt != nil {
		resultado.NewestAt = *resumen.NewestAt
	}
	if resumen.Total == 0 {
		resultado.TotalPages = 1
		return resultado, nil
	}

	filtrada := conn.WithContext(ctx).
		Model(&models.InventoryCompareSnapshot{}).
		Where("business_id = ? AND integration_id = ?", businessID, integrationID)

	if opts.OnlyDiff {
		filtrada = filtrada.Where("action = ?", ActionUpdate)
	}
	if buscar := strings.TrimSpace(opts.Search); buscar != "" {
		patron := "%" + strings.ToLower(buscar) + "%"
		filtrada = filtrada.Where("LOWER(sku) LIKE ? OR LOWER(name) LIKE ?", patron, patron)
	}
	if len(opts.SKUs) > 0 {
		buscados := make([]string, 0, len(opts.SKUs))
		for _, sku := range opts.SKUs {
			if limpio := strings.ToLower(strings.TrimSpace(sku)); limpio != "" {
				buscados = append(buscados, limpio)
			}
		}
		if len(buscados) > 0 {
			filtrada = filtrada.Where("LOWER(sku) IN ?", buscados)
		}
	}

	var total int64
	if err := filtrada.Count(&total).Error; err != nil {
		return nil, err
	}
	resultado.Total = int(total)
	resultado.TotalPages = TotalPages(int(total), pageSize)
	if total == 0 {
		return resultado, nil
	}

	var guardadas []models.InventoryCompareSnapshot
	err = filtrada.
		Order("sku ASC, external_variant_id ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&guardadas).Error
	if err != nil {
		return nil, err
	}

	filas := make([]Row, 0, len(guardadas))
	for _, fila := range guardadas {
		filas = append(filas, Row{
			ProductID:         fila.ProductID,
			SKU:               fila.SKU,
			Name:              fila.Name,
			ImageURL:          fila.ImageURL,
			ExternalItemID:    fila.ExternalItemID,
			ExternalVariantID: fila.ExternalVariantID,
			ProbabilityQty:    fila.ProbabilityQty,
			ChannelQty:        fila.ChannelQty,
			Delta:             fila.Delta,
			Action:            fila.Action,
			Reason:            fila.Reason,
		})
	}
	resultado.Rows = filas
	return resultado, nil
}
