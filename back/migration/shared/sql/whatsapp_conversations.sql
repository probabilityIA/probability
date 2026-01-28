-- Migration: WhatsApp Conversations
-- Description: Crea la tabla para persistir conversaciones activas de WhatsApp
-- Date: 2026-01-26

CREATE TABLE IF NOT EXISTS whatsapp_conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone_number VARCHAR(20) NOT NULL,
    order_number VARCHAR(100) NOT NULL,
    business_id INT NOT NULL,
    current_state VARCHAR(50) NOT NULL,
    last_message_id VARCHAR(255),
    last_template_id VARCHAR(100),
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    CONSTRAINT idx_phone_order UNIQUE (phone_number, order_number)
);

-- Índices para optimizar consultas
CREATE INDEX IF NOT EXISTS idx_whatsapp_conversations_phone_order
    ON whatsapp_conversations (phone_number, order_number);

CREATE INDEX IF NOT EXISTS idx_whatsapp_conversations_expires_at
    ON whatsapp_conversations (expires_at);

CREATE INDEX IF NOT EXISTS idx_whatsapp_conversations_business_id
    ON whatsapp_conversations (business_id);

CREATE INDEX IF NOT EXISTS idx_whatsapp_conversations_current_state
    ON whatsapp_conversations (current_state);

-- Trigger para actualizar updated_at automáticamente
CREATE OR REPLACE FUNCTION update_whatsapp_conversations_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_whatsapp_conversations_updated_at
    BEFORE UPDATE ON whatsapp_conversations
    FOR EACH ROW
    EXECUTE FUNCTION update_whatsapp_conversations_updated_at();

-- Comentarios para documentación
COMMENT ON TABLE whatsapp_conversations IS 'Almacena conversaciones activas de WhatsApp con su estado actual y metadata';
COMMENT ON COLUMN whatsapp_conversations.phone_number IS 'Número de teléfono del cliente en formato internacional (+573001234567)';
COMMENT ON COLUMN whatsapp_conversations.order_number IS 'Número de orden asociado a la conversación';
COMMENT ON COLUMN whatsapp_conversations.current_state IS 'Estado actual del flujo conversacional (START, AWAITING_CONFIRMATION, etc)';
COMMENT ON COLUMN whatsapp_conversations.expires_at IS 'Timestamp de expiración de la ventana de servicio (24h desde creación)';
COMMENT ON COLUMN whatsapp_conversations.metadata IS 'Datos adicionales en formato JSON (variables del pedido, contexto, etc)';
