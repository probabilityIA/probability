package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

type cancelBatchRequestPayload struct {
	Orders []cancelBatchOrderPayload `json:"orders"`
}

type cancelBatchOrderPayload struct {
	TrackingCode string `json:"trackingCode"`
	Motivo       string `json:"motivo"`
}

// CancelBatchShipments godoc
// @Summary      Cancelar múltiples envíos
// @Description  Cancela un lote de envíos en la API del carrier (async via cola)
// @Tags         Shipments
// @Accept       json
// @Produce      json
// @Param        body              body      cancelBatchRequestPayload  true  "Lista de envíos a cancelar"
// @Security     BearerAuth
// @Success      202  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /shipments/cancel-batch [post]
func (h *Handlers) CancelBatchShipments(c *gin.Context) {
	var req cancelBatchRequestPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Payload inválido",
		})
		return
	}

	if len(req.Orders) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Se requiere al menos un envío para cancelar",
		})
		return
	}

	// Resolvemos business_id del primer envío como punto de anclaje
	firstID := req.Orders[0].TrackingCode
	businessID, err := h.resolveBusinessIDFromShipment(c, firstID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error resolviendo el primer envío: " + err.Error()})
		return
	}

	// Resolvemos carrier activo
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

	// Convertimos el payload local al generic dict para RabbitMQ
	payloadData := make([]map[string]interface{}, len(req.Orders))
	for i, o := range req.Orders {
		payloadData[i] = map[string]interface{}{
			"trackingCode": o.TrackingCode,
			"motivo":       o.Motivo,
		}
	}

	msg := &domain.TransportRequestMessage{
		Provider:          carrier.ProviderCode,
		IntegrationTypeID: carrier.IntegrationTypeID,
		Operation:         "cancel_batch",
		CorrelationID:     correlationID,
		BusinessID:        businessID,
		IntegrationID:     carrier.IntegrationID,
		BaseURL:           effectiveBaseURL,
		IsTest:            carrier.IsTesting,
		Timestamp:         time.Now(),
		Payload: map[string]interface{}{
			"orders": payloadData,
		},
	}

	if err := h.transportPub.PublishTransportRequest(c.Request.Context(), msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al enviar solicitud de cancelación masiva: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success":        true,
		"message":        "Solicitud de cancelación masiva enviada. Será procesada en breve.",
		"correlation_id": correlationID,
	})
}
