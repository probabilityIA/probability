-- ============================================
-- Mapeos de Estados Shopify → Probability
-- ============================================
-- Este archivo contiene los INSERT statements para crear los mapeos iniciales
-- entre los estados de Shopify y los estados de Probability.
--
-- Uso: Ejecutar este script en PostgreSQL para poblar la tabla order_status_mappings
-- con los mapeos básicos de Shopify.
--
-- IMPORTANTE: Este script asume que:
--   1. La tabla order_statuses ya existe y está poblada (ejecutar primero: order_statuses.sql)
--   2. La tabla integration_types tiene el registro de Shopify con ID = 1
--
-- Tipos de estados de Shopify:
--   - Order Status: any, open, closed, cancelled
--   - Financial Status: any, authorized, pending, paid, partially_paid, refunded, voided, partially_refunded, unpaid
--   - Fulfillment Status: any, shipped, partial, unfulfilled, unshipped
--   - Shipment Status (dentro de fulfillments): confirmed, success, delivered, failure, cancelled (PRIORIDAD 2)
--
-- Estados de Probability (order_statuses):
--   - pending, processing, shipped, delivered, completed, cancelled, refunded, failed, on_hold
--
-- ID de IntegrationType para Shopify: 1
-- ============================================

-- Mapeos de Order Status (Estado de Orden)
INSERT INTO order_status_mappings (integration_type_id, original_status, order_status_id, is_active, priority, description, created_at, updated_at)
VALUES 
  (1, 'open', (SELECT id FROM order_statuses WHERE code = 'processing' LIMIT 1), true, 1, 'Orden abierta en Shopify → En procesamiento en Probability', NOW(), NOW()),
  (1, 'closed', (SELECT id FROM order_statuses WHERE code = 'completed' LIMIT 1), true, 1, 'Orden cerrada en Shopify → Completada en Probability', NOW(), NOW()),
  (1, 'cancelled', (SELECT id FROM order_statuses WHERE code = 'cancelled' LIMIT 1), true, 1, 'Orden cancelada en Shopify → Cancelada en Probability', NOW(), NOW()),
  (1, 'any', (SELECT id FROM order_statuses WHERE code = 'pending' LIMIT 1), true, 0, 'Estado "any" en Shopify → Pendiente en Probability (por defecto)', NOW(), NOW())
ON CONFLICT (integration_type_id, original_status) DO NOTHING;

-- Mapeos de Financial Status (Estado Financiero)
INSERT INTO order_status_mappings (integration_type_id, original_status, order_status_id, is_active, priority, description, created_at, updated_at)
VALUES 
  (1, 'pending', (SELECT id FROM order_statuses WHERE code = 'pending' LIMIT 1), true, 1, 'Pago pendiente en Shopify → Pendiente en Probability', NOW(), NOW()),
  (1, 'authorized', (SELECT id FROM order_statuses WHERE code = 'processing' LIMIT 1), true, 1, 'Pago autorizado en Shopify → En procesamiento en Probability', NOW(), NOW()),
  (1, 'paid', (SELECT id FROM order_statuses WHERE code = 'processing' LIMIT 1), true, 1, 'Pago completado en Shopify → En procesamiento en Probability', NOW(), NOW()),
  (1, 'partially_paid', (SELECT id FROM order_statuses WHERE code = 'processing' LIMIT 1), true, 1, 'Pago parcial en Shopify → En procesamiento en Probability', NOW(), NOW()),
  (1, 'refunded', (SELECT id FROM order_statuses WHERE code = 'refunded' LIMIT 1), true, 1, 'Reembolsado en Shopify → Reembolsado en Probability', NOW(), NOW()),
  (1, 'partially_refunded', (SELECT id FROM order_statuses WHERE code = 'refunded' LIMIT 1), true, 1, 'Reembolso parcial en Shopify → Reembolsado en Probability', NOW(), NOW()),
  (1, 'voided', (SELECT id FROM order_statuses WHERE code = 'cancelled' LIMIT 1), true, 1, 'Pago anulado en Shopify → Cancelado en Probability', NOW(), NOW()),
  (1, 'unpaid', (SELECT id FROM order_statuses WHERE code = 'pending' LIMIT 1), true, 1, 'No pagado en Shopify → Pendiente en Probability', NOW(), NOW()),
  (1, 'any', (SELECT id FROM order_statuses WHERE code = 'pending' LIMIT 1), true, 0, 'Estado financiero "any" en Shopify → Pendiente en Probability (por defecto)', NOW(), NOW())
