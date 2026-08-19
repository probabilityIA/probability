package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/infra/secondary/client/response"
)

const (
	productsPageSize = 200
	maxProductPages  = 200
)

func (c *TiendanubeClient) GetProducts(ctx context.Context, cred domain.Credential) ([]domain.TiendanubeProduct, error) {
	all := make([]domain.TiendanubeProduct, 0)

	for page := 1; page <= maxProductPages; page++ {
		query := url.Values{}
		query.Set("page", strconv.Itoa(page))
		query.Set("per_page", strconv.Itoa(productsPageSize))

		raw, _, err := c.do(ctx, cred, http.MethodGet, "/products", query, nil)
		if err != nil {
			return nil, err
		}

		var items []response.Product
		if err := json.Unmarshal(raw, &items); err != nil {
			return nil, fmt.Errorf("tiendanube client: parsing products: %w", err)
		}

		if len(items) == 0 {
			break
		}

		for _, item := range items {
			all = append(all, item.ToDomain())
		}

		if len(items) < productsPageSize {
			break
		}
	}

	return all, nil
}

func (c *TiendanubeClient) ResolveStockTarget(ctx context.Context, cred domain.Credential, sku string) (*domain.StockTarget, error) {
	trimmed := strings.TrimSpace(sku)
	if trimmed == "" {
		return &domain.StockTarget{Found: false}, nil
	}

	query := url.Values{}
	query.Set("q", trimmed)
	query.Set("per_page", "50")

	raw, _, err := c.do(ctx, cred, http.MethodGet, "/products", query, nil)
	if err != nil {
		return nil, err
	}

	var items []response.Product
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, fmt.Errorf("tiendanube client: parsing product search: %w", err)
	}

	for _, item := range items {
		product := item.ToDomain()
		for _, variant := range product.Variants {
			if strings.EqualFold(variant.SKU, trimmed) {
				return &domain.StockTarget{ProductID: product.ID, VariantID: variant.ID, Found: true}, nil
			}
		}
	}

	return &domain.StockTarget{Found: false}, nil
}

type variantStockBody struct {
	Stock interface{} `json:"stock"`
}

func (c *TiendanubeClient) SetVariantStock(ctx context.Context, cred domain.Credential, productID, variantID int64, stock int) error {
	if productID <= 0 || variantID <= 0 {
		return fmt.Errorf("tiendanube client: product_id y variant_id son requeridos para escribir stock")
	}
	path := fmt.Sprintf("/products/%d/variants/%d", productID, variantID)
	_, _, err := c.do(ctx, cred, http.MethodPut, path, nil, variantStockBody{Stock: stock})
	return err
}

type createVariantBody struct {
	SKU     string   `json:"sku,omitempty"`
	Barcode string   `json:"barcode,omitempty"`
	Price   string   `json:"price"`
	Stock   *int     `json:"stock"`
	Weight  *float64 `json:"weight,omitempty"`
	Depth   *float64 `json:"depth,omitempty"`
	Width   *float64 `json:"width,omitempty"`
	Height  *float64 `json:"height,omitempty"`
}

type createProductBody struct {
	Name        map[string]string   `json:"name"`
	Description map[string]string   `json:"description,omitempty"`
	Published   bool                `json:"published"`
	Variants    []createVariantBody `json:"variants"`
	Images      []createImageBody   `json:"images,omitempty"`
}

type createImageBody struct {
	Src string `json:"src"`
}

func formatPrice(value float64) string {
	return strconv.FormatFloat(value, 'f', 2, 64)
}

func (c *TiendanubeClient) CreateProduct(ctx context.Context, cred domain.Credential, input domain.CreateProductInput) (int64, int64, error) {
	variant := createVariantBody{
		SKU:     input.SKU,
		Barcode: input.Barcode,
		Price:   formatPrice(input.Price),
		Weight:  input.Weight,
		Depth:   input.Length,
		Width:   input.Width,
		Height:  input.Height,
	}
	if input.ManageStock {
		stock := input.StockQuantity
		variant.Stock = &stock
	}

	body := createProductBody{
		Name:      map[string]string{"es": input.Name},
		Published: true,
		Variants:  []createVariantBody{variant},
	}
	if input.Description != "" {
		body.Description = map[string]string{"es": input.Description}
	}
	if input.ImageURL != "" {
		body.Images = []createImageBody{{Src: input.ImageURL}}
	}

	raw, _, err := c.do(ctx, cred, http.MethodPost, "/products", nil, body)
	if err != nil {
		return 0, 0, err
	}

	var created response.Product
	if err := json.Unmarshal(raw, &created); err != nil {
		return 0, 0, fmt.Errorf("tiendanube client: parsing created product: %w", err)
	}

	variantID := int64(0)
	if len(created.Variants) > 0 {
		variantID = created.Variants[0].ID
	}

	return created.ID, variantID, nil
}

