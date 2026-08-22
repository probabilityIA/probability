package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orderscompare/internal/domain/dtos"
)

func (h *Handlers) Compare(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c, nil)
	if !ok {
		return
	}

	integrationID, err := strconv.ParseUint(c.Query("integration_id"), 10, 64)
	if err != nil || integrationID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "integration_id es requerido"})
		return
	}

	query := dtos.CompareQuery{
		BusinessID:    businessID,
		IntegrationID: uint(integrationID),
		From:          parseDate(c.Query("from"), false),
		To:            parseDate(c.Query("to"), true),
		Limit:         parseInt(c.Query("limit")),
		Page:          parseInt(c.Query("page")),
		PageSize:      parseInt(c.Query("page_size")),
		OnlyDiff:      c.Query("only_diff") == "true",
		Search:        c.Query("q"),
	}

	page, err := h.useCase.Compare(c.Request.Context(), query)
	if err != nil {
		h.logger.Error(c.Request.Context()).Err(err).
			Uint("integration_id", query.IntegrationID).
			Uint("business_id", businessID).
			Msg("No se pudo comparar las ordenes del canal")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": page})
}

func parseDate(value string, endOfDay bool) *time.Time {
	if value == "" {
		return nil
	}
	if parsed, err := time.Parse(time.RFC3339, value); err == nil {
		return &parsed
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil
	}
	if endOfDay {
		parsed = parsed.Add(24*time.Hour - time.Second)
	}
	return &parsed
}

func parseInt(value string) int {
	if value == "" {
		return 0
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return parsed
}
