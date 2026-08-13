package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
	"github.com/secamc93/probability/back/central/shared/productmatch"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func normalizeSKU(sku string) string {
	return strings.ToLower(strings.TrimSpace(sku))
}

func wooExternalRef(w domain.WooProduct) string {
	if w.ParentID != "" {
		return w.ParentID + ":" + w.ID
	}
	return w.ID
}

func wooRefs(w domain.WooProduct) productmatch.ExternalRefs {
	variantID := ""
	if w.ParentID != "" {
		variantID = w.ID
	}
	return productmatch.ExternalRefs{
		ProductID: wooExternalRef(w),
		VariantID: variantID,
		SKU:       w.SKU,
		Barcode:   w.Barcode,
	}
}

func probabilityItems(products []domain.ProductForSync) []productmatch.Item {
	items := make([]productmatch.Item, len(products))
	for i, p := range products {
		items[i] = p.MatchItem()
	}
	return items
}

func wooItems(products []domain.WooProduct) []productmatch.Item {
	items := make([]productmatch.Item, len(products))
	for i, p := range products {
		items[i] = p.MatchItem()
	}
	return items
}

type reconcileContext struct {
	storeURL     string
	ck           string
	cs           string
	probProducts []domain.ProductForSync
	wooProducts  []domain.WooProduct
	rules        []productmatch.Rule
	outcome      productmatch.Outcome
}

func fullImageURL(imageURL string) string {
	imageURL = strings.TrimSpace(imageURL)
	if imageURL == "" || strings.HasPrefix(imageURL, "http://") || strings.HasPrefix(imageURL, "https://") {
		return imageURL
	}
	base := strings.TrimRight(os.Getenv("URL_BASE_DOMAIN_S3"), "/")
	if base == "" {
		return ""
	}
	return base + "/" + strings.TrimLeft(imageURL, "/")
}

type providerUpsertMsg struct {
	BusinessID        uint              `json:"business_id"`
	IntegrationID     uint              `json:"integration_id"`
	SKU               string            `json:"sku"`
	Name              string            `json:"name"`
	TrackInventory    bool              `json:"track_inventory"`
	Price             float64           `json:"price"`
	ExternalID        string            `json:"external_id"`
	ImageURL          string            `json:"image_url,omitempty"`
	FamilyName        string            `json:"family_name,omitempty"`
	FamilyRef         string            `json:"family_ref,omitempty"`
	VariantLabel      string            `json:"variant_label,omitempty"`
	VariantAttributes map[string]string `json:"variant_attributes,omitempty"`
}

func (uc *wooCommerceUseCase) loadReconcileData(ctx context.Context, integrationID string, businessID uint) (*reconcileContext, error) {
	storeURL, ck, cs, err := uc.resolveStoreCreds(ctx, integrationID)
	if err != nil {
		return nil, err
	}
	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("getting integration: %w", err)
	}
	probProducts, err := uc.productRepo.ListProductsByBusiness(ctx, businessID)
	if err != nil {
		return nil, fmt.Errorf("listing probability products: %w", err)
	}
	wooProducts, err := uc.client.GetProducts(ctx, storeURL, ck, cs)
	if err != nil {
		return nil, fmt.Errorf("listing woocommerce products: %w", err)
	}
	rules := productmatch.Sanitize(integration.ProductMatchRules)
	return &reconcileContext{
		storeURL:     storeURL,
		ck:           ck,
		cs:           cs,
		probProducts: probProducts,
		wooProducts:  wooProducts,
		rules:        rules,
		outcome:      productmatch.Reconcile(rules, probabilityItems(probProducts), wooItems(wooProducts)),
	}, nil
}

