package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *handler) DeleteImage(c *gin.Context) {
	announcementID, ok := h.parseIDParam(c)
	if !ok {
		return
	}

	imageID, err := strconv.ParseUint(c.Param("imageId"), 10, 64)
	if err != nil || imageID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid image id"})
		return
	}

	if err := h.uc.DeleteImage(c.Request.Context(), announcementID, uint(imageID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "image deleted"})
}
