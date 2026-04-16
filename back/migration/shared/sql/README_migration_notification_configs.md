# Migración: Tabla Intermedia para Estados de Orden en Notificaciones

## Resumen

Esta migración **NO elimina ninguna tabla**. Solo agrega una nueva tabla intermedia para relacionar configuraciones de notificaciones con estados de orden de forma más directa.

## Tablas afectadas

### ✅ Tablas que SE MANTIENEN (sin cambios)
- `business_notification_configs` - **NO se elimina, NO se modifica su estructura**
- `order_statuses` - **NO se elimina, NO se modifica**

### ➕ Nueva tabla que se CREA
- `business_notification_config_order_statuses` - Tabla intermedia many-to-many

## Pasos de migración

### Opción 1: Usando GORM AutoMigrate (Recomendado)

El modelo ya está agregado a `AutoMigrate`, así que solo ejecuta:

```bash
# Ejecutar la migración de GORM
# (Esto creará automáticamente la nueva tabla si no existe)
```

### Opción 2: Ejecutar SQL manualmente

Si prefieres ejecutar el SQL manualmente:

```bash
# Ejecutar el script de creación de tabla
psql -d tu_base_de_datos -f back/migration/shared/sql/business_notification_config_order_statuses.sql

# O usar el script de migración completo que incluye la creación
psql -d tu_base_de_datos -f back/migration/shared/sql/migrate_notification_configs_to_order_statuses.sql
```

## Migración de datos (Opcional)

Si tenías estados en el campo JSON `filters` de la tabla `business_notification_configs`, puedes migrarlos usando el script comentado en `migrate_notification_configs_to_order_statuses.sql`.

**Ejemplo de datos a migrar:**
```json
{
  "statuses": ["pending", "processing", "delivered"]
}
```

Esto se convertiría en registros en `business_notification_config_order_statuses`:
```
business_notification_config_id | order_status_id
--------------------------------|----------------
1                               | 1  (pending)
1                               | 2  (processing)
1                               | 4  (delivered)
```

## Verificación

Después de la migración, puedes verificar con:

```sql
-- Ver la nueva tabla
SELECT * FROM business_notification_config_order_statuses;

-- Ver configuraciones con sus estados asociados
SELECT 
    bnc.id as config_id,
    bnc.event_type,
    os.code as status_code,
    os.name as status_name
FROM business_notification_configs bnc
LEFT JOIN business_notification_config_order_statuses bcos 
    ON bnc.id = bcos.business_notification_config_id
LEFT JOIN order_statuses os 
    ON bcos.order_status_id = os.id
WHERE bnc.event_type = 'order.status_changed'
ORDER BY bnc.id, os.code;
```

## Rollback (si es necesario)

Si necesitas revertir (no recomendado, pero por si acaso):

```sql
-- Eliminar la tabla intermedia
DROP TABLE IF EXISTS business_notification_config_order_statuses;
```

**NOTA:** Esto solo eliminaría los datos de la relación. La tabla `business_notification_configs` seguiría intacta.
