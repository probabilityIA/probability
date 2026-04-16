// Package routes defines HTTP route handlers for geocoding and address search.
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

// handleGeocode es un proxy server-side hacia Google Geocoding API.
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

		apiKey := cfg.Get("GOOGLE_MAPS_API_KEY")
		if apiKey == "" {
			c.JSON(http.StatusOK, GeocodingResult{Found: false})
			return
		}

		// Intento 1: dirección completa
		if address != "" {
			query := fmt.Sprintf("%s, %s, Colombia", address, city)
			lat, lon, ok := googleGeocode(query, apiKey)
			if ok {
				c.JSON(http.StatusOK, GeocodingResult{Lat: lat, Lon: lon, Found: true, Fallback: false})
				return
			}
		}

		// Intento 2 (fallback): solo ciudad
		lat, lon, ok := googleGeocode(fmt.Sprintf("%s, Colombia", city), apiKey)
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

// handleAddressSearch retorna un handler que usa Google Geocoding API como proxy.
// Una sola llamada HTTP por búsqueda (más eficiente que Places Autocomplete + Details).
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
		apiKey := cfg.Get("GOOGLE_MAPS_API_KEY")
		if apiKey == "" {
			c.JSON(http.StatusOK, []AddressSearchResult{})
			return
		}

		// Si hay ciudad, la añadimos al query para mejor contexto
		searchInput := q
		if city != "" {
			searchInput = q + ", " + city
		}

		// Google Geocoding API - una sola llamada, devuelve coordenadas + componentes
		geocodeURL := fmt.Sprintf(
			"https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s&language=es&components=country:%s",
			url.QueryEscape(searchInput),
			apiKey,
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

		var geoResp googleGeocodeResponse
		if err := json.Unmarshal(body, &geoResp); err != nil || geoResp.Status != "OK" {
			c.JSON(http.StatusOK, []AddressSearchResult{})
			return
		}

		results := make([]AddressSearchResult, 0, len(geoResp.Results))
		limit := len(geoResp.Results)
		if limit > 5 {
			limit = 5
		}

		for _, res := range geoResp.Results[:limit] {
			result := AddressSearchResult{
				DisplayName: res.FormattedAddress,
				PlaceID:     res.PlaceID,
				Lat:         res.Geometry.Location.Lat,
				Lon:         res.Geometry.Location.Lng,
			}

			for _, comp := range res.AddressComponents {
				for _, t := range comp.Types {
					switch t {
					case "locality":
						result.City = comp.LongName
					case "administrative_area_level_1":
						result.State = comp.LongName
					case "neighborhood", "sublocality_level_1", "sublocality":
						if result.Neighbourhood == "" {
							result.Neighbourhood = comp.LongName
						}
					case "postal_code":
						result.Postcode = comp.LongName
					}
				}
			}

			results = append(results, result)
		}

		c.JSON(http.StatusOK, results)
	}
}

// placesSearchResponse represents the Google Places API response.
type placesSearchResponse struct {
	Status  string               `json:"status"`
	Results []placesSearchResult `json:"results"`
}

type placesSearchResult struct {
	Name             string `json:"name"`
	FormattedAddress string `json:"formatted_address"`
	PlaceID          string `json:"place_id"`
	Geometry         struct {
		Location struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"location"`
	} `json:"geometry"`
}

// handlePlacesSearch proxies a text search to Google Places API
// GET /api/v1/places-search?query=Oficina+Coordinadora+Bogota
func handlePlacesSearch(cfg env.IConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("query")
		if query == "" {
			c.JSON(http.StatusOK, []AddressSearchResult{})
			return
		}

		apiKey := cfg.Get("GOOGLE_MAPS_API_KEY")
		if apiKey == "" {
			c.JSON(http.StatusOK, []AddressSearchResult{})
			return
		}

		placesURL := fmt.Sprintf(
			"https://maps.googleapis.com/maps/api/place/textsearch/json?query=%s&key=%s&language=es",
			url.QueryEscape(query),
			apiKey,
		)

		resp, err := http.Get(placesURL)
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

		var pResp placesSearchResponse
		if err := json.Unmarshal(body, &pResp); err != nil || pResp.Status != "OK" {
			c.JSON(http.StatusOK, []AddressSearchResult{})
			return
		}

		results := make([]AddressSearchResult, 0, len(pResp.Results))
		for _, res := range pResp.Results {
			results = append(results, AddressSearchResult{
				DisplayName: fmt.Sprintf("%s (%s)", res.Name, res.FormattedAddress),
				PlaceID:     res.PlaceID,
				Lat:         res.Geometry.Location.Lat,
				Lon:         res.Geometry.Location.Lng,
			})
		}

		c.JSON(http.StatusOK, results)
	}
}


// googleGeocodeResponse represents the Google Geocoding API response.
type googleGeocodeResponse struct {
	Status  string              `json:"status"`
	Results []googleGeocodeResult `json:"results"`
}

type googleGeocodeResult struct {
	FormattedAddress  string `json:"formatted_address"`
	PlaceID           string `json:"place_id"`
	Geometry          struct {
		Location struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"location"`
	} `json:"geometry"`
	AddressComponents []struct {
		LongName  string   `json:"long_name"`
		ShortName string   `json:"short_name"`
		Types     []string `json:"types"`
	} `json:"address_components"`
}

// googleGeocode performs a forward geocoding search using the Google Geocoding API.
func googleGeocode(query, apiKey string) (float64, float64, bool) {
	endpoint := fmt.Sprintf(
		"https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s&language=es",
		url.QueryEscape(query),
		apiKey,
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

	var result googleGeocodeResponse
	if err := json.Unmarshal(body, &result); err != nil || result.Status != "OK" || len(result.Results) == 0 {
		return 0, 0, false
	}

	loc := result.Results[0].Geometry.Location
	return loc.Lat, loc.Lng, true
}
