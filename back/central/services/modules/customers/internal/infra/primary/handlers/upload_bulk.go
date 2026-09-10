package handlers

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/infra/primary/handlers/response"
	"github.com/xuri/excelize/v2"
)

type bulkClientRow struct {
	rowNumber int
	name      string
	email     *string
	phone     string
	dni       *string
	address   *string
	city      *string
	notes     *string
}

func normalizeClientHeader(h string) string {
	s := strings.ToLower(strings.TrimSpace(h))
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.Trim(s, "\"'")
	return s
}

func (h *Handlers) UploadBulkClients(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Se requiere seleccionar un negocio (business_id)",
		})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "No se recibio un archivo valido",
			"error":   err.Error(),
		})
		return
	}
	defer file.Close()

	filename := header.Filename
	ext := strings.ToLower(filename[strings.LastIndex(filename, ".")+1:])

	var rows []bulkClientRow
	var parseErr error

	switch ext {
	case "csv":
		rows, parseErr = parseClientsCSV(file)
	case "xlsx", "xls":
		rows, parseErr = parseClientsExcel(file)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Formato no soportado. Use CSV o Excel (.xlsx, .xls)",
		})
		return
	}

	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("Error al parsear el archivo: %s", parseErr.Error()),
			"error":   parseErr.Error(),
		})
		return
	}

	result := response.BulkClientResultResponse{
		TotalRows: len(rows),
		Results:   []response.BulkClientRowResponse{},
	}

	for _, row := range rows {
		dto := dtos.CreateClientDTO{
			BusinessID: businessID,
			Name:       row.name,
			Email:      row.email,
			Phone:      row.phone,
			Dni:        row.dni,
			Address:    row.address,
			City:       row.city,
			Notes:      row.notes,
		}

		_, err := h.uc.CreateClient(c.Request.Context(), dto)
		if err != nil {
			result.FailedCount++
			result.Results = append(result.Results, response.BulkClientRowResponse{
				Row:     row.rowNumber,
				Name:    row.name,
				Success: false,
				Error:   err.Error(),
			})
			continue
		}
		result.SuccessCount++
		result.Results = append(result.Results, response.BulkClientRowResponse{
			Row:     row.rowNumber,
			Name:    row.name,
			Success: true,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Carga completada: %d exitosos, %d fallidos de %d totales", result.SuccessCount, result.FailedCount, result.TotalRows),
		"data":    result,
	})
}

func buildClientRowsFromRecords(records [][]string) ([]bulkClientRow, error) {
	if len(records) < 2 {
		return nil, fmt.Errorf("el archivo esta vacio o solo contiene encabezados")
	}

	headerMap := make(map[string]int)
	for i, head := range records[0] {
		headerMap[normalizeClientHeader(head)] = i
	}

	required := []string{"nombre", "apellido", "cedula", "correo", "telefono"}
	for _, col := range required {
		if _, exists := headerMap[col]; !exists {
			return nil, fmt.Errorf("falta la columna requerida: %s", col)
		}
	}

	rows := []bulkClientRow{}
	for i, record := range records[1:] {
		getValue := func(key string) string {
			idx, ok := headerMap[key]
			if !ok || idx >= len(record) {
				return ""
			}
			return strings.TrimSpace(record[idx])
		}

		nombre := getValue("nombre")
		apellido := getValue("apellido")
		if nombre == "" && apellido == "" {
			continue
		}

		fullName := strings.TrimSpace(nombre + " " + apellido)

		row := bulkClientRow{
			rowNumber: i + 2,
			name:      fullName,
			phone:     getValue("telefono"),
		}

		if v := getValue("correo"); v != "" {
			row.email = &v
		}
		if v := getValue("cedula"); v != "" {
			row.dni = &v
		}
		if v := getValue("direccion"); v != "" {
			row.address = &v
		}
		if v := getValue("ciudad"); v != "" {
			row.city = &v
		}
		if v := getValue("notas"); v != "" {
			row.notes = &v
		}

		rows = append(rows, row)
	}

	return rows, nil
}

func parseClientsCSV(file io.Reader) ([]bulkClientRow, error) {
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("fallo al leer bytes: %w", err)
	}

	content = bytes.TrimPrefix(content, []byte("\xef\xbb\xbf"))

	delimiter := ','
	firstLine := ""
	scanner := bufio.NewScanner(bytes.NewReader(content))
	if scanner.Scan() {
		firstLine = scanner.Text()
	}
	if strings.Count(firstLine, ";") > strings.Count(firstLine, ",") {
		delimiter = ';'
	}

	reader := csv.NewReader(bytes.NewReader(content))
	reader.Comma = delimiter
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error de formato CSV: %w", err)
	}

	return buildClientRowsFromRecords(records)
}

func parseClientsExcel(file io.Reader) ([]bulkClientRow, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("error al abrir Excel: %w", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("el Excel no tiene hojas")
	}

	records, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("error al leer celdas: %w", err)
	}

	return buildClientRowsFromRecords(records)
}
