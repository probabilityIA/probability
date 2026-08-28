package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/secamc93/probability/back/central/services/modules/tours/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/tours/internal/infra/primary/handlers/mappers"
)

func (h *Handler) ListProgress(c *gin.Context) {
	userID, businessID, ok := h.resolveIdentity(c)
	if !ok {
		return
	}

	resultado, err := h.uc.ListProgress(c.Request.Context(), dtos.ListProgressInput{
		UserID:     userID,
		BusinessID: businessID,
	})
	if err != nil {
		h.log.Error(c.Request.Context()).Err(err).Uint("user_id", userID).Msg("Error al consultar el progreso de tours")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": mappers.ToResponseList(resultado.Items)})
}
