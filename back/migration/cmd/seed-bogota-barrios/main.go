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
		fmt.Println("uso: seed-bogota-barrios <barrios.geojson>")
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

	created, skipped, errs := 0, 0, 0
	for _, f := range fc.Features {
		name := str(f.Properties, "SCANOMBRE")
		bCode := str(f.Properties, "SCACODIGO")
		if name == "" || bCode == "" {
			skipped++
			continue
		}
		code := fmt.Sprintf("BOG-BAR-%s", bCode)

		var existing uint
		if err := database.Conn(ctx).Raw(`SELECT id FROM geozones WHERE business_id = 0 AND type = 'barrio' AND code = ? AND deleted_at IS NULL LIMIT 1`, code).Scan(&existing).Error; err != nil {
			errs++
			continue
		}
		if existing > 0 {
			skipped++
			continue
		}

		props := fmt.Sprintf(`{"source":"catastrobogota.arcgis","sca_code":"%s","sca_tipo":%v}`, bCode, f.Properties["SCATIPO"])

		err := database.Conn(ctx).Exec(`
			INSERT INTO geozones (business_id, parent_id, type, code, name, geometry, centroid, properties, is_active)
			VALUES (
				0, ?, 'barrio', ?, ?,
				ST_Multi(ST_SetSRID(ST_GeomFromGeoJSON(?), 4326))::geometry(MultiPolygon, 4326),
				ST_PointOnSurface(ST_SetSRID(ST_GeomFromGeoJSON(?), 4326))::geography,
				?::jsonb,
				TRUE
			)
		`, bogotaID, code, name, string(f.Geometry), string(f.Geometry), props).Error
		if err != nil {
			errs++
			if errs < 10 {
				logger.Error(ctx).Str("name", name).Err(err).Msg("insert failed")
			}
			continue
		}
		created++
		if (created)%200 == 0 {
			logger.Info().Int("created", created).Msg("progress")
		}
	}
	logger.Info().Int("created", created).Int("skipped", skipped).Int("errors", errs).Msg("Barrios inserted")

	logger.Info().Msg("Reassigning parent_id from Bogota city to UPZ (neighborhood)...")
	res := database.Conn(ctx).Exec(`
		UPDATE geozones b
		SET parent_id = upz.id, updated_at = NOW()
		FROM geozones upz
		WHERE b.type = 'barrio'
		  AND b.business_id = 0
		  AND b.deleted_at IS NULL
		  AND upz.type = 'neighborhood'
		  AND upz.business_id = 0
		  AND upz.deleted_at IS NULL
		  AND upz.properties->>'kind' = 'upz'
		  AND ST_Contains(upz.geometry::geometry, b.centroid::geometry)
	`)
	if res.Error != nil {
		logger.Error(ctx).Err(res.Error).Msg("reassign parent")
	} else {
		logger.Info().Int64("rows", res.RowsAffected).Msg("Barrios assigned to UPZ")
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
