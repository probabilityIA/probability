package usecasecreateorder

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

func (uc *UseCaseCreateOrder) geocodeOrderIfNeeded(ctx context.Context, order *entities.ProbabilityOrder) {
	if order.ShippingLat != nil && order.ShippingLng != nil {
		return
	}
	if uc.geocoder == nil {
		return
	}

	query := buildGeocodeQuery(order.ShippingStreet, order.ShippingCity, order.ShippingState)
	if query == "" {
		return
	}

	res := uc.geocoder.GeocodeDetailed(ctx, query)
	if !res.Found {
		order.ShippingGeoConfidence = "low"
		return
	}

	order.ShippingLat = &res.Lat
	order.ShippingLng = &res.Lng
	order.ShippingGeoConfidence = geoConfidence(res.LocationType, res.PartialMatch)

	uc.logger.Info(ctx).
		Str("order_id", order.ID).
		Float64("lat", res.Lat).
		Float64("lng", res.Lng).
		Str("confidence", order.ShippingGeoConfidence).
		Msg("Order address geocoded")
}

func geoConfidence(locationType string, partialMatch bool) string {
	if partialMatch {
		return "low"
	}
	if locationType == "ROOFTOP" {
		return "high"
	}
	return "medium"
}

func buildGeocodeQuery(street, city, state string) string {
	parts := make([]string, 0, 5)

	if seg := streetSegment(street, 0); seg != "" {
		parts = append(parts, seg)
	}
	if seg := streetSegment(street, 2); seg != "" {
		parts = append(parts, seg)
	}
	if c := strings.TrimSpace(city); c != "" {
		parts = append(parts, c)
	}
	if s := strings.TrimSpace(state); s != "" {
		parts = append(parts, s)
	}

	if len(parts) == 0 {
		return ""
	}

	parts = append(parts, "Colombia")
	return strings.Join(parts, ", ")
}

func streetSegment(street string, index int) string {
	segments := strings.Split(street, "|")
	if index >= len(segments) {
		return ""
	}
	return strings.TrimSpace(segments[index])
}
