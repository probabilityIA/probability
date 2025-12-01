package usecasebusinesstype

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/domain"
)

// GetBusinessTypes obtiene todos los tipos de negocio
func (uc *BusinessTypeUseCase) GetBusinessTypes(ctx context.Context) ([]domain.BusinessTypeResponse, error) {
	uc.log.Info().Msg("Obteniendo tipos de negocio")

	businessTypes, err := uc.repository.GetBusinessTypes(ctx)
	if err != nil {
		uc.log.Error().Err(err).Msg("Error al obtener tipos de negocio")
		return nil, fmt.Errorf("error al obtener tipos de negocio: %w", err)
	}

	// Convertir entidades a DTOs
	response := make([]domain.BusinessTypeResponse, len(businessTypes))
	for i, bt := range businessTypes {
		response[i] = domain.BusinessTypeResponse{
			ID:          bt.ID,
			Name:        bt.Name,
			Code:        bt.Code,
			Description: bt.Description,
			Icon:        bt.Icon,
			IsActive:    bt.IsActive,
			CreatedAt:   bt.CreatedAt,
			UpdatedAt:   bt.UpdatedAt,
		}
	}

	uc.log.Info().Int("count", len(response)).Msg("Tipos de negocio obtenidos exitosamente")
	return response, nil
}
