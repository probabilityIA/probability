package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const productSyncProgressBatch = 10

type providerUpsertMsg struct {
	BusinessID     uint    `json:"business_id"`
	IntegrationID  uint    `json:"integration_id"`
	SKU            string  `json:"sku"`
	Name           string  `json:"name"`
	TrackInventory bool    `json:"track_inventory"`
	Price          float64 `json:"price"`
	ExternalID     string  `json:"external_id"`
}

func normalizeSKU(sku string) string {
	return strings.ToLower(strings.TrimSpace(sku))
}

func shopifyExternalRef(p domain.ShopifyProductForSync) string {
	if p.ProductID == "" {
		return ""
	}
	return p.ProductID + ":" + p.SKU
}

func (uc *SyncOrdersUseCase) emitProductEvent(ctx context.Context, integrationID uint, businessID uint, eventType string, data map[string]interface{}) {
	if uc.syncEventPublisher == nil {
		return
	}
	b := businessID
	uc.syncEventPublisher.PublishSyncEvent(ctx, integrationID, &b, eventType, data)
}

func (uc *SyncOrdersUseCase) maybeProductProgress(ctx context.Context, integrationID, businessID uint, correlationID string, processed, total, created, updated, failed int) {
	if processed%productSyncProgressBatch != 0 && processed != total {
		return
	}
	uc.emitProductEvent(ctx, integrationID, businessID, "shopify.product.sync.progress", map[string]interface{}{
		"correlation_id": correlationID,
		"processed":      processed,
		"total":          total,
		"created":        created,
		"updated":        updated,
		"failed":         failed,
	})
}

func (uc *SyncOrdersUseCase) loadReconcileData(ctx context.Context, integrationID string, businessID uint) ([]domain.ProductForSync, []domain.ShopifyProductForSync, error) {
	integration, err := uc.integrationService.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, nil, fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil {
		return nil, nil, fmt.Errorf("integration not found")
	}
	if integration.BusinessID == nil || *integration.BusinessID != businessID {
		return nil, nil, fmt.Errorf("la integracion no pertenece al negocio")
	}
	storeDomain, accessToken, err := uc.resolveStoreAndToken(ctx, integration, integrationID)
	if err != nil {
		return nil, nil, err
	}
	probProducts, err := uc.productRepo.ListProductsByBusiness(ctx, businessID)
	if err != nil {
		return nil, nil, fmt.Errorf("listing probability products: %w", err)
	}
	shopifyProducts, err := uc.shopifyClient.ListProducts(ctx, storeDomain, accessToken)
	if err != nil {
		return nil, nil, fmt.Errorf("listing shopify products: %w", err)
	}
	return probProducts, shopifyProducts, nil
}

func (uc *SyncOrdersUseCase) ReconcileProducts(ctx context.Context, integrationID string, businessID uint) (*domain.ReconcileResult, error) {
	probProducts, shopifyProducts, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return nil, err
	}

	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)
	mapped, merr := uc.inventoryRepo.ListMappedItems(ctx, uint(integIDUint))
	if merr != nil {
		return nil, fmt.Errorf("listing mapped items: %w", merr)
	}
	associatedSKUs := make(map[string]bool)
	for _, m := range mapped {
		if k := normalizeSKU(m.SKU); k != "" {
			associatedSKUs[k] = true
		}
	}

	result := &domain.ReconcileResult{
		MatchedNotAssociated: []domain.ProductBrief{},
		OnlyInProbability:    []domain.ProductBrief{},
		OnlyInShopify:        []domain.ProductBrief{},
	}

	probBySKU := make(map[string]domain.ProductForSync)
	for _, p := range probProducts {
		key := normalizeSKU(p.SKU)
		if key == "" {
			result.ProbabilityNoSKU++
			continue
		}
		probBySKU[key] = p
	}

	shopifySKUs := make(map[string]bool)
	for _, s := range shopifyProducts {
		key := normalizeSKU(s.SKU)
		if key == "" {
			result.ShopifyNoSKU++
			continue
		}
		shopifySKUs[key] = true
		if _, ok := probBySKU[key]; ok {
			if associatedSKUs[key] {
				result.Matched++
			} else {
				result.MatchedNotAssociated = append(result.MatchedNotAssociated, domain.ProductBrief{SKU: s.SKU, Name: s.Name})
			}
		} else {
			result.OnlyInShopify = append(result.OnlyInShopify, domain.ProductBrief{SKU: s.SKU, Name: s.Name})
		}
	}

	for key, p := range probBySKU {
		if !shopifySKUs[key] {
			result.OnlyInProbability = append(result.OnlyInProbability, domain.ProductBrief{SKU: p.SKU, Name: p.Name})
		}
	}

	return result, nil
}

