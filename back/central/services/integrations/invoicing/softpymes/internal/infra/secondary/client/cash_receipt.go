package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// SendCashReceiptFromDocument envía un recibo de caja a Softpymes usando datos
// del documento completo retornado por GetDocumentByNumber y la configuración
// de la integración.
//
// Extrae del fullDocument: prefix, documentNumber, total, customerIdentification,
// branchCode, documentDate.
// Lee de config: payment_type, payment_bank_account_id, payment_financial_entity_id,
// payment_bonus_code, payment_bank_name, payment_account_number.
func (c *Client) SendCashReceiptFromDocument(
	ctx context.Context,
	apiKey, apiSecret, referer, baseURL string,
	fullDocument map[string]interface{},
	config map[string]interface{},
) (map[string]interface{}, error) {
	// 1. Authenticate
	token, err := c.authenticate(ctx, apiKey, apiSecret, referer, baseURL)
	if err != nil {
		return nil, fmt.Errorf("cash receipt auth failed: %w", err)
	}

	// 2. Extract fields from fullDocument
	prefix, _ := fullDocument["prefix"].(string)
	documentNumber, _ := fullDocument["documentNumber"].(string)
	customerNit, _ := fullDocument["customerIdentification"].(string)
	branchCode, _ := fullDocument["branchCode"].(string)
	customerBranchCode, _ := fullDocument["customerBranchCode"].(string)
	documentDate, _ := fullDocument["documentDate"].(string)

	// Parse total from document (may be string or float64)
	var amount float64
	switch v := fullDocument["total"].(type) {
	case float64:
		amount = v
	case string:
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			amount = parsed
		}
	}

	// Fallbacks
	if branchCode == "" {
		branchCode = "001"
		if bc, ok := config["branch_code"].(string); ok && bc != "" {
			branchCode = bc
		}
	}
	// Primero intentar leer de config si no vino del documento
	if customerBranchCode == "" || customerBranchCode == "000" {
		if cb, ok := config["customer_branch_code"].(string); ok && cb != "" && cb != "000" {
			customerBranchCode = cb
		} else {
			// Softpymes cash_receipt no acepta "000" (da 500), default a "001"
			customerBranchCode = "001"
		}
	}
	if documentDate == "" {
		loc, _ := time.LoadLocation("America/Bogota")
		documentDate = time.Now().In(loc).Format("2006-01-02")
	}

	if documentNumber == "" {
		return nil, fmt.Errorf("cash receipt: documentNumber is empty in full document")
	}
	if amount <= 0 {
		return nil, fmt.Errorf("cash receipt: total is 0 or negative in full document")
	}

	// 3. Read payment config
	paymentType, _ := config["payment_type"].(string)
	if paymentType == "" {
		paymentType = "EF"
	}

	// 4. Build payment body per type
	payment := buildPaymentBody(paymentType, amount, config, documentNumber, prefix)

	// 5. Build and send request
	// Split prefix from documentNumber if it's combined (e.g. "FEV0000000026")
	docPrefix := prefix
	docNumber := documentNumber
	if docPrefix == "" {
		docPrefix, docNumber = splitDocumentNumber(documentNumber)
	}

	body := map[string]interface{}{
		"documents": []map[string]interface{}{
			{
				"documentNumber": docNumber,
				"prefix":         docPrefix,
			},
		},
		"documentDate":       documentDate,
		"branchCode":         branchCode,
		"customerNit":        customerNit,
		"customerBranchCode": customerBranchCode,
		"payment":            []map[string]interface{}{payment},
	}

	c.log.Info(ctx).
		Str("document_number", docNumber).
		Str("prefix", docPrefix).
		Str("payment_type", paymentType).
		Float64("amount", amount).
		Interface("payment_body", payment).
		Interface("request_body", body).
		Msg("Sending cash receipt to Softpymes")

	requestURL := c.resolveURL(baseURL, "/app/integration/cash_receipt/")
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Referer", referer).
		SetBody(body).
		Post(requestURL)

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Cash receipt request failed")
		return nil, fmt.Errorf("cash receipt request failed: %w", err)
	}

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("response", string(resp.Body())).
			Msg("Cash receipt creation failed")
		// Retornar datos de audit incluso en error para trazabilidad
		failedData := map[string]interface{}{
			"status":                "failed",
			"error":                 fmt.Sprintf("cash receipt failed with status %d: %s", resp.StatusCode(), string(resp.Body())),
			"audit_request_url":     requestURL,
			"audit_request_payload": body,
			"audit_response_status": resp.StatusCode(),
			"audit_response_body":   string(resp.Body()),
		}
		return failedData, fmt.Errorf("cash receipt failed with status %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	// Softpymes puede devolver:
	// 1. Objeto: {"message": "...", "error": ""}
	// 2. Array con error embebido: [[{"message":"..."}], 400]
	// 3. Array con éxito: [{"message":"..."}, 200]
	respBody := resp.Body()
	message, errMsg := parseCashReceiptResponse(respBody)

	if errMsg != "" {
		c.log.Error(ctx).
			Str("error", errMsg).
			Str("raw_response", string(respBody)).
			Msg("Cash receipt API returned error")
		// Retornar datos de audit incluso en error embebido en HTTP 200
		failedData := map[string]interface{}{
			"status":                "failed",
			"error":                 fmt.Sprintf("cash receipt error: %s", errMsg),
			"audit_request_url":     requestURL,
			"audit_request_payload": body,
			"audit_response_status": resp.StatusCode(),
			"audit_response_body":   string(respBody),
		}
		return failedData, fmt.Errorf("cash receipt error: %s", errMsg)
	}

	c.log.Info(ctx).
		Str("document_number", docNumber).
		Str("message", message).
		Msg("Cash receipt sent successfully - payment registered in Softpymes")

	receiptData := map[string]interface{}{
		"status":          "success",
		"message":         message,
		"payment_type":    paymentType,
		"amount":          amount,
		"document_number": docNumber,
		"prefix":          docPrefix,
		"customer_nit":    customerNit,
		"document_date":   documentDate,
		"request_body":    body,
		"raw_response":    string(respBody),
		// Audit data para trazabilidad separada
		"audit_request_url":     requestURL,
		"audit_request_payload": body,
		"audit_response_status": resp.StatusCode(),
		"audit_response_body":   string(respBody),
	}

	return receiptData, nil
}

