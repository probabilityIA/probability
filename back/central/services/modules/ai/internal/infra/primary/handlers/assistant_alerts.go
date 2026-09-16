package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/infra/primary/handlers/mappers"
)

func (h *Handlers) alertScope(c *gin.Context) (dtos.AccessScope, uint, bool) {
	scope, ok := accessScope(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "code": "unauthorized", "message": "Tu sesión no es válida. Vuelve a iniciar sesión."})
		return scope, 0, false
	}
	businessID := scope.TokenBusinessID
	if businessID == 0 {
		businessID = scope.RequestedBusinessID
	}
	if businessID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "code": "no_business", "message": "Selecciona un negocio para ver sus alertas."})
		return scope, 0, false
	}
	return scope, businessID, true
}

func (h *Handlers) ListAssistantAlerts(c *gin.Context) {
	scope, businessID, ok := h.alertScope(c)
	if !ok {
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.uc.ListAlerts(c.Request.Context(), dtos.AlertQuery{BusinessID: businessID, UserID: scope.UserID, Page: page, PageSize: pageSize})
	if err != nil {
		h.log.Error(c.Request.Context()).Err(err).Uint("business_id", businessID).Msg("[ai.assistant] error listando alertas")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "code": "internal_error", "message": "No se pudieron leer las alertas."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": mappers.FromAlerts(result)})
}

func (h *Handlers) GetAssistantAlertsUnread(c *gin.Context) {
	scope, businessID, ok := h.alertScope(c)
	if !ok {
		return
	}
	unread, err := h.uc.GetAlertsUnread(c.Request.Context(), businessID, scope.UserID)
	if err != nil {
		h.log.Error(c.Request.Context()).Err(err).Uint("business_id", businessID).Msg("[ai.assistant] error contando alertas")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "code": "internal_error", "message": "No se pudieron leer las alertas."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": mappers.FromAlertsUnread(unread)})
}

func (h *Handlers) MarkAssistantAlertsSeen(c *gin.Context) {
	scope, businessID, ok := h.alertScope(c)
	if !ok {
		return
	}
	if err := h.uc.MarkAlertsSeen(c.Request.Context(), businessID, scope.UserID); err != nil {
		h.log.Error(c.Request.Context()).Err(err).Uint("business_id", businessID).Msg("[ai.assistant] error marcando alertas vistas")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "code": "internal_error", "message": "No se pudieron marcar las alertas."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
