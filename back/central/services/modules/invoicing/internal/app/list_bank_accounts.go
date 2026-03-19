package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
)

// RequestListBankAccounts inicia una solicitud asíncrona de cuentas bancarias del proveedor.
// Retorna el correlationID que el frontend usará para correlacionar el evento SSE con el resultado.
func (uc *useCase) RequestListBankAccounts(ctx context.Context, dto *dtos.ListBankAccountsRequestDTO) (string, error) {
	// 1. Obtener configuración del negocio
	config, err := uc.repo.GetAnyConfigByBusiness(ctx, dto.BusinessID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("business_id", dto.BusinessID).Msg("Failed to get invoicing config")
		return "", errors.ErrProviderNotConfigured
	}
	if config == nil {
		uc.log.Warn(ctx).Uint("business_id", dto.BusinessID).Msg("No invoicing config found for business")
		return "", errors.ErrProviderNotConfigured
	}

	// 2. Determinar integrationID (mismo dual-read que create_invoice.go)
	var integrationID uint
	if config.InvoicingIntegrationID != nil {
		integrationID = *config.InvoicingIntegrationID
	} else if config.InvoicingProviderID != nil {
		integrationID = *config.InvoicingProviderID
	} else {
		uc.log.Error(ctx).Msg("No invoicing integration configured in active config")
		return "", errors.ErrProviderNotConfigured
	}

	// 3. Resolver proveedor dinámicamente
	provider, err := uc.resolveProvider(ctx, integrationID)
	if err != nil {
		uc.log.Warn(ctx).Err(err).Uint("integration_id", integrationID).Msg("Failed to resolve provider, defaulting to softpymes")
		provider = dtos.ProviderSoftpymes
	}

	// 4. Preparar config map para el consumer
	bankAccountsConfig := make(map[string]interface{})
	if config.InvoiceConfig != nil {
		for k, v := range config.InvoiceConfig {
			bankAccountsConfig[k] = v
		}
	}
	bankAccountsConfig["is_testing"] = config.IsTesting
	bankAccountsConfig["base_url"] = config.BaseURL
	bankAccountsConfig["base_url_test"] = config.BaseURLTest
	bankAccountsConfig["business_id"] = dto.BusinessID

	// 5. Generar correlationID único
	correlationID := uuid.New().String()

	// 6. Construir mensaje de request
	requestMessage := &dtos.InvoiceRequestMessage{
		InvoiceID: 0,
		Provider:  provider,
		Operation: dtos.OperationListBankAccounts,
		InvoiceData: dtos.InvoiceData{
			IntegrationID: integrationID,
			Config:        bankAccountsConfig,
		},
		CorrelationID: correlationID,
		Timestamp:     time.Now(),
	}

	// 7. Publicar request a RabbitMQ (async)
	if err := uc.invoiceRequestPub.PublishInvoiceRequest(ctx, requestMessage); err != nil {
		uc.log.Error(ctx).
			Err(err).
			Str("provider", provider).
			Str("correlation_id", correlationID).
			Msg("Failed to publish list bank accounts request to queue")
		return "", fmt.Errorf("failed to publish list bank accounts request: %w", err)
	}

	uc.log.Info(ctx).
		Uint("business_id", dto.BusinessID).
		Str("provider", provider).
		Str("correlation_id", correlationID).
		Msg("🏦 List bank accounts request published")

	return correlationID, nil
}

// GetListBankAccountsResult recupera el resultado de cuentas bancarias almacenado en Redis.
// Retorna nil si no existe (aún no listo o expirado).
func (uc *useCase) GetListBankAccountsResult(ctx context.Context, correlationID string) (*dtos.BankAccountsResponseData, error) {
	if correlationID == "" {
		return nil, fmt.Errorf("correlation_id is required")
	}

	result, err := uc.compareCache.GetBankAccountsResult(ctx, correlationID)
	if err != nil {
		uc.log.Error(ctx).Err(err).
			Str("correlation_id", correlationID).
			Msg("Failed to get bank accounts result from cache")
		return nil, err
	}

	return result, nil
}
