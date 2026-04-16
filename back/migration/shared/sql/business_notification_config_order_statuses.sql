-- ============================================
-- TABLA INTERMEDIA: business_notification_config_order_statuses
-- ============================================
-- Esta tabla relaciona las configuraciones de notificaciones con los estados de orden
-- Permite especificar qué estados de orden deben disparar notificaciones
-- ============================================

-- Crear tabla intermedia
CREATE TABLE IF NOT EXISTS business_notification_config_order_statuses (
    business_notification_config_id INTEGER NOT NULL,
    order_status_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    PRIMARY KEY (business_notification_config_id, order_status_id),
    
    -- Foreign keys
    CONSTRAINT fk_business_notification_config 
        FOREIGN KEY (business_notification_config_id) 
        REFERENCES business_notification_configs(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,
    
    CONSTRAINT fk_order_status 
        FOREIGN KEY (order_status_id) 
        REFERENCES order_statuses(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE
);

-- Crear índices
CREATE INDEX IF NOT EXISTS idx_bncos_notification_config 
    ON business_notification_config_order_statuses(business_notification_config_id);
    
CREATE INDEX IF NOT EXISTS idx_bncos_order_status 
    ON business_notification_config_order_statuses(order_status_id);

-- ============================================
-- NOTAS:
-- ============================================
-- 1. Si una configuración NO tiene estados asociados (tabla vacía),
--    significa que se notifican TODOS los estados para ese event_type.
--
-- 2. Si una configuración tiene estados asociados, solo se notifican esos estados.
--
-- 3. Esta relación solo es relevante para event_type = 'order.status_changed'
--    Para otros tipos de eventos (order.created, etc.), esta tabla no se usa.
