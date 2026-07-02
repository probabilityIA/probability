package client

import (
	"context"
	"fmt"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
)

type listProductsResponse struct {
	Pagination struct {
		Page         int `json:"page"`
		PageSize     int `json:"page_size"`
		TotalResults int `json:"total_results"`
	} `json:"pagination"`
	Results []struct {
		ID          string `json:"id"`
		Code        string `json:"code"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Prices      []struct {
			PriceList []struct {
				Position int     `json:"position"`
				Value    float64 `json:"value"`
			} `json:"price_list"`
		} `json:"prices"`
		StockControl      bool    `json:"stock_control"`
		AvailableQuantity float64 `json:"available_quantity"`
		Warehouses        []struct {
			ID       int     `json:"id"`
			Name     string  `json:"name"`
			Quantity float64 `json:"quantity"`
		} `json:"warehouses"`
	} `json:"results"`
}

func (c *Client) ListProducts(ctx context.Context, credentials dtos.Credentials, page, pageSize int) ([]dtos.ProductItem, error) {
	c.log.Info(ctx).
		Int("page", page).
		Int("page_size", pageSize).
		Msg("Listing Siigo products")

	token, err := c.authenticate(ctx, credentials.Username, credentials.AccessKey, credentials.AccountID, credentials.PartnerID, credentials.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with Siigo: %w", err)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 100
	}

	var listResp listProductsResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Partner-Id", credentials.PartnerID).
		SetQueryParam("page", strconv.Itoa(page)).
		SetQueryParam("page_size", strconv.Itoa(pageSize)).
		SetResult(&listResp).
		Get(c.endpointURL(credentials.BaseURL, "/v1/products"))

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Siigo list products request failed - network error")
		return nil, fmt.Errorf("error de red al listar productos en Siigo: %w", err)
	}

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("body", string(resp.Body())).
			Msg("Siigo list products failed")
		return nil, fmt.Errorf("error al listar productos en Siigo (codigo %d)", resp.StatusCode())
	}

	items := make([]dtos.ProductItem, 0, len(listResp.Results))
	for _, r := range listResp.Results {
		price := 0.0
		if len(r.Prices) > 0 && len(r.Prices[0].PriceList) > 0 {
			price = r.Prices[0].PriceList[0].Value
		}
		warehouses := make([]dtos.ProductWarehouseStock, 0, len(r.Warehouses))
		for _, w := range r.Warehouses {
			warehouses = append(warehouses, dtos.ProductWarehouseStock{
				ID:       w.ID,
				Name:     w.Name,
				Quantity: w.Quantity,
			})
		}
		items = append(items, dtos.ProductItem{
			ID:                r.ID,
			Code:              r.Code,
			Name:              r.Name,
			Description:       r.Description,
			Price:             price,
			StockControl:      r.StockControl,
			AvailableQuantity: r.AvailableQuantity,
			Warehouses:        warehouses,
		})
	}

	c.log.Info(ctx).
		Int("count", len(items)).
		Int("total_results", listResp.Pagination.TotalResults).
		Msg("Siigo products listed")

	return items, nil
}
