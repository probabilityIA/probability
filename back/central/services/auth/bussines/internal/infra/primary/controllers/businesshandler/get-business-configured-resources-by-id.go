package businesshandler

import (
	"net/http"
	"strconv"

	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/gin-gonic/gin"
)

// GetBusinessConfiguredResourcesByIDHandler obtiene la configuración de recursos de un business específico por ID
//
//	@Summary		Obtener configuración de recursos de un business por ID
//	@Description	Obtiene la configuración de recursos de un business específico por su ID, incluyendo todos los recursos asociados y su estado (activo/inactivo)
//	@Tags			Business Resources
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int						true	"ID del business"
//	@Success		200	{object}	map[string]interface{}	"Configuración de recursos del business obtenida exitosamente"
//	@Failure		400	{object}	map[string]interface{}	"Parámetros inválidos"
//	@Failure		401	{object}	map[string]interface{}	"No autorizado"
//	@Failure		403	{object}	map[string]interface{}	"Sin permisos"
//	@Failure		404	{object}	map[string]interface{}	"Business no encontrado"
//	@Failure		500	{object}	map[string]interface{}	"Error interno del servidor"
//	@Router			/businesses/{id}/configured-resources [get]
//	@Security		BearerAuth
func (h *BusinessHandler) GetBusinessConfiguredResourcesByIDHandler(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "GetBusinessConfiguredResourcesByIDHandler")

	// Parsear ID del path
	idStr := c.Param("id")
	businessID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error(ctx).Err(err).Str("id", idStr).Msg("Error al parsear ID de business")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID de business inválido",
			"error":   "El parámetro 'id' debe ser un número válido",
		})
		return
	}

	// Aplicar lógica de super admin
	isSuperAdmin := middleware.IsSuperAdmin(c)
	if !isSuperAdmin {
		// Usuario normal: validar que el business_id del token coincida con el ID solicitado
		tokenBusinessID, ok := middleware.GetBusinessID(c)
		if !ok || tokenBusinessID == 0 {
			h.logger.Error(ctx).Msg("Usuario normal sin business_id en token")
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Sin permisos para acceder a esta configuración",
				"error":   "Usuario sin business asignado",
			})
			return
		}

		if uint(businessID) != tokenBusinessID {
			h.logger.Warn(ctx).
				Uint("requested_business_id", uint(businessID)).
				Uint("token_business_id", tokenBusinessID).
				Msg("Usuario normal intentando acceder a business diferente al de su token")
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Sin permisos para acceder a esta configuración",
				"error":   "Solo puedes acceder a la configuración de tu propio business",
			})
			return
		}

		h.logger.Info(ctx).Uint("business_id", tokenBusinessID).Msg("Usuario normal: accediendo a configuración de su business")
	} else {
		h.logger.Info(ctx).Uint("business_id", uint(businessID)).Msg("Super admin: accediendo a configuración de business")
	}

	// Llamar al método del use case
	business, err := h.usecase.GetBusinessConfiguredResourcesByID(ctx, uint(businessID))
	if err != nil {
		if err.Error() == "business no encontrado" {
			h.logger.Warn(ctx).Uint("business_id", uint(businessID)).Msg("Business no encontrado")
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Business no encontrado",
				"error":   err.Error(),
			})
			return
		}

		h.logger.Error(ctx).Err(err).Uint("business_id", uint(businessID)).Msg("Error al obtener business con recursos configurados")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error interno del servidor",
			"error":   err.Error(),
		})
		return
	}

	h.logger.Info(ctx).
		Uint("business_id", uint(businessID)).
		Int("resources_count", len(business.Resources)).
		Bool("is_super_admin", isSuperAdmin).
		Msg("Business con recursos configurados obtenido exitosamente")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Configuración de recursos del business obtenida exitosamente",
		"data":    business,
	})
}
