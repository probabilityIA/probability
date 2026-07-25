package usecasecreateorder

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
)

func (uc *UseCaseCreateOrder) resolveCodIncludesShipping(ctx context.Context, dto *dtos.ProbabilityOrderDTO) {
	dto.CodIncludesShipping = true

	if dto.IntegrationID == 0 {
		return
	}

	includes, err := uc.repo.GetIntegrationCodIncludesShipping(ctx, dto.IntegrationID)
	if err != nil {
		uc.logger.Warn(ctx).Err(err).
			Uint("integration_id", dto.IntegrationID).
			Msg("No se pudo leer cod_includes_shipping de la integracion, se asume que el total ya incluye el envio")
		return
	}

	dto.CodIncludesShipping = includes
}
