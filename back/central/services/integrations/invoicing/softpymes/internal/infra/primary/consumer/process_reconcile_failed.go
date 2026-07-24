package consumer

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/queue"
)

const (
	reconcileMaxDaysBack    = 35
	reconcilePageSize       = 100
	reconcileMaxPagesPerDay = 60
)

type reconcilePair struct {
	invoiceID uint
	orderID   string
}

func (c *InvoiceRequestConsumer) processReconcileFailed(ctx context.Context, request *InvoiceRequestMessage) error {
	pairs := parseReconcilePairs(request.InvoiceData.Config)
	if len(pairs) == 0 {
		c.log.Warn(ctx).Msg("reconcile_failed sin facturas para procesar")
		return nil
	}

	dateFrom, _ := request.InvoiceData.Config["date_from"].(string)

	c.log.Info(ctx).
		Int("invoices", len(pairs)).
		Str("date_from", dateFrom).
		Msg("Starting reconcile of failed invoices against Softpymes")

	integrationIDStr := fmt.Sprintf("%d", request.InvoiceData.IntegrationID)
	integration, err := c.integrationCore.GetIntegrationByID(ctx, integrationIDStr)
	if err != nil {
		return c.publishReconcileErrors(ctx, request, pairs, "integration_not_found: "+err.Error())
	}

	apiKey, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_key")
	if err != nil {
		return c.publishReconcileErrors(ctx, request, pairs, "decryption_failed api_key")
	}
	apiSecret, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_secret")
	if err != nil {
		return c.publishReconcileErrors(ctx, request, pairs, "decryption_failed api_secret")
	}

	combinedConfig := make(map[string]interface{})
	for k, v := range integration.Config {
		combinedConfig[k] = v
	}
	for k, v := range request.InvoiceData.Config {
		combinedConfig[k] = v
	}
	referer, _ := combinedConfig["referer"].(string)

	effectiveURL := integration.BaseURL
	if integration.IsTesting && integration.BaseURLTest != "" {
		effectiveURL = integration.BaseURLTest
	}
	if effectiveURL == "" {
		return c.publishReconcileErrors(ctx, request, pairs, "missing_base_url")
	}

	loc, _ := time.LoadLocation("America/Bogota")
	now := time.Now().In(loc)

	daysBack := reconcileMaxDaysBack
	if parsed, parseErr := time.ParseInLocation("2006-01-02", dateFrom, loc); parseErr == nil {
		diff := int(now.Sub(parsed).Hours() / 24)
		if diff >= 0 && diff < reconcileMaxDaysBack {
			daysBack = diff
		}
	}

	docsByOrder := make(map[string]ports.ListedDocument)
	pageSize := strconv.Itoa(reconcilePageSize)

	for d := 0; d <= daysBack; d++ {
		day := now.AddDate(0, 0, -d).Format("2006-01-02")
		for page := 1; page <= reconcileMaxPagesPerDay; page++ {
			pageStr := strconv.Itoa(page)
			docs, listErr := c.softpymesClient.ListDocuments(ctx, apiKey, apiSecret, referer, ports.ListDocumentsParams{
				DateFrom: day,
				DateTo:   day,
				PageSize: &pageSize,
				Page:     &pageStr,
			}, effectiveURL)
			if listErr != nil {
				c.log.Warn(ctx).Err(listErr).
					Str("day", day).
					Int("page", page).
					Msg("Reconcile aborted: Softpymes document sweep failed")
				return c.publishReconcileErrors(ctx, request, pairs,
					"barrido de documentos fallo: "+listErr.Error())
			}

			for i := range docs {
				if docs[i].Annuled {
					continue
				}
				if orderID := extractOrderIDFromComment(docs[i].Comment); orderID != "" {
					if _, exists := docsByOrder[orderID]; !exists {
						docsByOrder[orderID] = docs[i]
					}
				}
			}

			if len(docs) < reconcilePageSize {
				break
			}
		}
	}

	found := 0
	requeued := 0
	for _, pair := range pairs {
		doc, exists := docsByOrder[pair.orderID]
		if !exists {
			requeued++
			c.publishReconcileResponse(ctx, &queue.InvoiceResponseMessage{
				InvoiceID:     pair.invoiceID,
				Provider:      "softpymes",
				Operation:     "reconcile_failed",
				Status:        "retry_requeued",
				CorrelationID: request.CorrelationID,
				Timestamp:     time.Now(),
				Error:         "sin documento en Softpymes - re-encolada para nuevo intento",
			})
			continue
		}

		found++
		var fullDocument map[string]interface{}
		docNumber := doc.DocumentNumber
		if fullDoc, docErr := c.softpymesClient.GetDocumentByNumber(ctx, apiKey, apiSecret, referer, docNumber, effectiveURL); docErr == nil && fullDoc != nil {
			fullDocument = fullDoc
			if canonical, ok := fullDocument["documentNumber"].(string); ok && canonical != "" {
				docNumber = canonical
			}
		}

		var issuedAt *time.Time
		if doc.DocumentDate != "" {
			if parsed, parseErr := time.Parse("2006-01-02", doc.DocumentDate); parseErr == nil {
				issuedAt = &parsed
			}
		}

		c.publishReconcileResponse(ctx, &queue.InvoiceResponseMessage{
			InvoiceID:     pair.invoiceID,
			Provider:      "softpymes",
			Operation:     "reconcile_failed",
			Status:        "success",
			InvoiceNumber: docNumber,
			ExternalID:    docNumber,
			IssuedAt:      issuedAt,
			DocumentJSON:  fullDocument,
			CorrelationID: request.CorrelationID,
			Timestamp:     time.Now(),
		})
	}

	c.log.Info(ctx).
		Int("total", len(pairs)).
		Int("found_in_provider", found).
		Int("requeued_for_creation", requeued).
		Int("documents_swept", len(docsByOrder)).
		Msg("Reconcile of failed invoices completed")

	return nil
}

