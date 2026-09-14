# Estado del loop: autorizacion en el backend

Plan: `.claude/plan/autorizacion-backend.md`
Rama: `feat/autorizacion-backend` (sin push)
Entorno: local, BD `127.0.0.1:5434`, backend :3050, front :3000

## Decisiones tomadas por defecto (recomendadas en el plan, sin respuesta del usuario)

1. Los modulos de un negocio los manda el plan de suscripcion (+ overrides vigentes);
   `business_resource_configured` solo restringe.
2. Negocio sin recursos configurados: recibe lo que diga su plan, no todo.
3. Rol Administrador: no edita roles/permisos globales (sigue RequireSuperAdmin).
4. Acciones: read/create/update/delete + especificas.
5. Auditoria: en local se corre en modo audit, se revisan denegaciones y luego enforce.

## Pendientes fuera del loop

- Crear el ticket en produccion (el loop no escribe en prod). Tickets locales no sirven.
- Push y PR: los decide el usuario.

## Fases

| Fase | Estado | Commit |
|---|---|---|
| 0 - Huecos explotables | hecha 2026-09-14 | 2cc1bf19 |
| 1 - Catalogo y motor | hecha 2026-09-14 | ver abajo |
| 2 - Enforcement audit | pendiente | |
| 3 - Enforce | pendiente | |
| 4 - Front y app obedecen | pendiente | |
| 5 - Limpieza y docs | pendiente | |
| Pruebas E2E multi-rol | pendiente | |

## Fase 0 - evidencia (local, usuario demo negocio 26)

- Orden sin negocio `779e57d7`: raw, history, status, update, delete -> 403.
- Orden propia `003bcef0` history -> 200.
- `/ai/recommendation` y `/notify/sse/order-notify` sin sesion -> 401.
- SSE del demo pidiendo `business_id=30` -> conecta al 26.
- `pause-ai` / `resume-ai` con body `business_id=30` -> usa el negocio del token.
- `go test` orders, whatsapp, events, ai: OK. `flutter analyze sse_client.dart`: OK.
- Pendiente de verificar en navegador: SSE con `withCredentials` (se hace en las pruebas E2E).

## Fase 1 - decisiones y evidencia

Desvio del plan: los codigos estables viven en `shared/authz/catalog.go` y se
mapean desde el nombre actual del recurso en BD (`LegacyNames`), sin migracion.
Motivo: `Migrate()` de `back/migration` tiene migraciones ajenas encoladas y
correrlo aplicaria cosas que no son de este trabajo. Un test obliga a que los 34
recursos de BD tengan codigo. Si un super admin renombra un recurso, el motor lo
ignora con un Warn en el log.

Cache: Redis `authz:access:v1:<user>:<business>`, TTL 60 s (sin invalidacion
explicita: un cambio de rol se ve en maximo un minuto).

`/auth/me/access` en local:
- Sin sesion -> 401.
- demo@probability.com (Administrador, negocio 26, plan basico): modulos
  iam, orders, shipments, invoicing, customers, wallet, integrations. 37 permisos.
  Navegacion: home, orders, products, shipments, customers, invoicing, wallet,
  integrations, users, website_config, subscription, profile.
- Super admin: 28 entradas de navegacion; con `?business_id=26` devuelve el negocio.
- Cache verificada en Redis con TTL.

Cambio de comportamiento a revisar en la fase de auditoria: con el plan basico,
Demo deja de ver Notificaciones, Bodegas, Inventario y Ultima milla (hoy el
sidebar los muestra solo por permiso). Es la decision 1 (manda el plan).

Hallazgo: el rol Administrador NO tiene `Productos` create/update/delete ni
`Ordenes` create/delete en BD. Con enforce, un administrador no podria crear
ordenes. La auditoria de la fase 2 lo tiene que resolver antes de bloquear.

Tambien: `business_module_overrides` ahora respeta `expires_at`.

## Fase 2 - decisiones y evidencia

Desvio del plan: en vez de un helper por archivo de rutas (829 rutas en ~80
modulos), la politica vive en una tabla central `shared/authz/route_policies.go`
(exactas + prefijos, accion por metodo HTTP) y un middleware global en
`/api/v1` la aplica con `c.FullPath()`. Al arrancar, el backend reporta las
rutas sin politica (`[authz] rutas sin politica declarada`); con
`AUTHZ_ROUTES_DUMP=<archivo>` escribe el listado completo. Hoy: 0 sin declarar
(78 public, 94 authenticated, 166 super_admin, 491 permission).

No hay test de CI que recorra `engine.Routes()`: construir el router exige BD,
Redis y RabbitMQ. La cobertura se verifica al arrancar (log) y con el volcado.

Modo por variables de entorno leidas con `os.Getenv` (no estan en el struct de
`shared/env`): `AUTHZ_MODE=audit|enforce`, `AUTHZ_ENFORCE_MODULES=orders,...`.

Auditoria en local (admin demo navegando orders, shipments, invoicing,
customers, integrations, wallet, users y 19 pantallas por iframe):
- `/businesses/:id/configured-resources` y `/businesses/simple` las pide todo
  usuario al cargar el layout -> pasan a `authenticated`. El handler de
  configured-resources ya valida propiedad; `/businesses/simple` NO la validaba
  (devolvia todos los negocios a cualquier JWT): ahora un usuario de negocio
  solo recibe el suyo. Verificado: demo -> [26], super admin -> [30, 26].
- `/storefront/catalog` denegado: correcto, Demo no tiene Storefront.

Analisis de permisos del Administrador contra las 491 rutas por permiso:
los modulos fuera del plan basico (inventory, warehouses, delivery,
notifications, storefront) se niegan por plan, como se decidio. Huecos reales
del rol: `orders.create/delete`, `products.create/update/delete`,
`integrations.update/delete`. Corregido POR API en local:
- `POST /permissions/bulk` -> permisos 129-134 (Ordenes Delete, Productos
  Create/Update/Delete, Integraciones Update/Delete).
- `POST /roles/4/permissions` con [128..134].

Tropiezo: `POST /roles/:id/permissions` REEMPLAZA los permisos del rol, no
los agrega. La primera llamada con [128..134] dejo al Administrador local con 7
permisos. Se restauro por el mismo API con la lista completa (1-21, 38-85,
128-134 = 76) y se verifico en BD. En produccion hay que mandar SIEMPRE la lista
completa.

PENDIENTE PARA PRODUCCION: esos 6 permisos y la asignacion al rol
Administrador no existen en prod. Hay que crearlos con el mismo API (o un seed)
ANTES de activar enforce, o los administradores no podran crear ordenes.