type updateProductBody struct {
	Name        map[string]string `json:"name,omitempty"`
	Description map[string]string `json:"description,omitempty"`
}

func (c *TiendanubeClient) UpdateProduct(ctx context.Context, cred domain.Credential, productID int64, input domain.UpdateProductInput) error {
	if productID <= 0 {
		return fmt.Errorf("tiendanube client: product_id es requerido para actualizar")
	}

	body := updateProductBody{}
	if input.Name != "" {
		body.Name = map[string]string{"es": input.Name}
	}
	if input.Description != "" {
		body.Description = map[string]string{"es": input.Description}
	}
	if body.Name == nil && body.Description == nil {
		return nil
	}

	_, _, err := c.do(ctx, cred, http.MethodPut, fmt.Sprintf("/products/%d", productID), nil, body)
	return err
}

type updateVariantBody struct {
	Price   *string  `json:"price,omitempty"`
	Barcode string   `json:"barcode,omitempty"`
	Weight  *float64 `json:"weight,omitempty"`
	Depth   *float64 `json:"depth,omitempty"`
	Width   *float64 `json:"width,omitempty"`
	Height  *float64 `json:"height,omitempty"`
}

func (c *TiendanubeClient) UpdateVariant(ctx context.Context, cred domain.Credential, productID, variantID int64, input domain.UpdateVariantInput) error {
	if productID <= 0 || variantID <= 0 {
		return fmt.Errorf("tiendanube client: product_id y variant_id son requeridos para actualizar la variante")
	}

	body := updateVariantBody{
		Barcode: input.Barcode,
		Weight:  input.Weight,
		Depth:   input.Depth,
		Width:   input.Width,
		Height:  input.Height,
	}
	if input.Price != nil {
		price := formatPrice(*input.Price)
		body.Price = &price
	}

	_, _, err := c.do(ctx, cred, http.MethodPut, fmt.Sprintf("/products/%d/variants/%d", productID, variantID), nil, body)
	return err
}

func (c *TiendanubeClient) GetProductsStock(ctx context.Context, cred domain.Credential, externalIDs []string) ([]domain.ChannelStock, error) {
	if len(externalIDs) == 0 {
		return []domain.ChannelStock{}, nil
	}

	wanted := make(map[string]bool, len(externalIDs))
	productIDs := make([]string, 0, len(externalIDs))
	seen := make(map[string]bool, len(externalIDs))
	for _, raw := range externalIDs {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			continue
		}
		wanted[trimmed] = true
		productPart, _, _ := strings.Cut(trimmed, ":")
		if productPart != "" && !seen[productPart] {
			seen[productPart] = true
			productIDs = append(productIDs, productPart)
		}
	}

	found := make(map[string]domain.ChannelStock, len(wanted))
	const idsPerRequest = 25

	for start := 0; start < len(productIDs); start += idsPerRequest {
		end := start + idsPerRequest
		if end > len(productIDs) {
			end = len(productIDs)
		}

		query := url.Values{}
		query.Set("ids", strings.Join(productIDs[start:end], ","))
		query.Set("per_page", strconv.Itoa(idsPerRequest))

		raw, _, err := c.do(ctx, cred, http.MethodGet, "/products", query, nil)
		if err != nil {
			return nil, err
		}

		var items []response.Product
		if err := json.Unmarshal(raw, &items); err != nil {
			return nil, fmt.Errorf("tiendanube client: parsing products stock: %w", err)
		}

		for _, item := range items {
			product := item.ToDomain()
			for _, variant := range product.Variants {
				key := strconv.FormatInt(product.ID, 10) + ":" + strconv.FormatInt(variant.ID, 10)
				entry := domain.ChannelStock{
					ExternalID:  key,
					Quantity:    variant.Stock,
					ManageStock: variant.StockManagement,
					Found:       true,
				}
				found[key] = entry
				if _, ok := found[strconv.FormatInt(product.ID, 10)]; !ok {
					plain := entry
					plain.ExternalID = strconv.FormatInt(product.ID, 10)
					found[plain.ExternalID] = plain
				}
			}
		}
	}

	out := make([]domain.ChannelStock, 0, len(wanted))
	for id := range wanted {
		if entry, ok := found[id]; ok {
			entry.ExternalID = id
			out = append(out, entry)
			continue
		}
		out = append(out, domain.ChannelStock{ExternalID: id})
	}
	return out, nil
}
