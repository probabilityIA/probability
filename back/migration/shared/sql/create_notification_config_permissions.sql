-- ============================================
-- CREAR RECURSO Y PERMISOS PARA NOTIFICACIONES
-- ============================================
-- Recurso: Notificaciones (CRUD con scope business)
-- Los permisos de configuracion (Read/Update) son para usuarios de negocio
-- Los canales y tipos de evento son solo para super admin (no necesitan permiso, se controla en UI)
-- ============================================

-- 1. Crear recurso
INSERT INTO resource (name, description, created_at, updated_at)
VALUES ('Notificaciones', 'Configuracion de notificaciones por integracion', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- 2. Obtener el ID del recurso recien creado
DO $$
DECLARE
    resource_notif_id BIGINT;
    scope_business_id BIGINT := 2; -- scope "Business"
BEGIN
    SELECT id INTO resource_notif_id FROM resource WHERE name = 'Notificaciones' LIMIT 1;

    IF resource_notif_id IS NULL THEN
        RAISE EXCEPTION 'No se encontro el recurso Notificaciones';
    END IF;

    -- 3. Crear permisos CRUD
    -- Create
    INSERT INTO permission (name, description, resource_id, action_id, scope_id, created_at, updated_at)
    VALUES ('Create Notificaciones', 'Crear configuraciones de notificacion', resource_notif_id, 1, scope_business_id, NOW(), NOW())
    ON CONFLICT DO NOTHING;

    -- Read
    INSERT INTO permission (name, description, resource_id, action_id, scope_id, created_at, updated_at)
    VALUES ('Read Notificaciones', 'Ver configuraciones de notificacion y auditoria de mensajes', resource_notif_id, 2, scope_business_id, NOW(), NOW())
    ON CONFLICT DO NOTHING;

    -- Update
    INSERT INTO permission (name, description, resource_id, action_id, scope_id, created_at, updated_at)
    VALUES ('Update Notificaciones', 'Editar reglas de notificacion', resource_notif_id, 3, scope_business_id, NOW(), NOW())
    ON CONFLICT DO NOTHING;

    -- Delete
    INSERT INTO permission (name, description, resource_id, action_id, scope_id, created_at, updated_at)
    VALUES ('Delete Notificaciones', 'Eliminar configuraciones de notificacion', resource_notif_id, 4, scope_business_id, NOW(), NOW())
    ON CONFLICT DO NOTHING;

    RAISE NOTICE 'Recurso Notificaciones (id=%) creado con 4 permisos CRUD', resource_notif_id;
END $$;

-- ============================================
-- VERIFICAR
-- ============================================
-- SELECT r.name as recurso, p.name as permiso, a.name as accion, s.name as scope
-- FROM permission p
-- JOIN resource r ON p.resource_id = r.id
-- JOIN action a ON p.action_id = a.id
-- JOIN scope s ON p.scope_id = s.id
-- WHERE r.name = 'Notificaciones'
-- ORDER BY p.id;
