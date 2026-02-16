-- Migration to add "Envios" and "Direcciones de Origen" permissions to Business role (Administrador)
-- Date: 2026-02-16
DO $$
DECLARE read_action_id INT;
business_scope_id INT;
envios_resource_id INT;
admin_role_id INT;
envios_read_permission_id INT;
BEGIN -- 1. Get Action ID for 'Read'
SELECT id INTO read_action_id
FROM actions
WHERE name = 'Read'
LIMIT 1;
IF read_action_id IS NULL THEN RAISE NOTICE 'Action "Read" not found. Please ensure basic actions are seeded.';
END IF;
-- 2. Get Scope ID for 'business'
SELECT id INTO business_scope_id
FROM scopes
WHERE code = 'business'
    OR name = 'Business'
LIMIT 1;
IF business_scope_id IS NULL THEN RAISE NOTICE 'Scope "business" not found. Please ensure basic scopes are seeded.';
END IF;
-- 3. Ensure "Envios" Resource exists
INSERT INTO resources (name, description, created_at, updated_at)
SELECT 'Envios',
    'Gestión de Envíos y Despachos',
    NOW(),
    NOW()
WHERE NOT EXISTS (
        SELECT 1
        FROM resources
        WHERE name = 'Envios'
    )
RETURNING id INTO envios_resource_id;
IF envios_resource_id IS NULL THEN
SELECT id INTO envios_resource_id
FROM resources
WHERE name = 'Envios'
LIMIT 1;
END IF;
-- 4. Ensure "Envios:Read" Permission exists
IF read_action_id IS NOT NULL
AND business_scope_id IS NOT NULL
AND envios_resource_id IS NOT NULL THEN
INSERT INTO permissions (
        name,
        description,
        resource_id,
        action_id,
        scope_id,
        created_at,
        updated_at
    )
SELECT 'Envios:Read',
    'Permiso para leer envíos',
    envios_resource_id,
    read_action_id,
    business_scope_id,
    NOW(),
    NOW()
WHERE NOT EXISTS (
        SELECT 1
        FROM permissions
        WHERE name = 'Envios:Read'
    )
RETURNING id INTO envios_read_permission_id;
IF envios_read_permission_id IS NULL THEN
SELECT id INTO envios_read_permission_id
FROM permissions
WHERE name = 'Envios:Read'
LIMIT 1;
END IF;
-- 5. Grant permission to "Administrador" role
SELECT id INTO admin_role_id
FROM roles
WHERE name = 'Administrador'
LIMIT 1;
IF admin_role_id IS NOT NULL
AND envios_read_permission_id IS NOT NULL THEN
INSERT INTO role_permissions (role_id, permission_id)
SELECT admin_role_id,
    envios_read_permission_id
WHERE NOT EXISTS (
        SELECT 1
        FROM role_permissions
        WHERE role_id = admin_role_id
            AND permission_id = envios_read_permission_id
    );
RAISE NOTICE 'Permission "Envios:Read" granted to "Administrador" role.';
ELSE RAISE NOTICE 'Role "Administrador" or Permission "Envios:Read" not found.';
END IF;
ELSE RAISE NOTICE 'Cannot create permission: missing Action, Scope or Resource dependency.';
END IF;
END $$;