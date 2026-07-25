package mappers

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/primary/handlers/response"
)

type shippingLineDetail struct {
	Title  string `json:"title"`
	Price  string `json:"price"`
	Source string `json:"source"`
	Code   string `json:"code"`
}

type shippingDetailsPayload struct {
	ShippingLines []shippingLineDetail `json:"shipping_lines"`
}

func buildQuotedShipping(raw []byte) *response.QuotedShipping {
	if len(raw) == 0 {
		return nil
	}

	var payload shippingDetailsPayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil
	}

	for _, line := range payload.ShippingLines {
		if line.Source != "probability" {
			continue
		}

		price, _ := strconv.ParseFloat(line.Price, 64)
		quoteID, rateIndex := splitQuoteCode(line.Code)

		return &response.QuotedShipping{
			Carrier:   carrierFromTitle(line.Title),
			Title:     line.Title,
			Price:     price,
			QuoteID:   quoteID,
			RateIndex: rateIndex,
		}
	}

	return nil
}

func splitQuoteCode(code string) (string, string) {
	parts := strings.Split(strings.TrimPrefix(code, "pq-"), "-")
	if len(parts) < 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

func carrierFromTitle(title string) string {
	if idx := strings.Index(title, " - "); idx > 0 {
		return strings.TrimSpace(title[:idx])
	}
	return strings.TrimSpace(title)
}
