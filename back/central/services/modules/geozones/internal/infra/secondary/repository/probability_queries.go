package repository

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/ports"
	"gorm.io/gorm"
)

func (r *Repository) AncestorsByOrderID(ctx context.Context, orderID string, businessID uint) (*entities.GeozoneAncestors, error) {
	type row struct {
		DestinationGeozoneID    *uint
		GeozoneCountryID        *uint
		GeozoneStateID          *uint
		GeozoneCityID           *uint
		GeozoneAdminDistrictID  *uint
		GeozoneLocalityID       *uint
		GeozoneNeighborhoodID   *uint
		GeozoneBarrioID         *uint
	}
	var rec row
	err := r.db.Conn(ctx).Raw(`
		SELECT destination_geozone_id,
		       geozone_country_id, geozone_state_id, geozone_city_id,
		       geozone_admin_district_id, geozone_locality_id,
		       geozone_neighborhood_id, geozone_barrio_id
		FROM orders
		WHERE id = ? AND (business_id = ? OR ? = 0)
		LIMIT 1
	`, orderID, businessID, businessID).Scan(&rec).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entities.GeozoneAncestors{
		DeepestID:       rec.DestinationGeozoneID,
		CountryID:       rec.GeozoneCountryID,
		StateID:         rec.GeozoneStateID,
		CityID:          rec.GeozoneCityID,
		AdminDistrictID: rec.GeozoneAdminDistrictID,
		LocalityID:      rec.GeozoneLocalityID,
		NeighborhoodID:  rec.GeozoneNeighborhoodID,
		BarrioID:        rec.GeozoneBarrioID,
	}, nil
}

var allowedLevelColumns = map[string]bool{
	"geozone_country_id":         true,
	"geozone_state_id":           true,
	"geozone_city_id":            true,
	"geozone_admin_district_id":  true,
	"geozone_locality_id":        true,
	"geozone_neighborhood_id":    true,
	"geozone_barrio_id":          true,
}

func (r *Repository) AggregateAtLevel(ctx context.Context, businessID uint, levelColumn string, geozoneID uint, carrier string) (ports.LevelAggregate, error) {
	if !allowedLevelColumns[levelColumn] {
		return ports.LevelAggregate{}, errors.New("invalid level column")
	}
	type agg struct {
		Total     int64
		Delivered int64
		Cancelled int64
		Returned  int64
		InTransit int64
	}
	var out agg
	args := []any{businessID, geozoneID}
	carrierFilter := ""
	if carrier != "" {
		carrierFilter = " AND s.carrier = ?"
		args = append(args, carrier)
	}
	q := `
		SELECT
		  COUNT(*) AS total,
		  COUNT(*) FILTER (WHERE s.status = 'delivered') AS delivered,
		  COUNT(*) FILTER (WHERE s.status = 'cancelled') AS cancelled,
		  COUNT(*) FILTER (WHERE s.status = 'returned')  AS returned,
		  COUNT(*) FILTER (WHERE s.status NOT IN ('delivered','cancelled','returned')) AS in_transit
		FROM shipments s
		JOIN orders o ON o.id = s.order_id
		WHERE s.deleted_at IS NULL
		  AND o.business_id = ?
		  AND s.` + levelColumn + ` = ?` + carrierFilter
	if err := r.db.Conn(ctx).Raw(q, args...).Scan(&out).Error; err != nil {
		return ports.LevelAggregate{}, err
	}
	return ports.LevelAggregate{
		Total:     out.Total,
		Delivered: out.Delivered,
		Cancelled: out.Cancelled,
		Returned:  out.Returned,
		InTransit: out.InTransit,
	}, nil
}

func (r *Repository) GeozoneNameAndType(ctx context.Context, geozoneID uint) (string, string, error) {
	type row struct {
		Name string
		Type string
	}
	var rec row
	if err := r.db.Conn(ctx).Raw(
		`SELECT name, type FROM geozones WHERE id = ? LIMIT 1`, geozoneID,
	).Scan(&rec).Error; err != nil {
		return "", "", err
	}
	return rec.Name, rec.Type, nil
}
