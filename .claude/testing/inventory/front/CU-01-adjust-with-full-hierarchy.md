# CU-01: Ajuste de stock con jerarquía completa, lotes y múltiples productos

## Objetivo

Validar que todo el flujo WMS funciona de extremo a extremo:
1. Crear una bodega con jerarquía física (zonas → pasillos → racks → niveles → posiciones)
2. Confirmar que los productos soportan `track_inventory`
3. Crear lotes con fecha de vencimiento
4. Ajustar stock seleccionando ubicación específica + lote específico
5. Verificar que los `inventory_levels` quedan diferenciados por `(product, warehouse, location, lot, state)`
6. Verificar que el `stock_movement` guarda la ubicación afectada

## Precondiciones

- Backend corriendo en `localhost:3050`
- Frontend corriendo en `localhost:3000`
- Login con `demo@probability.com / ProbabilityDemo` (business_id=26)
- Al menos 3 productos con `track_inventory=true` en el business (ver `shared/test_data.md`)

## Pasos

### Paso 1: Crear bodega

**UI**: `/warehouses` → botón "Nueva bodega"

- Código: `TEST-WH-01`
- Nombre: `Bodega Test Completa`
- Ciudad: `Medellin`
- Activa: sí

**Expected**: bodega aparece en el listado.

**Verificación DB**:
```sql
SELECT id, code, name FROM warehouses WHERE code = 'TEST-WH-01';
```

### Paso 2: Construir jerarquía física

**UI**: Click en "Ver detalle" de la bodega → `/warehouses/:id`

#### 2.1 Crear zona Picking

- Click "Crear primera zona"
- Código: `Z-PICK`, Nombre: `Zona Picking`, Propósito: Picking, Color: `#10b981`

#### 2.2 Crear zona Bulk

- Click "Nueva zona"
- Código: `Z-BULK`, Nombre: `Zona Bulk`, Propósito: Bulk, Color: `#f59e0b`

#### 2.3 Crear pasillos

- En Z-PICK → "Agregar pasillo" → `A-01` / Pasillo Frontal
- En Z-PICK → "Agregar pasillo" → `A-02` / Pasillo Trasero
- En Z-BULK → "Agregar pasillo" → `B-01` / Bulk principal

#### 2.4 Crear racks

- En A-01 → Rack `R-01` / Rack Principal / 3 niveles
- En A-01 → Rack `R-02` / Rack Secundario / 2 niveles
- En A-02 → Rack `R-03` / Rack Corto / 1 nivel
- En B-01 → Rack `R-04` / Rack Bulk Alto / 2 niveles

#### 2.5 Crear niveles

Para cada rack, agregar los niveles requeridos (código L-01, L-02, L-03 con ordinales 1, 2, 3).

**Expected**: el árbol muestra la jerarquía completa, contadores `Zonas=2 · Pasillos=3 · Racks=4 · Niveles=9`.

**Verificación DB**:
```sql
SELECT
  (SELECT COUNT(*) FROM warehouse_zones WHERE warehouse_id =
      (SELECT id FROM warehouses WHERE code='TEST-WH-01')) as zonas,
  (SELECT COUNT(*) FROM warehouse_aisles a JOIN warehouse_zones z ON z.id = a.zone_id
      WHERE z.warehouse_id = (SELECT id FROM warehouses WHERE code='TEST-WH-01')) as pasillos,
  (SELECT COUNT(*) FROM warehouse_racks r JOIN warehouse_aisles a ON a.id = r.aisle_id
      JOIN warehouse_zones z ON z.id = a.zone_id
      WHERE z.warehouse_id = (SELECT id FROM warehouses WHERE code='TEST-WH-01')) as racks,
  (SELECT COUNT(*) FROM warehouse_rack_levels l JOIN warehouse_racks r ON r.id = l.rack_id
      JOIN warehouse_aisles a ON a.id = r.aisle_id
      JOIN warehouse_zones z ON z.id = a.zone_id
      WHERE z.warehouse_id = (SELECT id FROM warehouses WHERE code='TEST-WH-01')) as niveles;
```

### Paso 3: Crear posiciones físicas

**UI**: desde la misma página, sección de posiciones (o vía `POST /warehouses/:id/locations`).

- `POS-A0101` → level_id del Rack R-01 Nivel L-01
- `POS-A0201` → level_id del Rack R-03 Nivel L-01
- `POS-B0101` → level_id del Rack R-04 Nivel L-02

**Expected**: 3 posiciones nuevas con código y level_id asignado.

**Verificación DB**:
```sql
SELECT code, name, level_id FROM warehouse_locations
WHERE code IN ('POS-A0101','POS-A0201','POS-B0101');
```

### Paso 4: Identificar productos con tracking

**Verificación DB**:
```sql
SELECT id, sku, name FROM products WHERE business_id = 26 AND track_inventory = true LIMIT 3;
```

Si no hay 3, habilitar tracking:
```sql
UPDATE products SET track_inventory = true
WHERE business_id = 26
  AND id IN (SELECT id FROM products WHERE business_id = 26 ORDER BY created_at LIMIT 3);
```

