package entities

import "time"

type BusinessModuleOverride struct {
	ID              uint
	BusinessID      uint
	ModuleCode      string
	GrantedByUserID uint
	Notes           *string
	CreatedAt       time.Time
}
