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
