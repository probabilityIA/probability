package app

import (
	"context"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
)

func (uc *UseCase) CreateClient(ctx context.Context, dto dtos.CreateClientDTO) (*entities.Client, error) {
	// Verificar email duplicado
	if dto.Email != nil && *dto.Email != "" {
		exists, err := uc.repo.ExistsByEmail(ctx, dto.BusinessID, *dto.Email, nil)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, domainerrors.ErrDuplicateEmail
		}
	}

	// Verificar DNI duplicado
	if dto.Dni != nil && *dto.Dni != "" {
		exists, err := uc.repo.ExistsByDni(ctx, dto.BusinessID, *dto.Dni, nil)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, domainerrors.ErrDuplicateDni
		}
	}

	// Normalizar: si email es puntero a string vacío, guardar como nil
	email := dto.Email
	if email != nil && *email == "" {
		email = nil
	}

	client := &entities.Client{
		BusinessID: dto.BusinessID,
		Name:       dto.Name,
		Email:      email,
		Phone:      dto.Phone,
		Dni:        dto.Dni,
	}

	return uc.repo.Create(ctx, client)
}
