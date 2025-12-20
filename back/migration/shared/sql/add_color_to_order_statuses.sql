-- ============================================
-- Agregar columna color a order_statuses
-- ============================================
-- Este script agrega la columna color a la tabla order_statuses
-- y actualiza los colores por defecto para cada estado
-- ============================================

-- Agregar columna color si no existe
ALTER TABLE order_statuses 
ADD COLUMN IF NOT EXISTS color VARCHAR(32);

-- Actualizar colores por defecto para cada estado
UPDATE order_statuses 
SET color = CASE 
    WHEN code = 'pending' THEN '#F59E0B'      -- Amarillo (amber-500)
    WHEN code = 'processing' THEN '#3B82F6'  -- Azul (blue-500)
    WHEN code = 'shipped' THEN '#8B5CF6'     -- Morado (purple-500)
    WHEN code = 'delivered' THEN '#10B981'   -- Verde (green-500)
    WHEN code = 'completed' THEN '#10B981'   -- Verde (green-500)
    WHEN code = 'cancelled' THEN '#EF4444'   -- Rojo (red-500)
    WHEN code = 'refunded' THEN '#F97316'    -- Naranja (orange-500)
    WHEN code = 'failed' THEN '#DC2626'      -- Rojo oscuro (red-600)
    WHEN code = 'on_hold' THEN '#6B7280'     -- Gris (gray-500)
    ELSE '#6B7280'                           -- Gris por defecto
END
WHERE color IS NULL OR color = '';

-- Crear índice si es necesario (opcional)
CREATE INDEX IF NOT EXISTS idx_order_statuses_color ON order_statuses(color);
