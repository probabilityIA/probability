package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/infra/primary/handlers/response"
)

func (h *Handlers) GetClientGroup(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	groupID, ok := parseUintParam(c, "id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	group, err := h.uc.GetClientGroup(c.Request.Context(), businessID, groupID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.FromGroupEntity(group))
}
