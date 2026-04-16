package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/primary/handlers/response"
)

// ListChannelStatuses godoc
// @Summary      Listar estados de un canal de integración
// @Description  Obtiene los estados nativos de un canal de integración (ej: Shopify, MercadoLibre).
// @Tags         Channel Statuses
// @Produce      json
// @Param        integration_type_id  query     int   true   "ID del tipo de integración"
// @Param        is_active            query     bool  false  "Filtrar por activo/inactivo"
// @Success      200  {object}  response.ChannelStatusListResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /channel-statuses [get]
func (h *handler) ListChannelStatuses(c *gin.Context) {
	integrationTypeIDStr := c.Query("integration_type_id")
	if integrationTypeIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "integration_type_id es requerido",
		})
		return
	}

	integrationTypeID, err := strconv.ParseUint(integrationTypeIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "integration_type_id inválido",
		})
		return
	}

	var isActive *bool
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if val, err := strconv.ParseBool(isActiveStr); err == nil {
			isActive = &val
		}
	}

	result, err := h.uc.ListChannelStatuses(c.Request.Context(), uint(integrationTypeID), isActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener estados del canal",
			"error":   err.Error(),
		})
		return
	}

	data := make([]response.ChannelStatusResponse, len(result))
	for i, s := range result {
		cs := response.ChannelStatusResponse{
			ID:                s.ID,
			IntegrationTypeID: s.IntegrationTypeID,
			Code:              s.Code,
			Name:              s.Name,
			Description:       s.Description,
			IsActive:          s.IsActive,
			DisplayOrder:      s.DisplayOrder,
		}
		if s.IntegrationType != nil {
			it := response.IntegrationTypeResponse{
				ID:       s.IntegrationType.ID,
				Code:     s.IntegrationType.Code,
				Name:     s.IntegrationType.Name,
				ImageURL: s.IntegrationType.ImageURL,
			}
			cs.IntegrationType = &it
		}
		data[i] = cs
	}

	c.JSON(http.StatusOK, response.ChannelStatusListResponse{
		Success: true,
		Message: "Estados del canal obtenidos exitosamente",
		Data:    data,
	})
}
