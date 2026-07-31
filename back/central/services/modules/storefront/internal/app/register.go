package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/errors"
)

func (uc *UseCase) Register(ctx context.Context, dto *dtos.RegisterDTO) error {
	business, err := uc.repo.GetBusinessByCode(ctx, dto.BusinessCode)
	if err != nil {
		return domainerrors.ErrBusinessNotFound
	}

	existingUserID, passwordHash, err := uc.repo.GetUserAuthByEmail(ctx, dto.Email)
	if err != nil {
		return fmt.Errorf("error verificando email: %w", err)
	}

	var userID uint
	if existingUserID > 0 {
		if !uc.repo.VerifyPassword(passwordHash, dto.Password) {
			return domainerrors.ErrPasswordMismatch
		}
		userID = existingUserID
	} else {
		newUser := &entities.NewUser{
			Name:     dto.Name,
			Email:    dto.Email,
			Password: dto.Password,
			Phone:    dto.Phone,
		}
		userID, err = uc.repo.CreateUser(ctx, newUser)
		if err != nil {
			return fmt.Errorf("error creando usuario: %w", err)
		}
	}

	roleID, err := uc.repo.GetClienteFinalRoleID(ctx)
	if err != nil {
		return domainerrors.ErrRoleNotFound
	}

	hasStaff, err := uc.repo.StaffExists(ctx, userID, business.ID)
	if err != nil {
		return fmt.Errorf("error verificando staff: %w", err)
	}
	if !hasStaff {
		if err := uc.repo.CreateBusinessStaff(ctx, userID, business.ID, roleID); err != nil {
			return fmt.Errorf("error creando business staff: %w", err)
		}
	}

	existingClient, err := uc.repo.GetClientByUserID(ctx, business.ID, userID)
	if err == nil && existingClient != nil {
		return nil
	}

	byEmail, err := uc.repo.GetClientByBusinessAndEmail(ctx, business.ID, dto.Email)
	if err != nil {
		return fmt.Errorf("error buscando cliente: %w", err)
	}
	if byEmail != nil {
		if byEmail.UserID != nil && *byEmail.UserID != userID {
			return domainerrors.ErrEmailAlreadyExists
		}
		if byEmail.UserID == nil {
			if err := uc.repo.LinkClientUser(ctx, byEmail.ID, userID); err != nil {
				return fmt.Errorf("error vinculando cliente: %w", err)
			}
		}
		uc.logger.Info(ctx).
			Uint("user_id", userID).
			Uint("business_id", business.ID).
			Str("email", dto.Email).
			Msg("Cliente existente vinculado a usuario en storefront")
		return nil
	}

	client := &entities.StorefrontClient{
		BusinessID: business.ID,
		UserID:     &userID,
		Name:       dto.Name,
		Email:      &dto.Email,
		Phone:      dto.Phone,
		Dni:        dto.Dni,
	}
	if err := uc.repo.CreateClient(ctx, client); err != nil {
		return fmt.Errorf("error creando cliente: %w", err)
	}

	uc.logger.Info(ctx).
		Uint("user_id", userID).
		Uint("business_id", business.ID).
		Str("email", dto.Email).
		Msg("Nuevo cliente registrado en storefront")

	return nil
}
