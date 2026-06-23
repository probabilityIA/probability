package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
)

// RequestComparison inicia una comparación asíncrona de facturas entre el sistema y el proveedor.
// Retorna el correlationID que el frontend usará para correlacionar el evento SSE con el resultado.
func (uc *useCase) RequestComparison(ctx context.Context, dto *dtos.CompareRequestDTO) (string, error) {
	// 1. Validar rango de fechas ≤ 7 días
	dateFrom, err := time.Parse("2006-01-02", dto.DateFrom)
	if err != nil {
		return "", fmt.Errorf("invalid date_from format, expected YYYY-MM-DD: %w", err)
	}
	dateTo, err := time.Parse("2006-01-02", dto.DateTo)
	if err != nil {
		return "", fmt.Errorf("invalid date_to format, expected YYYY-MM-DD: %w", err)
	}
	if dateTo.Before(dateFrom) {
		return "", fmt.Errorf("date_to must be after date_from")
	}
	if dateTo.Sub(dateFrom) > 7*24*time.Hour {
		return "", errors.ErrCompareDateRangeTooLarge
	}

	// 2. Obtener configuración del negocio (habilitada o no — la comparación es auditoría)
	config, err := uc.repo.GetAnyConfigByBusiness(ctx, dto.BusinessID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("business_id", dto.BusinessID).Msg("Failed to get invoicing config")
		return "", errors.ErrProviderNotConfigured
	}
	if config == nil {
		uc.log.Warn(ctx).Uint("business_id", dto.BusinessID).Msg("No invoicing config found for business")
		return "", errors.ErrProviderNotConfigured
	}

	// 3. Determinar integrationID (mismo dual-read que create_invoice.go)
	var integrationID uint
	if config.InvoicingIntegrationID != nil {
		integrationID = *config.InvoicingIntegrationID
	} else if config.InvoicingProviderID != nil {
		integrationID = *config.InvoicingProviderID
	} else {
		uc.log.Error(ctx).Msg("No invoicing integration configured in active config")
		return "", errors.ErrProviderNotConfigured
	}

	// 4. Resolver proveedor dinámicamente
	provider, err := uc.resolveProvider(ctx, integrationID)
	if err != nil {
		uc.log.Warn(ctx).Err(err).Uint("integration_id", integrationID).Msg("Failed to resolve provider, defaulting to softpymes")
		provider = dtos.ProviderSoftpymes
	}

	// 5. Preparar config map para el consumer (incluye URLs + parámetros de comparación)
	compareConfig := make(map[string]interface{})
	if config.InvoiceConfig != nil {
		for k, v := range config.InvoiceConfig {
			compareConfig[k] = v
		}
	}
	// Inyectar URLs dinámicas (las necesita el consumer para resolver la URL efectiva)
	compareConfig["is_testing"] = config.IsTesting
	compareConfig["base_url"] = config.BaseURL
	compareConfig["base_url_test"] = config.BaseURLTest
	// Inyectar parámetros de comparación
	compareConfig["date_from"] = dto.DateFrom
	compareConfig["date_to"] = dto.DateTo
	compareConfig["business_id"] = dto.BusinessID

	// 6. Generar correlationID único
	correlationID := uuid.New().String()

	// 7. Construir mensaje de request (InvoiceID=0 ya que no es una factura específica)
	requestMessage := &dtos.InvoiceRequestMessage{
		InvoiceID: 0,
		Provider:  provider,
		Operation: dtos.OperationCompare,
		InvoiceData: dtos.InvoiceData{
			IntegrationID: integrationID,
			Config:        compareConfig,
		},
		CorrelationID: correlationID,
		Timestamp:     time.Now(),
	}

	// 8. Publicar request a RabbitMQ (async)
	if err := uc.invoiceRequestPub.PublishInvoiceRequest(ctx, requestMessage); err != nil {
		uc.log.Error(ctx).
			Err(err).
			Str("provider", provider).
			Str("correlation_id", correlationID).
			Msg("Failed to publish compare request to queue")
		return "", fmt.Errorf("failed to publish compare request: %w", err)
	}

	uc.log.Info(ctx).
		Uint("business_id", dto.BusinessID).
		Str("date_from", dto.DateFrom).
		Str("date_to", dto.DateTo).
		Str("provider", provider).
		Str("correlation_id", correlationID).
		Msg("📊 Invoice comparison request published")

	return correlationID, nil
}

// GetCompareResult recupera el resultado de una comparación almacenado en Redis.
// Retorna nil si no existe (aún no listo o expirado).
func (uc *useCase) GetCompareResult(ctx context.Context, correlationID string) (*dtos.CompareResponseData, error) {
	if correlationID == "" {
		return nil, fmt.Errorf("correlation_id is required")
	}

	result, err := uc.compareCache.GetCompareResult(ctx, correlationID)
	if err != nil {
		uc.log.Error(ctx).Err(err).
			Str("correlation_id", correlationID).
			Msg("Failed to get compare result from cache")
		return nil, err
	}

	return result, nil
}
