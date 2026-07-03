package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/uuid"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
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

	total := len(products)
	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "woocommerce.product.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"total":          total,
	})

	created := 0
	updated := 0
	failed := 0

	for i, p := range products {
		if p.SKU == "" {
			failed++
			uc.maybeProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, created, updated, failed)
			continue
		}

		externalID, mapped, gerr := uc.productRepo.GetExternalProductID(ctx, p.ID, uint(integIDUint))
		if gerr != nil {
			failed++
			uc.maybeProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, created, updated, failed)
			continue
		}

		if mapped && externalID != "" {
			if perr := uc.client.UpdateProductStock(ctx, storeURL, consumerKey, consumerSecret, externalID, p.StockQuantity); perr != nil {
				uc.logger.Error(ctx).Err(perr).Str("sku", p.SKU).Msg("Error al actualizar producto en WooCommerce")
				failed++
			} else {
				updated++
			}
			uc.maybeProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, created, updated, failed)
			continue
		}

		newID, cerr := uc.client.CreateProduct(ctx, storeURL, consumerKey, consumerSecret, domain.CreateProductInput{
			Name:          p.Name,
			SKU:           p.SKU,
			Price:         p.Price,
			Description:   p.Description,
			StockQuantity: p.StockQuantity,
			ManageStock:   p.TrackInventory,
		})
		if cerr != nil {
			uc.logger.Error(ctx).Err(cerr).Str("sku", p.SKU).Msg("Error al crear producto en WooCommerce")
			failed++
			uc.maybeProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, created, updated, failed)
			continue
		}

		if merr := uc.productRepo.UpsertProductIntegrationMapping(ctx, p.ID, businessID, uint(integIDUint), newID); merr != nil {
			uc.logger.Error(ctx).Err(merr).Str("sku", p.SKU).Msg("Producto creado en Woo pero fallo el mapeo")
		}
		created++
		uc.maybeProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, created, updated, failed)
	}

	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "woocommerce.product.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"total":          total,
		"created":        created,
		"updated":        updated,
		"failed":         failed,
	})

	uc.logger.Info(ctx).
		Int("total", total).
		Int("created", created).
		Int("updated", updated).
		Int("failed", failed).
		Msg("Sincronizacion de productos a WooCommerce completada")

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
