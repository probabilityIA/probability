package app

import (
	"context"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
)

func (uc *UseCase) UpdateClient(ctx context.Context, dto dtos.UpdateClientDTO) (*entities.Client, error) {
	existing, err := uc.repo.GetByID(ctx, dto.BusinessID, dto.ID)
	if err != nil {
		return nil, err
	}

	if dto.Email != nil && *dto.Email != "" {
		existingEmail := ""
		if existing.Email != nil {
			existingEmail = *existing.Email
		}
		if *dto.Email != existingEmail {
			exists, err := uc.repo.ExistsByEmail(ctx, dto.BusinessID, *dto.Email, &dto.ID)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, domainerrors.ErrDuplicateEmail
			}
		}
	}

	if dto.Dni != nil && *dto.Dni != "" {
		existingDni := existing.Dni
		if existingDni == nil || *existingDni != *dto.Dni {
			exists, err := uc.repo.ExistsByDni(ctx, dto.BusinessID, *dto.Dni, &dto.ID)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, domainerrors.ErrDuplicateDni
			}
		}
	}

	email := dto.Email
	if email != nil && *email == "" {
		email = nil
	}

	existing.Name = dto.Name
	existing.Email = email
	existing.Phone = dto.Phone
	existing.Dni = dto.Dni
	existing.Address = dto.Address
	existing.City = dto.City
	existing.Notes = dto.Notes

	return uc.repo.Update(ctx, existing)
}