func (uc *wooCommerceUseCase) ReconcileProducts(ctx context.Context, integrationID string, businessID uint) (*domain.ReconcileResult, error) {
	rc, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return nil, err
	}

	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)
	mapped, merr := uc.productRepo.ListMappedItems(ctx, uint(integIDUint))
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
		OnlyInWoo:            []domain.ProductBrief{},
		ProbabilityNoSKU:     rc.outcome.ProbabilityUnmatchable,
		WooNoSKU:             rc.outcome.ChannelUnmatchable,
		MatchRules:           rc.rules,
	}

	snapshots := make([]productmatch.SnapshotEntry, 0, len(rc.outcome.Pairs))
	for _, pair := range rc.outcome.Pairs {
		prob := rc.probProducts[pair.ProbabilityIndex]
		channel := rc.wooProducts[pair.ChannelIndex]
		brief := domain.ProductBrief{
			SKU:          prob.SKU,
			Name:         channel.Name,
			MatchedBy:    pair.Rule.Key(),
			MatchedValue: channel.MatchItem().Values()[pair.Rule.Channel],
		}
		result.MatchedItems = append(result.MatchedItems, brief)
		snapshots = append(snapshots, productmatch.SnapshotEntry{
			ProbabilityID: prob.ID,
			Snapshot:      wooSnapshot(channel),
		})
		if associatedSKUs[normalizeSKU(prob.SKU)] {
			result.Matched++
		} else {
			result.MatchedNotAssociated = append(result.MatchedNotAssociated, brief)
		}
	}

	if serr := uc.productRepo.SaveChannelSnapshots(ctx, businessID, uint(integIDUint), snapshots); serr != nil {
		uc.logger.Warn(ctx).Err(serr).Uint("business_id", businessID).
			Msg("No se pudo guardar la foto del catalogo de WooCommerce")
	}

	for _, idx := range rc.outcome.OnlyInChannel {
		w := rc.wooProducts[idx]
		result.OnlyInWoo = append(result.OnlyInWoo, domain.ProductBrief{SKU: w.SKU, Name: w.Name})
	}
	for _, idx := range rc.outcome.OnlyInProbability {
		p := rc.probProducts[idx]
		result.OnlyInProbability = append(result.OnlyInProbability, domain.ProductBrief{SKU: p.SKU, Name: p.Name})
	}

	result.TypoSuspects = detectTypoSuspects(rc, productmatch.TypoOptions{Authority: uc.inventoryAuthority(ctx, businessID)})

	return result, nil
}

