package entities

type GeozoneAncestors struct {
	DeepestID        *uint
	Path             []uint
	CountryID        *uint
	StateID          *uint
	CityID           *uint
	AdminDistrictID  *uint
	LocalityID       *uint
	NeighborhoodID   *uint
	BarrioID         *uint
}
