# Plan WMS - Modulo Inventory - Estado Actual

> Rama: `feat/inventory-wms-phases` | 6 commits + cambios pendientes sesion 2026-04-21 | Build + tests backend OK | TS frontend sin errores nuevos

---

## Sesion 2026-04-21: Frontend UI completo + UX polish + E2E duro

### Entregables

1. **Frontend - Todas las paginas construidas**
   - `/warehouses` refactorizada: tabla expandible con jerarquia inline (WarehouseTreeTable).
   - `/warehouses/[id]` con arbol visual (ya existente).
   - `/inventory/lots`, `/inventory/serials`, `/inventory/uoms` unificadas en `/inventory/traceability` (3 tabs).
   - `/inventory/operations/{putaway,replenishment,cross-dock}` unificadas en `/inventory/operations` (3 tabs).
   - `/inventory/audit/{plans,tasks,discrepancies}` unificadas en `/inventory/audit` (3 tabs).
   - `/inventory/kardex` reporte con running balance + export CSV.
   - `/inventory/analytics/slotting` ABC con cards A/B/C.
   - `/inventory/lpn` CRUD + move/merge/dissolve.
   - `/inventory/mobile` scan 6-way con historial.
   - `/inventory/sync/logs` bitacora read-only.

2. **Sub-navbar unificado y scrollable**
   - Antes: 4 items. Ahora: 12 items en scroll horizontal, agrupados en 4 secciones (Catalogo, Inventario, Operaciones, Capture).
   - Color del item activo + "Tus Integraciones" usan `var(--color-primary)` del negocio.
   - Boton Guia icon-only con animacion `tour-pulse` (wiggle + ring + dot rojo) cuando no se ha visto.

3. **Tours guiados**
   - `InventoryTour` 13 pasos: explica cada sub-modulo del WMS con ejemplos.
   - `WarehouseTour` 11 pasos: explica jerarquia Zona > Pasillo > Rack > Nivel > Posicion + flags + LPN.
   - Boton Guia aparece en subnavbar y en header de bodegas, pulsa hasta que se abre la primera vez (`localStorage` flag).

4. **Modal Ajustar/Transferir stock consciente del contexto**
   - Carga lotes activos del producto seleccionado y muestra selector solo si aplica.
   - Carga UoMs del producto y muestra selector si hay mas de una.
   - Carga ubicaciones de la bodega y permite asignar posicion especifica.
   - Carga estados de inventario (default: available).
   - Backend: `AdjustStockRequest` / `TransferStockRequest` ahora aceptan `lot_id`, `state_id`, `uom_id`, `location_id` opcionales.
   - Tx (`AdjustStockTx`, `TransferStockTx`) refactorizados para usar `getOrCreateLevelKeyTx` con clave `(product, warehouse, location, lot, state)`.
   - Toast de exito detalla `+N uds en ubicacion X lote Y`.

5. **Theming del negocio aplicado globalmente**
   - `globals.css` define: `.btn-business-primary`, `.btn-business-primary-soft`, `.subnav-active`, `.tour-pulse`, `.text-business-primary`, `.border-business-primary` usando `var(--color-primary)`.
   - Tabla standard (`.table th`) usa `var(--color-primary, #7c3aed)` como fallback.
   - Removidos todos los `bg-gradient-to-r from-[#7c3aed]` hardcoded de los modulos inventory/warehouses.
   - Titulos h1 redundantes eliminados de todas las paginas de inventario (el subnavbar ya dice donde estamos).

6. **Tabla bodegas expandible**
   - Reemplazo de `WarehouseList` + detail page por `WarehouseTreeTable`.
   - Cada bodega se expande inline y muestra el arbol completo (zona > pasillo > rack > nivel > posicion).
   - Botones inline por nodo (hover): `+` agregar hijo, editar, eliminar.
   - Pildoras de stats agregados por bodega: Z N_zones P N_pasillos R N_racks N N_niveles POS N_posiciones.
   - Arranca colapsada mostrando solo el nivel zona; usuario expande a demanda.

### Bugs de produccion detectados y fixeados

1. **`idx_inventory_product_warehouse` legacy unique index** bloqueaba niveles multiples por producto/bodega.
   - Error: `duplicate key value violates unique constraint "idx_inventory_product_warehouse"`.
   - Causa: indice UNIQUE viejo `(product_id, warehouse_id)` de Fase 0 quedo huerfano cuando Fase 2 agrego lot_id/state_id al unique key.
   - Fix: `migrate_inventory_traceability.go` ahora ejecuta `DROP INDEX IF EXISTS idx_inventory_product_warehouse` (idempotente).
   - Sin este fix el WMS era inutilizable en la practica (un producto no podia estar en mas de una posicion).

2. **Capture module handlers sin JSON tags**
   - LPN/Scan/SyncLog devolvian PascalCase desde entities → frontend esperaba snake_case.
   - Fix: nuevo archivo `handlers/response/capture.go` con DTOs completos y mappers (`LicensePlateFromEntity`, `ScanResolutionFromEntity`, etc).
   - Handlers `CreateLPN/GetLPN/ListLPNs/UpdateLPN/MoveLPN/MergeLPN/AddToLPN/Scan/InboundSync/ListSyncLogs` ahora mapean antes de retornar.

