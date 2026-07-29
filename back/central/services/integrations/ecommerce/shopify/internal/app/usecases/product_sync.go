package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
	"github.com/secamc93/probability/back/central/shared/productmatch"
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
	ImageURL       string  `json:"image_url,omitempty"`
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

func shopifyRefs(p domain.ShopifyProductForSync) productmatch.ExternalRefs {
	return productmatch.ExternalRefs{
		ProductID: shopifyExternalRef(p),
		VariantID: p.VariantID,
		SKU:       p.SKU,
		Barcode:   p.Barcode,
	}
}

func probabilityItems(products []domain.ProductForSync) []productmatch.Item {
	items := make([]productmatch.Item, len(products))
	for i, p := range products {
		items[i] = p.MatchItem()
	}
	return items
}

func shopifyItems(products []domain.ShopifyProductForSync) []productmatch.Item {
	items := make([]productmatch.Item, len(products))
	for i, p := range products {
		items[i] = p.MatchItem()
	}
	return items
}

func matchedValue(item productmatch.Item, field string) string {
	return item.Values()[field]
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

type reconcileContext struct {
	probProducts    []domain.ProductForSync
	shopifyProducts []domain.ShopifyProductForSync
	rules           []productmatch.Rule
	outcome         productmatch.Outcome
	storeDomain     string
	accessToken     string
}

func (uc *SyncOrdersUseCase) loadReconcileData(ctx context.Context, integrationID string, businessID uint) (*reconcileContext, error) {
	integration, err := uc.integrationService.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil {
		return nil, fmt.Errorf("integration not found")
	}
	if integration.BusinessID == nil || *integration.BusinessID != businessID {
		return nil, fmt.Errorf("la integracion no pertenece al negocio")
	}
	storeDomain, accessToken, err := uc.resolveStoreAndToken(ctx, integration, integrationID)
	if err != nil {
		return nil, err
	}
	probProducts, err := uc.productRepo.ListProductsByBusiness(ctx, businessID)
	if err != nil {
		return nil, fmt.Errorf("listing probability products: %w", err)
	}
	shopifyProducts, err := uc.shopifyClient.ListProducts(ctx, storeDomain, accessToken)
	if err != nil {
		return nil, fmt.Errorf("listing shopify products: %w", err)
	}

	rules := productmatch.Sanitize(integration.ProductMatchRules)
	return &reconcileContext{
		probProducts:    probProducts,
		shopifyProducts: shopifyProducts,
		rules:           rules,
		outcome:         productmatch.Reconcile(rules, probabilityItems(probProducts), shopifyItems(shopifyProducts)),
		storeDomain:     storeDomain,
		accessToken:     accessToken,
	}, nil
}

func (uc *SyncOrdersUseCase) ReconcileProducts(ctx context.Context, integrationID string, businessID uint) (*domain.ReconcileResult, error) {
	rc, err := uc.loadReconcileData(ctx, integrationID, businessID)
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
		MatchedItems:         []domain.ProductBrief{},
		MatchedNotAssociated: []domain.ProductBrief{},
		OnlyInProbability:    []domain.ProductBrief{},
		OnlyInShopify:        []domain.ProductBrief{},
		ProbabilityNoSKU:     rc.outcome.ProbabilityUnmatchable,
		ShopifyNoSKU:         rc.outcome.ChannelUnmatchable,
		MatchRules:           rc.rules,
	}

	for _, pair := range rc.outcome.Pairs {
		prob := rc.probProducts[pair.ProbabilityIndex]
		channel := rc.shopifyProducts[pair.ChannelIndex]
		brief := domain.ProductBrief{
			SKU:          prob.SKU,
			Name:         channel.Name,
			MatchedBy:    pair.Rule.Key(),
			MatchedValue: matchedValue(channel.MatchItem(), pair.Rule.Channel),
		}
		result.MatchedItems = append(result.MatchedItems, brief)
		if associatedSKUs[normalizeSKU(prob.SKU)] {
			result.Matched++
		} else {
			result.MatchedNotAssociated = append(result.MatchedNotAssociated, brief)
		}
	}

	for _, idx := range rc.outcome.OnlyInChannel {
		s := rc.shopifyProducts[idx]
		result.OnlyInShopify = append(result.OnlyInShopify, domain.ProductBrief{SKU: s.SKU, Name: s.Name})
	}
	for _, idx := range rc.outcome.OnlyInProbability {
		p := rc.probProducts[idx]
		result.OnlyInProbability = append(result.OnlyInProbability, domain.ProductBrief{SKU: p.SKU, Name: p.Name})
	}

	return result, nil
}

func selectedSKUs(skus []string) map[string]bool {
	if len(skus) == 0 {
		return nil
	}
	set := make(map[string]bool, len(skus))
	for _, s := range skus {
		if key := normalizeSKU(s); key != "" {
			set[key] = true
		}
	}
	return set
}

