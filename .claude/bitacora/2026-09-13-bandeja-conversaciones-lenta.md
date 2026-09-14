# Bandeja de conversaciones: el listado tardaba mas de 120 s y apilaba consultas en RDS

**Ticket:** sin ticket (trabajo hecho en sesion, falta registrarlo)
**Negocio:** Mystic Rose - Official (36), afecta a todo negocio con volumen
**Canal:** WhatsApp, modulo Conversaciones

## Resumen

Tras desplegar el nombre del cliente en la bandeja, el listado de Mystic Rose
dejo de cargar (esqueletos infinitos) y quedaron ~8 consultas activas de minutos
en la base de produccion. Se reescribio la consulta en una sola pasada: de mas
de 120 s a 258 ms.

## Sintoma

- Pantalla de Conversaciones en "cargando" sin fin para Mystic Rose.
- `pg_stat_activity`: 9 consultas `WITH conv AS (...)` activas, la mas vieja
  con 3 min 41 s.
- `EXPLAIN ANALYZE` de la consulta real no termino en 120 s.

Volumen: 667 telefonos, 1.121 conversaciones, 2.175 mensajes, 500 clientes.

## Diagnostico

1. Hipotesis inicial: el `LEFT JOIN LATERAL` nuevo sobre `client` (regexp por
   fila, sin indice). El plan lo mostro barato: index scan por `business_id`,
   ~1 ms por telefono. **No era la causa principal.**
2. El plan real: el lateral `reciente` (ultimo mensaje) hacia un
   `Index Scan Backward` sobre todos los mensajes y, por cada mensaje, un
   `CTE Scan` del CTE `conv` filtrando por `phone_key`. Del orden de
   667 x 2.175 x 1.121 comparaciones. `agg` repetia el patron
   `conversation_id IN (SELECT ... FROM conv WHERE phone_key = k.phone_key)`.
3. Todos los laterales (orden, campana, nombre, baja) se evaluaban para los 667
   telefonos antes del `ORDER BY ... LIMIT`, no para los 20 de la pagina.
4. No se midio la version previa al nombre del cliente, asi que no consta si ya
   era lenta o si el lateral extra la empujo al plan malo.

## Causa raiz

Subconsultas correlacionadas sobre un CTE materializado (`IN (SELECT ... FROM
conv ...)`) dentro de `LATERAL`, reevaluadas por cada telefono y por cada
mensaje.

## Correccion

Commit `3411ed30` (`message_audit_queries.go`):

- `ultima` con `DISTINCT ON`, `primera` con `MIN`, `msgs` como un solo join
  conversaciones-mensajes, `agg` y `reciente` agregados por telefono.
- CTE `pagina` con `ORDER BY` + `OFFSET/LIMIT`; orden, campana, nombre y baja se
  calculan solo para las filas de la pagina.
- `CountUnreadConversations` con el mismo enfoque (`COUNT(DISTINCT phone_key)`).

## Verificacion

- Produccion, solo lectura: listado 258 ms, contador de sin leer 11 ms.
- Local: listado de Demo identico fila por fila al anterior, 65 ms.
- Tras el deploy: 0 consultas lentas activas. Las consultas atascadas terminaron
  solas, no hubo que cancelarlas.
- Contador de sin leer comparado solo en Demo (0 y 0): cobertura debil.

## Pendientes

- Registrar el ticket.
- Considerar indice funcional por telefono normalizado en `client` y
  `whatsapp_conversations` si crece el volumen.