3. **Date format en `LotFormModal`**
   - Frontend enviaba `"2026-01-01"` pero backend espera RFC3339 full.
   - Fix: conversion a `YYYY-MM-DDT00:00:00Z` antes de enviar.

### Caso E2E ejecutado

Documento: `.claude/testing/inventory/front/CU-01-adjust-with-full-hierarchy.md`

- Creada bodega TEST-WH-01 con jerarquia completa (2 zonas, 3 pasillos, 4 racks, 9 niveles) via API autenticada desde Playwright.
- 3 posiciones + 2 lotes creados.
- 5 ajustes de stock probados (1 via modal UI, 4 via API) con combinaciones de `(location, lot, state)`.
- 4 `inventory_levels` quedaron separados correctamente por el unique key compuesto.
- 5 `stock_movements` con location_id y lot_id correctos.
- El bug del indice viejo fue detectado y corregido en el paso 2 de 5 ajustes.

### Seed de demostracion

5 bodegas reales creadas via UI + fetch autenticado desde browser:

| Bodega | Codigo | Ciudad | Zonas/Pasillos/Racks/Niveles/Posiciones |
|--------|--------|--------|----|
| Bodega Norte | WH-NORTE | Bogota | 2/4/8/16/20 |
| Bodega Sur | WH-SUR | Medellin | 2/4/8/16/20 |
| Bodega Pacifico | WH-CALI | Cali | 2/4/8/16/20 |
| Bodega Caribe | WH-COSTA | Barranquilla | 2/4/8/16/20 |
| Bodega Andina | WH-RIO | Bucaramanga | 2/4/8/16/20 |

Total generado: 10 zonas + 20 pasillos + 40 racks + 80 niveles + 100 posiciones.

### Gaps pendientes identificados

1. **`CreateLocationRequest` no acepta `level_id`** → las posiciones quedan vinculadas a la bodega pero no aparecen anidadas bajo niveles en el arbol. Los contadores POS:N no se llenan hasta que se extienda el endpoint.
2. **`Product.TrackInventory` flag no expuesto en UI** → solo se puede toggle via DB/API. Para el modal adjust/transfer seria util mostrar un switch en la pagina de productos.
3. **FEFO/FIFO en reserve no integrado** → `ListLotsForReserve` existe pero `ReserveStockForOrder` aun opera sobre nivel agregado.
4. **Fase 6 sin empezar** (analitica, ocupacion, alertas). Ver seccion al final.

### Archivos modificados/creados en esta sesion

**Backend (12 archivos):**
- `back/central/services/modules/inventory/internal/app/mappers/adjust_stock.go`
- `back/central/services/modules/inventory/internal/app/request/adjust_stock.go`
- `back/central/services/modules/inventory/internal/domain/dtos/transactional.go`
- `back/central/services/modules/inventory/internal/infra/primary/handlers/adjust_stock.go`
- `back/central/services/modules/inventory/internal/infra/primary/handlers/capture.go`
- `back/central/services/modules/inventory/internal/infra/primary/handlers/request/inventory.go`
- `back/central/services/modules/inventory/internal/infra/primary/handlers/transfer_stock.go`
- `back/central/services/modules/inventory/internal/infra/primary/handlers/response/capture.go` (nuevo)
- `back/central/services/modules/inventory/internal/infra/secondary/repository/adjust_stock_tx.go`
- `back/central/services/modules/inventory/internal/infra/secondary/repository/inventory_level_queries.go` (nuevo helper `getOrCreateLevelKeyTx`)
- `back/central/services/modules/inventory/internal/infra/secondary/repository/transfer_stock_tx.go`
- `back/migration/internal/infra/repository/migrate_inventory_traceability.go` (drop indice legacy)
- `back/migration/internal/infra/repository/migrate_inventory_audit.go`

