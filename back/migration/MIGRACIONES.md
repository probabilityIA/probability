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

Antes de esta fecha no habia registro: todas las migraciones listadas en
`migrateHistorico()` se aplicaron corriendo la cadena completa.
