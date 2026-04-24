# Test Data - Inventario

## Usuario

- Email: `demo@probability.com`
- Password: `ProbabilityDemo`
- Business ID: 26

## Escenario de pruebas (CU-01)

### Bodega a crear

- Código: `TEST-WH-01`
- Nombre: `Bodega Test Completa`
- Ciudad: `Medellin`

### Jerarquía a construir

```
Bodega TEST-WH-01
├── Zona Z-PICK (Picking, color #10b981)
│   ├── Pasillo A-01 (Pasillo Frontal)
│   │   ├── Rack R-01 (Rack Principal, 3 niveles)
│   │   │   ├── Nivel L-01 (ordinal 1)
│   │   │   ├── Nivel L-02 (ordinal 2)
│   │   │   └── Nivel L-03 (ordinal 3)
│   │   └── Rack R-02 (Rack Secundario, 2 niveles)
│   │       ├── Nivel L-01
│   │       └── Nivel L-02
│   └── Pasillo A-02 (Pasillo Trasero)
│       └── Rack R-03 (Rack Corto, 1 nivel)
│           └── Nivel L-01
└── Zona Z-BULK (Bulk, color #f59e0b)
    └── Pasillo B-01 (Bulk principal)
        └── Rack R-04 (Rack Bulk Alto, 2 niveles)
            ├── Nivel L-01
            └── Nivel L-02
```

Resumen: 2 zonas, 4 pasillos, 4 racks, 9 niveles.

### Posiciones a crear

Desde `/warehouses/:id` tras construir la jerarquía. Cada posición va asociada a un nivel:

- POS-A0101 (Zona Picking, Rack R-01, Nivel L-01, tipo shelf)
- POS-A0201 (Zona Picking, Rack R-03, Nivel L-01, tipo pallet)
- POS-B0101 (Zona Bulk, Rack R-04, Nivel L-02, tipo pallet, alta capacidad)

### Productos a usar

Reutilizar 3 productos existentes del business 26 con `track_inventory=true`. Si no tienen tracking, habilitarlo vía consulta directa a la DB (sólo para propósitos de test).

Queries para identificar productos candidatos:

```sql
SELECT id, sku, name, track_inventory
FROM products
WHERE business_id = 26 AND track_inventory = true
LIMIT 3;
```

Si no hay, habilitar en 3 productos:

```sql
UPDATE products
SET track_inventory = true
WHERE business_id = 26
  AND id IN (
      SELECT id FROM products WHERE business_id = 26 ORDER BY id LIMIT 3
  );
```

### Lotes a crear

Para el producto 1 (alimento/químico):

- `LOT-TEST-A01` con vencimiento +60 días
- `LOT-TEST-A02` con vencimiento +120 días

Producto 2 y 3 → sin lotes (stock general).

### Movimientos esperados

1. **Ajuste positivo** Producto 1 en POS-A0101 con LOT-TEST-A01 → qty +30
2. **Ajuste positivo** Producto 1 en POS-A0201 con LOT-TEST-A02 → qty +50
3. **Ajuste positivo** Producto 2 en POS-B0101 → qty +100 (sin lote)
4. **Ajuste positivo** Producto 3 sin ubicación → qty +20 (stock general)
5. **Ajuste negativo** Producto 1 en POS-A0101 con LOT-TEST-A01 → qty -5 (merma)

Verificación final en DB:

```sql
SELECT il.id, p.sku, w.code as warehouse, wl.code as location, l.lot_code as lot,
       s.code as state, il.quantity, il.available_qty
FROM inventory_levels il
JOIN products p ON p.id = il.product_id
JOIN warehouses w ON w.id = il.warehouse_id
LEFT JOIN warehouse_locations wl ON wl.id = il.location_id
LEFT JOIN inventory_lots l ON l.id = il.lot_id
LEFT JOIN inventory_states s ON s.id = il.state_id
WHERE w.code = 'TEST-WH-01'
ORDER BY il.id;
```
