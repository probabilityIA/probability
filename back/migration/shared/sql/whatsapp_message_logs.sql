-- Migration: WhatsApp Message Logs
-- Description: Crea la tabla para persistir logs de mensajes enviados y recibidos en WhatsApp
-- Date: 2026-01-26

CREATE TABLE IF NOT EXISTS whatsapp_message_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL,
    direction VARCHAR(10) NOT NULL CHECK (direction IN ('outbound', 'inbound')),
    message_id VARCHAR(255) NOT NULL UNIQUE,
    template_name VARCHAR(100),
    content TEXT,
    status VARCHAR(20) NOT NULL CHECK (status IN ('sent', 'delivered', 'read', 'failed')),
    delivered_at TIMESTAMP,
    read_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_whatsapp_message_logs_conversation
        FOREIGN KEY (conversation_id)
        REFERENCES whatsapp_conversations(id)
        ON DELETE CASCADE
);

-- Índices para optimizar consultas
CREATE INDEX IF NOT EXISTS idx_whatsapp_message_logs_conversation
    ON whatsapp_message_logs (conversation_id);

CREATE INDEX IF NOT EXISTS idx_whatsapp_message_logs_message_id
    ON whatsapp_message_logs (message_id);

CREATE INDEX IF NOT EXISTS idx_whatsapp_message_logs_created_at
    ON whatsapp_message_logs (created_at);

CREATE INDEX IF NOT EXISTS idx_whatsapp_message_logs_status
    ON whatsapp_message_logs (status);

-- Comentarios para documentación
COMMENT ON TABLE whatsapp_message_logs IS 'Registro de todos los mensajes enviados y recibidos en conversaciones de WhatsApp';
COMMENT ON COLUMN whatsapp_message_logs.conversation_id IS 'ID de la conversación a la que pertenece el mensaje';
COMMENT ON COLUMN whatsapp_message_logs.direction IS 'Dirección del mensaje: outbound (enviado) o inbound (recibido)';
COMMENT ON COLUMN whatsapp_message_logs.message_id IS 'ID único del mensaje asignado por WhatsApp Cloud API';
COMMENT ON COLUMN whatsapp_message_logs.template_name IS 'Nombre de la plantilla usada (solo para mensajes outbound con templates)';
COMMENT ON COLUMN whatsapp_message_logs.content IS 'Contenido del mensaje (texto recibido del usuario o variables de la plantilla)';
COMMENT ON COLUMN whatsapp_message_logs.status IS 'Estado del mensaje: sent, delivered, read, failed';
COMMENT ON COLUMN whatsapp_message_logs.delivered_at IS 'Timestamp cuando el mensaje fue entregado al dispositivo del usuario';
COMMENT ON COLUMN whatsapp_message_logs.read_at IS 'Timestamp cuando el mensaje fue leído por el usuario';
