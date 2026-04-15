package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

func (h *Handlers) EnvioClickWebhook(c *gin.Context) {
	ctx := c.Request.Context()

	var payload domain.EnvioClickWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Payload inválido: " + err.Error(),
		})
		return
	}

	if payload.TrackingCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "trackingCode es requerido",
		})
		return
	}

	// Buscar el envío por número de tracking
	shipmentResp, err := h.uc.GetShipmentByTrackingNumber(ctx, payload.TrackingCode)
	if (err != nil || shipmentResp == nil) && payload.MyShipmentReference != "" {
		// Intentar por myShipmentReference como fallback
		shipmentResp, err = h.uc.GetShipmentByTrackingNumber(ctx, payload.MyShipmentReference)
	}
	if err != nil || shipmentResp == nil {
		// No encontrado: respondemos 200 para que EnvioClick no reintente
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Envío no encontrado en el sistema, ignorado",
		})
		return
	}

	var newStatus string
	var hasIncidence bool

	if len(payload.Events) > 0 {
		latestEvent := payload.Events[len(payload.Events)-1]
		hasIncidence = latestEvent.Incidence
		newStatus = domain.MapEnvioClickStatus(latestEvent.StatusStep, hasIncidence)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Sin eventos en el payload",
		})
		return
	}

	// Construir la solicitud de actualización
	updateReq := &domain.UpdateShipmentRequest{
		Status: &newStatus,
	}

	// Si tenemos fecha de entrega real, parsearla y guardarla
	if payload.RealDeliveryDate != "" && newStatus == "delivered" {
		for _, layout := range []string{"2006-01-02 15:04:05", "2006-01-02T15:04:05Z", "2006-01-02"} {
			if t, err := time.Parse(layout, payload.RealDeliveryDate); err == nil {
				updateReq.DeliveredAt = &t
				break
			}
		}
	}

	// Si tenemos fecha de recolección, parsearla
	if payload.RealPickupDate != "" && (newStatus == "in_transit" || newStatus == "delivered") {
		for _, layout := range []string{"2006-01-02 15:04:05", "2006-01-02T15:04:05Z", "2006-01-02"} {
			if t, err := time.Parse(layout, payload.RealPickupDate); err == nil {
				updateReq.ShippedAt = &t
				break
			}
		}
	}

	// Actualizar el envío
	if _, err := h.uc.UpdateShipment(ctx, shipmentResp.ID, updateReq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Error actualizando envío: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"message":       "Estado actualizado correctamente",
		"tracking":      payload.TrackingCode,
		"new_status":    newStatus,
		"has_incidence": hasIncidence,
	})
}
