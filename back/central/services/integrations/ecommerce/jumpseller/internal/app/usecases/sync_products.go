package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
	"github.com/secamc93/probability/back/central/shared/productmatch"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const (
	DirectionToJumpseller  = "to_jumpseller"
	DirectionToProbability = "to_probability"

	ModeCreate = "create"
	ModeUpdate = "update"
)

type providerUpsertMsg struct {
	BusinessID     uint     `json:"business_id"`
	IntegrationID  uint     `json:"integration_id"`
	SKU            string   `json:"sku"`
	Name           string   `json:"name"`
	TrackInventory bool     `json:"track_inventory"`
	Price          float64  `json:"price"`
	ExternalID     string   `json:"external_id"`
	Weight         *float64 `json:"weight,omitempty"`
	WeightUnit     string   `json:"weight_unit,omitempty"`
	Length         *float64 `json:"length,omitempty"`
	Width          *float64 `json:"width,omitempty"`
	Height         *float64 `json:"height,omitempty"`
	DimensionUnit  string   `json:"dimension_unit,omitempty"`
	ImageURL       string   `json:"image_url,omitempty"`
}

type jumpsellerSKU struct {
	SKU        string
	Barcode    string
	Name       string
	ImageURL   string
	Price      float64
	ExternalID string
	ProductID  int64
	VariantID  int64
	Weight     float64
	Height     float64
	Width      float64
	Length     float64
}

type reconcileData struct {
	cred            domain.Credential
	storeWeightUnit string
	probability     []domain.ProductForSync
	jumpseller      []jumpsellerSKU
	rules           []productmatch.Rule
	outcome         productmatch.Outcome
}

func (j jumpsellerSKU) MatchItem() productmatch.Item {
	variantID := ""
	if j.VariantID != 0 {
		variantID = strconv.FormatInt(j.VariantID, 10)
	}
	return productmatch.Item{
		SKU:        j.SKU,
		Barcode:    j.Barcode,
		ExternalID: strconv.FormatInt(j.ProductID, 10),
		VariantID:  variantID,
		Name:       j.Name,
	}
}

func jumpsellerRefs(j jumpsellerSKU) productmatch.ExternalRefs {
	variantID := ""
	if j.VariantID != 0 {
		variantID = strconv.FormatInt(j.VariantID, 10)
	}
	return productmatch.ExternalRefs{
		ProductID: j.ExternalID,
		VariantID: variantID,
		SKU:       j.SKU,
		Barcode:   j.Barcode,
	}
}

func probabilityItems(products []domain.ProductForSync) []productmatch.Item {
	items := make([]productmatch.Item, len(products))
	for i, p := range products {
		items[i] = p.MatchItem()
	}
	return items
}

func jumpsellerItems(skus []jumpsellerSKU) []productmatch.Item {
	items := make([]productmatch.Item, len(skus))
	for i, j := range skus {
		items[i] = j.MatchItem()
	}
	return items
}

func normalizeSKU(sku string) string {
	return strings.ToLower(strings.TrimSpace(sku))
}

func flattenProductSKUs(products []domain.JumpsellerProduct) []jumpsellerSKU {
	flat := make([]jumpsellerSKU, 0, len(products))
	for _, product := range products {
		if len(product.Variants) == 0 {
			flat = append(flat, jumpsellerSKU{
				SKU:        product.SKU,
				Barcode:    product.Barcode,
				Name:       product.Name,
				ImageURL:   product.ImageURL,
				Price:      product.Price,
				ExternalID: strconv.FormatInt(product.ID, 10),
				ProductID:  product.ID,
				Weight:     product.Weight,
				Height:     product.Height,
				Width:      product.Width,
				Length:     product.Length,
			})
			continue
		}
		for _, variant := range product.Variants {
			flat = append(flat, jumpsellerSKU{
				SKU:        variant.SKU,
				Barcode:    variant.Barcode,
				Name:       product.Name,
				ImageURL:   product.ImageURL,
				Price:      variant.Price,
				ExternalID: strconv.FormatInt(product.ID, 10) + ":" + strconv.FormatInt(variant.ID, 10),
				ProductID:  product.ID,
				VariantID:  variant.ID,
				Weight:     product.Weight,
				Height:     product.Height,
				Width:      product.Width,
				Length:     product.Length,
			})
		}
	}
	return flat
}

func matchedByProbabilityIndex(data *reconcileData) map[int]jumpsellerSKU {
	matched := make(map[int]jumpsellerSKU, len(data.outcome.Pairs))
	for _, pair := range data.outcome.Pairs {
		matched[pair.ProbabilityIndex] = data.jumpseller[pair.ChannelIndex]
	}
	return matched
}

