package handlers

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/xuri/excelize/v2"
)

type BulkUploadResult struct {
	TotalRows    int      `json:"total_rows"`
	SuccessCount int      `json:"success_count"`
	FailedCount  int      `json:"failed_count"`
	Errors       []string `json:"errors,omitempty"`
}

// normalizeHeader limpia encabezados para matcheo flexible (ignora espacios, guiones, mayúsculas)
func normalizeHeader(h string) string {
	s := strings.ToLower(strings.TrimSpace(h))
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.Trim(s, "\"'")
	return s
}

// parseRobustFloat maneja formatos como 1.200,50 o 1,200.50 o 50.00
func parseRobustFloat(val string) (float64, error) {
	s := strings.TrimSpace(val)
	if s == "" {
		return 0, nil
	}

	// Limpieza: solo permitir números, puntos, comas y signo menos
	s = strings.Map(func(r rune) rune {
		if (r >= '0' && r <= '9') || r == '.' || r == ',' || r == '-' {
			return r
		}
		return -1
	}, s)

	// Si tiene tanto coma como punto, eliminamos el que aparece primero (asumiendo que es separador de miles)
	if strings.Contains(s, ",") && strings.Contains(s, ".") {
		commaIdx := strings.Index(s, ",")
		dotIdx := strings.Index(s, ".")
		if commaIdx < dotIdx {
			s = strings.Replace(s, ",", "", 1)
		} else {
			s = strings.Replace(s, ".", "", 1)
		}
	}

	// Convertir coma decimal a punto
	s = strings.ReplaceAll(s, ",", ".")

	// Si quedaron múltiples puntos (miles + decimal), quitar todos menos el último
	if strings.Count(s, ".") > 1 {
		lastIdx := strings.LastIndex(s, ".")
		prefix := strings.ReplaceAll(s[:lastIdx], ".", "")
		s = prefix + "." + s[lastIdx+1:]
	}

	return strconv.ParseFloat(s, 64)
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
	// Obtener BusinessID del contexto de autenticación
	businessID, exists := middleware.GetBusinessID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "No se pudo identificar la empresa (Sesión expirada o inválida)",
		})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "No se recibió un archivo válido",
			"error":   err.Error(),
		})
		return
	}
	defer file.Close()

	// Tipo de archivo
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
			"message": "Formato no soportado. Use CSV o Excel (.xlsx, .xls)",
		})
		return
	}

	if parseErr != nil {
		// Loguear el error detallado en el servidor para el desarrollador
		fmt.Printf("[UPLOAD-BULK] ERROR parseando %s: %v\n", filename, parseErr)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("Error al parsear el archivo: %s", parseErr.Error()),
			"error":   parseErr.Error(),
		})
		return
	}

	// Procesar órdenes
	result := BulkUploadResult{
		TotalRows: len(orders),
		Errors:    []string{},
	}

	for i, orderReq := range orders {
		// Inyectar contexto
		orderReq.BusinessID = &businessID
		if orderReq.Platform == "" {
			orderReq.Platform = "manual"
		}

		_, err := h.orderCRUD.CreateOrder(c.Request.Context(), &orderReq)
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, fmt.Sprintf("Fila %d (Pedido %s): %s", i+2, orderReq.OrderNumber, err.Error()))
		} else {
			result.SuccessCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Carga completada: %d exitosas, %d fallidas de %d totales", result.SuccessCount, result.FailedCount, result.TotalRows),
		"data":    result,
	})
}

