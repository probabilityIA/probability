package businesshandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/infra/primary/controllers/businesshandler/mapper"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/shared/log"
)

// GetBusinesses godoc
//
//	@Summary		Obtener lista de negocios
//	@Description	Obtiene una lista paginada de todos los negocios del sistema
//	@Tags			businesses
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page				query		int						false	"Número de página (por defecto 1)"
//	@Param			per_page			query		int						false	"Elementos por página (por defecto 10)"
//	@Param			name				query		string					false	"Filtrar por nombre de negocio"
//	@Param			business_type_id	query		int						false	"Filtrar por tipo de negocio"
//	@Param			is_active			query		boolean					false	"Filtrar por estado activo/inactivo"
//	@Success		201					{object}	map[string]interface{}	"Negocios obtenidos exitosamente"
//	@Failure		400		{object}	map[string]interface{}	"Solicitud inválida"
//	@Failure		401		{object}	map[string]interface{}	"Token de acceso requerido"
//	@Failure		500		{object}	map[string]interface{}	"Error interno del servidor"
//	@Router			/businesses [get]
func (h *BusinessHandler) GetBusinesses(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "GetBusinesses")

	// Validar si es super admin
	isSuperAdmin := middleware.IsSuperAdmin(c)
	if !isSuperAdmin {
		h.logger.Warn(ctx).Msg("Intento de acceso no autorizado al endpoint de listar businesses")
		c.JSON(http.StatusForbidden, mapper.BuildErrorResponse("access_denied", "No tienes permisos para acceder a este endpoint"))
		return
	}

	h.logger.Info(ctx).Msg("Super admin accediendo al listado de businesses")

	// Obtener parámetros de paginación
	pageStr := c.DefaultQuery("page", "1")
	perPageStr := c.DefaultQuery("per_page", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage < 1 || perPage > 100 {
		perPage = 10
	}

	// Obtener filtros opcionales
	name := c.Query("name")
	businessTypeIDStr := c.Query("business_type_id")
	isActiveStr := c.Query("is_active")

	var businessTypeID *uint
	if businessTypeIDStr != "" {
		if id, err := strconv.ParseUint(businessTypeIDStr, 10, 32); err == nil {
			val := uint(id)
			businessTypeID = &val
		}
	}

	var isActive *bool
	if isActiveStr != "" {
		if active, err := strconv.ParseBool(isActiveStr); err == nil {
			isActive = &active
		}
	}

	// Ejecutar caso de uso
	businesses, total, err := h.usecase.GetBusinesses(c.Request.Context(), page, perPage, name, businessTypeID, isActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapper.BuildErrorResponse("internal_error", "Error interno del servidor"))
		return
	}

	// Construir respuesta exitosa con paginación
	response := mapper.BuildGetBusinessesResponseWithPagination(businesses, "Negocios obtenidos exitosamente", page, perPage, total)
	c.JSON(http.StatusOK, response)
}
