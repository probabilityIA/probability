package client

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/client/mappers"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/client/response"
)

func (c *Client) CreateInvoice(ctx context.Context, req *dtos.CreateInvoiceRequest) (*dtos.CreateInvoiceResult, error) {
	result := &dtos.CreateInvoiceResult{}

	c.log.Info(ctx).
		Str("order_id", req.OrderID).
		Str("customer_dni", req.Customer.DNI).
		Msg("Creating Siigo invoice")

	token, err := c.authenticate(
		ctx,
		req.Credentials.Username,
		req.Credentials.AccessKey,
		req.Credentials.AccountID,
		req.Credentials.PartnerID,
		req.Credentials.BaseURL,
	)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to authenticate with Siigo")
		return result, fmt.Errorf("failed to authenticate with Siigo: %w", err)
	}

	if req.OrderID != "" && idempotencyEnabled(req.Config) {
		lookback, ok := intFromConfig(req.Config, "idempotency_lookback_days")
		if !ok {
			lookback = defaultIdempotencyLookbackDays
		}
		existing, ferr := c.findExistingInvoiceByOrder(ctx, req.Credentials, req.OrderID, req.Customer.DNI, lookback)
		if ferr != nil {
			c.log.Warn(ctx).Err(ferr).
				Str("order_id", req.OrderID).
				Msg("Idempotency check failed, proceeding with creation")
		} else if existing != nil {
			c.log.Info(ctx).
				Str("order_id", req.OrderID).
				Str("invoice_number", existing.Number).
				Msg("Invoice already exists for this order in Siigo, skipping duplicate creation")
			result.InvoiceNumber = existing.Number
			result.ExternalID = existing.ID
			result.IssuedAt = existing.Date
			result.AlreadyExisted = true
			result.ProviderInfo = map[string]interface{}{
				"already_existed": true,
				"siigo_id":        existing.ID,
				"invoice_name":    existing.Number,
			}
			if detail, derr := c.GetInvoiceByID(ctx, req.Credentials, existing.ID); derr == nil && detail != nil {
				result.CUFE = detail.CUFE
				result.ProviderInfo["public_url"] = detail.PublicURL
			}
			return result, nil
		}
	}

	var customerSiigoID string
	if req.Customer.DNI != "" {
		existingCustomer, err := c.GetCustomerByIdentification(ctx, req.Credentials, req.Customer.DNI)
		if err != nil {
			c.log.Warn(ctx).Err(err).Msg("Error looking up customer, will create new")
		}
		if existingCustomer != nil {
			customerSiigoID = existingCustomer.ID
			c.log.Info(ctx).
				Str("customer_id", customerSiigoID).
				Msg("Existing Siigo customer found")
		}
	}

	if customerSiigoID == "" {
		newCustomer, err := c.CreateCustomer(ctx, req.Credentials, &dtos.CreateCustomerRequest{
			Identification: req.Customer.DNI,
			Name:           req.Customer.Name,
			Email:          req.Customer.Email,
			Phone:          req.Customer.Phone,
			Address:        req.Customer.Address,
			Credentials:    req.Credentials,
		})
		if err != nil {
			c.log.Error(ctx).Err(err).Msg("Failed to create Siigo customer")
			return result, fmt.Errorf("failed to create customer in Siigo: %w", err)
		}
		if newCustomer != nil {
			customerSiigoID = newCustomer.ID
			c.log.Info(ctx).
				Str("customer_id", customerSiigoID).
				Msg("Siigo customer created")
		}
	}

	invoiceReq := mappers.BuildCreateInvoiceRequest(req, customerSiigoID)

	endpoint := c.endpointURL(req.Credentials.BaseURL, "/v1/invoices")

	c.log.Info(ctx).
		Str("endpoint", endpoint).
		Int("items_count", len(invoiceReq.Items)).
		Msg("Sending invoice to Siigo API")

	var invoiceResp response.CreateInvoiceResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Partner-Id", req.Credentials.PartnerID).
		SetBody(invoiceReq).
		SetResult(&invoiceResp).
		Post(endpoint)

	result.AuditData = &dtos.AuditData{
		RequestURL:     endpoint,
		RequestPayload: invoiceReq,
	}

	if resp != nil {
		result.AuditData.ResponseStatus = resp.StatusCode()
		result.AuditData.ResponseBody = string(resp.Body())
	}

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Siigo invoice request failed - network error")
		return result, fmt.Errorf("error de red al crear factura en Siigo: %w", err)
	}

	c.log.Info(ctx).
		Int("status_code", resp.StatusCode()).
		Str("invoice_id", invoiceResp.ID).
		Str("invoice_name", invoiceResp.Name).
		Msg("Siigo invoice response received")

	if len(invoiceResp.Errors) > 0 {
		errMsg := invoiceResp.Errors[0].Message
		c.log.Error(ctx).
			Str("error_code", invoiceResp.Errors[0].Code).
			Str("error_msg", errMsg).
			Msg("Siigo returned business error")
		return result, fmt.Errorf("Siigo rechazo la factura: %s", errMsg)
	}

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("body", string(resp.Body())).
			Msg("Siigo invoice creation failed")
		return result, fmt.Errorf("error al crear factura en Siigo (codigo %d): %s", resp.StatusCode(), string(resp.Body()))
	}

	if invoiceResp.ID == "" && invoiceResp.Name == "" {
		return result, fmt.Errorf("Siigo no retorno datos de la factura creada")
	}

	result.InvoiceNumber = invoiceResp.Name
	result.ExternalID = invoiceResp.ID
	result.CUFE = invoiceResp.Metadata.CUFE
	result.QRCode = invoiceResp.Metadata.QR
	result.ProviderInfo = map[string]interface{}{
		"siigo_id":      invoiceResp.ID,
		"invoice_name":  invoiceResp.Name,
		"invoice_total": invoiceResp.Total,
		"public_url":    invoiceResp.PublicURL,
	}

	c.log.Info(ctx).
		Str("invoice_number", result.InvoiceNumber).
		Str("external_id", result.ExternalID).
		Str("cufe", result.CUFE).
		Msg("Siigo invoice created successfully")

	return result, nil
}
