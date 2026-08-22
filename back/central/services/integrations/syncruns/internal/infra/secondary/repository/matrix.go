package repository

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/domain"
)

const matrixColumnsQuery = `
SELECT i.id AS integration_id, i.name,
       COALESCE(it.code, '')                      AS code,
       COALESCE(it.image_url, '')                 AS image_url,
       COALESCE(ic.code = 'ecommerce', false)     AS is_sales
FROM integrations i
JOIN integration_types it ON it.id = i.integration_type_id AND it.deleted_at IS NULL
LEFT JOIN integration_categories ic ON ic.id = it.category_id
WHERE i.business_id = ? AND i.deleted_at IS NULL AND i.is_active = true
  AND EXISTS (
      SELECT 1 FROM product_business_integrations pbi
      WHERE pbi.integration_id = i.id AND pbi.deleted_at IS NULL
  )
ORDER BY i.name`

const matrixCellsQuery = `
SELECT pbi.product_id, pbi.integration_id,
       COALESCE(pbi.external_sku, '')        AS sku,
       COALESCE(pbi.external_barcode, '')    AS barcode,
       COALESCE(pbi.external_product_id, '') AS external_id,
       COALESCE(pbi.external_variant_id, '') AS variant_id
FROM product_business_integrations pbi
JOIN integrations i ON i.id = pbi.integration_id AND i.deleted_at IS NULL AND i.is_active = true
WHERE pbi.business_id = ? AND pbi.deleted_at IS NULL AND pbi.product_id IN ?`

func (r *repository) MatchMatrix(ctx context.Context, q domain.MatrixQuery) (*domain.MatrixPage, error) {
	q.Normalize()

	page := &domain.MatrixPage{
		Columns:  []domain.MatrixColumn{},
		Rows:     []domain.MatrixRow{},
		Page:     q.Page,
		PageSize: q.PageSize,
	}

	var columnas []struct {
		IntegrationID uint
		Name          string
		Code          string
		ImageURL      string
		IsSales       bool
	}
	if err := r.db.Conn(ctx).Raw(matrixColumnsQuery, q.BusinessID).Scan(&columnas).Error; err != nil {
		return nil, err
	}
	for _, c := range columnas {
		page.Columns = append(page.Columns, domain.MatrixColumn{
			IntegrationID: c.IntegrationID, Name: c.Name, Code: c.Code,
			ImageURL: c.ImageURL, IsSales: c.IsSales,
		})
	}

	base := r.db.Conn(ctx).
		Table("products p").
		Where("p.business_id = ? AND p.deleted_at IS NULL AND p.is_active = true", q.BusinessID)
	if q.Search != "" {
		like := "%" + q.Search + "%"
		switch q.SearchBy {
		case domain.MatrixSearchSKU:
			base = base.Where("p.sku ILIKE ?", like)
		case domain.MatrixSearchName:
			base = base.Where("p.name ILIKE ?", like)
		case domain.MatrixSearchBarcode:
			base = base.Where("p.barcode ILIKE ?", like)
		case domain.MatrixSearchAll:
			base = base.Where("p.sku ILIKE ? OR p.name ILIKE ? OR p.barcode ILIKE ?", like, like, like)
		default:
			return nil, fmt.Errorf("search_by invalido: %s", q.SearchBy)
		}
	}
	for _, id := range q.PresentIn {
		if id > 0 {
			base = base.Where("EXISTS ("+estaEnIntegracion+")", id)
		}
	}
	for _, id := range q.MissingIn {
		if id > 0 {
			base = base.Where("NOT EXISTS ("+estaEnIntegracion+")", id)
		}
	}

	if err := base.Session(&gorm.Session{}).Count(&page.Total).Error; err != nil {
		return nil, err
	}

	var filas []struct {
		ID       string
		SKU      string
		Barcode  string
		Name     string
		ImageURL string
	}
	listado := base.Session(&gorm.Session{}).
		Select("p.id, p.sku, COALESCE(p.barcode, '') AS barcode, p.name, COALESCE(p.image_url, '') AS image_url").
		Order("p.sku")
	if !q.All {
		listado = listado.Limit(q.PageSize).Offset(q.Offset())
	}
	if err := listado.Scan(&filas).Error; err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(filas))
	for _, f := range filas {
		ids = append(ids, f.ID)
	}

	celdas := map[string]map[uint]domain.MatrixCell{}
	if len(ids) > 0 {
		var rows []struct {
			ProductID     string
			IntegrationID uint
			SKU           string
			Barcode       string
			ExternalID    string
			VariantID     string
		}
		if err := r.db.Conn(ctx).Raw(matrixCellsQuery, q.BusinessID, ids).Scan(&rows).Error; err != nil {
			return nil, err
		}
		for _, row := range rows {
			if celdas[row.ProductID] == nil {
				celdas[row.ProductID] = map[uint]domain.MatrixCell{}
			}
			celdas[row.ProductID][row.IntegrationID] = domain.MatrixCell{
				IntegrationID: row.IntegrationID,
				Present:       true,
				SKU:           row.SKU,
				Barcode:       row.Barcode,
				ExternalID:    row.ExternalID,
				VariantID:     row.VariantID,
			}
		}
	}

	for _, f := range filas {
		fila := domain.MatrixRow{
			ProductID: f.ID, SKU: f.SKU, Barcode: f.Barcode, Name: f.Name,
			ImageURL: f.ImageURL,
			Cells:    make([]domain.MatrixCell, 0, len(page.Columns)),
		}
		for _, col := range page.Columns {
			celda, ok := celdas[f.ID][col.IntegrationID]
			if !ok {
				fila.Cells = append(fila.Cells, domain.MatrixCell{IntegrationID: col.IntegrationID})
				continue
			}
			if celda.SKU == "" {
				celda.SKU = f.SKU
				celda.SKUMatches = true
			} else {
				celda.SKUKnown = true
				celda.SKUMatches = mismoSKU(celda.SKU, f.SKU)
			}
			fila.Cells = append(fila.Cells, celda)
		}
		page.Rows = append(page.Rows, fila)
	}

	page.TotalPages = int((page.Total + int64(q.PageSize) - 1) / int64(q.PageSize))
	if q.All {
		page.TotalPages = 1
	}
	return page, nil
}

func mismoSKU(a, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}

const estaEnIntegracion = `
    SELECT 1 FROM product_business_integrations m
    WHERE m.product_id = p.id AND m.integration_id = ? AND m.deleted_at IS NULL`
