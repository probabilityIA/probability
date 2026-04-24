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
}

type UpdateAisleDTO struct {
	ID         uint
	BusinessID uint
	Code       string
	Name       string
	IsActive   *bool
}

type CreateRackDTO struct {
	AisleID     uint
	BusinessID  uint
	Code        string
	Name        string
	LevelsCount int
	IsActive    bool
}

type UpdateRackDTO struct {
	ID          uint
	BusinessID  uint
	Code        string
	Name        string
	LevelsCount *int
	IsActive    *bool
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
