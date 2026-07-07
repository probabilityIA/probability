package geocoder

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
)

type googleGeocodeResponse struct {
	Status  string `json:"status"`
	Results []struct {
		PartialMatch bool `json:"partial_match"`
		Geometry     struct {
			LocationType string `json:"location_type"`
			Location     struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
}

func (g *Geocoder) Geocode(ctx context.Context, query string) (float64, float64, bool) {
	r := g.GeocodeDetailed(ctx, query)
	return r.Lat, r.Lng, r.Found
}

func (g *Geocoder) GeocodeDetailed(ctx context.Context, query string) ports.GeoResult {
	if query == "" {
		return ports.GeoResult{}
	}

	endpoint := "https://maps.googleapis.com/maps/api/geocode/json?address=" +
		url.QueryEscape(query) + "&key=" + url.QueryEscape(g.apiKey) +
		"&language=es&components=country:co"

	reqCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, endpoint, nil)
	if err != nil {
		return ports.GeoResult{}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		g.logger.Warn(ctx).Err(err).Str("query", query).Msg("geocoder request failed")
		return ports.GeoResult{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ports.GeoResult{}
	}

	var parsed googleGeocodeResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return ports.GeoResult{}
	}
	if parsed.Status != "OK" || len(parsed.Results) == 0 {
		return ports.GeoResult{}
	}

	first := parsed.Results[0]
	loc := first.Geometry.Location
	if loc.Lat == 0 && loc.Lng == 0 {
		return ports.GeoResult{}
	}
	return ports.GeoResult{
		Lat:          loc.Lat,
		Lng:          loc.Lng,
		LocationType: first.Geometry.LocationType,
		PartialMatch: first.PartialMatch,
		Found:        true,
	}
}
