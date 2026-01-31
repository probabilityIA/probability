-- ============================================
-- MIGRATION: Notification System Refactor
-- Fecha: 2026-01-30
-- Descripción: Refactorización completa del sistema de notificaciones
--              con arquitectura jerárquica (NotificationType -> NotificationEventType)
-- ============================================

BEGIN;

-- ============================================
-- PASO 1: Insertar Notification Types (Tipos de Notificación)
-- ============================================

INSERT INTO notification_types (name, code, description, icon, is_active, config_schema, created_at, updated_at)
VALUES
    ('SSE', 'sse', 'Server-Sent Events para notificaciones en tiempo real en el panel administrativo', 'bell', true, '{"required_fields": [], "optional_fields": []}', NOW(), NOW()),
    ('WhatsApp', 'whatsapp', 'Mensajes de WhatsApp Business para clientes', 'message-circle', true, '{"required_fields": ["template_id"], "optional_fields": ["language"]}', NOW(), NOW()),
    ('Email', 'email', 'Notificaciones por correo electrónico', 'mail', true, '{"required_fields": ["template"], "optional_fields": ["subject", "reply_to"]}', NOW(), NOW()),
    ('SMS', 'sms', 'Mensajes de texto SMS', 'smartphone', false, '{"required_fields": ["message"], "optional_fields": []}', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- ============================================
-- PASO 2: Insertar Notification Event Types (Eventos por Tipo)
-- ============================================

-- Eventos para SSE (notificaciones internas del panel)
INSERT INTO notification_event_types (notification_type_id, event_code, event_name, description, is_active, created_at, updated_at)
SELECT
    nt.id,
    'order.created',
    'Nueva Orden',
    'Notificación cuando se crea una nueva orden',
    true,
    NOW(),
    NOW()
FROM notification_types nt WHERE nt.code = 'sse'
ON CONFLICT (notification_type_id, event_code) DO NOTHING;

INSERT INTO notification_event_types (notification_type_id, event_code, event_name, description, is_active, created_at, updated_at)
SELECT
    nt.id,
    'order.status_changed',
    'Cambio de Estado de Orden',
    'Notificación cuando cambia el estado de una orden',
    true,
    NOW(),
    NOW()
FROM notification_types nt WHERE nt.code = 'sse'
ON CONFLICT (notification_type_id, event_code) DO NOTHING;

-- Eventos para WhatsApp (notificaciones a clientes)
INSERT INTO notification_event_types (notification_type_id, event_code, event_name, description, template_config, is_active, created_at, updated_at)
SELECT
    nt.id,
    'order.created',
    'Confirmación de Pedido',
    'Confirmación enviada al cliente cuando se crea un pedido',
    '{"default_template_id": "order_confirmation", "template_language": "es"}',
    true,
    NOW(),
    NOW()
FROM notification_types nt WHERE nt.code = 'whatsapp'
ON CONFLICT (notification_type_id, event_code) DO NOTHING;

INSERT INTO notification_event_types (notification_type_id, event_code, event_name, description, template_config, is_active, created_at, updated_at)
SELECT
    nt.id,
    'order.shipped',
    'Pedido Enviado',
    'Notificación cuando el pedido ha sido enviado',
    '{"default_template_id": "order_shipped", "template_language": "es"}',
    true,
    NOW(),
    NOW()
FROM notification_types nt WHERE nt.code = 'whatsapp'
ON CONFLICT (notification_type_id, event_code) DO NOTHING;

INSERT INTO notification_event_types (notification_type_id, event_code, event_name, description, template_config, is_active, created_at, updated_at)
SELECT
    nt.id,
    'order.delivered',
    'Pedido Entregado',
    'Notificación cuando el pedido ha sido entregado',
    '{"default_template_id": "order_delivered", "template_language": "es"}',
    true,
    NOW(),
    NOW()
FROM notification_types nt WHERE nt.code = 'whatsapp'
ON CONFLICT (notification_type_id, event_code) DO NOTHING;

INSERT INTO notification_event_types (notification_type_id, event_code, event_name, description, template_config, is_active, created_at, updated_at)
SELECT
    nt.id,
    'order.canceled',
    'Pedido Cancelado',
    'Notificación cuando un pedido es cancelado',
    '{"default_template_id": "order_canceled", "template_language": "es"}',
    true,
    NOW(),
    NOW()
FROM notification_types nt WHERE nt.code = 'whatsapp'
ON CONFLICT (notification_type_id, event_code) DO NOTHING;

INSERT INTO notification_event_types (notification_type_id, event_code, event_name, description, template_config, is_active, created_at, updated_at)
SELECT
    nt.id,
    'invoice.created',
    'Factura Generada',
    'Notificación cuando se genera una factura',
    '{"default_template_id": "invoice_created", "template_language": "es"}',
    true,
    NOW(),
    NOW()
FROM notification_types nt WHERE nt.code = 'whatsapp'
ON CONFLICT (notification_type_id, event_code) DO NOTHING;

-- Eventos para Email
INSERT INTO notification_event_types (notification_type_id, event_code, event_name, description, template_config, is_active, created_at, updated_at)
SELECT
    nt.id,
    'order.created',
    'Confirmación de Pedido',
    'Email de confirmación cuando se crea un pedido',
    '{"default_subject": "Confirmación de tu pedido", "template_file": "order_confirmation.html"}',
    true,
    NOW(),
    NOW()
FROM notification_types nt WHERE nt.code = 'email'
ON CONFLICT (notification_type_id, event_code) DO NOTHING;

INSERT INTO notification_event_types (notification_type_id, event_code, event_name, description, template_config, is_active, created_at, updated_at)
SELECT
    nt.id,
    'order.shipped',
    'Pedido Enviado',
    'Email cuando el pedido ha sido enviado',
    '{"default_subject": "Tu pedido ha sido enviado", "template_file": "order_shipped.html"}',
    true,
    NOW(),
    NOW()
FROM notification_types nt WHERE nt.code = 'email'
ON CONFLICT (notification_type_id, event_code) DO NOTHING;

-- ============================================
-- PASO 3: Migrar datos existentes
-- ============================================

-- Primero, necesitamos agregar las nuevas columnas a business_notification_configs
-- NOTA: Estas columnas ya deberían existir por el AutoMigrate, pero las agregamos
--       manualmente por seguridad en caso de que AutoMigrate no se haya ejecutado

-- Agregar columna integration_id (FK a integrations - la integración origen)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'business_notification_configs'
        AND column_name = 'integration_id'
    ) THEN
        ALTER TABLE business_notification_configs
        ADD COLUMN integration_id BIGINT;
    END IF;
