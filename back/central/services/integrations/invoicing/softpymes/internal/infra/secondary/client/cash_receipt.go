package client

import (
	"context"
	"fmt"
	"strconv"
)

// cashReceiptRequest datos para crear un recibo de caja en Softpymes.
// El recibo de caja registra el pago de una o más facturas y mueve la cuenta
// contable de "cuentas por cobrar" (130505xx) a la cuenta de pago correspondiente.
// Endpoint: POST /app/integration/cash_receipt/
// Documentación: https://api-integracion.softpymes.com.co/doc/#api-Documentos-PostCashReceipt
type cashReceiptRequest struct {
	DocumentNumber     string  // Número del documento a pagar (ej: "0000000026" o "26")
	Prefix             string  // Prefijo del documento (ej: "FEV")
	BranchCode         string  // Código de sucursal donde se genera el recibo
	CustomerNit        string  // NIT del cliente (sin puntos, comas ni dígito verificación)
	CustomerBranchCode string  // Código de sucursal del cliente
	PaymentType        string  // Tipo de pago: "EF"=Efectivo, "TR"=Transferencia, "TC"=Tarjeta crédito, "TD"=Tarjeta débito, "CH"=Cheque, "BN"=Bonos
	Amount             float64 // Valor a pagar
	DocumentDate       string  // Fecha del recibo YYYY-MM-DD (zona horaria Colombia)
	AccountNumber      string  // Número de cuenta (OBLIGATORIO para TR/CH). Ver /app/integration/bank_account
	BankName           string  // Nombre del banco (OBLIGATORIO para CH)
}

// sendCashReceipt envía un recibo de caja a Softpymes para registrar el pago de una factura.
// Esto mueve la cuenta contable de "cuentas por cobrar" (130505xx) a la cuenta del medio de pago.
// El token de autenticación es el mismo obtenido durante la creación de la factura.
func (c *Client) sendCashReceipt(ctx context.Context, token, referer, baseURL string, req *cashReceiptRequest) error {
	payment := map[string]interface{}{
		"type":  req.PaymentType,
		"value": req.Amount,
	}
	if req.AccountNumber != "" {
		payment["accountNumber"] = req.AccountNumber
	}
	if req.BankName != "" {
		payment["bankName"] = req.BankName
	}

	body := map[string]interface{}{
		"documents": []map[string]interface{}{
			{
				"documentNumber": req.DocumentNumber,
				"prefix":         req.Prefix,
			},
		},
		"documentDate":       req.DocumentDate,
		"branchCode":         req.BranchCode,
		"customerNit":        req.CustomerNit,
		"customerBranchCode": req.CustomerBranchCode,
		"payment":            []map[string]interface{}{payment},
	}

	c.log.Info(ctx).
		Str("document_number", req.DocumentNumber).
		Str("prefix", req.Prefix).
		Str("payment_type", req.PaymentType).
		Float64("amount", req.Amount).
		Msg("Sending cash receipt to Softpymes")

	var receiptResp struct {
		Message string `json:"message"`
		Error   string `json:"error"`
	}

	requestURL := c.resolveURL(baseURL, "/app/integration/cash_receipt/")
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Referer", referer).
		SetBody(body).
		SetResult(&receiptResp).
		Post(requestURL)

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Cash receipt request failed")
		return fmt.Errorf("cash receipt request failed: %w", err)
	}

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("response", string(resp.Body())).
			Msg("Cash receipt creation failed")
		return fmt.Errorf("cash receipt failed with status %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	if receiptResp.Error != "" {
		c.log.Error(ctx).
			Str("error", receiptResp.Error).
			Msg("Cash receipt API returned error")
		return fmt.Errorf("cash receipt error: %s", receiptResp.Error)
	}

	c.log.Info(ctx).
		Str("document_number", req.DocumentNumber).
		Str("message", receiptResp.Message).
		Msg("Cash receipt sent successfully - payment registered in Softpymes")

	return nil
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