**Frontend (35+ archivos):**
- `front/central/src/app/globals.css` (utilities del negocio + animaciones)
- `front/central/src/app/(auth)/inventory/analytics/slotting/page.tsx`
- `front/central/src/app/(auth)/inventory/audit/page.tsx` (3 tabs unificados)
- `front/central/src/app/(auth)/inventory/kardex/page.tsx`
- `front/central/src/app/(auth)/inventory/lpn/page.tsx`
- `front/central/src/app/(auth)/inventory/mobile/page.tsx`
- `front/central/src/app/(auth)/inventory/operations/page.tsx` (3 tabs unificados)
- `front/central/src/app/(auth)/inventory/sync/logs/page.tsx`
- `front/central/src/app/(auth)/inventory/traceability/page.tsx` (3 tabs unificados)
- `front/central/src/app/(auth)/warehouses/[id]/page.tsx`
- `front/central/src/services/modules/inventory/domain/types.ts` (agregados lot_id/state_id/uom_id/location_id)
- `front/central/src/services/modules/inventory/ui/components/AdjustStockModal.tsx`
- `front/central/src/services/modules/inventory/ui/components/TransferStockModal.tsx`
- `front/central/src/services/modules/inventory/ui/components/CountPlanFormModal.tsx` (nuevo)
- `front/central/src/services/modules/inventory/ui/components/LotFormModal.tsx` (nuevo)
- `front/central/src/services/modules/inventory/ui/components/SerialFormModal.tsx` (nuevo)
- `front/central/src/services/modules/inventory/ui/components/PutawayRuleFormModal.tsx` (nuevo)
- `front/central/src/services/modules/inventory/ui/components/InventoryTour.tsx` (nuevo, 13 pasos)
- `front/central/src/services/modules/warehouses/ui/components/WarehouseManager.tsx`
- `front/central/src/services/modules/warehouses/ui/components/WarehouseTreeTable.tsx` (nuevo)
- `front/central/src/services/modules/warehouses/ui/components/WarehouseHierarchyTree.tsx` (nuevo, reusable)
- `front/central/src/services/modules/warehouses/ui/components/HierarchyNodeModal.tsx` (nuevo)
- `front/central/src/services/modules/warehouses/ui/components/WarehouseTour.tsx` (nuevo, 11 pasos)
- `front/central/src/services/modules/my-integrations/ui/components/MyIntegrationsButton.tsx` (usa `.subnav-active`)
- `front/central/src/shared/ui/inventory-subnavbar.tsx` (scroll + secciones + tour pulse)

**Testing (3 archivos):**
- `.claude/testing/inventory/README.md` (nuevo)
- `.claude/testing/inventory/shared/test_data.md` (nuevo)
- `.claude/testing/inventory/front/CU-01-adjust-with-full-hierarchy.md` (nuevo)
- `.claude/testing/inventory/front/RESULTS.md` (nuevo, con bitacora de la ejecucion)

---

## Resumen Ejecutivo

El modulo de inventario evoluciono de un CRUD basico de stock a un WMS operativo en 5 de 6 fases planeadas. La ultima fase (analitica, ocupacion, alertas) queda pendiente.

| Fase | Estado | Commit | Descripcion |
|------|--------|--------|-------------|
| 0 - Normalizacion hexagonal | ✅ | `5e47806f` | app/{request,response,mappers} + repository/mappers + FE hooks/revalidatePath |
| 1 - Jerarquia fisica + cubicaje | ✅ | `455981a2` | Zones > Aisles > Racks > Levels > Positions + dimensiones + flags |
| 2 - Lote/Serie + estados + UoM | ✅ | `43abfb80` | Traceability + 8 estados + 10 UoM + conversion + FEFO/FIFO base |
| 3 - Operaciones | ✅ | `19a367e9` | Put-away rules + suggest/confirm + replenishment + cross-dock + slotting ABC |
| 4 - Auditoria + kardex | ✅ | `872955f3` | Cycle counts + discrepancies + approve/reject tx + kardex con running balance |
| 5 - Capture (scan/LPN/sync) | ✅ | `65e2ad75` | Scan resolver 6-way + LPN CRUD/move/merge + Inbound sync idempotente (SHA-256) |
| 6 - Analitica + alertas | 🔲 | - | OccupancySnapshot + RotationReport + ReorderAlert + dashboard + SSE |

---

## Arquitectura Hexagonal Vigente

```
back/central/services/modules/inventory/
├── bundle.go
└── internal/
    ├── domain/
    │   ├── entities/        # 21 entidades puras (sin tags)
    │   ├── dtos/            # TxParams, ListParams, value objects compartidos
    │   ├── ports/           # IRepository (~100 metodos), ISyncPublisher, IInventoryEventPublisher
    │   └── errors/          # 40+ errores tipados
    ├── app/
    │   ├── constructor.go   # IUseCase con ~90 metodos
    │   ├── request/         # DTOs de entrada desde handlers
    │   ├── response/        # DTOs de salida a handlers
    │   ├── mappers/         # Conversion request -> TxParams
    │   └── *_usecases.go    # Un archivo por grupo: adjust, transfer, bulk, reserve, confirm, release, return, lot, serial, state, uom, putaway, replenishment, crossdock, slotting, count_plan, count_task, discrepancy, kardex, lpn, scan, sync, validate_cubing
    ├── infra/
    │   ├── primary/
    │   │   ├── handlers/
    │   │   │   ├── constructor.go + routes.go
    │   │   │   ├── request/         # Structs con binding tags
    │   │   │   ├── response/        # JSON response structs
    │   │   │   └── *.go             # Un archivo por grupo de endpoints
    │   │   └── queue/               # order_consumer + bulk_load_consumer (RabbitMQ)
    │   └── secondary/
    │       ├── queue/               # sync_publisher
    │       ├── redis/               # inventory_cache + event_publisher (SSE)
    │       └── repository/
    │           ├── constructor.go
    │           ├── mappers/         # Todos los model <-> entity
    │           └── *_queries.go
    └── mocks/                       # Repository + UseCase mocks (inyectables por Fn fields)

back/central/services/modules/warehouses/    # Jerarquia (Fase 1)
└── internal/{domain,app,infra/primary/handlers,infra/secondary/repository,mocks}

back/migration/
├── shared/models/                   # 30+ modelos GORM (fuente de verdad)
└── internal/infra/repository/       # migrate_*.go registrados en constructor.go
```

