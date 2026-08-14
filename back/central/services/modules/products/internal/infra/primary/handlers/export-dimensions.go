package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func (h *Handlers) ExportProductsDimensions(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		h.respondBusinessIDRequired(c)
		return
	}

	products, err := h.uc.ExportDimensions(c.Request.Context(), businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al generar el archivo de dimensiones",
			"error":   err.Error(),
		})
		return
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Dimensiones"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"sku", "id", "nombre", "peso_kg", "largo_cm", "ancho_cm", "alto_cm"}
	for col, head := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheet, cell, head)
	}

	for i, p := range products {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), p.SKU)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), p.ID)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), p.Name)
		if p.Weight != nil {
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), *p.Weight)
		}
		if p.Length != nil {
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), *p.Length)
		}
		if p.Width != nil {
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), *p.Width)
		}
		if p.Height != nil {
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), *p.Height)
		}
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al construir el archivo Excel",
			"error":   err.Error(),
		})
		return
	}

	filename := "dimensiones-productos-" + strconv.FormatUint(uint64(businessID), 10) + "-" + time.Now().Format("2006-01-02") + ".xlsx"
	c.Header("Content-Disposition", `attachment; filename="`+filename+`"`)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}
