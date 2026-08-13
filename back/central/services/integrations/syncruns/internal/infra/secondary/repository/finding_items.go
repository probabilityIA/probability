package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/domain"
)

var grupoDeHallazgo = map[string]string{
	domain.FindingChannelNoSKU:  "channel_no_sku",
	domain.FindingSKUChanged:    "sku_changed",
	domain.FindingSKUTypo:       "sku_typo",
	domain.FindingSKUSpacing:    "sku_spacing",
	domain.FindingNotAssociated: "not_associated",
}

const baseRuns = `
    SELECT r.id, i.name AS channel
    FROM integration_sync_runs r
    JOIN integrations i ON i.id = r.integration_id AND i.deleted_at IS NULL AND i.is_active = true
    WHERE r.business_id = ? AND r.kind = ? AND r.deleted_at IS NULL`

const baseSalesRuns = `
    SELECT r.id, i.name AS channel
    FROM integration_sync_runs r
    JOIN integrations i ON i.id = r.integration_id AND i.deleted_at IS NULL AND i.is_active = true
    JOIN integration_types it2 ON it2.id = i.integration_type_id AND it2.deleted_at IS NULL
    JOIN integration_categories ic ON ic.id = it2.category_id AND ic.code = 'ecommerce'
    WHERE r.business_id = ? AND r.kind = ? AND r.deleted_at IS NULL`

func (r *repository) FindingItems(ctx context.Context, q domain.FindingItemsQuery) (*domain.FindingItemsPage, error) {
	q.Normalize()

	sel, args, err := findingItemsQuery(q)
	if err != nil {
		return nil, err
	}

	page := &domain.FindingItemsPage{
		Items:    []domain.FindingItem{},
		Page:     q.Page,
		PageSize: q.PageSize,
	}

	countSQL := "SELECT COUNT(*) FROM (" + sel + ") AS conteo"
	if err := r.db.Conn(ctx).Raw(countSQL, args...).Scan(&page.Total).Error; err != nil {
		return nil, err
	}

	listSQL := sel + " ORDER BY sku, name"
	listArgs := args
	if !q.All {
		listSQL += " LIMIT ? OFFSET ?"
		listArgs = append(append([]interface{}{}, args...), q.PageSize, q.Offset())
	}

	var rows []struct {
		SKU             string
		Name            string
		Detail          string
		Channels        string
		CounterpartSKU  string
		CounterpartName string
		ChannelQty      *int
		OwnQty          *int
		FixSide         string
		Pattern         string
		PresentIn       string
		MissingIn       string
	}
	if err := r.db.Conn(ctx).Raw(listSQL, listArgs...).Scan(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		item := domain.FindingItem{
			SKU: row.SKU, Name: row.Name, Detail: row.Detail,
			CounterpartSKU: row.CounterpartSKU, CounterpartName: row.CounterpartName,
			ChannelQty: row.ChannelQty, OwnQty: row.OwnQty,
			FixSide: row.FixSide, Pattern: row.Pattern,
		}
		if row.Channels != "" {
			item.Channels = strings.Split(row.Channels, "|")
		}
		if row.PresentIn != "" {
			item.PresentIn = strings.Split(row.PresentIn, "|")
		}
		if row.MissingIn != "" {
			item.MissingIn = strings.Split(row.MissingIn, "|")
		}
		page.Items = append(page.Items, item)
	}

	page.TotalPages = int((page.Total + int64(q.PageSize) - 1) / int64(q.PageSize))
	if q.All {
		page.TotalPages = 1
	}
	return page, nil
}

func findingItemsQuery(q domain.FindingItemsQuery) (string, []interface{}, error) {
	if grupo, ok := grupoDeHallazgo[q.Code]; ok {
		return porGrupo(q, grupo)
	}
	switch q.Code {
	case domain.FindingNotPublished, domain.FindingSoldNotOwned, domain.FindingImbalance:
		return cruzado(q)
	}
	return "", nil, fmt.Errorf("hallazgo desconocido: %s", q.Code)
}

