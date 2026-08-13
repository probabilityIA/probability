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
	SKU             string   `json:"sku"`
	Name            string   `json:"name,omitempty"`
	Detail          string   `json:"detail,omitempty"`
	Channels        []string `json:"channels,omitempty"`
	CounterpartSKU  string   `json:"counterpart_sku,omitempty"`
	CounterpartName string   `json:"counterpart_name,omitempty"`
	ChannelQty      *int     `json:"channel_qty,omitempty"`
	OwnQty          *int     `json:"own_qty,omitempty"`
	FixSide         string   `json:"fix_side,omitempty"`
	Pattern         string   `json:"pattern,omitempty"`
	PresentIn       []string `json:"present_in,omitempty"`
	MissingIn       []string `json:"missing_in,omitempty"`
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
			CounterpartSKU: item.CounterpartSKU, CounterpartName: item.CounterpartName,
			ChannelQty: item.ChannelQty, OwnQty: item.OwnQty,
			FixSide: item.FixSide, Pattern: item.Pattern,
			PresentIn: item.PresentIn, MissingIn: item.MissingIn,
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

	_ = w.Write([]string{
		"SKU en el canal", "SKU en Probability", "Producto",
		"Stock en el canal", "Stock en Probability", "Que hacer",
		"Esta en", "Falta en", "Canales",
	})
	for _, item := range page.Items {
		_ = w.Write([]string{
			item.SKU, item.CounterpartSKU, item.Name,
			numero(item.ChannelQty), numero(item.OwnQty),
			item.Detail,
			strings.Join(item.PresentIn, ", "), strings.Join(item.MissingIn, ", "),
			strings.Join(item.Channels, ", "),
		})
	}
}

func numero(value *int) string {
	if value == nil {
		return ""
	}
	return strconv.Itoa(*value)
}
