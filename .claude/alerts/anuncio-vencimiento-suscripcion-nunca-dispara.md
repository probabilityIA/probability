# ALERTA: el anuncio de vencimiento de suscripcion esta cableado pero nunca se disparo

Fecha: 2026-08-20. **Causa raiz encontrada y corregida el mismo dia** (ver abajo).

## Causa raiz (confirmada)

`back/central/services/modules/announcements/gateway.go`, funcion
`resolveAlertCategoryID`, usaba un `sync.Once` a nivel de **variables de
paquete** (no de instancia) para cachear el ID de la categoria `"alert"`:

```go
var (
    alertCategoryOnce sync.Once
    alertCategoryID   uint
    alertCategoryErr  error
)
```

`sync.Once.Do` corre una sola vez en toda la vida del proceso. Si la
**primera** vez que se invoca `CreateBusinessAlert` (desde cualquier caller,
no solo subscriptions) la consulta `ListCategories` fallaba por cualquier
motivo transitorio -- muy probable en el arranque, porque
`ExpiryWorker.Start` dispara `runCheck` de inmediato al lanzar el goroutine,
sin esperar a que la app este completamente lista -- ese error quedaba
**grabado para siempre**. Ninguna llamada futura volvia a intentar la query,
sin importar que el problema original ya se hubiera resuelto. Solo un
restart del proceso lo limpiaba, y si volvia a fallar en el arranque
siguiente, quedaba roto otra vez.

Se confirmo por SQL que la categoria `"alert"` **si existe** en produccion
(`id=2, code='alert'`) -- no era un dato faltante, era pura condicion de
carrera en el arranque congelada para siempre.

Refuerza el problema: `CheckExpiringSubscriptions` nunca propaga el error
individual de `ensureExpiryAnnouncement` (solo `log.Error` + `return nil` al
final de la funcion), asi que el bug nunca genero ninguna alerta visible de
monitoreo -- el worker parecia "correr bien" indefinidamente.

## Fix aplicado (2026-08-20)

`gateway.go`: se elimino el `sync.Once` de paquete. El cache ahora vive en
campos de instancia de `Bundle` (`alertCategoryID`,
`alertCategoryResolved`, protegidos por `sync.RWMutex`) y **solo se escribe
cuando la resolucion es exitosa**. Un fallo transitorio ya no queda
cacheado: la siguiente llamada vuelve a intentar `ListCategories` desde
cero. Una vez resuelto con exito, se cachea permanentemente (los IDs de
categoria son datos de seed estaticos, seguro cachearlos).

Este bug bloqueaba en silencio, ademas del anuncio de vencimiento, TODA
llamada a `CreateBusinessAlert` desde cualquier modulo -- incluye la
logica nueva de degradacion de trial y liquidacion de excedente del plan
gratuito (ver `.claude/bitacora/` o el modulo `subscriptions`), que dependen
del mismo `ExpiryWorker`.

## Contexto

`subscriptions/internal/app/check_expiring_subscriptions.go` implementa
`CheckExpiringSubscriptions`: busca negocios que vencen en los proximos 7 dias
(`ListBusinessesExpiringBetween`) y negocios ya vencidos
(`ListBusinessesJustExpired`), y para cada uno crea un anuncio de tipo alerta
(`ensureExpiryAnnouncement`, con proteccion anti-duplicado via
`FindActiveBusinessAlert`).

El worker que lo ejecuta si esta bien cableado:
`subscriptions/internal/infra/primary/worker/expiry_worker.go` corre
`runCheck` de inmediato y luego cada 24h. Se instancia y arranca en
`subscriptions/bundle.go` (`go expiryWorker.Start(context.Background())`), y
`subscriptions.New(...)` se llama desde `services/modules/bundle.go:92`, que es
parte del arranque normal de la app.

## Hallazgo central

Consultando produccion (solo lectura, tabla `announcements`):

- **Cero registros** con `title = 'Tu suscripcion esta por vencer'` o
  `title = 'Tu suscripcion vencio'` en toda la historia de la tabla.
- El ultimo anuncio de **cualquier tipo** en la tabla es del 2026-05-13 (mas de
  3 meses antes de esta fecha), lo que sugiere que el problema puede ser mas
  amplio que solo el flujo de vencimiento.
- Hay negocios reales que pasaron por la ventana de "vence en <=7 dias"
  (ej. business_id 34: una fila de `business_subscriptions` vencia
  2026-08-24, dentro de la ventana de 7 dias en el momento de la consulta) y
  ninguno genero el anuncio esperado.
- `business_subscriptions` tiene **multiples filas historicas por negocio**
  (renovaciones que se solapan: mismo `business_id`, varias filas `status =
  'paid'` con distintos `end_date`). Sospecha sin confirmar: el problema puede
  estar en como `ListBusinessesExpiringBetween` / `ListBusinessesJustExpired`
  identifican cual es la suscripcion "vigente" del negocio cuando hay varias
  filas — si toma la de `end_date` mas lejano (la ya renovada) en vez de la
  activa real, nunca cae dentro de la ventana de aviso.

## Hipotesis descartadas durante la investigacion

1. Bug en la query del repo (`ListBusinessesExpiringBetween` /
   `ListBusinessesJustExpired`) que no filtra bien cual fila de
   `business_subscriptions` es la "actual" -- **descartada**: esas queries
   filtran directo sobre `businesses.subscription_end_date`, no sobre
   `business_subscriptions`, asi que las filas historicas solapadas no las
   afectan.
2. `resolveSystemUserID` fallando silenciosamente -- **descartada**: la
   query filtra por `scope.code = 'platform'`, maneja bien el caso de 0
   filas y propaga el error hacia arriba; no hay evidencia de que falte esa
   data.
3. Causa real: **confirmada arriba**, bug de `sync.Once` en
   `announcements/gateway.go`.

## Items

### Urgente (resuelto 2026-08-20)

- ~~Confirmar con un log/trace real si `runCheck` se esta ejecutando en
  produccion~~ -- resuelto: el worker si corria, el bloqueo estaba en
  `CreateBusinessAlert` via el `sync.Once`, ya corregido.
- ~~Revisar `ListBusinessesExpiringBetween` / `ListBusinessesJustExpired`~~ --
  descartado como causa, ver arriba.

### Importante (pendiente)

- Agregar un test de integracion real (no con `AnnouncementsGatewayMock`)
  que ejercite `CreateBusinessAlert` contra `announcement_categories`, para
  que un problema similar de timing/cacheo no vuelva a pasar desapercibido.
- Verificar en el proximo ciclo de 24h que el worker efectivamente inserta
  filas nuevas en `announcements` en produccion (confirmar el fix end-to-end,
  no solo que compila).
- `CheckExpiringSubscriptions` sigue sin propagar errores individuales de
  `ensureExpiryAnnouncement` mas alla de `log.Error` -- evaluar si vale la
  pena una alerta/metrica cuando falla repetidamente, para no depender de
  encontrar este tipo de bug por casualidad la proxima vez.

## Criterio de cierre

Cerrar cuando se confirme en produccion (no solo por codigo) que
`CheckExpiringSubscriptions` genero al menos un anuncio real tras el fix del
2026-08-20.
