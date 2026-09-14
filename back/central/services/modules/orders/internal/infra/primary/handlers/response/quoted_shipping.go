package response

type QuotedShipping struct {
	Carrier       string  `json:"carrier"`
	Title         string  `json:"title"`
	Price         float64 `json:"price"`
	QuoteID       string  `json:"quote_id"`
	RateIndex     string  `json:"rate_index"`
	CodCarrierFee float64 `json:"cod_carrier_fee,omitempty"`
}
