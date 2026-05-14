package errors

import "errors"

var (
	ErrGeozoneNotFound  = errors.New("geozona no encontrada")
	ErrInvalidGeometry  = errors.New("geometria invalida")
	ErrInvalidType      = errors.New("tipo de geozona invalido")
	ErrDuplicateGeozone = errors.New("ya existe una geozona con el mismo negocio/tipo/codigo")
)
