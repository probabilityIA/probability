package consumer

import (
	"context"
	"strconv"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/queue"
)

func (c *InvoiceRequestConsumer) processListBankAccountsRequest(
	ctx context.Context,
	request *InvoiceRequestMessage,
) error {
	businessID := businessIDFromConfig(request.InvoiceData.Config)

	c.log.Info(ctx).
		Uint("business_id", businessID).
		Str("correlation_id", request.CorrelationID).
		Msg("Starting Siigo list_bank_accounts request")

	publishErr := func(errMsg string) error {
		return c.responsePublisher.PublishListBankAccountsResponse(ctx, &queue.ListBankAccountsResponseMessage{
			Operation:     "list_bank_accounts",
			CorrelationID: request.CorrelationID,
			BusinessID:    businessID,
			Error:         errMsg,
			Timestamp:     time.Now(),
		})
	}

	ictx, _, err := c.resolveIntegration(ctx, request)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to resolve integration for list_bank_accounts")
		return publishErr("failed to resolve integration: " + err.Error())
	}

	paymentTypes, err := c.siigoClient.ListPaymentTypes(ctx, ictx.Credentials, "RC")
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to list Siigo payment types")
		return publishErr("failed to list payment types: " + err.Error())
	}

	items := make([]queue.BankAccountItem, 0, len(paymentTypes))
	for _, pt := range paymentTypes {
		items = append(items, queue.BankAccountItem{
			AccountNumber: strconv.Itoa(pt.ID),
			Name:          pt.Name,
			NameType:      pt.Type,
		})
	}

	c.log.Info(ctx).
		Int("total_accounts", len(items)).
		Str("correlation_id", request.CorrelationID).
		Msg("Siigo payment types fetched, publishing list_bank_accounts response")

	return c.responsePublisher.PublishListBankAccountsResponse(ctx, &queue.ListBankAccountsResponseMessage{
		Operation:     "list_bank_accounts",
		CorrelationID: request.CorrelationID,
		BusinessID:    businessID,
		Items:         items,
		Timestamp:     time.Now(),
	})
}
