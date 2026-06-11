package consumer

import (
	"context"
	"fmt"
	"strconv"
	"time"

	siigoDtos "github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/queue"
)

func (c *InvoiceRequestConsumer) processCompareRequest(
	ctx context.Context,
	request *InvoiceRequestMessage,
) error {
	dateFrom, _ := request.InvoiceData.Config["date_from"].(string)
	dateTo, _ := request.InvoiceData.Config["date_to"].(string)
	mode, _ := request.InvoiceData.Config["mode"].(string)
	businessID := businessIDFromConfig(request.InvoiceData.Config)

	c.log.Info(ctx).
		Str("date_from", dateFrom).
		Str("date_to", dateTo).
		Uint("business_id", businessID).
		Str("correlation_id", request.CorrelationID).
		Msg("Starting Siigo compare request")

	publishErr := func(errMsg string) error {
		return c.responsePublisher.PublishCompareResponse(ctx, &queue.CompareResponseMessage{
			Operation:     "compare",
			Mode:          mode,
			CorrelationID: request.CorrelationID,
			BusinessID:    businessID,
			DateFrom:      dateFrom,
			DateTo:        dateTo,
			Error:         errMsg,
			Timestamp:     time.Now(),
		})
	}

	if dateFrom == "" || dateTo == "" {
		c.log.Error(ctx).Msg("date_from or date_to missing in compare config")
		return publishErr("date_from and date_to are required in compare config")
	}

	ictx, _, err := c.resolveIntegration(ctx, request)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to resolve integration for compare")
		return publishErr("failed to resolve integration: " + err.Error())
	}

	allDocs := make([]queue.CompareDocument, 0)
	pageSize := 100

	for page := 1; ; page++ {
		c.log.Info(ctx).
			Int("page", page).
			Str("date_from", dateFrom).
			Str("date_to", dateTo).
			Msg("Fetching invoices page from Siigo")

		result, err := c.siigoClient.ListInvoices(ctx, ictx.Credentials, siigoDtos.ListInvoicesParams{
			Page:     page,
			PageSize: pageSize,
			DateFrom: dateFrom,
			DateTo:   dateTo,
		})
		if err != nil {
			c.log.Error(ctx).Err(err).Int("page", page).Msg("Failed to list Siigo invoices")
			return publishErr(fmt.Sprintf("failed to list invoices (page %d): %s", page, err.Error()))
		}

		for _, inv := range result.Items {
			allDocs = append(allDocs, queue.CompareDocument{
				DocumentNumber:     inv.Number,
				DocumentDate:       inv.Date,
				Total:              strconv.FormatFloat(inv.Total, 'f', 2, 64),
				CustomerNit:        inv.CustomerID,
				CustomerName:       inv.CustomerName,
				Prefix:             inv.Prefix,
				Annuled:            inv.Annulled,
				ElectronicDocument: inv.StampStatus != "",
			})
		}

		c.log.Info(ctx).
			Int("page", page).
			Int("page_count", len(result.Items)).
			Int("total_accumulated", len(allDocs)).
			Msg("Siigo invoices page fetched")

		if len(result.Items) < pageSize {
			break
		}
	}

	c.log.Info(ctx).
		Int("total_documents", len(allDocs)).
		Str("correlation_id", request.CorrelationID).
		Msg("All Siigo invoices fetched, publishing compare response")

	return c.responsePublisher.PublishCompareResponse(ctx, &queue.CompareResponseMessage{
		Operation:         "compare",
		Mode:              mode,
		CorrelationID:     request.CorrelationID,
		BusinessID:        businessID,
		DateFrom:          dateFrom,
		DateTo:            dateTo,
		ProviderDocuments: allDocs,
		Timestamp:         time.Now(),
	})
}
