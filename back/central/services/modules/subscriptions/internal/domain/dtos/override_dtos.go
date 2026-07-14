package dtos

type GrantOverrideDTO struct {
	BusinessID      uint
	ModuleCode      string
	Notes           *string
	GrantedByUserID uint
}
