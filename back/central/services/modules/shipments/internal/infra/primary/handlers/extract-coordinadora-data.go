package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/app/usecaseshipment"
)

func (h *Handlers) ExtractCoordinadoraData(c *gin.Context) {
	shipmentIDStr := c.Param("id")
	shipmentID, err := strconv.ParseUint(shipmentIDStr, 10, 64)
	if err != nil || shipmentID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shipment id"})
		return
	}

	shipment, err := h.uc.Repo().GetShipmentByID(c.Request.Context(), uint(shipmentID))
	if err != nil || shipment == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "shipment not found"})
		return
	}

	if shipment.GuideURL == nil || *shipment.GuideURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "shipment has no guide URL"})
		return
	}

	metadata, err := usecaseshipment.ExtractCoordinadoraMetadata(c.Request.Context(), *shipment.GuideURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to extract metadata: %v", err)})
		return
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode metadata"})
		return
	}

	shipment.Metadata = metadataJSON
	if err := h.uc.Repo().UpdateShipment(c.Request.Context(), shipment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save metadata"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "metadata extracted and saved",
		"data":    metadata,
	})
}
