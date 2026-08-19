package dtos

import "time"

type GrantOverrideDTO struct {
	BusinessID      uint
	ModuleCode      string
	Notes           *string
	ExpiresAt       *time.Time
	GrantedByUserID uint
}
