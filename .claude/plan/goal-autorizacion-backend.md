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
| 2 - Enforcement audit | hecha 2026-09-14 | 103e9c44 |
| 3 - Enforce | hecha 2026-09-14 (local) | ver abajo |
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

## Fase 3 - enforce en local

`back/central/.env` local: `AUTHZ_MODE=enforce` (en produccion NO se ha tocado
nada; alli la variable no existe y el default es `audit`).

Usuarios de prueba creados por API (super admin), negocio 26:
- 77 `authz.demo@probability-test.com`, rol `demo` (7).
- 78 `authz.envios@probability-test.com`, rol nuevo `Operador de envios` (8)
  con Ordenes Read + Envios Create/Read/Update (permisos 1, 6, 7, 9).
Las contrasenas las genero el API (quedan en el scratchpad de la sesion) y
ambos tienen `require_password_change = true`.

Matriz por API con enforce (codigo HTTP):

| endpoint | super | admin | rol demo | envios |
|---|---|---|---|---|
| GET /orders | 200 | 200 | 200 | 200 |
| POST /orders `{}` | 400 | 400 | 400 | 403 |
| GET /products | 200 | 200 | 200 | 403 |
| DELETE /products/999999 | 404 | 404 | 403 | 403 |
| GET /customers | 200 | 200 | 403 | 403 |
| GET /invoicing/invoices | 200 | 200 | 200 | 403 |
| GET /users | 200 | 200 | 403 | 403 |
| GET /integrations | 200 | 200 | 403 | 403 |
| GET /shipments | 200 | 200 | 200 | 200 |
| GET /inventory/lots | 200 | 403 (plan basico) | 403 | 403 |
| GET /tickets | 200 | 403 | 403 | 403 |
| GET /subscriptions/me | 200 | 200 | 200 | 200 |
| sin sesion GET /orders | 401 | | | |

400 y 404 significan que paso la autorizacion y fallo la validacion del handler.

Navegador como administrador con enforce: Ordenes, Clientes e Integraciones
cargan con datos y el log no registra ninguna denegacion para el usuario 8.

`RequireModuleAccess` de invoicing e inventory se deja: es redundante con el
motor (que ya cruza el plan) pero no estorba. Se retira en la fase 5.

## Fase 4 - front y app obedecen

Backend:
- Subrecursos de inventario heredan las acciones del padre si el negocio los
  tiene activos (`InheritFromParent`); los avanzados excluyen al rol `demo`
  (`ExcludeRoles`). Reemplaza la lista `DEMO_ALLOWED_RESOURCES` y el cruce con
  `useResourceConfig` que hacia el front.
- Navegacion: entradas `notification_channels` y `notification_event_types`.

Front (`front/central`):
- `services/auth/access`: tipos, repositorio y `getMyAccessAction` sobre
  `GET /auth/me/access`.
- `permissions-context.tsx` reescrito: expone `access`, `can(codigo)`,
  `hasNav(clave)`, `canAccessRoute`, `roleCode`. Se borraron `hasPermission`,
  `hasRouteAccess`, `getResourceActions`, `PermissionGate`, `useHasPermission`
  y `RESOURCE_ROUTE_MAP`.
- Permisos fuera de `localStorage`: `TokenStorage.getPermissions` devuelve una
  copia en memoria que llena el provider (los 55 archivos que la leen siguen
  funcionando) y borra la llave vieja.
- `AccessRouteGuard` en `app/(auth)/layout.tsx`: mientras carga muestra
  spinner; si la ruta no esta en la navegacion del backend muestra "No tienes
  acceso a este modulo". Es guard de cliente: la seguridad real es el backend.
- Sidebar, orders-sidebar, orders-subnavbar, inventory-subnavbar,
  storefront-subnavbar, pestanas de integraciones, tours y paginas de IAM
  consultan `hasNav`/`can`. Sin comparaciones por nombre de rol visible:
  quedan 5 por `roleCode` (`cliente_final` x4, `demo` en el cambio de tema),
  que es un codigo estable que manda el backend.
- LoginForm ya no pide `roles-permissions` ni escribe permisos.
- `tsc --noEmit`: 0 errores.

App movil:
- `LoginProvider` pide `/auth/me/access` al entrar y al restaurar sesion y
  expone `canNav`. `AppModules` cambia `resources` por `navKey`;
  `isRouteAllowed` y `visibleGroupsFor` filtran por la navegacion del backend.
- `flutter test` login + navegacion: 70 pruebas OK (4 nuevas).
