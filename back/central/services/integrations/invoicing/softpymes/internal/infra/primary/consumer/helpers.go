package consumer

import (
	"context"
	"encoding/json"
	"time"

	spDtos "github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/queue"
)

// mapItemsToClientDTOs convierte items del mensaje a DTOs del cliente Softpymes
func mapItemsToClientDTOs(items []invoiceItemData) []spDtos.ItemData {
	result := make([]spDtos.ItemData, 0, len(items))
	for _, item := range items {
		result = append(result, spDtos.ItemData{
			ProductID:             item.ProductID,
			SKU:                   item.SKU,
			Name:                  item.Name,
			Description:           item.Description,
			Quantity:              item.Quantity,
			UnitPrice:             item.UnitPrice,
			UnitPriceBase:         item.UnitPriceBase,
			TotalPrice:            item.TotalPrice,
			Tax:                   item.Tax,
			TaxRate:               item.TaxRate,
			Discount:              item.Discount,
			DiscountPercent:       item.DiscountPercent,
			UnitPricePresentment:      item.UnitPricePresentment,
			UnitPriceBasePresentment:  item.UnitPriceBasePresentment,
			TotalPricePresentment:     item.TotalPricePresentment,
			DiscountPresentment:       item.DiscountPresentment,
			TaxPresentment:            item.TaxPresentment,
		})
	}
	return result
}

// createErrorResponse crea una respuesta de error, opcionalmente con audit data
func (c *InvoiceRequestConsumer) createErrorResponse(
	request *InvoiceRequestMessage,
	errorCode string,
	errorMsg string,
	startTime time.Time,
	auditData *spDtos.AuditData,
) *queue.InvoiceResponseMessage {
	processingTime := time.Since(startTime).Milliseconds()

	resp := &queue.InvoiceResponseMessage{
		InvoiceID:      request.InvoiceID,
		Provider:       "softpymes",
		Status:         "error",
		Error:          errorMsg,
		ErrorCode:      errorCode,
		CorrelationID:  request.CorrelationID,
		Timestamp:      time.Now(),
		ProcessingTime: processingTime,
	}

	// Incluir audit data si está disponible (ej: cuando el HTTP request se hizo pero falló)
	if auditData != nil {
		resp.AuditRequestURL = auditData.RequestURL
		resp.AuditRequestPayload = toMapPayload(auditData.RequestPayload)
		resp.AuditResponseStatus = auditData.ResponseStatus
		resp.AuditResponseBody = auditData.ResponseBody
	}

	return resp
}

// sendCashReceiptIfConfigured envía un recibo de caja si la config lo tiene habilitado.
// Es non-fatal: si falla, se loguea el error pero no afecta el resultado de la factura.
func (c *InvoiceRequestConsumer) sendCashReceiptIfConfigured(
	ctx context.Context,
	fullDocument map[string]interface{},
	config map[string]interface{},
	apiKey, apiSecret, referer, baseURL string,
	invoiceID uint,
) {
	sendCashReceipt, _ := config["send_cash_receipt"].(bool)
	if !sendCashReceipt {
		return
	}

	if fullDocument == nil {
		c.log.Warn(ctx).
			Uint("invoice_id", invoiceID).
			Msg("Cash receipt configured but full document is nil — skipping")
		return
	}

	c.log.Info(ctx).
		Uint("invoice_id", invoiceID).
		Msg("Sending cash receipt (configured in integration)")

	if err := c.softpymesClient.SendCashReceiptFromDocument(ctx, apiKey, apiSecret, referer, baseURL, fullDocument, config); err != nil {
		c.log.Error(ctx).Err(err).
			Uint("invoice_id", invoiceID).
			Msg("Cash receipt failed — invoice created but payment not registered in Softpymes")
	} else {
		c.log.Info(ctx).
			Uint("invoice_id", invoiceID).
			Msg("Cash receipt sent successfully")
	}
}

// toMapPayload convierte cualquier valor (struct o map) a map[string]interface{} via JSON.
func toMapPayload(v interface{}) map[string]interface{} {
	if v == nil {
		return nil
	}
	if m, ok := v.(map[string]interface{}); ok {
		return m
	}
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil
	}
	return result
}
