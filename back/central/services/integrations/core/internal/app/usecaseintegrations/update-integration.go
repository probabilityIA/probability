package usecaseintegrations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"gorm.io/datatypes"
)

// UpdateIntegration actualiza una integración existente
func (uc *IntegrationUseCase) UpdateIntegration(ctx context.Context, id uint, dto domain.UpdateIntegrationDTO) (*domain.Integration, error) {
	ctx = log.WithFunctionCtx(ctx, "UpdateIntegration")

	// Obtener integración existente
	existing, err := uc.repo.GetIntegrationByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationNotFound, err)
	}

	// Actualizar campos si se proporcionan
	if dto.Name != nil {
		existing.Name = *dto.Name
	}
	if dto.Code != nil {
		// Validar que el nuevo código no exista
		exists, err := uc.repo.ExistsIntegrationByCode(ctx, *dto.Code, existing.BusinessID)
		if err != nil {
			return nil, fmt.Errorf("error al verificar código: %w", err)
		}
		if exists && *dto.Code != existing.Code {
			return nil, fmt.Errorf("%w: %s", domain.ErrIntegrationCodeExists, *dto.Code)
		}
		existing.Code = *dto.Code
	}
	if dto.IsActive != nil {
		existing.IsActive = *dto.IsActive
	}
	if dto.IsDefault != nil {
		existing.IsDefault = *dto.IsDefault
		if *dto.IsDefault {
			// Si se marca como default, desmarcar las demás
			if err := uc.repo.SetIntegrationAsDefault(ctx, id); err != nil {
				return nil, fmt.Errorf("error al marcar como default: %w", err)
			}
		}
	}
	if dto.Config != nil {
		existing.Config = *dto.Config
	}
	if dto.Credentials != nil {
		// Serializar credenciales (se encriptarán en el repository)
		credentialsBytes, err := json.Marshal(*dto.Credentials)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationCredentialsSerialize, err)
		}
		existing.Credentials = datatypes.JSON(credentialsBytes)
	}
	if dto.Description != nil {
		existing.Description = *dto.Description
	}

	// Solo actualizar UpdatedByID si es un ID válido (mayor que 0)
	if dto.UpdatedByID > 0 {
		existing.UpdatedByID = &dto.UpdatedByID
	}
	// Si UpdatedByID es 0, no actualizamos el campo (mantiene el valor existente o NULL)

	// Guardar cambios
	if err := uc.repo.UpdateIntegration(ctx, id, existing); err != nil {
		uc.log.Error(ctx).Err(err).Uint("id", id).Msg("Error al actualizar integración")
		return nil, fmt.Errorf("error al actualizar integración: %w", err)
	}

	uc.log.Info(ctx).Uint("id", id).Msg("Integración actualizada exitosamente")

	return existing, nil
}
