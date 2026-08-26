package usecases

import (
	"context"
	"fmt"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

func (uc *meliUseCase) GetShipmentLabel(ctx context.Context, shipmentID uint, businessID uint, responseType string) (*domain.ShipmentLabel, error) {
	if uc.orderLookupRepo == nil {
		return nil, domain.ErrShipmentNotFound
	}

	ref, err := uc.orderLookupRepo.GetMeliLabelRefByShipmentID(ctx, shipmentID)
	if err != nil {
		return nil, fmt.Errorf("looking up shipment %d: %w", shipmentID, err)
	}
	if ref == nil {
		return nil, domain.ErrShipmentNotFound
	}
	if ref.BusinessID != businessID {
		return nil, domain.ErrShipmentNotOwned
	}
	if ref.IntegrationID == 0 || ref.MeliShipmentID == 0 {
		return nil, domain.ErrLabelNotAvailable
	}

	integrationID := strconv.FormatUint(uint64(ref.IntegrationID), 10)

	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil {
		return nil, domain.ErrIntegrationNotFound
	}
	if integration.BusinessID == nil || *integration.BusinessID != businessID {
		return nil, domain.ErrShipmentNotOwned
	}

	token, err := uc.EnsureValidToken(ctx, integrationID)
	if err != nil {
		return nil, err
	}

	cli := uc.clientFor(ctx, integration)

	label, err := cli.GetShipmentLabel(ctx, token, ref.MeliShipmentID, responseType)
	if err == domain.ErrTokenExpired {
		newToken, refreshErr := uc.refreshAccessToken(ctx, integrationID, integration)
		if refreshErr != nil {
			uc.logger.Error(ctx).Err(refreshErr).
				Uint("shipment_id", shipmentID).
				Msg("No se pudo refrescar el token para descargar la etiqueta")
			return nil, domain.ErrTokenRefreshFailed
		}
		label, err = cli.GetShipmentLabel(ctx, newToken, ref.MeliShipmentID, responseType)
	}

	if err != nil {
		event := uc.logger.Error(ctx)
		if err == domain.ErrLabelAlreadyShipped || err == domain.ErrLabelNotAvailable {
			event = uc.logger.Warn(ctx)
		}
		event.Err(err).
			Uint("shipment_id", shipmentID).
			Int64("meli_shipment_id", ref.MeliShipmentID).
			Msg("No se pudo descargar la etiqueta del canal")
		return nil, err
	}

	return label, nil
}
