package consumer

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/queue"
)

// processListBankAccountsRequest obtiene cuentas bancarias del proveedor
// y publica un ListBankAccountsResponseMessage con las cuentas encontradas.
func (c *InvoiceRequestConsumer) processListBankAccountsRequest(
	ctx context.Context,
	request *InvoiceRequestMessage,
) error {
	// 1. Extraer parámetros del Config
	businessID := uint(0)
	if bid, ok := request.InvoiceData.Config["business_id"].(float64); ok {
		businessID = uint(bid)
	}

	c.log.Info(ctx).
		Uint("business_id", businessID).
		Str("correlation_id", request.CorrelationID).
		Msg("Starting list_bank_accounts request")

	// Helper para publicar error en el canal de list_bank_accounts
	publishErr := func(errMsg string) error {
		return c.responsePublisher.PublishListBankAccountsResponse(ctx, &queue.ListBankAccountsResponseMessage{
			Operation:     "list_bank_accounts",
			CorrelationID: request.CorrelationID,
			BusinessID:    businessID,
			Error:         errMsg,
			Timestamp:     time.Now(),
		})
	}

	// 2. Obtener integración y credenciales
	integrationID := request.InvoiceData.IntegrationID
	if integrationID == 0 {
		c.log.Error(ctx).Msg("integration_id is 0 in list_bank_accounts request")
		return publishErr("integration_id is 0")
	}

	integrationIDStr := fmt.Sprintf("%d", integrationID)
	integration, err := c.integrationCore.GetIntegrationByID(ctx, integrationIDStr)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to get integration for list_bank_accounts")
		return publishErr("failed to get integration: " + err.Error())
	}

	apiKey, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_key")
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to decrypt api_key")
		return publishErr("failed to decrypt api_key")
	}

	apiSecret, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_secret")
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to decrypt api_secret")
		return publishErr("failed to decrypt api_secret")
	}

	// 3. Combinar config de integración con config del mensaje
	combinedConfig := make(map[string]interface{})
	for k, v := range integration.Config {
		combinedConfig[k] = v
	}
	for k, v := range request.InvoiceData.Config {
		combinedConfig[k] = v
	}

	referer, _ := combinedConfig["referer"].(string)

	// 4. Resolver URL efectiva desde integration_type
	effectiveURL := integration.BaseURL
	if integration.IsTesting && integration.BaseURLTest != "" {
		effectiveURL = integration.BaseURLTest
	}
	if effectiveURL == "" {
		c.log.Error(ctx).
			Uint("integration_id", integrationID).
			Msg("base_url no configurada en el tipo de integración Softpymes")
		return publishErr("base_url no configurada en el tipo de integración Softpymes")
	}

	c.log.Info(ctx).
		Bool("is_testing", integration.IsTesting).
		Str("effective_url", effectiveURL).
		Msg("Resolved effective Softpymes URL for list_bank_accounts")

	// 5. Obtener cuentas bancarias (sin paginación, lista pequeña)
	accounts, err := c.softpymesClient.ListBankAccounts(ctx, apiKey, apiSecret, referer, effectiveURL, "001")
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to list bank accounts")
		return publishErr("failed to list bank accounts: " + err.Error())
	}

	// 6. Mapear a response items
	items := make([]queue.BankAccountItem, 0, len(accounts))
	for _, acc := range accounts {
		items = append(items, queue.BankAccountItem{
			AccountNumber: acc.AccountNumber,
			Name:          acc.Name,
			NameType:      acc.NameType,
		})
	}

	c.log.Info(ctx).
		Int("total_accounts", len(items)).
		Str("correlation_id", request.CorrelationID).
		Msg("Bank accounts fetched, publishing list_bank_accounts response")

	// 7. Publicar resultado
	return c.responsePublisher.PublishListBankAccountsResponse(ctx, &queue.ListBankAccountsResponseMessage{
		Operation:     "list_bank_accounts",
		CorrelationID: request.CorrelationID,
		BusinessID:    businessID,
		Items:         items,
		Timestamp:     time.Now(),
	})
}
