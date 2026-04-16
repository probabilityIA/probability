package app

import (
	"context"
	"encoding/json"
	"fmt"
)

type searchCustomerInput struct {
	Query string `json:"query"`
}

// executeSearchCustomer busca clientes existentes por DNI, email, telefono o nombre
func executeSearchCustomer(ctx context.Context, inputJSON string, deps *toolDeps) (string, error) {
	var input searchCustomerInput
	if err := parseToolInput(inputJSON, &input); err != nil {
		return "", fmt.Errorf("error parsing SearchCustomer input: %w", err)
	}

	if input.Query == "" {
		return `{"results": [], "message": "Debe proporcionar un termino de busqueda"}`, nil
	}

	customers, err := deps.customerRepo.SearchCustomers(ctx, deps.businessID, input.Query)
	if err != nil {
		return "", fmt.Errorf("error searching customers: %w", err)
	}

	if len(customers) == 0 {
		return `{"results": [], "message": "No se encontro ningun cliente con ese criterio de busqueda"}`, nil
	}

	type customerResult struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email,omitempty"`
		Phone string `json:"phone,omitempty"`
		DNI   string `json:"dni,omitempty"`
	}

	results := make([]customerResult, 0, len(customers))
	for _, c := range customers {
		results = append(results, customerResult{
			ID:    c.ID,
			Name:  c.Name,
			Email: c.Email,
			Phone: c.Phone,
			DNI:   c.DNI,
		})
	}

	response, err := json.Marshal(map[string]any{
		"results": results,
		"count":   len(results),
	})
	if err != nil {
		return "", fmt.Errorf("error marshaling customer results: %w", err)
	}

	return string(response), nil
}
