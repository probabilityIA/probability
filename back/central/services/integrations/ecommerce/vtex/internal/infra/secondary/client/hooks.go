package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
)

type hookConfigPayload struct {
	Filter struct {
		Type   string   `json:"type"`
		Status []string `json:"status"`
	} `json:"filter"`
	Hook struct {
		URL     string            `json:"url"`
		Headers map[string]string `json:"headers"`
	} `json:"hook"`
}

func (c *VTEXClient) GetOrderHook(ctx context.Context, cred domain.Credential) (*domain.HookConfig, error) {
	endpoint := fmt.Sprintf("%s/api/orders/hook/config", baseURL(cred))

	body, err := c.do(ctx, http.MethodGet, endpoint, cred, nil)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return nil, nil
		}
		return nil, err
	}

	trimmed := strings.TrimSpace(string(body))
	if trimmed == "" || trimmed == "null" {
		return nil, nil
	}

	var payload hookConfigPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("vtex client: parsing hook config: %w", err)
	}

	if strings.TrimSpace(payload.Hook.URL) == "" {
		return nil, nil
	}

	return &domain.HookConfig{
		URL:      payload.Hook.URL,
		Statuses: payload.Filter.Status,
		HasKey:   payload.Hook.Headers[domain.HookKeyHeader] != "",
	}, nil
}

func (c *VTEXClient) SetOrderHook(ctx context.Context, cred domain.Credential, url, hookKey string) error {
	endpoint := fmt.Sprintf("%s/api/orders/hook/config", baseURL(cred))

	payload := hookConfigPayload{}
	payload.Filter.Type = domain.HookFilterType
	payload.Filter.Status = domain.WebhookOrderStates
	payload.Hook.URL = url
	payload.Hook.Headers = map[string]string{
		domain.HookKeyHeader: hookKey,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("vtex client: building hook config: %w", err)
	}

	_, err = c.do(ctx, http.MethodPost, endpoint, cred, body)
	return err
}

func (c *VTEXClient) DeleteOrderHook(ctx context.Context, cred domain.Credential) error {
	endpoint := fmt.Sprintf("%s/api/orders/hook/config", baseURL(cred))

	_, err := c.do(ctx, http.MethodDelete, endpoint, cred, nil)
	if err != nil && errors.Is(err, domain.ErrProductNotFound) {
		return nil
	}
	return err
}
