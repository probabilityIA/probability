package handlerintegrationtype

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrationtype/response"
)

func (h *IntegrationTypeHandler) buildWebhookURLs(code string) map[string]string {
	base := strings.TrimRight(h.env.Get("WEBHOOK_BASE_URL"), "/")
	if base == "" {
		base = strings.TrimRight(h.env.Get("URL_BASE_SWAGGER"), "/")
	}
	if base == "" {
		return nil
	}
	switch code {
	case "bold_pay":
		return map[string]string{
			"production": base + "/api/v1/webhooks/bold",
			"sandbox":    base + "/api/v1/webhooks/bold/test",
		}
	}
	return nil
}

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

	intType, _ := h.usecase.GetIntegrationTypeByID(c.Request.Context(), uint(id))
	h.recordAudit(c, intType, uint(id))

	resp := gin.H{
		"success": true,
		"message": "Credenciales de plataforma obtenidas exitosamente",
		"data":    creds,
	}
	if intType != nil {
		if urls := h.buildWebhookURLs(intType.Code); urls != nil {
			resp["webhook_urls"] = urls
		}
	}
	c.JSON(http.StatusOK, resp)
}

func (h *IntegrationTypeHandler) GetPlatformCredentialsByCodeHandler(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "Código requerido",
			Error:   "El código del integration_type es requerido",
		})
		return
	}

	creds, intType, err := h.usecase.GetPlatformCredentialsByCode(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, response.IntegrationErrorResponse{
			Success: false,
			Message: "Integration type no encontrado o error al desencriptar",
			Error:   err.Error(),
		})
		return
	}

	h.recordAudit(c, intType, intType.ID)

	resp := gin.H{
		"success": true,
		"message": "Credenciales de plataforma obtenidas exitosamente",
		"data":    creds,
		"meta": gin.H{
			"integration_type_id": intType.ID,
			"code":                intType.Code,
			"name":                intType.Name,
		},
	}
	if urls := h.buildWebhookURLs(intType.Code); urls != nil {
		resp["webhook_urls"] = urls
	}
	c.JSON(http.StatusOK, resp)
}

func (h *IntegrationTypeHandler) recordAudit(c *gin.Context, intType *domain.IntegrationType, fallbackID uint) {
	userID, _ := middleware.GetUserID(c)
	businessID, _ := middleware.GetBusinessID(c)

	audit := &domain.CredentialRevealAudit{
		UserID:            userID,
		BusinessID:        businessID,
		IntegrationTypeID: fallbackID,
		IPAddress:         c.ClientIP(),
		UserAgent:         c.Request.UserAgent(),
	}
	if intType != nil {
		audit.IntegrationTypeID = intType.ID
		audit.IntegrationCode = intType.Code
	}

	if err := h.usecase.RecordRevealAudit(c.Request.Context(), audit); err != nil {
		h.logger.Warn().Err(err).Msg("audit reveal failed (non-blocking)")
	}
}
