package geocoder

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type Geocoder struct {
	apiKey string
	logger log.ILogger
}

func New(apiKey string, logger log.ILogger) ports.IGeocoder {
	if apiKey == "" {
		return nil
	}
	return &Geocoder{apiKey: apiKey, logger: logger}
}
