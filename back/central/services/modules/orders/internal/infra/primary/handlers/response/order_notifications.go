package response

import "time"

type OrderNotificationEvent struct {
	EventCode    string `json:"event_code"`
	Status       string `json:"status"`
	ResendAction string `json:"resend_action,omitempty"`
	CanResend    bool   `json:"can_resend"`
}

type OrderNotificationMessage struct {
	ID           string    `json:"id"`
	Channel      string    `json:"channel"`
	TemplateName string    `json:"template_name"`
	Direction    string    `json:"direction"`
	Status       string    `json:"status"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
}

type OrderNotifications struct {
	OrderID     string                     `json:"order_id"`
	OrderNumber string                     `json:"order_number"`
	Channel     string                     `json:"channel"`
	Sent        int                        `json:"sent"`
	Expected    int                        `json:"expected"`
	Failed      int                        `json:"failed"`
	State       string                     `json:"state"`
	Events      []OrderNotificationEvent   `json:"events"`
	Messages    []OrderNotificationMessage `json:"messages"`
}
