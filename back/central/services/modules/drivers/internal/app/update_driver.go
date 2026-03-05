package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/errors"
)

func (uc *UseCase) UpdateDriver(ctx context.Context, dto dtos.UpdateDriverDTO) (*entities.Driver, error) {
	existing, err := uc.repo.GetByID(ctx, dto.BusinessID, dto.ID)
	if err != nil {
		return nil, err
	}

	if dto.Identification != "" && dto.Identification != existing.Identification {
		exists, err := uc.repo.ExistsByIdentification(ctx, dto.BusinessID, dto.Identification, &dto.ID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, domainerrors.ErrDuplicateIdentification
		}
	}

	existing.FirstName = dto.FirstName
	existing.LastName = dto.LastName
	existing.Email = dto.Email
	existing.Phone = dto.Phone
	existing.Identification = dto.Identification
	existing.Status = dto.Status
	existing.PhotoURL = dto.PhotoURL
	existing.LicenseType = dto.LicenseType
	existing.LicenseExpiry = dto.LicenseExpiry
	existing.WarehouseID = dto.WarehouseID
	existing.Notes = dto.Notes

	return uc.repo.Update(ctx, existing)
}
