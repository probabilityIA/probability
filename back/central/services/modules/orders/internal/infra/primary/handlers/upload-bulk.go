package handlers

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

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

type bulkItemJSON struct {
	SKU      string  `json:"sku"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Discount float64 `json:"discount,omitempty"`
	Tax      float64 `json:"tax,omitempty"`
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
	// Para usuarios normales: business_id del JWT.
	// Para super admin: business_id del query param ?business_id=X.
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

		_, err := h.createUC.CreateManualOrder(c.Request.Context(), &orderReq)
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

	return buildOrdersFromRows(records)
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

	return buildOrdersFromRows(rows)
}

// buildOrdersFromRows agrupa filas por order_number: cada fila es un producto,
// las filas con el mismo order_number se consolidan en una sola orden con
// varios items. Los datos de orden (cliente, envío, montos) se toman de la
// primera fila de cada grupo.
func buildOrdersFromRows(rows [][]string) ([]dtos.CreateOrderRequest, error) {
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

	hasProductColumns := false
	for _, col := range []string{"sku", "quantity", "unit_price"} {
		if _, exists := headerMap[col]; exists {
			hasProductColumns = true
			break
		}
	}

	orderIndex := make(map[string]int)
	orders := []dtos.CreateOrderRequest{}
	items := make(map[string][]bulkItemJSON)

	for i, row := range rows[1:] {
		getValue := func(key string) string {
			idx, ok := headerMap[key]
			if !ok || idx >= len(row) {
				return ""
			}
			return strings.TrimSpace(row[idx])
		}

		orderNumber := getValue("order_number")
		if orderNumber == "" && getValue("customer_name") == "" {
			continue
		}

		if _, exists := orderIndex[orderNumber]; exists {
			if hasProductColumns {
				if item, ok, err := extractItem(getValue, i+2); err != nil {
					return nil, err
				} else if ok {
					items[orderNumber] = append(items[orderNumber], item)
				}
			}
			continue
		}

		order, err := buildOrderFromRow(getValue, i+2)
		if err != nil {
			return nil, err
		}

		orderIndex[orderNumber] = len(orders)
		orders = append(orders, order)

		if hasProductColumns {
			if item, ok, err := extractItem(getValue, i+2); err != nil {
				return nil, err
			} else if ok {
				items[orderNumber] = append(items[orderNumber], item)
			}
		}
	}

	for orderNumber, idx := range orderIndex {
		lineItems, ok := items[orderNumber]
		if !ok || len(lineItems) == 0 {
			continue
		}
		raw, err := json.Marshal(lineItems)
		if err != nil {
			return nil, fmt.Errorf("pedido %s: error armando los productos: %w", orderNumber, err)
		}
		orders[idx].Items = raw
	}

	return orders, nil
}

func extractItem(getValue func(string) string, rowNum int) (bulkItemJSON, bool, error) {
	sku := getValue("sku")
	quantityRaw := getValue("quantity")
	priceRaw := getValue("unit_price")

	if sku == "" && quantityRaw == "" && priceRaw == "" {
		return bulkItemJSON{}, false, nil
	}
	if sku == "" {
		return bulkItemJSON{}, false, fmt.Errorf("fila %d: falta el sku del producto", rowNum)
	}

	quantity := 1
	if quantityRaw != "" {
		q, err := strconv.Atoi(strings.TrimSpace(quantityRaw))
		if err != nil || q <= 0 {
			return bulkItemJSON{}, false, fmt.Errorf("fila %d: cantidad inválida '%s'", rowNum, quantityRaw)
		}
		quantity = q
	}

	price, err := parseRobustFloat(priceRaw)
	if err != nil {
		return bulkItemJSON{}, false, fmt.Errorf("fila %d: precio unitario inválido '%s'", rowNum, priceRaw)
	}

	item := bulkItemJSON{
		SKU:      sku,
		Name:     getValue("product_name"),
		Price:    price,
		Quantity: quantity,
	}

	if discountRaw := getValue("product_discount"); discountRaw != "" {
		discount, err := parseRobustFloat(discountRaw)
		if err != nil {
			return bulkItemJSON{}, false, fmt.Errorf("fila %d: descuento de producto inválido '%s'", rowNum, discountRaw)
		}
		item.Discount = discount
	}
	if taxRaw := getValue("product_tax"); taxRaw != "" {
		tax, err := parseRobustFloat(taxRaw)
		if err != nil {
			return bulkItemJSON{}, false, fmt.Errorf("fila %d: impuesto de producto inválido '%s'", rowNum, taxRaw)
		}
		item.Tax = tax
	}

	return item, true, nil
}

func buildOrderFromRow(getValue func(string) string, rowNum int) (dtos.CreateOrderRequest, error) {
	total, err := parseRobustFloat(getValue("total_amount"))
	if err != nil {
		return dtos.CreateOrderRequest{}, fmt.Errorf("fila %d: monto total inválido '%s'", rowNum, getValue("total_amount"))
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

	// Opcionales - Dimensiones
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

	// Opcionales - Financiera
	if v, err := parseRobustFloat(getValue("subtotal")); err == nil && getValue("subtotal") != "" {
		order.Subtotal = v
	}
	if v, err := parseRobustFloat(getValue("tax")); err == nil && getValue("tax") != "" {
		order.Tax = v
	}
	if v, err := parseRobustFloat(getValue("discount")); err == nil && getValue("discount") != "" {
		order.Discount = v
	}
	if v, err := parseRobustFloat(getValue("shipping_cost")); err == nil && getValue("shipping_cost") != "" {
		order.ShippingCost = v
	}
	if v, err := parseRobustFloat(getValue("shipping_discount")); err == nil && getValue("shipping_discount") != "" {
		order.ShippingDiscount = v
	}
	if c := getValue("currency"); c != "" {
		order.Currency = c
	}

	// Opcionales - Cliente
	if fn := getValue("customer_first_name"); fn != "" {
		order.CustomerFirstName = fn
	}
	if ln := getValue("customer_last_name"); ln != "" {
		order.CustomerLastName = ln
	}
	if dni := getValue("customer_dni"); dni != "" {
		order.CustomerDNI = dni
	}

	// Opcionales - Dirección
	if c := getValue("shipping_country"); c != "" {
		order.ShippingCountry = c
	}
	if pc := getValue("shipping_postal_code"); pc != "" {
		order.ShippingPostalCode = pc
	}
	if lat, err := parseRobustFloat(getValue("shipping_lat")); err == nil && getValue("shipping_lat") != "" {
		order.ShippingLat = &lat
	}
	if lng, err := parseRobustFloat(getValue("shipping_lng")); err == nil && getValue("shipping_lng") != "" {
		order.ShippingLng = &lng
	}

	// Opcionales - Estado y Pago
	if s := getValue("status"); s != "" {
		order.Status = s
	}
	if pm := getValue("payment_method_id"); pm != "" {
		if id, err := strconv.ParseUint(pm, 10, 32); err == nil {
			order.PaymentMethodID = uint(id)
		}
	}
	if isPaid := getValue("is_paid"); isPaid != "" {
		order.IsPaid = strings.ToLower(isPaid) == "true" || isPaid == "1" || strings.ToLower(isPaid) == "yes" || strings.ToLower(isPaid) == "si"
	}

	// Opcionales - Logística
	if tn := getValue("tracking_number"); tn != "" {
		order.TrackingNumber = &tn
	}
	if gid := getValue("guide_id"); gid != "" {
		order.GuideID = &gid
	}
	if wh := getValue("warehouse_name"); wh != "" {
		order.WarehouseName = wh
	}
	if dr := getValue("driver_name"); dr != "" {
		order.DriverName = dr
	}

	// Opcionales - Adicional
	if notes := getValue("notes"); notes != "" {
		order.Notes = &notes
	}
	if ot := getValue("order_type_name"); ot != "" {
		order.OrderTypeName = ot
	}
	if inv := getValue("invoiceable"); inv != "" {
		order.Invoiceable = strings.ToLower(inv) == "true" || inv == "1" || strings.ToLower(inv) == "yes" || strings.ToLower(inv) == "si"
	}

	if p := getValue("platform"); p != "" {
		order.Platform = p
	}

	if od := getValue("order_date"); od != "" {
		occurredAt, err := parseFlexibleDate(od)
		if err != nil {
			return dtos.CreateOrderRequest{}, fmt.Errorf("fila %d: fecha de orden inválida '%s' (use AAAA-MM-DD)", rowNum, od)
		}
		order.OccurredAt = occurredAt
	}

	return order, nil
}

func parseFlexibleDate(val string) (time.Time, error) {
	val = strings.TrimSpace(val)
	layouts := []string{"2006-01-02", "02/01/2006", "01/02/2006", "2006/01/02", "02-01-2006"}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, val); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("formato de fecha no reconocido: %s", val)
}
