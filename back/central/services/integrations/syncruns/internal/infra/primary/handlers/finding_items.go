package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/domain"
)

type findingItemResponse struct {
	SKU      string   `json:"sku"`
	Name     string   `json:"name,omitempty"`
	Detail   string   `json:"detail,omitempty"`
	Channels []string `json:"channels,omitempty"`
}

func (h *handler) FindingItems(c *gin.Context) {
	var provided *uint
	if raw := c.Query("business_id"); raw != "" {
		if parsed, err := strconv.ParseUint(raw, 10, 64); err == nil {
			value := uint(parsed)
			provided = &value
		}
	}

	businessID, ok := h.resolveBusinessID(c, provided)
	if !ok {
		return
	}

	code := strings.TrimSpace(c.Query("code"))
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "code es requerido"})
		return
	}

	descarga := c.Query("format") == "csv"

	query := domain.FindingItemsQuery{
		BusinessID: businessID,
		Code:       code,
		Search:     strings.TrimSpace(c.Query("q")),
		All:        descarga,
	}
	query.Page, _ = strconv.Atoi(c.Query("page"))
	query.PageSize, _ = strconv.Atoi(c.Query("page_size"))

	page, err := h.useCase.FindingItems(c.Request.Context(), query)
	if err != nil {
		h.logger.Error(c.Request.Context()).Err(err).
			Uint("business_id", businessID).Str("code", code).
			Msg("Error al listar el detalle del hallazgo")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "no se pudo obtener el detalle"})
		return
	}

	if descarga {
		h.escribirCSV(c, code, page)
		return
	}

	items := make([]findingItemResponse, 0, len(page.Items))
	for _, item := range page.Items {
		items = append(items, findingItemResponse{
			SKU: item.SKU, Name: item.Name, Detail: item.Detail, Channels: item.Channels,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"data":        items,
		"total":       page.Total,
		"page":        page.Page,
		"page_size":   page.PageSize,
		"total_pages": page.TotalPages,
	})
}

func (h *handler) escribirCSV(c *gin.Context, code string, page *domain.FindingItemsPage) {
	nombre := fmt.Sprintf("hallazgos-%s-%s.csv", code, time.Now().Format("2006-01-02"))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", `attachment; filename="`+nombre+`"`)

	// BOM para que Excel abra los acentos bien.
	if _, err := c.Writer.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
		return
	}

	w := csv.NewWriter(c.Writer)
	defer w.Flush()

	_ = w.Write([]string{"SKU", "Producto", "Detalle", "Canales"})
	for _, item := range page.Items {
		_ = w.Write([]string{item.SKU, item.Name, item.Detail, strings.Join(item.Channels, ", ")})
	}
}
