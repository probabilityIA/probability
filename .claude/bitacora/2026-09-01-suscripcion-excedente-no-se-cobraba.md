# Suscripciones: el excedente de uso nunca se cobraba al renovar

## Resumen

`PurchaseSubscription` (el cobro real via billetera cuando un negocio paga/
extiende su plan) calculaba el monto como `precio_plan x meses`, sin sumar
nunca el excedente de envios/facturas/ordenes del ciclo que termina. El
"pago pronosticado" que se veia en el dashboard (correcto, ej. $142.800 para
Mystic con 173/100 envios) era solo un calculo de exhibicion en otro archivo
que nunca se conectaba al cobro real: el negocio pagaba $99.000 (solo la
membresia) y el excedente quedaba perdonado en cada renovacion.

Reportado por el usuario con capturas de Mystic Rose (business real): plan
100 envios incluidos, 173 usados, $600 c/u de excedente = deberia cobrar
$142.800, el modal de compra solo mostraba y cobraba $99.000.

## Causa raiz (verificado en codigo)

`back/central/services/modules/subscriptions/internal/app/purchase_subscription.go`:

```go
amount := subType.Price * float64(dto.Months)
```

El unico lugar que si calculaba correctamente el excedente era
`forecastNextPayment` en `admin_businesses_repository.go`, usado solo para
el campo `forecasted_payment` de exhibicion (`GetSubscriptionUsage`,
`ListBusinessesForAdmin`). Nunca se le paso ese numero al flujo de cobro.

## Fix

- Nuevo `ports.IRepository.ComputeOverageAmount` (repo
  `overage_amount.go`): misma logica de excedente, ahora reutilizable.
- `forecastNextPayment` refactorizado para llamar a `ComputeOverageAmount`
  en vez de duplicar el calculo.
- `PurchaseSubscription` ahora consulta la suscripcion vigente antes de
  cobrar, calcula su excedente real (usando su propio plan y ciclo, no el
  plan nuevo que se este comprando) y lo suma a `amount` antes de validar
  saldo y debitar.
- Front (`subscription/page.tsx`, modal "Contratar Plan"): agrega linea
  "Excedente del ciclo actual" cuando aplica, y el total/saldo insuficiente
  ya lo incluyen. `usage` se pasa como prop nueva a `PlanCatalog`.
- 2 tests nuevos (`TestPurchaseSubscription_SumaElExcedenteDelCicloQueTermina`,
  `..._SinSuscripcionPrevia_NoSumaExcedente`) mas el fix de un test existente
  que se rompio por el refactor de `computeSubscriptionWindow`.

## Verificacion

Recreado en local: plan con 100 envios incluidos / $600 excedente, 173
envios reales insertados, `GetSubscriptionUsage` devolvio
`forecasted_payment: 142800` (igual que Mystic). Llamando
`POST /subscriptions/purchase` el negocio de prueba fue cobrado
exactamente $142.800 (antes del fix hubiera sido $99.000), y el saldo de
billetera bajo en esa misma cantidad. Datos de prueba limpiados despues.

## Pendiente

- No se corrigio retroactivamente ninguna renovacion pasada que ya cobro de
  menos (Mystic y cualquier otro negocio con excedente historico). Decidir
  con el usuario si se cobra el faltante acumulado o se deja como perdida
  asumida hasta hoy.
- Solo cubre el flujo de auto-servicio (`PurchaseSubscription`, boton
  "Pagar/Extender" del negocio). `RegisterPayment` (pago manual registrado
  por un admin) no se toco: ahi el monto lo escribe el admin a mano.
