package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
	"github.com/xuri/excelize/v2"
)

func normalizeDimensionHeader(h string) string {
	s := strings.ToLower(strings.TrimSpace(h))
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")
	return s
}

func parseDimensionFloat(val string) (*float64, error) {
	s := strings.TrimSpace(val)
	if s == "" {
		return nil, nil
	}
	s = strings.ReplaceAll(s, ",", ".")
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (h *Handlers) ImportProductsDimensions(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		h.respondBusinessIDRequired(c)
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "No se recibió un archivo válido",
			"error":   err.Error(),
		})
		return
	}
	defer file.Close()

	rows, err := parseDimensionRows(file, true)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("Error al leer el archivo: %s", err.Error()),
			"error":   err.Error(),
		})
		return
	}

	result, err := h.uc.BulkUpdateDimensions(c.Request.Context(), businessID, rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al procesar el archivo",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Carga completada: %d actualizados, %d fallidos de %d filas", result.SuccessCount, result.FailedCount, result.TotalRows),
		"data":    result,
	})
}

func parseDimensionRows(file io.Reader, requireDimension bool) ([]domain.DimensionRow, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("error al abrir el Excel: %w", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("el Excel no tiene hojas")
	}

	allRows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("error al leer celdas: %w", err)
	}
	if len(allRows) < 2 {
		return nil, fmt.Errorf("el archivo está vacío o solo contiene encabezados")
	}

	headerMap := make(map[string]int)
	for i, head := range allRows[0] {
		headerMap[normalizeDimensionHeader(head)] = i
	}

	skuIdx, ok := headerMap["sku"]
	if !ok {
		return nil, fmt.Errorf("falta la columna requerida: sku")
	}

	weightIdx, hasWeight := firstHeaderIndex(headerMap, "peso_kg", "weight")
	lengthIdx, hasLength := firstHeaderIndex(headerMap, "largo_cm", "length")
	widthIdx, hasWidth := firstHeaderIndex(headerMap, "ancho_cm", "width")
	heightIdx, hasHeight := firstHeaderIndex(headerMap, "alto_cm", "height")

	getValue := func(row []string, idx int, has bool) string {
		if !has || idx >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[idx])
	}

	out := make([]domain.DimensionRow, 0, len(allRows)-1)
	for _, row := range allRows[1:] {
		sku := ""
		if skuIdx < len(row) {
			sku = strings.TrimSpace(row[skuIdx])
		}
		if sku == "" {
			continue
		}

		weight, err := parseDimensionFloat(getValue(row, weightIdx, hasWeight))
		if err != nil {
			return nil, fmt.Errorf("SKU %s: peso inválido", sku)
		}
		length, err := parseDimensionFloat(getValue(row, lengthIdx, hasLength))
		if err != nil {
			return nil, fmt.Errorf("SKU %s: largo inválido", sku)
		}
		width, err := parseDimensionFloat(getValue(row, widthIdx, hasWidth))
		if err != nil {
			return nil, fmt.Errorf("SKU %s: ancho inválido", sku)
		}
		height, err := parseDimensionFloat(getValue(row, heightIdx, hasHeight))
		if err != nil {
			return nil, fmt.Errorf("SKU %s: alto inválido", sku)
		}

		if requireDimension && weight == nil && length == nil && width == nil && height == nil {
			continue
		}

		out = append(out, domain.DimensionRow{
			SKU:    sku,
			Weight: weight,
			Length: length,
			Width:  width,
			Height: height,
		})
	}

	return out, nil
}

func firstHeaderIndex(headerMap map[string]int, keys ...string) (int, bool) {
	for _, k := range keys {
		if idx, ok := headerMap[k]; ok {
			return idx, true
		}
	}
	return 0, false
}
