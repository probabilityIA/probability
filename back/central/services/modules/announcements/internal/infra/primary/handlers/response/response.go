package response

import "time"

type AnnouncementResponse struct {
	ID             uint                `json:"id"`
	BusinessID     *uint               `json:"business_id"`
	CategoryID     uint                `json:"category_id"`
	Category       *CategoryResponse   `json:"category,omitempty"`
	Title          string              `json:"title"`
	Message        string              `json:"message"`
	DisplayType    string              `json:"display_type"`
	FrequencyType  string              `json:"frequency_type"`
	Priority       int                 `json:"priority"`
	IsGlobal       bool                `json:"is_global"`
	Status         string              `json:"status"`
	StartsAt       *time.Time          `json:"starts_at"`
	EndsAt         *time.Time          `json:"ends_at"`
	ForceRedisplay bool                `json:"force_redisplay"`
	CreatedByID    uint                `json:"created_by_id"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
	Images         []ImageResponse     `json:"images"`
	Links          []LinkResponse      `json:"links"`
	Targets        []TargetResponse    `json:"targets"`
}

type CategoryResponse struct {
	ID    uint   `json:"id"`
	Code  string `json:"code"`
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
}

type ImageResponse struct {
	ID        uint   `json:"id"`
	ImageURL  string `json:"image_url"`
	SortOrder int    `json:"sort_order"`
}

type LinkResponse struct {
	ID        uint   `json:"id"`
	Label     string `json:"label"`
	URL       string `json:"url"`
	SortOrder int    `json:"sort_order"`
}

type TargetResponse struct {
	ID         uint `json:"id"`
	BusinessID uint `json:"business_id"`
}

type StatsResponse struct {
	TotalViews       int64 `json:"total_views"`
	UniqueUsers      int64 `json:"unique_users"`
	TotalClicks      int64 `json:"total_clicks"`
	TotalAcceptances int64 `json:"total_acceptances"`
	TotalClosed      int64 `json:"total_closed"`
}

type PaginatedAnnouncementsResponse struct {
	Data       []AnnouncementResponse `json:"data"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"page_size"`
	TotalPages int                    `json:"total_pages"`
}