END $$;

-- Agregar columna notification_type_id (FK a notification_types)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'business_notification_configs'
        AND column_name = 'notification_type_id'
    ) THEN
        ALTER TABLE business_notification_configs
        ADD COLUMN notification_type_id BIGINT;
    END IF;
END $$;

-- Agregar columna notification_event_type_id (FK a notification_event_types)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'business_notification_configs'
        AND column_name = 'notification_event_type_id'
    ) THEN
        ALTER TABLE business_notification_configs
        ADD COLUMN notification_event_type_id BIGINT;
    END IF;
END $$;

-- Migrar datos existentes:
-- 1. Asignar integration_id: Buscar la primera integración activa del business
-- 2. Asignar notification_type_id: Extraer del campo channels (si es ["sse"] -> tipo SSE)
-- 3. Asignar notification_event_type_id: Mapear event_type string a notification_event_types

-- Actualizar integration_id: asignar la primera integración del business
UPDATE business_notification_configs bnc
SET integration_id = (
    SELECT i.id
    FROM integrations i
    WHERE i.business_id = bnc.business_id
      AND i.is_active = true
      AND i.deleted_at IS NULL
    LIMIT 1
)
WHERE integration_id IS NULL;

-- Actualizar notification_type_id: extraer de channels
-- Si channels contiene "sse", asignar tipo SSE
-- Si channels contiene "whatsapp", asignar tipo WhatsApp
-- Por defecto: SSE
UPDATE business_notification_configs bnc
SET notification_type_id = (
    SELECT nt.id FROM notification_types nt
    WHERE nt.code = CASE
        WHEN bnc.channels::text LIKE '%"sse"%' THEN 'sse'
        WHEN bnc.channels::text LIKE '%"whatsapp"%' THEN 'whatsapp'
        WHEN bnc.channels::text LIKE '%"email"%' THEN 'email'
        ELSE 'sse'  -- Default
    END
    LIMIT 1
)
WHERE notification_type_id IS NULL;

-- Actualizar notification_event_type_id: mapear event_type a notification_event_types
-- Esto depende del notification_type_id asignado
UPDATE business_notification_configs bnc
SET notification_event_type_id = (
    SELECT net.id
    FROM notification_event_types net
    WHERE net.notification_type_id = bnc.notification_type_id
      AND net.event_code = bnc.event_type
    LIMIT 1
)
WHERE notification_event_type_id IS NULL;

