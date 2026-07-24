package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/infra/primary/handlers/request"
)

func (h *Handlers) GetDianConfig(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}
	status, err := h.uc.GetDianConfig(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": status})
}

func (h *Handlers) SaveDianConfig(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}
	var req request.DianConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	status, err := h.uc.SaveDianConfig(c.Request.Context(), dtos.DianConfigDTO{
		ApiURL:            req.ApiURL,
		ClientID:          req.ClientID,
		ClientSecret:      req.ClientSecret,
		Username:          req.Username,
		Password:          req.Password,
		NumberingRangeID:  req.NumberingRangeID,
		MunicipalityID:    req.MunicipalityID,
		DocumentIDType:    req.DocumentIDType,
		LegalOrgID:        req.LegalOrgID,
		PaymentMethodCode: req.PaymentMethodCode,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": status})
}
