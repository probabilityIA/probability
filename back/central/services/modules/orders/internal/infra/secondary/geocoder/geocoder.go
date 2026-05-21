package geocoder

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

type googleGeocodeResponse struct {
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

func (g *Geocoder) Geocode(ctx context.Context, query string) (float64, float64, bool) {
	if query == "" {
		return 0, 0, false
	}

	endpoint := "https://maps.googleapis.com/maps/api/geocode/json?address=" +
		url.QueryEscape(query) + "&key=" + url.QueryEscape(g.apiKey) +
		"&language=es&components=country:co"

	reqCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, 0, false
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		g.logger.Warn(ctx).Err(err).Str("query", query).Msg("geocoder request failed")
		return 0, 0, false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, false
	}

	var parsed googleGeocodeResponse
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
