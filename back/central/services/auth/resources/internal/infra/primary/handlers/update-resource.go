package handlers

import (
	"net/http"
	"strconv"

	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/auth/resources/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/resources/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/resources/internal/infra/primary/handlers/response"
	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/gin-gonic/gin"
)

// UpdateResourceHandler actualiza un recurso existente
//
//	@Summary		Actualizar recurso
//	@Description	Actualiza un recurso existente en el sistema
//	@Tags			Resources
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int								true	"ID del recurso"	minimum(1)
//	@Param			request	body		request.UpdateResourceRequest	true	"Datos del recurso a actualizar"
//	@Success		200		{object}	map[string]interface{}			"Recurso actualizado exitosamente"
//	@Failure		400		{object}	map[string]interface{}			"Datos de entrada inválidos"
//	@Failure		401		{object}	map[string]interface{}			"No autorizado"
//	@Failure		404		{object}	map[string]interface{}			"Recurso no encontrado"
//	@Failure		409		{object}	map[string]interface{}			"Conflicto con recurso existente"
//	@Failure		500		{object}	map[string]interface{}			"Error interno del servidor"
//	@Router			/resources/{id} [put]
//	@Security		BearerAuth
func (h *ResourceHandler) UpdateResourceHandler(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "UpdateResourceHandler")

	// Validar que el usuario sea super admin
	if !middleware.IsSuperAdmin(c) {
		h.logger.Warn(ctx).Msg("Intento de actualización de recurso por usuario no super admin")
		c.JSON(http.StatusForbidden, response.ErrorResponse{
			Success: false,
			Message: "Solo los super usuarios pueden actualizar recursos",
			Error:   "permisos insuficientes",
		})
		return
	}

	// Obtener el ID del recurso de los parámetros de la URL
	resourceIDStr := c.Param("id")
	resourceID, err := strconv.ParseUint(resourceIDStr, 10, 32)
	if err != nil {
		h.logger.Error(ctx).Err(err).Str("resource_id", resourceIDStr).Msg("Error al parsear ID del recurso")
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "ID de recurso inválido",
			Error:   "El ID del recurso debe ser un número válido",
		})
		return
	}

	h.logger.Info(ctx).Uint64("resource_id", resourceID).Msg("Iniciando actualización de recurso")

	// Parsear el cuerpo de la petición
	var req request.UpdateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error(ctx).Err(err).Uint64("resource_id", resourceID).Msg("Error al parsear el cuerpo de la petición")
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Datos de entrada inválidos",
			Error:   err.Error(),
		})
		return
	}

	h.logger.Info(ctx).
		Uint64("resource_id", resourceID).
		Str("name", req.Name).
		Msg("Datos de actualización de recurso recibidos")

	// Convertir a DTO de dominio
	updateDTO := domain.UpdateResourceDTO{
		Name:           req.Name,
		Description:    req.Description,
		BusinessTypeID: req.BusinessTypeID,
	}

	// Llamar al caso de uso
	result, err := h.usecase.UpdateResource(ctx, uint(resourceID), updateDTO)
	if err != nil {
		h.logger.Error(ctx).Err(err).Uint64("resource_id", resourceID).Msg("Error al actualizar recurso")

		// Determinar el tipo de error y el código de estado HTTP
		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"

		if err.Error() == "recurso con ID "+resourceIDStr+" no encontrado" {
			statusCode = http.StatusNotFound
			errorMessage = "Recurso no encontrado"
		} else if err.Error() == "ya existe otro recurso con el nombre '"+req.Name+"'" {
			statusCode = http.StatusConflict
			errorMessage = "Conflicto con recurso existente"
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
		Uint64("resource_id", resourceID).
		Str("name", result.Name).
		Msg("Recurso actualizado exitosamente")

	c.JSON(http.StatusOK, response.UpdateResourceResponse{
		Success: true,
		Message: "Recurso actualizado exitosamente",
		Data:    resourceResponse,
	})
}
