# RESULTS - Inventario (Frontend E2E)

## 2026-04-21 — CU-01-adjust-with-full-hierarchy

### Resultado: OK con 1 bug detectado y fixeado

### Ejecución

| Paso | Estado | Detalle |
|------|--------|---------|
| 1. Crear bodega TEST-WH-01 | ✅ | warehouse_id=5 |
| 2. Jerarquía (2 zonas, 3 pasillos, 4 racks, 9 niveles) | ✅ | Vía API con token JWT |
| 3. 3 posiciones (POS-A0101, POS-A0201, POS-B0102) | ✅ | id=1,2,3 (sin level_id por limitación actual del endpoint de locations) |
| 4. 2 lotes para PT01001 (`LOT-TEST-A01`, `LOT-TEST-A02`) | ✅ | id=2,3 |
| 5. Ajuste #1 desde modal UI (+30 en POS-A0101/LOT-A01) | ✅ | Modal mostró selectores condicionales de ubicación y lote |
| 6. Ajuste #2 via API (+50 en POS-A0201/LOT-A02) | ✅ tras fix | Bloqueado inicialmente por bug |
| 7. Ajuste #3 via API (+100 en POS-B0102 sin lote) | ✅ tras fix | |
| 8. Ajuste #4 via API (+20 sin ubicación, stock flotante) | ✅ tras fix | |
| 9. Ajuste #5 via API (-5 merma en POS-A0101/LOT-A01) | ✅ | |

### Estado final en DB (inventory_levels)

```
id | sku     | warehouse   | location   | lot           | state     | qty
---|---------|-------------|------------|---------------|-----------|----
 2 | PT01001 | TEST-WH-01  | POS-A0101  | LOT-TEST-A01  | available | 25
 6 | PT01001 | TEST-WH-01  | POS-A0201  | LOT-TEST-A02  | available | 50
 7 | PT01001 | TEST-WH-01  | POS-B0102  | (null)        | available | 100
 8 | PT01001 | TEST-WH-01  | (null)     | (null)        | available | 20
```

4 filas independientes por unique key `(product, warehouse, location, lot, state)`. El sistema NO mezcló el stock entre lotes ni ubicaciones.

### Movimientos (stock_movements)

```
id | sku     | location   | lot           | qty | prev | new | reason
---|---------|------------|---------------|-----|------|-----|-------------------
 2 | PT01001 | POS-A0101  | LOT-TEST-A01  | +30 |   0  |  30 | Recibo inicial E2E
 3 | PT01001 | POS-A0101  | LOT-TEST-A01  |  -5 |  30  |  25 | Merma por humedad
 4 | PT01001 | POS-A0201  | LOT-TEST-A02  | +50 |   0  |  50 | Recibo parcial lote A02
 5 | PT01001 | POS-B0102  | (null)        |+100 |   0  | 100 | Recibo bulk sin lote
 6 | PT01001 | (null)     | (null)        | +20 |   0  |  20 | Stock flotante de bodega
```

Cada movimiento guardó `location_id` y `lot_id` correctos cuando aplica.

### Bug detectado

**BUG**: El segundo ajuste con diferente `location_id`+`lot_id` sobre el mismo producto fallaba con error 500:

```
ERROR: duplicate key value violates unique constraint
       "idx_inventory_product_warehouse" (SQLSTATE 23505)
```

**Causa**: Existía un índice UNIQUE antiguo `idx_inventory_product_warehouse` (solo `product_id + warehouse_id`) de antes de la Fase 2. Cuando se refactorizó el modelo para incluir `lot_id` y `state_id` en el unique key (creando `idx_inventory_level_key`), el índice viejo quedó huérfano y bloqueaba crear múltiples niveles del mismo producto en la misma bodega.

**Fix**: Actualicé la migración `migrate_inventory_traceability.go` para dropear el índice legacy:

```go
if err := db.Exec(`DROP INDEX IF EXISTS idx_inventory_product_warehouse`).Error; err != nil {
    return fmt.Errorf("failed to drop legacy unique index: %w", err)
}
```

La migración es idempotente (`DROP IF EXISTS`) — se puede aplicar en producción sin riesgo.

### Capturas

- `adjust-with-location.png` — modal con 3 posiciones y 2 lotes detectados dinámicamente
- `adjust-full.png` — selectores de ubicación y lote visibles simultáneamente
- `adjust-success.png` — stock muestra PT01001 cantidad=30 después del primer ajuste

### Siguiente iteración

1. Agregar `level_id` al payload de `CreateLocationRequest` del backend (actualmente omitido, por eso las posiciones no quedan vinculadas a niveles específicos).
2. Agregar `track_inventory` a los productos del negocio demo para ampliar el escenario con 3 productos distintos.
3. Hacer los ajustes 2-5 vía UI (actualmente solo #1 se hizo via modal, los demás via API).
