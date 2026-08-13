package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const loteVariante = 200

const candidatosVarianteQuery = `
SELECT * FROM (
    SELECT DISTINCT ON (p.id)
           p.id                                              AS product_id,
           p.sku                                             AS sku,
           COALESCE(p.variant_label, '')                     AS label_actual,
           COALESCE(m.channel_snapshot->>'family_name', '')  AS family_name,
           COALESCE(m.channel_snapshot->>'color', '')        AS color,
           COALESCE(m.channel_snapshot->>'size', '')         AS size
    FROM product_business_integrations m
    JOIN products p ON p.id = m.product_id AND p.deleted_at IS NULL
    WHERE m.business_id = ? AND m.integration_id = ? AND m.deleted_at IS NULL
      AND m.channel_snapshot IS NOT NULL
      AND p.family_id IS NULL
      AND COALESCE(m.channel_snapshot->>'family_name', '') <> ''
      AND (COALESCE(m.channel_snapshot->>'color', '') <> '' OR COALESCE(m.channel_snapshot->>'size', '') <> '')
      %s
    ORDER BY p.id, m.snapshot_at DESC NULLS LAST, m.id DESC
) candidato
ORDER BY sku
LIMIT %d`

type candidatoVariante struct {
	ProductID   string
	SKU         string
	LabelActual string
	FamilyName  string
	Color       string
	Size        string
}

func (c candidatoVariante) atributos() map[string]string {
	out := make(map[string]string, 2)
	if v := strings.TrimSpace(c.Color); v != "" {
		out["color"] = v
	}
	if v := strings.TrimSpace(c.Size); v != "" {
		out["talla"] = v
	}
	return out
}

func (c candidatoVariante) etiqueta() string {
	partes := make([]string, 0, 2)
	if v := strings.TrimSpace(c.Color); v != "" {
		partes = append(partes, v)
	}
	if v := strings.TrimSpace(c.Size); v != "" {
		partes = append(partes, v)
	}
	return strings.Join(partes, " / ")
}

type varianteLista struct {
	ProductID  string
	FamilyID   uint
	Label      string
	Attributes datatypes.JSON
	Signature  string
	LabelPrev  string
}

func (r *Repository) candidatosDeVariante(ctx context.Context, req domain.DataApplyRequest) ([]candidatoVariante, error) {
	filtro := ""
	args := []interface{}{req.BusinessID, req.IntegrationID}
	if len(req.SKUs) > 0 {
		filtro = "AND UPPER(TRIM(p.sku)) IN (?)"
		args = append(args, mayusculas(req.SKUs))
	}

	var candidatos []candidatoVariante
	query := fmt.Sprintf(candidatosVarianteQuery, filtro, domain.MaxProductsPerApply)
	if err := r.db.Conn(ctx).Raw(query, args...).Scan(&candidatos).Error; err != nil {
		return nil, err
	}
	return candidatos, nil
}

func (r *Repository) ApplyChannelVariant(
	ctx context.Context,
	req domain.DataApplyRequest,
	batchID string,
	at time.Time,
	avance domain.ProgressFunc,
) (int, error) {
	candidatos, err := r.candidatosDeVariante(ctx, req)
	if err != nil {
		return 0, err
	}
	if len(candidatos) == 0 {
		return 0, nil
	}

	familias, err := r.resolverFamilias(ctx, req.BusinessID, candidatos)
	if err != nil {
		return 0, err
	}

	ocupadas, err := r.firmasOcupadas(ctx, familias)
	if err != nil {
		return 0, err
	}

	pendientes := make([]varianteLista, 0, len(candidatos))
	for _, candidato := range candidatos {
		atributos := candidato.atributos()
		if len(atributos) == 0 {
			continue
		}
		crudos, merr := json.Marshal(atributos)
		if merr != nil {
			continue
		}
		canonicos, firma, cerr := domain.CanonicalizeVariantAttributes(datatypes.JSON(crudos))
		if cerr != nil || firma == "" {
			continue
		}
		familiaID := familias[strings.ToLower(strings.TrimSpace(candidato.FamilyName))]
		if familiaID == 0 {
			continue
		}
		clave := firmaClave(familiaID, firma)
		if ocupadas[clave] {
			continue
		}
		ocupadas[clave] = true

		pendientes = append(pendientes, varianteLista{
			ProductID:  candidato.ProductID,
			FamilyID:   familiaID,
			Label:      candidato.etiqueta(),
			Attributes: canonicos,
			Signature:  firma,
			LabelPrev:  candidato.LabelActual,
		})
	}

	aplicados := 0
	for inicio := 0; inicio < len(pendientes); inicio += loteVariante {
		fin := inicio + loteVariante
		if fin > len(pendientes) {
			fin = len(pendientes)
		}
		lote := pendientes[inicio:fin]

		if err := r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
			if uerr := actualizarVariantes(tx, req.BusinessID, lote, at); uerr != nil {
				return uerr
			}
			return registrarVariantes(tx, req, lote, batchID, at)
		}); err != nil {
			return aplicados, err
		}

		aplicados += len(lote)
		if avance != nil {
			avance(aplicados, len(pendientes))
		}
	}

	if aplicados > 0 {
		_ = r.db.Conn(ctx).Exec(podarPlantilla, req.BusinessID, domain.FieldVariant, domain.MaxChangesPerField).Error
	}
	return aplicados, nil
}

