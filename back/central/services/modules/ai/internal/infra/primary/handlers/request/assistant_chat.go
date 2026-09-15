package request

type AssistantMessage struct {
	Role string `json:"role" binding:"required,oneof=user assistant"`
	Text string `json:"text" binding:"required"`
}

type AssistantChat struct {
	Messages []AssistantMessage `json:"messages" binding:"required,min=1,max=40,dive"`
}
