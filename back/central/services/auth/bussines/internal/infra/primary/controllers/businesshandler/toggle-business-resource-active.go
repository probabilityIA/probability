package businesshandler

import (
	"net/http"
	"strconv"

	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/gin-gonic/gin"
)

// ActivateBusinessResourceHandler activa un recurso para un business específico
//
//	@Summary		Activar recurso de business
//	@Description	Activa un recurso configurado para un business específico
//	@Tags			Business Resources
//	@Accept			json
//	@Produce		json
//	@Param			resource_id	path		int						true	"ID del recurso"							minimum(1)
//	@Param			business_id	query		int						false	"ID del business (solo super admin)"			minimum(1)
//	@Success		200			{object}	map[string]interface{}	"Recurso activado exitosamente"
//	@Failure		400			{object}	map[string]interface{}	"Parámetros inválidos"
//	@Failure		401			{object}	map[string]interface{}	"No autorizado"
//	@Failure		403			{object}	map[string]interface{}	"Sin permisos"
//	@Failure		404			{object}	map[string]interface{}	"Business o recurso no encontrado"
//	@Failure		500			{object}	map[string]interface{}	"Error interno del servidor"
//	@Router			/businesses/configured-resources/{resource_id}/activate [put]
//	@Security		BearerAuth
func (h *BusinessHandler) ActivateBusinessResourceHandler(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "ActivateBusinessResourceHandler")

	// Parsear resource_id
	resourceIDStr := c.Param("resource_id")
	resourceIDUint64, err := strconv.ParseUint(resourceIDStr, 10, 32)
	if err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error al parsear resource_id")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parámetros inválidos",
			"error":   "El resource_id debe ser un número válido",
		})
		return
	}
	resourceID := uint(resourceIDUint64)

	// Aplicar lógica de super admin
	isSuperAdmin := middleware.IsSuperAdmin(c)
	var businessID uint

	if !isSuperAdmin {
		// Usuario normal: forzar business_id del token
		tokenBusinessID, ok := middleware.GetBusinessID(c)
		if !ok || tokenBusinessID == 0 {
			h.logger.Error(ctx).Msg("Usuario normal sin business_id en token")
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Sin permisos para activar recursos",
				"error":   "Usuario sin business asignado",
			})
			return
		}
		businessID = tokenBusinessID
		h.logger.Info(ctx).Uint("business_id", businessID).Msg("Usuario normal: activando recurso de su business")
	} else {
		// Super admin: permitir business_id del query
		if businessIDStr := c.Query("business_id"); businessIDStr != "" {
			if id, err := strconv.ParseUint(businessIDStr, 10, 32); err == nil {
				businessID = uint(id)
				h.logger.Info(ctx).Uint("business_id", businessID).Msg("Super admin: activando recurso del business del query")
			} else {
				h.logger.Error(ctx).Err(err).Msg("Error al parsear business_id del query")
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Parámetros inválidos",
					"error":   "El business_id debe ser un número válido",
				})
				return
			}
		} else {
			h.logger.Error(ctx).Msg("Super admin sin business_id en query")
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Parámetro requerido",
				"error":   "El business_id es requerido en el query param",
			})
			return
		}
	}

	// Activar el recurso
	if err := h.usecase.ToggleBusinessResourceActive(ctx, businessID, resourceID, true); err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error al activar recurso")

		statusCode := http.StatusInternalServerError
		errorMessage := "Error al activar el recurso"

		if err.Error() == "business no encontrado" || err.Error() == "recurso no encontrado" {
			statusCode = http.StatusNotFound
			errorMessage = err.Error()
		} else if err.Error() == "el recurso no está permitido para este tipo de business" {
			statusCode = http.StatusBadRequest
			errorMessage = err.Error()
		} else if err.Error() == "la relación entre business y recurso no existe" {
			statusCode = http.StatusNotFound
			errorMessage = err.Error()
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": errorMessage,
			"error":   err.Error(),
		})
		return
	}

	h.logger.Info(ctx).Uint("business_id", businessID).Uint("resource_id", resourceID).Msg("Recurso activado exitosamente")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Recurso activado exitosamente",
	})
}

