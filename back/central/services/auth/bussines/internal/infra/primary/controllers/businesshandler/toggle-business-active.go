package businesshandler

import (
	"net/http"
	"strconv"

	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/gin-gonic/gin"
)

// ActivateBusinessHandler activa un business
//
//	@Summary		Activar business
//	@Description	Activa un business por su ID
//	@Tags			Businesses
//	@Produce		json
//	@Param			id	path		int						true	"ID del business"	minimum(1)
//	@Success		200	{object}	map[string]interface{}	"Business activado exitosamente"
//	@Failure		400	{object}	map[string]interface{}	"Parámetros inválidos"
//	@Failure		404	{object}	map[string]interface{}	"Business no encontrado"
//	@Failure		500	{object}	map[string]interface{}	"Error interno del servidor"
//	@Router			/businesses/{id}/activate [put]
//	@Security		BearerAuth
func (h *BusinessHandler) ActivateBusinessHandler(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "ActivateBusinessHandler")

	idStr := c.Param("id")
	idUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error al parsear id del business")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parámetros inválidos",
			"error":   "El id debe ser un número válido",
		})
		return
	}

	if err := h.usecase.ToggleBusinessActive(ctx, uint(idUint64), true); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "business no encontrado" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"success": false, "message": err.Error()})
		return
	}

	h.logger.Info(ctx).Uint("business_id", uint(idUint64)).Msg("Business activado exitosamente")
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Business activado exitosamente"})
}

// DeactivateBusinessHandler desactiva un business
//
//	@Summary		Desactivar business
//	@Description	Desactiva un business por su ID
//	@Tags			Businesses
//	@Produce		json
//	@Param			id	path		int						true	"ID del business"	minimum(1)
//	@Success		200	{object}	map[string]interface{}	"Business desactivado exitosamente"
//	@Failure		400	{object}	map[string]interface{}	"Parámetros inválidos"
//	@Failure		404	{object}	map[string]interface{}	"Business no encontrado"
//	@Failure		500	{object}	map[string]interface{}	"Error interno del servidor"
//	@Router			/businesses/{id}/deactivate [put]
//	@Security		BearerAuth
func (h *BusinessHandler) DeactivateBusinessHandler(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "DeactivateBusinessHandler")

	idStr := c.Param("id")
	idUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error al parsear id del business")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parámetros inválidos",
			"error":   "El id debe ser un número válido",
		})
		return
	}

	if err := h.usecase.ToggleBusinessActive(ctx, uint(idUint64), false); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "business no encontrado" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"success": false, "message": err.Error()})
		return
	}

	h.logger.Info(ctx).Uint("business_id", uint(idUint64)).Msg("Business desactivado exitosamente")
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Business desactivado exitosamente"})
}
