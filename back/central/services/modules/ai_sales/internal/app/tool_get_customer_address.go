package app

import (
	"context"
	"encoding/json"
	"fmt"
)

type getCustomerLastAddressInput struct {
	CustomerID uint `json:"customer_id"`
}

// executeGetCustomerLastAddress obtiene la ultima direccion de envio de un cliente
func executeGetCustomerLastAddress(ctx context.Context, inputJSON string, deps *toolDeps) (string, error) {
	var input getCustomerLastAddressInput
	if err := parseToolInput(inputJSON, &input); err != nil {
		return "", fmt.Errorf("error parsing GetCustomerLastAddress input: %w", err)
	}

	if input.CustomerID == 0 {
		return `{"error": "customer_id es requerido"}`, nil
	}

	address, err := deps.customerRepo.GetCustomerLastAddress(ctx, deps.businessID, input.CustomerID)
	if err != nil {
		return "", fmt.Errorf("error getting customer address: %w", err)
	}

	if address == nil {
		return `{"found": false, "message": "El cliente no tiene direcciones de envio anteriores"}`, nil
	}

	response, err := json.Marshal(map[string]any{
		"found":       true,
		"street":      address.Street,
		"city":        address.City,
		"state":       address.State,
		"country":     address.Country,
		"postal_code": address.PostalCode,
		"order_date":  address.OrderDate,
		"message":     fmt.Sprintf("Ultima direccion usada en pedido del %s", address.OrderDate),
	})
	if err != nil {
		return "", fmt.Errorf("error marshaling address: %w", err)
	}

	return string(response), nil
}