// DeactivateBusinessResourceHandler desactiva un recurso para un business específico
//
//	@Summary		Desactivar recurso de business
//	@Description	Desactiva un recurso configurado para un business específico
//	@Tags			Business Resources
//	@Accept			json
//	@Produce		json
//	@Param			resource_id	path		int						true	"ID del recurso"							minimum(1)
//	@Param			business_id	query		int						false	"ID del business (solo super admin)"			minimum(1)
//	@Success		200			{object}	map[string]interface{}	"Recurso desactivado exitosamente"
//	@Failure		400			{object}	map[string]interface{}	"Parámetros inválidos"
//	@Failure		401			{object}	map[string]interface{}	"No autorizado"
//	@Failure		403			{object}	map[string]interface{}	"Sin permisos"
//	@Failure		404			{object}	map[string]interface{}	"Business o recurso no encontrado"
//	@Failure		500			{object}	map[string]interface{}	"Error interno del servidor"
//	@Router			/businesses/configured-resources/{resource_id}/deactivate [put]
//	@Security		BearerAuth
func (h *BusinessHandler) DeactivateBusinessResourceHandler(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "DeactivateBusinessResourceHandler")

	// Parsear resource_id
	resourceIDStr := c.Param("resource_id")
	resourceIDUint64, err := strconv.ParseUint(resourceIDStr, 10, 32)
	if err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error al parsear resource_id")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parámetros inválidos",
			"error":   "El resource_id debe ser un número válido",
		})
		return
	}
	resourceID := uint(resourceIDUint64)

	// Aplicar lógica de super admin
	isSuperAdmin := middleware.IsSuperAdmin(c)
	var businessID uint

	if !isSuperAdmin {
		// Usuario normal: forzar business_id del token
		tokenBusinessID, ok := middleware.GetBusinessID(c)
		if !ok || tokenBusinessID == 0 {
			h.logger.Error(ctx).Msg("Usuario normal sin business_id en token")
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Sin permisos para desactivar recursos",
				"error":   "Usuario sin business asignado",
			})
			return
		}
		businessID = tokenBusinessID
		h.logger.Info(ctx).Uint("business_id", businessID).Msg("Usuario normal: desactivando recurso de su business")
	} else {
		// Super admin: permitir business_id del query
		if businessIDStr := c.Query("business_id"); businessIDStr != "" {
			if id, err := strconv.ParseUint(businessIDStr, 10, 32); err == nil {
				businessID = uint(id)
				h.logger.Info(ctx).Uint("business_id", businessID).Msg("Super admin: desactivando recurso del business del query")
			} else {
				h.logger.Error(ctx).Err(err).Msg("Error al parsear business_id del query")
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Parámetros inválidos",
					"error":   "El business_id debe ser un número válido",
				})
				return
			}
		} else {
			h.logger.Error(ctx).Msg("Super admin sin business_id en query")
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Parámetro requerido",
				"error":   "El business_id es requerido en el query param",
			})
			return
		}
	}

	// Desactivar el recurso
	if err := h.usecase.ToggleBusinessResourceActive(ctx, businessID, resourceID, false); err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error al desactivar recurso")

		statusCode := http.StatusInternalServerError
		errorMessage := "Error al desactivar el recurso"

		if err.Error() == "business no encontrado" || err.Error() == "recurso no encontrado" {
			statusCode = http.StatusNotFound
			errorMessage = err.Error()
		} else if err.Error() == "el recurso no está permitido para este tipo de business" {
			statusCode = http.StatusBadRequest
			errorMessage = err.Error()
		} else if err.Error() == "la relación entre business y recurso no existe" {
			statusCode = http.StatusNotFound
			errorMessage = err.Error()
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": errorMessage,
			"error":   err.Error(),
		})
		return
	}

	h.logger.Info(ctx).Uint("business_id", businessID).Uint("resource_id", resourceID).Msg("Recurso desactivado exitosamente")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Recurso desactivado exitosamente",
	})
}
