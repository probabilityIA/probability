package usecasebusinesstype

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/domain"
)

// CreateBusinessType crea un nuevo tipo de negocio
func (uc *BusinessTypeUseCase) CreateBusinessType(ctx context.Context, request domain.BusinessTypeRequest) (*domain.BusinessTypeResponse, error) {
	uc.log.Info().Str("name", request.Name).Str("code", request.Code).Msg("Creando tipo de negocio")

	// Validar que el código no exista
	existing, err := uc.repository.GetBusinessTypeByCode(ctx, request.Code)
	if err != nil && err.Error() != "tipo de negocio no encontrado" {
		uc.log.Error().Err(err).Str("code", request.Code).Msg("Error al verificar código existente")
		return nil, fmt.Errorf("error al verificar código existente: %w", err)
	}

	if existing != nil {
		uc.log.Warn().Str("code", request.Code).Msg("Código de tipo de negocio ya existe")
		return nil, fmt.Errorf("el código '%s' ya existe", request.Code)
	}

	// Crear entidad
	businessType := domain.BusinessType{
		Name:        request.Name,
		Code:        request.Code,
		Description: request.Description,
		Icon:        request.Icon,
		IsActive:    request.IsActive,
	}

	// Guardar en repositorio
	_, err = uc.repository.CreateBusinessType(ctx, businessType)
	if err != nil {
		uc.log.Error().Err(err).Str("name", request.Name).Msg("Error al crear tipo de negocio")
		return nil, fmt.Errorf("error al crear tipo de negocio: %w", err)
	}

	// Obtener el tipo de negocio creado
	created, err := uc.repository.GetBusinessTypeByCode(ctx, request.Code)
	if err != nil {
		uc.log.Error().Err(err).Str("code", request.Code).Msg("Error al obtener tipo de negocio creado")
		return nil, fmt.Errorf("error al obtener tipo de negocio creado: %w", err)
	}

	response := &domain.BusinessTypeResponse{
		ID:          created.ID,
		Name:        created.Name,
		Code:        created.Code,
		Description: created.Description,
		Icon:        created.Icon,
		IsActive:    created.IsActive,
		CreatedAt:   created.CreatedAt,
		UpdatedAt:   created.UpdatedAt,
	}

	uc.log.Info().Uint("id", created.ID).Str("name", request.Name).Msg("Tipo de negocio creado exitosamente")
	return response, nil
}
