# Alerta: MercadoLibre - persistencia de token y hardening de conexion

Fecha: 2026-07-14
Modulo: `back/central/services/integrations/ecommerce/meli`

## Contexto

Al completar la integracion de MercadoLibre (ordenes, inventario multi-bodega,
productos, webhooks por topic, push de estado) se detecto un bug pre-existente en
el flujo de refresh de token que impide que la integracion funcione mas alla del
primer ciclo de token (~6h post-conexion). No se corrigio en la misma sesion
porque el fix toca el path compartido y sensible de escritura de credenciales del
core de integraciones, y requiere sign-off.

## Items

### URGENTE (bloqueante de produccion)

`refresh_token.go` -> `refreshAccessToken` NO persiste:
- el nuevo `access_token`, ni
- el `refresh_token` rotado (ML rota el refresh_token en cada refresh).

Solo actualiza `config["token_expires_at"]` (+ `seller_id`) via `UpdateIntegrationConfig`.

Efecto: tras el primer refresh, `EnsureValidToken` toma la rama "token vigente"
(porque `token_expires_at` dice +6h) y hace `DecryptCredential("access_token")`,
que devuelve el token VIEJO (stale) -> 401 en toda llamada. Ademas ML puede
invalidar el refresh_token viejo tras la rotacion -> el proximo refresh falla.

Fix propuesto (aditivo):
1. Agregar `UpdateIntegrationCredentials(ctx, integrationID string, creds map[string]interface{}) error`
   a `core.IIntegrationCore` (facade en `services/integrations/core/bundle.go`),
   implementado con `useCase.UpdateIntegration(ctx, id, domain.UpdateIntegrationDTO{Credentials: &creds})`
   (el repo re-encripta; el update REEMPLAZA el mapa completo de credenciales).
2. Exponerlo en `meli/internal/domain/ports.go` (`IIntegrationService`) y en el
   adaptador `meli/internal/infra/secondary/core/integration_service.go`.
3. En `refreshAccessToken`, tras el refresh, construir el mapa completo
   `{access_token: nuevo, refresh_token: rotado (o el actual si viene vacio), client_secret: actual}`
   (ya se desencriptan client_secret y refresh_token en esa funcion) y persistirlo.

Verificar blast radius del cambio de interfaz (mocks de `IIntegrationCore`) antes.

### IMPORTANTE

`handlers/oauth_store.go`: el store de intercambio OAuth (state + token exchange)
es en memoria (`map` + `sync.Mutex`). No es multi-instancia y se pierde en reinicio.
Mover a Redis (el core ya recibe `redis.IRedis`; falta cablearlo al modulo meli).

### DESEABLE

Hook `OnIntegrationCreated(IntegrationTypeMercadoLibre, ...)` (patron Shopify) para
validar token/seller_id tras conectar y precalentar config.

### NOTA (no es bug)

En `refresh_token.go` se usa `config["app_id"]` como `client_id` del grant. En
MercadoLibre el `client_id` OAuth ES el app id, asi que es correcto. No cambiar.

## Criterio de cierre

- Los items URGENTE e IMPORTANTE resueltos y verificados: la integracion mantiene
  sesion valida mas alla de 6h (varios ciclos de refresh) y el store OAuth persiste
  entre reinicios.
- Actualizar este archivo marcando cada item resuelto con fecha (no borrar hasta
  cerrar todos los urgentes/importantes).
