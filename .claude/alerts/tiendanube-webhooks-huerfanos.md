# Webhooks huerfanos al eliminar una integracion (Tiendanube y hermanos)

Fecha: 2026-08-22
Origen: `.claude/bitacora/2026-08-22-tiendanube-oauth-5-ciclos.md`

## Contexto

Al crear una integracion de Tiendanube, `bundle.go` engancha `OnIntegrationCreated`
y crea 7 webhooks en la tienda apuntando a
`/api/v1/tiendanube/webhook?integration_id=<id>`.

Al **eliminar** la integracion no existe el gancho inverso: los webhooks quedan
vivos en la tienda apuntando a un `integration_id` que ya no existe. Cinco ciclos
de borrar/reconectar dejaron 35 webhooks, 28 de ellos huerfanos.

## Items

### Urgente

- [ ] Borrar los webhooks del canal cuando se elimina la integracion (equivalente
      a `OnIntegrationDeleted`). Aplica a Tiendanube y a cualquier canal que cree
      webhooks al conectar.
- [ ] Que el handler del webhook descarte con `Warn` + 200 cuando el
      `integration_id` no existe o esta borrado, en vez de tratarlo como error.
      Un webhook huerfano dispara en cada evento de la tienda y hoy es ruido de
      log permanente (ver `.claude/rules/colas-errores-permanentes.md`).

### Importante

- [ ] `CreateWebhooks` no deduplica: si se vuelve a conectar la misma tienda se
      suman juegos nuevos. Deberia listar y borrar los que apunten a este mismo
      host antes de crear.
- [ ] Modales de sincronizacion de productos de **Shopify** y **Jumpseller**:
      leen la respuesta 202 de `reconcile` como si fuera el diff y pintan
      "Todo sincronizado" con 0. Mismo bug que ya se corrigio en Tiendanube
      (escuchar `<canal>.product.reconcile.completed` y leer el detalle de
      `/internal/sync-run-items`). Woo y VTEX no aplican: su reconcile es sincrono.

### Deseable

- [ ] Boton en el gestor de webhooks para "limpiar webhooks que no son de esta
      integracion", que hoy hay que resolver a mano con
      `DELETE /api/v1/integrations/<id>/webhooks/<webhook_id>`.

## Criterio para cerrar

Eliminar una integracion de Tiendanube deja la tienda con cero webhooks de
Probability, y reconectar dos veces seguidas deja exactamente 7.
