package mappers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/primary/handlers/response"
)

type shippingLineMeta struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type shippingLineDetail struct {
	Title       string             `json:"title"`
	Price       any                `json:"price"`
	Source      string             `json:"source"`
	Code        string             `json:"code"`
	MethodID    string             `json:"method_id"`
	MethodTitle string             `json:"method_title"`
	Total       any                `json:"total"`
	MetaData    []shippingLineMeta `json:"meta_data"`
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
		if line.Source == "probability" {
			price, _ := strconv.ParseFloat(anyString(line.Price), 64)
			quoteID, rateIndex := splitQuoteCode(line.Code)
			return &response.QuotedShipping{
				Carrier:   carrierFromTitle(line.Title),
				Title:     line.Title,
				Price:     price,
				QuoteID:   quoteID,
				RateIndex: rateIndex,
			}
		}

		if line.MethodID == "probability_shipping" {
			meta := map[string]string{}
			for _, m := range line.MetaData {
				meta[m.Key] = anyString(m.Value)
			}
			carrier := meta["carrier"]
			if carrier == "" {
				carrier = carrierFromTitle(line.MethodTitle)
			}
			price, _ := strconv.ParseFloat(anyString(line.Total), 64)
			fee, _ := strconv.ParseFloat(meta["cod_carrier_fee"], 64)
			return &response.QuotedShipping{
				Carrier:       carrier,
				Title:         line.MethodTitle,
				Price:         price,
				QuoteID:       meta["quote_id"],
				RateIndex:     meta["rate_index"],
				CodCarrierFee: fee,
			}
		}
	}

	return nil
}

func anyString(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return strings.TrimSpace(t)
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	default:
		return strings.TrimSpace(fmt.Sprint(t))
	}
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
