# Alerta: el backend no valida permisos (RBAC decorativo)

Fecha: 2026-07-17
Contexto: auditoria de roles/permisos pedida por el usuario.

## Hallazgo central

Existe un RBAC completo en base de datos (role, permission, resource, action,
role_permissions) con su CRUD y su UI, pero **ninguna ruta HTTP lo consulta**.
De ~235 rutas autenticadas, la autorizacion se reduce a:

- `middleware.JWT()` -> solo autenticacion (tienes un token valido?)
- `middleware.IsSuperAdmin(c)` -> chequeos manuales dentro de ~40 handlers
- `middleware.RequireSuperAdmin()` -> ahora en 16 rutas (antes 2)

No existe `RequirePermission()` / `RequireResource()`. El gating por permisos es
hoy 100% de UI: ocultar un boton no protege el endpoint.

## RESUELTO el 2026-07-17

- [x] IDOR `GET /integrations/:id` -> filtraba credenciales DESENCRIPTADAS de
      cualquier negocio a cualquier JWT. Era el peor: con el autoregistro demo,
      cualquiera en internet sacaba las llaves de Shopify/Meli/Woo/Siigo/Bold de
      todos los tenants. Fix: validacion de propiedad + test que verifica que
      `GetIntegrationByIDWithCredentials` ni se llama.
- [x] IDOR `DELETE /businesses/:id` -> `RequireSuperAdmin()`.
- [x] IDOR `PUT /businesses/:id` -> valida que el `:id` sea el business del token.
- [x] `POST /businesses`, `PUT /:id/activate|deactivate`,
      `PUT /configured-resources/:resource_id/activate|deactivate` -> `RequireSuperAdmin()`.
- [x] IDOR `GET /users/:id` y `DELETE /users/:id` -> validan pertenencia al
      business del token (cerraba un TODO que admitia el fallo por escrito).
- [x] Escrituras de `/roles` y `/permissions` -> `RequireSuperAdmin()`.
      Los GET quedan abiertos a cualquier JWT (las pantallas los necesitan).
- [x] `PUT /businesses/:id` destruia el negocio con un PUT parcial: el struct de
      request usaba tipos valor y el mapper hacia `&req.Name`, asi que ningun
      puntero era nil y el merge del usecase pisaba todo con vacios. Fix: punteros.

## URGENTE (abierto)

- [ ] **`POST /whatsapp/*` (4 rutas) toman `business_id` del BODY** sin mirar el
      token. Cualquier usuario manda WhatsApps con la identidad de otro negocio.
      Viola literal la regla de multi-tenant-security.md.
      `services/integrations/messaging/whatsapp/internal/infra/primary/handlers/routes.go:13-16`
- [ ] **Modulo `orders`: 6 endpoints `:id` sin scoping por business**
      (get, get-raw, history, update, delete, change-status). Cero referencias a
      businessID. Agravante: `orders/constructor.go:24` YA define `resolveBusinessID`
      y solo se usa en `upload-bulk.go:88`. Esta escrito y sin cablear.
      `GET /orders/:id/raw` expone el payload crudo de la integracion.
- [ ] **`POST /orders/:id/request-confirmation` y `/send-guide-notification`**:
      mismo vacio, con efecto externo (mandan WhatsApp al cliente final de otro
      negocio).
- [ ] **Fail-open en la resolucion de permisos**: `user-roles-permissions.go:183-211`
      -> si el business tiene CERO filas en `business_resource_configured`, pasan
      TODOS los permisos del rol sin filtrar. Esta comentado como intencional.
      Hay demos en prod en esa condicion.

## IMPORTANTE (abierto)

- [ ] **`GET /businesses/simple`** devuelve todos los negocios activos con sus
      colores a cualquier JWT, sin filtrar por el business del usuario.
- [ ] **`RequireRole()` / `RequireAnyRole()` estan ROTOS**: leen `user_roles` del
      contexto, que solo escribe `APIKeyMiddleware` (`middleware.go:137`), que no
      esta cableado a ninguna ruta. `AuthMiddleware` (path JWT, `middleware.go:26-82`)
      nunca lo escribe. Si alguien los monta, dan 403 a todos. Falla cerrado, pero
      contamina codigo vivo: la rama `isAdmin` es inalcanzable en
      `pay/.../wallet_admin.go:18` y `codreport/.../helpers.go:36`.
- [ ] **Permisos en localStorage como JSON plano**, `is_super` incluido, confiado
      verbatim (`permissions-context.tsx:79`). Forjable desde devtools. Solo
      desbloquea UI (el session_token es HttpOnly), pero es la unica linea de
      defensa mientras el back no autorice.
- [ ] **`LoginForm.tsx:97-110`**: si falla el fetch de permisos, el catch escribe
      `is_super: true`. Solo aplica si el login ya devolvio `is_super_admin`, asi
      que NO es escalada de un usuario normal, pero falla abierto. Quitar.
- [ ] **52 de 59 paginas del front sin guard de ruta.** Ironia: `hasRouteAccess`
      (`permissions-context.tsx:100`) y `PermissionGate` (`:183`) ya estan escritos
      y NO se usan. Cablearlos en `app/(auth)/layout.tsx` cubre las 52 de un saque.
      Ojo: `RESOURCE_ROUTE_MAP` tiene fail-open para rutas no mapeadas (`:112`).

## DESEABLE

- [ ] Matching por substring: `roleName.includes('admin')` en
      `CodReportView.tsx:26` y `user-profile-modal.tsx`. Un rol "Subadministrador"
      o "no-admin" pasa el gate.
- [ ] `orders-subnavbar.tsx:22-23` falla ABIERTO mientras cargan los permisos;
      `inventory-subnavbar.tsx:59` falla cerrado. Unificar en fallar cerrado.
- [ ] Gates de rol hardcodeados en el front (~10 sitios + 31 ConfigForms que leen
      `is_super` directo de localStorage). Reemplazables por permisos del back.

## Decision pendiente con el usuario

El rol **"Administrador" (id 4) es scope `business`**, lo tienen 21 usuarios en
19 negocios, y le da CRUD sobre las tablas GLOBALES de Roles/Permisos/Recursos.
O sea: un admin de un negocio puede crear permisos globales y asignarselos.
Se bloquearon las ESCRITURAS (RequireSuperAdmin); esos 21 usuarios pierden esa
capacidad. Falta decidir si es lo correcto o si hacen falta roles por negocio.

## Orden recomendado (contraintuitivo, importa)

1. Blindar RBAC admin (HECHO) — antes de darle valor al RBAC.
2. Arreglar `AuthMiddleware` para poblar `user_roles`, sin lo cual cualquier gate
   de roles es inoperable.
3. IDOR restantes (whatsapp, orders).
4. Cablear el `resolveBusinessID` que ya existe en orders.
5. Recien ahi diseniar `RequirePermission()`.

Si se implementa `RequirePermission()` SIN blindar antes los endpoints de RBAC,
se crea una escalada de privilegios peor que la actual.

## Criterio para cerrar

Cuando no queden items URGENTE ni IMPORTANTE. Los IDOR de whatsapp y orders son
explotables hoy con un curl y un JWT de demo (que cualquiera obtiene en 30s por
el autoregistro).
