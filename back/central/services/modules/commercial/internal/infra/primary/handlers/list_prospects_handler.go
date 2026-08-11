package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/commercial/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/commercial/internal/infra/primary/handlers/response"
)

func (h *Handlers) ListProspects(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	filters := dtos.ListProspectsFilters{
		Platform: c.Query("platform"),
		Status:   c.Query("status"),
		City:     c.Query("city"),
		Search:   c.Query("search"),
		Page:     page,
		PageSize: pageSize,
	}

	if minScoreStr := c.Query("min_score"); minScoreStr != "" {
		if minScore, err := strconv.ParseFloat(minScoreStr, 64); err == nil {
			filters.MinScore = &minScore
		}
	}
	if maxScoreStr := c.Query("max_score"); maxScoreStr != "" {
		if maxScore, err := strconv.ParseFloat(maxScoreStr, 64); err == nil {
			filters.MaxScore = &maxScore
		}
	}
	if hasCODStr := c.Query("has_cod"); hasCODStr != "" {
		hasCOD := hasCODStr == "true"
		filters.HasCOD = &hasCOD
	}
	if seenStr := c.Query("seen"); seenStr != "" {
		seen := seenStr == "true"
		filters.Seen = &seen
	}

	prospects, result, err := h.uc.ListProspects(c.Request.Context(), filters)
	if err != nil {
		h.logger.Error(c.Request.Context()).Err(err).Msg("Error al listar prospectos comerciales")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"data":        response.FromEntities(prospects),
		"total":       result.Total,
		"page":        result.Page,
		"page_size":   result.PageSize,
		"total_pages": result.TotalPages,
	})
}
