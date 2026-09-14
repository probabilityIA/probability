# Plan: autorizacion 100% en el backend, el front solo obedece

Fecha: 2026-09-14
Origen: diseno del asistente de IA en la plataforma. Para que la IA solo ofrezca
lo que el usuario puede usar, primero hace falta que el backend sea la unica
fuente de verdad de los permisos. Hoy no lo es.
Alerta relacionada: `.claude/alerts/autorizacion-backend.md` (2026-07-17).

## 1. Diagnostico (verificado 2026-09-14)

### Backend

- El RBAC existe en BD (`resource` 34, `action` 13, `permission` 128, `role`,
  `role_permissions`, `business_staff`, `business_resource_configured`) pero
  **ninguna ruta lo consulta**. No existe `RequirePermission`.
- 952 rutas. ~85 publicas/webhooks, 45 con `RequireSuperAdmin`, **~822 solo con
  `JWT()`** (autenticado = autorizado a todo). Solo `invoicing` e `inventory`
  (119 rutas) pasan por `RequireModuleAccess`.
- JWT se aplica por modulo, a veces por grupo y a veces por ruta: no hay un
  punto unico donde poner un control.
- `RequireRole` / `RequireAnyRole` rotos (leen `user_roles`, que el camino JWT
  nunca escribe) y sin uso.
- `GET /auth/roles-permissions`: **fail-open** si el negocio no tiene filas en
  `business_resource_configured` (pasan todos los permisos del rol). No cruza
  con el plan de suscripcion.
- Dos catalogos desconectados: recursos en BD con nombres en espanol
  (`Ordenes`, `Facturacion`) y codigos de modulo en Go
  (`shared/moduleregistry/modules.go`: `orders`, `invoicing`...).
- `HasModuleAccess` no revisa `ExpiresAt` de los overrides.
- Sin cache: `my-modules` hace varias consultas por modulo en cada llamada.

Huecos concretos todavia abiertos:

| Hueco | Donde |
|---|---|
| `business_id` tomado del body | whatsapp `send-template`, `conversations/:id/reply`, `pause-ai`, `resume-ai` |
| `:id` sin scoping por negocio | orders `update`, `delete`, `raw`, `history`, `change-status` |
| Rutas sin auth (verificar) | `/notify/sse/order-notify[/:businessID]`, `/ai/recommendation` |
| Fail-open de permisos | `login/internal/app/user-roles-permissions.go:187-211` |

### Front (`front/central`)

- `proxy.ts` solo valida que exista cookie de sesion. **Ninguna ruta tiene guard
  de permisos**: escribir `/invoicing` en la URL abre la pagina.
- Permisos en `localStorage` como JSON plano (`is_super` incluido), editables
  desde DevTools. Si falla el refetch, queda el valor editado.
- `sidebar.tsx` decide con 29 `hasPermission` escritos a mano, nombres en
  espanol e ingles, y comparaciones por nombre de rol (`'Administrador'`,
  `'demo'`, `'cliente_final'`, `includes('admin')`).
- `isSuperAdmin` en 147 archivos (616 usos), 10 subnavbars con logica propia,
  `orders-sidebar` muestra todo mientras cargan los permisos (falla abierto).
- `hasRouteAccess`, `PermissionGate`, `getResourceActions`: escritos y sin uso.

### App movil (`mobile/mobile_central`)

- Pide `roles-permissions` pero `hasPermission()` no se llama nunca. No pide
  `my-modules`. El router solo bloquea modulos `superAdminOnly` (Integraciones).

## 2. Principios del diseno

1. **Una sola decision, en el backend.** El front y la app no deciden permisos:
   pintan lo que el backend dice y el backend rechaza lo que no corresponde.
   Ocultar un boton es UX, no seguridad.
2. **Denegar por defecto.** Una ruta sin politica declarada se rechaza. Un
   recurso no configurado no se concede. Un error consultando permisos niega.
3. **Acceso efectivo = interseccion**, calculado en un solo lugar:
   `permisos del rol ∩ modulos del plan (+ overrides vigentes) ∩ recursos activos
   del negocio ∩ estado de suscripcion`. Super admin: todo, pero siempre con un
   `business_id` elegido (regla multi-tenant vigente).
4. **Codigos estables, no nombres.** `orders.read`, no `Ordenes`/`Orders`. El
   nombre visible es un label.
5. **Autorizar no reemplaza validar propiedad.** El permiso dice "puede leer
   ordenes"; el scoping por `business_id` dice "de su negocio". Ambos siempre.
6. **Cambios sin dejar a nadie afuera**: modo auditoria antes de bloquear.

## 3. Arquitectura propuesta

### 3.1 Catalogo unico (Go, versionado en el repo)

`shared/authz/catalog.go`: la lista de modulos, recursos y acciones vive en
codigo y se sincroniza a BD con una migracion idempotente.

