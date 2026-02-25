package handlers

import (
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
		Payload:           raw,
	}

	if err := h.transportPub.PublishTransportRequest(c.Request.Context(), msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enviar solicitud de cotización: " + err.Error()})
		return
	}

	// If Redis is available, poll synchronously for the quote result
	if h.redisClient != nil {
		redisKey := fmt.Sprintf("shipment:quote:result:%s", correlationID)

		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()
		timeoutTimer := time.NewTimer(30 * time.Second)
		defer timeoutTimer.Stop()

		for {
			select {
			case <-c.Request.Context().Done():
				return
			case <-timeoutTimer.C:
				c.JSON(http.StatusRequestTimeout, gin.H{
					"success":        false,
					"message":        "La cotización tardó demasiado. Por favor intente nuevamente.",
					"correlation_id": correlationID,
				})
				return
			case <-ticker.C:
				val, err := h.redisClient.Get(c.Request.Context(), redisKey)
				if err != nil {
					continue // Key not yet available
				}

				var result struct {
					Status string                 `json:"status"`
					Data   map[string]interface{} `json:"data"`
					Error  string                 `json:"error"`
				}
				if err := json.Unmarshal([]byte(val), &result); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"success": false,
						"error":   "Error al procesar resultado de cotización",
					})
					return
				}

				h.redisClient.Delete(c.Request.Context(), redisKey)

				if result.Status == "error" {
					c.JSON(http.StatusOK, gin.H{
						"success":        false,
						"message":        result.Error,
						"correlation_id": correlationID,
					})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"success":        true,
					"message":        "Cotización exitosa",
					"correlation_id": correlationID,
					"data":           gin.H{"rates": extractRatesFromData(result.Data)},
				})
				return
			}
		}
	}

	// Fallback: async response when Redis is unavailable
	c.JSON(http.StatusAccepted, gin.H{
		"success":        true,
		"message":        "Solicitud de cotización enviada. Será procesada en breve.",
		"correlation_id": correlationID,
	})
}

// extractRatesFromData extracts the rates array from the transport provider response data.
// EnvioClick response format: { "status": "success", "data": { "rates": [...] } }
func extractRatesFromData(data map[string]interface{}) interface{} {
	if data == nil {
		return nil
	}
	if innerData, ok := data["data"].(map[string]interface{}); ok {
		if rates, ok := innerData["rates"]; ok {
			return rates
		}
	}
	return nil
}
