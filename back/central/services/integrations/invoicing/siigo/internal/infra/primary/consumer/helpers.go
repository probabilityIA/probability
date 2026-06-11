package consumer

import (
	"context"
	"fmt"
	"time"

	siigoDtos "github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
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

	accountID, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "account_id")
	if err != nil {
		return nil, "decryption_failed", fmt.Errorf("failed to decrypt account_id")
	}

	partnerID, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "partner_id")
	if err != nil {
		return nil, "decryption_failed", fmt.Errorf("failed to decrypt partner_id")
	}

	apiURL, _ := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_url")

	effectiveURL := apiURL
	if integration.IsTesting && integration.BaseURLTest != "" {
		effectiveURL = integration.BaseURLTest
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
