package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

func (c *OrderCreatedConsumer) maybeAutoGenerate(ctx context.Context, msg *orderCreatedMessage, quote *domain.SavedQuote, rateIdx int) {
	repo := c.uc.Repo()

	if msg.IntegrationID == nil || c.transportPub == nil {
		return
	}
	enabled, _ := repo.GetIntegrationConfigFlag(ctx, *msg.IntegrationID, autoGenerateGuideConfigKey)
	if !enabled {
		return
	}

	carrier, err := repo.GetActiveShippingCarrier(ctx, quote.BusinessID)
	if err != nil || carrier == nil {
		c.failQuote(ctx, quote, "el negocio no tiene una transportadora activa para generar la guia")
		return
	}

	idRate := quote.SelectedIDRate
	fresh := quote.ExpiresAt != nil && quote.ExpiresAt.After(time.Now())
	if idRate == nil || !fresh {
		matched, ok := c.requoteAndMatch(ctx, carrier, quote, rateIdx)
		if !ok {
			c.failQuote(ctx, quote, "la transportadora elegida en el checkout ("+quote.SelectedCarrier+") ya no esta disponible, genera la guia manualmente")
			return
		}
		idRate = matched
	}
	if idRate == nil {
		c.failQuote(ctx, quote, "no se pudo determinar la tarifa de la transportadora elegida")
		return
	}

	payload := cloneStringMap(quote.RequestPayload)
	payload["idRate"] = *idRate
	payload["carrier"] = quote.SelectedCarrier
	payload["order_uuid"] = msg.OrderID
	payload["external_order_id"] = msg.OrderID
	payload["myShipmentReference"] = msg.OrderID

	shipmentID, err := c.createOrReusePendingShipment(ctx, payload, carrier, msg.OrderID)
	if err != nil {
		c.log.Error(ctx).Err(err).Str("order_id", msg.OrderID).Msg("Failed to pre-create shipment for auto guide")
		return
	}

	effectiveBaseURL := carrier.BaseURL
	if carrier.IsTesting && carrier.BaseURLTest != "" {
		effectiveBaseURL = carrier.BaseURLTest
	}

	out := &domain.TransportRequestMessage{
		ShipmentID:        &shipmentID,
		Provider:          carrier.ProviderCode,
		IntegrationTypeID: carrier.IntegrationTypeID,
		Operation:         "generate",
		CorrelationID:     uuid.New().String(),
		BusinessID:        quote.BusinessID,
		IntegrationID:     carrier.IntegrationID,
		BaseURL:           effectiveBaseURL,
		IsTest:            carrier.IsTesting,
		Timestamp:         time.Now(),
		Payload:           payload,
	}
	if err := c.transportPub.PublishTransportRequest(ctx, out); err != nil {
		c.log.Error(ctx).Err(err).Str("order_id", msg.OrderID).Msg("Failed to publish auto generate request")
		return
	}

	quote.Status = domain.QuoteStatusGuideGenerated
	_ = repo.UpdateSavedQuote(ctx, quote)

	c.log.Info(ctx).
		Str("order_id", msg.OrderID).
		Str("carrier", quote.SelectedCarrier).
		Uint("shipment_id", shipmentID).
		Msg("Auto-generating guide with checkout carrier")
}

func (c *OrderCreatedConsumer) requoteAndMatch(ctx context.Context, carrier *domain.CarrierInfo, quote *domain.SavedQuote, rateIdx int) (*int64, bool) {
	wanted := rateAt(quote, rateIdx)
	if wanted == nil {
		return nil, false
	}
	wantCarrier := normalizeMatch(strFromAny(wanted["carrier"]))
	wantProduct := normalizeMatch(strFromAny(wanted["product"]))

	correlationID := uuid.New().String()
	effectiveBaseURL := carrier.BaseURL
	if carrier.IsTesting && carrier.BaseURLTest != "" {
		effectiveBaseURL = carrier.BaseURLTest
	}

	msg := &domain.TransportRequestMessage{
		Provider:          carrier.ProviderCode,
		IntegrationTypeID: carrier.IntegrationTypeID,
		Operation:         "quote",
		CorrelationID:     correlationID,
		BusinessID:        quote.BusinessID,
		IntegrationID:     carrier.IntegrationID,
		BaseURL:           effectiveBaseURL,
		IsTest:            carrier.IsTesting,
		Timestamp:         time.Now(),
		Payload:           quote.RequestPayload,
	}
	if err := c.transportPub.PublishTransportRequest(ctx, msg); err != nil {
		return nil, false
	}

	rates := c.pollQuoteRates(ctx, correlationID, 8*time.Second)

	var fallback *int64
	for _, r := range rates {
		if normalizeMatch(strFromAny(r["carrier"])) != wantCarrier {
			continue
		}
		idr := idRateFromMap(r)
		if idr == nil {
			continue
		}
		if normalizeMatch(strFromAny(r["product"])) == wantProduct {
			return idr, true
		}
		if fallback == nil {
			fallback = idr
		}
	}
	if fallback != nil {
		return fallback, true
	}
	return nil, false
}

