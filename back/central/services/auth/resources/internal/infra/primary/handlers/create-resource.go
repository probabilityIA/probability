package handlers

import (
	"net/http"

	"github.com/secamc93/probability/back/central/services/auth/middleware"

	"github.com/secamc93/probability/back/central/services/auth/resources/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/resources/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/resources/internal/infra/primary/handlers/response"
	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/gin-gonic/gin"
)

// CreateResourceHandler crea un nuevo recurso
//
//	@Summary		Crear recurso
//	@Description	Crea un nuevo recurso en el sistema con nombre y descripción únicos
//	@Tags			Resources
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.CreateResourceRequest	true	"Datos del recurso a crear"
//	@Success		201		{object}	map[string]interface{}			"Recurso creado exitosamente"
//	@Failure		400		{object}	map[string]interface{}			"Datos de entrada inválidos"
//	@Failure		401		{object}	map[string]interface{}			"No autorizado"
//	@Failure		409		{object}	map[string]interface{}			"Recurso ya existe"
//	@Failure		500		{object}	map[string]interface{}			"Error interno del servidor"
//	@Router			/resources [post]
//	@Security		BearerAuth
func (h *ResourceHandler) CreateResourceHandler(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "CreateResourceHandler")

	// Validar que el usuario sea super admin
	if !middleware.IsSuperAdmin(c) {
		h.logger.Warn(ctx).Msg("Intento de creación de recurso por usuario no super admin")
		c.JSON(http.StatusForbidden, response.ErrorResponse{
			Success: false,
			Message: "Solo los super usuarios pueden crear recursos",
			Error:   "permisos insuficientes",
		})
		return
	}

	h.logger.Info(ctx).Msg("Iniciando creación de recurso")

	// Parsear el cuerpo de la petición
	var req request.CreateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error al parsear el cuerpo de la petición")
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Datos de entrada inválidos",
			Error:   err.Error(),
		})
		return
	}

	h.logger.Info().Str("name", req.Name).Msg("Datos de creación de recurso recibidos")

	// Convertir a DTO de dominio
	createDTO := domain.CreateResourceDTO{
		Name:           req.Name,
		Description:    req.Description,
		BusinessTypeID: req.BusinessTypeID,
	}

	// Llamar al caso de uso
	result, err := h.usecase.CreateResource(ctx, createDTO)
	if err != nil {
		h.logger.Error(ctx).Err(err).Str("name", req.Name).Msg("Error al crear recurso")

		// Determinar el tipo de error y el código de estado HTTP
		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"

		if err.Error() == "ya existe un recurso con el nombre '"+req.Name+"'" {
			statusCode = http.StatusConflict
			errorMessage = "Recurso ya existe"
		} else if err.Error() == "el nombre del recurso es obligatorio" ||
			err.Error() == "el nombre del recurso no puede exceder 100 caracteres" ||
			err.Error() == "la descripción del recurso no puede exceder 500 caracteres" {
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
	resourceResponse := response.ResourceResponse{
		ID:               result.ID,
		Name:             result.Name,
		Description:      result.Description,
		BusinessTypeID:   result.BusinessTypeID,
		BusinessTypeName: result.BusinessTypeName,
		CreatedAt:        result.CreatedAt,
		UpdatedAt:        result.UpdatedAt,
	}

	h.logger.Info(ctx).
		Uint("resource_id", result.ID).
		Str("name", result.Name).
		Msg("Recurso creado exitosamente")

	c.JSON(http.StatusCreated, response.CreateResourceResponse{
		Success: true,
		Message: "Recurso creado exitosamente",
		Data:    resourceResponse,
	})
}
