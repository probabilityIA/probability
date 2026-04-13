package request

type CreateAnnouncementRequest struct {
	BusinessID    *uint  `json:"business_id"`
	CategoryID    uint   `json:"category_id" binding:"required"`
	Title         string `json:"title" binding:"required,max=255"`
	Message       string `json:"message"`
	DisplayType   string `json:"display_type" binding:"required"`
	FrequencyType string `json:"frequency_type" binding:"required"`
	Priority      int    `json:"priority"`
	IsGlobal      bool   `json:"is_global"`
	StartsAt      string `json:"starts_at"`
	EndsAt        string `json:"ends_at"`
	Links         []LinkRequest `json:"links"`
	TargetIDs     []uint        `json:"target_ids"`
}

type LinkRequest struct {
	Label     string `json:"label" binding:"required"`
	URL       string `json:"url" binding:"required"`
	SortOrder int    `json:"sort_order"`
}

type UpdateAnnouncementRequest struct {
	BusinessID    *uint  `json:"business_id"`
	CategoryID    uint   `json:"category_id" binding:"required"`
	Title         string `json:"title" binding:"required,max=255"`
	Message       string `json:"message"`
	DisplayType   string `json:"display_type" binding:"required"`
	FrequencyType string `json:"frequency_type" binding:"required"`
	Priority      int    `json:"priority"`
	IsGlobal      bool   `json:"is_global"`
	StartsAt      string `json:"starts_at"`
	EndsAt        string `json:"ends_at"`
	Links         []LinkRequest `json:"links"`
	TargetIDs     []uint        `json:"target_ids"`
}

type ChangeStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type RegisterViewRequest struct {
	Action string `json:"action" binding:"required"`
	LinkID *uint  `json:"link_id"`
}
