package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/infra/secondary/client/response"
)

func (c *JumpsellerClient) CreateHook(ctx context.Context, cred domain.Credential, event, hookURL string) (string, error) {
	body := response.CreateHookRequest{
		Hook: response.CreateHookFields{Event: event, URL: hookURL},
	}

	raw, err := c.do(ctx, cred, http.MethodPost, "/hooks.json", nil, body)
	if err != nil {
		return "", err
	}

	var envelope response.HookEnvelope
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return "", fmt.Errorf("jumpseller client: parsing hook: %w", err)
	}

	return strconv.FormatInt(envelope.Hook.ID, 10), nil
}

func (c *JumpsellerClient) ListHooks(ctx context.Context, cred domain.Credential) ([]domain.WebhookItem, error) {
	raw, err := c.do(ctx, cred, http.MethodGet, "/hooks.json", nil, nil)
	if err != nil {
		return nil, err
	}

	var envelopes []response.HookEnvelope
	if err := json.Unmarshal(raw, &envelopes); err != nil {
		return nil, fmt.Errorf("jumpseller client: parsing hooks: %w", err)
	}

	items := make([]domain.WebhookItem, 0, len(envelopes))
	for _, envelope := range envelopes {
		items = append(items, domain.WebhookItem{
			ID:        strconv.FormatInt(envelope.Hook.ID, 10),
			Address:   envelope.Hook.URL,
			Topic:     envelope.Hook.Event,
			Format:    "json",
			CreatedAt: envelope.Hook.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	return items, nil
}

func (c *JumpsellerClient) DeleteHook(ctx context.Context, cred domain.Credential, hookID string) error {
	_, err := c.do(ctx, cred, http.MethodDelete, fmt.Sprintf("/hooks/%s.json", hookID), nil, nil)
	return err
}
