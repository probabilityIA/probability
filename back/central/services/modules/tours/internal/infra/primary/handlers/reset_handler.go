package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/secamc93/probability/back/central/services/modules/tours/internal/domain/dtos"
)

func (h *Handler) ResetTour(c *gin.Context) {
	userID, businessID, ok := h.resolveIdentity(c)
	if !ok {
		return
	}

	tourKey := c.Param("tour_key")
	if err := h.uc.ResetTour(c.Request.Context(), dtos.ResetInput{
		UserID:     userID,
		BusinessID: businessID,
		TourKey:    tourKey,
	}); err != nil {
		h.log.Error(c.Request.Context()).Err(err).Uint("user_id", userID).Str("tour_key", tourKey).Msg("Error al reiniciar el tour")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) ResetAll(c *gin.Context) {
	userID, businessID, ok := h.resolveIdentity(c)
	if !ok {
		return
	}

	if err := h.uc.ResetAll(c.Request.Context(), dtos.ListProgressInput{
		UserID:     userID,
		BusinessID: businessID,
	}); err != nil {
		h.log.Error(c.Request.Context()).Err(err).Uint("user_id", userID).Msg("Error al reiniciar los tours")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
