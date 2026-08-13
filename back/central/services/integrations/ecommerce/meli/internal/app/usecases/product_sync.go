package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
	"github.com/secamc93/probability/back/central/shared/productmatch"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const productSyncProgressBatch = 10

func normalizeSKU(sku string) string {
	return strings.ToLower(strings.TrimSpace(sku))
}

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

const noSKUPlaceholder = "(sin SKU)"

func detectTypoSuspects(rc *reconcileContext, opts productmatch.TypoOptions) []productmatch.TypoSuspect {
	canal := make([]productmatch.TypoCandidate, 0, len(rc.outcome.OnlyInChannel))
	for _, idx := range rc.outcome.OnlyInChannel {
		m := rc.meliProducts[idx]
		canal = append(canal, productmatch.TypoCandidate{
			SKU:      m.SKU,
			Name:     m.Name,
			Quantity: m.StockQuantity,
			HasQty:   true,
			Ref:      parentRef(m),
			Label:    m.VariantLabel,
		})
	}

	propios := make([]productmatch.TypoCandidate, 0, len(rc.outcome.OnlyInProbability))
	for _, idx := range rc.outcome.OnlyInProbability {
		p := rc.probProducts[idx]
		propios = append(propios, productmatch.TypoCandidate{
			SKU:      p.SKU,
			Name:     p.Name,
			Quantity: p.StockQuantity,
			HasQty:   true,
		})
	}

	return productmatch.DetectTypos(canal, propios, opts)
}

func (uc *meliUseCase) inventoryAuthority(ctx context.Context, businessID uint) string {
	if uc.productRepo == nil {
		return productmatch.AuthorityUnknown
	}
	feeds, err := uc.productRepo.ERPFeedsInventory(ctx, businessID)
	if err != nil {
		uc.logger.Warn(ctx).Err(err).Uint("business_id", businessID).
			Msg("No se pudo determinar si el ERP alimenta el inventario")
		return productmatch.AuthorityUnknown
	}
	if feeds {
		return productmatch.AuthorityProbability
	}
	return productmatch.AuthorityUnknown
}

func channelKey(itemID, variantID string) string {
	return itemID + "|" + variantID
}

// detectSKUChanges encuentra los mapeos cuyo SKU cambio en el canal. El id de la
// variante nunca cambia, asi que el mapeo sigue "funcionando" y le empujariamos
// el stock del producto viejo a una variante que ya es otra cosa.
func detectSKUChanges(mapped []domain.MappedItem, channel []domain.MeliProduct) []domain.ProductBrief {
	actual := make(map[string]domain.MeliProduct, len(channel))
	for _, m := range channel {
		actual[channelKey(m.ID, m.VariationID)] = m
	}

	changed := make([]domain.ProductBrief, 0)
	for _, m := range mapped {
		if m.ExternalSKU == "" {
			continue
		}
		current, ok := actual[channelKey(m.ExternalItemID, m.ExternalVariantID)]
		if !ok || current.SKU == "" {
			continue
		}
		if normalizeSKU(current.SKU) == normalizeSKU(m.ExternalSKU) {
			continue
		}
		changed = append(changed, domain.ProductBrief{
			SKU:          m.SKU,
			Name:         current.Name,
			MatchedBy:    m.ExternalSKU,
			MatchedValue: current.SKU,
			ParentRef:    parentRef(current),
			ParentLabel:  current.ParentName,
			VariantLabel: current.VariantLabel,
		})
	}
	return changed
}

func parentRef(m domain.MeliProduct) string {
	if m.VariationID == "" {
		return ""
	}
	return m.ID
}

func meliRefs(m domain.MeliProduct) productmatch.ExternalRefs {
	return productmatch.ExternalRefs{
		ProductID:    m.ID,
		VariantID:    m.VariationID,
		SKU:          m.SKU,
		Barcode:      m.Barcode,
		LogisticType: m.LogisticType,
	}
}

func probabilityItems(products []domain.ProductForSync) []productmatch.Item {
	items := make([]productmatch.Item, len(products))
	for i, p := range products {
		items[i] = p.MatchItem()
	}
	return items
}

func meliItems(products []domain.MeliProduct) []productmatch.Item {
	items := make([]productmatch.Item, len(products))
	for i, p := range products {
		items[i] = p.MatchItem()
	}
	return items
}