**Reglas cumplidas:**
- Domain sin tags ni imports de infraestructura.
- Ports interface por modulo, consumidos via DI.
- Un constructor `New()` por capa.
- Un metodo HTTP por archivo handler + `routes.go` dedicado.
- Paginacion obligatoria en todo GET lista (default 10, max 100).
- Mappers extraidos a paquete propio.
- Multi-tenant via `resolveBusinessID` (JWT + `?business_id=X` para super admin).
- Aislamiento entre modulos via SELECT replicados (no cross-imports de repos).

---

## Inventario de Modelos GORM (30+)

### Jerarquia fisica (Fase 1)
- `Warehouse`, `WarehouseLocation` (refactor: level_id, dimensiones, flags JSONB, priority)
- `WarehouseZone`, `WarehouseAisle`, `WarehouseRack`, `WarehouseRackLevel`

### Stock y traceability (Fase 0 base + Fase 2)
- `InventoryLevel` (refactor: lot_id, state_id; unique key `(product, warehouse, location, lot, state)`)
- `StockMovement` (refactor: lot_id, serial_id, from_state_id, to_state_id, uom_id, qty_in_base_uom)
- `StockMovementType`
- `InventoryLot`, `InventorySerial`
- `InventoryState` (seed 8: available, reserved, on_hold, damaged, quarantine, expired, in_transit, returned)
- `UnitOfMeasure` (seed 10: UN, KG, G, LB, L, ML, CM, M, CAJA, PALETA)
- `ProductUoM`

### Operaciones (Fase 3)
- `PutawayRule`, `PutawaySuggestion`
- `ReplenishmentTask`
- `CrossDockLink`
- `ProductVelocity`

### Auditoria (Fase 4)
- `CycleCountPlan`, `CycleCountTask`, `CycleCountLine`
- `InventoryDiscrepancy`
- (+ seed movement type `count_adjustment`)

### Capture (Fase 5)
- `LicensePlate`, `LicensePlateLine`
- `ScanEvent`
- `InventorySyncLog` (unique index `(business_id, direction_key, payload_hash)` para idempotencia)

---

## Catalogo Completo de Endpoints

### Modulo `warehouses` - prefix `/api/v1`

#### Warehouses CRUD
- `GET /warehouses` - listar bodegas (paginado; filtros is_active, is_fulfillment, search)
- `POST /warehouses` - crear bodega con datos carrier (lat/lng, DANE, street, postal)
- `GET /warehouses/:id` - detalle de bodega
- `PUT /warehouses/:id` - actualizar bodega
- `DELETE /warehouses/:id` - soft delete de bodega
- `GET /warehouses/:id/tree` - arbol jerarquico completo (Zone > Aisle > Rack > Level > Position) en 5 queries batch
- `GET /warehouses/:id/zones` - listar zonas de una bodega

#### Warehouse Locations (posiciones fisicas, hojas de la jerarquia)
- `GET /warehouses/:id/locations` - listar posiciones (incluye level_id, dimensiones, flags JSONB, priority)
- `POST /warehouses/:id/locations` - crear posicion con cubing (max_weight_kg, max_volume_cm3, length/width/height)
- `PUT /warehouses/:id/locations/:locationId` - actualizar posicion
- `DELETE /warehouses/:id/locations/:locationId` - eliminar posicion

#### Zones CRUD
- `POST /zones` - crear zona (warehouse_id, code, name, purpose, color_hex)
- `GET /zones/:zoneId` - detalle de zona
- `PUT /zones/:zoneId` - actualizar zona
- `DELETE /zones/:zoneId` - eliminar zona
- `GET /zones/:zoneId/aisles` - listar pasillos de zona

#### Aisles CRUD
- `POST /aisles` - crear pasillo (zone_id, code, name)
- `GET /aisles/:aisleId` - detalle
- `PUT /aisles/:aisleId` - actualizar
- `DELETE /aisles/:aisleId` - eliminar
- `GET /aisles/:aisleId/racks` - listar racks de pasillo

#### Racks CRUD
- `POST /racks` - crear rack (aisle_id, code, name, levels_count)
- `GET /racks/:rackId` - detalle
- `PUT /racks/:rackId` - actualizar
- `DELETE /racks/:rackId` - eliminar
- `GET /racks/:rackId/levels` - listar niveles de rack

#### Rack Levels CRUD
- `POST /rack-levels` - crear nivel (rack_id, code, ordinal)
- `GET /rack-levels/:levelId` - detalle
- `PUT /rack-levels/:levelId` - actualizar
- `DELETE /rack-levels/:levelId` - eliminar

### Modulo `inventory` - prefix `/api/v1/inventory`

#### Stock basico (Fase 0)
- `GET /inventory/product/:productId` - niveles del producto en todas las bodegas
- `GET /inventory/warehouse/:warehouseId` - listar niveles de una bodega (paginado; search, low_stock filter)
- `POST /inventory/adjust` - ajuste manual (+/-) con AdjustStockTx (SELECT FOR UPDATE + StockMovement)
- `POST /inventory/transfer` - transferencia entre bodegas transaccional
- `POST /inventory/bulk-load` - carga masiva (sync o async via RabbitMQ si disponible)
- `GET /inventory/movements` - kardex filtrable por producto/bodega/tipo

