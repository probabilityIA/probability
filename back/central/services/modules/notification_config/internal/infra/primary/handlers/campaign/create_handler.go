package campaign

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
)

func (h *handler) Create(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido para super admin"})
		return
	}

	var body dtos.CreateCampaignDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	body.BusinessID = businessID
	if userID := c.GetUint("user_id"); userID > 0 {
		body.CreatedBy = &userID
	}

	campaign, err := h.useCase.Create(c.Request.Context(), body)
	if err != nil {
		h.logger.Error().Err(err).Str("name", body.Name).Msg("Error creando campana")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": campaign})
}
