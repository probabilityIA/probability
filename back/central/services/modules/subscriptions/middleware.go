package subscriptions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	authmiddleware "github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/app"
)

func RequireModuleAccess(uc app.IUseCase, moduleCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		businessID, ok := authmiddleware.GetBusinessIDFromContext(c)
		if !ok || businessID == 0 {
			c.Next()
			return
		}

		allowed, err := uc.HasModuleAccess(c.Request.Context(), businessID, moduleCode)
		if err != nil {
			c.Next()
			return
		}
		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Tu plan de suscripcion no incluye este modulo",
				"code":  "MODULE_NOT_INCLUDED",
			})
			return
		}

		c.Next()
	}
}
