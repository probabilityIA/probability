# Product Tour - plan de implementacion

Fecha: 2026-08-28
Decisiones tomadas: tour **hibrido** (modal conceptual + spotlight anclado),
persistencia en **backend por usuario**, componente **propio** (sin libreria nueva).

## 0. Punto de partida

Ya existe `front/central/src/services/modules/products/ui/components/ProductTour.tsx`
(266 lineas): modal de 7 pasos, navegacion con teclado, progreso, `localStorage`
(`products_tour_seen_v1`). Es el patron a generalizar, no a duplicar 20 veces.

Modulos a cubrir (por subnavbar / familia de rutas):

| Grupo | Rutas |
|---|---|
| home | `/home` |
| orders | `/orders` |
| products | `/products` |
| inventory | `/inventory` + audit, kardex, lpn, mobile, movements, operations, traceability |
| warehouses | `/warehouses`, `/warehouses/[id]` |
| shipments | `/shipments`, `/quotes`, `/generate`, `/cod`, `/shipping-margins` |
| delivery | drivers, geozones, routes, vehicles |
| wallet | `/wallet`, finanzas, saldos |
| accounting | `/accounting` + configuracion, facturas, movimientos |
| invoicing | configs, invoices, providers |
| customers | `/customers` |
| integrations | `/integrations` |
| storefront | catalogo, nuevo-pedido, pedidos |
| notifications | channels, config, event-types |
| iam | users, roles, permissions, resources |
| admin | businesses, subscription, announcements, tickets, marketing-leads, commercial, siigo-referrals, website-config |

~20 tours. No se hace uno por ruta: uno por **modulo**, con pasos que pueden
navegar entre subrutas del mismo modulo.

## 1. Backend - modulo `tours`

Modulo Go nuevo siguiendo `.claude/rules/architecture.md` (hexagonal completo).

### Modelo (en `back/migration/shared/models/`)

```
UserTourProgress
  id
  user_id      uint   idx
  business_id  uint   idx        (0 para super admin sin negocio)
  tour_key     string idx        'orders', 'products', ...
  version      int                version del contenido vista
  status       string             'pending' | 'in_progress' | 'completed' | 'skipped'
  step_index   int                ultimo paso alcanzado (drop-off)
  completed_at *time.Time
  created_at / updated_at / deleted_at
  UNIQUE (user_id, business_id, tour_key)
```

Migracion idempotente en `back/migration`, nomenclatura `XXX_user_tour_progress.go`.
Recordar: el deploy NO corre migraciones, se ejecuta el DDL a mano contra RDS.

### Endpoints (`/api/v1/tours`)

| Metodo | Ruta | Uso |
|---|---|---|
| GET | `/tours/progress` | todo el progreso del usuario actual (una llamada al entrar a la app) |
| PUT | `/tours/progress` | upsert `{tour_key, version, status, step_index}` |
| DELETE | `/tours/progress/:tour_key` | reset de un tour |
| POST | `/tours/progress/reset` | reset de todos (boton "ver tutoriales de nuevo") |

Reglas obligatorias:
- `business_id` se resuelve con el patron de `.claude/rules/multi-tenant-security.md`
  (token para usuario normal; query param obligatorio para super admin).
- `user_id` SIEMPRE del token, nunca del body.
- GET devuelve lista corta (< 50 filas por usuario): exento de paginado por catalogo.

Sin colas, sin eventos. Es escritura barata y sincrona.

## 2. Frontend - arquitectura

```
front/central/src/services/modules/tours/
  domain/
    types.ts        TourDefinition, TourStep, TourProgress
    ports.ts
  app/
    use-cases.ts
  infra/
    repository/tours.repository.ts
    actions/tours.actions.ts        ('use server' para el upsert)
  ui/
    components/
      TourProvider.tsx      contexto global: registry + progreso + estado activo
      TourModal.tsx         paso conceptual (generalizacion del ProductTour actual)
      TourSpotlight.tsx     overlay con recorte + popover anclado
      TourLauncher.tsx      boton "?" en el navbar
      TourStepFrame.tsx     chrome comun: progreso, prev/next, saltar, teclado
    hooks/
      useTour.ts            start(key), next(), prev(), skip(), complete()
      useTourAnchor.ts      helper para registrar/resolver targets
  content/
    index.ts       TOUR_REGISTRY: Record<TourKey, TourDefinition>
    orders.ts  products.ts  shipments.ts  wallet.ts  ...
```

`TourProvider` se monta en `app/(auth)/layout.tsx`, **dentro** de
`PermissionsProvider` y `SidebarProvider` (necesita permisos y el estado del
sidebar), envolviendo a `LayoutContent`.

### Tipos

```ts
type TourStep =
  | { kind: 'concept'; id; icon; title; subtitle; whatIs; whenToUse; example; highlight? }
  | { kind: 'spotlight'; id; title; body; target: string;
      placement?: 'top'|'bottom'|'left'|'right';
      route?: string;            // navega antes de mostrar el paso
      optional?: boolean;        // si el target no existe, se salta en silencio
      waitFor?: string }         // selector a esperar (modal que se abre)

interface TourDefinition {
  key: string;                   // 'orders'
  version: number;               // subirlo re-muestra el tour a todos
  routes: string[];              // rutas donde autoarranca / donde el "?" lo ofrece
  title: string;
  resource?: string;             // gate por permisos (permissions-context)
  autoStart: boolean;
  steps: TourStep[];
}
```

### Anclaje

Atributo `data-tour="<modulo>.<elemento>"` en la UI real:
`data-tour="orders.create"`, `data-tour="orders.filters"`, `data-tour="orders.row-actions"`.

