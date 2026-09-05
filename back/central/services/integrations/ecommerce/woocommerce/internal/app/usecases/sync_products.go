package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/uuid"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
	"github.com/secamc93/probability/back/central/shared/productmatch"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const productSyncProgressBatch = 10

type productSyncRequest struct {
	IntegrationID uint   `json:"integration_id"`
	BusinessID    uint   `json:"business_id"`
	CorrelationID string `json:"correlation_id"`
}

func (uc *wooCommerceUseCase) RequestProductSync(ctx context.Context, integrationID uint, businessID uint) (string, error) {
	if integrationID == 0 || businessID == 0 {
		return "", fmt.Errorf("integration_id y business_id son requeridos")
	}
	if uc.rabbit == nil {
		return "", fmt.Errorf("cola no disponible")
	}

	correlationID := uuid.New().String()

	msg := productSyncRequest{
		IntegrationID: integrationID,
		BusinessID:    businessID,
		CorrelationID: correlationID,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}

	if err := uc.rabbit.DeclareQueue(rabbitmq.QueueWooProductSyncRequests, true); err != nil {
		return "", err
	}
	if err := uc.rabbit.Publish(ctx, rabbitmq.QueueWooProductSyncRequests, data); err != nil {
		return "", err
	}

	return correlationID, nil
}

func (uc *wooCommerceUseCase) SyncProducts(ctx context.Context, integrationID string, businessID uint, correlationID string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)

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
	storeURL = resolveEffectiveStoreURL(integration, storeURL)
	consumerKey, err := uc.service.DecryptCredential(ctx, integrationID, "consumer_key")
	if err != nil {
		return fmt.Errorf("decrypting consumer_key: %w", err)
	}
	consumerSecret, err := uc.service.DecryptCredential(ctx, integrationID, "consumer_secret")
	if err != nil {
		return fmt.Errorf("decrypting consumer_secret: %w", err)
	}

	products, err := uc.productRepo.ListProductsByBusiness(ctx, businessID)
	if err != nil {
		return fmt.Errorf("listing products: %w", err)
	}

	matchRules := productmatch.Sanitize(integration.ProductMatchRules)
	wooRefsByProduct := make(map[int]productmatch.ExternalRefs)
	wooProducts, werr := uc.client.GetProducts(ctx, storeURL, consumerKey, consumerSecret)
	if werr != nil {
		uc.logger.Error(ctx).Err(werr).
			Uint("integration_id", uint(integIDUint)).
			Msg("No se pudo listar el catalogo de WooCommerce, se aborta la sincronizacion de stock")
		return fmt.Errorf("listing woocommerce products: %w", werr)
	}
	outcome := productmatch.Reconcile(matchRules, probabilityItems(products), wooItems(wooProducts))
	for _, pair := range outcome.Pairs {
		if refs := wooRefs(wooProducts[pair.ChannelIndex]); refs.ProductID != "" {
			wooRefsByProduct[pair.ProbabilityIndex] = refs
		}
	}

	total := len(products)
	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "woocommerce.product.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"total":          total,
	})

	skipped := 0
	updated := 0
	failed := 0

	for i, p := range products {
		if p.SKU == "" {
			failed++
			uc.emitProductItem(ctx, businessID, uint(integIDUint), correlationID, p.SKU, p.Name, p.StockQuantity, "failed")
			uc.maybeStockProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, skipped, updated, failed)
			continue
		}

		externalID, mapped, gerr := uc.productRepo.GetExternalProductID(ctx, p.ID, uint(integIDUint))
		if gerr != nil {
			failed++
			uc.emitProductItem(ctx, businessID, uint(integIDUint), correlationID, p.SKU, p.Name, p.StockQuantity, "failed")
			uc.maybeStockProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, skipped, updated, failed)
			continue
		}

		if !mapped || externalID == "" {
			if refs, ok := wooRefsByProduct[i]; ok {
				if merr := uc.productRepo.UpsertProductIntegrationMapping(ctx, p.ID, businessID, uint(integIDUint), refs); merr != nil {
					uc.logger.Error(ctx).Err(merr).Str("sku", p.SKU).Msg("Error al mapear producto existente de WooCommerce")
					failed++
					uc.emitProductItem(ctx, businessID, uint(integIDUint), correlationID, p.SKU, p.Name, p.StockQuantity, "failed")
					uc.maybeStockProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, skipped, updated, failed)
					continue
				}
				externalID = refs.ProductID
				mapped = true
			}
		}

		if mapped && externalID != "" {
			if perr := uc.client.UpdateProductStock(ctx, storeURL, consumerKey, consumerSecret, externalID, p.StockQuantity); perr != nil {
				uc.logger.Error(ctx).Err(perr).Str("sku", p.SKU).Msg("Error al actualizar producto en WooCommerce")
				failed++
				uc.emitProductItem(ctx, businessID, uint(integIDUint), correlationID, p.SKU, p.Name, p.StockQuantity, "failed")
			} else {
				updated++
				uc.emitProductItem(ctx, businessID, uint(integIDUint), correlationID, p.SKU, p.Name, p.StockQuantity, "updated")
			}
			uc.maybeStockProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, skipped, updated, failed)
			continue
		}

		skipped++
		uc.emitProductItem(ctx, businessID, uint(integIDUint), correlationID, p.SKU, p.Name, p.StockQuantity, "skipped")
		uc.maybeStockProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, skipped, updated, failed)
	}

	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "woocommerce.product.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"total":          total,
		"skipped":        skipped,
		"updated":        updated,
		"failed":         failed,
	})

	uc.logger.Info(ctx).
		Int("total", total).
		Int("skipped", skipped).
		Int("updated", updated).
		Int("failed", failed).
		Msg("Sincronizacion de stock a WooCommerce completada")

	return nil
}

