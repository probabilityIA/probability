package authhandler

import (
	"net/http"

	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/gin-gonic/gin"
)

// VerifyHandler verifica la autenticación del usuario
func (h *AuthHandler) VerifyHandler(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "VerifyHandler")

	h.logger.Info(ctx).Msg("Verificación de autenticación solicitada")

	// Obtener información de autenticación desde el middleware
	authInfo, exists := middleware.GetAuthInfo(c)
	if !exists {
		h.logger.Error(ctx).Msg("Información de autenticación no encontrada en contexto")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "No autorizado",
		})
		return
	}

	// Verificar que el tipo de autenticación sea JWT
	if authInfo.Type != middleware.AuthTypeJWT {
		h.logger.Error(ctx).Str("auth_type", string(authInfo.Type)).Msg("Tipo de autenticación no válido")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Tipo de autenticación no válido",
		})
		return
	}

	// Log de información de autenticación
	h.logger.Info(ctx).
		Uint("user_id", authInfo.UserID).
		Str("user_email", authInfo.Email).
		Strs("user_roles", authInfo.Roles).
		Uint("business_id", authInfo.BusinessID).
		Msg("Usuario autenticado correctamente")

	// Retornar información del usuario autenticado
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Usuario autenticado correctamente",
		"data": gin.H{
			"user_id":     authInfo.UserID,
			"email":       authInfo.Email,
			"roles":       authInfo.Roles,
			"business_id": authInfo.BusinessID,
		},
	})
}
