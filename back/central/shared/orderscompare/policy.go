package orderscompare

import "strings"

const (
	InventoryMoves = "moves"
	InventorySkips = "skips"
)

const (
	ReasonNone      = ""
	ReasonDelivered = "la orden ya viene entregada o completada del canal: el stock ya salio de la bodega"
	ReasonShipped   = "la orden ya viene despachada del canal: el stock ya salio de la bodega"
	ReasonCancelled = "la orden viene cancelada del canal: nunca hubo salida de stock que reservar"
	ReasonRefunded  = "la orden viene devuelta o reembolsada del canal: el movimiento ya lo hizo el canal"
	ReasonAbandoned = "la orden esta abandonada en el canal: no llego a ser una venta"
)

func SkipsInventoryFor(status, fulfillment string) (bool, string) {
	if skips, reason := SkipsInventory(status); skips {
		return true, reason
	}
	return SkipsInventory(fulfillment)
}

func SkipsInventory(status string) (bool, string) {
	s := strings.ToLower(strings.TrimSpace(status))
	switch {
	case s == "":
		return false, ReasonNone
	case strings.Contains(s, "delivered") || strings.Contains(s, "completed") || strings.Contains(s, "entregad"):
		return true, ReasonDelivered
	case strings.Contains(s, "shipped") || strings.Contains(s, "fulfilled") || strings.Contains(s, "despachad"):
		return true, ReasonShipped
	case strings.Contains(s, "cancel"):
		return true, ReasonCancelled
	case strings.Contains(s, "refund") || strings.Contains(s, "returned") || strings.Contains(s, "devuelt"):
		return true, ReasonRefunded
	case strings.Contains(s, "abandon"):
		return true, ReasonAbandoned
	}
	return false, ReasonNone
}
