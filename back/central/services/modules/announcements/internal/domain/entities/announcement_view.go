package entities

import "time"

type ViewAction string

const (
	ViewActionViewed      ViewAction = "viewed"
	ViewActionClosed      ViewAction = "closed"
	ViewActionClickedLink ViewAction = "clicked_link"
	ViewActionAccepted    ViewAction = "accepted"
)

type AnnouncementView struct {
	ID             uint
	AnnouncementID uint
	UserID         uint
	BusinessID     uint
	Action         ViewAction
	LinkID         *uint
	ViewedAt       time.Time
	CreatedAt      time.Time
}
