package tester

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/domain"
)

type shopifyTester struct {
	client domain.ShopifyClient
}

func New(client domain.ShopifyClient) core.ITestIntegration {
	return &shopifyTester{
		client: client,
	}
}

func (t *shopifyTester) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	storeName, ok := config["store_name"].(string)
	if !ok || storeName == "" {
		return fmt.Errorf("store_name is required in config")
	}

	accessToken, ok := credentials["access_token"].(string)
	if !ok || accessToken == "" {
		return fmt.Errorf("access_token is required in credentials")
	}

	valid, _, err := t.client.ValidateToken(ctx, storeName, accessToken)
	if err != nil {
		return fmt.Errorf("failed to validate token: %w", err)
	}

	if !valid {
		return fmt.Errorf("invalid credentials or store name")
	}

	return nil
}
