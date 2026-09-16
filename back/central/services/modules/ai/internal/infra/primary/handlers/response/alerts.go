package response

import "time"

type Alert struct {
	ID          string       `json:"id"`
	EventType   string       `json:"event_type"`
	Severity    string       `json:"severity"`
	Title       string       `json:"title"`
	Body        string       `json:"body"`
	Destination *Destination `json:"destination"`
	Reference   *Reference   `json:"reference"`
	Unread      bool         `json:"unread"`
	CreatedAt   time.Time    `json:"created_at"`
}

type Reference struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type Alerts struct {
	Data       []Alert `json:"data"`
	Total      int64   `json:"total"`
	Page       int     `json:"page"`
	PageSize   int     `json:"page_size"`
	TotalPages int     `json:"total_pages"`
}

type AlertsUnread struct {
	Count  int64  `json:"count"`
	Latest *Alert `json:"latest"`
}
