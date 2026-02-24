package handlers

import (
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

	c.JSON(http.StatusAccepted, gin.H{
		"success":        true,
		"message":        "Solicitud de cotización enviada. Será procesada en breve.",
		"correlation_id": correlationID,
	})
}
