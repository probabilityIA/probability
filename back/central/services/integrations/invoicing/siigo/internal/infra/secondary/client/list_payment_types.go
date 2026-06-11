package client

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
)

type paymentTypeResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Active   bool   `json:"active"`
	DueDate  bool   `json:"due_date"`
}

func (c *Client) ListPaymentTypes(ctx context.Context, credentials dtos.Credentials, documentType string) ([]dtos.PaymentTypeItem, error) {
	c.log.Info(ctx).
		Str("document_type", documentType).
		Msg("Listing Siigo payment types")

	token, err := c.authenticate(ctx, credentials.Username, credentials.AccessKey, credentials.AccountID, credentials.PartnerID, credentials.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with Siigo: %w", err)
	}

	var listResp []paymentTypeResponse

	req := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Partner-Id", credentials.PartnerID).
		SetResult(&listResp)

	if documentType != "" {
		req = req.SetQueryParam("document_type", documentType)
	}

	resp, err := req.Get(c.endpointURL(credentials.BaseURL, "/v1/payment-types"))

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Siigo list payment types request failed - network error")
		return nil, fmt.Errorf("error de red al listar medios de pago en Siigo: %w", err)
	}

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("body", string(resp.Body())).
			Msg("Siigo list payment types failed")
		return nil, fmt.Errorf("error al listar medios de pago en Siigo (codigo %d)", resp.StatusCode())
	}

	items := make([]dtos.PaymentTypeItem, 0, len(listResp))
	for _, r := range listResp {
		if !r.Active {
			continue
		}
		items = append(items, dtos.PaymentTypeItem{
			ID:   r.ID,
			Name: r.Name,
			Type: r.Type,
		})
	}

	c.log.Info(ctx).
		Int("count", len(items)).
		Msg("Siigo payment types listed")

	return items, nil
}
