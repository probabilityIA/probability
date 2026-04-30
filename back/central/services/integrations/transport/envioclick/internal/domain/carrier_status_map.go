package domain

import "strings"

type carrierStatusKey struct {
	carrier      string
	statusDetail string
}

var carrierStatusTable = map[carrierStatusKey]ProbabilityShipmentStatus{
	{"envia", "fecha de captura en el sistema"}:               StatusPending,
	{"envia", "guia generada en espera de transportadora"}:    StatusPending,
	{"envia", "fecha de recoleccion del envio"}:               StatusPickedUp,
	{"envia", "salida a ruta"}:                                StatusInTransit,
	{"envia", "en ruta de entrega"}:                           StatusOutForDelivery,
	{"envia", "envio entregado"}:                              StatusDelivered,
	{"envia", "cancelado"}:                                    StatusCancelled,
	{"envia", "en espera de ruta poblacion aledana"}:          StatusOnHold,
	{"envia", "no conocen destinatario en direccion destino"}: StatusOnHold,
	{"envia", "direccion incorrecta/ insuficiente"}:           StatusOnHold,
	{"envia", "direccion incorrecta insuficiente"}:            StatusOnHold,

	{"interrapidisimo", "envio admitido"}:           StatusPickedUp,
	{"interrapidisimo", "ingresado a bodega"}:       StatusInTransit,
	{"interrapidisimo", "despachado para bodega"}:   StatusInTransit,
	{"interrapidisimo", "viajando en ruta nacional"}: StatusInTransit,
	{"interrapidisimo", "viajando en ruta regional"}: StatusInTransit,

	{"deprisa", "guia generada en espera de transportadora"}: StatusPending,
	{"deprisa", "cancelado"}:                                 StatusCancelled,

	{"servientrega", "en alistamiento del cliente"}: StatusPending,
	{"servientrega", "cancelado"}:                   StatusCancelled,
}

func normalizeCarrier(c string) string {
	return strings.ToLower(strings.TrimSpace(c))
}

func MapCarrierStatusDetail(carrier, statusDetail string) (ProbabilityShipmentStatus, bool) {
	c := normalizeCarrier(carrier)
	d := normalize(statusDetail)
	if c == "" || d == "" {
		return "", false
	}
	if s, ok := carrierStatusTable[carrierStatusKey{c, d}]; ok {
		return s, true
	}
	return "", false
}

func MapEnvioClickEvent(carrier, status, statusStep, statusDetail string, incidence bool) (probStatus ProbabilityShipmentStatus, mappedFromTable bool, unknown bool) {
	if incidence {
		combined := normalize(status + " " + statusStep + " " + statusDetail)
		if strings.Contains(combined, "entregad") {
			return StatusDelivered, false, false
		}
		return StatusOnHold, false, false
	}

	if s, ok := MapCarrierStatusDetail(carrier, statusDetail); ok {
		return s, true, false
	}

	step := statusStep
	if step == "" {
		step = status
	}
	step = ApiStatusToStep(step, statusDetail)
	s, u := MapStatusStepToProbability(step, false)
	return s, false, u
}
