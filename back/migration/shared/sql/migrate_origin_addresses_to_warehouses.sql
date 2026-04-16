-- ============================================================
-- Migración: OriginAddress → Warehouse para Mystic Rose (business_id=36)
-- Fecha: 2026-03-01
-- Descripción: Copia los datos de origin_addresses al módulo warehouses
-- ============================================================

-- Paso 1: Verificar datos de origen (solo lectura, ejecutar antes)
-- SELECT id, alias, company, first_name, last_name, email, phone,
--        street, suburb, city_dane_code, city, state, postal_code, is_default
-- FROM origin_address
-- WHERE business_id = 36 AND deleted_at IS NULL;

-- Paso 2: Insertar como warehouses (idempotente - evita duplicados por nombre+business)
-- Nota: la tabla es "origin_address" (singular)
INSERT INTO warehouses (
    business_id, name, code, address, city, state, country, zip_code,
    phone, contact_name, contact_email,
    is_active, is_default, is_fulfillment,
    company, first_name, last_name, email,
    suburb, city_dane_code, postal_code, street,
    created_at, updated_at
)
SELECT
    oa.business_id,
    oa.alias AS name,
    'BOD-' || LPAD(ROW_NUMBER() OVER (ORDER BY oa.id)::text, 3, '0') AS code,
    oa.street AS address,
    oa.city,
    oa.state,
    'CO' AS country,
    COALESCE(oa.postal_code, '') AS zip_code,
    oa.phone,
    CONCAT(oa.first_name, ' ', oa.last_name) AS contact_name,
    oa.email AS contact_email,
    true AS is_active,
    oa.is_default,
    false AS is_fulfillment,
    oa.company,
    oa.first_name,
    oa.last_name,
    oa.email,
    oa.suburb,
    oa.city_dane_code,
    COALESCE(oa.postal_code, '') AS postal_code,
    oa.street,
    NOW() AS created_at,
    NOW() AS updated_at
FROM origin_address oa
WHERE oa.business_id = 36
  AND oa.deleted_at IS NULL
  AND NOT EXISTS (
      SELECT 1 FROM warehouses w
      WHERE w.business_id = oa.business_id
        AND w.name = oa.alias
        AND w.deleted_at IS NULL
  );

-- Paso 3: Vincular órdenes existentes sin bodega a la bodega default
UPDATE probability_orders po
SET warehouse_id = w.id,
    warehouse_name = w.name
FROM warehouses w
WHERE po.business_id = 36
  AND w.business_id = 36
  AND w.is_default = true
  AND w.deleted_at IS NULL
  AND po.warehouse_id IS NULL;

-- Paso 4: Verificar resultados
-- SELECT id, name, code, city, is_default, company, street, city_dane_code
-- FROM warehouses WHERE business_id = 36 AND deleted_at IS NULL;
--
-- SELECT COUNT(*) AS orders_linked
-- FROM probability_orders
-- WHERE business_id = 36 AND warehouse_id IS NOT NULL;
