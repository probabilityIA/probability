package entities

import "time"

type MessageRole string

const (
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
)

type ChatMessage struct {
	Role MessageRole
	Text string
}

type Destination struct {
	Key         string
	Label       string
	Route       string
	Description string
	Guide       string
}

type NavigationCatalog struct {
	Allowed []Destination
	Denied  []string
}

func (c NavigationCatalog) Find(key string) (*Destination, bool) {
	for i := range c.Allowed {
		if c.Allowed[i].Key == key {
			return &c.Allowed[i], true
		}
	}
	return nil, false
}

func (c NavigationCatalog) Keys() []string {
	keys := make([]string, 0, len(c.Allowed))
	for _, d := range c.Allowed {
		keys = append(keys, d.Key)
	}
	return keys
}

type AssistantReply struct {
	Message     string
	Destination *Destination
}

type Usage struct {
	Count   int
	ResetAt *time.Time
}

type AssistantState struct {
	IntroSeen bool
	Limit     int
	Remaining int
	ResetAt   *time.Time
}
