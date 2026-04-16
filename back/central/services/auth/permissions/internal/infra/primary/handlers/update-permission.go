package permissionhandler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/infra/primary/handlers/mapper"
	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/infra/primary/handlers/response"
	"gorm.io/gorm"
)

// UpdatePermissionHandler maneja la solicitud de actualizar un permiso existente
//
//	@Summary		Actualizar permiso
//	@Description	Actualiza un permiso existente en el sistema
//	@Tags			Permissions
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id			path		int									true	"ID del permiso"	minimum(1)
//	@Param			permission	body		request.UpdatePermissionRequest		true	"Datos del permiso a actualizar"
//	@Success		200			{object}	response.PermissionMessageResponse	"Permiso actualizado exitosamente"
//	@Failure		400			{object}	response.PermissionErrorResponse	"Datos de entrada inválidos"
//	@Failure		401			{object}	response.PermissionErrorResponse	"Token de acceso requerido"
//	@Failure		404			{object}	response.PermissionErrorResponse	"Permiso no encontrado"
//	@Failure		409			{object}	response.PermissionErrorResponse	"Permiso con código duplicado"
//	@Failure		500			{object}	response.PermissionErrorResponse	"Error interno del servidor"
//	@Router			/permissions/{id} [put]
func (h *PermissionHandler) UpdatePermissionHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error().Str("id", idStr).Err(err).Msg("ID de permiso inválido")
		c.JSON(http.StatusBadRequest, response.PermissionErrorResponse{
			Error: "ID de permiso inválido",
		})
		return
	}

	var req request.UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Error al validar datos de entrada para actualizar permiso")
		c.JSON(http.StatusBadRequest, response.PermissionErrorResponse{
			Error: "Datos de entrada inválidos: " + err.Error(),
		})
		return
	}

	h.logger.Info().
		Uint64("id", id).
		Str("name", req.Name).
		Str("code", req.Code).
		Msg("Iniciando solicitud para actualizar permiso")

	permissionDTO := mapper.ToUpdatePermissionDTO(req)

	result, err := h.usecase.UpdatePermission(c.Request.Context(), uint(id), permissionDTO)
	if err != nil {
		h.logger.Error().Uint64("id", id).Err(err).Msg("Error al actualizar permiso desde el caso de uso")

		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"

		errMsg := err.Error()
		if errors.Is(err, gorm.ErrRecordNotFound) ||
			strings.Contains(errMsg, "record not found") ||
			strings.Contains(errMsg, "permiso no encontrado") {
			statusCode = http.StatusNotFound
			errorMessage = "Permiso no encontrado"
		} else if strings.HasPrefix(errMsg, "ya existe otro permiso con el código: ") {
			statusCode = http.StatusConflict
			errorMessage = "Ya existe otro permiso con este código"
		}

		c.JSON(statusCode, response.PermissionErrorResponse{
			Error: errorMessage,
		})
		return
	}

	h.logger.Info().Uint64("id", id).Str("result", result).Msg("Permiso actualizado exitosamente")
	c.JSON(http.StatusOK, response.PermissionMessageResponse{
		Success: true,
		Message: result,
	})
}
