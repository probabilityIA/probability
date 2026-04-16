package domain

// ConversationRepository almacena las conversaciones simuladas en memoria
type ConversationRepository struct {
	conversations map[string]*Conversation
	messageLogs   map[string][]*MessageLog
}

// NewConversationRepository crea una nueva instancia del repositorio
func NewConversationRepository() *ConversationRepository {
	return &ConversationRepository{
		conversations: make(map[string]*Conversation),
		messageLogs:   make(map[string][]*MessageLog),
	}
}

// SaveConversation guarda o actualiza una conversación
func (r *ConversationRepository) SaveConversation(conv *Conversation) {
	r.conversations[conv.PhoneNumber] = conv
}

// GetConversation obtiene una conversación por número de teléfono
func (r *ConversationRepository) GetConversation(phoneNumber string) (*Conversation, bool) {
	conv, exists := r.conversations[phoneNumber]
	return conv, exists
}

// SaveMessage guarda un mensaje en el log
func (r *ConversationRepository) SaveMessage(msg *MessageLog) {
	r.messageLogs[msg.ConversationID] = append(r.messageLogs[msg.ConversationID], msg)
}

// GetMessages obtiene todos los mensajes de una conversación
func (r *ConversationRepository) GetMessages(conversationID string) []*MessageLog {
	return r.messageLogs[conversationID]
}

// GetAllConversations retorna todas las conversaciones
func (r *ConversationRepository) GetAllConversations() []*Conversation {
	conversations := make([]*Conversation, 0, len(r.conversations))
	for _, conv := range r.conversations {
		conversations = append(conversations, conv)
	}
	return conversations
}