#### Cubicaje (Fase 1)
- `POST /inventory/positions/validate-cubing` - valida si qty*dimensiones cabe en `location_id`; devuelve `fits, weight_needed/max, volume_needed/max, occupied_qty`

#### Lots (Fase 2)
- `GET /inventory/lots` - listar (filtros product_id, status, expiring_in_days)
- `POST /inventory/lots` - crear lote (lot_code, manufacture_date, expiration_date, supplier_id)
- `GET /inventory/lots/:id` - detalle
- `PUT /inventory/lots/:id` - actualizar
- `DELETE /inventory/lots/:id` - eliminar

#### Serials (Fase 2)
- `GET /inventory/serials` - listar (filtros product_id, lot_id, state_id, location_id)
- `POST /inventory/serials` - crear numero de serie
- `GET /inventory/serials/:id` - detalle
- `PUT /inventory/serials/:id` - actualizar lot_id/location/state

#### States (Fase 2)
- `GET /inventory/states` - listar catalogo de 8 estados
- `POST /inventory/state-transitions` - ChangeInventoryState transaccional (split entre levels por estado)

#### UoM y conversion (Fase 2)
- `GET /inventory/uoms` - catalogo global de UoM
- `GET /inventory/products/:productId/uoms` - UoM asignados a un producto
- `POST /inventory/products/:productId/uoms` - asignar UoM (conversion_factor, is_base, barcode)
- `DELETE /inventory/product-uoms/:id` - eliminar asignacion
- `POST /inventory/uoms/convert` - convertir cantidad entre UoMs usando factor base

#### Movement Types (Fase 0)
- `GET /inventory/movement-types` - listar tipos
- `POST /inventory/movement-types` - crear tipo (code, name, direction)
- `PUT /inventory/movement-types/:id` - actualizar
- `DELETE /inventory/movement-types/:id` - eliminar

#### Put-away (Fase 3)
- `GET /inventory/putaway-rules` - listar reglas (filtro active_only)
- `POST /inventory/putaway-rules` - crear regla (product_id?, category_id?, target_zone_id, priority, strategy)
- `PUT /inventory/putaway-rules/:id` - actualizar
- `DELETE /inventory/putaway-rules/:id` - eliminar
- `POST /inventory/putaway/suggest` - recibe items [{product_id, quantity}], busca rule aplicable + pick location libre en zona, crea PutawaySuggestion
- `POST /inventory/putaway/suggestions/:id/confirm` - confirma put-away con actual_location_id
- `GET /inventory/putaway/suggestions` - listar sugerencias (filtro status)

#### Replenishment (Fase 3)
- `GET /inventory/replenishment/tasks` - listar tareas (filtros warehouse, status, assigned_to)
- `POST /inventory/replenishment/tasks` - crear tarea manual
- `POST /inventory/replenishment/tasks/:id/assign` - asignar a usuario (status -> in_progress)
- `POST /inventory/replenishment/tasks/:id/complete` - completar tarea
- `POST /inventory/replenishment/tasks/:id/cancel` - cancelar
- `POST /inventory/replenishment/detect` - escanea inventory_levels donde available_qty < reorder_point, crea tasks automaticas

#### Cross-dock (Fase 3)
- `GET /inventory/cross-dock/links` - listar (filtros outbound_order_id, status)
- `POST /inventory/cross-dock/links` - crear link (inbound_shipment_id?, outbound_order_id, product_id, qty)
- `POST /inventory/cross-dock/links/:id/execute` - ejecutar cross-dock (status -> executed + executed_at)

#### Slotting (Fase 3)
- `POST /inventory/slotting/run` - computa velocities de ultimos N dias (7d/30d/90d/180d/365d) y clasifica ABC (80%/95% acumulado)
- `GET /inventory/slotting/velocities` - listar productos por velocity (filtros period, rank, limit)

#### Cycle Count Plans (Fase 4)
- `GET /inventory/cycle-count-plans` - listar (filtros warehouse_id, active_only)
- `POST /inventory/cycle-count-plans` - crear plan (warehouse_id, name, strategy abc|zone|random|full, frequency_days, next_run_at)
- `PUT /inventory/cycle-count-plans/:id` - actualizar
- `DELETE /inventory/cycle-count-plans/:id` - eliminar

#### Cycle Count Tasks (Fase 4)
- `GET /inventory/cycle-count-tasks` - listar (filtros warehouse_id, plan_id, status)
- `POST /inventory/cycle-count-tasks/generate` - genera task + materializa lineas desde inventory_levels segun strategy (abc join ProductVelocity rank=A, zone join hierarchy con scope_id, random LIMIT 50, full todo)
- `POST /inventory/cycle-count-tasks/:id/start` - inicia (status -> in_progress, started_at, assigned_to_id)
- `POST /inventory/cycle-count-tasks/:id/finish` - finaliza task
- `GET /inventory/cycle-count-tasks/:taskId/lines` - listar lineas del conteo

