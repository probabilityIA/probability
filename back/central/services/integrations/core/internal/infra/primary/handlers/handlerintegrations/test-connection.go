package handlerintegrations

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations/response"
)

type TestConnectionRequest struct {
	TypeCode    string                 `json:"type_code" binding:"required"`
	Config      map[string]interface{} `json:"config"`
	Credentials map[string]interface{} `json:"credentials"`
}

// TestConnectionRawHandler prueba la conexión con datos proporcionados
//
//	@Summary		Probar conexión (Raw)
//	@Description	Prueba la conexión con credenciales y configuración proporcionadas sin guardar
//	@Tags			Integrations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		TestConnectionRequest	true	"Datos de prueba"
//	@Success		200		{object}	response.IntegrationMessageResponse
//	@Failure		400		{object}	response.IntegrationErrorResponse
//	@Failure		401		{object}	response.IntegrationErrorResponse
//	@Failure		500		{object}	response.IntegrationErrorResponse
//	@Router			/integrations/test [post]
func (h *IntegrationHandler) TestConnectionRawHandler(c *gin.Context) {
	// Solo super admins pueden probar integraciones
	if !middleware.IsSuperAdmin(c) {
		h.logger.Error().Str("endpoint", "/integrations/test").Str("method", "POST").Msg("Intento de probar integración sin permisos de super admin")
		c.JSON(http.StatusForbidden, response.IntegrationErrorResponse{
			Success: false,
			Message: "Solo los super usuarios pueden probar integraciones",
			Error:   "permisos insuficientes",
		})
		return
	}

	var req TestConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Error al bindear JSON para TestConnectionRaw")
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "Datos inválidos",
			Error:   err.Error(),
		})
		return
	}

	if err := h.usecase.TestConnectionRaw(c.Request.Context(), req.TypeCode, req.Config, req.Credentials); err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := "Error al probar conexión"

		if err == domain.ErrIntegrationTestFailed {
			statusCode = http.StatusBadRequest
			errorMsg = "La prueba de conexión falló"
		} else if err == domain.ErrIntegrationAccessTokenNotFound {
			statusCode = http.StatusBadRequest
			errorMsg = "Falta el token de acceso"
		}

		h.logger.Error().Err(err).Str("type_code", req.TypeCode).Msg("Fallo en TestConnectionRaw")
		c.JSON(statusCode, response.IntegrationErrorResponse{
			Success: false,
			Message: errorMsg,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.IntegrationMessageResponse{
		Success: true,
		Message: "Conexión probada exitosamente",
	})
}
