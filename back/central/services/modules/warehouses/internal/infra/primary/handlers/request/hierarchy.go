package request

type CreateZoneRequest struct {
	WarehouseID uint   `json:"warehouse_id" binding:"required,min=1"`
	Code        string `json:"code" binding:"required,min=1,max=50"`
	Name        string `json:"name" binding:"required,min=1,max=255"`
	Purpose     string `json:"purpose"`
	ColorHex    string `json:"color_hex"`
	IsActive    *bool  `json:"is_active"`
}

type UpdateZoneRequest struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Purpose  string `json:"purpose"`
	ColorHex string `json:"color_hex"`
	IsActive *bool  `json:"is_active"`
}

type CreateAisleRequest struct {
	ZoneID   uint   `json:"zone_id" binding:"required,min=1"`
	Code     string `json:"code" binding:"required,min=1,max=50"`
	Name     string `json:"name" binding:"required,min=1,max=255"`
	IsActive *bool  `json:"is_active"`
}

type UpdateAisleRequest struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	IsActive *bool  `json:"is_active"`
}

type CreateRackRequest struct {
	AisleID     uint   `json:"aisle_id" binding:"required,min=1"`
	Code        string `json:"code" binding:"required,min=1,max=50"`
	Name        string `json:"name" binding:"required,min=1,max=255"`
	LevelsCount int    `json:"levels_count"`
	IsActive    *bool  `json:"is_active"`
}

type UpdateRackRequest struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	LevelsCount *int   `json:"levels_count"`
	IsActive    *bool  `json:"is_active"`
}

type CreateRackLevelRequest struct {
	RackID   uint   `json:"rack_id" binding:"required,min=1"`
	Code     string `json:"code" binding:"required,min=1,max=50"`
	Ordinal  int    `json:"ordinal"`
	IsActive *bool  `json:"is_active"`
}

type UpdateRackLevelRequest struct {
	Code     string `json:"code"`
	Ordinal  *int   `json:"ordinal"`
	IsActive *bool  `json:"is_active"`
}
