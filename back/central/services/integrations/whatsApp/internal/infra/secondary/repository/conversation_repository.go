package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

type ConversationRepository struct {
	db  db.IDatabase
	log log.ILogger
}

// Create crea una nueva conversación en la base de datos
func (r *ConversationRepository) Create(ctx context.Context, conversation *entities.Conversation) error {
	r.log.Info(ctx).
		Str("phone_number", conversation.PhoneNumber).
		Str("order_number", conversation.OrderNumber).
		Str("state", string(conversation.CurrentState)).
		Msg("[WhatsApp Repository] - creando conversación")

	// Convertir entidad del dominio a modelo GORM
	model, err := mappers.ConversationDomainToModel(conversation)
	if err != nil {
		r.log.Error(ctx).Err(err).Msg("[WhatsApp Repository] - error convirtiendo conversación")
		return fmt.Errorf("error al convertir conversación: %w", err)
	}

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		r.log.Error(ctx).Err(err).
			Str("phone_number", conversation.PhoneNumber).
			Str("order_number", conversation.OrderNumber).
			Msg("[WhatsApp Repository] - error creando conversación")
		return fmt.Errorf("error al crear conversación: %w", err)
	}

	// Actualizar la entidad con el ID generado
	conversation.ID = model.ID.String()

	r.log.Info(ctx).
		Str("id", conversation.ID).
		Str("phone_number", conversation.PhoneNumber).
		Msg("[WhatsApp Repository] - conversación creada exitosamente")

	return nil
}

// GetByID obtiene una conversación por su ID
func (r *ConversationRepository) GetByID(ctx context.Context, id string) (*entities.Conversation, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("ID inválido: %w", err)
	}

	var model models.WhatsAppConversation

	if err := r.db.Conn(ctx).
		Where("id = ?", parsedID).
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.log.Debug(ctx).Str("id", id).Msg("[WhatsApp Repository] - conversación no encontrada")
			return nil, fmt.Errorf("conversación no encontrada")
		}
		r.log.Error(ctx).Err(err).
			Str("id", id).
			Msg("[WhatsApp Repository] - error obteniendo conversación")
		return nil, fmt.Errorf("error al obtener conversación: %w", err)
	}

	return mappers.ConversationModelToDomain(&model), nil
}

// GetByPhoneAndOrder obtiene una conversación por número de teléfono y número de orden
func (r *ConversationRepository) GetByPhoneAndOrder(ctx context.Context, phoneNumber, orderNumber string) (*entities.Conversation, error) {
	var model models.WhatsAppConversation

	if err := r.db.Conn(ctx).
		Where("phone_number = ? AND order_number = ?", phoneNumber, orderNumber).
		Order("created_at DESC").
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.log.Debug(ctx).
				Str("phone_number", phoneNumber).
				Str("order_number", orderNumber).
				Msg("[WhatsApp Repository] - conversación no encontrada")
			return nil, &errors.ErrConversationNotFound{
				PhoneNumber: phoneNumber,
				OrderNumber: orderNumber,
			}
		}
		return nil, fmt.Errorf("error al buscar conversación: %w", err)
	}

	return mappers.ConversationModelToDomain(&model), nil
}

// GetActiveByPhone obtiene la conversación activa más reciente por número de teléfono
func (r *ConversationRepository) GetActiveByPhone(ctx context.Context, phoneNumber string) (*entities.Conversation, error) {
	var model models.WhatsAppConversation

	if err := r.db.Conn(ctx).
		Where("phone_number = ? AND expires_at > ? AND current_state NOT IN (?, ?)",
			phoneNumber,
			time.Now(),
			string(entities.StateCompleted),
			string(entities.StateHandoffToHuman)).
		Order("created_at DESC").
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.log.Debug(ctx).
				Str("phone_number", phoneNumber).
				Msg("[WhatsApp Repository] - no hay conversación activa")
			return nil, fmt.Errorf("no hay conversación activa")
		}
		return nil, fmt.Errorf("error al buscar conversación activa: %w", err)
	}

	return mappers.ConversationModelToDomain(&model), nil
}

// Update actualiza una conversación existente
func (r *ConversationRepository) Update(ctx context.Context, conversation *entities.Conversation) error {
	r.log.Info(ctx).
		Str("id", conversation.ID).
		Str("state", string(conversation.CurrentState)).
		Msg("[WhatsApp Repository] - actualizando conversación")

	// Convertir entidad del dominio a modelo GORM
	model, err := mappers.ConversationDomainToModel(conversation)
	if err != nil {
		r.log.Error(ctx).Err(err).Msg("[WhatsApp Repository] - error convirtiendo conversación")
		return fmt.Errorf("error al convertir conversación: %w", err)
	}

	if err := r.db.Conn(ctx).Save(model).Error; err != nil {
		r.log.Error(ctx).Err(err).
			Str("id", conversation.ID).
			Msg("[WhatsApp Repository] - error actualizando conversación")
		return fmt.Errorf("error al actualizar conversación: %w", err)
	}

	return nil
}

// Expire marca una conversación como expirada
func (r *ConversationRepository) Expire(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("ID inválido: %w", err)
	}

	r.log.Info(ctx).
		Str("id", id).
		Msg("[WhatsApp Repository] - expirando conversación")

	if err := r.db.Conn(ctx).
		Model(&models.WhatsAppConversation{}).
		Where("id = ?", parsedID).
		Update("expires_at", time.Now()).Error; err != nil {
		r.log.Error(ctx).Err(err).
			Str("id", id).
			Msg("[WhatsApp Repository] - error expirando conversación")
		return fmt.Errorf("error al expirar conversación: %w", err)
	}

	return nil
}

// Delete elimina una conversación por su ID
func (r *ConversationRepository) Delete(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("ID inválido: %w", err)
	}

	r.log.Info(ctx).
		Str("id", id).
		Msg("[WhatsApp Repository] - eliminando conversación")

	if err := r.db.Conn(ctx).
		Where("id = ?", parsedID).
		Delete(&models.WhatsAppConversation{}).Error; err != nil {
		r.log.Error(ctx).Err(err).
			Str("id", id).
			Msg("[WhatsApp Repository] - error eliminando conversación")
		return fmt.Errorf("error al eliminar conversación: %w", err)
	}

	return nil
}
