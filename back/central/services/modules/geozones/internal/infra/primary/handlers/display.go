package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/dtos"
)

func parseBbox(s string) *dtos.Bbox {
	parts := strings.Split(s, ",")
	if len(parts) != 4 {
		return nil
	}
	vals := make([]float64, 4)
	for i, p := range parts {
		v, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
		if err != nil {
			return nil
		}
		vals[i] = v
	}
	return &dtos.Bbox{MinLng: vals[0], MinLat: vals[1], MaxLng: vals[2], MaxLat: vals[3]}
}

func (h *Handlers) Display(c *gin.Context) {
	geozoneType := c.Query("type")
	zoom, _ := strconv.Atoi(c.DefaultQuery("zoom", "6"))
	bbox := parseBbox(c.Query("bbox"))

	var parentID *uint
	if v := c.Query("parent_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 64); err == nil && id > 0 {
			p := uint(id)
			parentID = &p
		}
	}

	if bbox == nil {
		if ifNoneMatch := c.GetHeader("If-None-Match"); ifNoneMatch != "" {
			expected := buildETagFromQuery(geozoneType, zoom, parentID)
			if ifNoneMatch == expected {
				c.Header("ETag", expected)
				c.Header("Cache-Control", "public, max-age=86400")
				c.Status(http.StatusNotModified)
				return
			}
		}
	}

	payload, etag, err := h.uc.GetForDisplay(c.Request.Context(), geozoneType, zoom, bbox, parentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if bbox == nil {
		c.Header("ETag", etag)
		c.Header("Cache-Control", "no-cache")
	} else {
		c.Header("Cache-Control", "no-cache")
	}
	c.Header("Vary", "Accept-Encoding")
	c.Data(http.StatusOK, "application/json; charset=utf-8", payload)
}

func (h *Handlers) FlushDisplayCache(c *gin.Context) {
	if err := h.uc.FlushDisplayCache(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "cache de display invalidado"})
}

func buildETagFromQuery(geozoneType string, zoom int, parentID *uint) string {
	t := geozoneType
	if t == "" {
		t = "all"
	}
	bucket := "zfull"
	switch {
	case zoom <= 7:
		bucket = "z7"
	case zoom <= 9:
		bucket = "z9"
	}
	parentSeg := "p0"
	if parentID != nil {
		parentSeg = strconv.FormatUint(uint64(*parentID), 10)
		parentSeg = "p" + parentSeg
	}
	return `"` + t + "-" + parentSeg + "-" + bucket + `-v7"`
}