func (c *OrderCreatedConsumer) pollQuoteRates(ctx context.Context, correlationID string, timeout time.Duration) []map[string]interface{} {
	if c.redisClient == nil {
		return nil
	}
	key := "shipment:quote:result:" + correlationID

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			return nil
		case <-ticker.C:
			val, err := c.redisClient.Get(ctx, key)
			if err != nil {
				continue
			}
			var result struct {
				Status string                 `json:"status"`
				Data   map[string]interface{} `json:"data"`
				Error  string                 `json:"error"`
			}
			if err := json.Unmarshal([]byte(val), &result); err != nil {
				return nil
			}
			_ = c.redisClient.Delete(ctx, key)
			if result.Status == "error" {
				return nil
			}
			return extractRatesFromData(result.Data)
		}
	}
}

func (c *OrderCreatedConsumer) createOrReusePendingShipment(ctx context.Context, payload map[string]interface{}, carrier *domain.CarrierInfo, orderUUID string) (uint, error) {
	existing, _ := c.uc.Repo().GetShipmentsByOrderID(ctx, orderUUID)
	for i := range existing {
		s := &existing[i]
		if s.Status == "cancelled" || s.Status == "failed" {
			continue
		}
		if (s.TrackingNumber != nil && *s.TrackingNumber != "") || (s.GuideURL != nil && *s.GuideURL != "") {
			return 0, fmt.Errorf("order %s already has an active guide (shipment %d)", orderUUID, s.ID)
		}
	}
	for i := range existing {
		s := existing[i]
		if s.Status == "pending" && (s.TrackingNumber == nil || *s.TrackingNumber == "") && (s.GuideURL == nil || *s.GuideURL == "") {
			return s.ID, nil
		}
	}

	req := buildGenerateShipmentReq(payload, carrier, orderUUID)
	resp, err := c.uc.CreateShipment(ctx, req)
	if err != nil {
		return 0, err
	}
	return resp.ID, nil
}

func (c *OrderCreatedConsumer) failQuote(ctx context.Context, quote *domain.SavedQuote, reason string) {
	quote.Status = domain.QuoteStatusFailed
	if err := c.uc.Repo().UpdateSavedQuote(ctx, quote); err != nil {
		c.log.Error(ctx).Err(err).Uint("quote_id", quote.ID).Msg("Failed to mark saved quote as failed")
	}
	c.log.Warn(ctx).Uint("quote_id", quote.ID).Str("reason", reason).Msg("Auto guide generation skipped")
}

func buildGenerateShipmentReq(payload map[string]interface{}, carrier *domain.CarrierInfo, orderUUID string) *domain.CreateShipmentRequest {
	provider := carrier.ProviderCode
	req := &domain.CreateShipmentRequest{
		Status:      "pending",
		CarrierCode: &provider,
		OrderID:     &orderUUID,
	}
	if v := strFromAny(payload["carrier"]); v != "" {
		req.Carrier = &v
	}
	if dest, ok := payload["destination"].(map[string]interface{}); ok {
		first := strFromAny(dest["firstName"])
		last := strFromAny(dest["lastName"])
		req.ClientName = strings.TrimSpace(first + " " + last)
		req.DestinationAddress = strFromAny(dest["address"])
		req.DestinationSuburb = strFromAny(dest["suburb"])
	}
	if pkgs, ok := payload["packages"].([]interface{}); ok && len(pkgs) > 0 {
		if pkg, ok := pkgs[0].(map[string]interface{}); ok {
			req.Weight = floatPtrFromAny(pkg["weight"])
			req.Height = floatPtrFromAny(pkg["height"])
			req.Width = floatPtrFromAny(pkg["width"])
			req.Length = floatPtrFromAny(pkg["length"])
		}
	}
	return req
}

func extractRatesFromData(data map[string]interface{}) []map[string]interface{} {
	if data == nil {
		return nil
	}
	inner, ok := data["data"].(map[string]interface{})
	if !ok {
		return nil
	}
	rawList, ok := inner["rates"].([]interface{})
	if !ok {
		return nil
	}
	out := make([]map[string]interface{}, 0, len(rawList))
	for _, raw := range rawList {
		if m, ok := raw.(map[string]interface{}); ok {
			out = append(out, m)
		}
	}
	return out
}

func cloneStringMap(in map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(in)+4)
	for k, v := range in {
		out[k] = v
	}
	return out
}

func normalizeMatch(s string) string {
	r := strings.NewReplacer(" ", "", "-", "", "_", "")
	return r.Replace(strings.ToUpper(strings.TrimSpace(s)))
}

func strFromAny(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

func floatPtrFromAny(v interface{}) *float64 {
	switch n := v.(type) {
	case float64:
		return &n
	case float32:
		f := float64(n)
		return &f
	case int:
		f := float64(n)
		return &f
	case int64:
		f := float64(n)
		return &f
	}
	return nil
}