func (uc *jumpsellerUseCase) loadReconcileData(ctx context.Context, integrationID string, businessID uint) (*reconcileData, error) {
	_, cred, err := uc.resolveIntegrationForBusiness(ctx, integrationID, businessID)
	if err != nil {
		return nil, err
	}

	storeInfo, err := uc.client.GetStoreInfo(ctx, cred)
	if err != nil {
		return nil, fmt.Errorf("obteniendo informacion de la tienda Jumpseller: %w", err)
	}

	probProducts, err := uc.productRepo.ListProductsByBusiness(ctx, businessID)
	if err != nil {
		return nil, fmt.Errorf("listing probability products: %w", err)
	}

	jsProducts, err := uc.client.GetProducts(ctx, cred)
	if err != nil {
		return nil, fmt.Errorf("listing jumpseller products: %w", err)
	}

	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("getting integration: %w", err)
	}

	flat := flattenProductSKUs(jsProducts)
	rules := productmatch.Sanitize(integration.ProductMatchRules)
	return &reconcileData{
		cred:            cred,
		storeWeightUnit: storeInfo.WeightUnit,
		probability:     probProducts,
		jumpseller:      flat,
		rules:           rules,
		outcome:         productmatch.Reconcile(rules, probabilityItems(probProducts), jumpsellerItems(flat)),
	}, nil
}

func (uc *jumpsellerUseCase) associatedSKUs(ctx context.Context, integrationID uint) (map[string]bool, error) {
	mapped, err := uc.productRepo.ListMappedItems(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("listing mapped items: %w", err)
	}
	associated := make(map[string]bool, len(mapped))
	for _, item := range mapped {
		if key := normalizeSKU(item.SKU); key != "" {
			associated[key] = true
		}
	}
	return associated, nil
}

func (uc *jumpsellerUseCase) ReconcileProducts(ctx context.Context, integrationID string, businessID uint) (*domain.ReconcileResult, error) {
	data, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return nil, err
	}

	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)
	associated, err := uc.associatedSKUs(ctx, uint(integIDUint))
	if err != nil {
		return nil, err
	}

	result := &domain.ReconcileResult{
		MatchedItems:         []domain.ProductBrief{},
		MatchedNotAssociated: []domain.ProductBrief{},
		OnlyInProbability:    []domain.ProductBrief{},
		OnlyInJumpseller:     []domain.ProductBrief{},
		ProbabilityNoSKU:     data.outcome.ProbabilityUnmatchable,
		JumpsellerNoSKU:      data.outcome.ChannelUnmatchable,
		MatchRules:           data.rules,
	}

	for _, pair := range data.outcome.Pairs {
		prob := data.probability[pair.ProbabilityIndex]
		channel := data.jumpseller[pair.ChannelIndex]
		brief := domain.ProductBrief{
			SKU:          prob.SKU,
			Name:         channel.Name,
			MatchedBy:    pair.Rule.Key(),
			MatchedValue: channel.MatchItem().Values()[pair.Rule.Channel],
		}
		result.MatchedItems = append(result.MatchedItems, brief)
		if associated[normalizeSKU(prob.SKU)] {
			result.Matched++
		} else {
			result.MatchedNotAssociated = append(result.MatchedNotAssociated, brief)
		}
	}

	for _, idx := range data.outcome.OnlyInChannel {
		j := data.jumpseller[idx]
		result.OnlyInJumpseller = append(result.OnlyInJumpseller, domain.ProductBrief{SKU: j.SKU, Name: j.Name})
	}
	for _, idx := range data.outcome.OnlyInProbability {
		p := data.probability[idx]
		result.OnlyInProbability = append(result.OnlyInProbability, domain.ProductBrief{SKU: p.SKU, Name: p.Name})
	}

	return result, nil
}

func (uc *jumpsellerUseCase) probabilityWeightForStore(ctx context.Context, p domain.ProductForSync, storeUnit string) *float64 {
	if p.Weight == nil || *p.Weight <= 0 {
		return nil
	}

	weightKg := *p.Weight
	if unit := strings.ToLower(strings.TrimSpace(p.WeightUnit)); unit != "" && unit != probabilityWeightUnit {
		factor, known := weightFactor(unit)
		if !known {
			uc.logger.Warn(ctx).
				Str("sku", p.SKU).
				Str("weight_unit", p.WeightUnit).
				Msg("Unidad de peso desconocida en el producto de Probability, no se envia el peso a Jumpseller")
			return nil
		}
		weightKg = *p.Weight * factor
	}

	converted, ok := convertKgToStoreUnit(weightKg, storeUnit)
	if !ok {
		uc.logger.Warn(ctx).
			Str("sku", p.SKU).
			Str("store_weight_unit", storeUnit).
			Msg("La tienda Jumpseller no reporta una unidad de peso conocida, no se envia el peso")
		return nil
	}
	return positive(converted)
}