ON CONFLICT (integration_type_id, original_status) DO NOTHING;

-- Mapeos de Fulfillment Status (Estado de Cumplimiento/Envío)
INSERT INTO order_status_mappings (integration_type_id, original_status, order_status_id, is_active, priority, description, created_at, updated_at)
VALUES 
  (1, 'shipped', (SELECT id FROM order_statuses WHERE code = 'shipped' LIMIT 1), true, 1, 'Enviado en Shopify → Enviado en Probability', NOW(), NOW()),
  (1, 'partial', (SELECT id FROM order_statuses WHERE code = 'processing' LIMIT 1), true, 1, 'Envío parcial en Shopify → En procesamiento en Probability', NOW(), NOW()),
  (1, 'unfulfilled', (SELECT id FROM order_statuses WHERE code = 'pending' LIMIT 1), true, 1, 'Sin cumplir en Shopify → Pendiente en Probability', NOW(), NOW()),
  (1, 'unshipped', (SELECT id FROM order_statuses WHERE code = 'on_hold' LIMIT 1), true, 1, 'No enviado en Shopify → En espera en Probability', NOW(), NOW()),
  (1, 'any', (SELECT id FROM order_statuses WHERE code = 'pending' LIMIT 1), true, 0, 'Estado de envío "any" en Shopify → Pendiente en Probability (por defecto)', NOW(), NOW())
ON CONFLICT (integration_type_id, original_status) DO NOTHING;

-- Mapeos de Shipment Status (Estado de Envío dentro de fulfillments) - PRIORIDAD 2
-- Estos tienen mayor prioridad que fulfillment_status y financial_status
INSERT INTO order_status_mappings (integration_type_id, original_status, order_status_id, is_active, priority, description, created_at, updated_at)
VALUES 
  (1, 'confirmed', (SELECT id FROM order_statuses WHERE code = 'processing' LIMIT 1), true, 2, 'Confirmado en Shopify → En Procesamiento en Probability', NOW(), NOW()),
  (1, 'success', (SELECT id FROM order_statuses WHERE code = 'shipped' LIMIT 1), true, 2, 'Exitoso (en tránsito) en Shopify → Enviada en Probability', NOW(), NOW()),
  (1, 'delivered', (SELECT id FROM order_statuses WHERE code = 'delivered' LIMIT 1), true, 2, 'Entregado en Shopify → Entregada en Probability', NOW(), NOW()),
  (1, 'failure', (SELECT id FROM order_statuses WHERE code = 'failed' LIMIT 1), true, 2, 'Fallido en Shopify → Fallida en Probability', NOW(), NOW()),
  (1, 'cancelled', (SELECT id FROM order_statuses WHERE code = 'cancelled' LIMIT 1), true, 2, 'Cancelado (shipment) en Shopify → Cancelada en Probability', NOW(), NOW())
ON CONFLICT (integration_type_id, original_status) DO UPDATE SET 
  order_status_id = EXCLUDED.order_status_id,
  priority = EXCLUDED.priority,
  description = EXCLUDED.description,
  updated_at = NOW();

-- ============================================
-- Notas:
-- ============================================
-- 1. Estos mapeos son sugerencias iniciales y pueden ser modificados desde el CRUD
--    disponible en /api/v1/order-status-mappings
--
-- 2. Los estados con priority=0 son valores por defecto cuando no hay un mapeo específico
--
-- 3. Se puede consultar todos los mapeos de Shopify con:
--    GET /api/v1/order-status-mappings?integration_type_id=1
--
-- 4. Para gestionar estos mapeos desde el frontend, usar los endpoints del CRUD:
--    - POST   /api/v1/order-status-mappings (crear)
--    - GET    /api/v1/order-status-mappings/:id (obtener)
--    - PUT    /api/v1/order-status-mappings/:id (actualizar)
--    - DELETE /api/v1/order-status-mappings/:id (eliminar)
--    - PATCH  /api/v1/order-status-mappings/:id/toggle (activar/desactivar)
-- ============================================
