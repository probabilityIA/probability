package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/domain"
)

// grupoDeHallazgo mapea los hallazgos que salen tal cual del detalle de cada
// canal. Los demas se calculan cruzando canales y llevan su propia consulta.
var grupoDeHallazgo = map[string]string{
	domain.FindingChannelNoSKU:  "channel_no_sku",
	domain.FindingSKUChanged:    "sku_changed",
	domain.FindingSKUTypo:       "sku_typo",
	domain.FindingNotAssociated: "not_associated",
}

// baseRuns limita a las corridas de comparacion del negocio, sobre canales activos.
const baseRuns = `
    SELECT r.id, i.name AS channel
    FROM integration_sync_runs r
    JOIN integrations i ON i.id = r.integration_id AND i.deleted_at IS NULL AND i.is_active = true
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

	listSQL := sel + " ORDER BY sku"
	listArgs := args
	if !q.All {
		listSQL += " LIMIT ? OFFSET ?"
		listArgs = append(append([]interface{}{}, args...), q.PageSize, q.Offset())
	}

	var rows []struct {
		SKU      string
		Name     string
		Detail   string
		Channels string
	}
	if err := r.db.Conn(ctx).Raw(listSQL, listArgs...).Scan(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		item := domain.FindingItem{SKU: row.SKU, Name: row.Name, Detail: row.Detail}
		if row.Channels != "" {
			item.Channels = strings.Split(row.Channels, "|")
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

// porGrupo agrupa por SKU: el mismo problema puede estar en varios canales y se
// muestra una sola fila con todos.
func porGrupo(q domain.FindingItemsQuery, grupo string) (string, []interface{}, error) {
	sel := `
SELECT it.sku AS sku,
       COALESCE(MAX(NULLIF(it.parent_label, '')), '') AS name,
       COALESCE(MAX(it.label), '') AS detail,
       string_agg(DISTINCT runs.channel, '|') AS channels
FROM integration_sync_run_items it
JOIN (` + baseRuns + `) AS runs ON runs.id = it.run_id
WHERE it.group_code = ?`

	args := []interface{}{q.BusinessID, domain.KindProducts, grupo}
	if q.Search != "" {
		sel += " AND (it.sku ILIKE ? OR it.label ILIKE ? OR it.parent_label ILIKE ?)"
		like := "%" + q.Search + "%"
		args = append(args, like, like, like)
	}
	sel += " GROUP BY it.sku"
	return sel, args, nil
}

// cruzado resuelve los hallazgos que solo existen al mirar todos los canales a la
// vez. Por cada SKU arma en que canales esta y en cuales falta.
func cruzado(q domain.FindingItemsQuery) (string, []interface{}, error) {
	var having, detail string
	switch q.Code {
	case domain.FindingNotPublished:
		having = "presente = 0 AND solo_canal = 0 AND ausente > 0"
		detail = "'no esta en ningun canal de venta'"
	case domain.FindingSoldNotOwned:
		having = "solo_canal > 0 AND presente = 0"
		detail = "'se vende en: ' || COALESCE(en_canal, 'ninguno')"
	case domain.FindingImbalance:
		having = "presente > 0 AND ausente > 0"
		detail = "'esta en: ' || COALESCE(presentes, '-') || '  ·  falta en: ' || COALESCE(faltantes, '-')"
	}

	sel := `
SELECT sku,
       COALESCE(nombre, '') AS name,
       ` + detail + ` AS detail,
       COALESCE(todos, '') AS channels
FROM (
    SELECT UPPER(TRIM(it.sku)) AS sku,
        MAX(NULLIF(it.parent_label, '')) AS nombre,
        COUNT(*) FILTER (WHERE it.group_code IN ('both', 'not_associated')) AS presente,
        COUNT(*) FILTER (WHERE it.group_code = 'only_probability')          AS ausente,
        COUNT(*) FILTER (WHERE it.group_code = 'only_channel')              AS solo_canal,
        string_agg(DISTINCT runs.channel, '|') AS todos,
        string_agg(DISTINCT runs.channel, ', ') FILTER (WHERE it.group_code IN ('both', 'not_associated')) AS presentes,
        string_agg(DISTINCT runs.channel, ', ') FILTER (WHERE it.group_code = 'only_probability')          AS faltantes,
        string_agg(DISTINCT runs.channel, ', ') FILTER (WHERE it.group_code = 'only_channel')              AS en_canal
    FROM integration_sync_run_items it
    JOIN (` + baseRuns + `) AS runs ON runs.id = it.run_id
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
