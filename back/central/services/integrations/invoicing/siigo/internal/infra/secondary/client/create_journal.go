package client

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/client/mappers"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/client/response"
)

// CreateJournal crea un comprobante contable en Siigo
// Endpoint: POST /v1/journals
func (c *Client) CreateJournal(ctx context.Context, req *dtos.CreateJournalRequest) (*dtos.CreateJournalResult, error) {
	result := &dtos.CreateJournalResult{}

	c.log.Info(ctx).
		Str("order_id", req.OrderID).
		Int("items_count", len(req.Items)).
		Msg("📋 Creating Siigo journal entry")

	// 1. Autenticar
	token, err := c.authenticate(
		ctx,
		req.Credentials.Username,
		req.Credentials.AccessKey,
		req.Credentials.AccountID,
		req.Credentials.PartnerID,
		req.Credentials.BaseURL,
	)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("❌ Failed to authenticate with Siigo")
		return result, fmt.Errorf("failed to authenticate with Siigo: %w", err)
	}

	// 2. Construir request del journal
	journalReq := mappers.BuildCreateJournalRequest(req)

	// Endpoint
	endpoint := c.endpointURL(req.Credentials.BaseURL, "/v1/journals")

	c.log.Info(ctx).
		Str("endpoint", endpoint).
		Int("items_count", len(journalReq.Items)).
		Msg("🚀 Sending journal to Siigo API")

	// 3. Llamar a la API de Siigo
	var journalResp response.CreateJournalResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Partner-Id", req.Credentials.PartnerID).
		SetBody(journalReq).
		SetResult(&journalResp).
		Post(endpoint)

	// Capturar audit data siempre (incluso en error)
	result.AuditData = &dtos.AuditData{
		RequestURL:     endpoint,
		RequestPayload: journalReq,
	}

	if resp != nil {
		result.AuditData.ResponseStatus = resp.StatusCode()
		result.AuditData.ResponseBody = string(resp.Body())
	}

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("❌ Siigo journal request failed - network error")
		return result, fmt.Errorf("error de red al crear journal en Siigo: %w", err)
	}

	c.log.Info(ctx).
		Int("status_code", resp.StatusCode()).
		Str("journal_id", journalResp.ID).
		Str("journal_name", journalResp.Name).
		Msg("📥 Siigo journal response received")

	// 4. Verificar errores de negocio en la respuesta
	if len(journalResp.Errors) > 0 {
		errMsg := journalResp.Errors[0].Message
		c.log.Error(ctx).
			Str("error_code", journalResp.Errors[0].Code).
			Str("error_msg", errMsg).
			Msg("❌ Siigo returned business error")
		return result, fmt.Errorf("Siigo rechazó el journal: %s", errMsg)
	}

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("body", string(resp.Body())).
			Msg("❌ Siigo journal creation failed")
		return result, fmt.Errorf("error al crear journal en Siigo (código %d): %s", resp.StatusCode(), string(resp.Body()))
	}

	if journalResp.ID == "" && journalResp.Name == "" {
		return result, fmt.Errorf("Siigo no retornó datos del journal creado")
	}

	// 5. Poblar resultado exitoso
	result.JournalName = journalResp.Name
	result.JournalID = journalResp.ID
	result.Number = journalResp.Number
	result.Date = journalResp.Date
	result.ProviderInfo = map[string]interface{}{
		"siigo_id":      journalResp.ID,
		"journal_name":  journalResp.Name,
		"journal_total": journalResp.Total,
	}

	c.log.Info(ctx).
		Str("journal_name", result.JournalName).
		Str("journal_id", result.JournalID).
		Msg("✅ Siigo journal created successfully")

	return result, nil
}
