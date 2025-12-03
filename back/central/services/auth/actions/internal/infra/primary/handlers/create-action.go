package handlers

import (
	"net/http"

	"github.com/secamc93/probability/back/central/services/auth/actions/internal/domain"

	"github.com/secamc93/probability/back/central/services/auth/actions/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/actions/internal/infra/primary/handlers/response"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/gin-gonic/gin"
)

// CreateActionHandler crea un nuevo action
//
//	@Summary		Crear action
//	@Description	Crea un nuevo action en el sistema con nombre y descripción únicos
//	@Tags			Actions
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.CreateActionRequest	true	"Datos del action a crear"
//	@Success		201		{object}	map[string]interface{}			"Action creado exitosamente"
//	@Failure		400		{object}	map[string]interface{}			"Datos de entrada inválidos"
//	@Failure		401		{object}	map[string]interface{}			"No autorizado"
//	@Failure		403		{object}	map[string]interface{}			"Solo super usuarios pueden crear actions"
//	@Failure		409		{object}	map[string]interface{}			"Action ya existe"
//	@Failure		500		{object}	map[string]interface{}			"Error interno del servidor"
//	@Router			/actions [post]
//	@Security		BearerAuth
func (h *ActionHandler) CreateActionHandler(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "CreateActionHandler")

	// Validar que el usuario sea super admin
	if !middleware.IsSuperAdmin(c) {
		h.logger.Warn(ctx).Msg("Intento de creación de action por usuario no super admin")
		c.JSON(http.StatusForbidden, response.ErrorResponse{
			Success: false,
			Message: "Solo los super usuarios pueden crear actions",
			Error:   "permisos insuficientes",
		})
		return
	}

	h.logger.Info(ctx).Msg("Iniciando creación de action")

	// Parsear el cuerpo de la petición
	var req request.CreateActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error al parsear el cuerpo de la petición")
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Datos de entrada inválidos",
			Error:   err.Error(),
		})
		return
	}

	h.logger.Info(ctx).Str("name", req.Name).Msg("Datos de creación de action recibidos")

	// Convertir a DTO de dominio
	createDTO := domain.CreateActionDTO{
		Name:        req.Name,
		Description: req.Description,
	}

	// Llamar al caso de uso
	result, err := h.usecase.CreateAction(ctx, createDTO)
	if err != nil {
		h.logger.Error(ctx).Err(err).Str("name", req.Name).Msg("Error al crear action")

		// Determinar el tipo de error y el código de estado HTTP
		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"

		if err.Error() == "ya existe un action con el nombre '"+req.Name+"'" {
			statusCode = http.StatusConflict
			errorMessage = "Action ya existe"
		} else if err.Error() == "el nombre del action es obligatorio" ||
			err.Error() == "el nombre del action no puede exceder 20 caracteres" ||
			err.Error() == "la descripción del action no puede exceder 255 caracteres" {
			statusCode = http.StatusBadRequest
			errorMessage = "Datos de entrada inválidos"
		}

		c.JSON(statusCode, response.ErrorResponse{
			Success: false,
			Message: errorMessage,
			Error:   err.Error(),
		})
		return
	}

	// Convertir a respuesta HTTP
	actionResponse := response.ActionResponse{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}

	h.logger.Info(ctx).
		Uint("action_id", result.ID).
		Str("name", result.Name).
		Msg("Action creado exitosamente")

	c.JSON(http.StatusCreated, response.CreateActionResponse{
		Success: true,
		Message: "Action creado exitosamente",
		Data:    actionResponse,
	})
}
