package request

type CreateZoneDTO struct {
	WarehouseID uint
	BusinessID  uint
	Code        string
	Name        string
	Purpose     string
	ColorHex    string
	IsActive    bool
}

type UpdateZoneDTO struct {
	ID         uint
	BusinessID uint
	Code       string
	Name       string
	Purpose    string
	ColorHex   string
	IsActive   *bool
}

type CreateAisleDTO struct {
	ZoneID     uint
	BusinessID uint
	Code       string
	Name       string
	IsActive   bool
	WidthCm    float64
}

type UpdateAisleDTO struct {
	ID         uint
	BusinessID uint
	Code       string
	Name       string
	IsActive   *bool
	WidthCm    *float64
}

type CreateRackDTO struct {
	AisleID     uint
	BusinessID  uint
	Code        string
	Name        string
	LevelsCount int
	IsActive    bool
	WidthCm     float64
	DepthCm     float64
	HeightCm    float64
}

type UpdateRackDTO struct {
	ID          uint
	BusinessID  uint
	Code        string
	Name        string
	LevelsCount *int
	IsActive    *bool
	WidthCm     *float64
	DepthCm     *float64
	HeightCm    *float64
}

type CreateRackLevelDTO struct {
	RackID     uint
	BusinessID uint
	Code       string
	Ordinal    int
	IsActive   bool
}

type UpdateRackLevelDTO struct {
	ID         uint
	BusinessID uint
	Code       string
	Ordinal    *int
	IsActive   *bool
}
