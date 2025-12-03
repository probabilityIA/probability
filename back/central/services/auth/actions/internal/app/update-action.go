package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/auth/actions/internal/domain"
)

// UpdateAction actualiza un action existente
func (uc *ActionUseCase) UpdateAction(ctx context.Context, id uint, updateDTO domain.UpdateActionDTO) (*domain.ActionDTO, error) {
	uc.logger.Info().Uint("action_id", id).Str("name", updateDTO.Name).Msg("Iniciando actualización de action")

	// Validar datos de entrada
	if err := uc.validateUpdateAction(updateDTO); err != nil {
		uc.logger.Error().Err(err).Uint("action_id", id).Msg("Validación fallida para actualizar action")
		return nil, err
	}

	// Verificar que el action existe
	existingAction, err := uc.repository.GetActionByID(ctx, id)
	if err != nil {
		uc.logger.Error().Err(err).Uint("action_id", id).Msg("Error al obtener action para actualizar")
		return nil, fmt.Errorf("action con ID %d no encontrado", id)
	}

	if existingAction == nil {
		uc.logger.Warn().Uint("action_id", id).Msg("Action no encontrado para actualizar")
		return nil, fmt.Errorf("action con ID %d no encontrado", id)
	}

	// Verificar que no existe otro action con el mismo nombre (si se está cambiando el nombre)
	if strings.TrimSpace(updateDTO.Name) != existingAction.Name {
		duplicateAction, err := uc.repository.GetActionByName(ctx, updateDTO.Name)
		if err == nil && duplicateAction != nil && duplicateAction.ID != id {
			uc.logger.Warn().Str("name", updateDTO.Name).Uint("action_id", id).Msg("Otro action ya existe con ese nombre")
			return nil, fmt.Errorf("ya existe otro action con el nombre '%s'", updateDTO.Name)
		}
	}

	// Crear entidad de dominio con los datos actualizados
	action := domain.Action{
		ID:          id,
		Name:        strings.TrimSpace(updateDTO.Name),
		Description: strings.TrimSpace(updateDTO.Description),
	}

	// Actualizar action en el repositorio
	_, err = uc.repository.UpdateAction(ctx, id, action)
	if err != nil {
		uc.logger.Error().Err(err).Uint("action_id", id).Msg("Error al actualizar action")
		return nil, fmt.Errorf("error al actualizar action: %w", err)
	}

	// Obtener el action actualizado para devolver el DTO completo
	updatedAction, err := uc.repository.GetActionByID(ctx, id)
	if err != nil {
		uc.logger.Error().Err(err).Uint("action_id", id).Msg("Error al obtener action actualizado")
		return nil, fmt.Errorf("error al obtener action actualizado: %w", err)
	}

	// Convertir a DTO
	actionDTO := &domain.ActionDTO{
		ID:          updatedAction.ID,
		Name:        updatedAction.Name,
		Description: updatedAction.Description,
		CreatedAt:   updatedAction.CreatedAt,
		UpdatedAt:   updatedAction.UpdatedAt,
	}

	uc.logger.Info().
		Uint("action_id", id).
		Str("name", updateDTO.Name).
		Msg("Action actualizado exitosamente")

	return actionDTO, nil
}

// validateUpdateAction valida los datos para actualizar un action
func (uc *ActionUseCase) validateUpdateAction(updateDTO domain.UpdateActionDTO) error {
	if strings.TrimSpace(updateDTO.Name) == "" {
		return fmt.Errorf("el nombre del action es obligatorio")
	}

	if len(updateDTO.Name) > 20 {
		return fmt.Errorf("el nombre del action no puede exceder 20 caracteres")
	}

	if len(updateDTO.Description) > 255 {
		return fmt.Errorf("la descripción del action no puede exceder 255 caracteres")
	}

	return nil
}
