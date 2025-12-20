-- ============================================
-- DATOS INICIALES PARA NOTIFICATION CONFIGS
-- ============================================
-- Esta tabla (business_notification_configs) se usa para configuraciones de 
-- notificaciones internas (SSE) por negocio.
-- Las notificaciones se activan cuando ocurren eventos relacionados con órdenes de Probability
-- ============================================

-- ============================================
-- 1. CONFIGURACIÓN PARA ÓRDENES CREADAS
-- ============================================
INSERT INTO business_notification_configs (
    business_id,
    event_type,
    enabled,
    channels,
    filters,
    description,
    created_at,
    updated_at
) VALUES
(
    1, -- business_id (ajusta según tu negocio)
    'order.created',
    true,
    '["sse"]'::jsonb,
    '{}'::jsonb,
    'Notificaciones cuando se crea una nueva orden',
    NOW(),
    NOW()
)
ON CONFLICT (business_id, event_type) DO UPDATE 
SET enabled = EXCLUDED.enabled,
    channels = EXCLUDED.channels,
    filters = EXCLUDED.filters,
    description = EXCLUDED.description,
    updated_at = NOW();

-- ============================================
-- 2. CONFIGURACIÓN PARA CAMBIOS DE ESTADO
-- ============================================
-- Para este tipo de evento, especificaremos qué estados deben disparar notificaciones
-- usando la tabla business_notification_config_order_statuses
INSERT INTO business_notification_configs (
    business_id,
    event_type,
    enabled,
    channels,
    filters,
    description,
    created_at,
    updated_at
) VALUES
(
    1, -- business_id (ajusta según tu negocio)
    'order.status_changed',
    true,
    '["sse"]'::jsonb,
    '{}'::jsonb,
    'Notificaciones para cambios de estado de órdenes',
    NOW(),
    NOW()
)
ON CONFLICT (business_id, event_type) DO UPDATE 
SET enabled = EXCLUDED.enabled,
    channels = EXCLUDED.channels,
    filters = EXCLUDED.filters,
    description = EXCLUDED.description,
    updated_at = NOW();

-- ============================================
-- 3. ASOCIAR ESTADOS A LA CONFIGURACIÓN
-- ============================================
-- Obtener el ID de la configuración que acabamos de insertar/actualizar
-- y asociar los estados de orden que queremos notificar

-- Primero, obtenemos el ID de la configuración
DO $$
DECLARE
    config_id INTEGER;
BEGIN
    -- Obtener el ID de la configuración de order.status_changed para business_id = 1
    SELECT id INTO config_id 
    FROM business_notification_configs 
    WHERE business_id = 1 AND event_type = 'order.status_changed'
    LIMIT 1;

    -- Si la configuración existe, asociar los estados que queremos notificar
    IF config_id IS NOT NULL THEN
        -- OPCIÓN 1: Notificar TODOS los estados (comentar las siguientes líneas)
        -- En este caso, NO insertamos nada en business_notification_config_order_statuses
        
        -- OPCIÓN 2: Notificar solo estados específicos (descomentar y ajustar)
        -- Asociar estados específicos (ejemplo: solo entregada, completada y cancelada)
        /*
        INSERT INTO business_notification_config_order_statuses (
            business_notification_config_id,
            order_status_id
        )
        SELECT config_id, os.id
        FROM order_statuses os
        WHERE os.code IN ('delivered', 'completed', 'cancelled')
        ON CONFLICT DO NOTHING;
        */
        
        -- OPCIÓN 3: Notificar múltiples estados importantes
        -- Descomentar para notificar: pending, processing, shipped, delivered, completed, cancelled
        INSERT INTO business_notification_config_order_statuses (
            business_notification_config_id,
            order_status_id
        )
        SELECT config_id, os.id
        FROM order_statuses os
        WHERE os.code IN ('pending', 'processing', 'shipped', 'delivered', 'completed', 'cancelled')
        ON CONFLICT DO NOTHING;
        
    END IF;
END $$;

-- ============================================
-- QUERIES ÚTILES PARA VERIFICAR
-- ============================================

-- Ver todas las configuraciones
SELECT 
    bnc.id,
    bnc.business_id,
    bnc.event_type,
    bnc.enabled,
    bnc.channels,
    bnc.description,
    COUNT(DISTINCT bcos.order_status_id) as estados_asociados
FROM business_notification_configs bnc
LEFT JOIN business_notification_config_order_statuses bcos 
    ON bnc.id = bcos.business_notification_config_id
WHERE bnc.business_id = 1
GROUP BY bnc.id, bnc.business_id, bnc.event_type, bnc.enabled, bnc.channels, bnc.description
ORDER BY bnc.event_type;

-- Ver configuraciones con sus estados asociados (detallado)
SELECT 
    bnc.id as config_id,
    bnc.event_type,
    bnc.enabled,
    os.code as status_code,
    os.name as status_name
FROM business_notification_configs bnc
LEFT JOIN business_notification_config_order_statuses bcos 
    ON bnc.id = bcos.business_notification_config_id
LEFT JOIN order_statuses os 
    ON bcos.order_status_id = os.id
WHERE bnc.business_id = 1 
    AND bnc.event_type = 'order.status_changed'
ORDER BY os.code;

-- ============================================
-- NOTAS IMPORTANTES:
-- ============================================
-- 1. Si NO hay registros en business_notification_config_order_statuses para una configuración,
--    significa que se notifican TODOS los cambios de estado.
--
-- 2. Si HAY registros, solo se notifican los estados específicos asociados.
--
-- 3. Los estados disponibles están en la tabla order_statuses:
--    - pending: Pendiente
--    - processing: En Procesamiento  
--    - shipped: Enviada
--    - delivered: Entregada
--    - completed: Completada
--    - cancelled: Cancelada
--    - refunded: Reembolsada
--    - failed: Fallida
--    - on_hold: En Espera
--
-- 4. Para cambiar qué estados notificar:
--    DELETE FROM business_notification_config_order_statuses 
--    WHERE business_notification_config_id = <config_id>;
--    
--    INSERT INTO business_notification_config_order_statuses 
--    SELECT <config_id>, id FROM order_statuses WHERE code IN ('estado1', 'estado2', ...);
