package usecaseordermapping

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/orders/domain"
)

// GetOrCreateCustomer verifica si el cliente existe, si no, lo crea
func (uc *UseCaseOrderMapping) GetOrCreateCustomer(ctx context.Context, businessID uint, dto *domain.CanonicalOrderDTO) (*domain.Client, error) {
	if dto.CustomerEmail == "" {
		return nil, nil // No hay email para validar
	}

	// 1. Buscar cliente existente por email
	client, err := uc.repo.GetClientByEmail(ctx, businessID, dto.CustomerEmail)
	if err != nil {
		return nil, fmt.Errorf("error searching client by email: %w", err)
	}

	if client != nil {
		return client, nil
	}

	// 2. Si no existe por email y hay DNI, buscar por DNI antes de crear
	if dto.CustomerDNI != "" {
		clientByDNI, err := uc.repo.GetClientByDNI(ctx, businessID, dto.CustomerDNI)
		if err != nil {
			return nil, fmt.Errorf("error searching client by DNI: %w", err)
		}

		if clientByDNI != nil {
			return clientByDNI, nil
		}
	}

	// 3. Crear nuevo cliente solo si no existe ni por email ni por DNI
	// Handle DNI: Save as NULL if empty to avoid unique constraint violations
	var dni *string
	if dto.CustomerDNI != "" {
		dni = &dto.CustomerDNI
	}

	newClient := &domain.Client{
		BusinessID: businessID,
		Name:       dto.CustomerName,
		Email:      dto.CustomerEmail,
		Phone:      dto.CustomerPhone,
		Dni:        dni,
	}

	if err := uc.repo.CreateClient(ctx, newClient); err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	return newClient, nil
}
