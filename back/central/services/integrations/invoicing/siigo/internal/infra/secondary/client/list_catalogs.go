package client

import (
	"context"
	"fmt"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
)

type documentTypeResponse struct {
	ID            int    `json:"id"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Active        bool   `json:"active"`
	ElectronicTyp string `json:"electronic_type"`
}

type userResponse struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Active    bool   `json:"active"`
}

type taxResponse struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	Percentage float64 `json:"percentage"`
	Active     bool    `json:"active"`
}

type costCenterResponse struct {
	ID     int    `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

func (c *Client) ListDocumentTypes(ctx context.Context, credentials dtos.Credentials, documentType string) ([]dtos.CatalogItem, error) {
	var listResp []documentTypeResponse
	if err := c.getCatalog(ctx, credentials, "/v1/document-types", map[string]string{"type": documentType}, &listResp); err != nil {
		return nil, err
	}

	items := make([]dtos.CatalogItem, 0, len(listResp))
	for _, r := range listResp {
		if !r.Active {
			continue
		}
		detalle := r.Description
		if r.ElectronicTyp != "" {
			if r.ElectronicTyp == "NoElectronic" {
				detalle = detalle + " (no electronico)"
			} else {
				detalle = detalle + " (electronico)"
			}
		}
		items = append(items, dtos.CatalogItem{
			ID:     r.ID,
			Code:   r.Code,
			Name:   r.Name,
			Detail: detalle,
		})
	}
	return items, nil
}

func (c *Client) ListSellers(ctx context.Context, credentials dtos.Credentials) ([]dtos.CatalogItem, error) {
	var paged struct {
		Results []userResponse `json:"results"`
	}
	if err := c.getCatalog(ctx, credentials, "/v1/users", map[string]string{"page_size": "100"}, &paged); err != nil {
		return nil, err
	}

	items := make([]dtos.CatalogItem, 0, len(paged.Results))
	for _, r := range paged.Results {
		if !r.Active {
			continue
		}
		name := r.FirstName
		if r.LastName != "" {
			name = name + " " + r.LastName
		}
		if name == " " || name == "" {
			name = r.Username
		}
		items = append(items, dtos.CatalogItem{
			ID:     r.ID,
			Name:   name,
			Detail: r.Email,
		})
	}
	return items, nil
}

func (c *Client) ListTaxes(ctx context.Context, credentials dtos.Credentials) ([]dtos.CatalogItem, error) {
	var listResp []taxResponse
	if err := c.getCatalog(ctx, credentials, "/v1/taxes", nil, &listResp); err != nil {
		return nil, err
	}

	items := make([]dtos.CatalogItem, 0, len(listResp))
	for _, r := range listResp {
		if !r.Active {
			continue
		}
		items = append(items, dtos.CatalogItem{
			ID:      r.ID,
			Name:    r.Name,
			Detail:  r.Type,
			Percent: strconv.FormatFloat(r.Percentage, 'f', -1, 64),
		})
	}
	return items, nil
}

func (c *Client) ListCostCenters(ctx context.Context, credentials dtos.Credentials) ([]dtos.CatalogItem, error) {
	var listResp []costCenterResponse
	if err := c.getCatalog(ctx, credentials, "/v1/cost-centers", nil, &listResp); err != nil {
		return nil, err
	}

	items := make([]dtos.CatalogItem, 0, len(listResp))
	for _, r := range listResp {
		if !r.Active {
			continue
		}
		items = append(items, dtos.CatalogItem{
			ID:   r.ID,
			Code: r.Code,
			Name: r.Name,
		})
	}
	return items, nil
}

func (c *Client) getCatalog(ctx context.Context, credentials dtos.Credentials, path string, query map[string]string, out interface{}) error {
	token, err := c.authenticate(ctx, credentials.Username, credentials.AccessKey, credentials.AccountID, credentials.PartnerID, credentials.BaseURL)
	if err != nil {
		return fmt.Errorf("failed to authenticate with Siigo: %w", err)
	}

	req := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Partner-Id", credentials.PartnerID).
		SetResult(out)

	for k, v := range query {
		if v != "" {
			req = req.SetQueryParam(k, v)
		}
	}

	resp, err := req.Get(c.endpointURL(credentials.BaseURL, path))
	if err != nil {
		c.log.Error(ctx).Err(err).Str("path", path).Msg("Siigo catalog request failed - network error")
		return fmt.Errorf("error de red al consultar %s en Siigo: %w", path, err)
	}

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("path", path).
			Str("body", string(resp.Body())).
			Msg("Siigo catalog request failed")
		return fmt.Errorf("error al consultar %s en Siigo (codigo %d)", path, resp.StatusCode())
	}

	return nil
}
