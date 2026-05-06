package repository

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"gorm.io/gorm"
)

type Repository struct {
	db db.IDatabase
}

func New(database db.IDatabase) ports.IRepository {
	return &Repository{db: database}
}

type geozoneRow struct {
	ID         uint
	BusinessID uint
	ParentID   *uint
	Type       string
	Code       *string
	Name       string
	GeometryJS *string
	CentroidJS *string
	Properties []byte
	IsActive   bool
	CreatedAt  any
	UpdatedAt  any
}

func (r *Repository) Create(ctx context.Context, dto dtos.CreateGeozoneDTO) (*entities.Geozone, error) {
	props := dto.Properties
	if len(props) == 0 {
		props = json.RawMessage(`{}`)
	}

	var id uint
	err := r.db.Conn(ctx).Raw(`
		INSERT INTO geozones (business_id, parent_id, type, code, name, geometry, centroid, properties, is_active)
		VALUES (
			?, ?, ?, ?, ?,
			ST_Multi(ST_SetSRID(ST_GeomFromGeoJSON(?), 4326))::geometry(MultiPolygon, 4326),
			ST_PointOnSurface(ST_SetSRID(ST_GeomFromGeoJSON(?), 4326))::geography,
			?::jsonb,
			TRUE
		)
		RETURNING id
	`, dto.BusinessID, dto.ParentID, dto.Type, dto.Code, dto.Name,
		string(dto.Geometry), string(dto.Geometry), string(props),
	).Scan(&id).Error

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, domainerrors.ErrDuplicateGeozone
		}
		if strings.Contains(err.Error(), "GeoJSON") || strings.Contains(err.Error(), "geometry") {
			return nil, domainerrors.ErrInvalidGeometry
		}
		return nil, err
	}

	return r.GetByID(ctx, id, false)
}

func (r *Repository) GetByID(ctx context.Context, id uint, includeGeom bool) (*entities.Geozone, error) {
	geomCol := "NULL::text AS geometry_js"
	if includeGeom {
		geomCol = "ST_AsGeoJSON(geometry) AS geometry_js"
	}

	var row geozoneRow
	err := r.db.Conn(ctx).Raw(`
		SELECT id, business_id, parent_id, type, code, name,
		       `+geomCol+`,
		       ST_AsGeoJSON(centroid::geometry) AS centroid_js,
		       properties, is_active, created_at, updated_at
		FROM geozones
		WHERE id = ? AND deleted_at IS NULL
	`, id).Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.ID == 0 {
		return nil, domainerrors.ErrGeozoneNotFound
	}
	return rowToEntity(&row), nil
}

func (r *Repository) GetByCode(ctx context.Context, businessID uint, geozoneType, code string) (*entities.Geozone, error) {
	var row geozoneRow
	err := r.db.Conn(ctx).Raw(`
		SELECT id, business_id, parent_id, type, code, name,
		       NULL::text AS geometry_js,
		       NULL::text AS centroid_js,
		       properties, is_active, created_at, updated_at
		FROM geozones
		WHERE business_id = ? AND type = ? AND code = ? AND deleted_at IS NULL
		LIMIT 1
	`, businessID, geozoneType, code).Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return rowToEntity(&row), nil
}

func (r *Repository) List(ctx context.Context, params dtos.ListGeozonesParams) ([]entities.Geozone, int64, error) {
	geomCol := "NULL::text AS geometry_js"
	if params.IncludeGeom {
		geomCol = "ST_AsGeoJSON(geometry) AS geometry_js"
	}

	where := []string{"deleted_at IS NULL", "(business_id = ? OR business_id = 0)"}
	args := []any{params.BusinessID}

	if params.Type != "" {
		where = append(where, "type = ?")
		args = append(args, params.Type)
	}
	if params.ParentID != nil {
		where = append(where, "parent_id = ?")
		args = append(args, *params.ParentID)
	}
	if params.Code != "" {
		where = append(where, "code = ?")
		args = append(args, params.Code)
	}
	if params.Search != "" {
		where = append(where, "name ILIKE ?")
		args = append(args, "%"+params.Search+"%")
	}

	whereSQL := strings.Join(where, " AND ")

	var total int64
	if err := r.db.Conn(ctx).Raw(`SELECT COUNT(*) FROM geozones WHERE `+whereSQL, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	listArgs := append(args, params.PageSize, params.Offset())
	var rows []geozoneRow
	err := r.db.Conn(ctx).Raw(`
		SELECT id, business_id, parent_id, type, code, name,
		       `+geomCol+`,
		       ST_AsGeoJSON(centroid::geometry) AS centroid_js,
		       properties, is_active, created_at, updated_at
		FROM geozones
		WHERE `+whereSQL+`
		ORDER BY type, name
		LIMIT ? OFFSET ?
	`, listArgs...).Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}

	out := make([]entities.Geozone, len(rows))
	for i := range rows {
		out[i] = *rowToEntity(&rows[i])
	}
	return out, total, nil
}

func (r *Repository) LookupContaining(ctx context.Context, params dtos.LookupParams) ([]entities.Geozone, error) {
	where := []string{"deleted_at IS NULL", "is_active = TRUE", "(business_id = ? OR business_id = 0)",
		"ST_Contains(geometry, ST_SetSRID(ST_MakePoint(?, ?), 4326))"}
	args := []any{params.BusinessID, params.Lng, params.Lat}

	if params.Type != "" {
		where = append(where, "type = ?")
		args = append(args, params.Type)
	}

	var rows []geozoneRow
	err := r.db.Conn(ctx).Raw(`
		SELECT id, business_id, parent_id, type, code, name,
		       NULL::text AS geometry_js,
		       NULL::text AS centroid_js,
		       properties, is_active, created_at, updated_at
		FROM geozones
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY CASE type
		    WHEN 'neighborhood' THEN 1
		    WHEN 'locality' THEN 2
		    WHEN 'city' THEN 3
		    WHEN 'state' THEN 4
		    WHEN 'country' THEN 5
		    ELSE 6 END
	`, args...).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	out := make([]entities.Geozone, len(rows))
	for i := range rows {
		out[i] = *rowToEntity(&rows[i])
	}
	return out, nil
}

func (r *Repository) Delete(ctx context.Context, id uint) error {
	res := r.db.Conn(ctx).Exec(`UPDATE geozones SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domainerrors.ErrGeozoneNotFound
	}
	return nil
}

func rowToEntity(r *geozoneRow) *entities.Geozone {
	g := &entities.Geozone{
		ID:         r.ID,
		BusinessID: r.BusinessID,
		ParentID:   r.ParentID,
		Type:       r.Type,
		Code:       r.Code,
		Name:       r.Name,
		Properties: json.RawMessage(r.Properties),
		IsActive:   r.IsActive,
	}
	if r.GeometryJS != nil {
		g.Geometry = json.RawMessage(*r.GeometryJS)
	}
	if r.CentroidJS != nil {
		g.Centroid = json.RawMessage(*r.CentroidJS)
	}
	return g
}
