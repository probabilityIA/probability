package domain

// IncomingMessageDTO mensaje recibido de la cola whatsapp.ai.incoming
type IncomingMessageDTO struct {
	PhoneNumber string
	MessageText string
	MessageID   string
	MessageType string
	BusinessID  uint
	Timestamp   int64
}

// AIResponseDTO respuesta a publicar en whatsapp.ai.response
type AIResponseDTO struct {
	PhoneNumber  string
	ResponseText string
	BusinessID   uint
	SessionID    string
	Timestamp    int64
}
