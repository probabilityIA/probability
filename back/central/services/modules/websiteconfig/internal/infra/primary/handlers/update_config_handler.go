package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/websiteconfig/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/websiteconfig/internal/infra/primary/handlers/response"
)

func (h *Handlers) UpdateConfig(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id es requerido"})
		return
	}

	var req request.UpdateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "datos invalidos"})
		return
	}

	dto := req.ToDTO()
	config, err := h.uc.UpdateConfig(c.Request.Context(), businessID, dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.ConfigFromEntity(config))
}
