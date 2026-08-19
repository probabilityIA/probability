package usecases

import (
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
)

const (
	probabilityWeightUnit    = "kg"
	probabilityDimensionUnit = "cm"
)

var weightFactors = map[string]float64{
	"kg":       1,
	"kgs":      1,
	"kilogram": 1,
	"g":        0.001,
	"gr":       0.001,
	"gram":     0.001,
	"lb":       0.45359237,
	"lbs":      0.45359237,
	"pound":    0.45359237,
	"oz":       0.028349523125,
}

func weightFactor(unit string) (float64, bool) {
	factor, ok := weightFactors[strings.ToLower(strings.TrimSpace(unit))]
	return factor, ok
}

func positive(value float64) *float64 {
	if value <= 0 {
		return nil
	}
	v := value
	return &v
}

func normalizeSKU(sku string) string {
	return strings.ToLower(strings.TrimSpace(sku))
}

func parseExternalProductID(externalID string) (int64, int64, error) {
	productPart, variantPart, hasVariant := strings.Cut(strings.TrimSpace(externalID), ":")

	productID, err := strconv.ParseInt(productPart, 10, 64)
	if err != nil {
		return 0, 0, domain.ErrIntegrationNotFound
	}

	var variantID int64
	if hasVariant && variantPart != "" {
		variantID, err = strconv.ParseInt(variantPart, 10, 64)
		if err != nil {
			return 0, 0, domain.ErrIntegrationNotFound
		}
	}

	return productID, variantID, nil
}
