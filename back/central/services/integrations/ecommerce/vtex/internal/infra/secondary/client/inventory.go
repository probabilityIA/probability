package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
)

func (c *VTEXClient) GetWarehouses(ctx context.Context, cred domain.Credential) ([]domain.Warehouse, error) {
	endpoint := fmt.Sprintf("%s/api/logistics/pvt/configuration/warehouses", baseURL(cred))

	body, err := c.do(ctx, http.MethodGet, endpoint, cred, nil)
	if err != nil {
		return nil, err
	}

	var raw []struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		IsActive bool   `json:"isActive"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("vtex client: parsing warehouses: %w", err)
	}

	warehouses := make([]domain.Warehouse, 0, len(raw))
	for _, w := range raw {
		warehouses = append(warehouses, domain.Warehouse{
			ID:       w.ID,
			Name:     w.Name,
			IsActive: w.IsActive,
		})
	}

	return warehouses, nil
}

func (c *VTEXClient) UpdateSKUInventory(ctx context.Context, cred domain.Credential, skuID, warehouseID string, quantity int) error {
	if quantity < 0 {
		quantity = 0
	}

	endpoint := fmt.Sprintf("%s/api/logistics/pvt/inventory/skus/%s/warehouses/%s", baseURL(cred), skuID, warehouseID)
	payload := []byte(fmt.Sprintf(`{"quantity": %d, "unlimitedQuantity": false}`, quantity))

	_, err := c.do(ctx, http.MethodPut, endpoint, cred, payload)
	return err
}
