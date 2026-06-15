package response

import "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"

type OccupancyLocationResponse struct {
	LocationID uint `json:"location_id"`
	Quantity   int  `json:"quantity"`
	Reserved   int  `json:"reserved"`
	Capacity   *int `json:"capacity"`
}

type OccupancyResponse struct {
	WarehouseID uint                         `json:"warehouse_id"`
	Locations   []OccupancyLocationResponse `json:"locations"`
}

func OccupancyFromEntities(warehouseID uint, items []entities.OccupancyItem) OccupancyResponse {
	locations := make([]OccupancyLocationResponse, len(items))
	for i, it := range items {
		locations[i] = OccupancyLocationResponse{
			LocationID: it.LocationID,
			Quantity:   it.Quantity,
			Reserved:   it.Reserved,
			Capacity:   it.Capacity,
		}
	}
	return OccupancyResponse{WarehouseID: warehouseID, Locations: locations}
}
