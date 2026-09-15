package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func resolveBusinessID(c *gin.Context, requestedBusinessID uint) (uint, bool) {
	businessID, ok := middleware.GetBusinessIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "message": "contexto de negocio no encontrado"})
		return 0, false
	}
	if businessID > 0 {
		return businessID, true
	}
	if requestedBusinessID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": "business_id es requerido para super admin"})
		return 0, false
	}
	return requestedBusinessID, true
}
