package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/domain"
)

type matrixColumnResponse struct {
	IntegrationID uint   `json:"integration_id"`
	Name          string `json:"name"`
	Code          string `json:"code"`
	ImageURL      string `json:"image_url,omitempty"`
	IsSales       bool   `json:"is_sales"`
}

type matrixCellResponse struct {
	IntegrationID uint   `json:"integration_id"`
	Present       bool   `json:"present"`
	SKU           string `json:"sku,omitempty"`
	Barcode       string `json:"barcode,omitempty"`
	ExternalID    string `json:"external_id,omitempty"`
	VariantID     string `json:"variant_id,omitempty"`
	SKUMatches    bool   `json:"sku_matches"`
	SKUKnown      bool   `json:"sku_known"`
}

type matrixRowResponse struct {
	ProductID string               `json:"product_id"`
	SKU       string               `json:"sku"`
	Barcode   string               `json:"barcode,omitempty"`
	Name      string               `json:"name,omitempty"`
	ImageURL  string               `json:"image_url,omitempty"`
	Cells     []matrixCellResponse `json:"cells"`
}

func (h *handler) MatchMatrix(c *gin.Context) {
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

	descarga := c.Query("format") == "csv"

	query := domain.MatrixQuery{
		BusinessID: businessID,
		Search:     strings.TrimSpace(c.Query("q")),
		SearchBy:   strings.TrimSpace(c.Query("search_by")),
		All:        descarga,
	}
	query.Page, _ = strconv.Atoi(c.Query("page"))
	query.PageSize, _ = strconv.Atoi(c.Query("page_size"))
	query.PresentIn = uintListQuery(c, "present_in")
	query.MissingIn = uintListQuery(c, "missing_in")

	page, err := h.useCase.MatchMatrix(c.Request.Context(), query)
	if err != nil {
		h.logger.Error(c.Request.Context()).Err(err).Uint("business_id", businessID).
			Msg("Error al construir el resumen de match")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "no se pudo obtener el resumen"})
		return
	}

	if descarga {
		h.escribirMatrizCSV(c, page)
		return
	}

	columnas := make([]matrixColumnResponse, 0, len(page.Columns))
	for _, col := range page.Columns {
		columnas = append(columnas, matrixColumnResponse{
			IntegrationID: col.IntegrationID, Name: col.Name, Code: col.Code,
			ImageURL: h.urlAbsoluta(col.ImageURL), IsSales: col.IsSales,
		})
	}

	filas := make([]matrixRowResponse, 0, len(page.Rows))
	for _, row := range page.Rows {
		celdas := make([]matrixCellResponse, 0, len(row.Cells))
		for _, cell := range row.Cells {
			celdas = append(celdas, matrixCellResponse{
				IntegrationID: cell.IntegrationID, Present: cell.Present,
				SKU: cell.SKU, Barcode: cell.Barcode,
				ExternalID: cell.ExternalID, VariantID: cell.VariantID,
				SKUMatches: cell.SKUMatches, SKUKnown: cell.SKUKnown,
			})
		}
		filas = append(filas, matrixRowResponse{
			ProductID: row.ProductID, SKU: row.SKU, Barcode: row.Barcode,
			Name: row.Name, ImageURL: h.urlAbsoluta(row.ImageURL), Cells: celdas,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"columns":     columnas,
		"data":        filas,
		"total":       page.Total,
		"page":        page.Page,
		"page_size":   page.PageSize,
		"total_pages": page.TotalPages,
	})
}

func (h *handler) urlAbsoluta(path string) string {
	if path == "" || strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	base := strings.TrimRight(os.Getenv("URL_BASE_DOMAIN_S3"), "/")
	if base == "" {
		return path
	}
	return base + "/" + strings.TrimLeft(path, "/")
}

func uintListQuery(c *gin.Context, key string) []uint {
	raw := strings.TrimSpace(c.Query(key))
	if raw == "" {
		return nil
	}
	out := make([]uint, 0, 4)
	for _, part := range strings.Split(raw, ",") {
		parsed, err := strconv.ParseUint(strings.TrimSpace(part), 10, 64)
		if err != nil || parsed == 0 {
			continue
		}
		out = append(out, uint(parsed))
	}
	return out
}

func (h *handler) escribirMatrizCSV(c *gin.Context, page *domain.MatrixPage) {
	nombre := fmt.Sprintf("resumen-match-%s.csv", time.Now().Format("2006-01-02"))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", `attachment; filename="`+nombre+`"`)

	if _, err := c.Writer.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
		return
	}

	w := csv.NewWriter(c.Writer)
	defer w.Flush()

	encabezado := []string{"SKU Probability", "Codigo de barras", "Producto", "ID Probability"}
	for _, col := range page.Columns {
		encabezado = append(encabezado,
			col.Name+" - SKU", col.Name+" - codigo de barras", col.Name+" - ID")
	}
	_ = w.Write(encabezado)

	for _, row := range page.Rows {
		fila := []string{row.SKU, row.Barcode, row.Name, row.ProductID}
		for _, cell := range row.Cells {
			if !cell.Present {
				fila = append(fila, "no esta", "", "")
				continue
			}
			fila = append(fila, cell.SKU, cell.Barcode, idDeCanal(cell))
		}
		_ = w.Write(fila)
	}
}

func idDeCanal(cell domain.MatrixCell) string {
	if cell.VariantID != "" {
		return cell.ExternalID + " / " + cell.VariantID
	}
	return cell.ExternalID
}
