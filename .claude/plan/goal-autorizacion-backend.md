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
| 0 - Huecos explotables | hecha 2026-09-14 | ver git log |
| 1 - Catalogo y motor | pendiente | |
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
