package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *ProductRepository) GetShippingQuoteRate(ctx context.Context, quoteID uint, rateIndex int) (*domain.ShippingQuoteRate, error) {
	var quote models.ShippingQuote
	if err := r.db.Conn(ctx).First(&quote, quoteID).Error; err != nil {
		return nil, err
	}

	var rates []map[string]interface{}
	if err := json.Unmarshal(quote.Rates, &rates); err != nil {
		return nil, err
	}
	if rateIndex < 0 || rateIndex >= len(rates) {
		return nil, fmt.Errorf("rate_index %d fuera de rango para shipping_quote %d", rateIndex, quoteID)
	}

	rate := rates[rateIndex]
	cod, _ := rate["cod"].(bool)

	return &domain.ShippingQuoteRate{
		Flete:                toFloatJSON(rate["flete"]),
		MinimumInsurance:     toFloatJSON(rate["minimumInsurance"]),
		COD:                  cod,
		CODCarrierFee:        toFloatJSON(rate["codCarrierFee"]),
		CODProbabilityMargin: toFloatJSON(rate["codProbabilityMargin"]),
	}, nil
}

func toFloatJSON(v interface{}) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case int:
		return float64(n)
	case int64:
		return float64(n)
	}
	return 0
}
