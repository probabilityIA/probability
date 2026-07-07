package domain

import "context"

type GeocodeResult struct {
	Lat              float64
	Lng              float64
	FormattedAddress string
	Locality         string
	AdminArea2       string
	Department       string
	LocationType     string
	PartialMatch     bool
	Found            bool
}

type IGeocoder interface {
	Geocode(ctx context.Context, address string) (GeocodeResult, error)
}