type reconcileContext struct {
	accessToken  string
	sellerID     int64
	integration  *domain.Integration
	probProducts []domain.ProductForSync
	meliProducts []domain.MeliProduct
	rules        []productmatch.Rule
	outcome      productmatch.Outcome
}

func resolveSellerID(integration *domain.Integration) (int64, error) {
	if integration.StoreID != "" {
		if id, err := strconv.ParseInt(strings.TrimSpace(integration.StoreID), 10, 64); err == nil && id > 0 {
			return id, nil
		}
	}
	if v, ok := integration.Config["seller_id"]; ok {
		switch val := v.(type) {
		case float64:
			return int64(val), nil
		case int64:
			return val, nil
		case int:
			return int64(val), nil
		case string:
			if id, err := strconv.ParseInt(strings.TrimSpace(val), 10, 64); err == nil {
				return id, nil
			}
		}
	}
	return 0, domain.ErrSellerIDNotFound
}

func (uc *meliUseCase) loadReconcileData(ctx context.Context, integrationID string, businessID uint) (*reconcileContext, error) {
	integration, ierr := uc.service.GetIntegrationByID(ctx, integrationID)
	if ierr != nil {
		return nil, fmt.Errorf("getting integration: %w", ierr)
	}
	if integration == nil {
		return nil, domain.ErrIntegrationNotFound
	}

	sellerID, err := resolveSellerID(integration)
	if err != nil {
		return nil, err
	}

	accessToken, err := uc.EnsureValidToken(ctx, integrationID)
	if err != nil {
		return nil, err
	}

	probProducts, err := uc.productRepo.ListProductsByBusiness(ctx, businessID)
	if err != nil {
		return nil, fmt.Errorf("listing probability products: %w", err)
	}

	meliProducts, err := uc.clientFor(ctx, integration).GetProducts(ctx, accessToken, sellerID)
	if err != nil {
		return nil, fmt.Errorf("listing meli products: %w", err)
	}

	rules := productmatch.Sanitize(integration.ProductMatchRules)
	return &reconcileContext{
		accessToken:  accessToken,
		sellerID:     sellerID,
		integration:  integration,
		probProducts: probProducts,
		meliProducts: meliProducts,
		rules:        rules,
		outcome:      productmatch.Reconcile(rules, probabilityItems(probProducts), meliItems(meliProducts)),
	}, nil
}

