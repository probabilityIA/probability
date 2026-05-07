package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

const probabilityMinSampleScore = 20

var levelChain = []struct {
	column string
	level  string
}{
	{"geozone_barrio_id", "barrio"},
	{"geozone_neighborhood_id", "neighborhood"},
	{"geozone_admin_district_id", "admin_district"},
	{"geozone_locality_id", "locality"},
	{"geozone_city_id", "city"},
	{"geozone_state_id", "state"},
	{"geozone_country_id", "country"},
}

func (r *Repository) GetGeozoneDeliveryRateForOrder(ctx context.Context, orderID string) (*float64, string, *uint, error) {
	type orderRow struct {
		BusinessID            *uint
		GeozoneCountryID      *uint
		GeozoneStateID        *uint
		GeozoneCityID         *uint
		GeozoneAdminDistrictID *uint
		GeozoneLocalityID     *uint
		GeozoneNeighborhoodID *uint
		GeozoneBarrioID       *uint
	}
	var ord orderRow
	if err := r.db.Conn(ctx).Raw(`
		SELECT business_id,
		       geozone_country_id, geozone_state_id, geozone_city_id,
		       geozone_admin_district_id, geozone_locality_id,
		       geozone_neighborhood_id, geozone_barrio_id
		FROM orders WHERE id = ? LIMIT 1
	`, orderID).Scan(&ord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", nil, nil
		}
		return nil, "", nil, err
	}
	if ord.BusinessID == nil || *ord.BusinessID == 0 {
		return nil, "", nil, nil
	}
	getter := map[string]*uint{
		"geozone_barrio_id":         ord.GeozoneBarrioID,
		"geozone_neighborhood_id":   ord.GeozoneNeighborhoodID,
		"geozone_admin_district_id": ord.GeozoneAdminDistrictID,
		"geozone_locality_id":       ord.GeozoneLocalityID,
		"geozone_city_id":           ord.GeozoneCityID,
		"geozone_state_id":          ord.GeozoneStateID,
		"geozone_country_id":        ord.GeozoneCountryID,
	}
	for _, lvl := range levelChain {
		gid := getter[lvl.column]
		if gid == nil || *gid == 0 {
			continue
		}
		type agg struct {
			Total     int64
			Delivered int64
		}
		var out agg
		q := `
			WITH live AS (
			    SELECT COUNT(*) FILTER (WHERE s.status IN ('delivered','failed','returned')) AS total,
			           COUNT(*) FILTER (WHERE s.status = 'delivered') AS delivered
			    FROM shipments s
			    WHERE s.deleted_at IS NULL
			      AND s.` + lvl.column + ` = ?
			), historical AS (
			    SELECT COALESCE(SUM(m.delivered + m.failed + m.returned),0) AS total,
			           COALESCE(SUM(m.delivered),0) AS delivered
			    FROM geozone_monthly_stats m
			    WHERE m.geozone_id = ?
			)
			SELECT l.total + h.total AS total, l.delivered + h.delivered AS delivered
			FROM live l, historical h
		`
		if err := r.db.Conn(ctx).Raw(q, *gid, *gid).Scan(&out).Error; err != nil {
			return nil, "", nil, err
		}
		if out.Total < probabilityMinSampleScore {
			continue
		}
		rate := float64(out.Delivered) / float64(out.Total) * 100
		gidCopy := *gid
		return &rate, lvl.level, &gidCopy, nil
	}
	return nil, "", nil, nil
}
