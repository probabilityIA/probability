# Plan: Modulo `shipping_margins` (margen y seguro por business + carrier)

## Context

Hoy en `back/central/services/modules/shipments/internal/infra/primary/queue/consumer/response_consumer.go:659` hay un markup hardcodeado:

- `const serviceFeeAmount = 2290.0` (no 2000) se suma al `flete` de **toda** cotizacion, sin distinguir business.
- Para `interrapidisimo` y `coordinadora` (lineas 704-720), ademas del 2290 se suma el `minimumInsurance` del carrier al flete (doble cobro de seguro como ganancia extra).
- EnvioClick es intermediario (multi-carrier) y NO debe llevar margen propio: el margen aplica al carrier real de cada rate (Servientrega, Interrapidisimo, Coordinadora, MiPaquete, Enviame, etc.).

El equipo comercial necesita poder negociar un margen distinto **por business + por carrier**. Este plan reemplaza el hardcode por un CRUD configurable por super admin, con datos persistidos en DB y servidos via Redis (cache-aside) al modulo `shipments` para mantener el rendimiento del flujo de cotizacion.

## Decisiones acordadas

| Punto | Decision |
|---|---|
| Granularidad | `(business_id, carrier_code)` |
| Tipo de margen | Monto fijo COP sobre el flete |
| Seguro | Monto fijo COP sumado al `minimumInsurance` (no flag, no porcentaje) |
| Modulo | Nuevo: `/back/central/services/modules/shipping_margins/` |
| Cache | Redis cache-aside; el nuevo modulo escribe, `shipments` solo lee |
| Default sin config | margin_amount=0, insurance_margin=0 (no se cobra extra) |
| Acceso | Solo super admin (`business_id=0` + `?business_id=X`) |

## Modelo de datos

`/back/migration/shared/models/shipping_margin.go`

```go
type ShippingMargin struct {
    ID              uint      `gorm:"primaryKey"`
    BusinessID      uint      `gorm:"not null;uniqueIndex:idx_business_carrier"`
    CarrierCode     string    `gorm:"size:50;not null;uniqueIndex:idx_business_carrier"`
    CarrierName     string    `gorm:"size:100;not null"`
    MarginAmount    float64   `gorm:"not null;default:0"`
    InsuranceMargin float64   `gorm:"not null;default:0"`
    CreatedAt       time.Time
    UpdatedAt       time.Time
    DeletedAt       gorm.DeletedAt `gorm:"index"`
}
```

`carrier_code` normalizado en lowercase (`servientrega`, `interrapidisimo`, `coordinadora`, `mipaquete`, `enviame`). Lista cerrada compartida con el frontend.

Migracion: `back/migration/migrations/XXX_create_shipping_margins.go` con `AutoMigrate(&models.ShippingMargin{})`. **Sin seed automatico** — comercial define caso por caso (cumple acuerdo: si no hay registro, no se suma nada).

## Modulo backend `shipping_margins`

Estructura hexagonal estandar (referencia: `/back/central/services/modules/invoicing/`):

```
back/central/services/modules/shipping_margins/
  bundle.go
  internal/
    domain/
      entities/shipping_margin.go
      dtos/                       (request/response/pagination)
      ports/                      (repository.go, cache_writer.go)
      errors/errors.go
    app/
      constructor.go
      create_usecase.go
      update_usecase.go
      list_usecase.go
      get_usecase.go
      delete_usecase.go
      request/  response/  mappers/
    infra/
      primary/handlers/
        constructor.go  routes.go
        create_handler.go  update_handler.go  list_handler.go
        get_handler.go  delete_handler.go
        request/  response/  mappers/
      secondary/
        repository/
          constructor.go  repository.go  mappers/
        cache/
          constructor.go  redis_writer.go
```

Endpoints (montados en `/api/v1/shipping-margins`, todos requieren super admin + `?business_id=X`):

