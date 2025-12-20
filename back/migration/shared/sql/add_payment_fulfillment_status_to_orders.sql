-- ============================================
-- Agregar columnas payment_status_id y fulfillment_status_id a orders
-- ============================================
-- Este script agrega las columnas payment_status_id y fulfillment_status_id
-- a la tabla orders y migra los datos existentes.
-- ============================================

-- Agregar columnas si no existen
ALTER TABLE orders 
ADD COLUMN IF NOT EXISTS payment_status_id INTEGER,
ADD COLUMN IF NOT EXISTS fulfillment_status_id INTEGER;

-- Crear índices
CREATE INDEX IF NOT EXISTS idx_orders_payment_status_id ON orders(payment_status_id);
CREATE INDEX IF NOT EXISTS idx_orders_fulfillment_status_id ON orders(fulfillment_status_id);

-- Agregar foreign keys
DO $$
BEGIN
    -- Agregar foreign key para payment_status_id
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = 'fk_orders_payment_status'
    ) THEN
        ALTER TABLE orders
        ADD CONSTRAINT fk_orders_payment_status
        FOREIGN KEY (payment_status_id) 
        REFERENCES payment_statuses(id) 
        ON UPDATE CASCADE 
        ON DELETE SET NULL;
    END IF;

    -- Agregar foreign key para fulfillment_status_id
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = 'fk_orders_fulfillment_status'
    ) THEN
        ALTER TABLE orders
        ADD CONSTRAINT fk_orders_fulfillment_status
        FOREIGN KEY (fulfillment_status_id) 
        REFERENCES fulfillment_statuses(id) 
        ON UPDATE CASCADE 
        ON DELETE SET NULL;
    END IF;
END $$;

-- Migrar datos existentes: mapear is_paid=true a payment_status_id de "paid" si existe
UPDATE orders
SET payment_status_id = (SELECT id FROM payment_statuses WHERE code = 'paid' LIMIT 1)
WHERE is_paid = true 
  AND payment_status_id IS NULL
  AND EXISTS (SELECT 1 FROM payment_statuses WHERE code = 'paid');

-- Para órdenes no pagadas, establecer como 'unpaid' si existe
UPDATE orders
SET payment_status_id = (SELECT id FROM payment_statuses WHERE code = 'unpaid' LIMIT 1)
WHERE is_paid = false 
  AND payment_status_id IS NULL
  AND EXISTS (SELECT 1 FROM payment_statuses WHERE code = 'unpaid');

-- Para fulfillment, establecer como 'unfulfilled' por defecto si no hay tracking_number
UPDATE orders
SET fulfillment_status_id = (SELECT id FROM fulfillment_statuses WHERE code = 'unfulfilled' LIMIT 1)
WHERE fulfillment_status_id IS NULL
  AND (tracking_number IS NULL OR tracking_number = '')
  AND EXISTS (SELECT 1 FROM fulfillment_statuses WHERE code = 'unfulfilled');

-- Si hay tracking_number, establecer como 'shipped' o 'in_transit'
UPDATE orders
SET fulfillment_status_id = (SELECT id FROM fulfillment_statuses WHERE code = 'shipped' LIMIT 1)
WHERE fulfillment_status_id IS NULL
  AND tracking_number IS NOT NULL
  AND tracking_number != ''
  AND delivered_at IS NULL
  AND EXISTS (SELECT 1 FROM fulfillment_statuses WHERE code = 'shipped');

-- Si hay delivered_at, establecer como 'delivered'
UPDATE orders
SET fulfillment_status_id = (SELECT id FROM fulfillment_statuses WHERE code = 'delivered' LIMIT 1)
WHERE fulfillment_status_id IS NULL
  AND delivered_at IS NOT NULL
  AND EXISTS (SELECT 1 FROM fulfillment_statuses WHERE code = 'delivered');

-- ============================================
-- Notas:
-- ============================================
-- 1. Los valores NULL son permitidos para mantener flexibilidad
--    en órdenes que aún no tienen estados asignados
--
-- 2. La migración establece valores por defecto basados en is_paid
--    y tracking_number/delivered_at, pero estos pueden ser refinados
--    con datos más precisos desde las integraciones
-- ============================================
