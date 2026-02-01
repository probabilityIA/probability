package response

import "gorm.io/datatypes"

// OrderRaw representa la respuesta HTTP con datos crudos de una orden
// ✅ DTO HTTP - CON TAGS (json + datatypes.JSON)
type OrderRaw struct {
	OrderID       string         `json:"order_id"`
	ChannelSource string         `json:"channel_source"`
	RawData       datatypes.JSON `json:"raw_data"`
}
