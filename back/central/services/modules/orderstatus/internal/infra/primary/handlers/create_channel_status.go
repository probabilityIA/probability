package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/primary/handlers/response"
)

func (h *handler) CreateChannelStatus(c *gin.Context) {
	var req request.CreateChannelStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Datos inválidos",
			"error":   err.Error(),
		})
		return
	}

	status := &entities.ChannelStatusInfo{
		IntegrationTypeID: req.IntegrationTypeID,
		Code:              req.Code,
		Name:              req.Name,
		Description:       req.Description,
		IsActive:          req.IsActive,
		DisplayOrder:      req.DisplayOrder,
	}

	result, err := h.uc.CreateChannelStatus(c.Request.Context(), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al crear estado del canal",
			"error":   err.Error(),
		})
		return
	}

	data := response.ChannelStatusResponse{
		ID:                result.ID,
		IntegrationTypeID: result.IntegrationTypeID,
		Code:              result.Code,
		Name:              result.Name,
		Description:       result.Description,
		IsActive:          result.IsActive,
		DisplayOrder:      result.DisplayOrder,
	}
	if result.IntegrationType != nil {
		it := response.IntegrationTypeResponse{
			ID:       result.IntegrationType.ID,
			Code:     result.IntegrationType.Code,
			Name:     result.IntegrationType.Name,
			ImageURL: result.IntegrationType.ImageURL,
		}
		data.IntegrationType = &it
	}

	c.JSON(http.StatusCreated, response.ChannelStatusSingleResponse{
		Success: true,
		Message: "Estado del canal creado exitosamente",
		Data:    data,
	})
}
