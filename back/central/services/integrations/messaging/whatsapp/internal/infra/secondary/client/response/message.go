package response

type Message struct {
	ID string `json:"id"`
}

type SendMessageResponse struct {
	Messages []Message `json:"messages"`
}
