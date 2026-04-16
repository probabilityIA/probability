-- ============================================
-- Agregar colores a payment_statuses
-- ============================================
-- Este script actualiza los colores por defecto para cada estado de pago
-- ============================================

-- Actualizar colores por defecto para cada estado de pago
UPDATE payment_statuses 
SET color = CASE 
    WHEN code = 'pending' THEN '#F59E0B'           -- Amarillo (amber-500) - Pendiente
    WHEN code = 'authorized' THEN '#3B82F6'       -- Azul (blue-500) - Autorizado
    WHEN code = 'paid' THEN '#10B981'             -- Verde (green-500) - Pagado
    WHEN code = 'partially_paid' THEN '#F97316'   -- Naranja (orange-500) - Parcialmente pagado
    WHEN code = 'refunded' THEN '#EF4444'         -- Rojo (red-500) - Reembolsado
    WHEN code = 'partially_refunded' THEN '#F59E0B' -- Amarillo (amber-500) - Parcialmente reembolsado
    WHEN code = 'voided' THEN '#6B7280'           -- Gris (gray-500) - Anulado
    WHEN code = 'failed' THEN '#DC2626'           -- Rojo oscuro (red-600) - Fallido
    WHEN code = 'unpaid' THEN '#EF4444'           -- Rojo (red-500) - No pagado
    ELSE '#6B7280'                                -- Gris por defecto
END
WHERE color IS NULL OR color = '';
