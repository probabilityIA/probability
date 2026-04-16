package permissionhandler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/infra/primary/handlers/mapper"
	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/infra/primary/handlers/request"
)

// BulkCreatePermissionsHandler maneja la creación masiva de permisos
//
//	@Summary		Crear permisos en masa
//	@Description	Crea múltiples permisos en una sola operación
//	@Tags			Permissions
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			permissions	body		request.BulkCreatePermissionRequest	true	"Lista de permisos a crear"
//	@Success		200			{object}	map[string]interface{}				"Resultado de la creación masiva"
//	@Failure		400			{object}	map[string]interface{}				"Datos inválidos"
//	@Router			/permissions/bulk [post]
func (h *PermissionHandler) BulkCreatePermissionsHandler(c *gin.Context) {
	var req request.BulkCreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Error al validar datos de entrada para creación masiva")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Datos de entrada inválidos: " + err.Error(),
		})
		return
	}

	dtos := make([]domain.CreatePermissionDTO, len(req.Permissions))
	for i, p := range req.Permissions {
		dtos[i] = mapper.ToBulkItemToCreateDTO(p)
	}

	h.logger.Info().Int("count", len(dtos)).Msg("Iniciando creación masiva de permisos")

	results, err := h.usecase.BulkCreatePermissions(c.Request.Context(), dtos)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	successCount := 0
	for _, r := range results {
		if r.Success {
			successCount++
		}
	}

	h.logger.Info().Int("success", successCount).Int("total", len(results)).Msg("Creación masiva completada")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("%d de %d permisos creados exitosamente", successCount, len(results)),
		"results": results,
	})
}
