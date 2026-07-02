package core

import (
	"context"

	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/ports"
)

// SiigoCore adapta el use case de Siigo al contrato core.IIntegrationContract.
type SiigoCore struct {
	integrationcore.BaseIntegration
	useCase ports.IInvoiceUseCase
}

func New(useCase ports.IInvoiceUseCase) *SiigoCore {
	return &SiigoCore{useCase: useCase}
}

func (s *SiigoCore) TestConnection(ctx context.Context, config, credentials map[string]interface{}) error {
	return s.useCase.TestConnection(ctx, config, credentials)
}

func (s *SiigoCore) ListWebhooks(ctx context.Context, integrationID string) ([]interface{}, error) {
	webhooks, err := s.useCase.ListWebhooks(ctx, integrationID)
	if err != nil {
		return nil, err
	}
	result := make([]interface{}, 0, len(webhooks))
	for _, w := range webhooks {
		result = append(result, w)
	}
	return result, nil
}

func (s *SiigoCore) DeleteWebhook(ctx context.Context, integrationID, webhookID string) error {
	return s.useCase.DeleteWebhook(ctx, integrationID, webhookID)
}

func (s *SiigoCore) VerifyWebhooksByURL(ctx context.Context, integrationID string, baseURL string) ([]interface{}, error) {
	webhooks, err := s.useCase.VerifyWebhooksByURL(ctx, integrationID, baseURL)
	if err != nil {
		return nil, err
	}
	result := make([]interface{}, 0, len(webhooks))
	for _, w := range webhooks {
		result = append(result, w)
	}
	return result, nil
}

func (s *SiigoCore) CreateWebhook(ctx context.Context, integrationID string, baseURL string) (interface{}, error) {
	return s.useCase.CreateWebhooks(ctx, integrationID, baseURL)
}
