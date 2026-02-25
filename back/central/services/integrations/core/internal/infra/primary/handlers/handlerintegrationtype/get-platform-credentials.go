package handlerintegrationtype

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrationtype/response"
)

// GetPlatformCredentialsHandler decrypts and returns the platform credentials for the given integration type.
// Admin only — requires JWT.
func (h *IntegrationTypeHandler) GetPlatformCredentialsHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "ID inválido",
			Error:   "El ID debe ser un número válido",
		})
		return
	}

	creds, err := h.usecase.GetPlatformCredentials(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.IntegrationErrorResponse{
			Success: false,
			Message: "Error al obtener credenciales de plataforma",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Credenciales de plataforma obtenidas exitosamente",
		"data":    creds,
	})
}
