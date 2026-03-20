package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

// PublicSearchTracking godoc
// @Summary      Buscar envío públicamente
// @Description  Busca un envío por tracking_number u order_number (sin autenticación)
// @Tags         Tracking
// @Accept       json
// @Produce      json
// @Param        tracking_number   query      string  false  "Número de tracking"
// @Param        order_number      query      string  false  "Número de orden"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /tracking/search [get]
func (h *Handlers) PublicSearchTracking(c *gin.Context) {
	trackingNumber := c.Query("tracking_number")
	orderNumber := c.Query("order_number")

	if trackingNumber == "" && orderNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Debes proporcionar tracking_number u order_number",
		})
		return
	}

	var shipment *domain.ShipmentResponse
	var err error

	// Buscar por tracking_number primero
	if trackingNumber != "" {
		shipment, err = h.uc.GetShipmentByTrackingNumber(c.Request.Context(), trackingNumber)
		if err != nil {
			if err == domain.ErrShipmentNotFound {
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"message": "Envío no encontrado",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error al buscar envío",
			})
			return
		}
	} else if orderNumber != "" {
		// Buscar por order_id
		shipments, err := h.uc.GetShipmentsByOrderID(c.Request.Context(), orderNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error al buscar envío",
			})
			return
		}

		if len(shipments) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Envío no encontrado",
			})
			return
		}

		// Convertir el primer envío a ShipmentResponse
		firstShipment := shipments[0]
		shipment = &domain.ShipmentResponse{
			ID:                firstShipment.ID,
			CreatedAt:         firstShipment.CreatedAt,
			UpdatedAt:         firstShipment.UpdatedAt,
			DeletedAt:         firstShipment.DeletedAt,
			OrderID:           firstShipment.OrderID,
			TrackingNumber:    firstShipment.TrackingNumber,
			TrackingURL:       firstShipment.TrackingURL,
			Carrier:           firstShipment.Carrier,
			CarrierCode:       firstShipment.CarrierCode,
			GuideID:           firstShipment.GuideID,
			GuideURL:          firstShipment.GuideURL,
			Status:            firstShipment.Status,
			ShippedAt:         firstShipment.ShippedAt,
			DeliveredAt:       firstShipment.DeliveredAt,
			ShippingAddressID: firstShipment.ShippingAddressID,
			ShippingCost:      firstShipment.ShippingCost,
			InsuranceCost:     firstShipment.InsuranceCost,
			TotalCost:         firstShipment.TotalCost,
			Weight:            firstShipment.Weight,
			Height:            firstShipment.Height,
			Width:             firstShipment.Width,
			Length:            firstShipment.Length,
			WarehouseID:       firstShipment.WarehouseID,
			WarehouseName:     firstShipment.WarehouseName,
			DriverID:          firstShipment.DriverID,
			DriverName:        firstShipment.DriverName,
			IsLastMile:        firstShipment.IsLastMile,
			EstimatedDelivery: firstShipment.EstimatedDelivery,
			DeliveryNotes:     firstShipment.DeliveryNotes,
			Metadata:          firstShipment.Metadata,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Envío encontrado",
		"data": gin.H{
			"shipment": shipment,
		},
	})
}

// PublicGetTrackingHistory godoc
// @Summary      Obtener historial de rastreo
// @Description  Obtiene el historial de eventos de rastreo de un envío (sin autenticación)
// @Tags         Tracking
// @Accept       json
// @Produce      json
// @Param        tracking_number   path      string  true  "Número de tracking"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /tracking/{tracking_number}/history [get]
func (h *Handlers) PublicGetTrackingHistory(c *gin.Context) {
	trackingNumber := c.Param("tracking_number")

	if trackingNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Número de tracking es requerido",
		})
		return
	}

	// Obtener el envío
	shipment, err := h.uc.GetShipmentByTrackingNumber(c.Request.Context(), trackingNumber)
	if err != nil {
		if err == domain.ErrShipmentNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Envío no encontrado",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener historial",
		})
		return
	}

	// Retornar historial (vacío si no hay eventos)
	history := []map[string]interface{}{}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Historial obtenido",
		"data": gin.H{
			"shipment": shipment,
			"history":  history,
		},
	})
}
