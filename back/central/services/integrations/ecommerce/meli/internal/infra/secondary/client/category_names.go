package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

type meliCategoryDetail struct {
	Name         string `json:"name"`
	PathFromRoot []struct {
		Name string `json:"name"`
	} `json:"path_from_root"`
}

func (c *MeliClient) categoryName(ctx context.Context, accessToken, categoryID string) string {
	if categoryID == "" {
		return ""
	}

	endpoint := fmt.Sprintf("%s/categories/%s", c.baseURL, categoryID)
	req, err := c.newAuthorizedRequest(ctx, http.MethodGet, endpoint, accessToken)
	if err != nil {
		return ""
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var detalle meliCategoryDetail
	if err := json.NewDecoder(resp.Body).Decode(&detalle); err != nil {
		return ""
	}
	return detalle.Name
}

func (c *MeliClient) resolveCategoryNames(ctx context.Context, accessToken string, products []domain.MeliProduct) {
	nombres := make(map[string]string)
	for i := range products {
		id := products[i].CategoryID
		if id == "" {
			continue
		}
		nombre, visto := nombres[id]
		if !visto {
			nombre = c.categoryName(ctx, accessToken, id)
			nombres[id] = nombre
		}
		if nombre != "" {
			products[i].Category = nombre
		}
	}
}
