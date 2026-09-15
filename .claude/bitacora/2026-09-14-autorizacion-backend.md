# Autorizacion movida al backend: huecos, motor y tropiezos

**Ticket:** pendiente de crear en produccion (el trabajo se hizo en un loop
local sin escritura en prod).
**Rama:** `feat/autorizacion-backend` (sin push). Commits 2cc1bf19, d82903b5,
103e9c44, b951cc27, 986811b2 y el de la fase 5.
**Plan:** `.claude/plan/autorizacion-backend.md`. Estado y evidencia por fase:
`.claude/plan/goal-autorizacion-backend.md`. Regla: `.claude/rules/autorizacion.md`.

## Que se pidio

Un asistente de IA en la plataforma que ofrezca y navegue solo a lo que el
usuario puede usar. Al revisar como se deciden los permisos se encontro que el
backend no autorizaba nada: el front decidia con listas escritas a mano.

## Diagnostico (verificado en codigo y en BD local)

- 952 rutas (829 bajo `/api/v1`); ~822 solo con `JWT()`.
- Ningun `RequirePermission`. `RequireRole`/`RequireAnyRole` rotos y sin uso.
- `GET /auth/roles-permissions`: si el negocio no tenia filas en
  `business_resource_configured`, pasaban todos los permisos del rol.
- Front: permisos en `localStorage` (incluido `is_super`), 29 `hasPermission`
  en el sidebar con nombres en espanol e ingles, 52 paginas sin guard.
- App movil: `hasPermission()` nunca se llamaba.

Huecos explotables con un JWT de demo (cerrados en 2cc1bf19):

| Hueco | Evidencia local |
|---|---|
| SSE sin autenticacion: cualquiera escuchaba eventos de cualquier negocio | sin sesion ahora 401; demo pidiendo negocio 30 queda conectado al 26 |
| WhatsApp send-template/reply/pause/resume con `business_id` del body | ahora usa el del token |
| Orders raw/history/update/delete/status sin scoping | orden sin negocio `779e57d7` -> 403 en los 5 |
| `/businesses/simple` devolvia todos los negocios | demo -> [26], super admin -> [30, 26] |

## Rol Administrador incompleto en BD

Con el motor nuevo se vio que el rol Administrador (id 4) NO tenia
`Ordenes` Create/Delete, `Productos` Create/Update/Delete ni `Integraciones`
Update/Delete. Hoy funciona porque nadie valida; con enforce, un administrador
no podria crear ordenes. En local se crearon los permisos 129-134 y se
asignaron por API. **En produccion no se ha hecho.**

## Tropiezo: asignar permisos a un rol los reemplaza

`POST /roles/4/permissions` con `[128..134]` dejo al Administrador local con 7
permisos: el endpoint reemplaza la lista, no la agrega. Se restauro con la
lista completa (1-21, 38-85, 128-134 = 76) y se verifico en BD. Lo mismo va a
pasar en produccion si se manda solo lo nuevo.

## Tropiezo: la cookie de prueba se piso

Al hacer login por curl de un segundo usuario en el mismo archivo de cookies,
la matriz de acceso atribuyo al "administrador" los resultados del rol demo.
Cada usuario necesita su propio archivo de cookies.

## Hipotesis y decisiones descartadas

- **Declarar la politica con un helper en cada archivo de rutas**: 829 rutas en
  ~80 modulos. Se prefirio una tabla central con prefijos y un middleware
  global que usa `c.FullPath()`.
- **Test de CI que recorra el router real**: construir el router exige BD,
  Redis y RabbitMQ. La cobertura se reporta al arrancar.
- **Migracion para agregar `code` a `resource`**: `Migrate()` tenia migraciones
  ajenas encoladas. Los codigos se mapean desde el nombre en Go.
- **Retirar `RequireModuleAccess`**: en produccion el modo por defecto es audit,
  y sin ese middleware facturacion e inventario quedarian abiertos por plan.

## Verificacion final

Ver seccion "Pruebas E2E" de `goal-autorizacion-backend.md`.
