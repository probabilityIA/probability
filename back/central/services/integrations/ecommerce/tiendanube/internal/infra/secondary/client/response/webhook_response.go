package response

type Webhook struct {
	ID        int64  `json:"id"`
	Event     string `json:"event"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
