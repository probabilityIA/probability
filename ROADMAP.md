# Roadmap

Este archivo SOLO ordena. El detalle de cada punto vive en su alerta o documento;
aca no se duplica contenido, se decide que sigue.

Ultima revision: 2026-08-10 (14 alertas abiertas en `.claude/alerts/`).

Como usarlo:
- Se atiende de arriba hacia abajo. P0 antes que P1, P1 antes que P2.
- Al cerrar un punto: actualizar la alerta (marcar item con fecha), y cuando sus
  items urgentes e importantes esten resueltos, borrar la alerta y crear la
  entrada de bitacora. Recien ahi se saca la linea de aca.
- Si algo urgente nuevo no entra en P0, revisar si de verdad es urgente.

---

## P0 - Seguridad y dinero (abierto hoy)

| # | Que | Donde | Estado |
|---|-----|-------|--------|
| 1 | **Endpoints sin scoping por business (IDOR multi-tenant).** WhatsApp toma `business_id` del body; 6 endpoints `:id` de `orders` sin filtrar; fail-open cuando el negocio no tiene `business_resource_configured`. | `.claude/alerts/autorizacion-backend.md` | 12 items abiertos, 7 resueltos el 2026-07-17 |
| 2 | **API key de OpenRouter viva en el codigo y en el historial de git.** No esta revocada: reemplazarla sin revocar no sirve. | `.claude/alerts/openrouter-api-key-hardcodeada.md` | 6 items, ninguno hecho |
| 3 | **Consumidores que mueren sin reconectar.** El fix esta escrito y testeado pero NO desplegado; hasta que se despliegue se puede repetir la caida de 9h sin facturar del 2026-08-05. | `.claude/alerts/consumidor-muerto-por-canal-cerrado.md` | fix listo, falta deploy + verificar la ventana 17:47-03:01 |
| 4 | **Reembolsos pendientes por guias duplicadas.** 3 ordenes con doble debito; al menos 16.353 confirmados sin acreditar. Ademas el handler de cancelacion no reembolsa el wallet. | `.claude/alerts/guias-duplicadas-doble-cobro.md` | codigo cerrado, plata sin devolver |
| 5 | **Mensajes de RabbitMQ sin `DeliveryMode: Persistent`.** Colas durable pero mensajes no: un restart del broker se come lo que este en vuelo. Documentado como bug de hoy, no como mejora futura. | `.claude/docs/escalabilidad-1m-ordenes-mes.md` (Fase 1) | diagnosticado, nada implementado |

Criterio de P0: o hay datos de otro negocio expuestos, o hay plata que se pierde,
o hay un incidente que ya ocurrio y puede repetirse identico.

---

## P1 - Endurecer lo que ya funciona

| # | Que | Donde |
|---|-----|-------|
| 6 | Indice unico parcial en `transaction(shipment_id) where type='USAGE'`. La idempotencia del cobro es por SELECT, hay ventana de carrera. Ademas el worker de reconciliacion queda en monitoreo permanente: que cobre significa que el flujo inline fallo. | `.claude/alerts/wallet-cobro-guias-no-atomico.md` |
| 7 | SoftPymes: nack/requeue cuando no se puede procesar (hoy ACKea y el estado queda congelado) y cerrar las facturas pendientes del backfill. | `.claude/alerts/softpymes-timeouts-masivo.md` |
| 8 | Alerta de "cola con mensajes y cero consumidores". Es la deteccion que hubiera avisado a las 17:47 en vez de a las 03:00. | `.claude/alerts/consumidor-muerto-por-canal-cerrado.md` |
| 9 | Siigo: cerrar los pendientes criticos y validar E2E las 6 operaciones. | `.claude/alerts/siigo-pendientes.md` |
| 10 | MELI: token a Redis (el core ya recibe `redis.IRedis`, falta cablearlo) y verificacion E2E de >6h con varios ciclos de refresh. | `.claude/alerts/meli-token-persistence.md` |
| 11 | Bancolombia QR: el webhook falla ABIERTO si no hay `webhook_secret` cargado. | `.claude/alerts/bancolombia-qr-spec-pendiente.md` |

---

## P2 - Producto y features

| # | Que | Donde |
|---|-----|-------|
| 12 | Inventario saliente filtrado por canal (`product_business_integrations`). | `.claude/alerts/inventario-saliente-por-canal.md` |
| 13 | Sync de inventario Siigo -> Probability -> WooCommerce. Fase 1 hecha; Fases 2/3 pendientes, falta confirmar el shape del payload de Siigo. | `.claude/alerts/woocommerce-inventory-push.md` |
| 14 | WooCommerce checkout controlado, Fase 2 (ciudad como dropdown, validacion de direccion). | `.claude/alerts/woocommerce-checkout-controlado.md` |
| 15 | Tienda publica: pago en linea con credenciales propias del negocio. | `.claude/alerts/tienda-pago-online-credenciales-propias.md` |
| 16 | Autoregistro demo. | `PLAN-demo.md` |

---

## Higiene de documentacion

No es trabajo de producto pero mantiene util el resto del sistema:

- **Alertas candidatas a cerrar y pasar a bitacora.** `meli-token-persistence`
  (urgentes e importantes marcados resueltos y verificados) y `meli-guia-pdf-proxy`
  (items resueltos, solo falta la verificacion E2E, que hoy no se puede hacer
  porque no existe ningun shipment de MELI con guia). Confirmar y moverlas.
- **La bitacora tiene 1 entrada frente a 14 alertas.** Varias alertas son
  incidentes de produccion ya diagnosticados (las 9h sin facturar, los timeouts
  masivos del business 34). Ese material es bitacora una vez cerrado el pendiente.
- **Escalabilidad**: `.claude/docs/escalabilidad-1m-ordenes-mes.md` mezcla un bug
  de hoy (punto 5 de P0) con mejoras a futuro. Solo la Fase 1 es P0.
