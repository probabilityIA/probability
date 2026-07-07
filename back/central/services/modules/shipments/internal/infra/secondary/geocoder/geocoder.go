package geocoder

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

type googleResponse struct {
	Status  string `json:"status"`
	Results []struct {
		FormattedAddress  string `json:"formatted_address"`
		PartialMatch      bool   `json:"partial_match"`
		AddressComponents []struct {
			LongName string   `json:"long_name"`
			Types    []string `json:"types"`
		} `json:"address_components"`
		Geometry struct {
			LocationType string `json:"location_type"`
			Location     struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
}

type Geocoder struct {
	apiKey string
	client *http.Client
}

func New(apiKey string) domain.IGeocoder {
	if apiKey == "" {
		return nil
	}
	return &Geocoder{apiKey: apiKey, client: &http.Client{Timeout: 8 * time.Second}}
}

func (g *Geocoder) Geocode(ctx context.Context, address string) (domain.GeocodeResult, error) {
	endpoint := "https://maps.googleapis.com/maps/api/geocode/json?address=" +
		url.QueryEscape(address) + "&components=country:CO&language=es&key=" + url.QueryEscape(g.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return domain.GeocodeResult{}, err
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return domain.GeocodeResult{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.GeocodeResult{}, err
	}

	var parsed googleResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return domain.GeocodeResult{}, err
	}

	if parsed.Status != "OK" || len(parsed.Results) == 0 {
		return domain.GeocodeResult{Found: false}, nil
	}

	r := parsed.Results[0]
	out := domain.GeocodeResult{
		Found:            true,
		Lat:              r.Geometry.Location.Lat,
		Lng:              r.Geometry.Location.Lng,
		FormattedAddress: r.FormattedAddress,
		LocationType:     r.Geometry.LocationType,
		PartialMatch:     r.PartialMatch,
	}

	for _, comp := range r.AddressComponents {
		for _, t := range comp.Types {
			switch t {
			case "locality":
				out.Locality = comp.LongName
			case "administrative_area_level_2":
				out.AdminArea2 = comp.LongName
			case "administrative_area_level_1":
				out.Department = comp.LongName
			}
		}
	}

	return out, nil
}
