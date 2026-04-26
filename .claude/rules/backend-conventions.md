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
