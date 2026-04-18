package dtos

import "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"

type ListZonesParams struct {
	BusinessID  uint
	WarehouseID uint
	ActiveOnly  bool
	Page        int
	PageSize    int
}

func (p ListZonesParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

type ListAislesParams struct {
	BusinessID uint
	ZoneID     uint
	Page       int
	PageSize   int
}

func (p ListAislesParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

type ListRacksParams struct {
	BusinessID uint
	AisleID    uint
	Page       int
	PageSize   int
}

func (p ListRacksParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

type ListRackLevelsParams struct {
	BusinessID uint
	RackID     uint
	Page       int
	PageSize   int
}

func (p ListRackLevelsParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

type PositionNode struct {
	entities.WarehouseLocation
}

type RackLevelNode struct {
	entities.WarehouseRackLevel
	Positions []PositionNode `json:"positions"`
}

type RackNode struct {
	entities.WarehouseRack
	Levels []RackLevelNode `json:"levels"`
}

type AisleNode struct {
	entities.WarehouseAisle
	Racks []RackNode `json:"racks"`
}

type ZoneNode struct {
	entities.WarehouseZone
	Aisles []AisleNode `json:"aisles"`
}

type WarehouseTreeDTO struct {
	WarehouseID uint       `json:"warehouse_id"`
	Zones       []ZoneNode `json:"zones"`
}
