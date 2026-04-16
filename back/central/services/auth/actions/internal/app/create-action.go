package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/auth/actions/internal/domain"
)

// CreateAction crea un nuevo action
func (uc *ActionUseCase) CreateAction(ctx context.Context, createDTO domain.CreateActionDTO) (*domain.ActionDTO, error) {
	uc.logger.Info().Str("name", createDTO.Name).Msg("Iniciando creación de action")

	// Validar datos de entrada
	if err := uc.validateCreateAction(createDTO); err != nil {
		uc.logger.Error().Err(err).Str("name", createDTO.Name).Msg("Validación fallida para crear action")
		return nil, err
	}

	// Verificar que no existe un action con el mismo nombre
	existingAction, err := uc.repository.GetActionByName(ctx, createDTO.Name)
	if err == nil && existingAction != nil {
		uc.logger.Warn().Str("name", createDTO.Name).Msg("Action ya existe con ese nombre")
		return nil, fmt.Errorf("ya existe un action con el nombre '%s'", createDTO.Name)
	}

	// Crear entidad de dominio
	action := domain.Action{
		Name:        strings.TrimSpace(createDTO.Name),
		Description: strings.TrimSpace(createDTO.Description),
	}

	// Crear action en el repositorio
	actionID, err := uc.repository.CreateAction(ctx, action)
	if err != nil {
		uc.logger.Error().Err(err).Str("name", createDTO.Name).Msg("Error al crear action")
		return nil, fmt.Errorf("error al crear action: %w", err)
	}

	// Obtener el action creado para devolver el DTO completo
	createdAction, err := uc.repository.GetActionByID(ctx, actionID)
	if err != nil {
		uc.logger.Error().Err(err).Uint("action_id", actionID).Msg("Error al obtener action creado")
		return nil, fmt.Errorf("error al obtener action creado: %w", err)
	}

	// Convertir a DTO
	actionDTO := &domain.ActionDTO{
		ID:          createdAction.ID,
		Name:        createdAction.Name,
		Description: createdAction.Description,
		CreatedAt:   createdAction.CreatedAt,
		UpdatedAt:   createdAction.UpdatedAt,
	}

	uc.logger.Info().
		Uint("action_id", actionID).
		Str("name", createDTO.Name).
		Msg("Action creado exitosamente")

	return actionDTO, nil
}

// validateCreateAction valida los datos para crear un action
func (uc *ActionUseCase) validateCreateAction(createDTO domain.CreateActionDTO) error {
	if strings.TrimSpace(createDTO.Name) == "" {
		return fmt.Errorf("el nombre del action es obligatorio")
	}

	if len(createDTO.Name) > 20 {
		return fmt.Errorf("el nombre del action no puede exceder 20 caracteres")
	}

	if len(createDTO.Description) > 255 {
		return fmt.Errorf("la descripción del action no puede exceder 255 caracteres")
	}

	return nil
}
