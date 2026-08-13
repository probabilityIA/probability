package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
	"github.com/secamc93/probability/back/central/shared/productmatch"
	"gorm.io/gorm"
)

const aplicarPlantilla = `
WITH objetivo AS (
    SELECT DISTINCT ON (p.id) p.id AS product_id, %s AS actual, %s AS nuevo
    FROM product_business_integrations m
    JOIN products p ON p.id = m.product_id AND p.deleted_at IS NULL
    WHERE m.business_id = ? AND m.integration_id = ? AND m.deleted_at IS NULL
      AND m.channel_snapshot IS NOT NULL
      %s
      AND (%s)
    ORDER BY p.id, m.snapshot_at DESC NULLS LAST, m.id DESC
    LIMIT %d
),
historial AS (
    INSERT INTO product_field_changes
        (product_id, field, business_id, old_value, truncated, source, integration_id, batch_id, created_at)
    SELECT product_id, ?, ?, LEFT(actual, %d), LENGTH(actual) > %d, 'channel', ?, ?, ?
    FROM objetivo
    RETURNING product_id
),
origen AS (
    INSERT INTO product_field_origins
        (product_id, field, business_id, source, integration_id, user_id, changed_at, created_at, updated_at)
    SELECT product_id, ?, ?, 'channel', ?, NULL, ?, ?, ?
    FROM objetivo
    ON CONFLICT (product_id, field) DO UPDATE SET
        source = EXCLUDED.source,
        integration_id = EXCLUDED.integration_id,
        user_id = NULL,
        changed_at = EXCLUDED.changed_at,
        updated_at = EXCLUDED.updated_at,
        deleted_at = NULL
    RETURNING product_id
)
UPDATE products
SET %s = objetivo.nuevo%s, updated_at = ?
FROM objetivo
WHERE products.id = objetivo.product_id`

const podarPlantilla = `
DELETE FROM product_field_changes
WHERE id IN (
    SELECT id FROM (
        SELECT id, ROW_NUMBER() OVER (
            PARTITION BY product_id, field ORDER BY created_at DESC, id DESC
        ) AS posicion
        FROM product_field_changes
        WHERE business_id = ? AND field = ?
    ) ordenadas
    WHERE posicion > ?
)`

func (r *Repository) ApplyChannelField(ctx context.Context, req domain.DataApplyRequest, batchID string, at time.Time) (int, error) {
	campo, ok := productmatch.DataFieldByKey(req.Field)
	if !ok {
		return 0, domain.ErrInvalidProductData
	}

	condicion := campo.FillCondition()
	if req.Mode == domain.ModeOverwrite {
		condicion = "(" + campo.FillCondition() + ") OR (" + campo.ReplaceCondition() + ")"
	}

	filtroSKU := ""
	args := []interface{}{req.BusinessID, req.IntegrationID}
	if len(req.SKUs) > 0 {
		filtroSKU = "AND UPPER(TRIM(p.sku)) IN (?)"
		args = append(args, mayusculas(req.SKUs))
	}

	query := fmt.Sprintf(aplicarPlantilla,
		campo.Own, campo.Channel,
		filtroSKU, condicion,
		domain.MaxProductsPerApply,
		domain.MaxOldValueLen, domain.MaxOldValueLen,
		campo.Column, campo.Cast,
	)

	args = append(args,
		req.Field, req.BusinessID, req.IntegrationID, batchID, at,
		req.Field, req.BusinessID, req.IntegrationID, at, at, at,
		at,
	)

	var afectados int64
	err := r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		resultado := tx.Exec(query, args...)
		if resultado.Error != nil {
			return resultado.Error
		}
		afectados = resultado.RowsAffected

		return tx.Exec(podarPlantilla, req.BusinessID, req.Field, domain.MaxChangesPerField).Error
	})
	if err != nil {
		return 0, err
	}
	return int(afectados), nil
}

type cambioRevertible struct {
	ID        uint
	ProductID string
	Field     string
	OldValue  string
}

const deshacerObjetivo = `
SELECT c.id, c.product_id, c.field, c.old_value
FROM product_field_changes c
WHERE c.business_id = ? AND c.batch_id = ?`

