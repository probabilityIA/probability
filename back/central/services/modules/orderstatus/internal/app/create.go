package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
	domainErrors "github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/errors"
)

func (uc *useCase) CreateOrderStatusMapping(ctx context.Context, mapping *entities.OrderStatusMapping) (*entities.OrderStatusMapping, error) {
	// Verificar si ya existe
	exists, err := uc.repo.Exists(ctx, mapping.IntegrationTypeID, mapping.OriginalStatus)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domainErrors.ErrMappingAlreadyExists
	}

	// Asegurar que esté activo por defecto
	mapping.IsActive = true

	if err := uc.repo.Create(ctx, mapping); err != nil {
		return nil, err
	}

	// Recargar para obtener relaciones
	return uc.repo.GetByID(ctx, mapping.ID)
}
