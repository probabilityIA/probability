BEGIN;

-- 1. Restaurar columna category antigua
ALTER TABLE integration_types
    ADD COLUMN category VARCHAR(50);

-- 2. Migrar datos de vuelta
UPDATE integration_types SET category = 'ecommerce'
WHERE category_id = (SELECT id FROM integration_categories WHERE code = 'ecommerce');

UPDATE integration_types SET category = 'messaging'
WHERE category_id = (SELECT id FROM integration_categories WHERE code = 'messaging');

-- 3. Eliminar columna category_id
ALTER TABLE integration_types
    DROP COLUMN category_id;

-- 4. Eliminar tabla de categorías
DROP TABLE integration_categories CASCADE;

COMMIT;
