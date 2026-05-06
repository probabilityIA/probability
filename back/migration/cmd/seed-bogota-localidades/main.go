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
		fmt.Println("uso: seed-bogota-localidades <localidades.geojson>")
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
		logger.Fatal(ctx).Msg("bogota city not found in geozones")
	}
	logger.Info().Uint("bogota_id", bogotaID).Msg("Bogota city found")

	created, skipped, errs := 0, 0, 0
	for _, f := range fc.Features {
		name := str(f.Properties, "Nombre_de_la_localidad")
		locID := str(f.Properties, "Identificador_unico_de_la_localidad")

		if name == "" || locID == "" {
			skipped++
			continue
		}

		code := fmt.Sprintf("BOG-LOC-%s", locID)

		var existing uint
		if err := database.Conn(ctx).Raw(`SELECT id FROM geozones WHERE business_id = 0 AND type = 'admin_district' AND code = ? AND deleted_at IS NULL LIMIT 1`, code).Scan(&existing).Error; err != nil {
			errs++
			continue
		}
		if existing > 0 {
			skipped++
			continue
		}

		props := fmt.Sprintf(`{"source":"datosabiertos.bogota.gov.co/SDHT","city_code":"%s","localidad_id":"%s"}`, bogotaCityCode, locID)
		err := database.Conn(ctx).Exec(`
			INSERT INTO geozones (business_id, parent_id, type, code, name, geometry, centroid, properties, is_active)
			VALUES (
				0, ?, 'admin_district', ?, ?,
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

	logger.Info().Int("created", created).Int("skipped", skipped).Int("errors", errs).Msg("DONE")
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
