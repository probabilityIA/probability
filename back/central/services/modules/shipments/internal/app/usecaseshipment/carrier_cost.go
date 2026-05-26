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

	codCustomerCharge := 0.0
	if shipment.CodCustomerCharge != nil {
		codCustomerCharge = *shipment.CodCustomerCharge
	}

	codAppliedMargin := 0.0
	codCarrierCost := codCustomerCharge
	if codCustomerCharge > 0 && margin.CODMarginPercent > 0 {
		codBase := codCustomerCharge / (1 + margin.CODMarginPercent/100.0)
		codAppliedMargin = codCustomerCharge - codBase
		codCarrierCost = codBase
	}

	fleteCustomer := *shipment.TotalCost - codCustomerCharge
	if fleteCustomer < 0 {
		fleteCustomer = 0
	}

	fleteMargin := margin.MarginAmount + margin.InsuranceMargin
	fleteCarrierCost := fleteCustomer - fleteMargin
	if fleteCarrierCost < 0 {
		fleteCarrierCost = 0
	}

	totalCarrierCost := fleteCarrierCost + codCarrierCost
	shipment.CarrierCost = &totalCarrierCost
	shipment.AppliedMargin = &fleteMargin
	if codCustomerCharge > 0 {
		applied := codAppliedMargin
		shipment.CodAppliedMargin = &applied
	}
}
