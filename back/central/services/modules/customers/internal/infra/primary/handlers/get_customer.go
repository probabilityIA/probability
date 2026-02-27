package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/infra/primary/handlers/response"
)

func (h *Handlers) GetClient(c *gin.Context) {
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

	client, err := h.uc.GetClient(c.Request.Context(), businessID, uint(clientID))
	if err != nil {
		if errors.Is(err, domainerrors.ErrClientNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.DetailFromEntity(client))
}
