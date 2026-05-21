package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/entities"
)

func (h *Handlers) CarrierConfigs(c *gin.Context) {
	businessID, err := resolveBusinessID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	configs, err := h.uc.CarrierConfigs(c.Request.Context(), businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener la configuracion de transportadoras",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Configuracion de transportadoras obtenida exitosamente",
		"data":    mapCarrierConfigs(configs),
	})
}

type saveCarrierConfigRequest struct {
	CarrierName        string  `json:"carrier_name"`
	DiscountPercentage float64 `json:"discount_percentage"`
	IsActive           *bool   `json:"is_active"`
}

func (h *Handlers) SaveCarrierConfig(c *gin.Context) {
	if !isAdminUser(c) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Solo un administrador puede configurar transportadoras",
		})
		return
	}

	businessID, err := resolveBusinessID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	var req saveCarrierConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Cuerpo de la solicitud invalido"})
		return
	}
	if strings.TrimSpace(req.CarrierName) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "El nombre de la transportadora es requerido"})
		return
	}
	if req.DiscountPercentage < 0 || req.DiscountPercentage > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "El porcentaje debe estar entre 0 y 100"})
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	cfg, err := h.uc.SaveCarrierConfig(c.Request.Context(), dtos.SaveCarrierConfigDTO{
		BusinessID:         businessID,
		CarrierName:        req.CarrierName,
		DiscountPercentage: req.DiscountPercentage,
		IsActive:           isActive,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al guardar la configuracion de transportadora",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Configuracion de transportadora guardada exitosamente",
		"data":    mapCarrierConfigs([]entities.CarrierConfig{*cfg})[0],
	})
}