func (uc *SyncOrdersUseCase) ApplyProductsToShopify(ctx context.Context, integrationID string, businessID uint, correlationID string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)
	probProducts, shopifyProducts, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		uc.emitProductEvent(ctx, uint(integIDUint), businessID, "shopify.product.sync.completed", map[string]interface{}{
			"correlation_id": correlationID, "direction": "to_shopify", "total": 0, "created": 0, "updated": 0, "failed": 0, "error": err.Error(),
		})
		return err
	}

	integration, err := uc.integrationService.GetIntegrationByID(ctx, integrationID)
	if err != nil || integration == nil {
		return err
	}
	storeDomain, accessToken, err := uc.resolveStoreAndToken(ctx, integration, integrationID)
	if err != nil {
		return err
	}

	shopifyBySKU := make(map[string]string)
	for _, s := range shopifyProducts {
		if key := normalizeSKU(s.SKU); key != "" {
			shopifyBySKU[key] = shopifyExternalRef(s)
		}
	}

	targets := make([]domain.ProductForSync, 0)
	for _, p := range probProducts {
		if normalizeSKU(p.SKU) == "" {
			continue
		}
		targets = append(targets, p)
	}

	total := len(targets)
	uc.emitProductEvent(ctx, uint(integIDUint), businessID, "shopify.product.sync.started", map[string]interface{}{
		"correlation_id": correlationID, "direction": "to_shopify", "total": total,
	})

	created, updated, failed := 0, 0, 0
	for i, p := range targets {
		key := normalizeSKU(p.SKU)
		if ref, ok := shopifyBySKU[key]; ok && ref != "" {
			if merr := uc.productRepo.UpsertProductIntegrationMapping(ctx, p.ID, businessID, uint(integIDUint), ref); merr != nil {
				failed++
			} else {
				updated++
			}
			uc.maybeProductProgress(ctx, uint(integIDUint), businessID, correlationID, i+1, total, created, updated, failed)
			continue
		}

		newRef, cerr := uc.shopifyClient.CreateProduct(ctx, storeDomain, accessToken, domain.CreateProductInput{
			Name:          p.Name,
			SKU:           p.SKU,
			Price:         p.Price,
			Description:   p.Description,
			StockQuantity: p.StockQuantity,
		})
		if cerr != nil {
			uc.log.Error(ctx).Err(cerr).Str("sku", p.SKU).Msg("Error al crear producto en Shopify")
			failed++
			uc.maybeProductProgress(ctx, uint(integIDUint), businessID, correlationID, i+1, total, created, updated, failed)
			continue
		}
		if merr := uc.productRepo.UpsertProductIntegrationMapping(ctx, p.ID, businessID, uint(integIDUint), newRef); merr != nil {
			failed++
			uc.maybeProductProgress(ctx, uint(integIDUint), businessID, correlationID, i+1, total, created, updated, failed)
			continue
		}
		created++
		uc.maybeProductProgress(ctx, uint(integIDUint), businessID, correlationID, i+1, total, created, updated, failed)
	}

	uc.emitProductEvent(ctx, uint(integIDUint), businessID, "shopify.product.sync.completed", map[string]interface{}{
		"correlation_id": correlationID, "direction": "to_shopify", "total": total, "created": created, "updated": updated, "failed": failed,
	})
	return nil
}