#### Cycle Count Lines (Fase 4)
- `POST /inventory/cycle-count-lines/:id/submit` - registra counted_qty, computa variance; si variance != 0 abre InventoryDiscrepancy automaticamente

#### Discrepancies (Fase 4)
- `GET /inventory/discrepancies` - listar (filtros task_id, status open|approved|rejected)
- `POST /inventory/discrepancies/:id/approve` - ApproveDiscrepancyTx transaccional: SELECT FOR UPDATE nivel + aplica delta (counted - expected) + crea StockMovement tipo count_adjustment + resuelve line
- `POST /inventory/discrepancies/:id/reject` - marca rejected con reason

#### Kardex (Fase 4)
- `GET /inventory/kardex/export?product_id=&warehouse_id=&from=&to=` - kardex con running balance acumulado + totales in/out/final_balance; join con movement_types

#### Scan (Fase 5)
- `POST /inventory/scan` - recibe code, resuelve 6-way: LPN code -> WarehouseLocation code -> InventorySerial -> InventoryLot -> ProductUoM barcode -> Product (SKU o ID). Registra ScanEvent con CodeType detectado. Retorna `{resolved, resolution, event}`

#### LPN (Fase 5)
- `GET /inventory/lpn` - listar (filtros lpn_type, status, location_id)
- `POST /inventory/lpn` - crear LPN (code, lpn_type pallet|case|tote, location_id)
- `GET /inventory/lpn/:id` - detalle (con Lines preload)
- `PUT /inventory/lpn/:id` - actualizar
- `DELETE /inventory/lpn/:id` - eliminar
- `POST /inventory/lpn/:id/lines` - agregar producto (product_id, lot_id?, serial_id?, qty)
- `POST /inventory/lpn/:id/move` - mover a nueva ubicacion
- `POST /inventory/lpn/:id/dissolve` - disolver (status -> dissolved)
- `POST /inventory/lpn/:id/merge` - mezclar con otra LPN (re-parenta lineas y disuelve fuente)

#### Sync / Integrations (Fase 5)
- `POST /inventory/sync/inbound/:integrationId` - sync entrante idempotente (hash SHA-256 del payload; retorna duplicate=true si ya procesado)
- `GET /inventory/sync/logs` - logs de sincronizacion (filtros integration_id, direction in|out, status)

---

## Eventos RabbitMQ y Redis

### Estado actual
- **Consumer activo**: `OrderConsumer` (orders_to_inventory queue)
- **Consumer activo**: `BulkLoadConsumer` (inventory.bulk_load queue)
- **Publisher activo**: `SyncPublisher` (outbound a integraciones)
- **Redis cache**: `InventoryCache` por producto/bodega (TTL configurable)
- **SSE**: `EventPublisher` por canal `probability:inventory:state:events`

### Eventos emitidos hoy
- `inventory.reserved`, `inventory.confirmed`, `inventory.released`, `inventory.returned`, `inventory.insufficient`
- `bulk_load.failed`, `bulk_load.completed`

### Eventos declarados en plan pero no cableados (candidatos a Fase 6 o integraciones futuras)
- `inventory.lot.created|expired`, `inventory.serial.scanned`, `inventory.state.changed`
- `inventory.replenishment.requested|assigned|completed`, `inventory.putaway.suggested|confirmed`, `inventory.slotting.run_completed`, `inventory.cross_dock.linked|executed`
- `inventory.count.task_started`, `count.line_submitted`, `discrepancy.opened|approved|rejected`
- `inventory.lpn.created|moved|dissolved|merged`, `inventory.scan.recorded`, `inventory.sync.inbound_received|outbound_sent`

---

## Frontend Vigente

### Modulo `inventory` (frontend/central)
```
domain/
├── types.ts                 # InventoryLevel, StockMovement, MovementType, PaginationParams, PaginatedResponse<T>
├── traceability-types.ts    # Lot, Serial, UoM, ProductUoM, ConvertUoMResult, InventoryState
├── operations-types.ts      # PutawayRule, PutawaySuggestion, ReplenishmentTask, CrossDockLink, ProductVelocity, SlottingRunResult
├── audit-types.ts           # CycleCountPlan/Task/Line, InventoryDiscrepancy, KardexEntry, KardexExportResult
└── capture-types.ts         # LicensePlate/Line, ScanResolution/Result, ScanEvent, InventorySyncLog

infra/
├── repository/
│   ├── api-repository.ts
│   ├── traceability-api-repository.ts
│   ├── operations-api-repository.ts
│   ├── audit-api-repository.ts
│   └── capture-api-repository.ts
└── actions/
    ├── index.ts             # adjust, transfer, bulk, get* (con revalidatePath en mutaciones)
    ├── traceability.ts
    ├── operations.ts
    ├── audit.ts
    └── capture.ts

ui/
├── components/              # Componentes pre-WMS (InventoryLevelList, StockMovementList, AdjustStockModal, TransferStockModal, BulkLoadInventoryModal, InventoryManager)
└── hooks/
    ├── index.ts
    ├── useInventoryLevels.ts
    ├── useMovements.ts
    ├── useMovementTypes.ts
    ├── useLots.ts
    ├── useSerials.ts
    ├── useUoMs.ts            # useUoMs + useProductUoMs + useUoMConverter
    ├── useOperations.ts      # usePutawayRules + usePutawaySuggestions + useReplenishmentTasks + useCrossDockLinks + useVelocities
    ├── useAudit.ts           # useCountPlans + useCountTasks + useCountLines + useDiscrepancies + useKardex
    └── useCapture.ts         # useLPNs + useScan + useSyncLogs
```

