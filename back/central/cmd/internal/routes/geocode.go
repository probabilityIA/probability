package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/shared/env"
)

type GeocodingResult struct {
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Found    bool    `json:"found"`
	Fallback bool    `json:"fallback"`
}

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

		if address != "" {
			query := fmt.Sprintf("%s, %s, Colombia", address, city)
			lat, lon, ok := googleGeocode(query, apiKey, city)
			if ok {
				c.JSON(http.StatusOK, GeocodingResult{Lat: lat, Lon: lon, Found: true, Fallback: false})
				return
			}
		}

		lat, lon, ok := googleGeocode(fmt.Sprintf("%s, Colombia", city), apiKey, "")
		if ok {
			c.JSON(http.StatusOK, GeocodingResult{Lat: lat, Lon: lon, Found: true, Fallback: true})
			return
		}

		c.JSON(http.StatusOK, GeocodingResult{Found: false})
	}
}

type AddressSearchResult struct {
	DisplayName   string  `json:"display_name"`
	PlaceID       string  `json:"place_id"`
	Lat           float64 `json:"lat"`
	Lon           float64 `json:"lon"`
	City          string  `json:"city"`
	State         string  `json:"state"`
	Neighbourhood string  `json:"neighbourhood"`
	Postcode      string  `json:"postcode"`
	DistanceKm    float64 `json:"distance_km,omitempty"`
}

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

		searchInput := q
		if city != "" {
			searchInput = q + ", " + city
		}

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
			displayName := res.FormattedAddress
			displayName = strings.TrimSuffix(displayName, ", Colombia")
			displayName = strings.TrimSuffix(displayName, ",Colombia")
			result := AddressSearchResult{
				DisplayName: displayName,
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

const officeMaxDistanceKm = 30.0

func normalizeForMatch(text string) string {
	replacer := strings.NewReplacer(
		"\u00e1", "a", "\u00e9", "e", "\u00ed", "i", "\u00f3", "o", "\u00fa", "u", "\u00fc", "u", "\u00f1", "n",
		"\u00c1", "a", "\u00c9", "e", "\u00cd", "i", "\u00d3", "o", "\u00da", "u", "\u00dc", "u", "\u00d1", "n",
		" ", "", ".", "", "-", "", "_", "",
	)
	return replacer.Replace(strings.ToLower(text))
}

var cityCenterCache sync.Map

func cityCenter(city, state, apiKey string) (float64, float64, bool) {
	key := strings.ToLower(strings.TrimSpace(city) + "|" + strings.TrimSpace(state))
	if cached, ok := cityCenterCache.Load(key); ok {
		point := cached.([2]float64)
		return point[0], point[1], point[0] != 0 || point[1] != 0
	}

	query := city
	if state != "" {
		query += ", " + state
	}
	lat, lng, ok := googleGeocode(query+", Colombia", apiKey, "")
	if !ok {
		cityCenterCache.Store(key, [2]float64{0, 0})
		return 0, 0, false
	}
	cityCenterCache.Store(key, [2]float64{lat, lng})
	return lat, lng, true
}

func distanceKm(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadiusKm = 6371.0
	toRad := func(deg float64) float64 { return deg * math.Pi / 180 }
	dLat := toRad(lat2 - lat1)
	dLng := toRad(lng2 - lng1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRad(lat1))*math.Cos(toRad(lat2))*math.Sin(dLng/2)*math.Sin(dLng/2)
	return earthRadiusKm * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

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

		city := strings.TrimSpace(c.Query("city"))
		state := strings.TrimSpace(c.Query("state"))
		carrier := normalizeForMatch(c.Query("carrier"))

		var centerLat, centerLng float64
		hasCenter := false
		if city != "" {
			centerLat, centerLng, hasCenter = cityCenter(city, state, apiKey)
		}

		placesURL := fmt.Sprintf(
			"https://maps.googleapis.com/maps/api/place/textsearch/json?query=%s&key=%s&language=es",
			url.QueryEscape(query),
			apiKey,
		)
		if hasCenter {
			placesURL += fmt.Sprintf("&location=%f,%f&radius=%d", centerLat, centerLng, int(officeMaxDistanceKm*1000))
		}

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
		if err := json.Unmarshal(body, &pResp); err != nil {
			fmt.Printf("❌ Places API JSON unmarshal error: %v\n", err)
			fmt.Printf("Response body: %s\n", string(body))
			c.JSON(http.StatusOK, []AddressSearchResult{})
			return
		}

		fmt.Printf("\U0001F4CD Places API Response - Status: %s, Results: %d\n", pResp.Status, len(pResp.Results))
		if pResp.Status != "OK" {
			fmt.Printf("⚠️ Google Places API returned: %s\n", pResp.Status)
			fmt.Printf("Full response: %s\n", string(body))
			c.JSON(http.StatusOK, []AddressSearchResult{})
			return
		}

		results := make([]AddressSearchResult, 0, len(pResp.Results))
		for _, res := range pResp.Results {
			item := AddressSearchResult{
				DisplayName: fmt.Sprintf("%s (%s)", res.Name, res.FormattedAddress),
				PlaceID:     res.PlaceID,
				Lat:         res.Geometry.Location.Lat,
				Lon:         res.Geometry.Location.Lng,
			}

			if carrier != "" && !strings.Contains(normalizeForMatch(res.Name), carrier) {
				continue
			}

			if hasCenter {
				km := distanceKm(centerLat, centerLng, item.Lat, item.Lon)
				if km > officeMaxDistanceKm {
					continue
				}
				item.DistanceKm = math.Round(km*10) / 10
			}

			results = append(results, item)
		}

		if hasCenter {
			sort.Slice(results, func(i, j int) bool {
				return results[i].DistanceKm < results[j].DistanceKm
			})
		}

		c.JSON(http.StatusOK, results)
	}
}

type googleGeocodeResponse struct {
	Status  string                `json:"status"`
	Results []googleGeocodeResult `json:"results"`
}

type googleGeocodeResult struct {
	FormattedAddress string `json:"formatted_address"`
	PlaceID          string `json:"place_id"`
	Geometry         struct {
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

func googleGeocode(query, apiKey, expectedCity string) (float64, float64, bool) {
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

	best := result.Results[0]
	if expectedCity != "" && !resultMatchesCity(best, expectedCity) {
		return 0, 0, false
	}

	loc := best.Geometry.Location
	return loc.Lat, loc.Lng, true
}

func resultMatchesCity(result googleGeocodeResult, expectedCity string) bool {
	expectedCityName := expectedCity
	if idx := strings.Index(expectedCity, ","); idx >= 0 {
		expectedCityName = expectedCity[:idx]
	}
	expectedNorm := normalizeLocationText(expectedCityName)
	if expectedNorm == "" {
		return true
	}

	for _, comp := range result.AddressComponents {
		for _, t := range comp.Types {
			if t == "locality" || t == "administrative_area_level_2" {
				if normalizeLocationText(comp.LongName) == expectedNorm {
					return true
				}
			}
		}
	}
	return false
}

var accentReplacer = strings.NewReplacer(
	"á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u",
	"Á", "A", "É", "E", "Í", "I", "Ó", "O", "Ú", "U",
	"ñ", "n", "Ñ", "N",
)

func normalizeLocationText(text string) string {
	normalized := accentReplacer.Replace(strings.ToUpper(strings.TrimSpace(text)))
	return strings.Join(strings.Fields(normalized), " ")
}
