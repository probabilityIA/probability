package entities

import "time"

type DisplayType string

const (
	DisplayTypeModalImage DisplayType = "modal_image"
	DisplayTypeModalText  DisplayType = "modal_text"
	DisplayTypeTicker     DisplayType = "ticker"
)

type FrequencyType string

const (
	FrequencyOnce              FrequencyType = "once"
	FrequencyDaily             FrequencyType = "daily"
	FrequencyAlways            FrequencyType = "always"
	FrequencyRequiresAcceptance FrequencyType = "requires_acceptance"
)

type AnnouncementStatus string

const (
	StatusDraft     AnnouncementStatus = "draft"
	StatusScheduled AnnouncementStatus = "scheduled"
	StatusActive    AnnouncementStatus = "active"
	StatusInactive  AnnouncementStatus = "inactive"
)

type Announcement struct {
	ID             uint
	BusinessID     *uint
	CategoryID     uint
	Category       *AnnouncementCategory
	Title          string
	Message        string
	DisplayType    DisplayType
	FrequencyType  FrequencyType
	Priority       int
	IsGlobal       bool
	Status         AnnouncementStatus
	StartsAt       *time.Time
	EndsAt         *time.Time
	ForceRedisplay bool
	CreatedByID    uint
	CreatedAt      time.Time
	UpdatedAt      time.Time

	Images  []AnnouncementImage
	Links   []AnnouncementLink
	Targets []AnnouncementTarget
}

type AnnouncementImage struct {
	ID             uint
	AnnouncementID uint
	ImageURL       string
	SortOrder      int
}

type AnnouncementLink struct {
	ID             uint
	AnnouncementID uint
	Label          string
	URL            string
	SortOrder      int
}

type AnnouncementTarget struct {
	ID             uint
	AnnouncementID uint
	BusinessID     uint
}
