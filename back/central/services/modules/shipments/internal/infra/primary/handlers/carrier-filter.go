package handlers

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

func (h *Handlers) filterRatesByBusinessCarriers(ctx context.Context, businessID uint, warehouseID *uint, isCOD bool, ratesList []map[string]interface{}) []map[string]interface{} {
	if businessID == 0 || len(ratesList) == 0 {
		return ratesList
	}

	settings, err := h.uc.Repo().GetBusinessCarrierSettings(ctx, businessID, warehouseID)
	if err != nil || len(settings) == 0 {
		return ratesList
	}

	out := make([]map[string]interface{}, 0, len(ratesList))
	for _, rate := range ratesList {
		if domain.CarrierAllowed(settings, toStr(rate["carrier"]), isCOD) {
			out = append(out, rate)
		}
	}
	return out
}