### Modulo `warehouses` (frontend/central)
```
domain/
├── types.ts                   # Warehouse, WarehouseLocation + DTOs originales
└── hierarchy-types.ts         # Zone, Aisle, Rack, RackLevel, Tree*, CubingCheckResult

infra/
├── repository/
│   ├── api-repository.ts
│   └── hierarchy-api-repository.ts
└── actions/
    ├── index.ts
    └── hierarchy.ts           # Zone/Aisle/Rack/Level CRUD + validateCubing + tree

ui/hooks/
├── index.ts
├── useWarehouseTree.ts
└── useCubing.ts
```

### Paginas App Router existentes
- `/inventory` y `/inventory/movements` (pre-WMS, usan componentes actuales)

### Paginas pendientes de implementar
- `/inventory/warehouses`, `/inventory/warehouses/[id]` (arbol visual + editor posiciones)
- `/inventory/lots`, `/inventory/serials`
- `/inventory/operations/{putaway,replenishment,cross-dock}`
- `/inventory/analytics/slotting`
- `/inventory/audit/{plans,tasks,discrepancies}`
- `/inventory/kardex`
- `/inventory/mobile` (PWA para scan)
- `/inventory/lpn`, `/inventory/sync/logs`

Los hooks y actions estan listos; solo falta la capa de UI.

---

## Que Queda: Fase 6 - Analitica, Ocupacion y Alertas

### Modelos GORM nuevos
- `OccupancySnapshot` - warehouse_id, zone_id, used_volume_cm3, total_volume_cm3, used_weight_kg, occupied_positions, total_positions, snapshot_at
- `RotationReport` - product_id, warehouse_id, period, turns, days_on_hand, is_obsolete_flag, computed_at
- `ReorderAlert` - product_id, warehouse_id, level, threshold, current_qty, status (open|acknowledged|resolved), triggered_at, acknowledged_at, resolved_at

### Migracion
- `migrate_inventory_analytics.go` con AutoMigrate (DDL idempotente)

### Domain + ports
- Entidades + DTOs en inventory module
- Extender `IRepository` con: `CreateOccupancySnapshot`, `ListOccupancySnapshots`, `RefreshOccupancyFor`, `ComputeRotationFor`, `ListRotation`, `CreateAlert`, `ListAlerts`, `UpdateAlertStatus`, `EvaluateReorderThresholdsFor`
- Errores: `ErrAlertNotFound`, `ErrAlertAlreadyResolved`

### App use cases
- `RefreshOccupancy(warehouse_id)` - agrega volumen/peso ocupado por zone/warehouse desde inventory_levels + warehouse_locations
- `ComputeRotation(warehouse_id, period)` - calcula turns (ventas / stock promedio) y days_on_hand; flag obsolete si days_on_hand > threshold
- `EvaluateReorderThresholds()` - escanea inventory_levels por business, crea o reabre ReorderAlert si cruza threshold
- `AcknowledgeAlert(alert_id, user_id)` / `ResolveAlert(alert_id, user_id)`
- `GetDashboardMetrics(warehouse_id)` - agrega ocupacion + top 10 rotacion + obsoletos + alertas open (cacheable Redis 60s)

### Redis
- `probability:inventory:dashboard:{warehouse_id}` TTL 60s
- `probability:inventory:rotation:{warehouse_id}:{period}` TTL 300s

### Eventos RabbitMQ
- `inventory.occupancy.refreshed`
- `inventory.rotation.computed`
- `inventory.alert.triggered|acknowledged|resolved`

### Handlers + rutas
```
GET  /inventory/dashboard/occupancy?warehouse_id=
GET  /inventory/dashboard/rotation?warehouse_id=&period=
GET  /inventory/dashboard/obsolescence?warehouse_id=
GET  /inventory/dashboard/metrics?warehouse_id=   # agregado
GET  /inventory/alerts
POST /inventory/alerts/:id/acknowledge
POST /inventory/alerts/:id/resolve
POST /inventory/occupancy/refresh                  # job manual
POST /inventory/rotation/compute                   # job manual
POST /inventory/alerts/evaluate                    # job manual
```

### Notificaciones
- Publicar a cola del modulo `events` (no duplicar logica email/push en inventory).

