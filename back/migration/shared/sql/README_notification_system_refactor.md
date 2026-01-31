# Migración: Refactorización del Sistema de Notificaciones

**Fecha:** 2026-01-30
**Versión:** 1.0
**Autor:** Claude Code

## Resumen

Esta migración refactoriza completamente el sistema de configuración de notificaciones para soportar:

1. **Tipos de notificación** (WhatsApp, SSE, Email, SMS) en tabla propia con CRUD
2. **Eventos de notificación** por tipo (order.created, order.canceled, invoice.created, etc.)
3. **Integración origen** única por configuración (la que genera el evento)
4. **Checklist de estados de orden** para configurar qué estados disparan notificaciones

## Arquitectura Anterior vs Nueva

### Antes (Arquitectura Plana)
```
BusinessNotificationConfig
├── business_id (FK)
├── event_type (string: "order.created")
├── channels (JSON: ["sse", "whatsapp"])
├── filters (JSON)
└── order_statuses (M2M)
```

**Problemas:**
- `channels` como JSON array (difícil de mantener)
- `event_type` como string libre (sin validación)
- NO hay relación con la integración origen
- Configuraciones duplicadas y difíciles de gestionar

### Después (Arquitectura Jerárquica)
```
NotificationType (ej: WhatsApp, SSE)
└── NotificationEventType (ej: order.created, order.shipped)
    └── BusinessNotificationConfig
        ├── business_id (FK)
        ├── integration_id (FK) - NUEVO - La integración origen
        ├── notification_type_id (FK) - NUEVO - Canal de salida
        ├── notification_event_type_id (FK) - NUEVO - Tipo de evento
        ├── enabled (bool)
        ├── filters (JSON)
        └── order_statuses (M2M)
```

**Ventajas:**
- Tipos de notificación normalizados y administrables vía CRUD
- Eventos tipados y validados
- Relación clara con integración origen (ej: Shopify, WhatsApp)
- Configuración más granular y escalable
- Facilita agregar nuevos tipos/eventos sin cambios de código

## Cambios en Base de Datos

### Nuevas Tablas

#### `notification_types`
```sql
CREATE TABLE notification_types (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    description VARCHAR(500),
    icon VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    config_schema JSONB,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
```

**Datos iniciales:**
- SSE (Server-Sent Events)
- WhatsApp
- Email
- SMS

#### `notification_event_types`
```sql
CREATE TABLE notification_event_types (
    id BIGSERIAL PRIMARY KEY,
    notification_type_id BIGINT NOT NULL,
    event_code VARCHAR(100) NOT NULL,
    event_name VARCHAR(200) NOT NULL,
    description TEXT,
    template_config JSONB,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (notification_type_id) REFERENCES notification_types(id),
    UNIQUE(notification_type_id, event_code)
);
```

**Eventos iniciales por tipo:**

**SSE:**
- order.created → "Nueva Orden"
- order.status_changed → "Cambio de Estado de Orden"

**WhatsApp:**
- order.created → "Confirmación de Pedido"
- order.shipped → "Pedido Enviado"
- order.delivered → "Pedido Entregado"
- order.canceled → "Pedido Cancelado"
- invoice.created → "Factura Generada"

**Email:**
- order.created → "Confirmación de Pedido"
- order.shipped → "Pedido Enviado"

### Modificaciones a `business_notification_configs`

**Columnas AGREGADAS:**
- `integration_id` (BIGINT, FK a `integrations`) - La integración origen del evento
- `notification_type_id` (BIGINT, FK a `notification_types`) - Canal de salida
- `notification_event_type_id` (BIGINT, FK a `notification_event_types`) - Tipo de evento

**Columnas DEPRECADAS:**
- `channels` (JSONB) - Eliminada (ahora se usa `notification_type_id`)
- `event_type` (VARCHAR) - Ahora nullable (se usa `notification_event_type_id`)

**Índices AGREGADOS:**
- `idx_bnc_integration_id`
- `idx_bnc_notification_type_id`
- `idx_bnc_notification_event_type_id`
- `idx_bnc_unique_config` (UNIQUE: integration_id, notification_type_id, notification_event_type_id)

**Índices ELIMINADOS:**
- `idx_business_event_type` (UNIQUE: business_id, event_type)

## Estrategia de Migración de Datos

### 1. Asignar `integration_id`
```sql
UPDATE business_notification_configs bnc
SET integration_id = (
    SELECT i.id FROM integrations i
    WHERE i.business_id = bnc.business_id
      AND i.is_active = true
    LIMIT 1
)
WHERE integration_id IS NULL;
```

**Criterio:** Asignar la primera integración activa del business.

### 2. Asignar `notification_type_id`
```sql
UPDATE business_notification_configs bnc
SET notification_type_id = (
    SELECT nt.id FROM notification_types nt
    WHERE nt.code = CASE
        WHEN bnc.channels::text LIKE '%"sse"%' THEN 'sse'
        WHEN bnc.channels::text LIKE '%"whatsapp"%' THEN 'whatsapp'
        WHEN bnc.channels::text LIKE '%"email"%' THEN 'email'
        ELSE 'sse'
    END
    LIMIT 1
)
WHERE notification_type_id IS NULL;
```

**Criterio:** Extraer el primer canal del array `channels`.

### 3. Asignar `notification_event_type_id`
```sql
UPDATE business_notification_configs bnc
SET notification_event_type_id = (
    SELECT net.id
    FROM notification_event_types net
    WHERE net.notification_type_id = bnc.notification_type_id
      AND net.event_code = bnc.event_type
    LIMIT 1
)
WHERE notification_event_type_id IS NULL;
```

