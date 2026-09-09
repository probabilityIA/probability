package errors

import "errors"

var (
	ErrRouteNotFound       = errors.New("route not found")
	ErrStopNotFound        = errors.New("route stop not found")
	ErrInvalidTransition   = errors.New("invalid status transition")
	ErrRouteNotPlanned     = errors.New("route must be in planned status")
	ErrRouteNotInProgress  = errors.New("route must be in in_progress status")
	ErrDriverNotFound      = errors.New("driver not found")
	ErrVehicleNotFound     = errors.New("vehicle not found")
	ErrOrderNotFound       = errors.New("order not found")
	ErrStopIDsMismatch     = errors.New("stop IDs do not match route stops")

	ErrOptimizerNotConfigured = errors.New("la optimizacion de rutas no esta configurada")
	ErrNotEnoughStops         = errors.New("se necesitan al menos dos paradas con coordenadas para optimizar")
	ErrNoRouteFound           = errors.New("no se encontro una ruta entre las paradas")
	ErrRouteNotEditable       = errors.New("solo se puede optimizar una ruta en planeacion")
	ErrOriginMissing          = errors.New("la ruta no tiene coordenadas de origen")
)
