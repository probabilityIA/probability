package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/domain/dtos"
)

// CreateInvoice procesa una solicitud de facturación con Factus.
// Es el único responsable de la lógica de negocio:
//  1. Obtiene la integración desde el core
//  2. Descifra las credenciales de la base de datos
//  3. Combina la configuración
//  4. Construye el request tipado para el cliente HTTP
//  5. Llama al cliente HTTP de Factus
//
// El consumer (adapter primario) solo deserializa el mensaje y delega aquí.
// El cliente HTTP (adapter secundario) solo ejecuta la llamada HTTP.
func (uc *invoicingUseCase) CreateInvoice(ctx context.Context, req *dtos.ProcessInvoiceRequest) (*dtos.ProcessInvoiceResult, error) {
	integrationIDStr := fmt.Sprintf("%d", req.IntegrationID)

	uc.log.Info(ctx).
		Uint("invoice_id", req.InvoiceID).
		Str("order_id", req.OrderID).
		Str("integration_id", integrationIDStr).
		Msg("Processing Factus invoice request")

	// 1. Obtener integración desde IntegrationCore
	integration, err := uc.integrationCore.GetIntegrationByID(ctx, integrationIDStr)
	if err != nil {
		uc.log.Error(ctx).Err(err).Str("integration_id", integrationIDStr).Msg("Failed to get integration")
		return nil, fmt.Errorf("factus: integration not found (%s): %w", integrationIDStr, err)
	}

	// 2. Descifrar credenciales desde la base de datos
	clientID, err := uc.integrationCore.DecryptCredential(ctx, integrationIDStr, "client_id")
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to decrypt client_id")
		return nil, fmt.Errorf("factus: failed to decrypt client_id: %w", err)
	}

	clientSecret, err := uc.integrationCore.DecryptCredential(ctx, integrationIDStr, "client_secret")
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to decrypt client_secret")
		return nil, fmt.Errorf("factus: failed to decrypt client_secret: %w", err)
	}

	username, err := uc.integrationCore.DecryptCredential(ctx, integrationIDStr, "username")
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to decrypt username")
		return nil, fmt.Errorf("factus: failed to decrypt username: %w", err)
	}

	password, err := uc.integrationCore.DecryptCredential(ctx, integrationIDStr, "password")
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to decrypt password")
		return nil, fmt.Errorf("factus: failed to decrypt password: %w", err)
	}

	// api_url viene de las credenciales en la DB (campo opcional)
	apiURL, _ := uc.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_url")

	// 3. Combinar config: primero la de la integración, luego la del request (prioridad)
	combinedConfig := make(map[string]interface{})
	for k, v := range integration.Config {
		combinedConfig[k] = v
	}
	for k, v := range req.Config {
		combinedConfig[k] = v
	}

	// 4. Construir request tipado para el cliente HTTP
	clientReq := &dtos.CreateInvoiceRequest{
		Customer:     req.Customer,
		Items:        req.Items,
		Total:        req.Total,
		Subtotal:     req.Subtotal,
		Tax:          req.Tax,
		Discount:     req.Discount,
		ShippingCost: req.ShippingCost,
		Currency:     req.Currency,
		OrderID:      req.OrderID,
		Credentials: dtos.Credentials{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Username:     username,
			Password:     password,
			BaseURL:      apiURL,
		},
		Config: combinedConfig,
	}

	// 5. Llamar al cliente HTTP de Factus (adapter secundario)
	uc.log.Info(ctx).
		Uint("invoice_id", req.InvoiceID).
		Str("order_id", req.OrderID).
		Msg("Calling Factus HTTP client")

	result, err := uc.factusClient.CreateInvoice(ctx, clientReq)
	if err != nil {
		uc.log.Error(ctx).Err(err).
			Uint("invoice_id", req.InvoiceID).
			Msg("Factus HTTP client returned error")

		ucResult := &dtos.ProcessInvoiceResult{}
		if result != nil {
			ucResult.AuditData = result.AuditData
		}
		return ucResult, fmt.Errorf("factus: API call failed: %w", err)
	}

	uc.log.Info(ctx).
		Uint("invoice_id", req.InvoiceID).
		Str("invoice_number", result.InvoiceNumber).
		Str("cufe", result.CUFE).
		Msg("Factus invoice created successfully")

	return &dtos.ProcessInvoiceResult{
		InvoiceNumber: result.InvoiceNumber,
		ExternalID:    result.ExternalID,
		CUFE:          result.CUFE,
		QRCode:        result.QRCode,
		Total:         result.Total,
		IssuedAt:      result.IssuedAt,
		AuditData:     result.AuditData,
	}, nil
}