| Metodo | Ruta | Accion |
|---|---|---|
| GET | `/shipping-margins?business_id=X&page=1&page_size=10` | Lista paginada de carriers configurados |
| GET | `/shipping-margins/:id?business_id=X` | Detalle |
| POST | `/shipping-margins?business_id=X` | Crear (upsert por carrier_code) |
| PUT | `/shipping-margins/:id?business_id=X` | Actualizar |
| DELETE | `/shipping-margins/:id?business_id=X` | Soft delete |

Tras cada mutacion exitosa (POST/PUT/DELETE), el usecase invoca `cache_writer.UpsertOrInvalidate(ctx, businessID, carrierCode, margin)` que escribe Redis. El handler resuelve `business_id` con el patron de `invoicing/internal/infra/primary/handlers/constructor.go:97-112` (`resolveBusinessID`).

## Cache Redis

**Estructura:** un hash por business para evitar N keys por carrier:

```
key:    shipping_margins:business:{business_id}
type:   HASH
fields: {carrier_code} -> JSON{"margin_amount":2290,"insurance_margin":0}
TTL:    24h (refresca en cada miss)
```

**Operaciones del writer (modulo nuevo):**
- Create/Update: `HSET key carrier_code <json>` + `EXPIRE key 86400`
- Delete: `HDEL key carrier_code` (si hash queda vacio: DEL key)
- Update masivo: opcional `Refresh(businessID)` que lee DB y reescribe el hash completo

**Operaciones del reader (modulo `shipments`):**
- `Get(businessID, carrierCode)`:
  1. `HGET key carrier_code` -> hit -> retorna
  2. miss -> consulta DB (`SELECT margin_amount, insurance_margin FROM shipping_margins WHERE business_id=? AND carrier_code=? AND deleted_at IS NULL`)
  3. si DB tiene registro: rehidrata Redis (`HSET` + `EXPIRE`)
  4. si DB no tiene: retorna `Margin{0,0}` (sin extra) y NO cachea (evita poison cache si luego se crea)

## Refactor del modulo `shipments`

**Nuevo port** `back/central/services/modules/shipments/internal/domain/ports/shipping_margin_reader.go`:

```go
type ShippingMarginReader interface {
    Get(ctx context.Context, businessID uint, carrierCode string) (Margin, error)
}
type Margin struct { Amount, InsuranceAmount float64 }
```

**Nuevo adapter** `internal/infra/secondary/cache/shipping_margin_reader.go` con la logica cache-aside descrita arriba (lee Redis, fallback DB).

**Bundle** `back/central/services/modules/shipments/bundle.go`: inyectar el reader al `ResponseConsumer` (campo nuevo `marginReader ports.ShippingMarginReader`).

**Refactor de `response_consumer.go:657-747`:**

```go
// ELIMINAR: const serviceFeeAmount = 2290.0
// ELIMINAR: bloque de extraInsuranceProfit por carrier nombre

func (c *ResponseConsumer) applyServiceFeeToQuoteData(
    ctx context.Context, data map[string]interface{}, provider string, businessID uint,
) {
    rates := extractRates(data) // helper trivial existente inline
    for _, rate := range rates {
        carrierCode := normalize(rate["carrier"])  // lowercase, sin espacios
        m, err := c.marginReader.Get(ctx, businessID, carrierCode)
        if err != nil || (m.Amount == 0 && m.InsuranceAmount == 0) {
            continue
        }
        if v, ok := toFloat(rate["flete"]); ok {
            rate["flete"] = v + m.Amount
        }
        if m.InsuranceAmount > 0 {
            if v, ok := toFloat(rate["minimumInsurance"]); ok {
                rate["minimumInsurance"] = v + m.InsuranceAmount
            }
        }
    }
}
```

`businessID` ya esta disponible en el flujo (ver `resolveBusinessID` en linea 644 del mismo archivo). El llamador en linea ~319 debe pasarlo.

## Frontend

`/front/central/src/services/modules/shipping_margins/` siguiendo `services/modules/invoicing/`:

