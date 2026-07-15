package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

type billingInfoResponse struct {
	Buyer struct {
		BillingInfo struct {
			DocType        string `json:"doc_type"`
			DocNumber      string `json:"doc_number"`
			AdditionalInfo []struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"additional_info"`
		} `json:"billing_info"`
	} `json:"buyer"`
}

func (c *MeliClient) GetBillingInfo(ctx context.Context, accessToken string, orderID int64) (*domain.MeliBillingInfo, error) {
	endpoint := fmt.Sprintf("%s/orders/%d/billing_info", c.baseURL, orderID)

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
		return nil, domain.ErrBillingInfoNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meli client: billing_info status %d: %s", resp.StatusCode, string(body))
	}

	var parsed billingInfoResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("meli client: parsing billing_info: %w", err)
	}

	info := &domain.MeliBillingInfo{
		DocType:   parsed.Buyer.BillingInfo.DocType,
		DocNumber: parsed.Buyer.BillingInfo.DocNumber,
	}

	if info.DocNumber == "" {
		for _, ai := range parsed.Buyer.BillingInfo.AdditionalInfo {
			switch ai.Type {
			case "DOC_TYPE":
				info.DocType = ai.Value
			case "DOC_NUMBER":
				info.DocNumber = ai.Value
			}
		}
	}

	if info.DocNumber == "" {
		return nil, domain.ErrBillingInfoNotFound
	}

	return info, nil
}
