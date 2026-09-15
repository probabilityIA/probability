package response

import "time"

type Destination struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Route       string `json:"route"`
	Description string `json:"description"`
}

type AssistantReply struct {
	Message     string       `json:"message"`
	Destination *Destination `json:"destination"`
}

type AssistantState struct {
	IntroSeen bool       `json:"intro_seen"`
	Limit     int        `json:"limit"`
	Remaining int        `json:"remaining"`
	ResetAt   *time.Time `json:"reset_at"`
}