func (r *Repository) UndoBatch(ctx context.Context, businessID uint, batchID string, userID uint, at time.Time) (int, error) {
	var filas []cambioRevertible
	if err := r.db.Conn(ctx).Raw(deshacerObjetivo, businessID, batchID).Scan(&filas).Error; err != nil {
		return 0, err
	}
	if len(filas) == 0 {
		return 0, nil
	}

	revertidos := 0
	err := r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		for _, fila := range filas {
			campo, ok := productmatch.DataFieldByKey(fila.Field)
			if !ok {
				continue
			}

			if fila.Field == domain.FieldVariant {
				if err := tx.Exec(`
UPDATE products
SET family_id = NULL, variant_label = ?, variant_attributes = NULL, variant_signature = '', updated_at = ?
WHERE id = ? AND business_id = ?`, fila.OldValue, at, fila.ProductID, businessID).Error; err != nil {
					return err
				}
				revertidos++
				continue
			}

			valor := interface{}(fila.OldValue)
			asignacion := fmt.Sprintf("%s = ?", campo.Column)
			if campo.Cast != "" {
				asignacion = fmt.Sprintf("%s = NULLIF(?, '')%s", campo.Column, campo.Cast)
			} else if fila.OldValue == "" && campo.Column == "barcode" {
				asignacion = fmt.Sprintf("%s = NULL", campo.Column)
				valor = nil
			}

			update := fmt.Sprintf("UPDATE products SET %s, updated_at = ? WHERE id = ? AND business_id = ?", asignacion)
			var args []interface{}
			if valor == nil {
				args = []interface{}{at, fila.ProductID, businessID}
			} else {
				args = []interface{}{valor, at, fila.ProductID, businessID}
			}
			if err := tx.Exec(update, args...).Error; err != nil {
				return err
			}
			revertidos++
		}

		if err := tx.Exec(`DELETE FROM product_field_changes WHERE business_id = ? AND batch_id = ?`,
			businessID, batchID).Error; err != nil {
			return err
		}

		return tx.Exec(`
UPDATE product_field_origins
SET source = 'manual', integration_id = NULL, user_id = ?, changed_at = ?, updated_at = ?
WHERE business_id = ? AND product_id IN (?) AND field IN (?)`,
			nullableUser(userID), at, at, businessID, productIDs(filas), fields(filas)).Error
	})
	if err != nil {
		return 0, err
	}
	return revertidos, nil
}

const lotesQuery = `
SELECT c.batch_id,
       MIN(c.field)         AS field,
       MIN(c.integration_id) AS integration_id,
       COUNT(DISTINCT c.product_id) AS products,
       MIN(c.created_at)    AS created_at
FROM product_field_changes c
WHERE c.business_id = ? AND c.batch_id <> ''
GROUP BY c.batch_id
ORDER BY MIN(c.created_at) DESC
LIMIT ?`

func (r *Repository) ListDataBatches(ctx context.Context, businessID uint, limit int) ([]domain.DataBatch, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	var batches []domain.DataBatch
	if err := r.db.Conn(ctx).Raw(lotesQuery, businessID, limit).Scan(&batches).Error; err != nil {
		return nil, err
	}
	return batches, nil
}

func mayusculas(values []string) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		trimmed := strings.ToUpper(strings.TrimSpace(v))
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func nullableUser(userID uint) interface{} {
	if userID == 0 {
		return nil
	}
	return userID
}

func productIDs(filas []cambioRevertible) []string {
	vistos := make(map[string]bool, len(filas))
	out := make([]string, 0, len(filas))
	for _, f := range filas {
		if !vistos[f.ProductID] {
			vistos[f.ProductID] = true
			out = append(out, f.ProductID)
		}
	}
	return out
}

func fields(filas []cambioRevertible) []string {
	vistos := make(map[string]bool, len(filas))
	out := make([]string, 0, 4)
	for _, f := range filas {
		if !vistos[f.Field] {
			vistos[f.Field] = true
			out = append(out, f.Field)
		}
	}
	return out
}
