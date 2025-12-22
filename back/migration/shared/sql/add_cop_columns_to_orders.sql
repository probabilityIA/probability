-- ============================================
-- Agregar columnas presentment (presentment_money) a orders
-- ============================================
-- Este script agrega las columnas para almacenar precios en moneda local (presentment_money)
-- manteniendo las columnas existentes en USD (shop_money)
-- ============================================

-- Agregar columnas si no existen
ALTER TABLE orders 
ADD COLUMN IF NOT EXISTS subtotal_presentment DECIMAL(12,2) NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS tax_presentment DECIMAL(12,2) NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS discount_presentment DECIMAL(12,2) NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS shipping_cost_presentment DECIMAL(12,2) NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS total_amount_presentment DECIMAL(12,2) NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS currency_presentment VARCHAR(10);

-- ============================================
-- Notas:
-- ============================================
-- 1. Las columnas *_presentment almacenan valores de presentment_money de Shopify
--    mientras que las columnas existentes (sin _presentment) almacenan shop_money (USD)
--
-- 2. currency_presentment almacena la moneda presentment (moneda local, puede ser COP, EUR, etc.)
--    mientras que currency almacena la moneda shop (generalmente "USD")
--
-- 3. Los valores por defecto son 0 para mantener compatibilidad con órdenes
--    que no tienen presentment_money disponible
--
-- 4. Si total_amount_presentment es 0 o NULL, el sistema debe usar total_amount (USD)
--    como fallback para mantener compatibilidad
-- ============================================





