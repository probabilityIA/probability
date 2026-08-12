# Migraciones

## Como funciona

`Migrate()` en `internal/infra/repository/constructor.go` esta **en cero**:
ejecutar `go run cmd/main.go` no migra nada.

Eso es a proposito. Antes se re-ejecutaban las 60+ migraciones historicas en cada
corrida: lento, ruidoso, y sin ninguna utilidad porque ya estaban aplicadas.

## Flujo para migrar algo

1. Escribir la migracion en su archivo (`XXX_descripcion_corta.go`).
2. Agregar la llamada dentro de `Migrate()`:

```go
func (r *Repository) Migrate(ctx context.Context) error {
    return r.migrateLoQueSea(ctx)
}
```

3. Correr `cd back/migration && go run cmd/main.go`.
4. Verificar el efecto en la base.
5. **Dejar `Migrate()` en cero otra vez** (`return nil`).
6. Registrar la corrida en la tabla de abajo.

Los DDL (AutoMigrate, `CREATE ... IF NOT EXISTS`) se conservan como archivo.
Los DML y seeds puntuales se pueden borrar una vez aplicados en produccion.

`migrateHistorico()` conserva el orden original de las migraciones ya aplicadas.
No se llama desde ningun lado; sirve como referencia y para reconstruir un
entorno desde cero si algun dia hace falta.

## Historico

| Fecha | Migracion | Que hizo | Entorno |
|-------|-----------|----------|---------|
| 2026-08-06 | `fixVigaCodCalibracionFallida` | Ajusto `cod_total` y `cod_carrier_fee` de VIG-0071, VIG-0072 y VIG-0069 al valor que liquida EnvioClick (3 ordenes) | produccion |
| 2026-08-12 | `migrateSyncRunItemParent` | Agrego `parent_ref`, `parent_label` y `variant_label` a `integration_sync_run_items` (+ indice en `parent_ref`) para agrupar variantes por publicacion en el comparativo | produccion |
| 2026-08-12 | `migrateSyncRunChannelNoSKU` | Agrego `channel_no_sku` a `integration_sync_runs` para contar los items del canal que no tienen SKU y no se pueden emparejar | produccion |
| 2026-08-12 | `migrateSyncRunSKUTypo` | Agrego `sku_typo` a `integration_sync_runs` para contar los posibles errores de digitacion detectados | produccion |
| 2026-08-12 | `migrateProductIntegrationLastPushedQty` | Agrego `last_pushed_qty` a `product_business_integrations` para no repetir el push cuando el stock no cambio | produccion |
| 2026-08-12 | `migrateSyncRunSKUChanged` | Agrego `sku_changed` a `integration_sync_runs` para contar los mapeos cuyo SKU cambio en el canal | produccion |
| 2026-08-12 | `migrateProductIntegrationLogisticType` | Agrego `external_logistic_type` a `product_business_integrations` para saber si la publicacion es de fulfillment (ML administra el stock y no se puede empujar). AutoMigrate arrastro tambien un DEFAULT '[]' en `subscription_types.features`, que ya coincidia con el modelo | produccion |
| 2026-08-12 | `migrateProductIntegrationVariantUnique` | Reemplazo el unico `idx_product_integration` de `(product_id, integration_id)` por `(product_id, integration_id, COALESCE(external_variant_id, ''))`, para que un producto pueda mapearse a varias variantes del mismo canal. Se quitaron los tags `uniqueIndex` del modelo: ahora el indice lo maneja este SQL | produccion |

Antes de esta fecha no habia registro: todas las migraciones listadas en
`migrateHistorico()` se aplicaron corriendo la cadena completa.
