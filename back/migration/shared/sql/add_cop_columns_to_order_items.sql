-- ============================================
-- Agregar columnas presentment (presentment_money) a order_items
-- ============================================
-- Este script agrega las columnas para almacenar precios en moneda local (presentment_money)
-- manteniendo las columnas existentes en USD (shop_money)
-- ============================================

-- Agregar columnas si no existen
ALTER TABLE order_items 
ADD COLUMN IF NOT EXISTS unit_price_presentment DECIMAL(12,2) NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS total_price_presentment DECIMAL(12,2) NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS discount_presentment DECIMAL(12,2) NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS tax_presentment DECIMAL(12,2) NOT NULL DEFAULT 0;

-- ============================================
-- Notas:
-- ============================================
-- 1. Las columnas *_presentment almacenan valores de presentment_money de Shopify
--    mientras que las columnas existentes (sin _presentment) almacenan shop_money (USD)
--
-- 2. Los valores por defecto son 0 para mantener compatibilidad con items
--    que no tienen presentment_money disponible
--
-- 3. Si los valores presentment son 0, el sistema debe usar los valores USD
--    como fallback para mantener compatibilidad
-- ============================================








