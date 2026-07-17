package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/infra/secondary/client/response"
)

const (
	productsPageSize = 100
	maxProductPages  = 100
)

func (c *JumpsellerClient) GetProducts(ctx context.Context, cred domain.Credential) ([]domain.JumpsellerProduct, error) {
	all := make([]domain.JumpsellerProduct, 0)

	for page := 1; page <= maxProductPages; page++ {
		query := url.Values{}
		query.Set("page", strconv.Itoa(page))
		query.Set("limit", strconv.Itoa(productsPageSize))

		raw, err := c.do(ctx, cred, http.MethodGet, "/products.json", query, nil)
		if err != nil {
			return nil, err
		}

		var envelopes []response.ProductEnvelope
		if err := json.Unmarshal(raw, &envelopes); err != nil {
			return nil, fmt.Errorf("jumpseller client: parsing products: %w", err)
		}

		if len(envelopes) == 0 {
			break
		}

		for _, envelope := range envelopes {
			all = append(all, envelope.Product.ToDomain())
		}

		if len(envelopes) < productsPageSize {
			break
		}
	}

	return all, nil
}

func (c *JumpsellerClient) searchProductsBySKU(ctx context.Context, cred domain.Credential, sku string) ([]domain.JumpsellerProduct, error) {
	query := url.Values{}
	query.Set("query", sku)
	query.Set("fields", "sku")

	raw, err := c.do(ctx, cred, http.MethodGet, "/products/search.json", query, nil)
	if err != nil {
		return nil, err
	}

	var envelopes []response.ProductEnvelope
	if err := json.Unmarshal(raw, &envelopes); err != nil {
		return nil, fmt.Errorf("jumpseller client: parsing product search: %w", err)
	}

	products := make([]domain.JumpsellerProduct, 0, len(envelopes))
	for _, envelope := range envelopes {
		products = append(products, envelope.Product.ToDomain())
	}
	return products, nil
}

func (c *JumpsellerClient) ResolveStockTarget(ctx context.Context, cred domain.Credential, sku string) (*domain.StockTarget, error) {
	products, err := c.searchProductsBySKU(ctx, cred, sku)
	if err != nil {
		return nil, err
	}

	for _, product := range products {
		if product.SKU == sku {
			return &domain.StockTarget{ProductID: product.ID, Found: true}, nil
		}
		for _, variant := range product.Variants {
			if variant.SKU == sku {
				return &domain.StockTarget{ProductID: product.ID, VariantID: variant.ID, Found: true}, nil
			}
		}
	}

	return &domain.StockTarget{Found: false}, nil
}

func (c *JumpsellerClient) SetProductStock(ctx context.Context, cred domain.Credential, productID int64, stock int) error {
	body := response.UpdateProductStockRequest{
		Product: response.UpdateStockFields{Stock: stock},
	}
	_, err := c.do(ctx, cred, http.MethodPut, fmt.Sprintf("/products/%d.json", productID), nil, body)
	return err
}

func (c *JumpsellerClient) SetVariantStock(ctx context.Context, cred domain.Credential, productID, variantID int64, stock int) error {
	body := response.UpdateVariantStockRequest{
		Variant: response.UpdateStockFields{Stock: stock},
	}
	_, err := c.do(ctx, cred, http.MethodPut, fmt.Sprintf("/products/%d/variants/%d.json", productID, variantID), nil, body)
	return err
}

func (c *JumpsellerClient) CreateProduct(ctx context.Context, cred domain.Credential, input domain.CreateProductInput) (string, error) {
	body := response.CreateProductRequest{
		Product: response.CreateProductFields{
			Name:           input.Name,
			SKU:            input.SKU,
			Price:          input.Price,
			Description:    input.Description,
			Stock:          input.StockQuantity,
			StockUnlimited: !input.ManageStock,
			Status:         "available",
			Weight:         input.Weight,
			Height:         input.Height,
			Width:          input.Width,
			Length:         input.Length,
		},
	}

	raw, err := c.do(ctx, cred, http.MethodPost, "/products.json", nil, body)
	if err != nil {
		return "", err
	}

	var envelope response.ProductEnvelope
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return "", fmt.Errorf("jumpseller client: parsing created product: %w", err)
	}

	return strconv.FormatInt(envelope.Product.ID, 10), nil
}

func (c *JumpsellerClient) UpdateProduct(ctx context.Context, cred domain.Credential, productID int64, input domain.UpdateProductInput) error {
	body := response.UpdateProductRequest{
		Product: response.UpdateProductFields{
			Name:        input.Name,
			Price:       input.Price,
			Description: input.Description,
			Weight:      input.Weight,
			Height:      input.Height,
			Width:       input.Width,
			Length:      input.Length,
		},
	}

	_, err := c.do(ctx, cred, http.MethodPut, fmt.Sprintf("/products/%d.json", productID), nil, body)
	return err
}
