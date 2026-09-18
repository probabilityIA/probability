package errors

import "errors"

var (
	ErrRouteNotFound      = errors.New("route not found")
	ErrStopNotFound       = errors.New("route stop not found")
	ErrInvalidTransition  = errors.New("invalid status transition")
	ErrRouteNotPlanned    = errors.New("route must be in planned status")
	ErrRouteNotInProgress = errors.New("route must be in in_progress status")
	ErrDriverNotFound     = errors.New("driver not found")
	ErrVehicleNotFound    = errors.New("vehicle not found")
	ErrOrderNotFound      = errors.New("order not found")
	ErrStopIDsMismatch    = errors.New("stop IDs do not match route stops")

	ErrOptimizerNotConfigured = errors.New("la optimizacion de rutas no esta configurada")
	ErrNotEnoughStops         = errors.New("se necesitan al menos dos paradas con coordenadas para optimizar")
	ErrNoRouteFound           = errors.New("no hay un camino por carretera entre las paradas: revisa que todas est\u00e9n bien ubicadas")
	ErrMapsProvider           = errors.New("Google Maps no pudo calcular la ruta, intenta de nuevo en unos minutos")
	ErrRouteNotEditable       = errors.New("solo se puede optimizar una ruta en planeacion")
	ErrOriginMissing          = errors.New("la ruta no tiene bodega de origen: crea una bodega para el negocio")
	ErrOriginWithoutLocation  = errors.New("la bodega de origen no tiene ubicaci\u00f3n: edita la bodega y marca su ubicaci\u00f3n en el mapa")
)

type UnreachableStopsError struct {
	Stops []string
}

func (e *UnreachableStopsError) Error() string {
	msg := "no hay camino por carretera hasta: "
	for i, s := range e.Stops {
		if i > 0 {
			msg += "; "
		}
		msg += s
	}
	return msg + ". Corrige su ubicaci\u00f3n o qu\u00edtala de la ruta"
}

func (e *UnreachableStopsError) Is(target error) bool {
	return target == ErrNoRouteFound
}
