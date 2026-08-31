package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/secamc93/probability/back/central/services/modules/tours/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/tours/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/tours/internal/infra/primary/handlers/request"
)

func (h *Handler) SaveProgress(c *gin.Context) {
	userID, businessID, ok := h.resolveIdentity(c)
	if !ok {
		return
	}

	var body request.SaveProgress
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "tour_key, version y status son requeridos"})
		return
	}

	guardado, err := h.uc.SaveProgress(c.Request.Context(), dtos.SaveProgressInput{
		UserID:     userID,
		BusinessID: businessID,
		TourKey:    body.TourKey,
		Version:    body.Version,
		Status:     body.Status,
		StepIndex:  body.StepIndex,
	})
	if err != nil {
		h.log.Error(c.Request.Context()).Err(err).Uint("user_id", userID).Str("tour_key", body.TourKey).Msg("Error al guardar el progreso del tour")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": mappers.ToResponse(*guardado)})
}