```
module   orders      label "Ordenes"   route "/orders"   plan_code "orders"
  resource orders            actions read create update delete
  resource shipments         actions read create update
module   invoicing   label "Facturacion" route "/invoicing/invoices" plan_code "invoicing"
  resource invoices          actions read create update delete
```

- Migracion: `resource` gana `code` (unico), `module_code`, `label`;
  `action` gana `code`. Se mapean los 34 recursos actuales a codigos y se
  consolidan duplicados espanol/ingles. Los `permission` existentes se conservan
  (mismo `resource_id`), asi que los roles no pierden nada.
- El catalogo tambien alimenta la navegacion (ruta, label, descripcion, icono) y
  al asistente de IA.

### 3.2 Motor de acceso efectivo

`services/auth/authz` (modulo hexagonal):

- `GetEffectiveAccess(ctx, userID, businessID) -> Access`
  (`is_super`, `business`, `role`, `modules[]`, `permissions set`,
  `subscription_status`).
- `Can(access, "orders", "update") bool`.
- Cache en Redis `authz:access:<user>:<business>`, TTL 5 min, e invalidacion
  explicita al cambiar rol, permisos del rol, recursos del negocio, plan u
  overrides (publicar evento o borrar la llave desde esos casos de uso).
- Corrige de paso: fail-open de `business_resource_configured`, `ExpiresAt` de
  overrides.

### 3.3 Enforcement en las rutas

Dos piezas:

**a) Declarar la politica junto a la ruta.** Helper que registra la ruta en gin
y su politica en un registro central:

```go
authz.GET(g, "/orders/:id", authz.Perm("orders", "read"), h.GetOrder)
authz.POST(g, "/webhook", authz.Public(authz.Webhook), h.Receive)
authz.GET(g, "/tickets", authz.SuperAdmin(), h.List)
```

Politicas: `Public` (webhook, oauth, publico, health), `Authenticated` (solo
sesion, p.ej. `/auth/me/access`, perfil), `Perm(resource, action)`,
`SuperAdmin`, `StorefrontCustomer` (clientes finales, audiencia aparte).

**b) Un middleware global en `/api/v1`** que busca la politica por
`c.FullPath()`, autentica si hace falta, calcula el acceso (cacheado) y decide.
Ruta sin politica registrada = 403 y log `authz_undeclared_route`.

**c) Test de cobertura en CI**: recorre `engine.Routes()` y falla si alguna ruta
no tiene politica. Asi una ruta nueva no puede nacer sin decidir quien la usa.

Modo por variable `AUTHZ_MODE`:

- `audit`: evalua todo, **no bloquea**, loguea cada denegacion con usuario,
  negocio, rol, ruta y permiso faltante.
- `enforce`: bloquea (403 `{"success":false,"error":"forbidden","permission":"orders.update"}`).
- Se puede pasar a `enforce` por modulo (`AUTHZ_ENFORCE_MODULES=orders,invoicing`).

### 3.4 Contrato para el front y la app

`GET /api/v1/auth/me/access` (reemplaza `roles-permissions` + `my-modules`):

```json
{
  "is_super": false,
  "business": {"id": 26, "name": "Demo"},
  "role": {"id": 4, "code": "administrador", "label": "Administrador"},
  "subscription": {"status": "active", "end_date": "..."},
  "navigation": [
    {"module": "orders", "label": "Ordenes", "route": "/orders", "icon": "orders",
     "children": [{"resource": "shipments", "label": "Envios", "route": "/shipments"}]}
  ],
  "permissions": ["orders.read", "orders.create", "shipments.read"]
}
```

La navegacion llega **ya filtrada**: el front no sabe que existe un modulo que
el usuario no puede ver.

### 3.5 Front: obedecer

- `AccessProvider` hidratado desde un Server Component en `app/(auth)/layout.tsx`
  (fetch server-side con la cookie). Nada en `localStorage`.
- Guard de ruta server-side: el layout compara la ruta pedida contra
  `navigation`; si no esta, `notFound()`/redirect. Una sola regla para las ~38
  secciones.
- `sidebar.tsx`, `orders-sidebar.tsx` y los 10 subnavbars se pintan desde
  `navigation`. Se borran los `canViewX`, `RESOURCE_ROUTE_MAP` y las
  comparaciones por nombre de rol.
- Botones: `can('orders.update')` sobre la lista `permissions`. Es UX; el
  backend igual lo exige.
- `isSuperAdmin` se queda solo para el selector de negocio, no para conceder.
- Borrar el fail-open de `LoginForm.tsx` (`is_super: true` en el catch).
- 403 del backend = toast unico "No tienes permiso para esta accion".

### 3.6 App movil

Mismo endpoint. `AppModules` deja de declarar recursos propios: los modulos
visibles salen de `navigation`; el router de GoRouter redirige lo que no este.

