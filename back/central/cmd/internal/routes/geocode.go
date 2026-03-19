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

// handleGeocode es un proxy server-side hacia Nominatim.
// El frontend no puede llamar a Nominatim directamente por restricciones de CORS/User-Agent,
// pero el backend sí puede. Este endpoint actúa como intermediario.
//
// GET /geocode?address=Calle 98 62-37&city=Bogotá
func handleGeocode(c *gin.Context) {
	address := c.Query("address")
	city := c.Query("city")

	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "el campo 'city' es requerido"})
		return
	}

	// Intento 1: dirección completa
	if address != "" {
		query := fmt.Sprintf("%s, %s, Colombia", address, city)
		lat, lon, ok := nominatimSearch(query)
		if ok {
			c.JSON(http.StatusOK, GeocodingResult{Lat: lat, Lon: lon, Found: true, Fallback: false})
			return
		}
	}

	// Intento 2 (fallback): solo ciudad
	lat, lon, ok := nominatimSearch(fmt.Sprintf("%s, Colombia", city))
	if ok {
		c.JSON(http.StatusOK, GeocodingResult{Lat: lat, Lon: lon, Found: true, Fallback: true})
		return
	}

	c.JSON(http.StatusOK, GeocodingResult{Found: false})
}

// AddressSearchResult representa una sugerencia de dirección.
type AddressSearchResult struct {
	DisplayName string  `json:"display_name"`
	PlaceID     string  `json:"place_id"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
}

// handleAddressSearch retorna un handler que usa Google Places Autocomplete como proxy.
// La API key se lee del config (cargada desde .env), nunca se expone al browser.
//
// GET /address-search?q=avenida+calle+145+128-40+bogota&country=co
func handleAddressSearch(cfg env.IConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		q := c.Query("q")
		if q == "" || len(q) < 4 {
			c.JSON(http.StatusOK, []AddressSearchResult{})
			return
		}

		country := c.DefaultQuery("country", "co")
		apiKey := cfg.Get("GOOGLE_MAPS_API_KEY")
		if apiKey == "" {
			c.JSON(http.StatusOK, []AddressSearchResult{})
			return
		}

	// Step 1: Google Places Autocomplete
	autocompleteURL := fmt.Sprintf(
		"https://maps.googleapis.com/maps/api/place/autocomplete/json?input=%s&components=country:%s&language=es&types=address&key=%s",
		url.QueryEscape(q),
		url.QueryEscape(country),
		apiKey,
	)

	resp, err := http.Get(autocompleteURL)
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

	var autocompleteResp struct {
		Status      string `json:"status"`
		Predictions []struct {
			Description string `json:"description"`
			PlaceID     string `json:"place_id"`
		} `json:"predictions"`
	}

	if err := json.Unmarshal(body, &autocompleteResp); err != nil || autocompleteResp.Status != "OK" {
		c.JSON(http.StatusOK, []AddressSearchResult{})
		return
	}

	// Step 2: For each prediction, get coordinates via Place Details (only first 5)
	results := make([]AddressSearchResult, 0, len(autocompleteResp.Predictions))
	limit := len(autocompleteResp.Predictions)
	if limit > 5 {
		limit = 5
	}

	for _, pred := range autocompleteResp.Predictions[:limit] {
		result := AddressSearchResult{
			DisplayName: pred.Description,
			PlaceID:     pred.PlaceID,
		}

		// Get coordinates from Place Details
		detailsURL := fmt.Sprintf(
			"https://maps.googleapis.com/maps/api/place/details/json?place_id=%s&fields=geometry&key=%s",
			url.QueryEscape(pred.PlaceID),
			apiKey,
		)
		detResp, err := http.Get(detailsURL)
		if err == nil {
			detBody, _ := io.ReadAll(detResp.Body)
			detResp.Body.Close()
			var details struct {
				Result struct {
					Geometry struct {
						Location struct {
							Lat float64 `json:"lat"`
							Lng float64 `json:"lng"`
						} `json:"location"`
					} `json:"geometry"`
				} `json:"result"`
			}
			if json.Unmarshal(detBody, &details) == nil {
				result.Lat = details.Result.Geometry.Location.Lat
				result.Lon = details.Result.Geometry.Location.Lng
			}
		}

		results = append(results, result)
	}

	c.JSON(http.StatusOK, results)
	}
}

// nominatimSearch realiza la búsqueda en Nominatim y retorna lat, lon y si encontró resultados.
func nominatimSearch(query string) (float64, float64, bool) {
	endpoint := fmt.Sprintf(
		"https://nominatim.openstreetmap.org/search?format=json&addressdetails=1&limit=1&q=%s",
		url.QueryEscape(query),
	)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, 0, false
	}
	// Nominatim requiere un User-Agent válido con app/contacto
	req.Header.Set("User-Agent", "ProbabilityApp/1.0 (contact@probability.com.co)")
	req.Header.Set("Accept-Language", "es")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return 0, 0, false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, false
	}

	var results []struct {
		Lat string `json:"lat"`
		Lon string `json:"lon"`
	}
	if err := json.Unmarshal(body, &results); err != nil || len(results) == 0 {
		return 0, 0, false
	}

	var lat, lon float64
	fmt.Sscanf(results[0].Lat, "%f", &lat)
	fmt.Sscanf(results[0].Lon, "%f", &lon)
	return lat, lon, true
}
