package consumer

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	siigoDtos "github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/entities"
	siigoerrors "github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/queue"
)

type integrationContext struct {
	Credentials siigoDtos.Credentials
	Config      map[string]interface{}
	IsTesting   bool
}

func (c *InvoiceRequestConsumer) resolveIntegration(
	ctx context.Context,
	request *InvoiceRequestMessage,
) (*integrationContext, string, error) {
	integrationID := request.InvoiceData.IntegrationID
	if integrationID == 0 {
		return nil, "missing_integration_id", fmt.Errorf("integration_id is 0")
	}

	integrationIDStr := fmt.Sprintf("%d", integrationID)
	integration, err := c.integrationCore.GetIntegrationByID(ctx, integrationIDStr)
	if err != nil {
		return nil, "integration_not_found", err
	}

	username, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "username")
	if err != nil {
		return nil, "decryption_failed", fmt.Errorf("failed to decrypt username")
	}

	accessKey, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "access_key")
	if err != nil {
		return nil, "decryption_failed", fmt.Errorf("failed to decrypt access_key")
	}

	accountID, _ := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "account_id")

	partnerID, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "partner_id")
	if err != nil {
		return nil, "decryption_failed", fmt.Errorf("failed to decrypt partner_id")
	}

	apiURL, _ := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_url")

	effectiveURL := entities.ResolveSiigoBaseURL(integration.IsTesting, integration.BaseURLTest, apiURL, integration.BaseURL)
	if effectiveURL == "" {
		return nil, "missing_base_url", fmt.Errorf("URL de Siigo no configurada en el tipo de integracion (base_url o base_url_test)")
	}

	combinedConfig := make(map[string]interface{})
	for k, v := range integration.Config {
		combinedConfig[k] = v
	}
	for k, v := range request.InvoiceData.Config {
		combinedConfig[k] = v
	}

	return &integrationContext{
		Credentials: siigoDtos.Credentials{
			Username:  username,
			AccessKey: accessKey,
			AccountID: accountID,
			PartnerID: partnerID,
			BaseURL:   effectiveURL,
		},
		Config:    combinedConfig,
		IsTesting: integration.IsTesting,
	}, "", nil
}

func (c *InvoiceRequestConsumer) createOperationErrorResponse(
	request *InvoiceRequestMessage,
	operation string,
	errorCode string,
	errorMsg string,
	startTime time.Time,
	auditData *siigoDtos.AuditData,
) *queue.InvoiceResponseMessage {
	resp := c.createErrorResponse(request, errorCode, errorMsg, startTime, auditData)
	resp.Operation = operation
	return resp
}

func resultAudit(result *siigoDtos.AnnulInvoiceResult) *siigoDtos.AuditData {
	if result == nil {
		return nil
	}
	return result.AuditData
}

func businessIDFromConfig(config map[string]interface{}) uint {
	if bid, ok := config["business_id"].(float64); ok {
		return uint(bid)
	}
	return 0
}

func defaultCustomerDNI(config map[string]interface{}) string {
	if nit, ok := config["default_customer_nit"].(string); ok {
		if trimmed := strings.TrimSpace(nit); trimmed != "" {
			return trimmed
		}
	}
	return "222222222222"
}

func camposDeConfiguracionFaltantes(config map[string]interface{}) []string {
	requeridos := []struct {
		clave    string
		etiqueta string
	}{
		{"document_id", "Tipo de documento (FV)"},
		{"seller_id", "Vendedor"},
		{"payment_method_id", "Medio de pago"},
	}

	faltantes := make([]string, 0, len(requeridos))
	for _, r := range requeridos {
		if getIntFromConfigMap(config, r.clave) <= 0 {
			faltantes = append(faltantes, r.etiqueta)
		}
	}
	return faltantes
}

func getIntFromConfigMap(config map[string]interface{}, key string) int {
	switch v := config[key].(type) {
	case float64:
		return int(v)
	case int:
		return v
	case int64:
		return int(v)
	case string:
		n, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return 0
		}
		return n
	}
	return 0
}

func codigoDeError(err error) string {
	return siigoerrors.CodeFromError(err)
}
