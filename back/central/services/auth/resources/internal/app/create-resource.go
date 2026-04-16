package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/auth/resources/internal/domain"
)

// CreateResource crea un nuevo recurso
func (uc *ResourceUseCase) CreateResource(ctx context.Context, createDTO domain.CreateResourceDTO) (*domain.ResourceDTO, error) {
	uc.logger.Info().Str("name", createDTO.Name).Msg("Iniciando creación de recurso")

	// Validar datos de entrada
	if err := uc.validateCreateResource(createDTO); err != nil {
		uc.logger.Error().Err(err).Str("name", createDTO.Name).Msg("Validación fallida para crear recurso")
		return nil, err
	}

	// Verificar que no existe un recurso con el mismo nombre
	existingResource, err := uc.repository.GetResourceByName(ctx, createDTO.Name)
	if err == nil && existingResource != nil {
		uc.logger.Warn().Str("name", createDTO.Name).Msg("Recurso ya existe con ese nombre")
		return nil, fmt.Errorf("ya existe un recurso con el nombre '%s'", createDTO.Name)
	}

	// Crear entidad de dominio
	resource := domain.Resource{
		Name:        strings.TrimSpace(createDTO.Name),
		Description: strings.TrimSpace(createDTO.Description),
	}

	// Agregar business_type_id si está presente
	if createDTO.BusinessTypeID != nil {
		resource.BusinessTypeID = *createDTO.BusinessTypeID
	}

	// Crear recurso en el repositorio
	resourceID, err := uc.repository.CreateResource(ctx, resource)
	if err != nil {
		uc.logger.Error().Err(err).Str("name", createDTO.Name).Msg("Error al crear recurso")
		return nil, fmt.Errorf("error al crear recurso: %w", err)
	}

	// Obtener el recurso creado para devolver el DTO completo
	createdResource, err := uc.repository.GetResourceByID(ctx, resourceID)
	if err != nil {
		uc.logger.Error().Err(err).Uint("resource_id", resourceID).Msg("Error al obtener recurso creado")
		return nil, fmt.Errorf("error al obtener recurso creado: %w", err)
	}

	// Convertir a DTO
	resourceDTO := &domain.ResourceDTO{
		ID:               createdResource.ID,
		Name:             createdResource.Name,
		Description:      createdResource.Description,
		BusinessTypeID:   createdResource.BusinessTypeID,
		BusinessTypeName: createdResource.BusinessTypeName,
		CreatedAt:        createdResource.CreatedAt,
		UpdatedAt:        createdResource.UpdatedAt,
	}

	uc.logger.Info().
		Uint("resource_id", resourceID).
		Str("name", createDTO.Name).
		Msg("Recurso creado exitosamente")

	return resourceDTO, nil
}

// validateCreateResource valida los datos para crear un recurso
func (uc *ResourceUseCase) validateCreateResource(createDTO domain.CreateResourceDTO) error {
	if strings.TrimSpace(createDTO.Name) == "" {
		return fmt.Errorf("el nombre del recurso es obligatorio")
	}

	if len(createDTO.Name) > 100 {
		return fmt.Errorf("el nombre del recurso no puede exceder 100 caracteres")
	}

	if len(createDTO.Description) > 500 {
		return fmt.Errorf("la descripción del recurso no puede exceder 500 caracteres")
	}

	return nil
}
