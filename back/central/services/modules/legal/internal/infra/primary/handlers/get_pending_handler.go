package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/legal/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/legal/internal/infra/primary/handlers/response"
)

func (h *Handler) GetPending(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "usuario no identificado"})
		return
	}

	resultado, err := h.uc.GetPendingDocuments(c.Request.Context(), userID)
	if err != nil {
		h.log.Error(c.Request.Context()).Err(err).Uint("user_id", userID).Msg("Error al consultar documentos legales pendientes")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": response.PendingDocuments{
			RequiresAcceptance: resultado.RequiresAcceptance,
			Documents:          mappers.ToResponseDocuments(resultado.Documents),
		},
	})
}
