package entities

type OccupancyItem struct {
	LocationID uint
	Quantity   int
	Reserved   int
	Capacity   *int
}