func (c *InvoiceRequestConsumer) publishReconcileErrors(ctx context.Context, request *InvoiceRequestMessage, pairs []reconcilePair, errMsg string) error {
	for _, pair := range pairs {
		c.publishReconcileResponse(ctx, &queue.InvoiceResponseMessage{
			InvoiceID:     pair.invoiceID,
			Provider:      "softpymes",
			Operation:     "reconcile_failed",
			Status:        "error",
			CorrelationID: request.CorrelationID,
			Timestamp:     time.Now(),
			Error:         "reconcile fallo (proveedor no disponible): " + errMsg,
		})
	}
	return nil
}

func (c *InvoiceRequestConsumer) publishReconcileResponse(ctx context.Context, resp *queue.InvoiceResponseMessage) {
	if err := c.responsePublisher.PublishResponse(ctx, resp); err != nil {
		c.log.Error(ctx).Err(err).
			Uint("invoice_id", resp.InvoiceID).
			Msg("Failed to publish reconcile response")
	}
}

func parseReconcilePairs(config map[string]interface{}) []reconcilePair {
	raw, ok := config["reconcile_invoices"].([]interface{})
	if !ok {
		return nil
	}
	pairs := make([]reconcilePair, 0, len(raw))
	for _, item := range raw {
		entry, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		orderID, _ := entry["order_id"].(string)
		var invoiceID uint
		switch v := entry["invoice_id"].(type) {
		case float64:
			invoiceID = uint(v)
		case int:
			invoiceID = uint(v)
		}
		if invoiceID > 0 && orderID != "" {
			pairs = append(pairs, reconcilePair{invoiceID: invoiceID, orderID: orderID})
		}
	}
	return pairs
}

func extractOrderIDFromComment(comment string) string {
	idx := strings.Index(comment, "order:")
	if idx < 0 {
		return ""
	}
	rest := comment[idx+len("order:"):]
	if len(rest) < 36 {
		return ""
	}
	return rest[:36]
}
