package usecaseorder

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

// getOrCreateClient busca un cliente existente por email o DNI; si no existe, lo crea.
// Esto garantiza que las órdenes manuales también generen registros en la tabla clients.
func (uc *UseCaseOrder) getOrCreateClient(ctx context.Context, businessID uint, req *dtos.CreateOrderRequest) (*entities.Client, error) {
	// 1. Buscar por email (si hay email)
	if req.CustomerEmail != "" {
		client, err := uc.repo.GetClientByEmail(ctx, businessID, req.CustomerEmail)
		if err != nil {
			return nil, fmt.Errorf("error searching client by email: %w", err)
		}
		if client != nil {
			return client, nil
		}
	}

	// 2. Buscar por DNI si está disponible
	if req.CustomerDNI != "" {
		client, err := uc.repo.GetClientByDNI(ctx, businessID, req.CustomerDNI)
		if err != nil {
			return nil, fmt.Errorf("error searching client by DNI: %w", err)
		}
		if client != nil {
			return client, nil
		}
	}

	// 3. Crear nuevo cliente
	var dni *string
	if req.CustomerDNI != "" {
		dni = &req.CustomerDNI
	}

	var email *string
	if req.CustomerEmail != "" {
		email = &req.CustomerEmail
	}

	name := req.CustomerName
	if name == "" {
		name = fmt.Sprintf("%s %s", req.CustomerFirstName, req.CustomerLastName)
	}

	newClient := &entities.Client{
		BusinessID: businessID,
		Name:       name,
		Email:      email,
		Phone:      req.CustomerPhone,
		Dni:        dni,
	}

	if err := uc.repo.CreateClient(ctx, newClient); err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	return newClient, nil
}
