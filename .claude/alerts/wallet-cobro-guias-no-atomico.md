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

## Criterio para cerrar

Cerrar cuando el fix + worker esten verificados en produccion y se agregue el indice
unico parcial anti-doble-cobro.
