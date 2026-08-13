package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/domain"
	"github.com/secamc93/probability/back/central/shared/productmatch"
)

func camposSQL() []productmatch.DataFieldSQL {
	return productmatch.DataFieldsSQL()
}

func campoPorKey(key string) (productmatch.DataFieldSQL, bool) {
	return productmatch.DataFieldByKey(key)
}

func condicionLlenar(c productmatch.DataFieldSQL) string {
	return c.FillCondition()
}

func condicionReemplazar(c productmatch.DataFieldSQL) string {
	return c.ReplaceCondition()
}

const baseDiff = `
FROM product_business_integrations m
JOIN products p ON p.id = m.product_id AND p.deleted_at IS NULL
JOIN integrations i ON i.id = m.integration_id AND i.deleted_at IS NULL
LEFT JOIN integration_types it ON it.id = i.integration_type_id AND it.deleted_at IS NULL
WHERE m.business_id = ? AND m.deleted_at IS NULL AND m.channel_snapshot IS NOT NULL`

func (r *repository) DataSummary(ctx context.Context, businessID uint) (*domain.DataSummary, error) {
	campos := camposSQL()

	selects := []string{
		"m.integration_id",
		"COALESCE(i.name,'') AS integration_name",
		"COALESCE(it.code,'') AS channel_code",
		"MAX(m.snapshot_at) AS snapshot_at",
	}
	for _, c := range campos {
		selects = append(selects,
			fmt.Sprintf("COUNT(DISTINCT p.id) FILTER (WHERE %s) AS %s_fill", condicionLlenar(c), c.Key),
			fmt.Sprintf("COUNT(DISTINCT p.id) FILTER (WHERE %s) AS %s_over", condicionReemplazar(c), c.Key),
		)
	}

	query := "SELECT " + strings.Join(selects, ",\n       ") + baseDiff +
		"\nGROUP BY m.integration_id, i.name, it.code ORDER BY COALESCE(i.name,'')"

	rows, err := r.db.Conn(ctx).Raw(query, businessID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type canal struct {
		id     uint
		nombre string
		code   string
		fill   map[string]int
		over   map[string]int
	}

	canales := make([]canal, 0, 4)
	var ultimo *time.Time

	for rows.Next() {
		destino := make([]interface{}, 0, 4+len(campos)*2)
		var id uint
		var nombre, code string
		var snapshotAt *time.Time
		destino = append(destino, &id, &nombre, &code, &snapshotAt)

		valores := make([]int64, len(campos)*2)
		for i := range valores {
			destino = append(destino, &valores[i])
		}

		if err := rows.Scan(destino...); err != nil {
			return nil, err
		}

		if snapshotAt != nil && (ultimo == nil || snapshotAt.After(*ultimo)) {
			ultimo = snapshotAt
		}

		c := canal{id: id, nombre: nombre, code: code, fill: map[string]int{}, over: map[string]int{}}
		for i, campo := range campos {
			c.fill[campo.Key] = int(valores[i*2])
			c.over[campo.Key] = int(valores[i*2+1])
		}
		canales = append(canales, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	summary := &domain.DataSummary{Rows: make([]domain.DataSummaryRow, 0, len(campos)), Snapshotal: ultimo}
	for _, definicion := range domain.DataFields() {
		fila := domain.DataSummaryRow{Field: definicion.Key, Label: definicion.Label, Note: definicion.Note}
		total := 0
		for _, c := range canales {
			fila.Cells = append(fila.Cells, domain.DataSummaryCell{
				IntegrationID:   c.id,
				IntegrationName: c.nombre,
				ChannelCode:     c.code,
				CanFill:         c.fill[definicion.Key],
				CanOverwrite:    c.over[definicion.Key],
			})
			total += c.fill[definicion.Key] + c.over[definicion.Key]
		}
		if total == 0 {
			continue
		}
		summary.Rows = append(summary.Rows, fila)
	}
	return summary, nil
}

const conflictosQuery = `
SELECT o.source,
       o.integration_id,
       COALESCE(i.name, '') AS integration_name,
       COUNT(*)             AS total,
       MAX(o.changed_at)    AS ultimo
FROM product_field_origins o
LEFT JOIN integrations i ON i.id = o.integration_id AND i.deleted_at IS NULL
WHERE o.business_id = ? AND o.field = ? AND o.deleted_at IS NULL
  AND o.product_id IN (%s)
  AND (o.source <> 'channel' OR o.integration_id IS DISTINCT FROM ?)
GROUP BY o.source, o.integration_id, i.name
ORDER BY total DESC`

func (r *repository) DataPreview(ctx context.Context, query domain.DataPreviewQuery) (*domain.DataPreview, error) {
	campo, ok := campoPorKey(query.Field)
	if !ok {
		return nil, domain.ErrInvalidDataField
	}

	condicion := condicionLlenar(campo)
	if query.Mode == domain.ModeOverwrite {
		condicion = "(" + condicionLlenar(campo) + ") OR (" + condicionReemplazar(campo) + ")"
	}

	args := []interface{}{query.BusinessID, query.IntegrationID}
	filtroSKU := ""
	if len(query.SKUs) > 0 {
		filtroSKU = " AND UPPER(TRIM(p.sku)) IN (?)"
		args = append(args, upperTrim(query.SKUs))
	}

	base := baseDiff + " AND m.integration_id = ?" + filtroSKU + " AND (" + condicion + ")"

	var conteo struct {
		Llenar     int64
		Reemplazar int64
	}
	countQuery := fmt.Sprintf(
		"SELECT COUNT(DISTINCT p.id) FILTER (WHERE %s) AS llenar, COUNT(DISTINCT p.id) FILTER (WHERE %s) AS reemplazar %s",
		condicionLlenar(campo), condicionReemplazar(campo), base)
	if err := r.db.Conn(ctx).Raw(countQuery, args...).Scan(&conteo).Error; err != nil {
		return nil, err
	}

	var muestras []domain.DataSample
	sampleQuery := fmt.Sprintf(
		"SELECT * FROM (SELECT DISTINCT ON (p.id) p.id AS product_id, COALESCE(p.sku,'') AS sku, COALESCE(p.name,'') AS name, %s AS current, %s AS incoming %s ORDER BY p.id, m.snapshot_at DESC NULLS LAST, m.id DESC) muestra ORDER BY sku LIMIT 8",
		campo.Own, campo.Channel, base)
	if err := r.db.Conn(ctx).Raw(sampleQuery, args...).Scan(&muestras).Error; err != nil {
		return nil, err
	}

	preview := &domain.DataPreview{
		Field:        query.Field,
		Mode:         query.Mode,
		WouldFill:    int(conteo.Llenar),
		WouldReplace: int(conteo.Reemplazar),
		Samples:      muestras,
		Conflicts:    []domain.DataConflict{},
	}

	idsQuery := "SELECT p.id " + base
	conflictos := fmt.Sprintf(conflictosQuery, idsQuery)
	conflictArgs := append([]interface{}{query.BusinessID, query.Field}, args...)
	conflictArgs = append(conflictArgs, query.IntegrationID)

	var filas []struct {
		Source          string
		IntegrationID   *uint
		IntegrationName string
		Total           int64
		Ultimo          *time.Time
	}
	if err := r.db.Conn(ctx).Raw(conflictos, conflictArgs...).Scan(&filas).Error; err != nil {
		return nil, err
	}
	for _, fila := range filas {
		preview.Conflicts = append(preview.Conflicts, domain.DataConflict{
			Source:          fila.Source,
			IntegrationID:   fila.IntegrationID,
			IntegrationName: fila.IntegrationName,
			Count:           int(fila.Total),
			LastChangeAt:    fila.Ultimo,
		})
	}

	return preview, nil
}

func upperTrim(values []string) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		trimmed := strings.ToUpper(strings.TrimSpace(v))
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
