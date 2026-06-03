package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

func (h *Handlers) QuoteShipment(c *gin.Context) {
	// 1. Resolve business_id (JWT for normal users, order DB lookup for super admin)
	businessID, err := h.resolveBusinessIDFromOrder(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Resolve active shipping carrier
	carrier, err := h.resolveCarrier(c, businessID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Parse request body
	var raw map[string]interface{}
	if err := c.ShouldBindJSON(&raw); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, ok := raw["origin"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "origin es requerido"})
		return
	}
	if _, ok := raw["destination"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "destination es requerido"})
		return
	}
	if _, ok := raw["packages"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "packages es requerido"})
		return
	}

	correlationID := uuid.New().String()

	result, err := h.runQuote(c.Request.Context(), carrier, businessID, raw, correlationID, 30*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enviar solicitud de cotización: " + err.Error()})
		return
	}

	switch result.Status {
	case quoteStatusTimeout:
		c.JSON(http.StatusRequestTimeout, gin.H{
			"success":        false,
			"message":        "La cotización tardó demasiado. Por favor intente nuevamente.",
			"correlation_id": correlationID,
		})
	case quoteStatusAccepted:
		c.JSON(http.StatusAccepted, gin.H{
			"success":        true,
			"message":        "Solicitud de cotización enviada. Será procesada en breve.",
			"correlation_id": correlationID,
		})
	case quoteStatusError:
		c.JSON(http.StatusOK, gin.H{
			"success":        false,
			"message":        result.Error,
			"correlation_id": correlationID,
		})
	default:
		ratesList := toRatesList(getRatesFromData(result.Data))
		if len(ratesList) > 0 {
			orderRef, _ := raw["order_uuid"].(string)
			_, _ = h.uc.Quotes.SaveQuote(c.Request.Context(), domain.SaveQuoteInput{
				BusinessID:       businessID,
				IntegrationID:    carrier.IntegrationID,
				Source:           domain.QuoteSourcePanel,
				CorrelationID:    correlationID,
				ExternalOrderRef: orderRef,
				RequestPayload:   raw,
				Rates:            ratesList,
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"success":        true,
			"message":        "Cotización exitosa",
			"correlation_id": correlationID,
			"data":           gin.H{"rates": getRatesFromData(result.Data)},
		})
	}
}

const (
	quoteStatusSuccess  = "success"
	quoteStatusError    = "error"
	quoteStatusTimeout  = "timeout"
	quoteStatusAccepted = "accepted"
)

type quoteResult struct {
	Status string
	Data   map[string]interface{}
	Error  string
}

func (h *Handlers) runQuote(ctx context.Context, carrier *domain.CarrierInfo, businessID uint, payload map[string]interface{}, correlationID string, timeout time.Duration) (*quoteResult, error) {
	effectiveBaseURL := carrier.BaseURL
	if carrier.IsTesting && carrier.BaseURLTest != "" {
		effectiveBaseURL = carrier.BaseURLTest
	}

	msg := &domain.TransportRequestMessage{
		Provider:          carrier.ProviderCode,
		IntegrationTypeID: carrier.IntegrationTypeID,
		Operation:         "quote",
		CorrelationID:     correlationID,
		BusinessID:        businessID,
		IntegrationID:     carrier.IntegrationID,
		BaseURL:           effectiveBaseURL,
		IsTest:            carrier.IsTesting,
		Timestamp:         time.Now(),
		Payload:           payload,
	}

	if err := h.transportPub.PublishTransportRequest(ctx, msg); err != nil {
		return nil, err
	}

	if h.redisClient == nil {
		return &quoteResult{Status: quoteStatusAccepted}, nil
	}

	redisKey := fmt.Sprintf("shipment:quote:result:%s", correlationID)

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	timeoutTimer := time.NewTimer(timeout)
	defer timeoutTimer.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timeoutTimer.C:
			return &quoteResult{Status: quoteStatusTimeout}, nil
		case <-ticker.C:
			val, err := h.redisClient.Get(ctx, redisKey)
			if err != nil {
				continue
			}

			var result struct {
				Status string                 `json:"status"`
				Data   map[string]interface{} `json:"data"`
				Error  string                 `json:"error"`
			}
			if err := json.Unmarshal([]byte(val), &result); err != nil {
				return nil, err
			}

			h.redisClient.Delete(ctx, redisKey)

			if result.Status == quoteStatusError {
				return &quoteResult{Status: quoteStatusError, Error: result.Error}, nil
			}
			return &quoteResult{Status: quoteStatusSuccess, Data: result.Data}, nil
		}
	}
}

// getRatesFromData extracts the rates array from the transport provider response data.
// Transport response format: { "status": "success", "data": { "rates": [...] } }
func getRatesFromData(data map[string]interface{}) interface{} {
	if data == nil {
		return nil
	}
	innerData, ok := data["data"].(map[string]interface{})
	if !ok {
		return nil
	}
	rates, ok := innerData["rates"]
	if !ok {
		return nil
	}
	return rates
}
