package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/errors"
)

func (h *Handlers) DeleteClient(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	clientID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || clientID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client id"})
		return
	}

	if err := h.uc.DeleteClient(c.Request.Context(), businessID, uint(clientID)); err != nil {
		if errors.Is(err, domainerrors.ErrClientNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domainerrors.ErrClientHasOrders) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "client deleted successfully"})
}
