package permissionhandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/infra/primary/handlers/mapper"
	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/infra/primary/handlers/response"
)

// GetPermissionsHandler maneja la solicitud de obtener todos los permisos
//
//	@Summary		Obtener todos los permisos
//	@Description	Obtiene la lista completa de permisos del sistema con filtros opcionales por tipo de business, nombre y scope
//	@Tags			Permissions
//	@Accept			json
//	@Produce		json
//	@Param			business_type_id	query		int		false	"Filtrar por tipo de business (incluye genéricos)"
//	@Param			name				query		string	false	"Filtrar por nombre de permiso (búsqueda parcial)"
//	@Param			scope_id			query		int		false	"Filtrar por ID de scope"
//	@Security		BearerAuth
//	@Success		200					{object}	response.PermissionListResponse		"Lista de permisos obtenida exitosamente"
//	@Failure		401					{object}	response.PermissionErrorResponse	"Token de acceso requerido"
//	@Failure		500					{object}	response.PermissionErrorResponse	"Error interno del servidor"
//	@Router			/permissions [get]
func (h *PermissionHandler) GetPermissionsHandler(c *gin.Context) {
	h.logger.Info().Msg("Iniciando solicitud para obtener todos los permisos")

	// Leer query params
	var businessTypeID *uint
	var name *string
	var scopeID *uint

	if businessTypeIDStr := c.Query("business_type_id"); businessTypeIDStr != "" {
		if id, err := strconv.ParseUint(businessTypeIDStr, 10, 32); err == nil {
			id := uint(id)
			businessTypeID = &id
		}
	}

	if nameStr := c.Query("name"); nameStr != "" {
		name = &nameStr
	}

	if scopeIDStr := c.Query("scope_id"); scopeIDStr != "" {
		if id, err := strconv.ParseUint(scopeIDStr, 10, 32); err == nil {
			id := uint(id)
			scopeID = &id
		}
	}

	permissions, err := h.usecase.GetPermissions(c.Request.Context(), businessTypeID, name, scopeID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error al obtener permisos desde el caso de uso")
		c.JSON(http.StatusInternalServerError, response.PermissionErrorResponse{
			Error: "Error interno del servidor",
		})
		return
	}

	response := mapper.ToPermissionListResponse(permissions)

	h.logger.Info().Int("count", len(permissions)).Msg("Permisos obtenidos exitosamente")
	c.JSON(http.StatusOK, response)
}
