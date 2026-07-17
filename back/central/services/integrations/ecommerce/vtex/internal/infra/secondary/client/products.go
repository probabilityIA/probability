package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
)

const skuPageSize = 50

func (c *VTEXClient) ListSKUs(ctx context.Context, cred domain.Credential) ([]domain.VTEXSKU, error) {
	var all []domain.VTEXSKU
	page := 1

	for {
		endpoint := fmt.Sprintf("%s/api/catalog_system/pvt/sku/stockkeepingunitids?page=%d&pagesize=%d", baseURL(cred), page, skuPageSize)

		body, err := c.do(ctx, http.MethodGet, endpoint, cred, nil)
		if err != nil {
			return nil, err
		}

		var ids []int
		if err := json.Unmarshal(body, &ids); err != nil {
			return nil, fmt.Errorf("vtex client: parsing sku ids: %w", err)
		}

		if len(ids) == 0 {
			break
		}

		for _, id := range ids {
			sku, err := c.getSKUByID(ctx, cred, id)
			if err != nil {
				continue
			}
			all = append(all, *sku)
		}

		if len(ids) < skuPageSize {
			break
		}

		page++
		time.Sleep(300 * time.Millisecond)
	}

	return all, nil
}

func (c *VTEXClient) getSKUByID(ctx context.Context, cred domain.Credential, skuID int) (*domain.VTEXSKU, error) {
	endpoint := fmt.Sprintf("%s/api/catalog/pvt/stockkeepingunit/%d", baseURL(cred), skuID)

	body, err := c.do(ctx, http.MethodGet, endpoint, cred, nil)
	if err != nil {
		return nil, err
	}

	var raw struct {
		ID              int     `json:"Id"`
		ProductID       int     `json:"ProductId"`
		Name            string  `json:"Name"`
		RefID           string  `json:"RefId"`
		IsActive        bool    `json:"IsActive"`
		WeightKg        float64 `json:"WeightKg"`
		Length          float64 `json:"Length"`
		Width           float64 `json:"Width"`
		Height          float64 `json:"Height"`
		MeasurementUnit string  `json:"MeasurementUnit"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("vtex client: parsing sku %d: %w", skuID, err)
	}

	sku := &domain.VTEXSKU{
		ID:              fmt.Sprintf("%d", raw.ID),
		ProductID:       fmt.Sprintf("%d", raw.ProductID),
		Name:            raw.Name,
		RefID:           strings.TrimSpace(raw.RefID),
		IsActive:        raw.IsActive,
		MeasurementUnit: raw.MeasurementUnit,
	}
	if raw.WeightKg > 0 {
		w := raw.WeightKg
		sku.Weight = &w
	}
	if raw.Length > 0 {
		l := raw.Length
		sku.Length = &l
	}
	if raw.Width > 0 {
		w := raw.Width
		sku.Width = &w
	}
	if raw.Height > 0 {
		h := raw.Height
		sku.Height = &h
	}

	return sku, nil
}

func (c *VTEXClient) GetSKUIDByRefID(ctx context.Context, cred domain.Credential, refID string, isSeller bool) (string, error) {
	if strings.TrimSpace(refID) == "" {
		return "", domain.ErrSKUNotFound
	}
	if isSeller {
		return c.getSKUIDFromSeller(ctx, cred, refID)
	}
	return c.getSKUIDFromStore(ctx, cred, refID)
}

func (c *VTEXClient) getSKUIDFromStore(ctx context.Context, cred domain.Credential, refID string) (string, error) {
	endpoint := fmt.Sprintf("%s/api/catalog_system/pvt/sku/stockkeepingunitidbyrefid/%s", baseURL(cred), refID)

	body, err := c.do(ctx, http.MethodGet, endpoint, cred, nil)
	if err != nil {
		return "", err
	}

	trimmed := strings.TrimSpace(string(body))
	if trimmed == "" || trimmed == "null" {
		return "", domain.ErrSKUNotFound
	}

	var skuID string
	if err := json.Unmarshal(body, &skuID); err != nil {
		var numeric int
		if err2 := json.Unmarshal(body, &numeric); err2 != nil {
			return "", fmt.Errorf("vtex client: parsing sku id for ref %s: %w", refID, err)
		}
		skuID = fmt.Sprintf("%d", numeric)
	}

	if skuID == "" {
		return "", domain.ErrSKUNotFound
	}
	return skuID, nil
}

func (c *VTEXClient) getSKUIDFromSeller(ctx context.Context, cred domain.Credential, refID string) (string, error) {
	endpoint := fmt.Sprintf("%s/api/catalog-seller-portal/skus/_search?externalid=%s", baseURL(cred), refID)

	body, err := c.do(ctx, http.MethodGet, endpoint, cred, nil)
	if err != nil {
		return "", err
	}

	var result struct {
		Data []struct {
			ID         string `json:"id"`
			ExternalID string `json:"externalId"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("vtex client: parsing seller sku search: %w", err)
	}

	if len(result.Data) == 0 {
		return "", domain.ErrSKUNotFound
	}

	if !strings.EqualFold(result.Data[0].ExternalID, refID) {
		return "", domain.ErrSKUNotFound
	}

	return result.Data[0].ID, nil
}
