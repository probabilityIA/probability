package usecaseordermapping

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain"
)

// GetOrCreateCustomer verifica si el cliente existe, si no, lo crea
func (uc *UseCaseOrderMapping) GetOrCreateCustomer(ctx context.Context, businessID uint, dto *domain.CanonicalOrderDTO) (*domain.Client, error) {
	if dto.CustomerEmail == "" {
		return nil, nil // No hay email para validar
	}

	// 1. Buscar cliente existente
	client, err := uc.repo.GetClientByEmail(ctx, businessID, dto.CustomerEmail)
	if err != nil {
		return nil, fmt.Errorf("error searching client: %w", err)
	}

	if client != nil {
		return client, nil
	}

	// 2. Crear nuevo cliente si no existe
	newClient := &domain.Client{
		BusinessID: businessID,
		Name:       dto.CustomerName,
		Email:      dto.CustomerEmail,
		Phone:      dto.CustomerPhone,
		Dni:        &dto.CustomerDNI,
	}

	if err := uc.repo.CreateClient(ctx, newClient); err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	return newClient, nil
}
