package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

const (
	// TTL de 25h (1h buffer sobre ventana de 24h de WhatsApp)
	conversationTTL = 25 * time.Hour
	// TTL corto para conversaciones expiradas (cleanup)
	expiredTTL = 1 * time.Hour

	// Prefijos de claves Redis
	convKeyPrefix      = "whatsapp:conv:"
	convPhoneOrderIdx  = "whatsapp:conv:idx:po:"
	convActivePhoneIdx = "whatsapp:conv:idx:active:"
)

// conversationCache implementa IConversationCache usando Redis
type conversationCache struct {
	redis redisclient.IRedis
	log   log.ILogger
}

// newConversationCache crea una nueva instancia del cache de conversaciones
func newConversationCache(redis redisclient.IRedis, logger log.ILogger) ports.IConversationCache {
	return &conversationCache{
		redis: redis,
		log:   logger.WithModule("whatsapp-conversation-cache"),
	}
}

// convKey genera la clave principal: whatsapp:conv:{id}
func convKey(id string) string {
	return convKeyPrefix + id
}

// phoneOrderIdxKey genera el índice por teléfono+orden: whatsapp:conv:idx:po:{phone}:{order}
func phoneOrderIdxKey(phone, order string) string {
	return convPhoneOrderIdx + phone + ":" + order
}

// activePhoneIdxKey genera el índice por teléfono activo: whatsapp:conv:idx:active:{phone}
func activePhoneIdxKey(phone string) string {
	return convActivePhoneIdx + phone
}

// isTerminalState verifica si el estado es terminal
func isTerminalState(state entities.ConversationState) bool {
	return state == entities.StateCompleted || state == entities.StateHandoffToHuman
}

