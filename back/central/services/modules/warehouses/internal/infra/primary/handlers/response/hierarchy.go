package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
)

type ZoneResponse struct {
	ID          uint      `json:"id"`
	WarehouseID uint      `json:"warehouse_id"`
	BusinessID  uint      `json:"business_id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Purpose     string    `json:"purpose"`
	ColorHex    string    `json:"color_hex"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AisleResponse struct {
	ID         uint      `json:"id"`
	ZoneID     uint      `json:"zone_id"`
	BusinessID uint      `json:"business_id"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type RackResponse struct {
	ID          uint      `json:"id"`
	AisleID     uint      `json:"aisle_id"`
	BusinessID  uint      `json:"business_id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	LevelsCount int       `json:"levels_count"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RackLevelResponse struct {
	ID         uint      `json:"id"`
	RackID     uint      `json:"rack_id"`
	BusinessID uint      `json:"business_id"`
	Code       string    `json:"code"`
	Ordinal    int       `json:"ordinal"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ZoneListResponse struct {
	Data       []ZoneResponse `json:"data"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

type AisleListResponse struct {
	Data       []AisleResponse `json:"data"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

type RackListResponse struct {
	Data       []RackResponse `json:"data"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

type RackLevelListResponse struct {
	Data       []RackLevelResponse `json:"data"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	TotalPages int                 `json:"total_pages"`
}

type TreePositionResponse struct {
	ID       uint   `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	IsActive bool   `json:"is_active"`
	Priority int    `json:"priority"`
}

type TreeRackLevelResponse struct {
	RackLevelResponse
	Positions []TreePositionResponse `json:"positions"`
}

type TreeRackResponse struct {
	RackResponse
	Levels []TreeRackLevelResponse `json:"levels"`
}

type TreeAisleResponse struct {
	AisleResponse
	Racks []TreeRackResponse `json:"racks"`
}

type TreeZoneResponse struct {
	ZoneResponse
	Aisles []TreeAisleResponse `json:"aisles"`
}

type WarehouseTreeResponse struct {
	WarehouseID uint               `json:"warehouse_id"`
	Zones       []TreeZoneResponse `json:"zones"`
}

func ZoneFromEntity(e *entities.WarehouseZone) ZoneResponse {
	return ZoneResponse{
		ID:          e.ID,
		WarehouseID: e.WarehouseID,
		BusinessID:  e.BusinessID,
		Code:        e.Code,
		Name:        e.Name,
		Purpose:     e.Purpose,
		ColorHex:    e.ColorHex,
		IsActive:    e.IsActive,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func AisleFromEntity(e *entities.WarehouseAisle) AisleResponse {
	return AisleResponse{
		ID:         e.ID,
		ZoneID:     e.ZoneID,
		BusinessID: e.BusinessID,
		Code:       e.Code,
		Name:       e.Name,
		IsActive:   e.IsActive,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}

func RackFromEntity(e *entities.WarehouseRack) RackResponse {
	return RackResponse{
		ID:          e.ID,
		AisleID:     e.AisleID,
		BusinessID:  e.BusinessID,
		Code:        e.Code,
		Name:        e.Name,
		LevelsCount: e.LevelsCount,
		IsActive:    e.IsActive,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func RackLevelFromEntity(e *entities.WarehouseRackLevel) RackLevelResponse {
	return RackLevelResponse{
		ID:         e.ID,
		RackID:     e.RackID,
		BusinessID: e.BusinessID,
		Code:       e.Code,
		Ordinal:    e.Ordinal,
		IsActive:   e.IsActive,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}

func WarehouseTreeFromDTO(t *dtos.WarehouseTreeDTO) WarehouseTreeResponse {
	resp := WarehouseTreeResponse{WarehouseID: t.WarehouseID}
	resp.Zones = make([]TreeZoneResponse, 0, len(t.Zones))
	for _, z := range t.Zones {
		zResp := TreeZoneResponse{ZoneResponse: ZoneFromEntity(&z.WarehouseZone)}
		zResp.Aisles = make([]TreeAisleResponse, 0, len(z.Aisles))
		for _, a := range z.Aisles {
			aResp := TreeAisleResponse{AisleResponse: AisleFromEntity(&a.WarehouseAisle)}
			aResp.Racks = make([]TreeRackResponse, 0, len(a.Racks))
			for _, r := range a.Racks {
				rResp := TreeRackResponse{RackResponse: RackFromEntity(&r.WarehouseRack)}
				rResp.Levels = make([]TreeRackLevelResponse, 0, len(r.Levels))
				for _, lv := range r.Levels {
					lvResp := TreeRackLevelResponse{RackLevelResponse: RackLevelFromEntity(&lv.WarehouseRackLevel)}
					lvResp.Positions = make([]TreePositionResponse, 0, len(lv.Positions))
					for _, p := range lv.Positions {
						lvResp.Positions = append(lvResp.Positions, TreePositionResponse{
							ID:       p.ID,
							Code:     p.Code,
							Name:     p.Name,
							Type:     p.Type,
							IsActive: p.IsActive,
							Priority: p.Priority,
						})
					}
					rResp.Levels = append(rResp.Levels, lvResp)
				}
				aResp.Racks = append(aResp.Racks, rResp)
			}
			zResp.Aisles = append(zResp.Aisles, aResp)
		}
		resp.Zones = append(resp.Zones, zResp)
	}
	return resp
}
