package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

const (
	platformCredsRedisKey = "integration:platform_creds:2"
)

// aiForwarder reenvia mensajes de WhatsApp sin conversacion activa al modulo AI Sales
type aiForwarder struct {
	rabbit rabbitmq.IQueue
	redis  redisclient.IRedis
	log    log.ILogger
}

// NewAIForwarder crea un nuevo forwarder que publica mensajes a whatsapp.ai.incoming
func NewAIForwarder(rabbit rabbitmq.IQueue, redis redisclient.IRedis, logger log.ILogger) ports.IAIForwarder {
	return &aiForwarder{
		rabbit: rabbit,
		redis:  redis,
		log:    logger.WithModule("whatsapp-ai-forwarder"),
	}
}

// ForwardToAI lee platform_creds para verificar si AI esta habilitado y publica a whatsapp.ai.incoming
func (f *aiForwarder) ForwardToAI(ctx context.Context, phoneNumber, messageText, messageID, messageType string) error {
	// Solo procesar mensajes de texto
	if messageType != "text" {
		f.log.Debug(ctx).
			Str("type", messageType).
			Msg("Mensaje no-texto ignorado por AI forwarder")
		return nil
	}

	// Leer platform_creds de Redis para verificar ai_sales_enabled
	platJSON, err := f.redis.Get(ctx, platformCredsRedisKey)
	if err != nil {
		f.log.Debug(ctx).
			Err(err).
			Msg("No se pudieron leer platform_creds, AI forwarder deshabilitado")
		return nil
	}

	var platCreds map[string]any
	if err := json.Unmarshal([]byte(platJSON), &platCreds); err != nil {
		f.log.Error(ctx).Err(err).Msg("Error deserializando platform_creds para AI")
		return nil
	}

	// Verificar ai_sales_enabled
	enabled, _ := platCreds["ai_sales_enabled"].(bool)
	if !enabled {
		f.log.Debug(ctx).Msg("AI Sales no habilitado en platform_creds")
		return nil
	}

	// Obtener business_id demo
	var businessID uint
	switch v := platCreds["ai_sales_demo_business_id"].(type) {
	case float64:
		businessID = uint(v)
	default:
		businessID = 1 // Default fase 1
	}

	// Construir DTO para la cola
	dto := map[string]any{
		"PhoneNumber": phoneNumber,
		"MessageText": messageText,
		"MessageID":   messageID,
		"MessageType": messageType,
		"BusinessID":  businessID,
		"Timestamp":   time.Now().Unix(),
	}

	body, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("error serializando mensaje AI: %w", err)
	}

	if err := f.rabbit.Publish(ctx, rabbitmq.QueueWhatsAppAIIncoming, body); err != nil {
		f.log.Error(ctx).Err(err).
			Str("phone", phoneNumber).
			Msg("Error publicando mensaje a cola AI incoming")
		return fmt.Errorf("error publicando a %s: %w", rabbitmq.QueueWhatsAppAIIncoming, err)
	}

	f.log.Info(ctx).
		Str("phone", phoneNumber).
		Uint("business_id", businessID).
		Msg("Mensaje reenviado a AI Sales")

	return nil
}
