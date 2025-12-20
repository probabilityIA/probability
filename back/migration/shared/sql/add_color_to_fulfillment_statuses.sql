-- ============================================
-- Agregar colores a fulfillment_statuses
-- ============================================
-- Este script actualiza los colores por defecto para cada estado de fulfillment
-- ============================================

-- Actualizar colores por defecto para cada estado de fulfillment
UPDATE fulfillment_statuses 
SET color = CASE 
    WHEN code = 'unfulfilled' THEN '#EF4444'      -- Rojo (red-500) - No cumplida
    WHEN code = 'partial' THEN '#F59E0B'          -- Amarillo (amber-500) - Parcial
    WHEN code = 'fulfilled' THEN '#10B981'        -- Verde (green-500) - Cumplida
    WHEN code = 'shipped' THEN '#3B82F6'          -- Azul (blue-500) - Enviada
    WHEN code = 'in_transit' THEN '#8B5CF6'       -- Morado (purple-500) - En tránsito
    WHEN code = 'delivered' THEN '#10B981'        -- Verde (green-500) - Entregada
    WHEN code = 'cancelled' THEN '#DC2626'        -- Rojo oscuro (red-600) - Cancelada
    WHEN code = 'failed' THEN '#DC2626'           -- Rojo oscuro (red-600) - Fallido
    ELSE '#6B7280'                                -- Gris por defecto
END
WHERE color IS NULL OR color = '';
