# Convenciones Backend - Go

## 1. Aislamiento de Repositorios

Repos NO se comparten entre modulos. Si modulo A necesita datos de modulo B, replicar SOLO metodos SELECT en repo propio.

```go
// NUNCA:
import paymentstatusrepo "github.com/.../paymentstatus/infra/secondary/repository"

// Correcto: implementar localmente en status_queries.go
func (r *Repository) GetPaymentStatusIDByCode(ctx context.Context, code string) (*uint, error) {
    var result struct{ ID uint }
    err := r.db.Conn(ctx).Table("payment_statuses").Select("id").
        Where("code = ? AND deleted_at IS NULL", code).Limit(1).First(&result).Error
    if err == gorm.ErrRecordNotFound { return nil, nil }
    return &result.ID, err
}
```

Solo replicar CONSULTAS (GetByID, GetByCode, FindBy...). Comunicacion compleja entre modulos: RabbitMQ.

## 2. Migraciones

TODAS desde `/back/migration`. `cd /back/migration && go run cmd/main.go`
- Modelos GORM en `migration/shared/models/` = fuente de verdad. NUNCA `models/` en modulos. NUNCA `ALTER TABLE` desde modulos.
- Solo AutoMigrate el modelo que cambio.
- DDL: idempotente (`IF NOT EXISTS`), se mantienen. DML/seeds: se eliminan despues de produccion.
- Nomenclatura: `XXX_descripcion_corta.go`

## 3. Logging

zerolog, sistema dual:
- Normal: `.Info()`, `.Warn()`, `.Error()` -> consola (flujo operacional)
- Debug: `.DebugToFile()` -> `/back/central/log/app-YYYY-MM-DD.log` (activar: `ENABLE_DEBUG_FILE_LOGGING=true`)

## 4. Gestion de Procesos Backend

NUNCA iniciar/reiniciar/detener backend sin permiso explicito del usuario.
- Siempre: modificar codigo, compilar (`go build -o /tmp/test cmd/main.go`), matar zombies (`pkill -9 go`)
- Solo con permiso: iniciar (SIEMPRE foreground, NUNCA `&` ni `nohup`), reiniciar (`./scripts/dev-services.sh restart backend`)

## 5. Super Admin - Business ID

Super admins tienen `business_id = 0` en JWT. Requieren `?business_id=X` en query param; sin el = 400.

```go
func (h *Handlers) resolveBusinessID(c *gin.Context) (uint, bool) {
    businessID := c.GetUint("business_id")
    if businessID > 0 { return businessID, true }
    if param := c.Query("business_id"); param != "" {
        if id, err := strconv.ParseUint(param, 10, 64); err == nil && id > 0 {
            return uint(id), true
        }
    }
    return 0, false
}
```

POST/PUT/DELETE: `business_id` en query param, no en body.

**Frontend:** `isSuperAdmin` -> selector obligatorio -> sin negocio = gate/placeholder -> pasar `business_id` a todas las operaciones -> resetear al cambiar negocio -> usar `useBusinessesSimple` de `@/services/auth/business/ui/hooks/`

Modulos implementados: orders, invoicing, customers. Referencia: `services/modules/customers/`

## 6. Creacion de Productos desde Integraciones (regla de arquitectura)

Una integracion NUNCA crea productos escribiendo la tabla `products` ni llamando
al modulo `products` directamente. La UNICA via para materializar productos en
Probability desde una integracion es **publicar a la cola**
`rabbitmq.QueueProductsProviderUpsert` (`products.provider_upsert.requests`).

El mensaje debe incluir SIEMPRE el `integration_id`:

```go
type productUpsertMessage struct {
    BusinessID     uint    `json:"business_id"`
    IntegrationID  uint    `json:"integration_id"` // OBLIGATORIO
    SKU            string  `json:"sku"`
    Name           string  `json:"name"`
    TrackInventory bool    `json:"track_inventory"`
    Price          float64 `json:"price"`
    ExternalID     string  `json:"external_id"`
}
```

El consumer del modulo `products` (`UpsertFromProvider`) crea/actualiza el producto
y, con el `integration_id`, crea la relacion producto<->canal en
`product_business_integrations` (idempotente). Asi todo el catalogo queda asociado
a su canal de origen, lo que habilita a futuro empujar inventario solo a los
productos asociados a cada canal (ver `.claude/alerts/inventario-saliente-por-canal.md`).

- Publican HOY por esta via: WooCommerce, MercadoLibre, Siigo (manual y auto-sync).
- Las integraciones SI pueden escribir directamente `product_business_integrations`
  (mapping), pero la creacion del PRODUCTO va por la cola.
- Los repos de integraciones sobre `products` deben ser solo LECTURA (SELECT).

**Violacion critica:** integracion que crea/actualiza filas en `products` sin pasar
por la cola de upsert (rompe la asociacion producto<->canal y el aislamiento de modulos).
