package handlers

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

type shopifyRateAddress struct {
	Country     string `json:"country"`
	PostalCode  string `json:"postal_code"`
	Province    string `json:"province"`
	City        string `json:"city"`
	Name        string `json:"name"`
	Address1    string `json:"address1"`
	Address2    string `json:"address2"`
	Phone       string `json:"phone"`
	CompanyName string `json:"company_name"`
	Email       string `json:"email"`
}

type shopifyRateItem struct {
	Name     string `json:"name"`
	Sku      string `json:"sku"`
	Quantity int    `json:"quantity"`
	Grams    int    `json:"grams"`
	Price    int64  `json:"price"`
}

type shopifyRateRequest struct {
	Rate struct {
		Origin      shopifyRateAddress `json:"origin"`
		Destination shopifyRateAddress `json:"destination"`
		Items       []shopifyRateItem  `json:"items"`
		Currency    string             `json:"currency"`
		Locale      string             `json:"locale"`
	} `json:"rate"`
}

type shopifyRate struct {
	ServiceName     string `json:"service_name"`
	ServiceCode     string `json:"service_code"`
	TotalPrice      string `json:"total_price"`
	Currency        string `json:"currency"`
	Description     string `json:"description,omitempty"`
	MinDeliveryDate string `json:"min_delivery_date,omitempty"`
	MaxDeliveryDate string `json:"max_delivery_date,omitempty"`
}

func (h *Handlers) ShopifyShippingRates(c *gin.Context) {
	emptyRates := gin.H{"rates": []shopifyRate{}}

	integrationIDStr := c.Param("integration_id")
	integrationID64, err := strconv.ParseUint(integrationIDStr, 10, 64)
	if err != nil || integrationID64 == 0 {
		c.JSON(http.StatusOK, emptyRates)
		return
	}

	var req shopifyRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, emptyRates)
		return
	}

	ctx := c.Request.Context()

	businessID, err := h.uc.Repo().GetIntegrationBusinessID(ctx, uint(integrationID64))
	if err != nil || businessID == 0 {
		c.JSON(http.StatusOK, emptyRates)
		return
	}

	carrier, err := h.carrierResolver.GetActiveShippingCarrier(ctx, businessID)
	if err != nil || carrier == nil {
		c.JSON(http.StatusOK, emptyRates)
		return
	}

	origin, err := h.uc.Repo().GetDefaultWarehouseOrigin(ctx, businessID)
	if err != nil || origin == nil {
		origin, err = h.uc.Repo().GetDefaultOriginAddress(ctx, businessID)
	}
	if err != nil || origin == nil {
		c.JSON(http.StatusOK, emptyRates)
		return
	}

	destDane, _ := h.uc.Repo().GetCityDaneByName(ctx, req.Rate.Destination.City, req.Rate.Destination.Province)
	if destDane == "" {
		c.JSON(http.StatusOK, emptyRates)
		return
	}

	payload := buildShopifyQuotePayload(req, origin, destDane)

	correlationID := uuid.New().String()
	result, err := h.runQuote(ctx, carrier, businessID, payload, correlationID, 8*time.Second)
	if err != nil || result.Status != quoteStatusSuccess {
		c.JSON(http.StatusOK, emptyRates)
		return
	}

	currency := req.Rate.Currency
	if currency == "" {
		currency = "COP"
	}

	ratesList := toRatesList(getRatesFromData(result.Data))

	var quoteID uint
	if len(ratesList) > 0 {
		saved, saveErr := h.uc.Quotes.SaveQuote(ctx, domain.SaveQuoteInput{
			BusinessID:       businessID,
			IntegrationID:    uint(integrationID64),
			Source:           domain.QuoteSourceShopify,
			CorrelationID:    correlationID,
			ExternalOrderRef: req.Rate.Destination.Name,
			RequestPayload:   payload,
			Rates:            ratesList,
		})
		if saveErr == nil && saved != nil {
			quoteID = saved.ID
		}
	}

	rates := mapQuoteRatesToShopify(ratesList, currency, quoteID)
	c.JSON(http.StatusOK, gin.H{"rates": rates})
}

func toRatesList(ratesData interface{}) []map[string]interface{} {
	rawList, ok := ratesData.([]interface{})
	if !ok {
		return nil
	}
	out := make([]map[string]interface{}, 0, len(rawList))
	for _, raw := range rawList {
		if m, ok := raw.(map[string]interface{}); ok {
			out = append(out, m)
		}
	}
	return out
}