func (uc *SyncOrdersUseCase) ApplyProductsToProbability(ctx context.Context, integrationID string, businessID uint, correlationID string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)
	probProducts, shopifyProducts, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		uc.emitProductEvent(ctx, uint(integIDUint), businessID, "shopify.product.sync.completed", map[string]interface{}{
			"correlation_id": correlationID, "direction": "to_probability", "total": 0, "created": 0, "updated": 0, "failed": 0, "error": err.Error(),
		})
		return err
	}

	probSKUs := make(map[string]bool)
	for _, p := range probProducts {
		if key := normalizeSKU(p.SKU); key != "" {
			probSKUs[key] = true
		}
	}

	missing := make([]domain.ShopifyProductForSync, 0)
	seen := make(map[string]bool)
	for _, s := range shopifyProducts {
		key := normalizeSKU(s.SKU)
		if key == "" || probSKUs[key] || seen[key] {
			continue
		}
		seen[key] = true
		missing = append(missing, s)
	}

	total := len(missing)
	uc.emitProductEvent(ctx, uint(integIDUint), businessID, "shopify.product.sync.started", map[string]interface{}{
		"correlation_id": correlationID, "direction": "to_probability", "total": total,
	})

	if uc.rabbit != nil {
		_ = uc.rabbit.DeclareQueue(rabbitmq.QueueProductsProviderUpsert, true)
	}

	created, failed := 0, 0
	for i, s := range missing {
		msg := providerUpsertMsg{
			BusinessID:     businessID,
			IntegrationID:  uint(integIDUint),
			SKU:            s.SKU,
			Name:           s.Name,
			TrackInventory: true,
			Price:          0,
			ExternalID:     shopifyExternalRef(s),
		}
		data, merr := json.Marshal(msg)
		if merr != nil || uc.rabbit == nil {
			failed++
		} else if perr := uc.rabbit.Publish(ctx, rabbitmq.QueueProductsProviderUpsert, data); perr != nil {
			uc.log.Error(ctx).Err(perr).Str("sku", s.SKU).Msg("Error al publicar producto para crear en Probability")
			failed++
		} else {
			created++
		}
		uc.maybeProductProgress(ctx, uint(integIDUint), businessID, correlationID, i+1, total, created, 0, failed)
	}

	uc.emitProductEvent(ctx, uint(integIDUint), businessID, "shopify.product.sync.completed", map[string]interface{}{
		"correlation_id": correlationID, "direction": "to_probability", "total": total, "created": created, "updated": 0, "failed": failed,
	})
	return nil
}

func (uc *SyncOrdersUseCase) AssociateProducts(ctx context.Context, integrationID string, businessID uint, correlationID string, skus []string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)
	probProducts, shopifyProducts, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		uc.emitProductEvent(ctx, uint(integIDUint), businessID, "shopify.product.sync.completed", map[string]interface{}{
			"correlation_id": correlationID, "direction": "associate", "total": 0, "created": 0, "updated": 0, "failed": 0, "error": err.Error(),
		})
		return err
	}

	shopifyBySKU := make(map[string]string)
	for _, s := range shopifyProducts {
		if k := normalizeSKU(s.SKU); k != "" {
			shopifyBySKU[k] = shopifyExternalRef(s)
		}
	}
	probBySKU := make(map[string]domain.ProductForSync)
	for _, p := range probProducts {
		if k := normalizeSKU(p.SKU); k != "" {
			probBySKU[k] = p
		}
	}

	mapped, err := uc.inventoryRepo.ListMappedItems(ctx, uint(integIDUint))
	if err != nil {
		return err
	}
	associated := make(map[string]bool)
	for _, m := range mapped {
		if k := normalizeSKU(m.SKU); k != "" {
			associated[k] = true
		}
	}

	targets := make([]string, 0)
	if len(skus) > 0 {
		for _, s := range skus {
			if k := normalizeSKU(s); k != "" {
				targets = append(targets, k)
			}
		}
	} else {
		for k := range probBySKU {
			if shopifyBySKU[k] != "" && !associated[k] {
				targets = append(targets, k)
			}
		}
	}

	total := len(targets)
	uc.emitProductEvent(ctx, uint(integIDUint), businessID, "shopify.product.sync.started", map[string]interface{}{
		"correlation_id": correlationID, "direction": "associate", "total": total,
	})

	updated, failed := 0, 0
	for i, k := range targets {
		p, okP := probBySKU[k]
		ref, okS := shopifyBySKU[k]
		if !okP || !okS || ref == "" || associated[k] {
			uc.maybeProductProgress(ctx, uint(integIDUint), businessID, correlationID, i+1, total, 0, updated, failed)
			continue
		}
		if merr := uc.productRepo.UpsertProductIntegrationMapping(ctx, p.ID, businessID, uint(integIDUint), ref); merr != nil {
			uc.log.Error(ctx).Err(merr).Str("sku", p.SKU).Msg("Error al asociar producto a Shopify")
			failed++
		} else {
			associated[k] = true
			updated++
		}
		uc.maybeProductProgress(ctx, uint(integIDUint), businessID, correlationID, i+1, total, 0, updated, failed)
	}

	uc.emitProductEvent(ctx, uint(integIDUint), businessID, "shopify.product.sync.completed", map[string]interface{}{
		"correlation_id": correlationID, "direction": "associate", "total": total, "created": 0, "updated": updated, "failed": failed,
	})
	return nil
}
