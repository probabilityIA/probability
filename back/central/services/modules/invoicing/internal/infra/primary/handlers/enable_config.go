package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/mappers"
)

// EnableConfig activa una configuración de facturación
// @Summary Activar configuración de facturación
// @Description Activa una configuración de facturación existente
// @Tags Invoicing Config
// @Accept json
// @Produce json
// @Param id path int true "Config ID"
// @Success 200 {object} response.Config
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /invoicing/configs/{id}/enable [post]
func (h *handler) EnableConfig(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid config ID"})
		return
	}

	enabled := true
	dto := &dtos.UpdateConfigDTO{
		Enabled: &enabled,
	}

	config, err := h.useCase.UpdateConfig(c.Request.Context(), uint(id), dto)
	if err != nil {
		h.log.Error(c.Request.Context()).Err(err).Uint("config_id", uint(id)).Msg("Failed to enable config")
		handleDomainError(c, err, "config_enable_failed")
		return
	}

	baseURL, bucket := h.getS3Config()
	c.JSON(http.StatusOK, mappers.ConfigToResponse(config, baseURL, bucket))
}
