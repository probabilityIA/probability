package usecaseordermapping

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

// GetOrCreateCustomer verifica si el cliente existe, si no, lo crea
func (uc *UseCaseOrderMapping) GetOrCreateCustomer(ctx context.Context, businessID uint, dto *dtos.ProbabilityOrderDTO) (*entities.Client, error) {
	// 1. Buscar cliente existente por email (si hay email)
	if dto.CustomerEmail != "" {
		client, err := uc.repo.GetClientByEmail(ctx, businessID, dto.CustomerEmail)
		if err != nil {
			return nil, fmt.Errorf("error searching client by email: %w", err)
		}
		if client != nil {
			return client, nil
		}
	}

	// 2. Buscar por DNI si está disponible
	if dto.CustomerDNI != "" {
		clientByDNI, err := uc.repo.GetClientByDNI(ctx, businessID, dto.CustomerDNI)
		if err != nil {
			return nil, fmt.Errorf("error searching client by DNI: %w", err)
		}
		if clientByDNI != nil {
			return clientByDNI, nil
		}
	}

	// 3. Si no hay email ni DNI ni nombre, no se puede crear cliente
	if dto.CustomerName == "" && dto.CustomerFirstName == "" {
		return nil, nil
	}

	// 4. Crear nuevo cliente
	var dni *string
	if dto.CustomerDNI != "" {
		dni = &dto.CustomerDNI
	}

	var email *string
	if dto.CustomerEmail != "" {
		email = &dto.CustomerEmail
	}

	name := dto.CustomerName
	if name == "" {
		name = fmt.Sprintf("%s %s", dto.CustomerFirstName, dto.CustomerLastName)
	}

	newClient := &entities.Client{
		BusinessID: businessID,
		Name:       name,
		Email:      email,
		Phone:      dto.CustomerPhone,
		Dni:        dni,
	}

	if err := uc.repo.CreateClient(ctx, newClient); err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	return newClient, nil
}
