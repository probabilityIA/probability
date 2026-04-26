package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

func (h *Handlers) CollectCOD(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID de envio invalido",
			"error":   "invalid_id",
		})
		return
	}

	var req domain.CollectCODRequest
	if c.Request.ContentLength > 0 {
		_ = c.ShouldBindJSON(&req)
	}

	resp, err := h.uc.CollectCOD(c.Request.Context(), uint(id), req.Notes)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, domain.ErrShipmentNotFound):
			status = http.StatusNotFound
		case errors.Is(err, domain.ErrShipmentNotDelivered),
			errors.Is(err, domain.ErrOrderAlreadyPaid),
			errors.Is(err, domain.ErrOrderNotCOD),
			errors.Is(err, domain.ErrOrderIDRequired):
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cobro contra entrega registrado exitosamente",
		"data":    resp,
	})
}
