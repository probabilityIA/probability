package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

type wooRateDestination struct {
	Country  string `json:"country"`
	State    string `json:"state"`
	City     string `json:"city"`
	Postcode string `json:"postcode"`
	Address1 string `json:"address_1"`
	Address2 string `json:"address_2"`
	Name     string `json:"name"`
	Company  string `json:"company"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

type wooRateItem struct {
	Name        string  `json:"name"`
	Sku         string  `json:"sku"`
	Quantity    int     `json:"quantity"`
	WeightGrams float64 `json:"weight_grams"`
	Price       float64 `json:"price"`
}

type wooRateRequest struct {
	Destination wooRateDestination `json:"destination"`
	Contents    []wooRateItem      `json:"contents"`
	Currency    string             `json:"currency"`
}

type wooRate struct {
	ID           string                 `json:"id"`
	Label        string                 `json:"label"`
	Cost         string                 `json:"cost"`
	Currency     string                 `json:"currency"`
	DeliveryDays int                    `json:"delivery_days,omitempty"`
	MetaData     map[string]interface{} `json:"meta_data,omitempty"`
}

func (h *Handlers) WooCommerceShippingRates(c *gin.Context) {
	emptyRates := gin.H{"rates": []wooRate{}}

	integrationIDStr := c.Param("integration_id")
	integrationID64, err := strconv.ParseUint(integrationIDStr, 10, 64)
	if err != nil || integrationID64 == 0 {
		c.JSON(http.StatusOK, emptyRates)
		return
	}

	ctx := c.Request.Context()

	resolved, err := h.resolveWoo(ctx, uint(integrationID64))
	if err != nil {
		c.JSON(http.StatusOK, emptyRates)
		return
	}

	if !resolved.Found || resolved.Revoked || !h.wooTokenMatches(uint(integrationID64), resolved.Salt, c.GetHeader("X-Probability-Token")) {
		c.JSON(http.StatusUnauthorized, gin.H{"rates": []wooRate{}, "error": "invalid_token"})
		return
	}

	var req wooRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, emptyRates)
		return
	}

	if resolved.FreeShippingEnabled && resolved.FreeShippingMin > 0 {
		var subtotal float64
		for _, it := range req.Contents {
			qty := it.Quantity
			if qty <= 0 {
				qty = 1
			}
			subtotal += it.Price * float64(qty)
		}
		if subtotal >= resolved.FreeShippingMin {
			currency := req.Currency
			if currency == "" {
				currency = "COP"
			}
			c.JSON(http.StatusOK, gin.H{"rates": []wooRate{{
				ID:       "probability_free_shipping",
				Label:    "Envio gratis",
				Cost:     "0",
				Currency: currency,
				MetaData: map[string]interface{}{"free_shipping": true, "threshold": resolved.FreeShippingMin},
			}}})
			return
		}
	}

	if resolved.BusinessID == 0 || resolved.Carrier == nil || resolved.Origin == nil {
		c.JSON(http.StatusOK, emptyRates)
		return
	}

	businessID := resolved.BusinessID
	carrier := resolved.Carrier

	destDane := h.daneCached(ctx, req.Destination.City, req.Destination.State)
	if destDane == "" {
		c.JSON(http.StatusOK, emptyRates)
		return
	}

	payload := buildWooQuotePayload(req, resolved.Origin, normalizeDaneCode(destDane))

	correlationID := uuid.New().String()
	result, err := h.runQuote(ctx, carrier, businessID, payload, correlationID, 12*time.Second)
	if err != nil || result.Status != quoteStatusSuccess {
		c.JSON(http.StatusOK, emptyRates)
		return
	}

	currency := req.Currency
	if currency == "" {
		currency = "COP"
	}

	ratesList := toRatesList(getRatesFromData(result.Data))

	var quoteID uint
	if len(ratesList) > 0 {
		saved, saveErr := h.uc.Quotes.SaveQuote(ctx, domain.SaveQuoteInput{
			BusinessID:       businessID,
			IntegrationID:    uint(integrationID64),
			Source:           domain.QuoteSourceWooCommerce,
			CorrelationID:    correlationID,
			ExternalOrderRef: req.Destination.Name,
			RequestPayload:   payload,
			Rates:            ratesList,
		})
		if saveErr == nil && saved != nil {
			quoteID = saved.ID
		}
	}

	rates := mapQuoteRatesToWoo(ratesList, currency, quoteID)
	c.JSON(http.StatusOK, gin.H{"rates": rates})
}

func normalizeDaneCode(code string) string {
	code = strings.TrimSpace(code)
	if l := len(code); l >= 5 && l < 8 {
		return code + strings.Repeat("0", 8-l)
	}
	return code
}

func buildWooQuotePayload(req wooRateRequest, origin *domain.OriginAddress, destDane string) map[string]interface{} {
	dest := req.Destination

	firstName, lastName := splitName(dest.Name)

	street := strings.TrimSpace(dest.Address1)
	if dest.Address2 != "" {
		street = strings.TrimSpace(street + " " + dest.Address2)
	}

	var totalGrams float64
	var contentValue float64
	for _, it := range req.Contents {
		qty := it.Quantity
		if qty <= 0 {
			qty = 1
		}
		totalGrams += it.WeightGrams * float64(qty)
		contentValue += it.Price * float64(qty)
	}

	weightKg := totalGrams / 1000.0
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
			"company":   dest.Company,
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

func mapQuoteRatesToWoo(ratesList []map[string]interface{}, currency string, quoteID uint) []wooRate {
	out := make([]wooRate, 0)

	for i, rate := range ratesList {
		carrierName := toStr(rate["carrier"])
		product := toStr(rate["product"])
		flete := toFloat(rate["flete"])
		if carrierName == "" || flete <= 0 {
			continue
		}

		label := carrierName
		if product != "" {
			label = carrierName + " - " + product
		}

		var id string
		if quoteID > 0 {
			id = "pq-" + strconv.FormatUint(uint64(quoteID), 10) + "-" + strconv.Itoa(i)
		} else {
			id = slugify(carrierName)
			if product != "" {
				id += "_" + slugify(product)
			}
			id += "_" + strconv.Itoa(i)
		}

		wr := wooRate{
			ID:       id,
			Label:    label,
			Cost:     strconv.FormatFloat(flete, 'f', -1, 64),
			Currency: currency,
			MetaData: map[string]interface{}{
				"quote_id":     quoteID,
				"rate_index":   i,
				"carrier":      carrierName,
				"product":      product,
				"service_code": toStr(rate["serviceCode"]),
				"id_rate":      rate["idRate"],
			},
		}

		if days := int(toFloat(rate["deliveryDays"])); days > 0 {
			wr.DeliveryDays = days
		}

		out = append(out, wr)
	}

	return out
}
