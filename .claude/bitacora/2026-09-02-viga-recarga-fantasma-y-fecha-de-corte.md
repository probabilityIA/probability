# Viga: recarga fantasma por acceso equivocado, y por que "editar fechas" no desbloqueaba la cuenta

Un compañero entro por error a la cuenta de Viga (creyendo entrar a la propia) y
alcanzo a iniciar una recarga de $323.907 que nunca se completo. Aparte, cuando
se intento desbloquear la cuenta de Viga (vencida) editando la fecha de corte,
la cuenta seguia bloqueada.

## Sintoma

- Transaccion `328a6235-8584-465a-ba20-6dcad56fe496`, tipo `RECHARGE`, estado
  `PENDING`, $323.907, creada 2026-09-02 03:54 UTC (10:54 p.m. hora Colombia del
  2026-09-01), business_id 46 (Viga ropa deportiva).
- El equipo confirmo que nunca se llego a abrir el checkout de Bold: el
  compañero se dio cuenta del error antes de completar nada.
- Por separado: se edito la fecha de corte de Viga (`PUT /edit-dates`) para
  desbloquear la cuenta, pero `GET /me` y el guard de la app seguian mostrando
  la cuenta suspendida.

## Diagnostico

**Recarga fantasma**: `bold_recharge.go` crea la fila `PENDING` en el momento en
que el frontend pide la firma de Bold (`getBoldSignatureAction`), antes de que
el usuario siquiera vea el boton de pago. Como el checkout nunca se abrio, Bold
jamas recibio esa orden (`order_id` `WLTb82e12a823864683a2b52d138`) y nunca va a
mandar un webhook de confirmacion o rechazo para ella. Se verifico que crear una
recarga `PENDING` no toca `wallet.balance` (eso solo pasa cuando el webhook
confirma `COMPLETED`), asi que no hubo perdida ni movimiento real de dinero.

**Fecha de corte sin desbloquear**: revisando `edit_subscription_dates.go` y
`extend_courtesy.go`, ninguno de los dos tocaba `business.subscription_status`.
Solo mueven `end_date`/`courtesy_until`. El bloqueo real (`SubscriptionGuard` en
el frontend) lee `subscription_status` desde el JWT, y ese campo solo lo cambia
`CreateSubscriptionAndActivate` (un pago real) o `ReactivateSubscription` (el
boton aparte). Editar fechas sin ademas registrar un pago o reactivar
manualmente deja la fecha corregida pero la cuenta sigue `expired`/`cancelled`.

## Correccion

- **Recarga fantasma**: se confirmo con el usuario que era segura de borrar (sin
  riesgo de webhook tardio porque Bold nunca vio la orden) y se elimino la fila
  con un `DELETE` puntual en produccion, autorizado explicitamente. `wallet.balance`
  de Viga no se toco (verificado antes y despues: sin cambio).
- **Fecha de corte**: se agrego reactivacion automatica a `EditSubscriptionDates`
  y `ExtendCourtesy` (y una nueva columna `subscription_courtesy_until` para que
  la cortesia posponga el bloqueo sin mover el rango de facturacion). Ver
  `.claude/testing/subscriptions/back/RESULTS.md` para el detalle y las pruebas.

## Pendiente

- Ninguno del lado de Viga: cuenta reactivada, balance intacto.
- El fix de reactivacion automatica va empaquetado junto a otros bugs de
  suscripciones encontrados en la sesion de testing del 2026-09-02
  (ver `2026-09-02-suscripciones-control-de-acceso-y-planes.md`), sin desplegar
  aun.
