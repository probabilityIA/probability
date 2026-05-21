package app

import (
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/infra/secondary/repository"
)

var carrierBaselines = map[string]float64{
	"SERVIENTREGA":     0.92,
	"COORDINADORA":     0.90,
	"ENVIA":            0.88,
	"INTERRAPIDISIMO":  0.89,
	"TCC":              0.87,
	"DEPRISA":          0.88,
	"ENVIOCLICK":       0.85,
	"ENVIAME":          0.86,
	"MIPAQUETE":        0.85,
	"FEDEX":            0.93,
	"DHL":              0.93,
	"UPS":              0.92,
	"99MINUTOS":        0.91,
	"PIBOX":            0.89,
	"SPEED":            0.87,
	"SPEEDCARGO":       0.87,
}

const defaultBaseline = 0.85

func baselineForCarrier(carrier string) float64 {
	key := repository.NormalizeCarrierKey(carrier)
	if v, ok := carrierBaselines[key]; ok {
		return v
	}
	return defaultBaseline
}