func detectTypoSuspects(rc *reconcileContext, opts productmatch.TypoOptions) []productmatch.TypoSuspect {
	canal := make([]productmatch.TypoCandidate, 0, len(rc.outcome.OnlyInChannel))
	for _, idx := range rc.outcome.OnlyInChannel {
		w := rc.wooProducts[idx]
		canal = append(canal, productmatch.TypoCandidate{
			SKU:      w.SKU,
			Name:     w.Name,
			Quantity: w.StockQuantity,
			HasQty:   true,
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

func (uc *wooCommerceUseCase) inventoryAuthority(ctx context.Context, businessID uint) string {
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

func (uc *wooCommerceUseCase) ApplyProductsToWoo(ctx context.Context, integrationID string, businessID uint, correlationID string, skus ...string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)
	rc, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	matchedRefs := make(map[int]productmatch.ExternalRefs, len(rc.outcome.Pairs))
	for _, pair := range rc.outcome.Pairs {
		matchedRefs[pair.ProbabilityIndex] = wooRefs(rc.wooProducts[pair.ChannelIndex])
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
	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "woocommerce.product.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "to_woo",
		"total":          total,
	})

	created, updated, failed := 0, 0, 0
	for n, idx := range targets {
		p := rc.probProducts[idx]
		if refs, ok := matchedRefs[idx]; ok && refs.ProductID != "" {
			if merr := uc.productRepo.UpsertProductIntegrationMapping(ctx, p.ID, businessID, uint(integIDUint), refs); merr != nil {
				uc.logger.Error(ctx).Err(merr).Str("sku", p.SKU).Msg("Error al mapear producto existente de WooCommerce")
				failed++
			} else {
				updated++
			}
			uc.maybeProgress(ctx, businessID, uint(integIDUint), correlationID, n+1, total, created, updated, failed)
			continue
		}

		newID, cerr := uc.client.CreateProduct(ctx, rc.storeURL, rc.ck, rc.cs, domain.CreateProductInput{
			Name:          p.Name,
			SKU:           p.SKU,
			Price:         p.Price,
			Description:   p.Description,
			StockQuantity: p.StockQuantity,
			ManageStock:   p.TrackInventory,
			ImageURL:      fullImageURL(p.ImageURL),
		})
		if cerr != nil {
			uc.logger.Error(ctx).Err(cerr).Str("sku", p.SKU).Msg("Error al crear producto en WooCommerce")
			failed++
			uc.maybeProgress(ctx, businessID, uint(integIDUint), correlationID, n+1, total, created, updated, failed)
			continue
		}
		refs := productmatch.ExternalRefs{ProductID: newID, SKU: p.SKU, Barcode: p.Barcode}
		if merr := uc.productRepo.UpsertProductIntegrationMapping(ctx, p.ID, businessID, uint(integIDUint), refs); merr != nil {
			uc.logger.Error(ctx).Err(merr).Str("sku", p.SKU).Msg("Producto creado en Woo pero fallo el mapeo")
			failed++
			uc.maybeProgress(ctx, businessID, uint(integIDUint), correlationID, n+1, total, created, updated, failed)
			continue
		}
		created++
		uc.maybeProgress(ctx, businessID, uint(integIDUint), correlationID, n+1, total, created, updated, failed)
	}

	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "woocommerce.product.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "to_woo",
		"total":          total,
		"created":        created,
		"updated":        updated,
		"failed":         failed,
	})
	return nil
}

func (uc *wooCommerceUseCase) ApplyProductsToProbability(ctx context.Context, integrationID string, businessID uint, correlationID string, skus ...string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)
	rc, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	only := selectedSKUs(skus)
	missing := make([]domain.WooProduct, 0)
	for _, idx := range rc.outcome.OnlyInChannel {
		w := rc.wooProducts[idx]
		key := normalizeSKU(w.SKU)
		if key == "" || (only != nil && !only[key]) {
			continue
		}
		missing = append(missing, w)
	}

	total := len(missing)
	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "woocommerce.product.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "to_probability",
		"total":          total,
	})

	if uc.rabbit != nil {
		_ = uc.rabbit.DeclareQueue(rabbitmq.QueueProductsProviderUpsert, true)
	}

	created, failed := 0, 0
	for i, w := range missing {
		msg := providerUpsertMsg{
			BusinessID:     businessID,
			IntegrationID:  uint(integIDUint),
			SKU:            w.SKU,
			Name:           w.Name,
			TrackInventory: true,
			Price:          w.Price,
			ExternalID:     wooExternalRef(w),
			ImageURL:       w.ImageURL,
		}
		if len(w.VariantAttributes) > 0 && w.ParentID != "" {
			msg.FamilyName = nombreDeFamiliaWoo(w)
			msg.FamilyRef = "woo:" + w.ParentID
			msg.VariantLabel = w.VariantLabel
			msg.VariantAttributes = w.VariantAttributes
		}
		data, merr := json.Marshal(msg)
		if merr != nil || uc.rabbit == nil {
			failed++
		} else if perr := uc.rabbit.Publish(ctx, rabbitmq.QueueProductsProviderUpsert, data); perr != nil {
			uc.logger.Error(ctx).Err(perr).Str("sku", w.SKU).Msg("Error al publicar producto para crear en Probability")
			failed++
		} else {
			created++
		}
		uc.maybeProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, created, 0, failed)
	}

	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "woocommerce.product.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "to_probability",
		"total":          total,
		"created":        created,
		"updated":        0,
		"failed":         failed,
	})
	return nil
}

func (uc *wooCommerceUseCase) AssociateProducts(ctx context.Context, integrationID string, businessID uint, correlationID string, skus []string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)
	rc, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	mappedItems, err := uc.productRepo.ListMappedItems(ctx, uint(integIDUint))
	if err != nil {
		return err
	}
	associated := make(map[string]bool)
	for _, m := range mappedItems {
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
	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "woocommerce.product.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "associate",
		"total":          total,
	})

	updated, failed := 0, 0
	for i, pair := range targets {
		p := rc.probProducts[pair.ProbabilityIndex]
		refs := wooRefs(rc.wooProducts[pair.ChannelIndex])
		if refs.ProductID == "" {
			uc.maybeProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, 0, updated, failed)
			continue
		}
		if merr := uc.productRepo.UpsertProductIntegrationMapping(ctx, p.ID, businessID, uint(integIDUint), refs); merr != nil {
			uc.logger.Error(ctx).Err(merr).Str("sku", p.SKU).Msg("Error al asociar producto a WooCommerce")
			failed++
			uc.emitProductItem(ctx, businessID, uint(integIDUint), correlationID, p.SKU, p.Name, 0, "failed")
		} else {
			updated++
			uc.emitProductItem(ctx, businessID, uint(integIDUint), correlationID, p.SKU, p.Name, 0, "updated")
		}
		uc.maybeProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, 0, updated, failed)
	}

	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "woocommerce.product.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "associate",
		"total":          total,
		"created":        0,
		"updated":        updated,
		"failed":         failed,
	})
	return nil
}