Resolucion via `document.querySelector('[data-tour="..."]')` +
`getBoundingClientRect()`, recalculado en `resize` y `scroll`. Si el target no
aparece en 1.5 s, el paso se salta (nunca bloquear al usuario).

### Arranque automatico

Al montar una ruta: si existe un tour cuyo `routes` la cubre, `autoStart` es true,
el usuario tiene permiso sobre `resource`, y no hay progreso `completed`/`skipped`
con `version >= def.version`, abrir tras ~600 ms. Nunca dos tours simultaneos;
nunca si hay un modal abierto.

### Reentrada manual

- Boton "?" en el navbar: abre el tour de la ruta actual (siempre disponible).
- `/profile` -> seccion "Tutoriales": lista de tours con su estado y boton de reset
  individual y global.

## 3. Fases

**Fase 0 - fundaciones (bloqueante). HECHA 2026-08-28.**
Modelo + migracion + modulo Go + endpoints; modulo `tours` en el front con
provider, modal, spotlight y launcher; registry vacio. Migrar `products` al nuevo
sistema como piloto y borrar `ProductTour.tsx` (incluida la clave de
`localStorage`, con fallback: si existe `products_tour_seen_v1`, sembrar el
progreso como `completed` en el primer arranque).

**Fase 1 - modulos de mayor uso.**
home, orders, products (ya en fase 0), shipments (quotes + generate), wallet,
integrations.

**Fase 2 - operacion.**
inventory, warehouses, customers, invoicing, accounting, delivery, storefront.

**Fase 3 - administracion.**
iam (users/roles/permissions/resources), notifications, businesses, subscription,
announcements, tickets, marketing-leads, commercial, siigo-referrals,
website-config.

Cada tour de fase 1-3 es: escribir el contenido, sembrar los `data-tour` en los
componentes de ese modulo, y probar en local.

### Estado de la Fase 0

Backend: modelo `UserTourProgress`, `migrateUserTourProgress` (corrida en local,
pendiente en produccion), modulo `services/modules/tours` con los 4 endpoints,
registrado en `services/modules/bundle.go`. Verificado por curl: GET vacio, PUT
in_progress, PUT completed (upsert por unico), status invalido -> 400, DELETE,
sin token -> 401.

Frontend: `services/modules/tours/` con `TourProvider`, `TourRunner`, `TourModal`,
`TourSpotlight`, `TourLauncher` (boton flotante abajo a la derecha),
`TourSettings` (en `/profile`) y `use-target-rect`. Provider montado en
`app/(auth)/layout-content.tsx` dentro de `SelectedBusinessProvider`.

Piloto products: `ProductTour.tsx` eliminado, sus 7 pasos conceptuales movidos a
`content/products.ts` mas 4 pasos de spotlight anclados a `products.tabs`,
`products.search`, `products.integration-filter` y `products.create`. El
`localStorage` viejo (`products_tour_seen_v1`) se lee una vez y se siembra como
`completed` en backend antes de borrarse.

Pendiente antes de fase 1: correr la migracion en produccion.

## 4. Riesgos y decisiones ya tomadas

- **Estados vacios.** Un spotlight sobre "la primera fila de la tabla" no existe
  en una cuenta nueva, que es justo quien ve el tour. Todo paso anclado a datos
  va con `optional: true`; los pasos obligatorios se anclan a chrome fijo
  (botones, filtros, tabs, subnavbar).
- **Pasos dentro de modales.** Usar `waitFor` con MutationObserver y timeout. No
  se automatiza el clic del usuario: el paso pide la accion y espera.
- **z-index.** Los modales del repo usan `z-50` / `z-[60]`. El overlay del tour va
  en `z-[100]` y su popover en `z-[101]`.
- **Movil / sidebar colapsado.** Bajo 768 px el tour cae a modo solo-modal: el
  spotlight sobre elementos ocultos no aporta.
- **Permisos.** Si el usuario no tiene el recurso, el tour no se ofrece; si le
  falta un boton concreto, ese paso se salta.
- **Dark mode.** Usar tokens del theme, nunca colores fijos (el ProductTour actual
  ya usa `var(--color-primary)`; mantenerlo).
- **Multi-tenant.** El progreso es por `(user_id, business_id)`: un usuario con
  varios negocios ve el tour una vez por negocio. Deliberado, porque la
  configuracion cambia entre negocios.
- **Sin comentarios en TS.** Regla del repo, aplica a todo lo nuevo.
- **UTF-8.** Los archivos de contenido de tours van a superar las 500 lineas si se
  agrupan: mantener **un archivo por modulo** y sin acentos ni emojis en archivos
  largos. Los textos con acentos van con escapes `\u00XX` o el archivo se parte.
  El `ProductTour.tsx` actual (266 lineas, con emojis y acentos) esta al limite:
  al migrarlo, separar contenido de componente.

## 5. Metricas

Con el progreso en BD sale gratis:
- % de usuarios que completan cada tour,
- drop-off por `step_index` (donde se aburren),
- tours nunca abiertos (contenido inutil o launcher poco visible).

Consulta base:
```sql
SELECT tour_key, status, count(*)
FROM user_tour_progress
WHERE deleted_at IS NULL
GROUP BY 1,2 ORDER BY 1;
```

## 6. Criterio de listo

- [ ] Endpoint de progreso probado contra la copia local (business 26).
- [ ] Un tour completo (products) funcionando con modal + spotlight + persistencia.
- [ ] El "?" del navbar abre el tour correcto en todas las rutas de fase 1.
- [ ] Reset desde `/profile` vuelve a mostrar el tour.
- [ ] Ningun tour bloquea la UI cuando el target no existe.
- [ ] Probado en dark mode y en ancho movil.
