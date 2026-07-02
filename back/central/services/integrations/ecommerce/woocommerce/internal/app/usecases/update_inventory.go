package usecases

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
)

func (uc *wooCommerceUseCase) UpdateInventory(ctx context.Context, integrationID string, productExternalID string, quantity int) error {
	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil {
		return domain.ErrIntegrationNotFound
	}

	storeURL, err := extractString(integration.Config, "store_url")
	if err != nil {
		return domain.ErrMissingStoreURL
	}

	consumerKey, err := uc.service.DecryptCredential(ctx, integrationID, "consumer_key")
	if err != nil {
		return fmt.Errorf("decrypting consumer_key: %w", err)
	}
	consumerSecret, err := uc.service.DecryptCredential(ctx, integrationID, "consumer_secret")
	if err != nil {
		return fmt.Errorf("decrypting consumer_secret: %w", err)
	}

	if err := uc.client.UpdateProductStock(ctx, storeURL, consumerKey, consumerSecret, productExternalID, quantity); err != nil {
		uc.logger.Error(ctx).
			Err(err).
			Str("integration_id", integrationID).
			Str("external_product_id", productExternalID).
			Int("quantity", quantity).
			Msg("Error al actualizar stock en WooCommerce")
		return err
	}

	uc.logger.Info(ctx).
		Str("integration_id", integrationID).
		Str("external_product_id", productExternalID).
		Int("quantity", quantity).
		Msg("Stock actualizado en WooCommerce")

	return nil
}
