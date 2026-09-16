package entities

import "time"

const (
	FeedbackNone     = 0
	FeedbackPositive = 1
	FeedbackNegative = -1
)

const (
	ErrorCodeRateLimited = "rate_limited"
	ErrorCodeUnavailable = "assistant_unavailable"
)

type MessageRecord struct {
	ID               string
	ConversationID   string
	BusinessID       *uint
	UserID           uint
	Pathname         string
	Question         string
	Answer           string
	DestinationKey   string
	DestinationRoute string
	ErrorCode        string
	Model            string
	InputTokens      int
	OutputTokens     int
	LatencyMs        int
	CreatedAt        time.Time
}

type ReviewMessage struct {
	MessageRecord
	BusinessName string
	UserName     string
	UserEmail    string
	Feedback     int
	FeedbackAt   *time.Time
	ClickedAt    *time.Time
}

type DestinationCount struct {
	Key   string
	Count int64
}

type ReviewSummary struct {
	Messages         int64
	Conversations    int64
	Users            int64
	InputTokens      int64
	OutputTokens     int64
	NoDestination    int64
	Errors           int64
	Positive         int64
	Negative         int64
	WithDestination  int64
	Clicked          int64
	EstimatedCostUSD float64
	TopDestinations  []DestinationCount
}