func porGrupo(q domain.FindingItemsQuery, grupo string) (string, []interface{}, error) {
	sel := `
SELECT it.sku AS sku,
       COALESCE(MAX(NULLIF(it.parent_label, '')), '') AS name,
       COALESCE(MAX(it.label), '') AS detail,
       string_agg(DISTINCT runs.channel, '|') AS channels,
       COALESCE(MAX(it.counterpart_sku), '')  AS counterpart_sku,
       COALESCE(MAX(it.counterpart_name), '') AS counterpart_name,
       MAX(it.channel_qty) AS channel_qty,
       MAX(it.own_qty)     AS own_qty,
       COALESCE(MAX(it.fix_side), '') AS fix_side,
       COALESCE(MAX(it.pattern), '')  AS pattern,
       '' AS present_in,
       '' AS missing_in
FROM integration_sync_run_items it
JOIN (` + baseRuns + `) AS runs ON runs.id = it.run_id
WHERE it.group_code = ?`

	args := []interface{}{q.BusinessID, domain.KindProducts, grupo}
	if q.Search != "" {
		sel += " AND (it.sku ILIKE ? OR it.label ILIKE ? OR it.parent_label ILIKE ?)"
		like := "%" + q.Search + "%"
		args = append(args, like, like, like)
	}
	sel += " GROUP BY " + agrupacion(grupo)
	return sel, args, nil
}

func agrupacion(grupo string) string {
	if grupo == "channel_no_sku" {
		return "it.sku, it.parent_ref, it.variant_label"
	}
	return "it.sku"
}

func cruzado(q domain.FindingItemsQuery) (string, []interface{}, error) {
	var having, detail, estaEn, faltaEn string
	switch q.Code {
	case domain.FindingNotPublished:
		having = "presente = 0 AND solo_canal = 0 AND ausente > 0"
		detail = "'no esta en ningun canal de venta'"
		estaEn = "''"
		faltaEn = "COALESCE(faltantes, '')"
	case domain.FindingSoldNotOwned:
		having = "solo_canal > 0 AND presente = 0"
		detail = "'se vende en el canal pero no existe en Probability'"
		estaEn = "COALESCE(en_canal, '')"
		faltaEn = "'Probability'"
	case domain.FindingImbalance:
		having = "presente > 0 AND ausente > 0"
		detail = "''"
		estaEn = "COALESCE(presentes, '')"
		faltaEn = "COALESCE(faltantes, '')"
	}

	sel := `
SELECT sku,
       COALESCE(nombre, '') AS name,
       ` + detail + ` AS detail,
       COALESCE(todos, '') AS channels,
       ` + estaEn + ` AS present_in,
       ` + faltaEn + ` AS missing_in,
       '' AS counterpart_sku, '' AS counterpart_name,
       NULL::int AS channel_qty, NULL::int AS own_qty,
       '' AS fix_side, '' AS pattern
FROM (
    SELECT UPPER(TRIM(it.sku)) AS sku,
        MAX(NULLIF(it.parent_label, '')) AS nombre,
        COUNT(*) FILTER (WHERE it.group_code IN ('both', 'not_associated')) AS presente,
        COUNT(*) FILTER (WHERE it.group_code = 'only_probability')          AS ausente,
        COUNT(*) FILTER (WHERE it.group_code = 'only_channel')              AS solo_canal,
        string_agg(DISTINCT runs.channel, '|') AS todos,
        string_agg(DISTINCT runs.channel, '|') FILTER (WHERE it.group_code IN ('both', 'not_associated')) AS presentes,
        string_agg(DISTINCT runs.channel, '|') FILTER (WHERE it.group_code = 'only_probability')          AS faltantes,
        string_agg(DISTINCT runs.channel, '|') FILTER (WHERE it.group_code = 'only_channel')              AS en_canal
    FROM integration_sync_run_items it
    JOIN (` + baseSalesRuns + `) AS runs ON runs.id = it.run_id
    WHERE TRIM(it.sku) <> '' AND it.sku NOT LIKE '(%'
    GROUP BY UPPER(TRIM(it.sku))
) AS por_sku
WHERE ` + having

	args := []interface{}{q.BusinessID, domain.KindProducts}
	if q.Search != "" {
		sel += " AND sku ILIKE ?"
		args = append(args, "%"+q.Search+"%")
	}
	return sel, args, nil
}
