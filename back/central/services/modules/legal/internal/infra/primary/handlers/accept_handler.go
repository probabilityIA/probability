package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/legal/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/legal/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/legal/internal/infra/primary/handlers/response"
)

func (h *Handler) Accept(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "usuario no identificado"})
		return
	}

	var body request.AcceptDocuments
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "document_ids es requerido"})
		return
	}

	var businessID *uint
	if bid := c.GetUint("business_id"); bid > 0 {
		businessID = &bid
	}

	resultado, err := h.uc.AcceptDocuments(c.Request.Context(), dtos.AcceptDocumentsInput{
		UserID:      userID,
		BusinessID:  businessID,
		DocumentIDs: body.DocumentIDs,
		IPAddress:   c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
	})
	if err != nil {
		h.log.Error(c.Request.Context()).Err(err).Uint("user_id", userID).Msg("Error al registrar la aceptacion de documentos legales")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": response.AcceptResult{
			AcceptedAt:  resultado.AcceptedAt.Format(time.RFC3339),
			DocumentIDs: resultado.DocumentIDs,
		},
	})
}