func (r *Repository) resolverFamilias(ctx context.Context, businessID uint, candidatos []candidatoVariante) (map[string]uint, error) {
	vistos := make(map[string]string, 16)
	for _, c := range candidatos {
		limpio := strings.TrimSpace(c.FamilyName)
		clave := strings.ToLower(limpio)
		if clave == "" {
			continue
		}
		vistos[clave] = limpio
	}
	if len(vistos) == 0 {
		return map[string]uint{}, nil
	}

	claves := make([]string, 0, len(vistos))
	for clave := range vistos {
		claves = append(claves, clave)
	}

	var existentes []struct {
		ID   uint
		Name string
	}
	if err := r.db.Conn(ctx).Table("product_families").
		Select("id, name").
		Where("business_id = ? AND deleted_at IS NULL AND LOWER(TRIM(name)) IN ?", businessID, claves).
		Scan(&existentes).Error; err != nil {
		return nil, err
	}

	familias := make(map[string]uint, len(vistos))
	for _, e := range existentes {
		familias[strings.ToLower(strings.TrimSpace(e.Name))] = e.ID
	}

	valores := make([]string, 0, len(vistos))
	args := make([]interface{}, 0, len(vistos)*2)
	for clave, nombre := range vistos {
		if familias[clave] > 0 {
			continue
		}
		valores = append(valores, "(?, ?, 'active', true, NOW(), NOW())")
		args = append(args, businessID, nombre)
	}
	if len(valores) == 0 {
		return familias, nil
	}

	var creadas []struct {
		ID   uint
		Name string
	}
	insert := "INSERT INTO product_families (business_id, name, status, is_active, created_at, updated_at) VALUES " +
		strings.Join(valores, ", ") + " RETURNING id, name"
	if err := r.db.Conn(ctx).Raw(insert, args...).Scan(&creadas).Error; err != nil {
		return nil, err
	}
	for _, c := range creadas {
		familias[strings.ToLower(strings.TrimSpace(c.Name))] = c.ID
	}
	return familias, nil
}

func (r *Repository) firmasOcupadas(ctx context.Context, familias map[string]uint) (map[string]bool, error) {
	ids := make([]uint, 0, len(familias))
	for _, id := range familias {
		if id > 0 {
			ids = append(ids, id)
		}
	}
	ocupadas := make(map[string]bool)
	if len(ids) == 0 {
		return ocupadas, nil
	}

	var filas []struct {
		FamilyID         uint
		VariantSignature string
	}
	if err := r.db.Conn(ctx).Table("products").
		Select("family_id, variant_signature").
		Where("family_id IN ? AND deleted_at IS NULL AND variant_signature <> ''", ids).
		Scan(&filas).Error; err != nil {
		return nil, err
	}
	for _, fila := range filas {
		ocupadas[firmaClave(fila.FamilyID, fila.VariantSignature)] = true
	}
	return ocupadas, nil
}

func actualizarVariantes(tx *gorm.DB, businessID uint, lote []varianteLista, at time.Time) error {
	valores := make([]string, 0, len(lote))
	args := make([]interface{}, 0, len(lote)*5+2)
	args = append(args, at)
	for _, v := range lote {
		valores = append(valores, "(?, ?::bigint, ?, ?::jsonb, ?)")
		args = append(args, v.ProductID, v.FamilyID, v.Label, string(v.Attributes), v.Signature)
	}
	args = append(args, businessID)

	query := `
UPDATE products p
SET family_id = v.family_id,
    variant_label = v.label,
    variant_attributes = v.attrs,
    variant_signature = v.sig,
    updated_at = ?
FROM (VALUES ` + strings.Join(valores, ", ") + `) AS v(id, family_id, label, attrs, sig)
WHERE p.id = v.id AND p.business_id = ?`

	return tx.Exec(query, args...).Error
}

func registrarVariantes(tx *gorm.DB, req domain.DataApplyRequest, lote []varianteLista, batchID string, at time.Time) error {
	historia := make([]string, 0, len(lote))
	histArgs := make([]interface{}, 0, len(lote)*7)
	origenes := make([]string, 0, len(lote))
	origArgs := make([]interface{}, 0, len(lote)*7)

	for _, v := range lote {
		historia = append(historia, "(?, ?, ?, ?, false, 'channel', ?, ?, ?)")
		histArgs = append(histArgs, v.ProductID, domain.FieldVariant, req.BusinessID, v.LabelPrev, req.IntegrationID, batchID, at)

		origenes = append(origenes, "(?, ?, ?, 'channel', ?, ?, ?, ?)")
		origArgs = append(origArgs, v.ProductID, domain.FieldVariant, req.BusinessID, req.IntegrationID, at, at, at)
	}

	insertHistoria := `
INSERT INTO product_field_changes
    (product_id, field, business_id, old_value, truncated, source, integration_id, batch_id, created_at)
VALUES ` + strings.Join(historia, ", ")
	if err := tx.Exec(insertHistoria, histArgs...).Error; err != nil {
		return err
	}

	insertOrigen := `
INSERT INTO product_field_origins
    (product_id, field, business_id, source, integration_id, changed_at, created_at, updated_at)
VALUES ` + strings.Join(origenes, ", ") + `
ON CONFLICT (product_id, field) DO UPDATE SET
    source = EXCLUDED.source,
    integration_id = EXCLUDED.integration_id,
    user_id = NULL,
    changed_at = EXCLUDED.changed_at,
    updated_at = EXCLUDED.updated_at,
    deleted_at = NULL`
	return tx.Exec(insertOrigen, origArgs...).Error
}

func firmaClave(familiaID uint, firma string) string {
	return fmt.Sprintf("%d|%s", familiaID, firma)
}
