package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/push/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/push/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/push/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/push/internal/infra/primary/handlers/response"
)

func (h *handler) RegisterDevice(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "usuario no encontrado en el contexto"})
		return
	}

	var body request.RegisterDevice
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	businessID, ok := resolveBusinessID(c, body.BusinessID)
	if !ok {
		return
	}

	dto := dtos.RegisterDeviceDTO{
		UserID:     userID,
		BusinessID: businessID,
		Token:      body.Token,
		Platform:   body.Platform,
		AppVersion: body.AppVersion,
		DeviceName: body.DeviceName,
	}

	if err := h.uc.RegisterDevice(c.Request.Context(), dto); err != nil {
		status := http.StatusInternalServerError
		if domainerrors.IsNonRetryable(err) {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "dispositivo registrado"})
}

func (h *handler) UnregisterDevice(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "usuario no encontrado en el contexto"})
		return
	}

	token := c.Query("token")
	if token == "" {
		var body request.UnregisterDevice
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "token es requerido"})
			return
		}
		token = body.Token
	}

	if err := h.uc.UnregisterDevice(c.Request.Context(), userID, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "dispositivo dado de baja"})
}

func (h *handler) ListDevices(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "usuario no encontrado en el contexto"})
		return
	}

	devices, err := h.uc.ListDevices(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response.FromEntities(devices)})
}

func resolveBusinessID(c *gin.Context, bodyBusinessID *uint) (*uint, bool) {
	businessID, ok := middleware.GetBusinessID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "contexto de negocio no encontrado"})
		return nil, false
	}

	if businessID > 0 {
		return &businessID, true
	}

	if bodyBusinessID == nil || *bodyBusinessID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido para super admin"})
		return nil, false
	}
	return bodyBusinessID, true
}
