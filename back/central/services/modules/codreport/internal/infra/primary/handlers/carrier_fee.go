package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/dtos"
)

type updateCarrierFeeRequest struct {
	ShipmentID    uint     `json:"shipment_id"`
	CodCarrierFee *float64 `json:"cod_carrier_fee"`
}

func (h *Handlers) UpdateCarrierFee(c *gin.Context) {
	if !middleware.IsSuperAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Solo un super administrador puede corregir la comision de la transportadora",
		})
		return
	}

	businessID, err := resolveBusinessID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	var req updateCarrierFeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Cuerpo de la solicitud invalido"})
		return
	}
	if req.ShipmentID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "shipment_id es requerido"})
		return
	}
	if req.CodCarrierFee == nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "cod_carrier_fee es requerido"})
		return
	}

	userID, _ := middleware.GetUserID(c)

	result, err := h.uc.UpdateCarrierFee(c.Request.Context(), dtos.UpdateCarrierFeeDTO{
		BusinessID: businessID,
		ShipmentID: req.ShipmentID,
		Fee:        *req.CodCarrierFee,
		UserID:     userID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "No se pudo actualizar la comision de la transportadora",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Comision de la transportadora actualizada exitosamente",
		"data": gin.H{
			"shipment_id":     req.ShipmentID,
			"order_number":    result.OrderNumber,
			"previous_fee":    result.PreviousFee,
			"cod_carrier_fee": result.NewFee,
		},
	})
}
