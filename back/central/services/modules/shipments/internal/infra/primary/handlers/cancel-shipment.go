package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

// CancelShipment godoc
// @Summary      Cancelar envío
// @Description  Cancela un envío directamente en la API del carrier (async via cola)
// @Tags         Shipments
// @Accept       json
// @Produce      json
// @Param        id                path      string  true  "ID de Envío (shipment id o tracking)"
// @Security     BearerAuth
// @Success      202  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /shipments/:id/cancel [post]
func (h *Handlers) CancelShipment(c *gin.Context) {
	// 1. Extract shipment ID first (needed for super admin business resolution)
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID de envío es requerido",
		})
		return
	}

	// 2. Resolve business_id (JWT for normal users, shipment DB lookup for super admin)
	businessID, err := h.resolveBusinessIDFromShipment(c, id)
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
		Operation:         "cancel",
		CorrelationID:     correlationID,
		BusinessID:        businessID,
		IntegrationID:     carrier.IntegrationID,
		BaseURL:           effectiveBaseURL,
		IsTest:            carrier.IsTesting,
		Timestamp:         time.Now(),
		Payload: map[string]interface{}{
			"id_shipment": id,
		},
	}

	if err := h.transportPub.PublishTransportRequest(c.Request.Context(), msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al enviar solicitud de cancelación: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success":        true,
		"message":        "Solicitud de cancelación enviada. Será procesada en breve.",
		"correlation_id": correlationID,
	})
}
