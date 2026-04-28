-- ============================================
-- RECURSOS POR SUBMODULO DE INVENTARIO
-- ============================================
-- Permite gatear cada vista del modulo de inventario
-- (Stock, Movimientos, Trazabilidad, Kardex, Operaciones,
--  Slotting, Auditoria, LPN, Scan, Sync Logs) de forma
-- independiente por negocio.
-- ============================================

SELECT setval(pg_get_serial_sequence('resource',   'id'), GREATEST(COALESCE((SELECT MAX(id) FROM resource),   1), 1));
SELECT setval(pg_get_serial_sequence('permission', 'id'), GREATEST(COALESCE((SELECT MAX(id) FROM permission), 1), 1));

INSERT INTO resource (name, description, created_at, updated_at) VALUES
    ('Inventario-Stock',         'Vista de stock por bodega/SKU',                NOW(), NOW()),
    ('Inventario-Movimientos',   'Historial y registro de movimientos',          NOW(), NOW()),
    ('Inventario-Trazabilidad',  'Trazabilidad por lote/serial',                 NOW(), NOW()),
    ('Inventario-Kardex',        'Kardex contable por SKU',                      NOW(), NOW()),
    ('Inventario-Operaciones',   'Recepciones, despachos y conteos',             NOW(), NOW()),
    ('Inventario-Slotting',      'Analitica de slotting ABC',                    NOW(), NOW()),
    ('Inventario-Auditoria',     'Auditoria de inventario',                      NOW(), NOW()),
    ('Inventario-LPN',           'Gestion de License Plate Numbers',             NOW(), NOW()),
    ('Inventario-Scan',          'App movil de captura por escaneo',             NOW(), NOW()),
    ('Inventario-Sync-Logs',     'Logs de sincronizacion de inventario',         NOW(), NOW())
ON CONFLICT (name) DO NOTHING;

DO $$
DECLARE
    res_name TEXT;
    res_id   BIGINT;
    scope_business_id BIGINT := 2;
    sub_resources TEXT[] := ARRAY[
        'Inventario-Stock','Inventario-Movimientos','Inventario-Trazabilidad',
        'Inventario-Kardex','Inventario-Operaciones','Inventario-Slotting',
        'Inventario-Auditoria','Inventario-LPN','Inventario-Scan','Inventario-Sync-Logs'
    ];
BEGIN
    FOREACH res_name IN ARRAY sub_resources LOOP
        SELECT id INTO res_id FROM resource WHERE name = res_name LIMIT 1;
        IF res_id IS NULL THEN
            RAISE EXCEPTION 'No se encontro el recurso %', res_name;
        END IF;

        INSERT INTO permission (name, description, resource_id, action_id, scope_id, created_at, updated_at) VALUES
            ('Create ' || res_name, 'Crear en ' || res_name, res_id, 1, scope_business_id, NOW(), NOW()),
            ('Read '   || res_name, 'Ver '       || res_name, res_id, 2, scope_business_id, NOW(), NOW()),
            ('Update ' || res_name, 'Editar '    || res_name, res_id, 3, scope_business_id, NOW(), NOW()),
            ('Delete ' || res_name, 'Eliminar '  || res_name, res_id, 4, scope_business_id, NOW(), NOW())
        ON CONFLICT (name) DO NOTHING;
    END LOOP;

    RAISE NOTICE '10 recursos Inventario-* creados con permisos CRUD';
END $$;
