package app

import (
	"context"
	"fmt"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/entities"
)

func (uc *UseCase) Register(ctx context.Context, dto *dtos.RegisterDTO) error {
	// Find business by code
	business, err := uc.repo.GetBusinessByCode(ctx, dto.BusinessCode)
	if err != nil {
		return domainerrors.ErrBusinessNotFound
	}

	// Check email uniqueness
	exists, err := uc.repo.UserExistsByEmail(ctx, dto.Email)
	if err != nil {
		return fmt.Errorf("error verificando email: %w", err)
	}
	if exists {
		return domainerrors.ErrEmailAlreadyExists
	}

	// Get cliente_final role ID
	roleID, err := uc.repo.GetClienteFinalRoleID(ctx)
	if err != nil {
		return domainerrors.ErrRoleNotFound
	}

	// Create user (password will be hashed in repository)
	newUser := &entities.NewUser{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
		Phone:    dto.Phone,
	}
	userID, err := uc.repo.CreateUser(ctx, newUser)
	if err != nil {
		return fmt.Errorf("error creando usuario: %w", err)
	}

	// Create business staff record
	if err := uc.repo.CreateBusinessStaff(ctx, userID, business.ID, roleID); err != nil {
		return fmt.Errorf("error creando business staff: %w", err)
	}

	// Create client record linked to user
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
