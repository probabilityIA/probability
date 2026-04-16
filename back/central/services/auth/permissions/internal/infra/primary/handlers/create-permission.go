package permissionhandler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/infra/primary/handlers/mapper"
	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/infra/primary/handlers/response"
)

// CreatePermissionHandler maneja la solicitud de crear un nuevo permiso
//
//	@Summary		Crear nuevo permiso
//	@Description	Crea un nuevo permiso en el sistema
//	@Tags			Permissions
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			permission	body		request.CreatePermissionRequest		true	"Datos del permiso a crear"
//	@Success		201			{object}	response.PermissionMessageResponse	"Permiso creado exitosamente"
//	@Failure		400			{object}	response.PermissionErrorResponse	"Datos de entrada inválidos"
//	@Failure		401			{object}	response.PermissionErrorResponse	"Token de acceso requerido"
//	@Failure		409			{object}	response.PermissionErrorResponse	"Permiso con código o nombre duplicado"
//	@Failure		500			{object}	response.PermissionErrorResponse	"Error interno del servidor"
//	@Router			/permissions [post]
func (h *PermissionHandler) CreatePermissionHandler(c *gin.Context) {
	var req request.CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Error al validar datos de entrada para crear permiso")
		c.JSON(http.StatusBadRequest, response.PermissionErrorResponse{
			Error: "Datos de entrada inválidos: " + err.Error(),
		})
		return
	}

	h.logger.Info().
		Str("name", req.Name).
		Str("code", req.Code).
		Uint("resource_id", req.ResourceID).
		Uint("action_id", req.ActionID).
		Msg("Iniciando solicitud para crear nuevo permiso")

	permissionDTO := mapper.ToCreatePermissionDTO(req)

	result, err := h.usecase.CreatePermission(c.Request.Context(), permissionDTO)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"
		errMsg := err.Error()

		// Detectar errores de validación de negocio (deben loggearse como WARN)
		if errMsg == "ya existe un permiso con el código: "+req.Code {
			statusCode = http.StatusConflict
			errorMessage = "Ya existe un permiso con este código"
			h.logger.Warn().Str("code", req.Code).Msg("Intento de crear permiso con código duplicado")
		} else if strings.Contains(errMsg, "ya existe un permiso con el nombre") {
			// Detectar error de nombre duplicado (validación previa)
			statusCode = http.StatusConflict
			errorMessage = errMsg
			h.logger.Warn().Str("name", req.Name).Msg("Intento de crear permiso con nombre duplicado")
		} else if strings.Contains(errMsg, "duplicate key") && strings.Contains(errMsg, "uni_permission_name") {
			// Fallback: Detectar error de nombre duplicado desde la BD (no debería ocurrir)
			statusCode = http.StatusConflict
			errorMessage = "Ya existe un permiso con el nombre '" + req.Name + "'. Por favor, use un nombre diferente."
			h.logger.Error().Err(err).Str("name", req.Name).Msg("Error de constraint de BD - validación previa falló")
		} else if strings.Contains(errMsg, "duplicate key") && strings.Contains(errMsg, "SQLSTATE 23505") {
			// Fallback: Detectar cualquier otro error de clave duplicada desde la BD
			statusCode = http.StatusConflict
			errorMessage = "Ya existe un permiso con estos datos. Por favor, verifique los valores ingresados."
			h.logger.Error().Err(err).Msg("Error de constraint de BD - validación previa falló")
		} else {
			// Error real del sistema
			h.logger.Error().Err(err).Msg("Error al crear permiso desde el caso de uso")
		}

		c.JSON(statusCode, response.PermissionErrorResponse{
			Error: errorMessage,
		})
		return
	}

	h.logger.Info().Str("result", result).Msg("Permiso creado exitosamente")
	c.JSON(http.StatusCreated, response.PermissionMessageResponse{
		Success: true,
		Message: result,
	})
}
