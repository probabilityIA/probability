package handlers

import (
	"net/http"

	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/infra/primary/handlers/mapper"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/infra/primary/handlers/response"
	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/gin-gonic/gin"
)

// GetUsersHandler maneja la solicitud de obtener usuarios filtrados y paginados
//
//	@Summary		Obtener usuarios filtrados y paginados
//	@Description	Obtiene la lista filtrada y paginada de usuarios del sistema con sus roles y businesses
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page		query		int							false	"Número de página"	default(1)	minimum(1)
//	@Param			page_size	query		int							false	"Tamaño de página"	default(10)	minimum(1)	maximum(100)
//	@Param			name		query		string						false	"Filtrar por nombre (búsqueda parcial)"
//	@Param			email		query		string						false	"Filtrar por email (búsqueda parcial)"
//	@Param			phone		query		string						false	"Filtrar por teléfono (búsqueda parcial)"
//	@Param			user_ids	query		string						false	"Filtrar por IDs de usuarios separados por comas (ej: 1,2,3)"
//	@Param			is_active	query		bool						false	"Filtrar por estado activo"
//	@Param			role_id		query		int							false	"Filtrar por ID de rol"
//	@Param			business_id	query		int							false	"Filtrar por ID de business"
//	@Param			created_at	query		string						false	"Filtrar por fecha de creación (YYYY-MM-DD o YYYY-MM-DD,YYYY-MM-DD para rango)"
//	@Param			sort_by		query		string						false	"Campo para ordenar"		Enums(id, name, email, phone, is_active, created_at, updated_at)	default(created_at)
//	@Param			sort_order	query		string						false	"Orden de clasificación"	Enums(asc, desc)													default(desc)
//	@Success		200			{object}	response.UserListResponse	"Usuarios obtenidos exitosamente"
//	@Failure		400			{object}	response.UserErrorResponse	"Parámetros de filtro inválidos"
//	@Failure		401			{object}	response.UserErrorResponse	"Token de acceso requerido"
//	@Failure		500			{object}	response.UserErrorResponse	"Error interno del servidor"
//	@Router			/users [get]
func (h *handlers) GetUsersHandler(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "GetUsersHandler")

	// Crear struct de request y bindear parámetros de query
	var req request.GetUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error al parsear parámetros de query")
		c.JSON(http.StatusBadRequest, response.UserErrorResponse{
			Error: "Parámetros de filtro inválidos",
		})
		return
	}

	// Obtener información del usuario que hace la solicitud
	isSuperAdmin := middleware.IsSuperAdmin(c)
	requesterUserID, _ := middleware.GetUserID(c)
	requesterScope := middleware.GetScope(c) // "platform" o "business"

	// Si no es super admin (scope platform), restringir acceso
	if !isSuperAdmin {
		tokenBusinessID, ok := middleware.GetBusinessID(c)
		if ok && tokenBusinessID > 0 {
			req.BusinessID = &tokenBusinessID
			h.logger.Info(ctx).Uint("business_id", tokenBusinessID).Msg("Usando business_id del token para usuario normal")
		} else {
			h.logger.Error(ctx).Msg("Business ID no disponible en token")
			c.JSON(http.StatusUnauthorized, response.UserErrorResponse{
				Error: "Token inválido: business_id no disponible",
			})
			return
		}
	}

	// Convertir request a filtros del dominio
	filters := mapper.ToUserFilters(req)

	// Agregar información del requester para filtrado de seguridad
	filters.RequesterScope = requesterScope
	filters.RequesterUserID = requesterUserID

	h.logger.Info(ctx).
		Int("page", filters.Page).
		Int("page_size", filters.PageSize).
		Str("name", filters.Name).
		Str("email", filters.Email).
		Str("phone", filters.Phone).
		Bool("is_super_admin", isSuperAdmin).
		Str("requester_scope", requesterScope).
		Uint("requester_user_id", requesterUserID).
		Msg("Iniciando solicitud para obtener usuarios filtrados y paginados")

	userListDTO, err := h.usecase.GetUsers(ctx, filters)
	if err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error al obtener usuarios desde el caso de uso")
		c.JSON(http.StatusInternalServerError, response.UserErrorResponse{
			Error: "Error interno del servidor",
		})
		return
	}

	response := mapper.ToUserListResponse(userListDTO)

	h.logger.Info(ctx).
		Int("count", len(userListDTO.Users)).
		Int64("total", userListDTO.Total).
		Int("current_page", userListDTO.Page).
		Int("per_page", userListDTO.PageSize).
		Int("last_page", userListDTO.TotalPages).
		Bool("has_next", userListDTO.Page < userListDTO.TotalPages).
		Bool("has_prev", userListDTO.Page > 1).
		Msg("Usuarios obtenidos exitosamente con paginación")
	c.JSON(http.StatusOK, response)
}
