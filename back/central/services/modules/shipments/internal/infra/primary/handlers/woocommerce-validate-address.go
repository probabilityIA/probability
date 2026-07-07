package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type validateAddressRequest struct {
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
}

type validateAddressResponse struct {
	Found      bool    `json:"found"`
	Lat        float64 `json:"lat"`
	Lng        float64 `json:"lng"`
	City       string  `json:"city"`
	Department string  `json:"department"`
	DaneCode   string  `json:"dane_code"`
	Formatted  string  `json:"formatted_address"`
	Confidence string  `json:"confidence"`
}

func (h *Handlers) WooCommerceValidateAddress(c *gin.Context) {
	if _, ok := h.authWooPublic(c); !ok {
		return
	}

	var req validateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, validateAddressResponse{Found: false, Confidence: "low"})
		return
	}

	ctx := c.Request.Context()
	city := strings.TrimSpace(req.City)
	department := strings.TrimSpace(req.State)
	daneCode := ""

	if h.geocoder != nil {
		query := strings.TrimSpace(strings.Join([]string{req.Address, req.City, req.State, "Colombia"}, ", "))
		if geo, err := h.geocoder.Geocode(ctx, query); err == nil && geo.Found {
			resolvedCity := firstNonEmptyStr(geo.Locality, geo.AdminArea2, req.City)
			resolvedDept := firstNonEmptyStr(geo.Department, req.State)

			dane := h.daneCached(ctx, resolvedCity, resolvedDept)
			confidence := "low"
			if dane != "" {
				city = resolvedCity
				department = resolvedDept
				daneCode = dane
				switch {
				case geo.LocationType == "ROOFTOP" && !geo.PartialMatch:
					confidence = "high"
				case geo.PartialMatch:
					confidence = "low"
				default:
					confidence = "medium"
				}
			}

			c.JSON(http.StatusOK, validateAddressResponse{
				Found:      true,
				Lat:        geo.Lat,
				Lng:        geo.Lng,
				City:       city,
				Department: department,
				DaneCode:   daneCode,
				Formatted:  geo.FormattedAddress,
				Confidence: confidence,
			})
			return
		}
	}

	daneCode = h.daneCached(ctx, city, department)
	confidence := "low"
	if daneCode != "" {
		confidence = "medium"
	}
	c.JSON(http.StatusOK, validateAddressResponse{
		Found:      daneCode != "",
		City:       city,
		Department: department,
		DaneCode:   daneCode,
		Confidence: confidence,
	})
}

func firstNonEmptyStr(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
