# Autorizacion - el backend decide, el front obedece

Regla OBLIGATORIA desde 2026-09-14. Quien puede ver o hacer algo lo decide un
solo lugar del backend. El front y la app movil pintan lo que el backend dice;
ocultar un boton es experiencia de usuario, no seguridad.

## Piezas

| Pieza | Donde | Para que |
|---|---|---|
| Catalogo | `back/central/shared/authz/catalog.go` | recursos, acciones y navegacion con codigos estables (`orders.update`) |
| Politica por ruta | `back/central/shared/authz/route_policies.go` | que exige cada ruta de `/api/v1` |
| Motor | `back/central/services/auth/authz` | acceso efectivo del usuario, cache 60 s, `GET /auth/me/access` |
| Middleware global | `authz/internal/infra/primary/middleware/guard.go` | aplica la politica en cada request |
| Front | `front/central/src/shared/contexts/permissions-context.tsx` | `can('orders.update')`, `hasNav('orders')`, `AccessRouteGuard` |
| App movil | `LoginProvider.canNav`, `AppModules.navKey` | modulos y rutas visibles |

## Acceso efectivo

```
permisos del rol
  ∩ recursos activos del negocio (business_resource_configured, si tiene filas)
  ∩ modulos del plan de suscripcion (+ overrides vigentes)
  + herencia de subrecursos (inventory.stock hereda de inventory)
  - recursos solo super admin y exclusiones por rol
```

El super admin (`business_id = 0` en el token) puede todo, pero sigue
necesitando elegir negocio segun `multi-tenant-security.md`.

Un permiso de BD se traduce a codigo por el nombre del recurso
(`LegacyNames`). Si un recurso de BD no tiene codigo, el motor lo ignora y deja
un Warn `[authz] recurso de BD sin codigo`. Un test obliga a que los recursos
conocidos esten mapeados.

## Al crear una ruta nueva

1. Registrar la ruta en el modulo como siempre.
2. Declarar su politica en `route_policies.go`:
   - `publicRoutes`: webhooks, OAuth callbacks, login, tienda publica.
   - `authenticatedRoutes`: solo sesion (catalogos compartidos, perfil).
   - `superAdminRoutes` o prefijo con `SuperAdmin()`.
   - Por permiso: prefijo `Perm("recurso")` (accion por metodo: GET read, POST
     create, PUT/PATCH update, DELETE delete) o `PermAction` si la accion no
     sigue el metodo.
3. Al arrancar, el backend loguea `[authz] rutas sin politica declarada`. Con
   `AUTHZ_ROUTES_DUMP=/ruta/archivo` escribe el listado completo. Debe quedar en 0.
4. En modo `enforce`, una ruta sin politica responde 403.

La politica no reemplaza validar propiedad: toda consulta por `:id` sigue
verificando que el recurso sea del `business_id` resuelto.

## Al crear un modulo nuevo

1. Agregar el recurso a `Resources` en `catalog.go` (codigo, modulo del plan,
   nombre en BD).
2. Si tiene pantalla, agregar la entrada a `Navigation` con su regla.
3. Crear los permisos en BD por API (`POST /permissions/bulk`) y asignarlos a
   los roles.
4. En el front, usar `hasNav('clave')` para el menu y `can('recurso.accion')`
   para botones. Nunca comparar nombres de rol ni nombres de recurso.

## Errores de acceso en el front (401, 402, 403)

El backend responde siempre `{ code, message, permission }` con un `message`
ya redactado para el usuario ("No tienes permiso de creación en Órdenes.").
El front lo convierte en un modal global (`AccessDeniedModal`, montado en el
`PermissionsProvider`) sin que cada pantalla tenga que manejarlo:

| Camino | Donde se detecta |
|---|---|
| Server Action (fetch en el servidor) | `src/instrumentation.ts` intercepta el fetch al API y deja la cookie `pb_access_denied`; el provider la lee cada segundo |
| fetch desde el navegador | el provider envuelve `window.fetch` |
| Pantalla que muestra el error con toast o `alert` | `ToastProvider` y el `alert` envuelto reconocen el mensaje y abren el modal en su lugar |

Comportamiento: 403 muestra "No tienes permiso" y recarga el acceso (menú y
botones se ajustan); 402 lleva a Suscripción; 401 pide iniciar sesión.

La cookie existe porque en producción Next.js borra el texto de los errores
lanzados desde Server Actions: sin ella, el usuario vería un error genérico.

Al construir una pantalla nueva no hace falta nada: basta con no tragarse el
error. Para disparar el modal a mano: `notifyAccessDenied({ code, message })`.

## Modos

| Variable | Valores | Efecto |
|---|---|---|
| `AUTHZ_MODE` | `audit` (default) / `enforce` | audit solo loguea `[authz] denegacion en auditoria` |
| `AUTHZ_ENFORCE_MODULES` | `orders,invoicing,...` | bloquea esos modulos aunque el modo sea audit |

Se leen con `os.Getenv`, no estan en el struct de `shared/env`.

## Trampas conocidas

- `POST /roles/:id/permissions` **reemplaza** todos los permisos del rol. Mandar
  siempre la lista completa.
- El cache de acceso dura 60 s: un cambio de rol se nota en maximo un minuto.
- `RequireModuleAccess` de invoicing e inventory sigue montado a proposito:
  mientras produccion este en `audit`, es lo unico que bloquea esos modulos por
  plan. Se retira cuando produccion pase a `enforce`.
- `GET /auth/roles-permissions` y `GET /subscriptions/my-modules` siguen vivos
  porque la app movil y la pantalla de suscripcion los usan. No construir nada
  nuevo sobre ellos.

## Prohibido

- Decidir permisos en el front o en la app (listas de recursos, `role === 'X'`).
- Guardar permisos en `localStorage`.
- Ruta nueva sin politica declarada.
- Tomar `business_id` del body o query para un usuario normal.