func (uc *jumpsellerUseCase) upsertMsgFromJumpseller(ctx context.Context, businessID, integrationID uint, j jumpsellerSKU, storeUnit string) providerUpsertMsg {
	msg := providerUpsertMsg{
		BusinessID:     businessID,
		IntegrationID:  integrationID,
		ImageURL:       j.ImageURL,
		SKU:            j.SKU,
		Name:           j.Name,
		TrackInventory: true,
		Price:          j.Price,
		ExternalID:     j.ExternalID,
		Length:         positive(j.Length),
		Width:          positive(j.Width),
		Height:         positive(j.Height),
	}

	if msg.Length != nil || msg.Width != nil || msg.Height != nil {
		msg.DimensionUnit = probabilityDimensionUnit
	}

	if j.Weight > 0 {
		weightKg, ok := normalizeWeightToKg(j.Weight, storeUnit)
		if !ok {
			uc.logger.Warn(ctx).
				Str("sku", j.SKU).
				Str("store_weight_unit", storeUnit).
				Msg("La tienda Jumpseller no reporta una unidad de peso conocida, el peso no se importa")
		} else {
			msg.Weight = positive(weightKg)
			msg.WeightUnit = probabilityWeightUnit
		}
	}

	return msg
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

func (uc *jumpsellerUseCase) ApplyProductsToJumpseller(ctx context.Context, integrationID string, businessID uint, correlationID string, skus ...string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)

	data, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	matched := matchedByProbabilityIndex(data)

	only := selectedSKUs(skus)
	targets := make([]int, 0, len(data.probability))
	for i, p := range data.probability {
		key := normalizeSKU(p.SKU)
		if key == "" || (only != nil && !only[key]) {
			continue
		}
		targets = append(targets, i)
	}

	total := len(targets)
	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "jumpseller.product.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      DirectionToJumpseller,
		"mode":           ModeCreate,
		"total":          total,
	})

	fails := &failedSKUs{}
	created, updated := 0, 0
	for i, idx := range targets {
		p := data.probability[idx]

		if existing, ok := matched[idx]; ok {
			if merr := uc.productRepo.UpsertProductIntegrationMapping(ctx, p.ID, businessID, uint(integIDUint), jumpsellerRefs(existing)); merr != nil {
				uc.logger.Error(ctx).Err(merr).Str("sku", p.SKU).Msg("Error al mapear producto existente de Jumpseller")
				fails.add(p.SKU)
			} else {
				updated++
			}
			uc.maybeProductProgress(ctx, businessID, uint(integIDUint), correlationID, DirectionToJumpseller, i+1, total, created, updated, fails.count())
			continue
		}

		newID, cerr := uc.client.CreateProduct(ctx, data.cred, domain.CreateProductInput{
			Name:          p.Name,
			SKU:           p.SKU,
			Price:         p.Price,
			Description:   p.Description,
			StockQuantity: p.StockQuantity,
			ManageStock:   p.TrackInventory,
			Weight:        uc.probabilityWeightForStore(ctx, p, data.storeWeightUnit),
			Height:        p.Height,
			Width:         p.Width,
			Length:        p.Length,
		})
		if cerr != nil {
			uc.logger.Error(ctx).Err(cerr).Str("sku", p.SKU).Msg("Error al crear producto en Jumpseller")
			fails.add(p.SKU)
			uc.maybeProductProgress(ctx, businessID, uint(integIDUint), correlationID, DirectionToJumpseller, i+1, total, created, updated, fails.count())
			continue
		}

		newRefs := productmatch.ExternalRefs{ProductID: newID, SKU: p.SKU, Barcode: p.Barcode}
		if merr := uc.productRepo.UpsertProductIntegrationMapping(ctx, p.ID, businessID, uint(integIDUint), newRefs); merr != nil {
			uc.logger.Error(ctx).Err(merr).Str("sku", p.SKU).Msg("Producto creado en Jumpseller pero fallo el mapeo")
			fails.add(p.SKU)
		} else {
			created++
		}
		uc.maybeProductProgress(ctx, businessID, uint(integIDUint), correlationID, DirectionToJumpseller, i+1, total, created, updated, fails.count())
	}

	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "jumpseller.product.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      DirectionToJumpseller,
		"mode":           ModeCreate,
		"total":          total,
		"created":        created,
		"updated":        updated,
		"failed":         fails.count(),
		"failed_skus":    fails.list(),
		"failed_hidden":  fails.truncated(),
	})

	return nil
}

