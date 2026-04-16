package app

import (
	"context"
	"math"

	"github.com/secamc93/probability/back/central/services/auth/actions/internal/domain"
)

// GetActions obtiene todos los actions con filtros y paginación
func (uc *ActionUseCase) GetActions(ctx context.Context, page, pageSize int, name string) (*domain.ActionListDTO, error) {
	uc.logger.Info().Int("page", page).Int("page_size", pageSize).Str("name", name).Msg("Iniciando obtención de actions")

	// Configurar valores por defecto para paginación
	if pageSize <= 0 {
		pageSize = 10
	}
	if page <= 0 {
		page = 1
	}

	// Obtener actions del repositorio
	actions, total, err := uc.repository.GetActions(ctx, page, pageSize, name)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Error al obtener actions")
		return nil, err
	}

	// Convertir a DTOs
	var actionDTOs []domain.ActionDTO
	for _, action := range actions {
		actionDTOs = append(actionDTOs, domain.ActionDTO{
			ID:          action.ID,
			Name:        action.Name,
			Description: action.Description,
			CreatedAt:   action.CreatedAt,
			UpdatedAt:   action.UpdatedAt,
		})
	}

	// Calcular total de páginas
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	result := &domain.ActionListDTO{
		Actions:    actionDTOs,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	uc.logger.Info().
		Int64("total", total).
		Int("returned", len(actionDTOs)).
		Int("page", page).
		Int("total_pages", totalPages).
		Msg("Actions obtenidos exitosamente")

	return result, nil
}