## 4. Fases

### Fase 0 - Huecos explotables hoy (antes que todo, independiente)

- [ ] WhatsApp: `business_id` del token con `resolveBusinessID` en las 4 rutas.
- [ ] Orders: scoping por negocio en update, delete, raw, history, change-status.
- [ ] Verificar y cerrar auth de `/notify/sse/order-notify` y `/ai/recommendation`.
- [ ] Quitar fail-open de `LoginForm.tsx`.
- [ ] Tests por cada IDOR (usuario de negocio A pide recurso de B -> 404/403).

Criterio: curl con JWT de demo contra recursos de otro negocio devuelve 403/404.

### Fase 1 - Catalogo y motor (sin cambiar comportamiento)

- [ ] `shared/authz/catalog.go` + migracion de codigos en `resource`/`action`.
- [ ] Mapeo de los 34 recursos y consolidacion de duplicados; verificar en local
      que ningun rol pierde permisos (consulta antes/despues).
- [ ] Modulo `services/auth/authz`: acceso efectivo, cache, invalidacion.
- [ ] Fail-open y `ExpiresAt` corregidos dentro del motor.
- [ ] `GET /auth/me/access`. `roles-permissions` y `my-modules` quedan como
      alias temporales.
- [ ] Tests de la matriz: super admin, Administrador, demo, cliente_final,
      negocio sin plan, suscripcion vencida, override vencido.

### Fase 2 - Enforcement en modo auditoria

- [ ] Helpers `authz.GET/POST/...` y middleware global.
- [ ] Declarar politica en las 952 rutas, modulo por modulo (commits por modulo).
      Default de accion por metodo: GET read, POST create, PUT/PATCH update,
      DELETE delete; excepciones explicitas.
- [ ] Test de cobertura de rutas en CI.
- [ ] Desplegar con `AUTHZ_MODE=audit`. Revisar denegaciones reales una semana:
      cada una es un permiso faltante en un rol o una politica mal declarada.
- [ ] Ajustar roles con datos, no a ojo.

### Fase 3 - Enforce

- [ ] Encender por modulo: orders -> products -> shipments -> invoicing ->
      inventory -> resto. Un modulo por despliegue.
- [ ] Retirar `RequireModuleAccess` suelto (queda dentro del motor).
- [ ] `AUTHZ_MODE=enforce` global.

### Fase 4 - Front y app obedecen

- [ ] `AccessProvider` server-side + guard en layout.
- [ ] Sidebar y subnavbars desde `navigation`.
- [ ] Reemplazar `hasPermission`/comparaciones de rol por `can()`.
- [ ] Sacar permisos de `localStorage`.
- [ ] App movil sobre `/auth/me/access`.

### Fase 5 - Limpieza y documentacion

- [ ] Borrar `RequireRole`, `RequireAnyRole`, `RESOURCE_ROUTE_MAP`,
      `PermissionGate` sin uso, endpoints alias.
- [ ] Regla `.claude/rules/autorizacion.md` (como declarar la politica de una
      ruta nueva, como agregar un modulo al catalogo).
- [ ] Cerrar `.claude/alerts/autorizacion-backend.md`.

Despues de la fase 1, el asistente de IA ya puede construirse sobre
`GetEffectiveAccess` y `navigation`, sin esperar las fases 2-5.

## 5. Riesgos

| Riesgo | Mitigacion |
|---|---|
| Bloquear a usuarios que hoy trabajan bien | modo audit + enforce por modulo |
| Rol Administrador con permisos incompletos en BD | la semana de audit lo revela con datos |
| Rendimiento (952 rutas consultando permisos) | cache Redis por usuario/negocio |
| Cambio de rol que no se refleja | invalidacion explicita + TTL corto |
| Webhooks rotos por quedar sin politica publica | test de cobertura + lista de ~85 publicas ya inventariada |
| Deploy blue-green con dos versiones | el endpoint nuevo convive con los alias |

## 6. Decisiones pendientes del usuario

1. **Disponibilidad de modulos**: la manda el plan de suscripcion (+ overrides)
   y `business_resource_configured` solo restringe, o se mantiene
   `business_resource_configured` como fuente. Recomendado: el plan manda.
2. **Negocio sin recursos configurados**: hoy recibe todo. Recomendado: aplicar
   lo que diga su plan.
3. **Rol Administrador y tablas globales** de roles/permisos: roles por negocio
   (el admin crea roles para su equipo con permisos que el tiene) o solo super
   admin edita roles. Pendiente desde julio.
4. **Acciones**: hoy hay 13 (`Manage`, `Approve`, `Migrate`...). Recomendado:
   `read/create/update/delete` + las especificas que un modulo necesite de
   verdad (p.ej. `invoices.cancel`, `shipments.generate_guide`).
5. **Duracion del modo auditoria** en produccion. Recomendado: una semana.