// cachedConversation es la representación JSON de una conversación en Redis
type cachedConversation struct {
	ID             string                 `json:"id"`
	PhoneNumber    string                 `json:"phone_number"`
	OrderNumber    string                 `json:"order_number"`
	BusinessID     uint                   `json:"business_id"`
	CurrentState   string                 `json:"current_state"`
	LastMessageID  string                 `json:"last_message_id"`
	LastTemplateID string                 `json:"last_template_id"`
	Metadata       map[string]interface{} `json:"metadata"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	ExpiresAt      time.Time              `json:"expires_at"`
}

func toCached(c *entities.Conversation) *cachedConversation {
	return &cachedConversation{
		ID:             c.ID,
		PhoneNumber:    c.PhoneNumber,
		OrderNumber:    c.OrderNumber,
		BusinessID:     c.BusinessID,
		CurrentState:   string(c.CurrentState),
		LastMessageID:  c.LastMessageID,
		LastTemplateID: c.LastTemplateID,
		Metadata:       c.Metadata,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
		ExpiresAt:      c.ExpiresAt,
	}
}

func (cc *cachedConversation) toDomain() *entities.Conversation {
	return &entities.Conversation{
		ID:             cc.ID,
		PhoneNumber:    cc.PhoneNumber,
		OrderNumber:    cc.OrderNumber,
		BusinessID:     cc.BusinessID,
		CurrentState:   entities.ConversationState(cc.CurrentState),
		LastMessageID:  cc.LastMessageID,
		LastTemplateID: cc.LastTemplateID,
		Metadata:       cc.Metadata,
		CreatedAt:      cc.CreatedAt,
		UpdatedAt:      cc.UpdatedAt,
		ExpiresAt:      cc.ExpiresAt,
	}
}

// Save guarda una conversación en cache con clave principal e índices
func (c *conversationCache) Save(ctx context.Context, conversation *entities.Conversation) error {
	// Generar ID si no tiene
	if conversation.ID == "" {
		conversation.ID = uuid.New().String()
	}

	data, err := json.Marshal(toCached(conversation))
	if err != nil {
		return fmt.Errorf("error serializando conversación: %w", err)
	}

	// 1. Escribir clave principal
	if err := c.redis.Set(ctx, convKey(conversation.ID), string(data), conversationTTL); err != nil {
		return fmt.Errorf("error guardando conversación en cache: %w", err)
	}

	// 2. Escribir índice phone+order
	if err := c.redis.Set(ctx, phoneOrderIdxKey(conversation.PhoneNumber, conversation.OrderNumber), conversation.ID, conversationTTL); err != nil {
		c.log.Error(ctx).Err(err).Msg("Error guardando índice phone+order")
	}

	// 3. Escribir/eliminar índice active según estado
	if isTerminalState(conversation.CurrentState) {
		// Estado terminal: eliminar índice active
		if err := c.redis.Delete(ctx, activePhoneIdxKey(conversation.PhoneNumber)); err != nil {
			c.log.Error(ctx).Err(err).Msg("Error eliminando índice active")
		}
	} else {
		// Estado activo: escribir índice active
		if err := c.redis.Set(ctx, activePhoneIdxKey(conversation.PhoneNumber), conversation.ID, conversationTTL); err != nil {
			c.log.Error(ctx).Err(err).Msg("Error guardando índice active")
		}
	}

	return nil
}

// GetByID obtiene una conversación del cache por su ID
func (c *conversationCache) GetByID(ctx context.Context, id string) (*entities.Conversation, error) {
	data, err := c.redis.Get(ctx, convKey(id))
	if err != nil {
		return nil, fmt.Errorf("conversación no encontrada en cache: %s", id)
	}

	var cached cachedConversation
	if err := json.Unmarshal([]byte(data), &cached); err != nil {
		return nil, fmt.Errorf("error deserializando conversación: %w", err)
	}

	return cached.toDomain(), nil
}

// GetByPhoneAndOrder obtiene una conversación por teléfono + número de orden
func (c *conversationCache) GetByPhoneAndOrder(ctx context.Context, phoneNumber, orderNumber string) (*entities.Conversation, error) {
	convID, err := c.redis.Get(ctx, phoneOrderIdxKey(phoneNumber, orderNumber))
	if err != nil {
		return nil, fmt.Errorf("conversación no encontrada para phone=%s order=%s", phoneNumber, orderNumber)
	}

	return c.GetByID(ctx, convID)
}

// GetActiveByPhone obtiene la conversación activa de un teléfono
func (c *conversationCache) GetActiveByPhone(ctx context.Context, phoneNumber string) (*entities.Conversation, error) {
	convID, err := c.redis.Get(ctx, activePhoneIdxKey(phoneNumber))
	if err != nil {
		return nil, fmt.Errorf("no hay conversación activa para phone=%s", phoneNumber)
	}

	conv, err := c.GetByID(ctx, convID)
	if err != nil {
		// Índice stale: limpiar
		c.redis.Delete(ctx, activePhoneIdxKey(phoneNumber))
		return nil, err
	}

	// Validar que no ha expirado
	if conv.IsExpired() {
		c.redis.Delete(ctx, activePhoneIdxKey(phoneNumber))
		return nil, fmt.Errorf("conversación activa expirada para phone=%s", phoneNumber)
	}

	return conv, nil
}

// Expire marca una conversación como expirada y limpia índices activos
func (c *conversationCache) Expire(ctx context.Context, id string) error {
	conv, err := c.GetByID(ctx, id)
	if err != nil {
		return err
	}

	conv.ExpiresAt = time.Now()
	conv.UpdatedAt = time.Now()

	data, err := json.Marshal(toCached(conv))
	if err != nil {
		return fmt.Errorf("error serializando conversación expirada: %w", err)
	}

	// Actualizar con TTL corto
	if err := c.redis.Set(ctx, convKey(id), string(data), expiredTTL); err != nil {
		return fmt.Errorf("error expirando conversación en cache: %w", err)
	}

	// Eliminar índice active
	c.redis.Delete(ctx, activePhoneIdxKey(conv.PhoneNumber))

	return nil
}
