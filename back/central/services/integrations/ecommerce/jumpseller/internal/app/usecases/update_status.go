package usecases

import (
	"context"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
)

var probabilityToShipment = map[string]string{
	"picking":    domain.ShipmentRequested,
	"packing":    domain.ShipmentRequested,
	"processing": domain.ShipmentRequested,
	"paid":       domain.ShipmentRequested,
	"in_transit": domain.ShipmentInTransit,
	"shipped":    domain.ShipmentInTransit,
	"delivered":  domain.ShipmentDelivered,
	"fulfilled":  domain.ShipmentDelivered,
	"returned":   domain.ShipmentFailed,
	"failed":     domain.ShipmentFailed,
}

func convertProbabilityStatus(status string) (orderStatus string, shipmentStatus string, err error) {
	if status == "cancelled" || status == "canceled" {
		return domain.StatusCanceled, "", nil
	}
	if shipment, ok := probabilityToShipment[status]; ok {
		return "", shipment, nil
	}
	return "", "", domain.ErrStatusNotMapped
}

func (uc *jumpsellerUseCase) UpdateOrderStatus(ctx context.Context, integrationID string, externalOrderID string, probabilityStatus string, tracking domain.UpdateOrderFields) error {
	integration, cred, err := uc.resolveIntegration(ctx, integrationID)
	if err != nil {
		return err
	}

	if enabled, _ := integration.Config["status_sync_enabled"].(bool); !enabled {
		uc.logger.Info(ctx).
			Str("integration_id", integrationID).
			Msg("Sync de estados desactivado para la integracion Jumpseller, actualizacion omitida")
		return nil
	}

	orderID, parseErr := strconv.ParseInt(externalOrderID, 10, 64)
	if parseErr != nil {
		return domain.ErrNoOrdersFound
	}

	orderStatus, shipmentStatus, err := convertProbabilityStatus(probabilityStatus)
	if err != nil {
		uc.logger.Info(ctx).
			Str("integration_id", integrationID).
			Str("status", probabilityStatus).
			Msg("Estado sin homologacion a Jumpseller, actualizacion omitida")
		return nil
	}

	fields := domain.UpdateOrderFields{
		Status:          orderStatus,
		ShipmentStatus:  shipmentStatus,
		TrackingNumber:  tracking.TrackingNumber,
		TrackingCompany: tracking.TrackingCompany,
		AdditionalInfo:  tracking.AdditionalInfo,
	}

	if err := uc.client.UpdateOrder(ctx, cred, orderID, fields); err != nil {
		uc.logger.Error(ctx).Err(err).
			Str("integration_id", integrationID).
			Int64("order_id", orderID).
			Str("status", probabilityStatus).
			Msg("Error al actualizar el estado de la orden en Jumpseller")
		return err
	}

	uc.logger.Info(ctx).
		Str("integration_id", integrationID).
		Int64("order_id", orderID).
		Str("order_status", orderStatus).
		Str("shipment_status", shipmentStatus).
		Msg("Estado de la orden actualizado en Jumpseller")

	return nil
}
