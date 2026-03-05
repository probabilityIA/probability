package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/entities"
)

func (uc *UseCase) UpdateStop(ctx context.Context, dto dtos.UpdateStopDTO) (*entities.RouteStop, error) {
	_, err := uc.repo.GetRouteByID(ctx, dto.BusinessID, dto.RouteID)
	if err != nil {
		return nil, err
	}

	existing, err := uc.repo.GetStopByID(ctx, dto.RouteID, dto.ID)
	if err != nil {
		return nil, err
	}

	existing.Address = dto.Address
	existing.City = dto.City
	existing.Lat = dto.Lat
	existing.Lng = dto.Lng
	existing.CustomerName = dto.CustomerName
	existing.CustomerPhone = dto.CustomerPhone
	existing.DeliveryNotes = dto.DeliveryNotes

	return uc.repo.UpdateStop(ctx, existing)
}