Guardar los 3 product_id como P1, P2, P3.

### Paso 5: Crear lotes para P1

**UI**: `/inventory/traceability` → tab Lotes → "Nuevo lote"

- Producto: P1 (SKU)
- Código: `LOT-TEST-A01`, Vencimiento: +60 días
- Producto: P1 (SKU)
- Código: `LOT-TEST-A02`, Vencimiento: +120 días

**Expected**: 2 lotes activos para P1.

### Paso 6: Ajustes de stock (el caso duro)

**UI**: `/inventory` → seleccionar bodega `TEST-WH-01` → botón "Ajustar stock" (icono).

#### Ajuste 1 — P1 en POS-A0101 con LOT-TEST-A01 (+30)

- Producto: P1
- **Ubicación**: POS-A0101
- **Lote**: LOT-TEST-A01
- Cantidad: 30
- Estado: Disponible
- Razón: Recibo inicial

**Expected**:
- Toast: "Ajuste aplicado: +30 uds · ubicación POS-A0101 · lote LOT-TEST-A01"
- Stock = 30 en ese nivel específico

#### Ajuste 2 — P1 en POS-A0201 con LOT-TEST-A02 (+50)

- Producto: P1, Ubicación POS-A0201, Lote LOT-TEST-A02, qty 50

#### Ajuste 3 — P2 en POS-B0101 sin lote (+100)

- Producto: P2, Ubicación POS-B0101, sin lote, qty 100

#### Ajuste 4 — P3 sin ubicación (+20)

- Producto: P3, **sin ubicación** (stock general), qty 20

#### Ajuste 5 — P1 en POS-A0101 con LOT-TEST-A01 (−5)

- Producto: P1, Ubicación POS-A0101, Lote LOT-TEST-A01, qty -5, Razón: Merma

**Expected final**: el nivel de P1 en POS-A0101/LOT-TEST-A01 queda en qty=25.

### Paso 7: Verificar niveles en DB

```sql
SELECT il.id, p.sku, w.code AS warehouse, wl.code AS location, l.lot_code AS lot,
       s.code AS state, il.quantity, il.available_qty
FROM inventory_levels il
JOIN products p ON p.id = il.product_id
JOIN warehouses w ON w.id = il.warehouse_id
LEFT JOIN warehouse_locations wl ON wl.id = il.location_id
LEFT JOIN inventory_lots l ON l.id = il.lot_id
LEFT JOIN inventory_states s ON s.id = il.state_id
WHERE w.code = 'TEST-WH-01'
ORDER BY il.id;
```

**Expected**: 4 filas

| product | location | lot | state | quantity |
|---------|----------|-----|-------|----------|
| P1 | POS-A0101 | LOT-TEST-A01 | available | 25 |
| P1 | POS-A0201 | LOT-TEST-A02 | available | 50 |
| P2 | POS-B0101 | (null) | available | 100 |
| P3 | (null) | (null) | available | 20 |

Esto demuestra que el unique key `(product, warehouse, location, lot, state)` crea filas distintas y no mezcla el stock de lotes/posiciones.

### Paso 8: Verificar movimientos en DB

```sql
SELECT sm.id, p.sku, wl.code AS location, l.lot_code AS lot, sm.quantity,
       sm.previous_qty, sm.new_qty, sm.reason
FROM stock_movements sm
JOIN products p ON p.id = sm.product_id
JOIN warehouses w ON w.id = sm.warehouse_id
LEFT JOIN warehouse_locations wl ON wl.id = sm.location_id
LEFT JOIN inventory_lots l ON l.id = sm.lot_id
WHERE w.code = 'TEST-WH-01'
ORDER BY sm.id;
```

**Expected**: 5 movimientos, cada uno con `location_id` y `lot_id` correctos cuando aplica.

### Paso 9: Ver kardex

**UI**: `/inventory/kardex` → Producto P1, Bodega TEST-WH-01, sin rango.

**Expected**:
- 3 movimientos listados (2 entradas +30 y +50, 1 salida -5)
- Total entradas: 80, Total salidas: 5, Saldo final: 75

## Cleanup (opcional)

```sql
DELETE FROM stock_movements WHERE warehouse_id = (SELECT id FROM warehouses WHERE code='TEST-WH-01');
DELETE FROM inventory_levels WHERE warehouse_id = (SELECT id FROM warehouses WHERE code='TEST-WH-01');
DELETE FROM inventory_lots WHERE lot_code LIKE 'LOT-TEST-A%';
DELETE FROM warehouse_locations WHERE code LIKE 'POS-%' AND warehouse_id = (SELECT id FROM warehouses WHERE code='TEST-WH-01');
-- Los niveles, racks, pasillos, zonas se eliminan en cascada si se borra la bodega
DELETE FROM warehouses WHERE code = 'TEST-WH-01';
```

## Bugs que este caso detectó

(Actualizar al ejecutar)
