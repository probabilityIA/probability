package app

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/dtos"
)

const maxCarrierFee = 100000000.0

func (uc *UseCase) UpdateCarrierFee(ctx context.Context, d dtos.UpdateCarrierFeeDTO) (*dtos.UpdateCarrierFeeResult, error) {
	if d.BusinessID == 0 {
		return nil, errors.New("el negocio es requerido")
	}
	if d.ShipmentID == 0 {
		return nil, errors.New("el envio es requerido")
	}
	if d.Fee < 0 {
		return nil, errors.New("la comision no puede ser negativa")
	}
	if d.Fee > maxCarrierFee {
		return nil, errors.New("la comision excede el maximo permitido")
	}

	previous, orderNumber, err := uc.repo.ShipmentCarrierFee(ctx, d.BusinessID, d.ShipmentID)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.UpdateShipmentCarrierFee(ctx, d.BusinessID, d.ShipmentID, d.Fee); err != nil {
		return nil, err
	}

	uc.log.Info(ctx).
		Uint("business_id", d.BusinessID).
		Uint("shipment_id", d.ShipmentID).
		Str("order_number", orderNumber).
		Float64("comision_anterior", previous).
		Float64("comision_nueva", d.Fee).
		Uint("user_id", d.UserID).
		Str("user_name", d.UserName).
		Msg("Comision de transportadora corregida manualmente")

	return &dtos.UpdateCarrierFeeResult{
		OrderNumber: orderNumber,
		PreviousFee: previous,
		NewFee:      d.Fee,
	}, nil
}