// buildPaymentBody construye el objeto de pago según el tipo.
//
// Tipos soportados:
//   - EF (Efectivo): solo type + value
//   - BN (Bonos): type + value + code
//   - CH (Cheque): type + value + accountNumber + bankName
//   - TR (Transferencia): type + value + accountNumber + documentNumber + prefixNumber
//   - TC/TD (Tarjeta crédito/débito): type + value + finantialEntityId (int)
func buildPaymentBody(paymentType string, amount float64, config map[string]interface{}, docNumber, docPrefix string) map[string]interface{} {
	payment := map[string]interface{}{
		"type":  paymentType,
		"value": amount,
	}

	switch paymentType {
	case "BN":
		if code, ok := config["payment_bonus_code"].(string); ok && code != "" {
			payment["code"] = code
		}
	case "CH":
		if acct, ok := config["payment_account_number"].(string); ok && acct != "" {
			payment["accountNumber"] = acct
		}
		if bank, ok := config["payment_bank_name"].(string); ok && bank != "" {
			payment["bankName"] = bank
		}
	case "TR":
		// accountNumber puede venir como string o como float64 (JSON number)
		switch v := config["payment_bank_account_id"].(type) {
		case string:
			if v != "" {
				payment["accountNumber"] = v
			}
		case float64:
			payment["accountNumber"] = fmt.Sprintf("%.0f", v)
		case int:
			payment["accountNumber"] = fmt.Sprintf("%d", v)
		}
		payment["documentNumber"] = docNumber
		payment["prefixNumber"] = docPrefix
	case "TC", "TD":
		if entityID := getConfigInt(config, "payment_financial_entity_id"); entityID > 0 {
			payment["finantialEntityId"] = entityID
		}
	}

	return payment
}

// getConfigInt extracts an int from config (which may store float64 from JSON).
func getConfigInt(config map[string]interface{}, key string) int {
	if v, ok := config[key].(float64); ok {
		return int(v)
	}
	if v, ok := config[key].(int); ok {
		return v
	}
	return 0
}

// parseCashReceiptResponse parsea la respuesta del endpoint cash_receipt de Softpymes.
// Softpymes puede devolver un objeto simple o un array con status code embebido.
// Retorna (message, errorMsg). Si errorMsg != "" hubo error.
func parseCashReceiptResponse(body []byte) (message string, errorMsg string) {
	if len(body) == 0 {
		return "", ""
	}

	// Intentar como objeto simple: {"message": "...", "error": "..."}
	var obj struct {
		Message string `json:"message"`
		Error   string `json:"error"`
	}
	if err := json.Unmarshal(body, &obj); err == nil {
		return obj.Message, obj.Error
	}

	// Intentar como array: [{...}, statusCode] o [[{...}], statusCode]
	var arr []json.RawMessage
	if err := json.Unmarshal(body, &arr); err != nil {
		return "", ""
	}

	if len(arr) < 2 {
		return "", ""
	}

	// Verificar status code en el último elemento
	var statusCode int
	if err := json.Unmarshal(arr[len(arr)-1], &statusCode); err == nil && statusCode >= 400 {
		// Error embebido — extraer mensaje del primer elemento
		var errMsg string

		// Intentar como array de errores: [[{"message":"..."}], 400]
		var errArr []json.RawMessage
		if json.Unmarshal(arr[0], &errArr) == nil && len(errArr) > 0 {
			var errObj struct {
				Message string `json:"message"`
			}
			if json.Unmarshal(errArr[0], &errObj) == nil {
				errMsg = errObj.Message
			}
		}

		// Intentar como objeto directo: [{"message":"..."}, 400]
		if errMsg == "" {
			var errObj struct {
				Message string `json:"message"`
			}
			if json.Unmarshal(arr[0], &errObj) == nil {
				errMsg = errObj.Message
			}
		}

		if errMsg == "" {
			errMsg = string(arr[0])
		}
		return "", errMsg
	}

	// Status < 400 o no es un status code → éxito
	// Extraer message del primer elemento
	if json.Unmarshal(arr[0], &obj) == nil {
		return obj.Message, obj.Error
	}

	return "", ""
}

// splitDocumentNumber separa el número combinado que retorna Softpymes en creación (ej: "FEV26")
// en prefijo ("FEV") y número con ceros ("0000000026").
// Si no tiene prefijo alfabético, retorna ("", paddedNumber).
func splitDocumentNumber(documentNumber string) (prefix, paddedNumber string) {
	bare := documentNumber
	for i, ch := range documentNumber {
		if ch >= '0' && ch <= '9' {
			prefix = documentNumber[:i]
			bare = documentNumber[i:]
			break
		}
	}
	// Zero-pad a 10 dígitos (formato estándar de Softpymes)
	if n, err := strconv.ParseInt(bare, 10, 64); err == nil {
		paddedNumber = fmt.Sprintf("%010d", n)
	} else {
		paddedNumber = bare
	}
	return prefix, paddedNumber
}
