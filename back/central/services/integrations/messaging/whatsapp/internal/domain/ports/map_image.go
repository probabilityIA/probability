package ports

import (
	"context"
	"errors"
)

var ErrMapImageNotConfigured = errors.New("la generacion del mapa de la direccion no esta configurada")

type IMapImageGenerator interface {
	BuildAddressMap(ctx context.Context, lat, lng float64, referencia string) (string, error)
	IsConfigured() bool
}
