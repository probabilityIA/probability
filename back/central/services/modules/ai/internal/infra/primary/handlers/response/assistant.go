package response

import "time"

type Destination struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Route       string `json:"route"`
	Description string `json:"description"`
}

type AssistantReply struct {
	MessageID      string       `json:"message_id"`
	ConversationID string       `json:"conversation_id"`
	Message        string       `json:"message"`
	Destination    *Destination `json:"destination"`
}

type AssistantState struct {
	IntroSeen bool       `json:"intro_seen"`
	Limit     int        `json:"limit"`
	Remaining int        `json:"remaining"`
	ResetAt   *time.Time `json:"reset_at"`
}

type ReviewMessage struct {
	ID               string     `json:"id"`
	ConversationID   string     `json:"conversation_id"`
	BusinessID       *uint      `json:"business_id"`
	BusinessName     string     `json:"business_name"`
	UserID           uint       `json:"user_id"`
	UserName         string     `json:"user_name"`
	UserEmail        string     `json:"user_email"`
	Pathname         string     `json:"pathname"`
	Question         string     `json:"question"`
	Answer           string     `json:"answer"`
	DestinationKey   string     `json:"destination_key"`
	DestinationRoute string     `json:"destination_route"`
	ErrorCode        string     `json:"error_code"`
	Model            string     `json:"model"`
	InputTokens      int        `json:"input_tokens"`
	OutputTokens     int        `json:"output_tokens"`
	LatencyMs        int        `json:"latency_ms"`
	Feedback         int        `json:"feedback"`
	FeedbackAt       *time.Time `json:"feedback_at"`
	ClickedAt        *time.Time `json:"clicked_at"`
	CreatedAt        time.Time  `json:"created_at"`
}

type ReviewMessages struct {
	Data       []ReviewMessage `json:"data"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

type DestinationCount struct {
	Key   string `json:"key"`
	Count int64  `json:"count"`
}

type ReviewSummary struct {
	Messages         int64              `json:"messages"`
	Conversations    int64              `json:"conversations"`
	Users            int64              `json:"users"`
	InputTokens      int64              `json:"input_tokens"`
	OutputTokens     int64              `json:"output_tokens"`
	NoDestination    int64              `json:"no_destination"`
	Errors           int64              `json:"errors"`
	Positive         int64              `json:"positive"`
	Negative         int64              `json:"negative"`
	WithDestination  int64              `json:"with_destination"`
	Clicked          int64              `json:"clicked"`
	EstimatedCostUSD float64            `json:"estimated_cost_usd"`
	TopDestinations  []DestinationCount `json:"top_destinations"`
}
