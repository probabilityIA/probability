package response

type ConvertUoMResult struct {
	FromUomCode      string  `json:"from_uom_code"`
	ToUomCode        string  `json:"to_uom_code"`
	InputQuantity    float64 `json:"input_quantity"`
	ConvertedQty     float64 `json:"converted_quantity"`
	BaseUnitQuantity float64 `json:"base_unit_quantity"`
	BaseUomCode      string  `json:"base_uom_code"`
}
