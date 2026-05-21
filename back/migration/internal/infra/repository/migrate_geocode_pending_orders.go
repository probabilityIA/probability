package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func (r *Repository) backfillGeocodePendingOrders(ctx context.Context) error {
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		return nil
	}

	type pendingOrder struct {
		ID             string
		ShippingStreet string
		ShippingCity   string
		ShippingState  string
	}

	var rows []pendingOrder
	if err := r.db.Conn(ctx).Raw(`
		SELECT id,
		       COALESCE(shipping_street, '') AS shipping_street,
		       COALESCE(shipping_city, '')   AS shipping_city,
		       COALESCE(shipping_state, '')  AS shipping_state
		FROM orders
		WHERE deleted_at IS NULL
		  AND destination_geozone_id IS NULL
		  AND (shipping_lat IS NULL OR shipping_lng IS NULL)
		  AND COALESCE(shipping_city, '') <> ''
	`).Scan(&rows).Error; err != nil {
		return fmt.Errorf("failed to load pending orders for geocoding: %w", err)
	}

	geocoded := 0
	for _, row := range rows {
		query := buildMigrationGeocodeQuery(row.ShippingStreet, row.ShippingCity, row.ShippingState)
		if query == "" {
			continue
		}

		lat, lng, ok := googleGeocodeMigration(ctx, query, apiKey)
		if !ok {
			continue
		}

		if err := r.db.Conn(ctx).Exec(
			`UPDATE orders SET shipping_lat = ?, shipping_lng = ? WHERE id = ?`,
			lat, lng, row.ID,
		).Error; err != nil {
			return fmt.Errorf("failed to update coords for order %s: %w", row.ID, err)
		}
		geocoded++
	}

	fmt.Printf("backfillGeocodePendingOrders: %d pending, %d geocoded\n", len(rows), geocoded)
	return nil
}

func buildMigrationGeocodeQuery(street, city, state string) string {
	parts := make([]string, 0, 5)

	segments := strings.Split(street, "|")
	if first := strings.TrimSpace(segments[0]); first != "" {
		parts = append(parts, first)
	}
	if len(segments) >= 3 {
		if neighborhood := strings.TrimSpace(segments[2]); neighborhood != "" {
			parts = append(parts, neighborhood)
		}
	}
	if c := strings.TrimSpace(city); c != "" {
		parts = append(parts, c)
	}
	if s := strings.TrimSpace(state); s != "" {
		parts = append(parts, s)
	}

	if len(parts) == 0 {
		return ""
	}

	parts = append(parts, "Colombia")
	return strings.Join(parts, ", ")
}

func googleGeocodeMigration(ctx context.Context, query, apiKey string) (float64, float64, bool) {
	endpoint := "https://maps.googleapis.com/maps/api/geocode/json?address=" +
		url.QueryEscape(query) + "&key=" + url.QueryEscape(apiKey) +
		"&language=es&components=country:co"

	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, 0, false
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, 0, false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, false
	}

	var parsed struct {
		Status  string `json:"status"`
		Results []struct {
			Geometry struct {
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
			} `json:"geometry"`
		} `json:"results"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return 0, 0, false
	}
	if parsed.Status != "OK" || len(parsed.Results) == 0 {
		return 0, 0, false
	}

	loc := parsed.Results[0].Geometry.Location
	if loc.Lat == 0 && loc.Lng == 0 {
		return 0, 0, false
	}
	return loc.Lat, loc.Lng, true
}
