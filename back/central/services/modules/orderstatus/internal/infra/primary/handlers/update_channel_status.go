package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/primary/handlers/response"
)

func (h *handler) UpdateChannelStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID inválido",
		})
		return
	}

	var req request.UpdateChannelStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Datos inválidos",
			"error":   err.Error(),
		})
		return
	}

	status := &entities.ChannelStatusInfo{
		Code:         req.Code,
		Name:         req.Name,
		Description:  req.Description,
		IsActive:     req.IsActive,
		DisplayOrder: req.DisplayOrder,
	}

	result, err := h.uc.UpdateChannelStatus(c.Request.Context(), uint(id), status)
	if err != nil {
		if err == domainerrors.ErrChannelStatusNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Estado del canal no encontrado",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al actualizar estado del canal",
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

	c.JSON(http.StatusOK, response.ChannelStatusSingleResponse{
		Success: true,
		Message: "Estado del canal actualizado exitosamente",
		Data:    data,
	})
}