func (uc *jumpsellerUseCase) ApplyProductsToProbability(ctx context.Context, integrationID string, businessID uint, correlationID string, skus ...string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)

	data, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	only := selectedSKUs(skus)
	missing := make([]jumpsellerSKU, 0)
	for _, idx := range data.outcome.OnlyInChannel {
		j := data.jumpseller[idx]
		key := normalizeSKU(j.SKU)
		if key == "" || (only != nil && !only[key]) {
			continue
		}
		missing = append(missing, j)
	}

	return uc.publishUpserts(ctx, businessID, uint(integIDUint), correlationID, ModeCreate, missing, data.storeWeightUnit)
}

func (uc *jumpsellerUseCase) UpdateProductsToProbability(ctx context.Context, integrationID string, businessID uint, correlationID string, skus ...string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)

	data, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	only := selectedSKUs(skus)
	existing := make([]jumpsellerSKU, 0)
	for _, pair := range data.outcome.Pairs {
		j := data.jumpseller[pair.ChannelIndex]
		key := normalizeSKU(j.SKU)
		if key == "" || (only != nil && !only[key]) {
			continue
		}
		existing = append(existing, j)
	}

	return uc.publishUpserts(ctx, businessID, uint(integIDUint), correlationID, ModeUpdate, existing, data.storeWeightUnit)
}

func (uc *jumpsellerUseCase) publishUpserts(ctx context.Context, businessID, integrationID uint, correlationID, mode string, items []jumpsellerSKU, storeUnit string) error {
	total := len(items)
	uc.emitSyncEvent(ctx, businessID, integrationID, "jumpseller.product.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      DirectionToProbability,
		"mode":           mode,
		"total":          total,
	})

	if uc.rabbit == nil {
		uc.emitSyncEvent(ctx, businessID, integrationID, "jumpseller.product.sync.completed", map[string]interface{}{
			"correlation_id": correlationID,
			"direction":      DirectionToProbability,
			"mode":           mode,
			"total":          total,
			"created":        0,
			"failed":         total,
		})
		return fmt.Errorf("RabbitMQ no disponible: no se pueden sincronizar productos hacia Probability")
	}

	if derr := uc.rabbit.DeclareQueue(rabbitmq.QueueProductsProviderUpsert, true); derr != nil {
		return fmt.Errorf("declarando la cola de upsert de productos: %w", derr)
	}

	fails := &failedSKUs{}
	applied := 0
	for i, j := range items {
		msg := uc.upsertMsgFromJumpseller(ctx, businessID, integrationID, j, storeUnit)

		payload, merr := json.Marshal(msg)
		if merr != nil {
			fails.add(j.SKU)
		} else if perr := uc.rabbit.Publish(ctx, rabbitmq.QueueProductsProviderUpsert, payload); perr != nil {
			uc.logger.Error(ctx).Err(perr).Str("sku", j.SKU).Msg("Error al publicar producto hacia Probability")
			fails.add(j.SKU)
		} else {
			applied++
		}
		uc.maybeProductProgress(ctx, businessID, integrationID, correlationID, DirectionToProbability, i+1, total, applied, 0, fails.count())
	}

	uc.emitSyncEvent(ctx, businessID, integrationID, "jumpseller.product.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      DirectionToProbability,
		"mode":           mode,
		"total":          total,
		"created":        applied,
		"updated":        applied,
		"failed":         fails.count(),
		"failed_skus":    fails.list(),
		"failed_hidden":  fails.truncated(),
	})

	return nil
}