### Frontend
- Tipos `OccupancySnapshot`, `RotationReport`, `ReorderAlert`, `DashboardMetrics` en `domain/analytics-types.ts`
- Repository `analytics-api-repository.ts`
- Actions con revalidatePath en acknowledge/resolve
- Hooks: `useOccupancy`, `useRotation`, `useReorderAlerts`, `useDashboardMetrics`
- Paginas:
  - `/inventory/dashboard` con cards (% ocupacion por bodega, top 10 rotacion, obsoletos, alertas open)
  - `/inventory/alerts` con tabla + acciones acknowledge/resolve
- Componentes:
  - `OccupancyHeatmap` (recharts: ocupacion por zone)
  - `RotationChart` (recharts: top productos por turns)
  - `AlertsBanner` en dashboard home
- SSE subscription para refresh en vivo via `useSSE('/notify/sse/inventory-notify', { event_types: ['alert.triggered','occupancy.refreshed'] })`

### Estimacion Fase 6
- Backend: 2-3 dias (modelos + migracion + app + repo + handlers + mocks)
- Frontend: 2-3 dias (tipos + repo + actions + hooks + 2 paginas + 3 componentes)
- Testing E2E: 1 dia
- **Total: ~1 semana**

---

## Integraciones entre Modulos

| Modulo consumidor | Modulo proveedor | Metodo | Via |
|-------------------|-------------------|--------|-----|
| inventory | orders | Recibe eventos reserve/confirm/release/return | RabbitMQ (`orders_to_inventory`) |
| inventory | products | Consulta productos y dimensiones | SELECT replicado local |
| inventory | warehouses | Consulta locations para cubing y slotting | SELECT replicado local |
| inventory | integrations | Outbound sync stock deltas | RabbitMQ (SyncPublisher) |
| inventory | events (fase 6) | Publicar alert.triggered | RabbitMQ (pendiente) |

**No se importan repositorios entre modulos** - replicacion de SELECTs respeta regla `backend-conventions.md #1`.

---

## Deuda Tecnica Pendiente

1. **FEFO/FIFO en reserve completo**: `ListLotsForReserve` existe pero no esta integrado al `ReserveStockForOrder`. Requiere refactor de `reserve_stock_tx.go` para iterar lotes segun strategy.
2. **Consumers scheduled**: Jobs periodicos para `RefreshOccupancy`, `ComputeRotation`, `EvaluateReorderThresholds`, `DetectReplenishmentNeeds`, `ExpireLots` - se espera cron externo que publique a cola.
3. **Eventos no cableados**: 18+ routing keys declaradas pero sin publisher activo (ver seccion eventos).
4. **Paginas frontend**: Hooks listos pero paginas `/inventory/*` de Fases 1-5 no construidas (solo quedan las pre-WMS).
5. **Testing E2E**: No se crearon casos en `.claude/testing/inventory/`.
6. **Comentarios en modelos**: Algunos archivos en `migration/shared/models/` conservan acentos/comentarios que deberian limpiarse (regla cero-comentarios del proyecto).
7. **Unused function warnings**: Varios `migrate_inventory_*` aparecen como unused en linter hasta que se compila desde `cmd/main.go`.
8. **Helper `NoPutawayRuleFound`**: Metodo no usado en `putaway_usecases.go` (conveniencia que quedo declarada).

---

## Verificacion E2E por Fase

**Backend:**
```bash
cd /home/cam/Desktop/probability/back/central && go build ./... && go test ./services/modules/inventory/... && go test ./services/modules/warehouses/...
cd /home/cam/Desktop/probability/back/migration && go build ./...
```

**Frontend:**
```bash
cd /home/cam/Desktop/probability/front/central && pnpm exec tsc --noEmit 2>&1 | grep inventory
```

**Migracion (en desarrollo local):**
```bash
cd /home/cam/Desktop/probability/back/migration && go run cmd/main.go
```

Verificacion via `mcp__postgres-probability__query`:
```sql
SELECT tablename FROM pg_tables WHERE schemaname = 'public'
  AND tablename IN (
    'warehouse_zones','warehouse_aisles','warehouse_racks','warehouse_rack_levels',
    'inventory_states','units_of_measure','product_uoms','inventory_lots','inventory_serials',
    'putaway_rules','putaway_suggestions','replenishment_tasks','cross_dock_links','product_velocities',
    'cycle_count_plans','cycle_count_tasks','cycle_count_lines','inventory_discrepancies',
    'license_plates','license_plate_lines','scan_events','inventory_sync_logs'
  );
SELECT code FROM inventory_states;
SELECT code FROM units_of_measure;
SELECT code FROM stock_movement_types WHERE code = 'count_adjustment';
```

---

## Total de Codigo Agregado

Contando solo archivos nuevos/modificados de las 6 fases WMS:

- **Backend**: ~14,000 LOC (entidades, dtos, ports, use cases, handlers, repositorios, mappers, mocks)
- **Frontend**: ~2,500 LOC (tipos, repositorios, actions, hooks)
- **Migracion**: ~900 LOC (modelos, funciones de migracion, seeds)
- **Endpoints nuevos**: ~85 (30 warehouses jerarquia + 55+ inventory WMS)
- **Modelos GORM nuevos**: 18
- **Modelos GORM refactorizados**: 3 (InventoryLevel, StockMovement, WarehouseLocation)
- **Use cases backend**: ~90 metodos en IUseCase
