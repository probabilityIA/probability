package usecases

import "strings"

const (
	probabilityWeightUnit    = "kg"
	probabilityDimensionUnit = "cm"
)

var weightToKg = map[string]float64{
	"kg":        1,
	"kgs":       1,
	"kilo":      1,
	"kilos":     1,
	"kilogram":  1,
	"kilograms": 1,
	"g":         0.001,
	"gr":        0.001,
	"grs":       0.001,
	"gram":      0.001,
	"grams":     0.001,
	"lb":        0.45359237,
	"lbs":       0.45359237,
	"pound":     0.45359237,
	"pounds":    0.45359237,
	"oz":        0.028349523125,
	"ounce":     0.028349523125,
	"ounces":    0.028349523125,
}

func weightFactor(storeUnit string) (float64, bool) {
	unit := strings.ToLower(strings.TrimSpace(storeUnit))
	if unit == "" {
		return 0, false
	}
	factor, known := weightToKg[unit]
	return factor, known
}

func normalizeWeightToKg(weight float64, storeUnit string) (float64, bool) {
	if weight <= 0 {
		return 0, false
	}
	factor, known := weightFactor(storeUnit)
	if !known {
		return 0, false
	}
	return weight * factor, true
}

func convertKgToStoreUnit(weightKg float64, storeUnit string) (float64, bool) {
	if weightKg <= 0 {
		return 0, false
	}
	factor, known := weightFactor(storeUnit)
	if !known || factor == 0 {
		return 0, false
	}
	return weightKg / factor, true
}

func positive(value float64) *float64 {
	if value <= 0 {
		return nil
	}
	v := value
	return &v
}
