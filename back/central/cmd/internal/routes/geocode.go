package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
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
