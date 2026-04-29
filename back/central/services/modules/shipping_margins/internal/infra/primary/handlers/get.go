package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/infra/primary/handlers/response"
)

func (h *Handlers) Get(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	m, err := h.uc.Get(c.Request.Context(), businessID, uint(id))
	if err != nil {
		if errors.Is(err, domainerrors.ErrShippingMarginNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.FromEntity(m))
}
