package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type updateSeenRequest struct {
	Seen bool `json:"seen"`
}

func (h *Handlers) UpdateSeen(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id invalido"})
		return
	}

	var req updateSeenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if err := h.uc.UpdateSeen(c.Request.Context(), uint(id), req.Seen); err != nil {
		h.logger.Error(c.Request.Context()).Err(err).Msg("Error al actualizar visto de prospecto comercial")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
