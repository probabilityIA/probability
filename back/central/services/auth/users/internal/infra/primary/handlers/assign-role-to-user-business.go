package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/infra/primary/handlers/response"
)

// AssignRoleToUserBusinessHandler maneja la solicitud de asignar roles a un usuario en múltiples businesses
//
//	@Summary		Asignar roles a usuario en businesses
//	@Description	Asigna o actualiza roles de un usuario en múltiples businesses. El usuario debe estar previamente asociado a cada business. Solo se permite un rol por business y cada rol debe ser del mismo tipo de business que su business asociado.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id			path		int								true	"ID del usuario"
//	@Param			request		body		request.AssignRoleToUserBusinessRequest	true	"Lista de asignaciones (business_id y role_id por cada asignación)"
//	@Success		200			{object}	response.AssignRoleToUserBusinessResponse	"Roles asignados exitosamente"
//	@Failure		400			{object}	response.UserErrorResponse	"Datos inválidos"
//	@Failure		401			{object}	response.UserErrorResponse	"Token de acceso requerido"
//	@Failure		403			{object}	response.UserErrorResponse	"El usuario no está asociado a algún business o algún rol no corresponde al tipo de business"
//	@Failure		404			{object}	response.UserErrorResponse	"Usuario, business o rol no encontrado"
//	@Failure		500			{object}	response.UserErrorResponse	"Error interno del servidor"
//	@Router			/users/{id}/assign-role [post]
func (h *handlers) AssignRoleToUserBusinessHandler(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "AssignRoleToUserBusinessHandler")

	// Obtener user_id de la URL
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		h.logger.Error(ctx).Err(err).Str("user_id", userIDStr).Msg("Error al parsear user_id")
		c.JSON(http.StatusBadRequest, response.UserErrorResponse{
			Error: "ID de usuario inválido",
		})
		return
	}

	// Verificar si el usuario autenticado es super admin o si está asignando rol a sí mismo
	authenticatedUserID, exists := middleware.GetUserID(c)
	if !exists {
		h.logger.Error(ctx).Msg("Usuario no autenticado")
		c.JSON(http.StatusUnauthorized, response.UserErrorResponse{
			Error: "Usuario no autenticado",
		})
		return
	}

	isSuperAdmin := middleware.IsSuperAdmin(c)
	if !isSuperAdmin && authenticatedUserID != uint(userID) {
		h.logger.Warn(ctx).
			Uint("authenticated_user_id", authenticatedUserID).
			Uint("target_user_id", uint(userID)).
			Msg("Usuario no super admin intentando asignar roles a otro usuario")
		c.JSON(http.StatusForbidden, response.UserErrorResponse{
			Error: "No tienes permisos para asignar roles a otros usuarios",
		})
		return
	}

	// Parsear request body
	var req request.AssignRoleToUserBusinessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error al validar request")
		c.JSON(http.StatusBadRequest, response.UserErrorResponse{
			Error: "Datos de entrada inválidos: " + err.Error(),
		})
		return
	}

	if len(req.Assignments) == 0 {
		h.logger.Error(ctx).Msg("No se proporcionaron asignaciones")
		c.JSON(http.StatusBadRequest, response.UserErrorResponse{
			Error: "Debe proporcionar al menos una asignación",
		})
		return
	}

	// Convertir request a domain.BusinessRoleAssignment
	assignments := make([]domain.BusinessRoleAssignment, len(req.Assignments))
	for i, assignment := range req.Assignments {
		assignments[i] = domain.BusinessRoleAssignment{
			BusinessID: assignment.BusinessID,
			RoleID:     assignment.RoleID,
		}
	}

	h.logger.Info(ctx).
		Uint("user_id", uint(userID)).
		Int("assignments_count", len(assignments)).
		Msg("Iniciando asignación de roles a usuario en businesses")

	// Ejecutar caso de uso
	if err := h.usecase.AssignRoleToUserBusiness(ctx, uint(userID), assignments); err != nil {
		h.logger.Error(ctx).Err(err).
			Uint("user_id", uint(userID)).
			Int("assignments_count", len(assignments)).
			Msg("Error al asignar roles a usuario en businesses")

		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"

		errMsg := err.Error()
		if strings.Contains(errMsg, "usuario no encontrado") {
			statusCode = http.StatusNotFound
			errorMessage = "Usuario no encontrado"
		} else if strings.Contains(errMsg, "business") && strings.Contains(errMsg, "no encontrado") {
			statusCode = http.StatusNotFound
			errorMessage = "Algunos businesses no fueron encontrados"
		} else if strings.Contains(errMsg, "rol") && strings.Contains(errMsg, "no encontrado") {
			statusCode = http.StatusNotFound
			errorMessage = "Algunos roles no fueron encontrados"
		} else if strings.Contains(errMsg, "no está asociado al business") {
			statusCode = http.StatusForbidden
			errorMessage = errMsg
		} else if strings.Contains(errMsg, "no corresponde al tipo de business") {
			statusCode = http.StatusForbidden
			errorMessage = errMsg
		} else if strings.Contains(errMsg, "no se proporcionaron asignaciones") {
			statusCode = http.StatusBadRequest
			errorMessage = "No se proporcionaron asignaciones"
		}

		c.JSON(statusCode, response.UserErrorResponse{
			Error: errorMessage,
		})
		return
	}

	h.logger.Info(ctx).
		Uint("user_id", uint(userID)).
		Int("assignments_count", len(assignments)).
		Msg("Roles asignados exitosamente a usuario en businesses")

	c.JSON(http.StatusOK, response.AssignRoleToUserBusinessResponse{
		Success: true,
		Message: "Roles asignados exitosamente al usuario en los businesses",
	})
}
