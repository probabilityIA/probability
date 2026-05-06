package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/secamc93/probability/back/migration/shared/db"
	"github.com/secamc93/probability/back/migration/shared/env"
	"github.com/secamc93/probability/back/migration/shared/log"
)

type feature struct {
	Geometry   json.RawMessage        `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

type featureCollection struct {
	Features []feature `json:"features"`
}

const bogotaCityCode = "11001"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("uso: seed-bogota-upz <upz.geojson>")
		os.Exit(2)
	}
	logger := log.New()
	cfg := env.New(logger)
	database := db.New(logger, cfg)
	defer database.Close()
	ctx := context.Background()

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		logger.Fatal(ctx).Err(err).Msg("read file")
	}
	var fc featureCollection
	if err := json.Unmarshal(data, &fc); err != nil {
		logger.Fatal(ctx).Err(err).Msg("parse json")
	}
	logger.Info().Int("features", len(fc.Features)).Msg("loaded")

	var bogotaID uint
	if err := database.Conn(ctx).Raw(`SELECT id FROM geozones WHERE business_id = 0 AND type = 'city' AND code = ? AND deleted_at IS NULL LIMIT 1`, bogotaCityCode).Scan(&bogotaID).Error; err != nil {
		logger.Fatal(ctx).Err(err).Msg("find bogota")
	}
	if bogotaID == 0 {
		logger.Fatal(ctx).Msg("bogota city not found")
	}

	created, skipped, errs := 0, 0, 0
	for _, f := range fc.Features {
		name := str(f.Properties, "NOMBRE")
		upzCode := str(f.Properties, "CODIGO_UPZ")
		if name == "" || upzCode == "" {
			skipped++
			continue
		}

		code := fmt.Sprintf("BOG-UPZ-%s", upzCode)

		var existing uint
		if err := database.Conn(ctx).Raw(`SELECT id FROM geozones WHERE business_id = 0 AND type = 'neighborhood' AND code = ? AND deleted_at IS NULL LIMIT 1`, code).Scan(&existing).Error; err != nil {
			errs++
			continue
		}
		if existing > 0 {
			skipped++
			continue
		}

		props := fmt.Sprintf(`{"source":"catastrobogota.arcgis","upz_code":"%s","decreto_pot":"%s","kind":"upz"}`, upzCode, str(f.Properties, "DECRETO_POT"))

		err := database.Conn(ctx).Exec(`
			INSERT INTO geozones (business_id, parent_id, type, code, name, geometry, centroid, properties, is_active)
			VALUES (
				0, ?, 'neighborhood', ?, ?,
				ST_Multi(ST_SetSRID(ST_GeomFromGeoJSON(?), 4326))::geometry(MultiPolygon, 4326),
				ST_PointOnSurface(ST_SetSRID(ST_GeomFromGeoJSON(?), 4326))::geography,
				?::jsonb,
				TRUE
			)
		`, bogotaID, code, name, string(f.Geometry), string(f.Geometry), props).Error
		if err != nil {
			errs++
			logger.Error(ctx).Str("name", name).Err(err).Msg("insert failed")
			continue
		}
		created++
	}
	logger.Info().Int("created", created).Int("skipped", skipped).Int("errors", errs).Msg("UPZ inserted")

	logger.Info().Msg("Reassigning parent_id from Bogota city to its admin_district (localidad)...")
	res := database.Conn(ctx).Exec(`
		UPDATE geozones u
		SET parent_id = ad.id, updated_at = NOW()
		FROM geozones ad
		WHERE u.type = 'neighborhood'
		  AND u.business_id = 0
		  AND u.deleted_at IS NULL
		  AND u.properties->>'kind' = 'upz'
		  AND ad.type = 'admin_district'
		  AND ad.business_id = 0
		  AND ad.deleted_at IS NULL
		  AND ad.parent_id = ?
		  AND ST_Contains(ad.geometry::geometry, u.centroid::geometry)
	`, bogotaID)
	if res.Error != nil {
		logger.Error(ctx).Err(res.Error).Msg("reassign parent")
	} else {
		logger.Info().Int64("rows", res.RowsAffected).Msg("UPZs assigned to localidades")
	}
}

func str(m map[string]interface{}, k string) string {
	if v, ok := m[k]; ok && v != nil {
		if s, ok := v.(string); ok {
			return s
		}
		return fmt.Sprintf("%v", v)
	}
	return ""
}
