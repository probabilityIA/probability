package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/secamc93/probability/back/central/services/modules/tours/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/tours/internal/infra/primary/handlers/request"
)

func (h *Handler) SkipAll(c *gin.Context) {
	userID, businessID, ok := h.resolveIdentity(c)
	if !ok {
		return
	}

	var body request.SkipAll
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "tours es requerido"})
		return
	}

	tours := make([]dtos.SkipAllTour, 0, len(body.Tours))
	for _, tour := range body.Tours {
		tours = append(tours, dtos.SkipAllTour{TourKey: tour.TourKey, Version: tour.Version})
	}

	if err := h.uc.SkipAll(c.Request.Context(), dtos.SkipAllInput{
		UserID:     userID,
		BusinessID: businessID,
		Tours:      tours,
	}); err != nil {
		h.log.Error(c.Request.Context()).Err(err).Uint("user_id", userID).Msg("Error al omitir los tutoriales")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
