package handlers

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/xuri/excelize/v2"
)

type BulkUploadResult struct {
	TotalRows    int      `json:"total_rows"`
	SuccessCount int      `json:"success_count"`
	FailedCount  int      `json:"failed_count"`
	Errors       []string `json:"errors,omitempty"`
}

// UploadBulkOrders godoc
// @Summary      Carga masiva de órdenes
// @Description  Carga múltiples órdenes desde un archivo CSV o Excel
// @Tags         Orders
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "Archivo CSV o Excel"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /orders/upload-bulk [post]
func (h *Handlers) UploadBulkOrders(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "No se pudo leer el archivo",
			"error":   err.Error(),
		})
		return
	}
	defer file.Close()

	// Determine file type
	filename := header.Filename
	ext := strings.ToLower(filename[strings.LastIndex(filename, ".")+1:])

	var orders []dtos.CreateOrderRequest
	var parseErr error

	switch ext {
	case "csv":
		orders, parseErr = h.parseCSV(file)
	case "xlsx", "xls":
		orders, parseErr = h.parseExcel(file)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Formato de archivo no soportado. Use CSV o Excel (.xlsx, .xls)",
		})
		return
	}

	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Error al parsear el archivo",
			"error":   parseErr.Error(),
		})
		return
	}

	// Process orders
	result := BulkUploadResult{
		TotalRows: len(orders),
		Errors:    []string{},
	}

	for i, orderReq := range orders {
		// Set platform to manual if not specified
		if orderReq.Platform == "" {
			orderReq.Platform = "manual"
		}

		_, err := h.orderCRUD.CreateOrder(c.Request.Context(), &orderReq)
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, fmt.Sprintf("Fila %d: %s", i+2, err.Error()))
		} else {
			result.SuccessCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Procesadas %d órdenes: %d exitosas, %d fallidas", result.TotalRows, result.SuccessCount, result.FailedCount),
		"data":    result,
	})
}

func (h *Handlers) parseCSV(file io.Reader) ([]dtos.CreateOrderRequest, error) {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("el archivo debe contener al menos una fila de encabezados y una fila de datos")
	}

	// Parse header
	headers := records[0]
	headerMap := make(map[string]int)
	for i, h := range headers {
		headerMap[strings.TrimSpace(strings.ToLower(h))] = i
	}

	// Required columns
	requiredCols := []string{"order_number", "customer_name", "customer_email", "customer_phone", "shipping_street", "shipping_city", "shipping_state", "total_amount"}
	for _, col := range requiredCols {
		if _, exists := headerMap[col]; !exists {
			return nil, fmt.Errorf("columna requerida faltante: %s", col)
		}
	}

	// Parse rows
	orders := []dtos.CreateOrderRequest{}
	for i, record := range records[1:] {
		if len(record) != len(headers) {
			return nil, fmt.Errorf("fila %d tiene un número incorrecto de columnas", i+2)
		}

		totalAmount, err := strconv.ParseFloat(record[headerMap["total_amount"]], 64)
		if err != nil {
			return nil, fmt.Errorf("fila %d: total_amount inválido", i+2)
		}

		order := dtos.CreateOrderRequest{
			OrderNumber:    record[headerMap["order_number"]],
			CustomerName:   record[headerMap["customer_name"]],
			CustomerEmail:  record[headerMap["customer_email"]],
			CustomerPhone:  record[headerMap["customer_phone"]],
			ShippingStreet: record[headerMap["shipping_street"]],
			ShippingCity:   record[headerMap["shipping_city"]],
			ShippingState:  record[headerMap["shipping_state"]],
			TotalAmount:    totalAmount,
			Platform:       "manual",
		}

		// Optional columns
		if idx, exists := headerMap["weight"]; exists && record[idx] != "" {
			if weight, err := strconv.ParseFloat(record[idx], 64); err == nil {
				order.Weight = &weight
			}
		}
		if idx, exists := headerMap["height"]; exists && record[idx] != "" {
			if height, err := strconv.ParseFloat(record[idx], 64); err == nil {
				order.Height = &height
			}
		}
		if idx, exists := headerMap["width"]; exists && record[idx] != "" {
			if width, err := strconv.ParseFloat(record[idx], 64); err == nil {
				order.Width = &width
			}
		}
		if idx, exists := headerMap["length"]; exists && record[idx] != "" {
			if length, err := strconv.ParseFloat(record[idx], 64); err == nil {
				order.Length = &length
			}
		}
		if idx, exists := headerMap["platform"]; exists && record[idx] != "" {
			order.Platform = record[idx]
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (h *Handlers) parseExcel(file io.Reader) ([]dtos.CreateOrderRequest, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Get first sheet
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("el archivo Excel no contiene hojas")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, err
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("el archivo debe contener al menos una fila de encabezados y una fila de datos")
	}

	// Parse header
	headers := rows[0]
	headerMap := make(map[string]int)
	for i, h := range headers {
		headerMap[strings.TrimSpace(strings.ToLower(h))] = i
	}

	// Required columns
	requiredCols := []string{"order_number", "customer_name", "customer_email", "customer_phone", "shipping_street", "shipping_city", "shipping_state", "total_amount"}
	for _, col := range requiredCols {
		if _, exists := headerMap[col]; !exists {
			return nil, fmt.Errorf("columna requerida faltante: %s", col)
		}
	}

	// Parse rows
	orders := []dtos.CreateOrderRequest{}
	for i, row := range rows[1:] {
		if len(row) < len(requiredCols) {
			return nil, fmt.Errorf("fila %d tiene un número incorrecto de columnas", i+2)
		}

		totalAmount, err := strconv.ParseFloat(row[headerMap["total_amount"]], 64)
		if err != nil {
			return nil, fmt.Errorf("fila %d: total_amount inválido", i+2)
		}

		order := dtos.CreateOrderRequest{
			OrderNumber:    row[headerMap["order_number"]],
			CustomerName:   row[headerMap["customer_name"]],
			CustomerEmail:  row[headerMap["customer_email"]],
			CustomerPhone:  row[headerMap["customer_phone"]],
			ShippingStreet: row[headerMap["shipping_street"]],
			ShippingCity:   row[headerMap["shipping_city"]],
			ShippingState:  row[headerMap["shipping_state"]],
			TotalAmount:    totalAmount,
			Platform:       "manual",
		}

		// Optional columns
		if idx, exists := headerMap["weight"]; exists && len(row) > idx && row[idx] != "" {
			if weight, err := strconv.ParseFloat(row[idx], 64); err == nil {
				order.Weight = &weight
			}
		}
		if idx, exists := headerMap["height"]; exists && len(row) > idx && row[idx] != "" {
			if height, err := strconv.ParseFloat(row[idx], 64); err == nil {
				order.Height = &height
			}
		}
		if idx, exists := headerMap["width"]; exists && len(row) > idx && row[idx] != "" {
			if width, err := strconv.ParseFloat(row[idx], 64); err == nil {
				order.Width = &width
			}
		}
		if idx, exists := headerMap["length"]; exists && len(row) > idx && row[idx] != "" {
			if length, err := strconv.ParseFloat(row[idx], 64); err == nil {
				order.Length = &length
			}
		}
		if idx, exists := headerMap["platform"]; exists && len(row) > idx && row[idx] != "" {
			order.Platform = row[idx]
		}

		orders = append(orders, order)
	}

	return orders, nil
}