func (uc *wooCommerceUseCase) maybeProgress(ctx context.Context, businessID, integrationID uint, correlationID string, processed, total, created, updated, failed int) {
	if processed%productSyncProgressBatch != 0 && processed != total {
		return
	}
	uc.emitSyncEvent(ctx, businessID, integrationID, "woocommerce.product.sync.progress", map[string]interface{}{
		"correlation_id": correlationID,
		"processed":      processed,
		"total":          total,
		"created":        created,
		"updated":        updated,
		"failed":         failed,
	})
}

func (uc *wooCommerceUseCase) emitProductItem(ctx context.Context, businessID, integrationID uint, correlationID, sku, name string, quantity int, action string) {
	uc.emitSyncEvent(ctx, businessID, integrationID, "woocommerce.product.sync.item", map[string]interface{}{
		"correlation_id": correlationID,
		"sku":            sku,
		"name":           name,
		"quantity":       quantity,
		"action":         action,
	})
}

func (uc *wooCommerceUseCase) emitSyncEvent(ctx context.Context, businessID, integrationID uint, eventType string, data map[string]interface{}) {
	if uc.rabbit == nil {
		return
	}
	_ = rabbitmq.PublishEvent(ctx, uc.rabbit, rabbitmq.EventEnvelope{
		Type:          eventType,
		Category:      "woocommerce",
		BusinessID:    businessID,
		IntegrationID: integrationID,
		Data:          data,
	})
}

func (uc *wooCommerceUseCase) maybeStockProgress(ctx context.Context, businessID, integrationID uint, correlationID string, processed, total, skipped, updated, failed int) {
	if processed%productSyncProgressBatch != 0 && processed != total {
		return
	}
	uc.emitSyncEvent(ctx, businessID, integrationID, "woocommerce.product.sync.progress", map[string]interface{}{
		"correlation_id": correlationID,
		"processed":      processed,
		"total":          total,
		"skipped":        skipped,
		"updated":        updated,
		"failed":         failed,
	})
}
