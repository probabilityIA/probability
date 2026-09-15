package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/infra/primary/handlers/mappers"
)

const reviewDateLayout = "2006-01-02"

func (h *Handlers) ListAssistantReviewMessages(c *gin.Context) {
	filter, ok := reviewFilter(c)
	if !ok {
		return
	}

	page, err := h.uc.ListReviewMessages(c.Request.Context(), filter)
	if err != nil {
		h.log.Error(c.Request.Context()).Err(err).Msg("[ai.assistant] error listando conversaciones")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "code": "internal_error", "message": "No se pudieron cargar las conversaciones."})
		return
	}

	out := mappers.FromReviewMessages(page)
	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"data":        out.Data,
		"total":       out.Total,
		"page":        out.Page,
		"page_size":   out.PageSize,
		"total_pages": out.TotalPages,
	})
}

func (h *Handlers) GetAssistantReviewSummary(c *gin.Context) {
	filter, ok := reviewFilter(c)
	if !ok {
		return
	}

	summary, err := h.uc.GetReviewSummary(c.Request.Context(), filter)
	if err != nil {
		h.log.Error(c.Request.Context()).Err(err).Msg("[ai.assistant] error calculando el resumen")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "code": "internal_error", "message": "No se pudo calcular el resumen."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": mappers.FromReviewSummary(summary)})
}

func reviewFilter(c *gin.Context) (dtos.ReviewFilter, bool) {
	filter := dtos.ReviewFilter{
		Kind:   c.Query("kind"),
		Search: c.Query("search"),
	}

	if raw := c.Query("business_id"); raw != "" {
		parsed, err := strconv.ParseUint(raw, 10, 64)
		if err != nil || parsed == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "code": "invalid_business", "message": "El negocio no es v\u00e1lido."})
			return filter, false
		}
		id := uint(parsed)
		filter.BusinessID = &id
	}

	if raw := c.Query("from"); raw != "" {
		from, err := time.ParseInLocation(reviewDateLayout, raw, time.Local)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "code": "invalid_date", "message": "La fecha inicial no es v\u00e1lida (AAAA-MM-DD)."})
			return filter, false
		}
		filter.From = &from
	}

	if raw := c.Query("to"); raw != "" {
		to, err := time.ParseInLocation(reviewDateLayout, raw, time.Local)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "code": "invalid_date", "message": "La fecha final no es v\u00e1lida (AAAA-MM-DD)."})
			return filter, false
		}
		end := to.AddDate(0, 0, 1)
		filter.To = &end
	}

	filter.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	filter.PageSize, _ = strconv.Atoi(c.DefaultQuery("page_size", "20"))
	return filter, true
}
