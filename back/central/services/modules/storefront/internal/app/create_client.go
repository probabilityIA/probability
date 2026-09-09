package app

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/errors"
)

const clienteFinalMaxRoleLevel = 2

func (uc *UseCase) CreateClient(ctx context.Context, businessID, requesterUserID uint, dto *dtos.CreateClientDTO) (*entities.StorefrontClient, string, error) {
	active, err := uc.repo.IsIntegrationActiveOrMissing(ctx, businessID, tiendaIntegrationTypeID)
	if err != nil {
		return nil, "", err
	}
	if !active {
		return nil, "", domainerrors.ErrStorefrontNotActive
	}

	level, err := uc.repo.GetRoleLevelByUserAndBusiness(ctx, requesterUserID, businessID)
	if err != nil {
		return nil, "", err
	}
	if level > clienteFinalMaxRoleLevel {
		return nil, "", domainerrors.ErrRoleNotAllowed
	}

	exists, err := uc.repo.UserExistsByEmail(ctx, dto.Email)
	if err != nil {
		return nil, "", err
	}
	if exists {
		return nil, "", domainerrors.ErrEmailAlreadyExists
	}

	password := dto.Password
	generatedPassword := ""
	if password == "" {
		password, err = generateTempPassword()
		if err != nil {
			return nil, "", err
		}
		generatedPassword = password
	}

	userID, err := uc.repo.CreateUser(ctx, &entities.NewUser{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: password,
		Phone:    dto.Phone,
	})
	if err != nil {
		return nil, "", fmt.Errorf("error creando usuario: %w", err)
	}

	roleID, err := uc.repo.GetClienteFinalRoleID(ctx)
	if err != nil {
		return nil, "", domainerrors.ErrRoleNotFound
	}

	if err := uc.repo.CreateBusinessStaff(ctx, userID, businessID, roleID); err != nil {
		return nil, "", fmt.Errorf("error creando business staff: %w", err)
	}

	client := &entities.StorefrontClient{
		BusinessID: businessID,
		UserID:     &userID,
		Name:       dto.Name,
		Email:      &dto.Email,
		Phone:      dto.Phone,
		Dni:        dto.Dni,
	}
	if err := uc.repo.CreateClient(ctx, client); err != nil {
		return nil, "", fmt.Errorf("error creando cliente: %w", err)
	}

	uc.logger.Info(ctx).
		Uint("business_id", businessID).
		Uint("requester_user_id", requesterUserID).
		Uint("new_user_id", userID).
		Str("email", dto.Email).
		Msg("Cliente del catalogo creado por el negocio")

	return client, generatedPassword, nil
}

func generateTempPassword() (string, error) {
	buf := make([]byte, 9)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("error generando contrasena temporal: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
