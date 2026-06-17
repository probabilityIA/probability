# ALERTA: Cobro de guias al wallet no es atomico ni recuperable

**Fecha:** 2026-06-16
**Modulo:** shipments (consumer de respuestas transport) + pay/wallet
**Severidad:** URGENTE (perdida de dinero)

## Contexto

El cobro al wallet por generacion de guia ocurre en el consumer async de RabbitMQ
`handleGenerateResponse` (`services/modules/shipments/internal/infra/primary/queue/consumer/response_consumer.go`),
en un paso separado y "fire-and-forget" despues de guardar la guia.

Dos formas de perder el cobro sin recuperacion:
1. La respuesta del proveedor llega con `business_id = 0` -> el debito se saltaba en
   silencio (solo un Warn en `resolveBusinessID`). La guia igual quedaba creada.
2. `DebitWalletForGuide` falla (hipo de DB, race) -> solo se logueaba; sin retry,
   sin DLQ, sin rollback, sin reconciliacion.

## Impacto medido (al 2026-06-16)

24 guias generadas sin cobrar, **$368.864 COP** no debitados (negocios con wallet, sin tests):
- Mystic Rose (36): ~$300.861 (19 guias) — incluye 2 de hoy ($33.397): shipments 34673, 34674
- LaPerchaDel10 (37): $51.556 (4 guias)
- Demo (26): $16.447 (1 guia, cuenta interna)

Patron por rafagas: feb 25-26, mar 3-9, jun 16. Abril/mayo limpios. Exito ~97.4%.

## Items

### Resuelto (branch fix/wallet-cobro-guia-atomico)
- [x] Recuperar `business_id` desde la orden cuando la respuesta llega en 0
  (elimina el salto silencioso). `response_consumer.go`.
- [x] Loguear `Error` (no saltar silencioso) cuando business_id no se resuelve o
  el debito falla, con shipment_id para reconciliar.
- [x] `DebitWalletForGuide` idempotente: no cobra dos veces si ya existe USAGE para
  ese shipment_id (habilita reintentos y reconciliacion segura). `wallet_queries.go`.
- [x] Worker de reconciliacion periodico (cada 10 min): detecta guias con tracking +
  sin USAGE en las ultimas 24h (con gracia de 5 min para no chocar con el flujo
  inline) y las cobra de forma idempotente. `reconciliation_worker.go` + `bundle.go`.

### Decision del usuario (2026-06-16)
Los $368.864 historicos NO se regularizan desde aqui: soporte ya los corrige
manualmente. El worker solo mira una ventana de 24h, por lo que NO tocara el backlog
antiguo; y si soporte registra el cobro como USAGE, el worker lo respeta (idempotente).

### IMPORTANTE pendiente
- [ ] Garantia fuerte anti-doble-cobro a nivel DB: indice unico parcial en
  transaction(shipment_id) where type='USAGE' (la idempotencia actual es por SELECT,
  hay ventana de carrera teorica bajo concurrencia).

## Monitoreo continuo del worker de reconciliacion (DEJAR ABIERTO)

El worker `shipments.wallet_reconciliation` es la red de seguridad. Hay que vigilarlo:

- **Que cobre = que algo fallo aguas arriba.** Si en los logs aparece seguido
  `Wallet reconciliation: guide charged` o `cycle complete recovered>0`, significa que
  el cobro inline (en `handleGenerateResponse`) esta fallando y el worker esta tapando.
  Investigar la causa raiz de ese ciclo (business_id, error de DB, etc.), no solo confiar
  en el worker.
- **Que falle el worker.** Vigilar `Wallet reconciliation: failed to debit guide` y
  `failed to query uncharged guides`. Si se repiten, el worker no esta cumpliendo.
- **Backlog fuera de ventana.** El worker solo mira 24h. Revisar periodicamente que no
  haya guias con tracking sin USAGE mas viejas de 24h (se le escaparian).

Consulta de control (guias sin cobro, negocios con wallet, no test):
```sql
SELECT count(*), sum(s.total_cost)
FROM shipments s JOIN orders o ON o.id=s.order_id JOIN wallet w ON w.business_id=o.business_id
WHERE s.tracking_number<>'' AND s.total_cost>0 AND s.deleted_at IS NULL AND s.is_test IS NOT TRUE
  AND NOT EXISTS (SELECT 1 FROM transaction t WHERE t.shipment_id=s.id AND t.type='USAGE');
```
Idealmente devuelve 0. Si crece, el worker no esta alcanzando o hay un fallo nuevo.

## Criterio para cerrar

NO cerrar el item de monitoreo (es permanente). Los demas: cerrar cuando el fix + worker
esten verificados en produccion y se agregue el indice unico parcial anti-doble-cobro.
