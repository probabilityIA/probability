-- Query SQL para actualizar la columna status_id de las órdenes existentes
-- basándose en su original_status e integration_type usando los mapeos de order_status_mappings
-- 
-- Esta query actualiza todas las órdenes que tienen original_status y encuentra
-- el mapeo correspondiente en order_status_mappings para asignar el order_status_id correcto

UPDATE orders o
SET status_id = (
    SELECT osm.order_status_id
    FROM order_status_mappings osm
    WHERE 
        osm.original_status = o.original_status
        AND osm.integration_type_id = CASE 
            WHEN o.integration_type = 'shopify' THEN 1
            WHEN o.integration_type IN ('whatsapp', 'whatsap', 'whastap') THEN 2
            WHEN o.integration_type IN ('mercado_libre', 'mercadolibre') THEN 3
            WHEN o.integration_type IN ('woocommerce', 'woocormerce') THEN 4
            ELSE 0
        END
        AND osm.is_active = true
        AND osm.integration_type_id > 0
    ORDER BY osm.priority DESC, osm.created_at DESC
    LIMIT 1
)
WHERE 
    o.original_status IS NOT NULL 
    AND o.original_status != ''
    AND EXISTS (
        SELECT 1
        FROM order_status_mappings osm
        WHERE 
            osm.original_status = o.original_status
            AND osm.integration_type_id = CASE 
                WHEN o.integration_type = 'shopify' THEN 1
                WHEN o.integration_type IN ('whatsapp', 'whatsap', 'whastap') THEN 2
                WHEN o.integration_type IN ('mercado_libre', 'mercadolibre') THEN 3
                WHEN o.integration_type IN ('woocommerce', 'woocormerce') THEN 4
                ELSE 0
            END
            AND osm.is_active = true
            AND osm.integration_type_id > 0
    );

-- Verificar cuántas órdenes fueron actualizadas (ejecutar después para verificar)
-- SELECT COUNT(*) as total_actualizadas FROM orders WHERE status_id IS NOT NULL;

-- Ver órdenes que NO tienen status_id después de la actualización (para debugging)
-- SELECT id, integration_type, original_status, status_id 
-- FROM orders 
-- WHERE original_status IS NOT NULL AND original_status != '' AND status_id IS NULL
-- LIMIT 20;
