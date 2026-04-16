package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/errors"
)

func (uc *UseCase) CreateDriver(ctx context.Context, dto dtos.CreateDriverDTO) (*entities.Driver, error) {
	if dto.Identification != "" {
		exists, err := uc.repo.ExistsByIdentification(ctx, dto.BusinessID, dto.Identification, nil)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, domainerrors.ErrDuplicateIdentification
		}
	}

	status := dto.Status
	if status == "" {
		status = "active"
	}

	driver := &entities.Driver{
		BusinessID:     dto.BusinessID,
		FirstName:      dto.FirstName,
		LastName:       dto.LastName,
		Email:          dto.Email,
		Phone:          dto.Phone,
		Identification: dto.Identification,
		Status:         status,
		PhotoURL:       dto.PhotoURL,
		LicenseType:    dto.LicenseType,
		LicenseExpiry:  dto.LicenseExpiry,
		WarehouseID:    dto.WarehouseID,
		Notes:          dto.Notes,
	}

	return uc.repo.Create(ctx, driver)
}
