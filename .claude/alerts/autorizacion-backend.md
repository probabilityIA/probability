# Alerta: autorizacion en el backend (RBAC decorativo)

Fecha: 2026-07-17. Actualizada: 2026-09-14.
Contexto: auditoria de roles/permisos pedida por el usuario. El 2026-09-14 se
implemento el plan `.claude/plan/autorizacion-backend.md` en la rama
`feat/autorizacion-backend` (sin push), probado solo en local.
Regla vigente: `.claude/rules/autorizacion.md`.

## Hallazgo central (original)

Existia un RBAC completo en base de datos pero ninguna ruta HTTP lo consultaba.
De ~235 rutas autenticadas, la autorizacion se reducia a `JWT()`,
`IsSuperAdmin(c)` y `RequireSuperAdmin()`. El gating por permisos era 100% de UI.

## RESUELTO el 2026-07-17

- [x] IDOR `GET /integrations/:id` (credenciales desencriptadas de cualquier negocio).
- [x] IDOR `DELETE /businesses/:id`, `PUT /businesses/:id`, `GET/DELETE /users/:id`.
- [x] Escrituras de `/roles`, `/permissions`, `/businesses` -> `RequireSuperAdmin()`.
- [x] `PUT /businesses/:id` destruia el negocio con un PUT parcial.

## RESUELTO el 2026-09-14 (local, rama feat/autorizacion-backend)

- [x] `POST /whatsapp/*` tomaban `business_id` del body -> usan el del token
      (send-template, reply, pause-ai, resume-ai). Commit 2cc1bf19.
- [x] Orders sin scoping (raw, history, update, delete, change-status) y
      `/orders/map` -> validan que la orden sea del negocio del token. 2cc1bf19.
- [x] SSE `/notify/sse/order-notify` sin autenticacion (cualquiera escuchaba
      eventos de cualquier negocio) y `/ai/recommendation` abierto. 2cc1bf19.
- [x] Fail-open de permisos con negocio sin recursos configurados: el motor
      nuevo cruza siempre con el plan. d82903b5.
- [x] `GET /businesses/simple` devolvia todos los negocios a cualquier JWT ->
      un usuario de negocio solo recibe el suyo. 103e9c44.
- [x] `RequireRole` / `RequireAnyRole` rotos -> eliminados (fase 5).
- [x] Permisos en localStorage forjables -> ya no se guardan. 986811b2.
- [x] `LoginForm` escribia `is_super: true` si fallaba el fetch. 2cc1bf19.
- [x] 52 paginas sin guard de ruta -> `AccessRouteGuard` en el layout. 986811b2.
- [x] Gates por nombre de rol en el front -> `can()`/`hasNav()`; quedan 5 por
      `roleCode` estable del backend. 986811b2.
- [x] `RequirePermission`: middleware global con politica por ruta para las
      829 rutas de `/api/v1`, modos audit/enforce. 103e9c44.

## URGENTE (abierto) - para llevar a produccion

- [x] **Permisos del rol Administrador en produccion** (2026-09-14): por API de
      produccion se crearon los permisos 129-134 (Ordenes Delete, Productos
      Create/Update/Delete, Integraciones Update/Delete) y se asignaron al rol 4
      con la lista completa (1-21, 38-85, 128-134). Verificado en BD de prod:
      76 permisos = 69 previos + 7 nuevos.
- [ ] Desplegar en `AUTHZ_MODE=audit` y revisar en produccion los logs
      `[authz] denegacion en auditoria` al menos una semana antes de enforce.
- [ ] Versiones viejas de la app movil pierden el SSE (ahora exige token).

## IMPORTANTE (abierto)

- [ ] Rol Administrador (scope business) y tablas globales de roles/permisos:
      decidir si hacen falta roles por negocio. Hoy solo el super admin escribe.
- [ ] Retirar `RequireModuleAccess` de invoicing/inventory cuando produccion
      este en enforce (hoy es lo unico que bloquea esos modulos por plan en audit).
- [ ] `/geocode`, `/address-search`, `/places-search` quedaron publicos como
      estaban: confirmar si deben exigir sesion (consumen cuota de Google).

## DESEABLE

- [ ] `roleCode === 'cliente_final'` en 4 lugares del front: reemplazar por un
      campo `audience` en `/auth/me/access`.
- [ ] Guard de rutas del lado del servidor (hoy es de cliente; el backend ya
      bloquea los datos).
- [ ] Test de CI que recorra `engine.Routes()` (hoy la cobertura se ve al arrancar).

## Criterio para cerrar

Cuando los items URGENTE esten hechos en produccion y enforce lleve una semana
sin denegaciones legitimas.
