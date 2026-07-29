package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

func (c *MeliClient) GetShipmentLabel(ctx context.Context, accessToken string, shipmentID int64, responseType string) (*domain.ShipmentLabel, error) {
	if responseType == "" {
		responseType = "pdf"
	}

	endpoint := fmt.Sprintf("%s/shipment_labels?shipment_ids=%d&response_type=%s", c.baseURL, shipmentID, responseType)

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
		return nil, domain.ErrLabelNotAvailable
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meli client: unexpected status %d: %s", resp.StatusCode, string(body))
	}
	if len(body) == 0 {
		return nil, domain.ErrLabelNotAvailable
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" || strings.HasPrefix(contentType, "application/json") {
		contentType = labelContentType(responseType)
	}

	return &domain.ShipmentLabel{
		Content:     body,
		ContentType: contentType,
		Filename:    fmt.Sprintf("meli-%d.%s", shipmentID, labelExtension(responseType)),
	}, nil
}

func labelContentType(responseType string) string {
	switch strings.ToLower(responseType) {
	case "zpl2":
		return "application/zpl"
	default:
		return "application/pdf"
	}
}

func labelExtension(responseType string) string {
	switch strings.ToLower(responseType) {
	case "zpl2":
		return "zpl"
	default:
		return "pdf"
	}
}