-- ============================================
-- PASO 4: Agregar Foreign Keys
-- ============================================

-- FK integration_id -> integrations
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_business_notification_configs_integration'
    ) THEN
        ALTER TABLE business_notification_configs
        ADD CONSTRAINT fk_business_notification_configs_integration
        FOREIGN KEY (integration_id) REFERENCES integrations(id)
        ON DELETE CASCADE ON UPDATE CASCADE;
    END IF;
END $$;

-- FK notification_type_id -> notification_types
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_business_notification_configs_notification_type'
    ) THEN
        ALTER TABLE business_notification_configs
        ADD CONSTRAINT fk_business_notification_configs_notification_type
        FOREIGN KEY (notification_type_id) REFERENCES notification_types(id)
        ON DELETE RESTRICT ON UPDATE CASCADE;
    END IF;
END $$;

-- FK notification_event_type_id -> notification_event_types
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_business_notification_configs_notification_event_type'
    ) THEN
        ALTER TABLE business_notification_configs
        ADD CONSTRAINT fk_business_notification_configs_notification_event_type
        FOREIGN KEY (notification_event_type_id) REFERENCES notification_event_types(id)
        ON DELETE RESTRICT ON UPDATE CASCADE;
    END IF;
END $$;

-- ============================================
-- PASO 5: Crear índices
-- ============================================

CREATE INDEX IF NOT EXISTS idx_bnc_integration_id ON business_notification_configs(integration_id);
CREATE INDEX IF NOT EXISTS idx_bnc_notification_type_id ON business_notification_configs(notification_type_id);
CREATE INDEX IF NOT EXISTS idx_bnc_notification_event_type_id ON business_notification_configs(notification_event_type_id);

-- Índice único compuesto para evitar duplicados
-- Una configuración debe ser única por (integration_id, notification_type_id, notification_event_type_id)
CREATE UNIQUE INDEX IF NOT EXISTS idx_bnc_unique_config
ON business_notification_configs(integration_id, notification_type_id, notification_event_type_id)
WHERE deleted_at IS NULL;

-- ============================================
-- PASO 6: Eliminar columna channels (deprecada)
-- ============================================

-- NOTA: Comentado por seguridad - descomentar cuando se haya verificado que la migración funciona
-- ALTER TABLE business_notification_configs DROP COLUMN IF EXISTS channels;

-- ============================================
-- PASO 7: Hacer event_type nullable (deprecado pero mantener para referencia)
-- ============================================

ALTER TABLE business_notification_configs ALTER COLUMN event_type DROP NOT NULL;

-- ============================================
-- PASO 8: Eliminar unique index antiguo (business_id, event_type)
-- ============================================

DROP INDEX IF EXISTS idx_business_event_type;

COMMIT;

-- ============================================
-- VERIFICACIÓN POST-MIGRACIÓN
-- ============================================

-- Verificar notification_types
SELECT 'Notification Types:' as step;
SELECT id, name, code, is_active FROM notification_types ORDER BY id;

-- Verificar notification_event_types
SELECT 'Notification Event Types:' as step;
SELECT
    net.id,
    nt.name as notification_type,
    net.event_code,
    net.event_name,
    net.is_active
FROM notification_event_types net
JOIN notification_types nt ON net.notification_type_id = nt.id
ORDER BY nt.id, net.id;

-- Verificar business_notification_configs migradas
SELECT 'Business Notification Configs migradas:' as step;
SELECT
    bnc.id,
    bnc.business_id,
    i.name as integration,
    nt.name as notification_type,
    net.event_name,
    bnc.enabled,
    bnc.event_type as deprecated_event_type
FROM business_notification_configs bnc
LEFT JOIN integrations i ON bnc.integration_id = i.id
LEFT JOIN notification_types nt ON bnc.notification_type_id = nt.id
LEFT JOIN notification_event_types net ON bnc.notification_event_type_id = net.id
ORDER BY bnc.id;

-- Verificar configs sin migrar (integration_id NULL)
SELECT 'Configs sin migrar (ERROR):' as step;
SELECT id, business_id, event_type
FROM business_notification_configs
WHERE integration_id IS NULL
   OR notification_type_id IS NULL
   OR notification_event_type_id IS NULL;
