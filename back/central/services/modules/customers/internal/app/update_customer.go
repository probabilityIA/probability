package app

import (
	"context"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
)

func (uc *UseCase) UpdateClient(ctx context.Context, dto dtos.UpdateClientDTO) (*entities.Client, error) {
	// Verificar que existe
	existing, err := uc.repo.GetByID(ctx, dto.BusinessID, dto.ID)
	if err != nil {
		return nil, err
	}

	// Verificar email duplicado (excluyendo el propio cliente)
	if dto.Email != "" && dto.Email != existing.Email {
		exists, err := uc.repo.ExistsByEmail(ctx, dto.BusinessID, dto.Email, &dto.ID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, domainerrors.ErrDuplicateEmail
		}
	}

	// Verificar DNI duplicado (excluyendo el propio cliente)
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

	existing.Name = dto.Name
	existing.Email = dto.Email
	existing.Phone = dto.Phone
	existing.Dni = dto.Dni

	return uc.repo.Update(ctx, existing)
}
