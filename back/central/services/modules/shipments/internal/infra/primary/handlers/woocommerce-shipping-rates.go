package handlers

import (
	"context"
	"net/http"
	"net/url"
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

const defaultPackageDimCm = 10.0

type wooRateRequest struct {
	Destination wooRateDestination `json:"destination"`
	Contents    []wooRateItem      `json:"contents"`
	Currency    string             `json:"currency"`
	COD         bool               `json:"cod"`
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

	if req.COD && resolved.CODQuotingDisabled {
		c.JSON(http.StatusOK, gin.H{"rates": []wooRate{}, "cod_blocked": true})
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

	pkg := h.resolveWooPackageDimensions(ctx, businessID, resolved, req)

	payload := buildWooQuotePayload(req, resolved.Origin, normalizeDaneCode(destDane), pkg)

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

	allowedCarriers := resolved.AllowedCarriersPrepaid
	if req.COD {
		allowedCarriers = resolved.AllowedCarriersCOD
	}
	ratesList = filterRatesByAllowedCarriers(ratesList, allowedCarriers)

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

	rates := mapQuoteRatesToWoo(ratesList, currency, quoteID, h.pluginBaseURL, req.COD)

	if req.COD && len(rates) == 0 && len(ratesList) > 0 {
		c.JSON(http.StatusOK, gin.H{"rates": []wooRate{}, "cod_blocked": true})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rates": rates})
}

// filterRatesByAllowedCarriers deja pasar solo las tarifas cuyo carrier esta en la
// lista permitida. Si la lista viene vacia (el cliente no configuro un filtro para
// ese entorno), no se filtra nada y se mantiene el comportamiento actual.
func filterRatesByAllowedCarriers(ratesList []map[string]interface{}, allowed []string) []map[string]interface{} {
	if len(allowed) == 0 {
		return ratesList
	}
	allowedSet := make(map[string]bool, len(allowed))
	for _, c := range allowed {
		allowedSet[strings.ToUpper(strings.TrimSpace(c))] = true
	}
	out := make([]map[string]interface{}, 0, len(ratesList))
	for _, rate := range ratesList {
		carrierName := strings.ToUpper(strings.TrimSpace(toStr(rate["carrier"])))
		if allowedSet[carrierName] {
			out = append(out, rate)
		}
	}
	return out
}

func normalizeDaneCode(code string) string {
	code = strings.TrimSpace(code)
	if l := len(code); l >= 5 && l < 8 {
		return code + strings.Repeat("0", 8-l)
	}
	return code
}

type wooPackageDims struct {
	Weight float64
	Length float64
	Width  float64
	Height float64
}

func (h *Handlers) resolveWooPackageDimensions(ctx context.Context, businessID uint, resolved *wooResolved, req wooRateRequest) wooPackageDims {
	var totalGrams float64
	var totalQuantity int
	skus := make([]string, 0, len(req.Contents))
	for _, it := range req.Contents {
		qty := it.Quantity
		if qty <= 0 {
			qty = 1
		}
		totalGrams += it.WeightGrams * float64(qty)
		totalQuantity += qty
		if it.Sku != "" {
			skus = append(skus, it.Sku)
		}
	}

	weightKg := totalGrams / 1000.0
	if weightKg <= 0 {
		weightKg = 1
	}

	var maxLength, maxWidth, maxHeight float64
	if businessID > 0 && len(skus) > 0 {
		dims, err := h.uc.Repo().GetProductDimensionsBySKUs(ctx, businessID, skus)
		if err == nil {
			for _, d := range dims {
				if d.Length != nil && *d.Length > maxLength {
					maxLength = *d.Length
				}
				if d.Width != nil && *d.Width > maxWidth {
					maxWidth = *d.Width
				}
				if d.Height != nil && *d.Height > maxHeight {
					maxHeight = *d.Height
				}
			}
		}
	}

	if resolved.OriginIsWarehouse && resolved.PackageConfig != nil &&
		resolved.PackageConfig.Strategy == domain.ShippingPackageStrategyStandardBox {
		if box := resolved.PackageConfig.SelectBox(totalQuantity, maxLength, maxWidth, maxHeight); box != nil {
			out := wooPackageDims{Weight: weightKg, Length: maxLength, Width: maxWidth, Height: maxHeight}
			if box.Weight != nil {
				out.Weight = *box.Weight
			}
			if box.Length != nil {
				out.Length = *box.Length
			}
			if box.Width != nil {
				out.Width = *box.Width
			}
			if box.Height != nil {
				out.Height = *box.Height
			}
			if out.Length > 0 && out.Width > 0 && out.Height > 0 {
				return out
			}
		}
	}

	if maxLength > 0 && maxWidth > 0 && maxHeight > 0 {
		return wooPackageDims{Weight: weightKg, Length: maxLength, Width: maxWidth, Height: maxHeight}
	}

	return wooPackageDims{Weight: weightKg, Length: defaultPackageDimCm, Width: defaultPackageDimCm, Height: defaultPackageDimCm}
}

func buildWooQuotePayload(req wooRateRequest, origin *domain.OriginAddress, destDane string, pkgDims wooPackageDims) map[string]interface{} {
	dest := req.Destination

	firstName, lastName := splitName(dest.Name)

	street := strings.TrimSpace(dest.Address1)
	if dest.Address2 != "" {
		street = strings.TrimSpace(street + " " + dest.Address2)
	}

	var contentValue float64
	for _, it := range req.Contents {
		qty := it.Quantity
		if qty <= 0 {
			qty = 1
		}
		contentValue += it.Price * float64(qty)
	}

	pkg := map[string]interface{}{
		"weight": pkgDims.Weight,
		"height": pkgDims.Height,
		"width":  pkgDims.Width,
		"length": pkgDims.Length,
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

func mapQuoteRatesToWoo(ratesList []map[string]interface{}, currency string, quoteID uint, logoBaseURL string, isCOD bool) []wooRate {
	out := make([]wooRate, 0)

	for i, rate := range ratesList {
		carrierName := toStr(rate["carrier"])
		product := toStr(rate["product"])
		flete := toFloat(rate["flete"])
		if carrierName == "" || flete <= 0 {
			continue
		}

		minimumInsurance := toFloat(rate["minimumInsurance"])
		cost := flete + minimumInsurance
		codCarrierFee := 0.0
		codProbabilityMargin := 0.0
		if isCOD {
			supportsCOD, _ := rate["cod"].(bool)
			if !supportsCOD {
				continue
			}
			codCarrierFee = toFloat(rate["codCarrierFee"])
			codProbabilityMargin = toFloat(rate["codProbabilityMargin"])
			cost = flete + minimumInsurance + codCarrierFee + codProbabilityMargin
		}

		logoURL := ""
		if logoBaseURL != "" {
			logoURL = strings.TrimRight(logoBaseURL, "/") + "/api/v1/woocommerce/carrier-logo/" + url.PathEscape(carrierName)
		}

		displayCarrier := accentEsWords(carrierName)
		displayProduct := accentEsWords(product)
		label := displayCarrier
		if product != "" {
			label = displayCarrier + " - " + displayProduct
		}
		if days := int(toFloat(rate["deliveryDays"])); days > 0 {
			label += " (" + strconv.Itoa(days) + " d\u00edas h\u00e1biles)"
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

		meta := map[string]interface{}{
			"quote_id":     quoteID,
			"rate_index":   i,
			"carrier":      carrierName,
			"product":      product,
			"service_code": toStr(rate["serviceCode"]),
			"id_rate":      rate["idRate"],
			"logo_url":     logoURL,
			"cod":          isCOD,
			"flete":        flete,
			"insurance":    minimumInsurance,
		}
		if isCOD {
			meta["cod_carrier_fee"] = codCarrierFee
			meta["cod_probability_margin"] = codProbabilityMargin
		}

		wr := wooRate{
			ID:       id,
			Label:    label,
			Cost:     strconv.FormatFloat(cost, 'f', -1, 64),
			Currency: currency,
			MetaData: meta,
		}

		out = append(out, wr)
	}

	return out
}

var esAccentMap = map[string]string{
	"envia":           "Env\u00eda",
	"interrapidisimo": "Interrapid\u00edsimo",
	"mercancia":       "Mercanc\u00eda",
	"estandar":        "Est\u00e1ndar",
	"paqueteria":      "Paqueter\u00eda",
	"dia":             "D\u00eda",
	"economico":       "Econ\u00f3mico",
	"logistica":       "Log\u00edstica",
	"mensajeria":      "Mensajer\u00eda",
	"rapido":          "R\u00e1pido",
	"aereo":           "A\u00e9reo",
}

func accentEsWords(s string) string {
	if s == "" {
		return s
	}
	words := strings.Fields(s)
	for i, w := range words {
		if v, ok := esAccentMap[strings.ToLower(w)]; ok {
			words[i] = v
		}
	}
	return strings.Join(words, " ")
}
