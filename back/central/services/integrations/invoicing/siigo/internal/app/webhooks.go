package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/ports"
)

var siigoWebhookTopics = []string{
	"public.siigoapi.products.stock.update",
	"public.siigoapi.products.create",
	"public.siigoapi.products.update",
}

func (uc *invoicingUseCase) resolveWebhookCredentials(ctx context.Context, integrationID string) (dtos.Credentials, error) {
	integration, err := uc.integrationCore.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return dtos.Credentials{}, fmt.Errorf("error al obtener integracion: %w", err)
	}

	username, err := uc.integrationCore.DecryptCredential(ctx, integrationID, "username")
	if err != nil {
		return dtos.Credentials{}, fmt.Errorf("error al descifrar username: %w", err)
	}
	accessKey, err := uc.integrationCore.DecryptCredential(ctx, integrationID, "access_key")
	if err != nil {
		return dtos.Credentials{}, fmt.Errorf("error al descifrar access_key: %w", err)
	}
	accountID, _ := uc.integrationCore.DecryptCredential(ctx, integrationID, "account_id")
	partnerID, err := uc.integrationCore.DecryptCredential(ctx, integrationID, "partner_id")
	if err != nil {
		return dtos.Credentials{}, fmt.Errorf("error al descifrar partner_id: %w", err)
	}
	apiURL, _ := uc.integrationCore.DecryptCredential(ctx, integrationID, "api_url")

	effectiveURL := entities.ResolveSiigoBaseURL(integration.IsTesting, integration.BaseURLTest, apiURL, integration.BaseURL)
	if effectiveURL == "" {
		return dtos.Credentials{}, fmt.Errorf("URL de Siigo no configurada en el tipo de integracion (base_url o base_url_test)")
	}

	return dtos.Credentials{
		Username:  username,
		AccessKey: accessKey,
		AccountID: accountID,
		PartnerID: partnerID,
		BaseURL:   effectiveURL,
	}, nil
}

func (uc *invoicingUseCase) ListWebhooks(ctx context.Context, integrationID string) ([]dtos.WebhookItem, error) {
	credentials, err := uc.resolveWebhookCredentials(ctx, integrationID)
	if err != nil {
		return nil, err
	}
	return uc.siigoClient.ListWebhooks(ctx, credentials)
}

func (uc *invoicingUseCase) DeleteWebhook(ctx context.Context, integrationID string, webhookID string) error {
	credentials, err := uc.resolveWebhookCredentials(ctx, integrationID)
	if err != nil {
		return err
	}
	return uc.siigoClient.DeleteWebhook(ctx, credentials, webhookID)
}

func siigoWebhookURL(baseURL, integrationID string) string {
	baseURL = strings.TrimSuffix(strings.TrimSpace(baseURL), "/")
	baseURL = strings.TrimSuffix(baseURL, "/api/v1")
	return fmt.Sprintf("%s/api/v1/siigo/webhook?integration_id=%s", baseURL, integrationID)
}

func (uc *invoicingUseCase) VerifyWebhooksByURL(ctx context.Context, integrationID string, baseURL string) ([]dtos.WebhookItem, error) {
	credentials, err := uc.resolveWebhookCredentials(ctx, integrationID)
	if err != nil {
		return nil, err
	}
	all, err := uc.siigoClient.ListWebhooks(ctx, credentials)
	if err != nil {
		return nil, err
	}
	target := siigoWebhookURL(baseURL, integrationID)
	matched := make([]dtos.WebhookItem, 0, len(all))
	for _, w := range all {
		if w.URL == target {
			matched = append(matched, w)
		}
	}
	return matched, nil
}

func (uc *invoicingUseCase) CreateWebhooks(ctx context.Context, integrationID string, baseURL string) (*ports.WebhookCreateResult, error) {
	credentials, err := uc.resolveWebhookCredentials(ctx, integrationID)
	if err != nil {
		return nil, err
	}

	webhookURL := siigoWebhookURL(baseURL, integrationID)

	applicationID := credentials.PartnerID
	if applicationID == "" {
		applicationID = "Probability"
	}

	existing, err := uc.siigoClient.ListWebhooks(ctx, credentials)
	if err != nil {
		return nil, err
	}

	existingByTopic := make(map[string]dtos.WebhookItem)
	for _, w := range existing {
		if w.URL == webhookURL {
			existingByTopic[w.Topic] = w
		}
	}

	result := &ports.WebhookCreateResult{
		WebhookURL:       webhookURL,
		CreatedWebhooks:  []dtos.WebhookItem{},
		ExistingWebhooks: []dtos.WebhookItem{},
		Errors:           []string{},
	}

	for _, topic := range siigoWebhookTopics {
		if w, ok := existingByTopic[topic]; ok {
			result.ExistingWebhooks = append(result.ExistingWebhooks, w)
			continue
		}
		created, err := uc.siigoClient.CreateWebhook(ctx, credentials, dtos.CreateWebhookInput{
			ApplicationID: applicationID,
			URL:           webhookURL,
			Topic:         topic,
		})
		if err != nil {
			uc.log.Error(ctx).Err(err).Str("topic", topic).Msg("Error al crear webhook Siigo")
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", topic, err))
			continue
		}
		result.CreatedWebhooks = append(result.CreatedWebhooks, *created)
	}

	return result, nil
}
