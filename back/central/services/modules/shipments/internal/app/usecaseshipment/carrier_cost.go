package usecaseshipment

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

func (uc *UseCaseShipment) applyCarrierCost(ctx context.Context, shipment *domain.Shipment) {
	if shipment.TotalCost == nil || *shipment.TotalCost <= 0 {
		return
	}
	if shipment.CarrierCost != nil {
		return
	}
	if uc.marginReader == nil {
		return
	}
	if shipment.OrderID == nil || *shipment.OrderID == "" {
		return
	}

	carrierCode := ""
	if shipment.CarrierCode != nil {
		carrierCode = strings.TrimSpace(*shipment.CarrierCode)
	}
	if carrierCode == "" && shipment.Carrier != nil {
		carrierCode = strings.TrimSpace(*shipment.Carrier)
	}
	if carrierCode == "" {
		return
	}
	carrierCode = strings.ToLower(carrierCode)

	businessID, err := uc.repo.GetOrderBusinessID(ctx, *shipment.OrderID)
	if err != nil || businessID == 0 {
		return
	}

	margin, err := uc.marginReader.Get(ctx, businessID, carrierCode)
	if err != nil {
		return
	}

	totalMargin := margin.MarginAmount + margin.InsuranceMargin
	if totalMargin <= 0 {
		zero := 0.0
		applied := zero
		shipment.AppliedMargin = &applied
		carrierCost := *shipment.TotalCost
		shipment.CarrierCost = &carrierCost
		return
	}

	carrierCost := *shipment.TotalCost - totalMargin
	if carrierCost < 0 {
		carrierCost = 0
	}
	shipment.CarrierCost = &carrierCost
	shipment.AppliedMargin = &totalMargin
}
