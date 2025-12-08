package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

// GetShipmentsByOrderID godoc
// @Summary      Obtener envíos por ID de orden
// @Description  Obtiene todos los envíos asociados a una orden específica
// @Tags         Shipments
// @Accept       json
// @Produce      json
// @Param        order_id   path      string  true  "ID de la orden (UUID)"
// @Security     BearerAuth
// @Success      200  {array}  domain.ShipmentResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /shipments/order/{order_id} [get]
func (h *Handlers) GetShipmentsByOrderID(c *gin.Context) {
	orderID := c.Param("order_id")

	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Order ID es requerido",
			"error":   "El order_id es requerido",
		})
		return
	}

	// Llamar al caso de uso
	shipments, err := h.uc.GetShipmentsByOrderID(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener envíos",
			"error":   err.Error(),
		})
		return
	}

	// Convertir a respuestas
	shipmentResponses := make([]domain.ShipmentResponse, len(shipments))
	for i, shipment := range shipments {
		shipmentResponses[i] = domain.ShipmentResponse{
			ID:         shipment.ID,
			CreatedAt:  shipment.CreatedAt,
			UpdatedAt:  shipment.UpdatedAt,
			DeletedAt:  shipment.DeletedAt,
			OrderID:    shipment.OrderID,
			TrackingNumber: shipment.TrackingNumber,
			TrackingURL:    shipment.TrackingURL,
			Carrier:        shipment.Carrier,
			CarrierCode:    shipment.CarrierCode,
			GuideID:        shipment.GuideID,
			GuideURL:       shipment.GuideURL,
			Status:         shipment.Status,
			ShippedAt:      shipment.ShippedAt,
			DeliveredAt:    shipment.DeliveredAt,
			ShippingAddressID: shipment.ShippingAddressID,
			ShippingCost:   shipment.ShippingCost,
			InsuranceCost:  shipment.InsuranceCost,
			TotalCost:      shipment.TotalCost,
			Weight:         shipment.Weight,
			Height:         shipment.Height,
			Width:          shipment.Width,
			Length:         shipment.Length,
			WarehouseID:    shipment.WarehouseID,
			WarehouseName:  shipment.WarehouseName,
			DriverID:       shipment.DriverID,
			DriverName:     shipment.DriverName,
			IsLastMile:    shipment.IsLastMile,
			EstimatedDelivery: shipment.EstimatedDelivery,
			DeliveryNotes:     shipment.DeliveryNotes,
			Metadata:          shipment.Metadata,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Envíos obtenidos exitosamente",
		"data":    shipmentResponses,
		"total":   len(shipmentResponses),
	})
}

