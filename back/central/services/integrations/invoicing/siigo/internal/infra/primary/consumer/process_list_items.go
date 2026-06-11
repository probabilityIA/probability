package consumer

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/queue"
)

func (c *InvoiceRequestConsumer) processListItemsRequest(
	ctx context.Context,
	request *InvoiceRequestMessage,
) error {
	businessID := businessIDFromConfig(request.InvoiceData.Config)

	c.log.Info(ctx).
		Uint("business_id", businessID).
		Str("correlation_id", request.CorrelationID).
		Msg("Starting Siigo list_items request")

	publishErr := func(errMsg string) error {
		return c.responsePublisher.PublishListItemsResponse(ctx, &queue.ListItemsResponseMessage{
			Operation:     "list_items",
			CorrelationID: request.CorrelationID,
			BusinessID:    businessID,
			Error:         errMsg,
			Timestamp:     time.Now(),
		})
	}

	ictx, _, err := c.resolveIntegration(ctx, request)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to resolve integration for list_items")
		return publishErr("failed to resolve integration: " + err.Error())
	}

	allItems := make([]queue.ListItemsItem, 0)
	pageSize := 100

	for page := 1; ; page++ {
		c.log.Info(ctx).
			Int("page", page).
			Msg("Fetching products page from Siigo")

		products, err := c.siigoClient.ListProducts(ctx, ictx.Credentials, page, pageSize)
		if err != nil {
			c.log.Error(ctx).Err(err).Int("page", page).Msg("Failed to list Siigo products")
			return publishErr(fmt.Sprintf("failed to list products (page %d): %s", page, err.Error()))
		}

		for _, p := range products {
			allItems = append(allItems, queue.ListItemsItem{
				ItemCode:    p.Code,
				ItemName:    p.Name,
				ItemPrice:   p.Price,
				Description: p.Description,
			})
		}

		c.log.Info(ctx).
			Int("page", page).
			Int("page_count", len(products)).
			Int("total_accumulated", len(allItems)).
			Msg("Siigo products page fetched")

		if len(products) < pageSize {
			break
		}
	}

	c.log.Info(ctx).
		Int("total_items", len(allItems)).
		Str("correlation_id", request.CorrelationID).
		Msg("All Siigo products fetched, publishing list_items response")

	return c.responsePublisher.PublishListItemsResponse(ctx, &queue.ListItemsResponseMessage{
		Operation:     "list_items",
		CorrelationID: request.CorrelationID,
		BusinessID:    businessID,
		Items:         allItems,
		Timestamp:     time.Now(),
	})
}