func (uc *meliUseCase) ReconcileProducts(ctx context.Context, integrationID string, businessID uint) (*domain.ReconcileResult, error) {
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
		OnlyInMeli:           []domain.ProductBrief{},
		ProbabilityNoSKU:     rc.outcome.ProbabilityUnmatchable,
		MeliNoSKU:            rc.outcome.ChannelUnmatchable,
		MeliNoSKUItems:       []domain.ProductBrief{},
		SKUChangedItems:      []domain.ProductBrief{},
		MatchRules:           rc.rules,
	}

	result.SKUChangedItems = detectSKUChanges(mapped, rc.meliProducts)
	result.TypoSuspects = detectTypoSuspects(rc, productmatch.TypoOptions{Authority: uc.inventoryAuthority(ctx, businessID)})

	for _, pair := range rc.outcome.Pairs {
		prob := rc.probProducts[pair.ProbabilityIndex]
		channel := rc.meliProducts[pair.ChannelIndex]
		brief := domain.ProductBrief{
			SKU:          prob.SKU,
			Name:         channel.Name,
			MatchedBy:    pair.Rule.Key(),
			MatchedValue: channel.MatchItem().Values()[pair.Rule.Channel],
			ParentRef:    parentRef(channel),
			ParentLabel:  channel.ParentName,
			VariantLabel: channel.VariantLabel,
		}
		result.MatchedItems = append(result.MatchedItems, brief)
		if associatedSKUs[normalizeSKU(prob.SKU)] {
			result.Matched++
		} else {
			result.MatchedNotAssociated = append(result.MatchedNotAssociated, brief)
		}
	}

	for _, idx := range rc.outcome.ChannelNoKey {
		m := rc.meliProducts[idx]
		result.MeliNoSKUItems = append(result.MeliNoSKUItems, domain.ProductBrief{
			SKU:          noSKUPlaceholder,
			Name:         m.Name,
			ParentRef:    parentRef(m),
			ParentLabel:  m.ParentName,
			VariantLabel: m.VariantLabel,
		})
	}

	for _, idx := range rc.outcome.OnlyInChannel {
		m := rc.meliProducts[idx]
		result.OnlyInMeli = append(result.OnlyInMeli, domain.ProductBrief{
			SKU:          m.SKU,
			Name:         m.Name,
			ParentRef:    parentRef(m),
			ParentLabel:  m.ParentName,
			VariantLabel: m.VariantLabel,
		})
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

func missingPublishData(p domain.ProductForSync) string {
	missing := make([]string, 0, 3)
	if strings.TrimSpace(p.Name) == "" {
		missing = append(missing, "nombre")
	}
	if p.Price <= 0 {
		missing = append(missing, "precio")
	}
	if strings.TrimSpace(p.Brand) == "" {
		missing = append(missing, "marca")
	}
	if len(missing) == 0 {
		return ""
	}
	return "falta " + strings.Join(missing, ", ") + " en el producto"
}

const maxFailedItemsReported = 10

type failedItem struct {
	SKU   string `json:"sku"`
	Error string `json:"error"`
}

func appendFailure(items []failedItem, sku string, err error) []failedItem {
	if len(items) >= maxFailedItemsReported {
		return items
	}
	message := "error desconocido"
	if err != nil {
		message = err.Error()
	}
	if len(message) > 200 {
		message = message[:200]
	}
	return append(items, failedItem{SKU: sku, Error: message})
}

func (uc *meliUseCase) ApplyProductsToMeli(ctx context.Context, integrationID string, businessID uint, correlationID string, skus ...string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)
	rc, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	only := selectedSKUs(skus)
	missing := make([]domain.ProductForSync, 0)
	for _, idx := range rc.outcome.OnlyInProbability {
		p := rc.probProducts[idx]
		key := normalizeSKU(p.SKU)
		if key == "" || (only != nil && !only[key]) {
			continue
		}
		missing = append(missing, p)
	}

	total := len(missing)
	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "meli.product.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "to_meli",
		"total":          total,
	})

	siteID, currencyID, listingTypeID := uc.resolveProductPublishConfig(ctx, integrationID)
	cli := uc.clientFor(ctx, rc.integration)

	created, failed := 0, 0
	failures := make([]failedItem, 0)
	for i, p := range missing {
		if reason := missingPublishData(p); reason != "" {
			uc.logger.Warn(ctx).Str("sku", p.SKU).Str("motivo", reason).Msg("Producto sin datos para publicar en MercadoLibre")
			failures = appendFailure(failures, p.SKU, errors.New(reason))
			failed++
			uc.maybeProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, created, failed)
			continue
		}

		newID, cerr := cli.CreateProduct(ctx, rc.accessToken, domain.CreateProductInput{
			Name:          p.Name,
			SKU:           p.SKU,
			Price:         p.Price,
			Description:   p.Description,
			StockQuantity: p.StockQuantity,
			SiteID:        siteID,
			CurrencyID:    currencyID,
			ListingTypeID: listingTypeID,
			Brand:         p.Brand,
			ImageURL:      p.ImageURL,
			CategoryID:    p.MeliCategoryID,
			VariantAttrs:  p.VariantAttrs,
		})
		if cerr != nil {
			uc.logger.Error(ctx).Err(cerr).Str("sku", p.SKU).Msg("Error al crear producto en MercadoLibre")
			failures = appendFailure(failures, p.SKU, cerr)
			failed++
		} else {
			refs := productmatch.ExternalRefs{ProductID: newID, SKU: p.SKU, Barcode: p.Barcode}
			if merr := uc.productRepo.UpsertProductIntegrationMapping(ctx, p.ID, businessID, uint(integIDUint), refs); merr != nil {
				uc.logger.Error(ctx).Err(merr).Str("sku", p.SKU).Msg("Producto creado en MeLi pero fallo el mapeo")
			}
			created++
		}
		uc.maybeProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, created, failed)
	}

	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "meli.product.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "to_meli",
		"total":          total,
		"created":        created,
		"updated":        0,
		"failed":         failed,
		"failed_items":   failures,
	})
	return nil
}

