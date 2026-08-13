package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
)

type dataApplyRequest struct {
	IntegrationID uint     `json:"integration_id" binding:"required"`
	Field         string   `json:"field" binding:"required"`
	Mode          string   `json:"mode"`
	SKUs          []string `json:"skus"`
}

func (h *Handlers) ApplyChannelData(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		h.respondBusinessIDRequired(c)
		return
	}

	var req dataApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "integration_id y field son requeridos",
			"error":   err.Error(),
		})
		return
	}

	userID, _ := middleware.GetUserID(c)

	result, err := h.uc.ApplyChannelData(c.Request.Context(), domain.DataApplyRequest{
		BusinessID:    businessID,
		IntegrationID: req.IntegrationID,
		Field:         strings.TrimSpace(req.Field),
		Mode:          strings.TrimSpace(req.Mode),
		SKUs:          req.SKUs,
		UserID:        userID,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidProductData) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Datos invalidos para aplicar",
				"error":   err.Error(),
			})
			return
		}
		h.log.Error(c.Request.Context()).Err(err).
			Uint("business_id", businessID).Str("field", req.Field).
			Msg("Error al traer datos del canal a Probability")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al aplicar los datos del canal",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": "Actualizacion iniciada",
		"data":    result,
	})
}

type dataUndoRequest struct {
	BatchID string `json:"batch_id" binding:"required"`
}

func (h *Handlers) UndoChannelData(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		h.respondBusinessIDRequired(c)
		return
	}

	var req dataUndoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "batch_id es requerido",
			"error":   err.Error(),
		})
		return
	}

	userID, _ := middleware.GetUserID(c)

	result, err := h.uc.UndoChannelData(c.Request.Context(), businessID, req.BatchID, userID)
	if err != nil {
		h.log.Error(c.Request.Context()).Err(err).
			Uint("business_id", businessID).Str("batch_id", req.BatchID).
			Msg("Error al deshacer la aplicacion de datos")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al deshacer",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cambios revertidos",
		"data":    result,
	})
}

func (h *Handlers) ListDataBatches(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		h.respondBusinessIDRequired(c)
		return
	}

	limit, _ := strconv.Atoi(c.Query("limit"))

	batches, err := h.uc.ListDataBatches(c.Request.Context(), businessID, limit)
	if err != nil {
		h.log.Error(c.Request.Context()).Err(err).Uint("business_id", businessID).
			Msg("Error al listar los lotes de datos aplicados")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al listar los lotes",
			"error":   err.Error(),
		})
		return
	}

	if batches == nil {
		batches = []domain.DataBatch{}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": batches})
}