func (uc *SyncOrdersUseCase) ApplyProductsToShopify(ctx context.Context, integrationID string, businessID uint, correlationID string, skus ...string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)
	rc, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		uc.emitProductEvent(ctx, uint(integIDUint), businessID, "shopify.product.sync.completed", map[string]interface{}{
			"correlation_id": correlationID, "direction": "to_shopify", "total": 0, "created": 0, "updated": 0, "failed": 0, "error": err.Error(),
		})
		return err
	}

	matchedRefs := make(map[int]productmatch.ExternalRefs, len(rc.outcome.Pairs))
	for _, pair := range rc.outcome.Pairs {
		matchedRefs[pair.ProbabilityIndex] = shopifyRefs(rc.shopifyProducts[pair.ChannelIndex])
	}

	only := selectedSKUs(skus)
	targets := make([]int, 0, len(rc.probProducts))
	for i, p := range rc.probProducts {
		key := normalizeSKU(p.SKU)
		if key == "" || (only != nil && !only[key]) {
			continue
		}
		targets = append(targets, i)
	}

	total := len(targets)
	uc.emitProductEvent(ctx, uint(integIDUint), businessID, "shopify.product.sync.started", map[string]interface{}{
		"correlation_id": correlationID, "direction": "to_shopify", "total": total,
	})

	created, updated, failed := 0, 0, 0
	for n, idx := range targets {
		p := rc.probProducts[idx]
		if refs, ok := matchedRefs[idx]; ok && refs.ProductID != "" {
			if merr := uc.productRepo.UpsertProductIntegrationMapping(ctx, p.ID, businessID, uint(integIDUint), refs); merr != nil {
				failed++
			} else {
				updated++
			}
			uc.maybeProductProgress(ctx, uint(integIDUint), businessID, correlationID, n+1, total, created, updated, failed)
			continue
		}

		newRef, cerr := uc.shopifyClient.CreateProduct(ctx, rc.storeDomain, rc.accessToken, domain.CreateProductInput{
			Name:          p.Name,
			SKU:           p.SKU,
			Price:         p.Price,
			Description:   p.Description,
			StockQuantity: p.StockQuantity,
		})
		if cerr != nil {
			uc.log.Error(ctx).Err(cerr).Str("sku", p.SKU).Msg("Error al crear producto en Shopify")
			failed++
			uc.maybeProductProgress(ctx, uint(integIDUint), businessID, correlationID, n+1, total, created, updated, failed)
			continue
		}
		refs := productmatch.ExternalRefs{ProductID: newRef, SKU: p.SKU, Barcode: p.Barcode}
		if merr := uc.productRepo.UpsertProductIntegrationMapping(ctx, p.ID, businessID, uint(integIDUint), refs); merr != nil {
			failed++
			uc.maybeProductProgress(ctx, uint(integIDUint), businessID, correlationID, n+1, total, created, updated, failed)
			continue
		}
		created++
		uc.maybeProductProgress(ctx, uint(integIDUint), businessID, correlationID, n+1, total, created, updated, failed)
	}

	uc.emitProductEvent(ctx, uint(integIDUint), businessID, "shopify.product.sync.completed", map[string]interface{}{
		"correlation_id": correlationID, "direction": "to_shopify", "total": total, "created": created, "updated": updated, "failed": failed,
	})
	return nil
}

func (uc *SyncOrdersUseCase) ApplyProductsToProbability(ctx context.Context, integrationID string, businessID uint, correlationID string, skus ...string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)
	rc, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		uc.emitProductEvent(ctx, uint(integIDUint), businessID, "shopify.product.sync.completed", map[string]interface{}{
			"correlation_id": correlationID, "direction": "to_probability", "total": 0, "created": 0, "updated": 0, "failed": 0, "error": err.Error(),
		})
		return err
	}

	only := selectedSKUs(skus)
	missing := make([]domain.ShopifyProductForSync, 0)
	seen := make(map[string]bool)
	for _, idx := range rc.outcome.OnlyInChannel {
		s := rc.shopifyProducts[idx]
		key := normalizeSKU(s.SKU)
		if key == "" || seen[key] || (only != nil && !only[key]) {
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
			ImageURL:       s.ImageURL,
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
	rc, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		uc.emitProductEvent(ctx, uint(integIDUint), businessID, "shopify.product.sync.completed", map[string]interface{}{
			"correlation_id": correlationID, "direction": "associate", "total": 0, "created": 0, "updated": 0, "failed": 0, "error": err.Error(),
		})
		return err
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

	only := selectedSKUs(skus)
	targets := make([]productmatch.Pair, 0, len(rc.outcome.Pairs))
	for _, pair := range rc.outcome.Pairs {
		key := normalizeSKU(rc.probProducts[pair.ProbabilityIndex].SKU)
		if only != nil {
			if key == "" || !only[key] {
				continue
			}
		} else if associated[key] {
			continue
		}
		targets = append(targets, pair)
	}

	total := len(targets)
	uc.emitProductEvent(ctx, uint(integIDUint), businessID, "shopify.product.sync.started", map[string]interface{}{
		"correlation_id": correlationID, "direction": "associate", "total": total,
	})

	updated, failed := 0, 0
	for i, pair := range targets {
		p := rc.probProducts[pair.ProbabilityIndex]
		refs := shopifyRefs(rc.shopifyProducts[pair.ChannelIndex])
		if refs.ProductID == "" {
			uc.maybeProductProgress(ctx, uint(integIDUint), businessID, correlationID, i+1, total, 0, updated, failed)
			continue
		}
		if merr := uc.productRepo.UpsertProductIntegrationMapping(ctx, p.ID, businessID, uint(integIDUint), refs); merr != nil {
			uc.log.Error(ctx).Err(merr).Str("sku", p.SKU).Msg("Error al asociar producto a Shopify")
			failed++
		} else {
			updated++
		}
		uc.maybeProductProgress(ctx, uint(integIDUint), businessID, correlationID, i+1, total, 0, updated, failed)
	}

	uc.emitProductEvent(ctx, uint(integIDUint), businessID, "shopify.product.sync.completed", map[string]interface{}{
		"correlation_id": correlationID, "direction": "associate", "total": total, "created": 0, "updated": updated, "failed": failed,
	})
	return nil
}
