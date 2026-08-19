package entities

import "time"

type BusinessModuleOverride struct {
	ID              uint
	BusinessID      uint
	ModuleCode      string
	GrantedByUserID uint
	Notes           *string
	ExpiresAt       *time.Time
	CreatedAt       time.Time
}