func (uc *meliUseCase) ApplyProductsToProbability(ctx context.Context, integrationID string, businessID uint, correlationID string, skus ...string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)
	rc, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	only := selectedSKUs(skus)
	missing := make([]domain.MeliProduct, 0)
	for _, idx := range rc.outcome.OnlyInChannel {
		m := rc.meliProducts[idx]
		key := normalizeSKU(m.SKU)
		if key == "" || (only != nil && !only[key]) {
			continue
		}
		missing = append(missing, m)
	}

	total := len(missing)
	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "meli.product.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "to_probability",
		"total":          total,
	})

	if uc.rabbit != nil {
		_ = uc.rabbit.DeclareQueue(rabbitmq.QueueProductsProviderUpsert, true)
	}

	created, failed := 0, 0
	failures := make([]failedItem, 0)
	for i, m := range missing {
		msg := providerUpsertMsg{
			BusinessID:     businessID,
			IntegrationID:  uint(integIDUint),
			SKU:            m.SKU,
			Name:           m.Name,
			TrackInventory: true,
			Price:          m.Price,
			ExternalID:     m.ID,
			ImageURL:       m.ImageURL,
		}
		data, merr := json.Marshal(msg)
		if merr != nil || uc.rabbit == nil {
			failures = appendFailure(failures, m.SKU, merr)
			failed++
		} else if perr := uc.rabbit.Publish(ctx, rabbitmq.QueueProductsProviderUpsert, data); perr != nil {
			uc.logger.Error(ctx).Err(perr).Str("sku", m.SKU).Msg("Error al publicar producto para crear en Probability")
			failures = appendFailure(failures, m.SKU, perr)
			failed++
		} else {
			created++
		}
		uc.maybeProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, created, failed)
	}

	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "meli.product.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "to_probability",
		"total":          total,
		"created":        created,
		"updated":        0,
		"failed":         failed,
		"failed_items":   failures,
	})
	return nil
}

func (uc *meliUseCase) AssociateProducts(ctx context.Context, integrationID string, businessID uint, correlationID string, skus []string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)
	rc, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
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
	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "meli.product.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "associate",
		"total":          total,
	})

	created, failed := 0, 0
	for i, pair := range targets {
		p := rc.probProducts[pair.ProbabilityIndex]
		refs := meliRefs(rc.meliProducts[pair.ChannelIndex])
		if refs.ProductID == "" {
			uc.maybeProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, created, failed)
			continue
		}
		if merr := uc.productRepo.UpsertProductIntegrationMapping(ctx, p.ID, businessID, uint(integIDUint), refs); merr != nil {
			uc.logger.Error(ctx).Err(merr).Str("sku", p.SKU).Msg("Error al asociar producto a MercadoLibre")
			failed++
		} else {
			created++
		}
		uc.maybeProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, created, failed)
	}

	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "meli.product.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "associate",
		"total":          total,
		"created":        created,
		"updated":        0,
		"failed":         failed,
	})
	return nil
}

func (uc *meliUseCase) resolveProductPublishConfig(ctx context.Context, integrationID string) (string, string, string) {
	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil || integration == nil {
		return "", "", ""
	}
	site, _ := integration.Config["meli_site_id"].(string)
	currency, _ := integration.Config["meli_currency_id"].(string)
	listing, _ := integration.Config["meli_listing_type_id"].(string)
	return site, currency, listing
}

func (uc *meliUseCase) maybeProgress(ctx context.Context, businessID, integrationID uint, correlationID string, processed, total, created, failed int) {
	if processed%productSyncProgressBatch != 0 && processed != total {
		return
	}
	uc.emitSyncEvent(ctx, businessID, integrationID, "meli.product.sync.progress", map[string]interface{}{
		"correlation_id": correlationID,
		"processed":      processed,
		"total":          total,
		"created":        created,
		"updated":        0,
		"failed":         failed,
	})
}

func (uc *meliUseCase) emitSyncEvent(ctx context.Context, businessID, integrationID uint, eventType string, data map[string]interface{}) {
	if uc.rabbit == nil {
		return
	}
	_ = rabbitmq.PublishEvent(ctx, uc.rabbit, rabbitmq.EventEnvelope{
		Type:          eventType,
		Category:      "meli",
		BusinessID:    businessID,
		IntegrationID: integrationID,
		Data:          data,
	})
}
