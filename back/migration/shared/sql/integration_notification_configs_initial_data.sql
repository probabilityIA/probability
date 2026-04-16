-- ============================================
-- DATOS INICIALES PARA integration_notification_configs
-- ============================================
-- Esta tabla se usa para configuraciones de notificaciones por integración
-- (WhatsApp, Email, SMS) que se envían a través de integraciones externas

-- IMPORTANTE: Ajusta los integration_id según tus integraciones existentes
-- Puedes ver tus integraciones con: 
-- SELECT id, name, code, integration_type_id FROM integrations;

-- Ejemplo de configuración para WhatsApp cuando cambia el estado de orden
INSERT INTO integration_notification_configs (
    integration_id,
    notification_type,
    is_active,
    conditions,
    config,
    description,
    priority,
    created_at,
    updated_at
) VALUES
-- Configuración WhatsApp para orden en entrega (ajusta integration_id)
(
    1, -- integration_id (debe ser el ID de una integración de WhatsApp)
    'whatsapp',
    true,
    '{"trigger": "order_status_change", "statuses": ["en_entrega"]}'::jsonb,
    '{"template_id": "order_status_update", "language": "es", "recipient_type": "customer"}'::jsonb,
    'Enviar notificación WhatsApp cuando la orden esté en entrega',
    10,
    NOW(),
    NOW()
),
-- Configuración WhatsApp para orden entregada
(
    1, -- integration_id (debe ser el ID de una integración de WhatsApp)
    'whatsapp',
    true,
    '{"trigger": "order_status_change", "statuses": ["entregada"]}'::jsonb,
    '{"template_id": "order_delivered", "language": "es", "recipient_type": "customer"}'::jsonb,
    'Enviar notificación WhatsApp cuando la orden sea entregada',
    10,
    NOW(),
    NOW()
),
-- Configuración WhatsApp para orden creada
(
    1, -- integration_id
    'whatsapp',
    true,
    '{"trigger": "order_created"}'::jsonb,
    '{"template_id": "order_confirmation", "language": "es", "recipient_type": "customer"}'::jsonb,
    'Enviar notificación WhatsApp cuando se crea una orden',
    5,
    NOW(),
    NOW()
),
-- Configuración Email para orden creada (si tienes integración de email)
(
    2, -- integration_id (ajusta si tienes una integración de email)
    'email',
    true,
    '{"trigger": "order_created"}'::jsonb,
    '{"template": "order_confirmation", "subject": "Confirmación de tu orden", "recipient_type": "customer"}'::jsonb,
    'Enviar email de confirmación cuando se crea una orden',
    5,
    NOW(),
    NOW()
),
-- Configuración SMS para pago completado (si tienes integración SMS)
(
    3, -- integration_id (ajusta si tienes una integración SMS)
    'sms',
    true,
    '{"trigger": "payment_completed"}'::jsonb,
    '{"message_template": "Tu pago de ${{amount}} para la orden #{{order_number}} ha sido confirmado", "recipient_type": "customer"}'::jsonb,
    'Enviar SMS cuando se completa un pago',
    5,
    NOW(),
    NOW()
)
ON CONFLICT DO NOTHING;

-- ============================================
-- Query para verificar los datos insertados
-- ============================================
SELECT 
    inc.id,
    inc.integration_id,
    i.name as integration_name,
    i.code as integration_code,
    it.name as integration_type_name,
    inc.notification_type,
    inc.is_active,
    inc.conditions,
    inc.config,
    inc.description,
    inc.priority,
    inc.created_at,
    inc.updated_at
FROM integration_notification_configs inc
LEFT JOIN integrations i ON inc.integration_id = i.id
LEFT JOIN integration_types it ON i.integration_type_id = it.id
ORDER BY inc.integration_id, inc.priority DESC, inc.created_at DESC;

-- ============================================
-- Query para ver integraciones disponibles (para obtener los IDs correctos)
-- ============================================
SELECT 
    i.id,
    i.name,
    i.code,
    it.name as type_name,
    it.code as type_code,
    i.business_id,
    i.is_active
FROM integrations i
LEFT JOIN integration_types it ON i.integration_type_id = it.id
WHERE i.is_active = true
ORDER BY i.id;
