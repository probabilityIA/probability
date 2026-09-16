# Los estados elegidos en las reglas de notificacion nunca se guardaban

**Ticket:** TKT-000099. Afecta a todos los negocios, canal WhatsApp y Via.

## Resumen

El modal "Reglas de Notificacion" deja elegir estados de orden por regla, pero
ninguno llegaba a la base. Toda regla se comportaba como "Todos los estados".

## Sintoma

- Produccion: 62 reglas activas y **0 filas** en
  `business_notification_config_order_statuses`. 40 de esas reglas son de eventos
  que admiten filtro de estado.
- La API de sync respondia `order_status_ids: []` sin error, aunque el request
  llevaba IDs.

## Diagnostico

- Se descarto primero un problema de cache: Redis tenia los IDs porque se
  cachean desde el entity en memoria, no desde la base.
- Se descarto el endpoint de creacion: el de sync (el que usa el modal) tampoco
  guardaba, ni al crear ni al actualizar.
- El repositorio usaba `tx.Model(...).Association("OrderStatuses").Replace(...)`
  y no devolvia error.

## Causa raiz

`shared/db/db.go` abre la conexion global con
`d.conn.Omit(clause.Associations)`. Con ese Omit, `Association(...).Replace` y
`Clear` no escriben nada y tampoco fallan. El mismo patron estaba en los estados
permitidos de `notification_event_types`.

**Cualquier many2many escrito con `Association()` en este backend es un no-op
silencioso.** Hay que escribir la tabla intermedia con SQL.

## Correccion

- `notification_config/.../repository/sync_configs.go`: `replaceConfigStatuses`
  borra e inserta la tabla intermedia con SQL en crear, actualizar y borrar.
- `notification_event_type_repository.go`: `replaceAllowedStatuses`, mismo
  arreglo para los estados permitidos por evento.

## Verificacion

Local: sync con una regla actualizada (Via, cancelada + cliente solicita
cancelar) y una creada (WhatsApp, enviada) deja 2 y 1 filas en la tabla
intermedia. Un negocio creado por API nace con la regla predeterminada y sus
dos estados guardados.

## Pendientes

- Las reglas existentes en produccion no tienen estados: quien los haya elegido
  debe volver a elegirlos despues del deploy.
- Buscar otros `Association(...)` en el backend (hoy solo estaban estos dos).