```
domain/      types.ts (ShippingMargin, CarrierOption), ports.ts
infra/       actions/ (server actions con 'use server' + revalidatePath)
             repository/
ui/
  components/ ShippingMarginsTable.tsx  ShippingMarginForm.tsx
              CarrierSelect.tsx (lista cerrada)  BusinessGate.tsx
  hooks/      useShippingMargins.ts (con businessId)
```

Pagina: `/front/central/src/app/(admin)/shipping-margins/page.tsx` (super admin + selector de business via `useBusinessesSimple`). Sin business seleccionado: gate con placeholder. Patron exactamente igual a `customers` (ver convencion en CLAUDE.md modules backend conventions).

## Archivos criticos

**Crear:**
- `back/migration/shared/models/shipping_margin.go`
- `back/migration/migrations/XXX_create_shipping_margins.go`
- `back/central/services/modules/shipping_margins/**` (modulo completo)
- `back/central/services/modules/shipments/internal/domain/ports/shipping_margin_reader.go`
- `back/central/services/modules/shipments/internal/infra/secondary/cache/shipping_margin_reader.go`
- `front/central/src/services/modules/shipping_margins/**`
- `front/central/src/app/(admin)/shipping-margins/page.tsx`

**Modificar:**
- `back/central/services/modules/shipments/internal/infra/primary/queue/consumer/response_consumer.go` (eliminar 657-747 del esquema actual y usar `marginReader`; pasar `businessID` desde el call site ~319)
- `back/central/services/modules/shipments/bundle.go` (inyectar reader)
- `back/central/cmd/internal/server/init.go` (registrar bundle del nuevo modulo)
- `back/central/cmd/internal/routes/api_routes.go` (montar rutas)

**Reusables del codebase:**
- `resolveBusinessID` -> patron en `services/modules/invoicing/internal/infra/primary/handlers/constructor.go:97-112`
- `useBusinessesSimple` -> `front/central/src/services/auth/business/ui/hooks/`
- Conexion Redis -> `back/central/shared/redis/`
- `PaginationParams` / `PaginatedResponse` -> `domain/dtos/` patron estandar

## Verificacion end-to-end

1. **Migrar DB:** `cd back/migration && go run cmd/main.go` -> verificar tabla `shipping_margins` con MCP postgres.
2. **Levantar back:** `./scripts/dev-services.sh restart backend`.
3. **CRUD via API** con super admin (token de `.env.ai`):
   - `POST /api/v1/shipping-margins?business_id=1` body `{carrier_code:"servientrega",carrier_name:"Servientrega",margin_amount:2000,insurance_margin:0}`
   - `GET /api/v1/shipping-margins?business_id=1` confirma listado
4. **Validar Redis:** `redis-cli HGETALL shipping_margins:business:1` debe traer el JSON.
5. **Cotizar shipment** desde `/api/v1/shipments/quote` para business 1 con carrier Servientrega -> log "Service fee added" debe mostrar `margin=2000` (no 2290) y `final_flete = original + 2000`.
6. **Cache miss test:** `redis-cli DEL shipping_margins:business:1` -> volver a cotizar -> debe leer DB y rehidratar (verificar HGETALL otra vez).
7. **Sin config:** business 2 sin registros -> cotizar -> `flete` y `minimumInsurance` SIN cambios (margen 0, no se modifica el rate).
8. **Caso seguro:** crear `interrapidisimo` con `insurance_margin=500` -> cotizar -> `minimumInsurance = original + 500`.
9. **E2E frontend:** crear/editar/eliminar via UI con selector de business; mutaciones reflejadas en DB y Redis.
10. **Build + tests:** `go build ./...` y `go test ./...` en `/back/central` y `/back/migration`.

## Fuera de alcance

- Migracion automatica de "businesses existentes con 2290". Si comercial decide preservar la ganancia historica para algun business, debe crear los registros manualmente via la nueva UI.
- Margen porcentual / por nivel de servicio / por ruta. Solo monto fijo COP por carrier.
- Auditoria de cambios (quien cambio que cuando). Si se requiere, abrir issue separado.
