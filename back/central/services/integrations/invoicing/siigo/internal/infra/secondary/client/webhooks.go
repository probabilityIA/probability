package client

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
)

type webhookResponse struct {
	ID            string `json:"id"`
	ApplicationID string `json:"application_id"`
	URL           string `json:"url"`
	Topic         string `json:"topic"`
	CompanyKey    string `json:"company_key"`
	Active        bool   `json:"active"`
	CreatedAt     string `json:"created_at"`
}

func (c *Client) ListWebhooks(ctx context.Context, credentials dtos.Credentials) ([]dtos.WebhookItem, error) {
	token, err := c.authenticate(ctx, credentials.Username, credentials.AccessKey, credentials.AccountID, credentials.PartnerID, credentials.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with Siigo: %w", err)
	}

	var listResp []webhookResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Partner-Id", credentials.PartnerID).
		SetResult(&listResp).
		Get(c.endpointURL(credentials.BaseURL, "/v1/webhooks"))

	if err != nil {
		return nil, fmt.Errorf("error de red al listar webhooks en Siigo: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error al listar webhooks en Siigo (codigo %d): %s", resp.StatusCode(), string(resp.Body()))
	}

	items := make([]dtos.WebhookItem, 0, len(listResp))
	for _, w := range listResp {
		items = append(items, dtos.WebhookItem{
			ID:            w.ID,
			ApplicationID: w.ApplicationID,
			URL:           w.URL,
			Topic:         w.Topic,
			CompanyKey:    w.CompanyKey,
			Active:        w.Active,
			CreatedAt:     w.CreatedAt,
		})
	}

	return items, nil
}

func (c *Client) CreateWebhook(ctx context.Context, credentials dtos.Credentials, input dtos.CreateWebhookInput) (*dtos.WebhookItem, error) {
	token, err := c.authenticate(ctx, credentials.Username, credentials.AccessKey, credentials.AccountID, credentials.PartnerID, credentials.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with Siigo: %w", err)
	}

	body := map[string]string{
		"application_id": input.ApplicationID,
		"url":            input.URL,
		"topic":          input.Topic,
	}

	var created webhookResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Partner-Id", credentials.PartnerID).
		SetBody(body).
		SetResult(&created).
		Post(c.endpointURL(credentials.BaseURL, "/v1/webhooks"))

	if err != nil {
		return nil, fmt.Errorf("error de red al crear webhook en Siigo: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error al crear webhook en Siigo (codigo %d): %s", resp.StatusCode(), string(resp.Body()))
	}

	return &dtos.WebhookItem{
		ID:            created.ID,
		ApplicationID: created.ApplicationID,
		URL:           created.URL,
		Topic:         created.Topic,
		CompanyKey:    created.CompanyKey,
		Active:        created.Active,
		CreatedAt:     created.CreatedAt,
	}, nil
}

func (c *Client) DeleteWebhook(ctx context.Context, credentials dtos.Credentials, webhookID string) error {
	token, err := c.authenticate(ctx, credentials.Username, credentials.AccessKey, credentials.AccountID, credentials.PartnerID, credentials.BaseURL)
	if err != nil {
		return fmt.Errorf("failed to authenticate with Siigo: %w", err)
	}

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Partner-Id", credentials.PartnerID).
		Delete(c.endpointURL(credentials.BaseURL, "/v1/webhooks/"+webhookID))

	if err != nil {
		return fmt.Errorf("error de red al borrar webhook en Siigo: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("error al borrar webhook en Siigo (codigo %d): %s", resp.StatusCode(), string(resp.Body()))
	}

	return nil
}
