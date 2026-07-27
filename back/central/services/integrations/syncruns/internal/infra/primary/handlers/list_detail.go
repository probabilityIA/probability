package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/domain"
)

func (h *handler) ListDetail(c *gin.Context) {
	integrationID, err := strconv.ParseUint(c.Query("integration_id"), 10, 64)
	if err != nil || integrationID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "integration_id es requerido"})
		return
	}

	var provided *uint
	if raw := c.Query("business_id"); raw != "" {
		if parsed, parseErr := strconv.ParseUint(raw, 10, 64); parseErr == nil {
			value := uint(parsed)
			provided = &value
		}
	}

	businessID, ok := h.resolveBusinessID(c, provided)
	if !ok {
		return
	}

	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))

	query := domain.DetailQuery{
		BusinessID:    businessID,
		IntegrationID: uint(integrationID),
		Kind:          strings.TrimSpace(c.Query("kind")),
		Group:         strings.TrimSpace(c.Query("group")),
		Search:        strings.TrimSpace(c.Query("q")),
		Page:          page,
		PageSize:      pageSize,
	}

	result, err := h.useCase.ListDetail(c.Request.Context(), query)
	if err != nil {
		if errors.Is(err, domain.ErrIntegrationNotOwned) {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "la integracion no pertenece al negocio"})
			return
		}
		h.logger.Error(c.Request.Context()).Err(err).
			Uint64("integration_id", integrationID).
			Msg("Error al listar el detalle de la sincronizacion")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "no se pudo obtener el detalle"})
		return
	}

	items := make([]detailResponse, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, detailResponse{SKU: item.SKU, Label: item.Label, Tone: item.Tone, Group: item.Group})
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"data":        items,
		"total":       result.Total,
		"page":        result.Page,
		"page_size":   result.PageSize,
		"total_pages": result.TotalPages,
	})
}
