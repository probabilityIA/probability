package repository

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/domain"
)

func (r *repository) ChannelSummaries(ctx context.Context, businessID uint) ([]domain.ChannelSummary, error) {
	var rows []struct {
		IntegrationID   uint
		IntegrationName string
		ChannelCode     string
		Matched         int
		NotAssociated   int
		OnlyInChannel   int
		ChannelNoSKU    int
		SKUChanged      int
		SKUTypo         int
		SKUSpacing      int
		FinishedAt      *time.Time
	}

	err := r.db.Conn(ctx).
		Table("integration_sync_runs AS r").
		Select(`r.integration_id, i.name AS integration_name, it.code AS channel_code,
			r.matched, r.not_associated, r.only_in_channel,
			r.channel_no_sku, r.sku_changed, r.sku_typo, r.sku_spacing, r.finished_at`).
		Joins("JOIN integrations i ON i.id = r.integration_id AND i.deleted_at IS NULL AND i.is_active = true").
		Joins("JOIN integration_types it ON it.id = i.integration_type_id AND it.deleted_at IS NULL").
		Where("r.business_id = ? AND r.kind = ? AND r.deleted_at IS NULL", businessID, domain.KindProducts).
		Order("r.finished_at DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	out := make([]domain.ChannelSummary, 0, len(rows))
	for _, row := range rows {
		s := domain.ChannelSummary{
			IntegrationID:   row.IntegrationID,
			IntegrationName: row.IntegrationName,
			ChannelCode:     row.ChannelCode,
			Matched:         row.Matched,
			NotAssociated:   row.NotAssociated,
			OnlyInChannel:   row.OnlyInChannel,
			ChannelNoSKU:    row.ChannelNoSKU,
			SKUChanged:      row.SKUChanged,
			SKUTypo:         row.SKUTypo,
			SKUSpacing:      row.SKUSpacing,
		}
		if row.FinishedAt != nil {
			s.ComparedAt = row.FinishedAt.Format(time.RFC3339)
		}
		out = append(out, s)
	}
	return out, nil
}

func (r *repository) CrossChannel(ctx context.Context, businessID uint) (domain.CrossChannelCounts, error) {
	var out domain.CrossChannelCounts

	const query = `
WITH runs AS (
    SELECT r.id
    FROM integration_sync_runs r
    JOIN integrations i ON i.id = r.integration_id AND i.deleted_at IS NULL AND i.is_active = true
    JOIN integration_types it ON it.id = i.integration_type_id AND it.deleted_at IS NULL
    JOIN integration_categories ic ON ic.id = it.category_id AND ic.code = 'ecommerce'
    WHERE r.business_id = ? AND r.kind = ? AND r.deleted_at IS NULL
),
items AS (
    SELECT UPPER(TRIM(it.sku)) AS sku, it.group_code
    FROM integration_sync_run_items it
    JOIN runs ON runs.id = it.run_id
    WHERE TRIM(it.sku) <> '' AND it.sku NOT LIKE '(%'
),
por_sku AS (
    SELECT sku,
        COUNT(*) FILTER (WHERE group_code IN ('both', 'not_associated')) AS presente,
        COUNT(*) FILTER (WHERE group_code = 'only_probability')          AS ausente,
        COUNT(*) FILTER (WHERE group_code = 'only_channel')              AS solo_canal
    FROM items GROUP BY sku
)
SELECT
    COUNT(*) FILTER (WHERE presente = 0 AND solo_canal = 0 AND ausente > 0) AS not_published,
    COUNT(*) FILTER (WHERE solo_canal > 0 AND presente = 0)                 AS sold_not_owned,
    COUNT(*) FILTER (WHERE presente > 0 AND ausente > 0)                    AS imbalance
FROM por_sku`

	var row struct {
		NotPublished int
		SoldNotOwned int
		Imbalance    int
	}
	if err := r.db.Conn(ctx).Raw(query, businessID, domain.KindProducts).Scan(&row).Error; err != nil {
		return out, err
	}

	out.NotPublished = row.NotPublished
	out.SoldNotOwned = row.SoldNotOwned
	out.Imbalance = row.Imbalance
	return out, nil
}
