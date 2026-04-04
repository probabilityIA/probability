package app

import (
	"context"
	"encoding/json"
	"fmt"
)

type searchProductsInput struct {
	Query string `json:"query"`
	Limit int    `json:"limit"`
}

// executeSearchProducts busca productos en el catalogo del negocio
func executeSearchProducts(ctx context.Context, inputJSON string, deps *toolDeps) (string, error) {
	var input searchProductsInput
	if err := parseToolInput(inputJSON, &input); err != nil {
		return "", fmt.Errorf("error parsing SearchProducts input: %w", err)
	}

	if input.Limit <= 0 {
		input.Limit = 5
	}

	products, err := deps.productRepo.SearchProducts(ctx, deps.businessID, input.Query, input.Limit)
	if err != nil {
		return "", fmt.Errorf("error searching products: %w", err)
	}

	if len(products) == 0 {
		return `{"results": [], "message": "No se encontraron productos para la busqueda"}`, nil
	}

	type productResult struct {
		SKU       string  `json:"sku"`
		Name      string  `json:"name"`
		Price     float64 `json:"price"`
		Currency  string  `json:"currency"`
		Stock     int     `json:"stock"`
		Category  string  `json:"category"`
		Brand     string  `json:"brand"`
		Available bool    `json:"available"`
	}

	results := make([]productResult, 0, len(products))
	for _, p := range products {
		results = append(results, productResult{
			SKU:       p.SKU,
			Name:      p.Name,
			Price:     p.Price,
			Currency:  p.Currency,
			Stock:     p.StockQuantity,
			Category:  p.Category,
			Brand:     p.Brand,
			Available: p.StockQuantity > 0,
		})
	}

	response, err := json.Marshal(map[string]any{
		"results": results,
		"count":   len(results),
	})
	if err != nil {
		return "", fmt.Errorf("error marshaling search results: %w", err)
	}

	return string(response), nil
}
