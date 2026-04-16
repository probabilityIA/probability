-- ============================================
-- MIGRACIÓN: Agregar tabla intermedia para estados de orden
-- ============================================
-- Este script crea la nueva tabla intermedia para relacionar
-- configuraciones de notificaciones con estados de orden.
-- 
-- IMPORTANTE: NO elimina ninguna tabla existente, solo agrega la nueva.
-- ============================================

-- 1. Crear la tabla intermedia (si no existe)
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

-- 2. Crear índices
CREATE INDEX IF NOT EXISTS idx_bncos_notification_config 
    ON business_notification_config_order_statuses(business_notification_config_id);
    
CREATE INDEX IF NOT EXISTS idx_bncos_order_status 
    ON business_notification_config_order_statuses(order_status_id);

-- ============================================
-- OPCIONAL: Migrar datos del campo filters (JSON) a la nueva tabla
-- ============================================
-- Si tenías estados en el campo filters como:
-- {"statuses": ["pending", "processing", ...]}
-- puedes migrarlos con este script:

/*
DO $$
DECLARE
    config_record RECORD;
    status_code TEXT;
    status_id INTEGER;
BEGIN
    -- Iterar sobre todas las configuraciones con event_type = 'order.status_changed'
    FOR config_record IN 
        SELECT id, filters 
        FROM business_notification_configs 
        WHERE event_type = 'order.status_changed'
    LOOP
        -- Extraer estados del JSON filters
        -- Nota: Esto es un ejemplo simplificado, puede requerir ajustes según tu estructura JSON exacta
        FOR status_code IN 
            SELECT value::text
            FROM jsonb_array_elements_text(
                COALESCE(config_record.filters->'statuses', '[]'::jsonb)
            )
        LOOP
            -- Obtener el ID del estado por su código
            SELECT id INTO status_id 
            FROM order_statuses 
            WHERE code = TRIM(BOTH '"' FROM status_code);
            
            -- Insertar en la tabla intermedia si se encontró el estado
            IF status_id IS NOT NULL THEN
                INSERT INTO business_notification_config_order_statuses 
                    (business_notification_config_id, order_status_id)
                VALUES (config_record.id, status_id)
                ON CONFLICT (business_notification_config_id, order_status_id) DO NOTHING;
            END IF;
        END LOOP;
    END LOOP;
END $$;
*/

-- ============================================
-- NOTAS:
-- ============================================
-- 1. La tabla business_notification_configs NO se elimina ni modifica.
--    El campo filters puede seguir existiendo para otros tipos de filtros.
--
-- 2. Si quieres limpiar el campo filters de estados (opcional):
--    UPDATE business_notification_configs 
--    SET filters = filters - 'statuses'
--    WHERE event_type = 'order.status_changed';
--
-- 3. GORM AutoMigrate debería crear esta tabla automáticamente si ejecutas
--    la migración con el modelo actualizado, pero este script SQL te da
--    más control sobre la creación.