**Criterio:** Mapear `event_type` string al evento correspondiente del tipo asignado.

## Ejecución de la Migración

### Opción 1: AutoMigrate + SQL Manual

```bash
# 1. Ejecutar AutoMigrate (crea tablas y columnas)
cd /back/central
go run cmd/main.go migrate

# 2. Ejecutar script SQL (inserta datos y migra configs existentes)
psql -U postgres -d probability_db -f back/migration/shared/sql/migrate_notification_system_refactor.sql
```

### Opción 2: Solo SQL (si AutoMigrate falla)

```bash
psql -U postgres -d probability_db -f back/migration/shared/sql/migrate_notification_system_refactor.sql
```

## Verificación Post-Migración

### 1. Verificar tablas creadas
```sql
\dt notification_*

-- Debe mostrar:
-- notification_types
-- notification_event_types
```

### 2. Verificar datos iniciales
```sql
SELECT * FROM notification_types;
-- Debe mostrar: SSE, WhatsApp, Email, SMS

SELECT * FROM notification_event_types;
-- Debe mostrar eventos para cada tipo
```

### 3. Verificar configs migradas
```sql
SELECT
    bnc.id,
    b.name as business,
    i.name as integration,
    nt.name as notification_type,
    net.event_name,
    bnc.enabled
FROM business_notification_configs bnc
JOIN businesses b ON bnc.business_id = b.id
LEFT JOIN integrations i ON bnc.integration_id = i.id
LEFT JOIN notification_types nt ON bnc.notification_type_id = nt.id
LEFT JOIN notification_event_types net ON bnc.notification_event_type_id = net.id
ORDER BY bnc.id;
```

### 4. Verificar configs sin migrar (ERROR)
```sql
SELECT id, business_id, event_type
FROM business_notification_configs
WHERE integration_id IS NULL
   OR notification_type_id IS NULL
   OR notification_event_type_id IS NULL;

-- Si hay resultados, revisar por qué no se migraron
```

## Rollback

Si necesitas revertir la migración:

```sql
BEGIN;

-- Restaurar event_type como NOT NULL
ALTER TABLE business_notification_configs ALTER COLUMN event_type SET NOT NULL;

-- Restaurar columna channels (requiere backup de datos)
ALTER TABLE business_notification_configs ADD COLUMN channels JSONB;
UPDATE business_notification_configs SET channels = '["sse"]'; -- Valor por defecto

-- Eliminar FKs y columnas nuevas
ALTER TABLE business_notification_configs DROP CONSTRAINT IF EXISTS fk_business_notification_configs_integration;
ALTER TABLE business_notification_configs DROP CONSTRAINT IF EXISTS fk_business_notification_configs_notification_type;
ALTER TABLE business_notification_configs DROP CONSTRAINT IF EXISTS fk_business_notification_configs_notification_event_type;
ALTER TABLE business_notification_configs DROP COLUMN IF EXISTS integration_id;
ALTER TABLE business_notification_configs DROP COLUMN IF EXISTS notification_type_id;
ALTER TABLE business_notification_configs DROP COLUMN IF EXISTS notification_event_type_id;

-- Eliminar tablas nuevas
DROP TABLE IF EXISTS notification_event_types;
DROP TABLE IF EXISTS notification_types;

COMMIT;
```

## Impacto en el Código

### Backend

**Archivos nuevos:**
- `/back/migration/shared/models/notification_type.go`
- `/back/migration/shared/models/notification_event_type.go`
- `/back/central/services/modules/notification_config/internal/domain/entities/notification_type.go`
- `/back/central/services/modules/notification_config/internal/domain/entities/notification_event_type.go`
- Repositorios, Use Cases, Handlers para ambos

**Archivos modificados:**
- `/back/migration/shared/models/notification_config.go`
- `/back/migration/internal/infra/repository/constructor.go`
- `/back/central/services/modules/notification_config/bundle.go`

### Frontend

**Archivos nuevos:**
- `/front/central/src/services/modules/notification-config/infra/actions/notification-types.ts`
- `/front/central/src/services/modules/notification-config/infra/actions/notification-event-types.ts`
- `/front/central/src/app/(auth)/notification-types/page.tsx`
- `/front/central/src/app/(auth)/notification-event-types/page.tsx`

**Archivos modificados:**
- `/front/central/src/services/modules/notification-config/domain/types.ts`
- `/front/central/src/services/modules/notification-config/ui/components/NotificationConfigForm.tsx`

## Notas Importantes

1. **NO se elimina la columna `channels` automáticamente** - Está comentada en el script SQL por seguridad. Eliminar manualmente después de verificar que todo funciona.

2. **Configs existentes se migran automáticamente** - El script SQL intenta migrar todas las configs existentes asignando valores por defecto.

3. **Si un business no tiene integraciones** - La migración fallará para ese business. Revisar manualmente.

4. **El índice unique antiguo se elimina** - Ya no se puede tener múltiples configs con el mismo (business_id, event_type). Ahora es (integration_id, notification_type_id, notification_event_type_id).

## Soporte

Si hay problemas durante la migración:
1. Revisar logs de PostgreSQL: `tail -f /var/log/postgresql/postgresql-15-main.log`
2. Verificar que todas las integraciones estén activas
3. Ejecutar queries de verificación post-migración
4. Si falla, ejecutar rollback y contactar al equipo de desarrollo
