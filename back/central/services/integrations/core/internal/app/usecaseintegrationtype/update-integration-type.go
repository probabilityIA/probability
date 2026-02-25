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
	if dto.CategoryID != nil {
		existing.CategoryID = *dto.CategoryID
	}
	if dto.IsActive != nil {
		existing.IsActive = *dto.IsActive
	}
	if dto.InDevelopment != nil {
		existing.InDevelopment = *dto.InDevelopment
	}
	if dto.ConfigSchema != nil {
		existing.ConfigSchema = *dto.ConfigSchema
	}
	if dto.CredentialsSchema != nil {
		existing.CredentialsSchema = *dto.CredentialsSchema
	}
	if dto.BaseURL != nil {
		existing.BaseURL = *dto.BaseURL
	}
	if dto.BaseURLTest != nil {
		existing.BaseURLTest = *dto.BaseURLTest
	}

	// Procesar imagen si se proporciona una nueva
	if dto.ImageFile != nil {
		uc.log.Info(ctx).Uint("id", id).Msg("Subiendo nueva imagen del tipo de integración a S3")

		// Subir nueva imagen a S3 en la carpeta "integration-types"
		// Retorna el path relativo (ej: "integration-types/1234567890_logo.jpg")
		imagePath, err := uc.s3.UploadImage(ctx, dto.ImageFile, "integration-types")
		if err != nil {
			uc.log.Error(ctx).Err(err).Uint("id", id).Msg("Error al subir nueva imagen del tipo de integración")
			return nil, fmt.Errorf("%w: %v", domain.ErrIntegrationTypeImageUploadFailed, err)
		}

		// Guardar solo el path relativo en la base de datos
		existing.ImageURL = imagePath
		uc.log.Info(ctx).Uint("id", id).Str("image_path", imagePath).Msg("Nueva imagen del tipo de integración subida exitosamente")

		// Eliminar imagen anterior si existe y es diferente
		if existing.ImageURL != "" && existing.ImageURL != imagePath {
			// Verificar si la imagen anterior es un path relativo (no URL completa)
			if !strings.HasPrefix(existing.ImageURL, "http") {
				uc.log.Info(ctx).Uint("id", id).Str("old_image", existing.ImageURL).Msg("Eliminando imagen anterior del tipo de integración")
				if err := uc.s3.DeleteImage(ctx, existing.ImageURL); err != nil {
					uc.log.Warn(ctx).Err(err).Str("old_image", existing.ImageURL).Msg("Error al eliminar imagen anterior (no crítico)")
					// No fallar la actualización si no se puede eliminar la imagen anterior
				}
			}
		}
	} else if dto.RemoveImage {
		// Eliminar imagen solo si el cliente lo solicita explícitamente
		uc.log.Info(ctx).Uint("id", id).Str("old_image", existing.ImageURL).Msg("Eliminando imagen del tipo de integración")

		// Verificar si la imagen anterior es un path relativo (no URL completa)
		if !strings.HasPrefix(existing.ImageURL, "http") {
			if err := uc.s3.DeleteImage(ctx, existing.ImageURL); err != nil {
				uc.log.Warn(ctx).Err(err).Str("old_image", existing.ImageURL).Msg("Error al eliminar imagen anterior (no crítico)")
				// No fallar la actualización si no se puede eliminar la imagen
			}
		}
		existing.ImageURL = "" // Limpiar la URL
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

	// Invalidar caché de todas las integraciones que usan este tipo
	// para que en la próxima consulta se recarguen con las URLs actualizadas
	if integrations, err := uc.repo.ListIntegrationsByIntegrationTypeID(ctx, id); err == nil {
		for _, integration := range integrations {
			if cacheErr := uc.cache.InvalidateIntegration(ctx, integration.ID); cacheErr != nil {
				uc.log.Warn(ctx).
					Err(cacheErr).
					Uint("integration_id", integration.ID).
					Msg("Error al invalidar caché de integración tras actualizar tipo")
			}
		}
		uc.log.Info(ctx).
			Uint("type_id", id).
			Int("invalidated_count", len(integrations)).
			Msg("Caché invalidado para integraciones del tipo actualizado")
	} else {
		uc.log.Warn(ctx).Err(err).Uint("type_id", id).Msg("No se pudo obtener integraciones para invalidar caché")
	}

	return existing, nil
}
