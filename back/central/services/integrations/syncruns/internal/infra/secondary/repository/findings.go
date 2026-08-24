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

	// Los contadores de la corrida cuentan ocurrencias; la lista del hallazgo
	// agrupa por SKU, que es como se cruza con el canal. Se recalculan con la
	// misma agrupacion para que el numero del hallazgo y su lista coincidan.
	porSKU, err := r.contarHallazgosPorSKU(ctx, businessID)
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
		if grupos, ok := porSKU[row.IntegrationID]; ok && len(grupos) > 0 {
			s.NotAssociated = grupos[domain.FindingNotAssociated]
			s.ChannelNoSKU = grupos[domain.FindingChannelNoSKU]
			s.SKUChanged = grupos[domain.FindingSKUChanged]
			s.SKUTypo = grupos[domain.FindingSKUTypo]
			s.SKUSpacing = grupos[domain.FindingSKUSpacing]
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

const conteoDesdeItems = `
SELECT integration_id, group_code, COUNT(*) AS total FROM (
  SELECT r.integration_id, it.group_code, it.sku,
         CASE WHEN it.group_code = 'channel_no_sku'
              THEN COALESCE(it.parent_ref, '') || '|' || COALESCE(it.variant_label, '')
              ELSE '' END AS extra
  FROM integration_sync_run_items it
  JOIN integration_sync_runs r ON r.id = it.run_id AND r.deleted_at IS NULL
  JOIN integrations i ON i.id = r.integration_id AND i.deleted_at IS NULL AND i.is_active = true
  WHERE r.business_id = ? AND r.kind = ?
  GROUP BY r.integration_id, it.group_code, it.sku, extra
) g GROUP BY integration_id, group_code`

// Las corridas viejas no escribieron integration_sync_run_items: su detalle
// quedo como JSON en la propia fila.
const conteoDesdeDetalle = `
SELECT integration_id, grupo AS group_code, COUNT(*) AS total FROM (
  SELECT r.integration_id, it->>'group' AS grupo, it->>'sku' AS sku,
         CASE WHEN it->>'group' = 'channel_no_sku'
              THEN COALESCE(it->>'parent_ref', '') || '|' || COALESCE(it->>'variant_label', '')
              ELSE '' END AS extra
  FROM integration_sync_runs r
  JOIN integrations i ON i.id = r.integration_id AND i.deleted_at IS NULL AND i.is_active = true
  CROSS JOIN LATERAL jsonb_array_elements(r.detail) AS it
  WHERE r.business_id = ? AND r.kind = ? AND r.deleted_at IS NULL
    AND jsonb_typeof(r.detail) = 'array'
    AND NOT EXISTS (SELECT 1 FROM integration_sync_run_items x WHERE x.run_id = r.id)
  GROUP BY r.integration_id, grupo, sku, extra
) g GROUP BY integration_id, grupo`

func (r *repository) contarHallazgosPorSKU(ctx context.Context, businessID uint) (map[uint]map[string]int, error) {
	out := map[uint]map[string]int{}

	for _, sql := range []string{conteoDesdeItems, conteoDesdeDetalle} {
		var rows []struct {
			IntegrationID uint
			GroupCode     string
			Total         int
		}
		if err := r.db.Conn(ctx).Raw(sql, businessID, domain.KindProducts).Scan(&rows).Error; err != nil {
			return nil, err
		}
		for _, row := range rows {
			if _, ok := out[row.IntegrationID]; !ok {
				out[row.IntegrationID] = map[string]int{}
			}
			out[row.IntegrationID][row.GroupCode] = row.Total
		}
	}

	return out, nil
}
