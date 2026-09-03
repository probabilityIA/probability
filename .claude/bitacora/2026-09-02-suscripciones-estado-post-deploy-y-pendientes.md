# Suscripciones: estado real de cobertura tras el deploy del 2026-09-02

Nota de cierre para tener a mano en la proxima reunion: que tan probado quedo
el modulo de suscripciones despues del deploy, y que falta.

## Que se desplego

Los dos bugs graves y sus tests de
`2026-09-02-suscripciones-control-de-acceso-y-planes.md`:

1. `HasModuleAccess` ya no le da mas acceso a una cuenta deshabilitada que a un
   plan pagado (el fallback de "sin plan" solo aplica a cuentas activas).
2. Los planes personalizados nuevos ya nacen `payable=true` y se pueden
   autopagar/renovar.

Junto con el resto de la sesion del dia: dias de cortesia que ya no mueven el
rango de facturacion (`subscription_courtesy_until`), reactivacion automatica
al editar fechas o extender cortesia, y el fix del boton "Pagar/Extender" para
planes personalizados.

## Que quedo probado en vivo (CU-01 a CU-15)

CRUD completo de super admin, compra self-service, planes personalizados,
ownership entre negocios, overrides de modulos, cortesia/reactivacion, dia de
corte, y el frontend (boton Pagar/Extender + SubscriptionGuard) con Playwright
contra el backend real. Detalle en
`.claude/testing/subscriptions/back/RESULTS.md` y
`.claude/testing/subscriptions/front/RESULTS.md`.

## Que NO quedo con cobertura completa

- **Excedentes al renovar (CU-12)**: solo tests unitarios. No se armo el
  escenario real con mas de 50 envios en un ciclo para verlo cobrar de punta a
  punta.
- **Worker de vencimiento** (`CheckExpiringSubscriptions`, el cron que marca
  cuentas `expired`): no se disparo manualmente en esta sesion.
- Hallazgos menores documentados pero **sin corregir**:
  - `PUT /custom-plans/{id}` es reemplazo total; un body parcial borra
    `active`/`description`.
  - Borrar un custom plan mientras un negocio esta activo en el deja
    `business.subscription_type_id` huerfano (mismo sintoma del bug 1, pero la
    causa de fondo -- que dejarlo borrar sin reasignar es la puerta de entrada
    -- no se cerro).
  - `disable` recibe `business_id` por query param, `reactivate` por body.

## Efecto colateral aceptado

`business_id 38` ("Alejandra", negocio real en produccion sin plan asignado)
pierde el acceso amplio que tenia por el bug 1. El usuario confirmo que no es
prioridad remediarlo.

## Recomendacion

No se puede llamar "sin errores garantizado". Vigilar logs de produccion las
primeras horas/dias despues del deploy, en particular `register-payment`,
`extend-days`, `edit-dates` y `custom-plans` (los endpoints tocados). Si
aparece un negocio con acceso reducido inesperadamente, revisar primero si es
uno de los que dependia del fallback de "sin plan" (business 38 y los 4 de
prueba en produccion).
