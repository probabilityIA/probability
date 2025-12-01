package permissionhandler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/infra/primary/handlers/response"
	"gorm.io/gorm"
)

// DeletePermissionHandler maneja la solicitud de eliminar un permiso
//
//	@Summary		Eliminar permiso permanentemente
//	@Description	Elimina permanentemente un permiso del sistema (eliminación física, no se puede recuperar)
//	@Tags			Permissions
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int									true	"ID del permiso"	minimum(1)
//	@Success		200	{object}	response.PermissionMessageResponse	"Permiso eliminado exitosamente"
//	@Failure		400	{object}	response.PermissionErrorResponse	"ID inválido"
//	@Failure		401	{object}	response.PermissionErrorResponse	"Token de acceso requerido"
//	@Failure		404	{object}	response.PermissionErrorResponse	"Permiso no encontrado"
//	@Failure		500	{object}	response.PermissionErrorResponse	"Error interno del servidor"
//	@Router			/permissions/{id} [delete]
func (h *PermissionHandler) DeletePermissionHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error().Str("id", idStr).Err(err).Msg("ID de permiso inválido")
		c.JSON(http.StatusBadRequest, response.PermissionErrorResponse{
			Error: "ID de permiso inválido",
		})
		return
	}

	h.logger.Info().Uint64("id", id).Msg("Iniciando solicitud para eliminar permiso")

	result, err := h.usecase.DeletePermission(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error().Uint64("id", id).Err(err).Msg("Error al eliminar permiso desde el caso de uso")

		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"

		errMsg := err.Error()
		// Verificar si es un error de "record not found" o contiene ese mensaje
		if errors.Is(err, gorm.ErrRecordNotFound) ||
			strings.Contains(errMsg, "record not found") ||
			strings.Contains(errMsg, "permiso no encontrado") {
			statusCode = http.StatusNotFound
			errorMessage = "Permiso no encontrado"
		}

		c.JSON(statusCode, response.PermissionErrorResponse{
			Error: errorMessage,
		})
		return
	}

	h.logger.Info().Uint64("id", id).Str("result", result).Msg("Permiso eliminado exitosamente")
	c.JSON(http.StatusOK, response.PermissionMessageResponse{
		Success: true,
		Message: result,
	})
}
