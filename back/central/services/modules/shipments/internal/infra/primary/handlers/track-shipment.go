package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

// TrackShipment godoc
// @Summary      Rastrear envío
// @Description  Envía solicitud de rastreo a la cola de transporte (async)
// @Tags         Shipments
// @Accept       json
// @Produce      json
// @Param        tracking_number   path      string  true  "Número de tracking"
// @Security     BearerAuth
// @Success      202  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /shipments/tracking/{tracking_number}/track [post]
func (h *Handlers) TrackShipment(c *gin.Context) {
	// 1. Extract tracking number first (needed for super admin business resolution)
	trackingNumber := c.Param("tracking_number")
	if trackingNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Número de tracking es requerido",
		})
		return
	}

	// 2. Resolve business_id (JWT for normal users, shipment DB lookup for super admin)
	businessID, err := h.resolveBusinessIDFromShipment(c, trackingNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Resolve active shipping carrier
	carrier, err := h.resolveCarrier(c, businessID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		Operation:         "track",
		CorrelationID:     correlationID,
		BusinessID:        businessID,
		IntegrationID:     carrier.IntegrationID,
		BaseURL:           effectiveBaseURL,
		IsTest:            carrier.IsTesting,
		Timestamp:         time.Now(),
		Payload: map[string]interface{}{
			"tracking_number": trackingNumber,
		},
	}

	if err := h.transportPub.PublishTransportRequest(c.Request.Context(), msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al enviar solicitud de rastreo: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success":        true,
		"message":        "Solicitud de rastreo enviada. Será procesada en breve.",
		"correlation_id": correlationID,
	})
}
