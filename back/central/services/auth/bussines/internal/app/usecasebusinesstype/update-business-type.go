package usecasebusinesstype

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/domain"
)

// UpdateBusinessType actualiza un tipo de negocio existente
func (uc *BusinessTypeUseCase) UpdateBusinessType(ctx context.Context, id uint, request domain.BusinessTypeRequest) (*domain.BusinessTypeResponse, error) {
	uc.log.Info().Uint("id", id).Str("name", request.Name).Msg("Actualizando tipo de negocio")

	// Verificar que existe
	existing, err := uc.repository.GetBusinessTypeByID(ctx, id)
	if err != nil {
		uc.log.Error().Err(err).Uint("id", id).Msg("Error al obtener tipo de negocio para actualizar")
		return nil, fmt.Errorf("error al obtener tipo de negocio: %w", err)
	}

	if existing == nil {
		uc.log.Warn().Uint("id", id).Msg("Tipo de negocio no encontrado para actualizar")
		return nil, fmt.Errorf("tipo de negocio no encontrado")
	}

	// Verificar que el código no exista en otro tipo de negocio
	if request.Code != existing.Code {
		codeExists, err := uc.repository.GetBusinessTypeByCode(ctx, request.Code)
		if err != nil && err.Error() != "tipo de negocio no encontrado" {
			uc.log.Error().Err(err).Str("code", request.Code).Msg("Error al verificar código existente")
			return nil, fmt.Errorf("error al verificar código existente: %w", err)
		}

		if codeExists != nil {
			uc.log.Warn().Str("code", request.Code).Msg("Código de tipo de negocio ya existe")
			return nil, fmt.Errorf("el código '%s' ya existe", request.Code)
		}
	}

	// Actualizar entidad
	businessType := domain.BusinessType{
		Name:        request.Name,
		Code:        request.Code,
		Description: request.Description,
		Icon:        request.Icon,
		IsActive:    request.IsActive,
	}

	// Guardar en repositorio
	_, err = uc.repository.UpdateBusinessType(ctx, id, businessType)
	if err != nil {
		uc.log.Error().Err(err).Uint("id", id).Msg("Error al actualizar tipo de negocio")
		return nil, fmt.Errorf("error al actualizar tipo de negocio: %w", err)
	}

	// Obtener el tipo de negocio actualizado
	updated, err := uc.repository.GetBusinessTypeByID(ctx, id)
	if err != nil {
		uc.log.Error().Err(err).Uint("id", id).Msg("Error al obtener tipo de negocio actualizado")
		return nil, fmt.Errorf("error al obtener tipo de negocio actualizado: %w", err)
	}

	response := &domain.BusinessTypeResponse{
		ID:          updated.ID,
		Name:        updated.Name,
		Code:        updated.Code,
		Description: updated.Description,
		Icon:        updated.Icon,
		IsActive:    updated.IsActive,
		CreatedAt:   updated.CreatedAt,
		UpdatedAt:   updated.UpdatedAt,
	}

	uc.log.Info().Uint("id", id).Str("name", request.Name).Msg("Tipo de negocio actualizado exitosamente")
	return response, nil
}
