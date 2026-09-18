package entities

type OriginWarehouse struct {
	ID      uint
	Name    string
	Address string
	City    string
	Lat     *float64
	Lng     *float64
}
