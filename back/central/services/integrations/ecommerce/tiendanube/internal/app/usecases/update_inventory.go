package usecases

import (
	"context"
	"fmt"
	"strconv"
)

func (uc *tiendanubeUseCase) UpdateInventory(ctx context.Context, integrationID string, productExternalID string, quantity int) error {
	integration, err := uc.fetchIntegration(ctx, integrationID)
	if err != nil {
		return err
	}

	businessID := uint(0)
	if integration.BusinessID != nil {
		businessID = *integration.BusinessID
	}

	_, cred, err := uc.resolveIntegrationForBusiness(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	productID, variantID, perr := parseExternalProductID(productExternalID)
	if perr != nil {
		return fmt.Errorf("tiendanube: external_id invalido %q: %w", productExternalID, perr)
	}

	if variantID == 0 {
		return fmt.Errorf("tiendanube: el external_id %q no incluye variante, no se puede escribir stock", productExternalID)
	}

	if err := uc.client.SetVariantStock(ctx, cred, productID, variantID, quantity); err != nil {
		return err
	}

	uc.logger.Info(ctx).
		Str("integration_id", integrationID).
		Str("external_id", productExternalID).
		Str("variant_id", strconv.FormatInt(variantID, 10)).
		Int("quantity", quantity).
		Msg("Stock actualizado en Tiendanube")

	return nil
}
