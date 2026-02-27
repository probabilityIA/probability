package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

// resolveBusinessID retorna el business_id efectivo para la request.
// Super admin puede proveer ?business_id=X como query param; negocio normal usa el JWT.
func resolveBusinessID(c *gin.Context) (uint, bool) {
	if middleware.IsSuperAdmin(c) {
		param := c.Query("business_id")
		if param == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "business_id es requerido para super admin"})
			return 0, false
		}
		id, err := strconv.ParseUint(param, 10, 64)
		if err != nil || id == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "business_id inválido"})
			return 0, false
		}
		return uint(id), true
	}
	businessID, exists := middleware.GetBusinessID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return 0, false
	}
	return businessID, true
}
