package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/infra/primary/handlers/response"
)

func (h *Handlers) UpdateClientGroup(c *gin.Context) {
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

	var req request.SaveClientGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	group, err := h.uc.UpdateClientGroup(c.Request.Context(), dtos.SaveClientGroupDTO{
		ID:          groupID,
		BusinessID:  businessID,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    isActive,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.FromGroupEntity(group))
}