func buildShopifyQuotePayload(req shopifyRateRequest, origin *domain.OriginAddress, destDane string) map[string]interface{} {
	dest := req.Rate.Destination

	firstName, lastName := splitName(dest.Name)

	street := strings.TrimSpace(dest.Address1)
	if dest.Address2 != "" {
		street = strings.TrimSpace(street + " " + dest.Address2)
	}

	var totalGrams int
	var contentValue float64
	for _, it := range req.Rate.Items {
		qty := it.Quantity
		if qty <= 0 {
			qty = 1
		}
		totalGrams += it.Grams * qty
		contentValue += float64(it.Price) * float64(qty) / 100.0
	}

	weightKg := float64(totalGrams) / 1000.0
	if weightKg <= 0 {
		weightKg = 1
	}

	pkg := map[string]interface{}{
		"weight": weightKg,
		"height": 10.0,
		"width":  10.0,
		"length": 10.0,
	}

	return map[string]interface{}{
		"requestPickup": false,
		"insurance":     false,
		"description":   "Compra en linea",
		"contentValue":  contentValue,
		"packages":      []interface{}{pkg},
		"origin": map[string]interface{}{
			"company":   origin.Company,
			"firstName": origin.FirstName,
			"lastName":  origin.LastName,
			"email":     origin.Email,
			"phone":     origin.Phone,
			"address":   origin.Street,
			"suburb":    origin.Suburb,
			"daneCode":  origin.CityDaneCode,
		},
		"destination": map[string]interface{}{
			"company":   dest.CompanyName,
			"firstName": firstName,
			"lastName":  lastName,
			"email":     dest.Email,
			"phone":     dest.Phone,
			"address":   street,
			"suburb":    dest.City,
			"daneCode":  destDane,
		},
	}
}

func mapQuoteRatesToShopify(ratesList []map[string]interface{}, currency string, quoteID uint) []shopifyRate {
	out := make([]shopifyRate, 0)

	now := time.Now()

	for i, rate := range ratesList {
		carrierName := toStr(rate["carrier"])
		product := toStr(rate["product"])
		flete := toFloat(rate["flete"])
		if carrierName == "" || flete <= 0 {
			continue
		}

		serviceName := carrierName
		if product != "" {
			serviceName = carrierName + " - " + product
		}

		var serviceCode string
		if quoteID > 0 {
			serviceCode = "pq-" + strconv.FormatUint(uint64(quoteID), 10) + "-" + strconv.Itoa(i)
		} else {
			serviceCode = slugify(carrierName)
			if product != "" {
				serviceCode += "_" + slugify(product)
			}
			serviceCode += "_" + strconv.Itoa(i)
		}

		minimumInsurance := toFloat(rate["minimumInsurance"])
		totalPriceCents := int64(math.Round((flete + minimumInsurance) * 100))

		sr := shopifyRate{
			ServiceName: serviceName,
			ServiceCode: serviceCode,
			TotalPrice:  strconv.FormatInt(totalPriceCents, 10),
			Currency:    currency,
		}

		if days := int(toFloat(rate["deliveryDays"])); days > 0 {
			sr.Description = "Entrega estimada " + strconv.Itoa(days) + " dias habiles"
			sr.MinDeliveryDate = now.AddDate(0, 0, days).Format("2006-01-02 15:04:05 -0700")
			sr.MaxDeliveryDate = now.AddDate(0, 0, days+2).Format("2006-01-02 15:04:05 -0700")
		}

		out = append(out, sr)
	}

	return out
}

func splitName(full string) (string, string) {
	full = strings.TrimSpace(full)
	if full == "" {
		return "Cliente", "Shopify"
	}
	parts := strings.Fields(full)
	if len(parts) == 1 {
		return parts[0], parts[0]
	}
	return parts[0], strings.Join(parts[1:], " ")
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == ' ' || r == '-' || r == '_':
			b.WriteRune('_')
		}
	}
	out := b.String()
	if out == "" {
		return "carrier"
	}
	return out
}

func toStr(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func toFloat(v interface{}) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case float32:
		return float64(n)
	case int:
		return float64(n)
	case int64:
		return float64(n)
	case string:
		f, _ := strconv.ParseFloat(n, 64)
		return f
	}
	return 0
}
