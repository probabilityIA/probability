package usecasecreateorder

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

func (uc *UseCaseCreateOrder) assignWarehouseIfMissing(ctx context.Context, order *entities.ProbabilityOrder) {
	if order.WarehouseID != nil && *order.WarehouseID > 0 {
		return
	}
	if order.BusinessID == nil || *order.BusinessID == 0 {
		return
	}

	ref, err := uc.repo.ResolveOrderWarehouse(ctx, *order.BusinessID, order.IntegrationID)
	if err != nil {
		uc.logger.Warn(ctx).Err(err).Uint("business_id", *order.BusinessID).Msg("No se pudo resolver la bodega de la orden")
		return
	}
	if ref == nil {
		return
	}

	id := ref.ID
	order.WarehouseID = &id
	if order.WarehouseName == "" {
		order.WarehouseName = ref.Name
	}
}
