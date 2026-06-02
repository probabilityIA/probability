package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func (r *Repository) GetIntegrationBusinessID(ctx context.Context, integrationID uint) (uint, error) {
	var result struct {
		BusinessID *uint
	}
	err := r.db.Conn(ctx).
		Table("integrations").
		Select("business_id").
		Where("id = ? AND deleted_at IS NULL", integrationID).
		Limit(1).
		Scan(&result).Error
	if err != nil {
		return 0, err
	}
	if result.BusinessID == nil || *result.BusinessID == 0 {
		return 0, fmt.Errorf("integracion %d no esta asociada a un negocio", integrationID)
	}
	return *result.BusinessID, nil
}

func (r *Repository) GetCityDaneByName(ctx context.Context, city, province string) (string, error) {
	if city == "" {
		return "", fmt.Errorf("ciudad de destino vacia")
	}

	query := `
		SELECT g.code
		FROM geozones g
		JOIN geozones p ON p.id = g.parent_id AND p.type = 'state'
		WHERE g.type = 'city' AND g.deleted_at IS NULL AND g.code IS NOT NULL
		  AND (
		    unaccent(lower(g.name)) = unaccent(lower(@city))
		    OR unaccent(lower(g.name)) LIKE unaccent(lower(@city)) || ',%'
		    OR unaccent(lower(g.name)) LIKE unaccent(lower(@city)) || '%'
		  )
		ORDER BY
		  (CASE WHEN @province <> '' AND (
		      unaccent(lower(p.name)) LIKE unaccent(lower(@province)) || '%'
		      OR unaccent(lower(@province)) LIKE unaccent(lower(p.name)) || '%'
		  ) THEN 0 ELSE 1 END),
		  (unaccent(lower(g.name)) = unaccent(lower(@city))) DESC,
		  (unaccent(lower(g.name)) LIKE unaccent(lower(@city)) || ',%') DESC,
		  length(g.name) ASC
		LIMIT 1`

	var code string
	err := r.db.Conn(ctx).
		Raw(query, map[string]interface{}{"city": city, "province": province}).
		Scan(&code).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	return code, nil
}
