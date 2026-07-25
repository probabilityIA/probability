package response

type QuotedShipping struct {
	Carrier   string  `json:"carrier"`
	Title     string  `json:"title"`
	Price     float64 `json:"price"`
	QuoteID   string  `json:"quote_id"`
	RateIndex string  `json:"rate_index"`
}
