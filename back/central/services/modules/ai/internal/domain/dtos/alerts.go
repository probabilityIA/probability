package dtos

import "time"

type AlertEvent struct {
	ID         string
	Type       string
	BusinessID uint
	Timestamp  time.Time
	Data       map[string]any
}

type AlertQuery struct {
	BusinessID uint
	UserID     uint
	Page       int
	PageSize   int
}
