package client

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
)

const (
	defaultIdempotencyLookbackDays = 30
	maxIdempotencyPages            = 10
	idempotencyPageSize            = 100
)

type existingInvoice struct {
	ID     string
	Number string
	Date   string
}

func idempotencyEnabled(config map[string]interface{}) bool {
	return boolFromConfig(config, "idempotency_check", true)
}

func (c *Client) findExistingInvoiceByOrder(ctx context.Context, credentials dtos.Credentials, orderID, customerDNI string, lookbackDays int) (*existingInvoice, error) {
	if orderID == "" {
		return nil, nil
	}

	token, err := c.authenticate(ctx, credentials.Username, credentials.AccessKey, credentials.AccountID, credentials.PartnerID, credentials.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with Siigo: %w", err)
	}

	if lookbackDays <= 0 {
		lookbackDays = defaultIdempotencyLookbackDays
	}

	loc, locErr := time.LoadLocation("America/Bogota")
	if locErr != nil {
		loc = time.UTC
	}
	now := time.Now().In(loc)
	createdStart := now.AddDate(0, 0, -lookbackDays).Format(time.RFC3339)
	createdEnd := now.AddDate(0, 0, 1).Format(time.RFC3339)

	marker := "order:" + orderID

	c.log.Info(ctx).
		Str("order_id", orderID).
		Str("created_start", createdStart).
		Str("customer_identification", customerDNI).
		Msg("Checking for existing Siigo invoice (idempotency)")

	for page := 1; page <= maxIdempotencyPages; page++ {
		req := c.httpClient.R().
			SetContext(ctx).
			SetAuthToken(token).
			SetHeader("Partner-Id", credentials.PartnerID).
			SetQueryParam("created_start", createdStart).
			SetQueryParam("created_end", createdEnd).
			SetQueryParam("page", strconv.Itoa(page)).
			SetQueryParam("page_size", strconv.Itoa(idempotencyPageSize))

		if customerDNI != "" {
			req = req.SetQueryParam("customer_identification", customerDNI)
		}

		var listResp listInvoicesResponse
		resp, err := req.SetResult(&listResp).Get(c.endpointURL(credentials.BaseURL, "/v1/invoices"))
		if err != nil {
			return nil, fmt.Errorf("error de red al buscar factura existente en Siigo: %w", err)
		}
		if resp.IsError() {
			return nil, fmt.Errorf("error al buscar factura existente en Siigo (codigo %d)", resp.StatusCode())
		}

		for _, r := range listResp.Results {
			if !strings.Contains(r.Observations, marker) {
				continue
			}
			if isAnnulledStatus(r.Status) {
				continue
			}
			c.log.Info(ctx).
				Str("order_id", orderID).
				Str("siigo_id", r.ID).
				Str("invoice_number", r.Name).
				Msg("Found existing vigente Siigo invoice for order - skipping duplicate")
			return &existingInvoice{ID: r.ID, Number: r.Name, Date: r.Date}, nil
		}

		if len(listResp.Results) < idempotencyPageSize {
			break
		}
		if page == maxIdempotencyPages {
			c.log.Warn(ctx).
				Str("order_id", orderID).
				Int("max_pages", maxIdempotencyPages).
				Msg("Idempotency search hit page cap; older invoices not scanned")
		}
	}

	return nil, nil
}

func boolFromConfig(config map[string]interface{}, key string, defaultVal bool) bool {
	if v, ok := config[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return defaultVal
}
