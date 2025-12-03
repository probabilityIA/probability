package usecaseintegrationtype

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// UpdateIntegrationType actualiza un tipo de integración
func (uc *integrationTypeUseCase) UpdateIntegrationType(ctx context.Context, id uint, dto domain.UpdateIntegrationTypeDTO) (*domain.IntegrationType, error) {
	ctx = log.WithFunctionCtx(ctx, "UpdateIntegrationType")

	// Obtener el tipo de integración existente
	existing, err := uc.repo.GetIntegrationTypeByID(ctx, id)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("id", id).Msg("Error al obtener tipo de integración para actualizar")
		return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationTypeNotFound, err)
	}

	// Validar nombre si se está actualizando
	if dto.Name != nil && *dto.Name != existing.Name {
		nameExists, err := uc.repo.GetIntegrationTypeByName(ctx, *dto.Name)
		if err != nil && !strings.Contains(err.Error(), "no encontrado") {
			uc.log.Error(ctx).Err(err).Str("name", *dto.Name).Msg("Error al verificar si el nombre ya existe")
			return nil, fmt.Errorf("error al verificar disponibilidad del nombre: %w", err)
		}
		if nameExists != nil {
			uc.log.Warn(ctx).Str("name", *dto.Name).Msg("El nombre del tipo de integración ya está en uso")
			return nil, fmt.Errorf("%w: %s", domain.ErrIntegrationTypeNameExists, *dto.Name)
		}
		existing.Name = *dto.Name
	}

	// Validar código si se está actualizando
	if dto.Code != nil && *dto.Code != existing.Code {
		codeExists, err := uc.repo.GetIntegrationTypeByCode(ctx, *dto.Code)
		if err != nil && !strings.Contains(err.Error(), "no encontrado") {
			uc.log.Error(ctx).Err(err).Str("code", *dto.Code).Msg("Error al verificar si el código ya existe")
			return nil, fmt.Errorf("error al verificar disponibilidad del código: %w", err)
		}
		if codeExists != nil {
			uc.log.Warn(ctx).Str("code", *dto.Code).Msg("El código del tipo de integración ya está en uso")
			return nil, fmt.Errorf("%w: %s", domain.ErrIntegrationTypeCodeExists, *dto.Code)
		}
		existing.Code = *dto.Code
	}

	// Actualizar campos opcionales
	if dto.Description != nil {
		existing.Description = *dto.Description
	}
	if dto.Icon != nil {
		existing.Icon = *dto.Icon
	}
	if dto.Category != nil {
		existing.Category = *dto.Category
	}
	if dto.IsActive != nil {
		existing.IsActive = *dto.IsActive
	}
	if dto.ConfigSchema != nil {
		existing.ConfigSchema = *dto.ConfigSchema
	}
	if dto.CredentialsSchema != nil {
		existing.CredentialsSchema = *dto.CredentialsSchema
	}

	if err := uc.repo.UpdateIntegrationType(ctx, id, existing); err != nil {
		uc.log.Error(ctx).Err(err).
			Uint("id", id).
			Msg("Error al guardar los cambios del tipo de integración en la base de datos")
		return nil, fmt.Errorf("error al guardar los cambios del tipo de integración: %w", err)
	}

	uc.log.Info(ctx).
		Uint("id", id).
		Msg("Tipo de integración actualizado exitosamente")

	return existing, nil
}
