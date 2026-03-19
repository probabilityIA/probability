package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/shared/env"
)

// GeocodingResult representa el resultado estandarizado de geocodificación.
type GeocodingResult struct {
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Found    bool    `json:"found"`
	Fallback bool    `json:"fallback"` // true si se usó solo la ciudad como fallback
}

// handleGeocode es un proxy server-side hacia Mapbox Geocoding API.
// El frontend no puede llamar a APIs externas directamente por restricciones de CORS,
// pero el backend sí puede. Este endpoint actúa como intermediario.
//
// GET /geocode?address=Calle 98 62-37&city=Bogotá
func handleGeocode(cfg env.IConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Query("address")
		city := c.Query("city")

		if city == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "el campo 'city' es requerido"})
			return
		}

		token := cfg.Get("MAPBOX_ACCESS_TOKEN")
		if token == "" {
			c.JSON(http.StatusOK, GeocodingResult{Found: false})
			return
		}

		// Intento 1: dirección completa
		if address != "" {
			query := fmt.Sprintf("%s, %s, Colombia", address, city)
			lat, lon, ok := mapboxGeocode(query, token)
			if ok {
				c.JSON(http.StatusOK, GeocodingResult{Lat: lat, Lon: lon, Found: true, Fallback: false})
				return
			}
		}

		// Intento 2 (fallback): solo ciudad
		lat, lon, ok := mapboxGeocode(fmt.Sprintf("%s, Colombia", city), token)
		if ok {
			c.JSON(http.StatusOK, GeocodingResult{Lat: lat, Lon: lon, Found: true, Fallback: true})
			return
		}

		c.JSON(http.StatusOK, GeocodingResult{Found: false})
	}
}

// AddressSearchResult representa una sugerencia de dirección.
type AddressSearchResult struct {
	DisplayName   string  `json:"display_name"`
	PlaceID       string  `json:"place_id"`
	Lat           float64 `json:"lat"`
	Lon           float64 `json:"lon"`
	City          string  `json:"city"`
	State         string  `json:"state"`
	Neighbourhood string  `json:"neighbourhood"`
	Postcode      string  `json:"postcode"`
}

// handleAddressSearch retorna un handler que usa Mapbox Geocoding como proxy.
// La API key se lee del config (cargada desde .env), nunca se expone al browser.
//
// GET /address-search?q=avenida+calle+145+128-40+bogota&country=co
func handleAddressSearch(cfg env.IConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		q := c.Query("q")
		if q == "" || len(q) < 8 {
			c.JSON(http.StatusOK, []AddressSearchResult{})
			return
		}

		country := c.DefaultQuery("country", "co")
		city := c.Query("city")
		token := cfg.Get("MAPBOX_ACCESS_TOKEN")
		if token == "" {
			c.JSON(http.StatusOK, []AddressSearchResult{})
			return
		}

		// Si hay ciudad, la usamos como proximity y la añadimos al query
		searchInput := q
		if city != "" {
			searchInput = q + ", " + city
		}

		// Mapbox Geocoding v5 - search endpoint
		geocodeURL := fmt.Sprintf(
			"https://api.mapbox.com/geocoding/v5/mapbox.places/%s.json?access_token=%s&country=%s&language=es&types=address&limit=5",
			url.PathEscape(searchInput),
			token,
			url.QueryEscape(country),
		)

		resp, err := http.Get(geocodeURL)
		if err != nil {
			c.JSON(http.StatusOK, []AddressSearchResult{})
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusOK, []AddressSearchResult{})
			return
		}

		var mapboxResp mapboxFeatureCollection
		if err := json.Unmarshal(body, &mapboxResp); err != nil {
			c.JSON(http.StatusOK, []AddressSearchResult{})
			return
		}

		results := make([]AddressSearchResult, 0, len(mapboxResp.Features))
		for _, feat := range mapboxResp.Features {
			result := AddressSearchResult{
				DisplayName: feat.PlaceName,
				PlaceID:     feat.ID,
			}

			// Mapbox coordinates are [longitude, latitude]
			if len(feat.Center) == 2 {
				result.Lon = feat.Center[0]
				result.Lat = feat.Center[1]
			}

			// Extract address components from context
			for _, ctx := range feat.Context {
				switch {
				case containsType(ctx.ID, "place"):
					result.City = ctx.Text
				case containsType(ctx.ID, "region"):
					result.State = ctx.Text
				case containsType(ctx.ID, "neighborhood") || containsType(ctx.ID, "locality"):
					if result.Neighbourhood == "" {
						result.Neighbourhood = ctx.Text
					}
				case containsType(ctx.ID, "postcode"):
					result.Postcode = ctx.Text
				}
			}

			results = append(results, result)
		}

		c.JSON(http.StatusOK, results)
	}
}

// mapboxFeatureCollection represents the Mapbox Geocoding API response.
type mapboxFeatureCollection struct {
	Features []mapboxFeature `json:"features"`
}

type mapboxFeature struct {
	ID        string          `json:"id"`
	PlaceName string          `json:"place_name"`
	Center    []float64       `json:"center"` // [lon, lat]
	Context   []mapboxContext `json:"context"`
}

type mapboxContext struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// containsType checks if a Mapbox context ID contains a specific type prefix.
// Mapbox context IDs are formatted as "type.id" (e.g., "place.12345", "region.67890").
func containsType(id, typeName string) bool {
	return len(id) > len(typeName) && id[:len(typeName)+1] == typeName+"."
}

// mapboxGeocode performs a forward geocoding search using the Mapbox API.
func mapboxGeocode(query, token string) (float64, float64, bool) {
	endpoint := fmt.Sprintf(
		"https://api.mapbox.com/geocoding/v5/mapbox.places/%s.json?access_token=%s&limit=1&language=es",
		url.PathEscape(query),
		token,
	)

	resp, err := http.Get(endpoint)
	if err != nil || resp.StatusCode != http.StatusOK {
		return 0, 0, false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, false
	}

	var result mapboxFeatureCollection
	if err := json.Unmarshal(body, &result); err != nil || len(result.Features) == 0 {
		return 0, 0, false
	}

	center := result.Features[0].Center
	if len(center) != 2 {
		return 0, 0, false
	}

	// center = [longitude, latitude]
	return center[1], center[0], true
}
