package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/response"
)

// SetAsDefaultHandler marca una integración como default
//
//	@Summary		Marcar como default
//	@Description	Marca una integración como default para su tipo
//	@Tags			Integrations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int		true	"ID de la integración"
//	@Success		200	{object}	response.IntegrationMessageResponse
//	@Failure		400	{object}	response.IntegrationErrorResponse
//	@Failure		404	{object}	response.IntegrationErrorResponse
//	@Failure		500	{object}	response.IntegrationErrorResponse
//	@Router			/integrations/{id}/set-default [put]
func (h *IntegrationHandler) SetAsDefaultHandler(c *gin.Context) {
	if !middleware.IsSuperAdmin(c) {
		c.JSON(http.StatusForbidden, response.IntegrationErrorResponse{
			Success: false,
			Message: "Solo los super usuarios pueden marcar integraciones como default",
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "ID inválido",
		})
		return
	}

	if err := h.usecase.SetAsDefault(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, response.IntegrationErrorResponse{
			Success: false,
			Message: "Error al marcar integración como default",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.IntegrationMessageResponse{
		Success: true,
		Message: "Integración marcada como default exitosamente",
	})
}
