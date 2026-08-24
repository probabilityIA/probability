package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/mobile/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/mobile/internal/infra/primary/handlers/response"
)

func (h *Handlers) GetOrderFull(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "business_id es requerido para super admin",
		})
		return
	}

	orderID := strings.TrimSpace(c.Param("id"))
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id de orden invalido"})
		return
	}

	result, err := h.uc.GetOrderFull(c.Request.Context(), businessID, orderID)
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrOrderNotFound), errors.Is(err, domainerrors.ErrOrderNotInScope):
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": domainerrors.ErrOrderNotFound.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response.FromOrderFull(result),
	})
}
