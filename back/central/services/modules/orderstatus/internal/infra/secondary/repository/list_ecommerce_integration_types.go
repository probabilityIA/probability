package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
)

func (r *repository) ListEcommerceIntegrationTypes(ctx context.Context, businessID uint) ([]entities.IntegrationTypeInfo, error) {
	type row struct {
		ID       uint
		Code     string
		Name     string
		ImageURL string
	}
	var rows []row

	if businessID == 0 {
		// Super admin — todos los tipos de integración de categoría ecommerce
		err := r.db.Conn(ctx).Raw(`
			SELECT it.id, it.code, it.name, it.image_url
			FROM integration_types it
			JOIN integration_categories ic ON ic.id = it.category_id
			WHERE ic.code = 'ecommerce'
			  AND it.is_active = true
			  AND it.deleted_at IS NULL
			  AND ic.deleted_at IS NULL
			ORDER BY it.name ASC
		`).Scan(&rows).Error
		if err != nil {
			return nil, err
		}
	} else {
		// Scope business — solo los tipos donde el negocio tiene una integración activa
		err := r.db.Conn(ctx).Raw(`
			SELECT DISTINCT it.id, it.code, it.name, it.image_url
			FROM integration_types it
			JOIN integration_categories ic ON ic.id = it.category_id
			JOIN integrations i ON i.integration_type_id = it.id
			WHERE ic.code = 'ecommerce'
			  AND it.is_active = true
			  AND i.business_id = ?
			  AND it.deleted_at IS NULL
			  AND ic.deleted_at IS NULL
			  AND i.deleted_at IS NULL
			ORDER BY it.name ASC
		`, businessID).Scan(&rows).Error
		if err != nil {
			return nil, err
		}
	}

	result := make([]entities.IntegrationTypeInfo, len(rows))
	for i, row := range rows {
		result[i] = entities.IntegrationTypeInfo{
			ID:       row.ID,
			Code:     row.Code,
			Name:     row.Name,
			ImageURL: row.ImageURL,
		}
	}
	return result, nil
}
