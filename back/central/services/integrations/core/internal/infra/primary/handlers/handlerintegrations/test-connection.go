package handlerintegrations

import (
	"errors"
	"net/http"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations/response"
)

func humanizeTestError(err error) string {
	msg := err.Error()
	if errors.Is(err, domain.ErrIntegrationTestFailed) {
		msg = strings.TrimPrefix(msg, domain.ErrIntegrationTestFailed.Error()+": ")
	}
	if errors.Is(err, domain.ErrIntegrationAccessTokenNotFound) {
		return "Falta el token de acceso"
	}
	msg = strings.TrimSpace(msg)
	if msg == "" {
		return "No pudimos probar la conexion. Verifica los datos e intenta de nuevo"
	}
	r := []rune(msg)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

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
	var req TestConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("❌ Error al bindear JSON para TestConnectionRaw")
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "Datos inválidos",
			Error:   err.Error(),
		})
		return
	}

	h.logger.Info().
		Str("type_code", req.TypeCode).
		Interface("config", req.Config).
		Msg("📥 TestConnectionRaw - Request received")

	if err := h.usecase.TestConnectionRaw(c.Request.Context(), req.TypeCode, req.Config, req.Credentials); err != nil {
		statusCode := http.StatusBadRequest
		errorMsg := humanizeTestError(err)

		h.logger.Error().
			Err(err).
			Str("type_code", req.TypeCode).
			Int("status_code", statusCode).
			Msg("❌ Fallo en TestConnectionRaw")
		c.JSON(statusCode, response.IntegrationErrorResponse{
			Success: false,
			Message: errorMsg,
			Error:   err.Error(),
		})
		return
	}

	h.logger.Info().
		Str("type_code", req.TypeCode).
		Msg("✅ TestConnectionRaw - Successful")

	c.JSON(http.StatusOK, response.IntegrationMessageResponse{
		Success: true,
		Message: "Conexión probada exitosamente",
	})
}
