-- Agregar campos de contacto para carrier (APIs de transportadoras) al modelo Warehouse
-- Estos campos permiten usar bodegas como fuente de datos de remitente para guías de envío

ALTER TABLE warehouses ADD COLUMN IF NOT EXISTS company VARCHAR(100) DEFAULT '';
ALTER TABLE warehouses ADD COLUMN IF NOT EXISTS first_name VARCHAR(100) DEFAULT '';
ALTER TABLE warehouses ADD COLUMN IF NOT EXISTS last_name VARCHAR(100) DEFAULT '';
ALTER TABLE warehouses ADD COLUMN IF NOT EXISTS email VARCHAR(100) DEFAULT '';
ALTER TABLE warehouses ADD COLUMN IF NOT EXISTS suburb VARCHAR(100) DEFAULT '';
ALTER TABLE warehouses ADD COLUMN IF NOT EXISTS city_dane_code VARCHAR(10) DEFAULT '';
ALTER TABLE warehouses ADD COLUMN IF NOT EXISTS postal_code VARCHAR(20) DEFAULT '';
ALTER TABLE warehouses ADD COLUMN IF NOT EXISTS street VARCHAR(255) DEFAULT '';