func (uc *jumpsellerUseCase) UpdateProductsToJumpseller(ctx context.Context, integrationID string, businessID uint, correlationID string, skus ...string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)

	data, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	matched := matchedByProbabilityIndex(data)

	only := selectedSKUs(skus)
	targets := make([]int, 0)
	for i, p := range data.probability {
		key := normalizeSKU(p.SKU)
		if key == "" || (only != nil && !only[key]) {
			continue
		}
		if _, ok := matched[i]; ok {
			targets = append(targets, i)
		}
	}

	total := len(targets)
	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "jumpseller.product.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      DirectionToJumpseller,
		"mode":           ModeUpdate,
		"total":          total,
	})

	fails := &failedSKUs{}
	updated, skipped := 0, 0
	touchedParents := make(map[int64]bool)

	for i, idx := range targets {
		p := data.probability[idx]
		target := matched[idx]

		if target.VariantID > 0 {
			if touchedParents[target.ProductID] {
				skipped++
				uc.maybeProductProgress(ctx, businessID, uint(integIDUint), correlationID, DirectionToJumpseller, i+1, total, 0, updated, fails.count())
				continue
			}
			touchedParents[target.ProductID] = true
		}

		price := p.Price
		input := domain.UpdateProductInput{
			Name:        p.Name,
			Price:       &price,
			Description: p.Description,
			Weight:      uc.probabilityWeightForStore(ctx, p, data.storeWeightUnit),
			Height:      p.Height,
			Width:       p.Width,
			Length:      p.Length,
		}

		if target.VariantID > 0 {
			input.Price = nil
			input.Name = ""
		}

		if uerr := uc.client.UpdateProduct(ctx, data.cred, target.ProductID, input); uerr != nil {
			uc.logger.Error(ctx).Err(uerr).
				Str("sku", p.SKU).
				Int64("product_id", target.ProductID).
				Msg("Error al actualizar producto en Jumpseller")
			fails.add(p.SKU)
		} else {
			updated++
		}
		uc.maybeProductProgress(ctx, businessID, uint(integIDUint), correlationID, DirectionToJumpseller, i+1, total, 0, updated, fails.count())
	}

	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "jumpseller.product.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      DirectionToJumpseller,
		"mode":           ModeUpdate,
		"total":          total,
		"created":        0,
		"updated":        updated,
		"skipped":        skipped,
		"failed":         fails.count(),
		"failed_skus":    fails.list(),
		"failed_hidden":  fails.truncated(),
	})

	return nil
}

func (uc *jumpsellerUseCase) AssociateProducts(ctx context.Context, integrationID string, businessID uint, correlationID string, skus []string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)

	data, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	associated, err := uc.associatedSKUs(ctx, uint(integIDUint))
	if err != nil {
		return err
	}

	var only map[string]bool
	if len(skus) > 0 {
		only = make(map[string]bool, len(skus))
		for _, sku := range skus {
			if key := normalizeSKU(sku); key != "" {
				only[key] = true
			}
		}
	}

	targets := make([]productmatch.Pair, 0, len(data.outcome.Pairs))
	for _, pair := range data.outcome.Pairs {
		key := normalizeSKU(data.probability[pair.ProbabilityIndex].SKU)
		if only != nil {
			if key == "" || !only[key] {
				continue
			}
		} else if associated[key] {
			continue
		}
		if data.jumpseller[pair.ChannelIndex].ExternalID == "" {
			continue
		}
		targets = append(targets, pair)
	}

	total := len(targets)
	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "jumpseller.product.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "associate",
		"total":          total,
	})

	fails := &failedSKUs{}
	updated := 0
	for i, pair := range targets {
		p := data.probability[pair.ProbabilityIndex]
		j := data.jumpseller[pair.ChannelIndex]

		if merr := uc.productRepo.UpsertProductIntegrationMapping(ctx, p.ID, businessID, uint(integIDUint), jumpsellerRefs(j)); merr != nil {
			uc.logger.Error(ctx).Err(merr).Str("sku", p.SKU).Msg("Error al asociar producto con Jumpseller")
			fails.add(p.SKU)
		} else {
			updated++
		}
		uc.maybeProductProgress(ctx, businessID, uint(integIDUint), correlationID, "associate", i+1, total, 0, updated, fails.count())
	}

	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "jumpseller.product.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "associate",
		"total":          total,
		"created":        0,
		"updated":        updated,
		"failed":         fails.count(),
		"failed_skus":    fails.list(),
		"failed_hidden":  fails.truncated(),
	})

	return nil
}

func (uc *jumpsellerUseCase) SyncProducts(ctx context.Context, integrationID string, businessID uint, correlationID string) error {
	if err := uc.ApplyProductsToJumpseller(ctx, integrationID, businessID, correlationID); err != nil {
		return err
	}
	return uc.ApplyProductsToProbability(ctx, integrationID, businessID, correlationID)
}

func (uc *jumpsellerUseCase) maybeProductProgress(ctx context.Context, businessID, integrationID uint, correlationID, direction string, processed, total, created, updated, failed int) {
	if processed%syncProgressBatch != 0 && processed != total {
		return
	}
	uc.emitSyncEvent(ctx, businessID, integrationID, "jumpseller.product.sync.progress", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      direction,
		"processed":      processed,
		"total":          total,
		"created":        created,
		"updated":        updated,
		"failed":         failed,
	})
}
