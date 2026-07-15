package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

type claimResponse struct {
	ID         int64  `json:"id"`
	Resource   string `json:"resource"`
	ResourceID int64  `json:"resource_id"`
	Status     string `json:"status"`
	Reason     string `json:"reason_id"`
}

type claimMessagesResponse struct {
	Messages []struct {
		Message struct {
			Text string `json:"text"`
		} `json:"message"`
	} `json:"messages"`
}

func (c *MeliClient) GetClaim(ctx context.Context, accessToken string, claimID int64) (*domain.MeliClaim, error) {
	endpoint := fmt.Sprintf("%s/v1/claims/%d", c.baseURL, claimID)

	resp, body, err := c.do(ctx, func() (*http.Request, error) {
		return c.newAuthorizedRequest(ctx, http.MethodGet, endpoint, accessToken)
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, domain.ErrTokenExpired
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, domain.ErrOrderNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meli client: claim status %d: %s", resp.StatusCode, string(body))
	}

	var parsed claimResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("meli client: parsing claim: %w", err)
	}

	claim := &domain.MeliClaim{
		ID:           parsed.ID,
		ResourceType: parsed.Resource,
		ResourceID:   parsed.ResourceID,
		Reason:       parsed.Reason,
		Status:       parsed.Status,
	}

	claim.Messages = c.fetchClaimMessages(ctx, accessToken, claimID)
	return claim, nil
}

func (c *MeliClient) fetchClaimMessages(ctx context.Context, accessToken string, claimID int64) []string {
	endpoint := fmt.Sprintf("%s/v1/claims/%d/messages", c.baseURL, claimID)
	resp, body, err := c.do(ctx, func() (*http.Request, error) {
		return c.newAuthorizedRequest(ctx, http.MethodGet, endpoint, accessToken)
	})
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil
	}
	var parsed claimMessagesResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil
	}
	var texts []string
	for _, m := range parsed.Messages {
		if m.Message.Text != "" {
			texts = append(texts, m.Message.Text)
		}
	}
	return texts
}
