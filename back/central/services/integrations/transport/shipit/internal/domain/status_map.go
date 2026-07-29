package domain

import "strings"

var shipitStatusTable = map[string]ProbabilityShipmentStatus{
	"created":           StatusPending,
	"in_preparation":    StatusPending,
	"waiting_courier":   StatusPending,
	"ready_to_ship":     StatusPending,
	"ready_for_pickup":  StatusPending,
	"picked_up":         StatusPickedUp,
	"withdrawn":         StatusPickedUp,
	"shipped":           StatusInTransit,
	"dispatched":        StatusInTransit,
	"in_transit":        StatusInTransit,
	"in_process":        StatusInTransit,
	"processing":        StatusInTransit,
	"out_for_delivery":  StatusOutForDelivery,
	"in_route":          StatusOutForDelivery,
	"in_distribution":   StatusOutForDelivery,
	"delivered":         StatusDelivered,
	"successful":        StatusDelivered,
	"failed":            StatusFailed,
	"failed_delivery":   StatusFailed,
	"rejected":          StatusFailed,
	"returned":          StatusReturned,
	"in_return":         StatusReturned,
	"returning":         StatusReturned,
	"cancelled":         StatusCancelled,
	"canceled":          StatusCancelled,
	"annulled":          StatusCancelled,
	"on_hold":           StatusOnHold,
	"retained":          StatusOnHold,
	"incidence":         StatusOnHold,
	"claim":             StatusOnHold,
	"pending_pickup":    StatusPending,
	"pending":           StatusPending,
	"fulfillment_ready": StatusPending,
}

func normalizeShipitStatus(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")
	return s
}

func MapShipitStatus(status string) (ProbabilityShipmentStatus, bool) {
	normalized := normalizeShipitStatus(status)
	if normalized == "" {
		return StatusPending, false
	}
	if mapped, ok := shipitStatusTable[normalized]; ok {
		return mapped, true
	}
	return mapByHeuristics(normalized)
}

func mapByHeuristics(normalized string) (ProbabilityShipmentStatus, bool) {
	spaced := strings.ReplaceAll(normalized, "_", " ")
	switch {
	case strings.Contains(spaced, "entregad") || strings.Contains(spaced, "deliver"):
		return StatusDelivered, true
	case strings.Contains(spaced, "devol") || strings.Contains(spaced, "return"):
		return StatusReturned, true
	case strings.Contains(spaced, "cancel") || strings.Contains(spaced, "anulad"):
		return StatusCancelled, true
	case strings.Contains(spaced, "reparto") || strings.Contains(spaced, "ruta de entrega"):
		return StatusOutForDelivery, true
	case strings.Contains(spaced, "retir") || strings.Contains(spaced, "recolec") || strings.Contains(spaced, "pickup"):
		return StatusPickedUp, true
	case strings.Contains(spaced, "transito") || strings.Contains(spaced, "transit") || strings.Contains(spaced, "viaj") || strings.Contains(spaced, "despach"):
		return StatusInTransit, true
	case strings.Contains(spaced, "fallid") || strings.Contains(spaced, "fail") || strings.Contains(spaced, "rechaz"):
		return StatusFailed, true
	case strings.Contains(spaced, "reten") || strings.Contains(spaced, "incidenc") || strings.Contains(spaced, "hold"):
		return StatusOnHold, true
	case strings.Contains(spaced, "prepar") || strings.Contains(spaced, "cread") || strings.Contains(spaced, "generad"):
		return StatusPending, true
	}
	return StatusInTransit, false
}

func MapTrackingData(data *TrackingData) (ProbabilityShipmentStatus, string, bool) {
	if data == nil {
		return StatusPending, "", false
	}
	if data.IsDelivered {
		last := ""
		if len(data.Statuses) > 0 {
			last = data.Statuses[len(data.Statuses)-1].CourierStatus
		}
		return StatusDelivered, last, true
	}
	if len(data.Statuses) == 0 {
		return StatusPending, "", true
	}
	last := data.Statuses[len(data.Statuses)-1]
	status, known := MapShipitStatus(last.CourierStatus)
	return status, last.CourierStatus, known
}