func (h *Handlers) parseCSV(file io.Reader) ([]dtos.CreateOrderRequest, error) {
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("fallo al leer bytes: %w", err)
	}

	// Quitar BOM si existe
	content = bytes.TrimPrefix(content, []byte("\xef\xbb\xbf"))

	// Detectar delimitador
	delimiter := ','
	firstLine := ""
	scanner := bufio.NewScanner(bytes.NewReader(content))
	if scanner.Scan() {
		firstLine = scanner.Text()
	}
	// Si hay más puntos y coma que comas, probablemente es Excel Spanish locale
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

	if len(records) < 2 {
		return nil, fmt.Errorf("el archivo está vacío o solo contiene encabezados")
	}

	// Mapear encabezados normalizados
	headerMap := make(map[string]int)
	for i, head := range records[0] {
		headerMap[normalizeHeader(head)] = i
	}

	// Columnas requeridas
	required := []string{"order_number", "customer_name", "customer_email", "customer_phone", "shipping_street", "shipping_city", "shipping_state", "total_amount"}
	for _, col := range required {
		if _, exists := headerMap[col]; !exists {
			return nil, fmt.Errorf("falta la columna requerida: %s", col)
		}
	}

	orders := []dtos.CreateOrderRequest{}
	for i, record := range records[1:] {
		getValue := func(key string) string {
			idx, ok := headerMap[key]
			if !ok || idx >= len(record) {
				return ""
			}
			return strings.TrimSpace(record[idx])
		}

		// Saltar si no hay datos mínimos
		if getValue("order_number") == "" && getValue("customer_name") == "" {
			continue
		}

		total, err := parseRobustFloat(getValue("total_amount"))
		if err != nil {
			return nil, fmt.Errorf("fila %d: monto total inválido '%s'", i+2, getValue("total_amount"))
		}

		order := dtos.CreateOrderRequest{
			ExternalID:     getValue("order_number"),
			OrderNumber:    getValue("order_number"),
			CustomerName:   getValue("customer_name"),
			CustomerEmail:  getValue("customer_email"),
			CustomerPhone:  getValue("customer_phone"),
			ShippingStreet: getValue("shipping_street"),
			ShippingCity:   getValue("shipping_city"),
			ShippingState:  getValue("shipping_state"),
			TotalAmount:    total,
			Platform:       "manual",
		}

		// Opcionales
		if v, err := parseRobustFloat(getValue("weight")); err == nil && getValue("weight") != "" {
			order.Weight = &v
		}
		if v, err := parseRobustFloat(getValue("height")); err == nil && getValue("height") != "" {
			order.Height = &v
		}
		if v, err := parseRobustFloat(getValue("width")); err == nil && getValue("width") != "" {
			order.Width = &v
		}
		if v, err := parseRobustFloat(getValue("length")); err == nil && getValue("length") != "" {
			order.Length = &v
		}
		if p := getValue("platform"); p != "" {
			order.Platform = p
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (h *Handlers) parseExcel(file io.Reader) ([]dtos.CreateOrderRequest, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("error al abrir Excel: %w", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("el Excel no tiene hojas")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("error al leer celdas: %w", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("archivo Excel vacío")
	}

	// Mapear encabezados
	headerMap := make(map[string]int)
	for i, head := range rows[0] {
		headerMap[normalizeHeader(head)] = i
	}

	required := []string{"order_number", "customer_name", "customer_email", "customer_phone", "shipping_street", "shipping_city", "shipping_state", "total_amount"}
	for _, col := range required {
		if _, exists := headerMap[col]; !exists {
			return nil, fmt.Errorf("falta la columna requerida: %s", col)
		}
	}

	orders := []dtos.CreateOrderRequest{}
	for i, row := range rows[1:] {
		getValue := func(key string) string {
			idx, ok := headerMap[key]
			if !ok || idx >= len(row) {
				return ""
			}
			return strings.TrimSpace(row[idx])
		}

		if getValue("order_number") == "" && getValue("customer_name") == "" {
			continue
		}

		total, err := parseRobustFloat(getValue("total_amount"))
		if err != nil {
			return nil, fmt.Errorf("fila %d: monto total inválido '%s'", i+2, getValue("total_amount"))
		}

		order := dtos.CreateOrderRequest{
			ExternalID:     getValue("order_number"),
			OrderNumber:    getValue("order_number"),
			CustomerName:   getValue("customer_name"),
			CustomerEmail:  getValue("customer_email"),
			CustomerPhone:  getValue("customer_phone"),
			ShippingStreet: getValue("shipping_street"),
			ShippingCity:   getValue("shipping_city"),
			ShippingState:  getValue("shipping_state"),
			TotalAmount:    total,
			Platform:       "manual",
		}

		if v, err := parseRobustFloat(getValue("weight")); err == nil && getValue("weight") != "" {
			order.Weight = &v
		}
		if v, err := parseRobustFloat(getValue("height")); err == nil && getValue("height") != "" {
			order.Height = &v
		}
		if v, err := parseRobustFloat(getValue("width")); err == nil && getValue("width") != "" {
			order.Width = &v
		}
		if v, err := parseRobustFloat(getValue("length")); err == nil && getValue("length") != "" {
			order.Length = &v
		}
		if p := getValue("platform"); p != "" {
			order.Platform = p
		}

		orders = append(orders, order)
	}

	return orders, nil
}
