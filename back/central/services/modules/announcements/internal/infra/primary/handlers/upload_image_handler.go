package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *handler) UploadImage(c *gin.Context) {
	id, ok := h.parseIDParam(c)
	if !ok {
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "image file is required"})
		return
	}

	sortOrder := 0
	if s := c.PostForm("sort_order"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			sortOrder = v
		}
	}

	img, err := h.uc.UploadImage(c.Request.Context(), id, file, sortOrder)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": gin.H{
			"id":        img.ID,
			"image_url": img.ImageURL,
			"sort_order": img.SortOrder,
		},
	})
}
