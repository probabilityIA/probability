package usecases

import (
	"context"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
)

func (uc *jumpsellerUseCase) UpdateInventory(ctx context.Context, integrationID string, productExternalID string, quantity int) error {
	integration, cred, err := uc.resolveIntegration(ctx, integrationID)
	if err != nil {
		return err
	}

	if enabled, _ := integration.Config["inventory_sync_enabled"].(bool); !enabled {
		uc.logger.Info(ctx).
			Str("integration_id", integrationID).
			Msg("Sync de inventario desactivado para la integracion Jumpseller, push omitido")
		return nil
	}

	productID, variantID, err := parseExternalProductID(productExternalID)
	if err != nil {
		return err
	}

	if variantID > 0 {
		err = uc.client.SetVariantStock(ctx, cred, productID, variantID, quantity)
	} else {
		err = uc.client.SetProductStock(ctx, cred, productID, quantity)
	}
	if err != nil {
		uc.logger.Error(ctx).Err(err).
			Str("integration_id", integrationID).
			Str("external_product_id", productExternalID).
			Int("quantity", quantity).
			Msg("Error al actualizar stock en Jumpseller")
		return err
	}

	uc.logger.Info(ctx).
		Str("integration_id", integrationID).
		Str("external_product_id", productExternalID).
		Int("quantity", quantity).
		Msg("Stock actualizado en Jumpseller")

	return nil
}

func parseExternalProductID(externalID string) (productID int64, variantID int64, err error) {
	productPart, variantPart, hasVariant := strings.Cut(externalID, ":")

	productID, err = strconv.ParseInt(productPart, 10, 64)
	if err != nil {
		return 0, 0, domain.ErrProductNotFound
	}

	if hasVariant && variantPart != "" {
		variantID, err = strconv.ParseInt(variantPart, 10, 64)
		if err != nil {
			return 0, 0, domain.ErrProductNotFound
		}
	}

	return productID, variantID, nil
}
