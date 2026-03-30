-- =============================================
-- Migration: Order Statuses v2 - Flujo logístico completo
-- Agrega estados de última milla, normaliza códigos, migra datos
-- =============================================

BEGIN;

-- ═══════════════════════════════════════════
-- 1. INSERTAR NUEVOS ESTADOS
-- ═══════════════════════════════════════════

INSERT INTO order_statuses (code, name, description, category, is_active, priority, color, created_at, updated_at)
VALUES
    ('picking', 'Seleccionando productos', 'Seleccionando productos del inventario', 'active', true, 10, '#3B82F6', NOW(), NOW()),
    ('packing', 'Empacando', 'Empacando el pedido', 'active', true, 20, '#6366F1', NOW(), NOW()),
    ('ready_to_ship', 'Listo para despacho', 'Listo para despacho', 'active', true, 30, '#8B5CF6', NOW(), NOW()),
    ('assigned_to_driver', 'Asignado a piloto', 'Asignado a piloto/conductor', 'active', true, 40, '#A855F7', NOW(), NOW()),
    ('picked_up', 'Recogido', 'Recogido por el piloto', 'active', true, 50, '#D946EF', NOW(), NOW()),
    ('in_transit', 'En camino', 'En camino al destino', 'active', true, 60, '#EC4899', NOW(), NOW()),
    ('out_for_delivery', 'En reparto final', 'En reparto final (última milla)', 'active', true, 70, '#F43F5E', NOW(), NOW()),
    ('delivery_failed', 'Entrega fallida', 'Entrega fallida', 'issue', true, 76, '#EF4444', NOW(), NOW()),
    ('rejected', 'Rechazado', 'Rechazado por el cliente', 'issue', true, 77, '#DC2626', NOW(), NOW()),
    ('return_in_transit', 'Devolución en camino', 'Devolución en camino al almacén', 'return', true, 80, '#F59E0B', NOW(), NOW()),
    ('inventory_issue', 'Novedad de inventario', 'Novedad de inventario (sin stock, producto dañado)', 'issue', true, 15, '#FB923C', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- ═══════════════════════════════════════════
-- 2. NORMALIZAR ESTADOS EXISTENTES
-- ═══════════════════════════════════════════

-- Novelty (ID 10) → delivery_novelty
UPDATE order_statuses
SET code = 'delivery_novelty', name = 'Novedad de entrega', description = 'Novedad de entrega', category = 'issue', priority = 75, color = '#F97316', updated_at = NOW()
WHERE id = 10 AND code = 'Novelty';

-- Refund (ID 11) → returned
UPDATE order_statuses
SET code = 'returned', name = 'Devuelto', description = 'Devuelto al almacén', category = 'return', priority = 85, color = '#EAB308', updated_at = NOW()
WHERE id = 11 AND code = 'Refund';

-- ═══════════════════════════════════════════
-- 3. MIGRAR ÓRDENES DE ESTADOS DEPRECADOS
-- ═══════════════════════════════════════════

-- processing → picking (5,291 órdenes)
UPDATE orders
SET status = 'picking',
    status_id = (SELECT id FROM order_statuses WHERE code = 'picking' LIMIT 1),
    updated_at = NOW()
WHERE status = 'processing' AND deleted_at IS NULL;

-- shipped → in_transit (20 órdenes)
UPDATE orders
SET status = 'in_transit',
    status_id = (SELECT id FROM order_statuses WHERE code = 'in_transit' LIMIT 1),
    updated_at = NOW()
WHERE status = 'shipped' AND deleted_at IS NULL;

-- ═══════════════════════════════════════════
-- 4. ACTUALIZAR SHOPIFY MAPPINGS
-- ═══════════════════════════════════════════

-- Mappings que apuntaban a processing (ID 2) → ahora apuntan a picking
UPDATE order_status_mappings
SET order_status_id = (SELECT id FROM order_statuses WHERE code = 'picking' LIMIT 1),
    updated_at = NOW()
WHERE order_status_id = (SELECT id FROM order_statuses WHERE code = 'processing' LIMIT 1)
  AND deleted_at IS NULL;

-- Mappings que apuntaban a shipped (ID 3) → ahora apuntan a in_transit
UPDATE order_status_mappings
SET order_status_id = (SELECT id FROM order_statuses WHERE code = 'in_transit' LIMIT 1),
    updated_at = NOW()
WHERE order_status_id = (SELECT id FROM order_statuses WHERE code = 'shipped' LIMIT 1)
  AND deleted_at IS NULL;

-- ═══════════════════════════════════════════
-- 5. DESACTIVAR ESTADOS DEPRECADOS
-- ═══════════════════════════════════════════

UPDATE order_statuses
SET is_active = false, updated_at = NOW()
WHERE code IN ('processing', 'shipped');

COMMIT;
