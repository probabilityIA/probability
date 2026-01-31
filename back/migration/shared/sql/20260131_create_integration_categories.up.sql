BEGIN;

-- 1. Crear tabla de categorías
CREATE TABLE integration_categories (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    icon VARCHAR(255),
    color VARCHAR(20),
    display_order INTEGER DEFAULT 0,
    parent_category_id INTEGER REFERENCES integration_categories(id),
    is_active BOOLEAN DEFAULT true,
    is_visible BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_categories_code ON integration_categories(code);
CREATE INDEX idx_categories_active ON integration_categories(is_active);

-- 2. Insertar categorías iniciales
INSERT INTO integration_categories (code, name, description, icon, color, display_order, is_active) VALUES
    ('ecommerce', 'E-commerce', 'Plataformas de venta online', 'shopping-cart', '#3B82F6', 1, true),
    ('invoicing', 'Facturación Electrónica', 'Proveedores de facturación', 'receipt', '#10B981', 2, true),
    ('messaging', 'Mensajería', 'Canales de comunicación', 'message-circle', '#8B5CF6', 3, true),
    ('payment', 'Pagos', 'Pasarelas de pago', 'credit-card', '#F59E0B', 4, false),
    ('shipping', 'Logística', 'Operadores logísticos', 'truck', '#EF4444', 5, false);

-- 3. Agregar columna category_id a integration_types
ALTER TABLE integration_types
    ADD COLUMN category_id INTEGER REFERENCES integration_categories(id);

-- 4. Migrar datos existentes
UPDATE integration_types SET category_id = (SELECT id FROM integration_categories WHERE code = 'ecommerce')
WHERE code IN ('shopify', 'mercadolibre', 'amazon');

UPDATE integration_types SET category_id = (SELECT id FROM integration_categories WHERE code = 'messaging')
WHERE code IN ('whatsapp', 'whatsap');

-- 5. Hacer NOT NULL después de migrar
ALTER TABLE integration_types
    ALTER COLUMN category_id SET NOT NULL;

-- 6. Eliminar columna category antigua
ALTER TABLE integration_types
    DROP COLUMN category;

COMMIT;
