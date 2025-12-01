package rolehandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/infra/primary/handlers/mapper"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/infra/primary/handlers/response"
	"github.com/secamc93/probability/back/central/shared/log"
)

// GetRolesHandler maneja la solicitud de obtener todos los roles
//
//	@Summary		Obtener todos los roles
//	@Description	Obtiene la lista completa de roles del sistema con opciones de filtrado
//	@Tags			Roles
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			business_type_id	query		int		false	"Filtrar por tipo de business"
//	@Param			scope_id			query		int		false	"Filtrar por ID de scope"
//	@Param			is_system			query		bool	false	"Filtrar por rol de sistema (true/false)"
//	@Param			name				query		string	false	"Buscar en el nombre del rol (búsqueda parcial)"
//	@Param			level				query		int		false	"Filtrar por nivel del rol"
//	@Success		200					{object}	response.RoleListResponse	"Roles obtenidos exitosamente"
//	@Failure		401					{object}	response.RoleErrorResponse		"Token de acceso requerido"
//	@Failure		500					{object}	response.RoleErrorResponse		"Error interno del servidor"
//	@Router			/roles [get]
func (h *RoleHandler) GetRolesHandler(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "GetRolesHandler")

	// Leer todos los query parameters
	filters := domain.RoleFilters{}

	// Si no es super admin, obtener business_type_id del token
	isSuperAdmin := middleware.IsSuperAdmin(c)
	if !isSuperAdmin {
		// Obtener business_type_id del token
		tokenBusinessTypeID, ok := middleware.GetBusinessTypeID(c)
		if ok && tokenBusinessTypeID > 0 {
			filters.BusinessTypeID = &tokenBusinessTypeID
			h.logger.Info(ctx).Uint("business_type_id", tokenBusinessTypeID).Msg("Usuario normal: filtrando roles por business_type_id del token")
		}
	} else {
		// Super admin puede filtrar por business_type_id desde query param
		if businessTypeIDStr := c.Query("business_type_id"); businessTypeIDStr != "" {
			if id, err := strconv.ParseUint(businessTypeIDStr, 10, 32); err == nil {
				val := uint(id)
				filters.BusinessTypeID = &val
				h.logger.Info(ctx).Uint("business_type_id", val).Msg("Super admin: filtrando roles por business_type_id del query")
			}
		}
	}

	// Filtrar por scope_id
	if scopeIDStr := c.Query("scope_id"); scopeIDStr != "" {
		if id, err := strconv.ParseUint(scopeIDStr, 10, 32); err == nil {
			val := uint(id)
			filters.ScopeID = &val
		}
	}

	// Filtrar por is_system
	if isSystemStr := c.Query("is_system"); isSystemStr != "" {
		if val, err := strconv.ParseBool(isSystemStr); err == nil {
			filters.IsSystem = &val
		}
	}

	// Filtrar por name (búsqueda parcial)
	if nameStr := c.Query("name"); nameStr != "" {
		filters.Name = &nameStr
	}

	// Filtrar por level
	if levelStr := c.Query("level"); levelStr != "" {
		if level, err := strconv.Atoi(levelStr); err == nil {
			filters.Level = &level
		}
	}

	roles, err := h.usecase.GetRoles(ctx, filters)
	if err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error al obtener roles desde el caso de uso")
		c.JSON(http.StatusInternalServerError, response.RoleErrorResponse{
			Error: "Error interno del servidor",
		})
		return
	}

	response := mapper.ToRoleListResponse(roles)

	h.logger.Info(ctx).Int("count", len(roles)).Bool("is_super_admin", isSuperAdmin).Msg("Roles obtenidos exitosamente")
	c.JSON(http.StatusOK, response)
}
