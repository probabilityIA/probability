# Redis Registry - Uso y Ejemplos

## Descripción

El sistema de Redis Registry permite registrar y visualizar en los logs de startup:
1. **Prefijos de caché** - Patrones de keys que usa cada módulo
2. **Canales pub/sub** - Canales activos para comunicación en tiempo real

## Cómo Usar

### 1. Registrar Prefijos de Caché

Cuando un módulo utiliza Redis para caché, debe registrar su prefijo al inicializarse:

```go
// En el bundle.go de tu módulo
func New(router *gin.RouterGroup, database db.IDatabase, redisClient redis.IRedis, logger log.ILogger, config env.IConfig) {
    // Registrar el prefijo de caché que usará este módulo
    if redisClient != nil {
        redisClient.RegisterCachePrefix("probability:invoicing:config:*")
        redisClient.RegisterCachePrefix("probability:invoicing:retry:*")
    }

    // ... resto de la inicialización
}
```

### 2. Registrar Canales Pub/Sub

Si tu módulo usa Redis pub/sub para comunicación en tiempo real:

```go
// En el consumer o publisher que usa canales
func NewOrderEventPublisher(redisClient redis.IRedis, logger log.ILogger) *OrderEventPublisher {
    if redisClient != nil {
        // Registrar canales que este publisher usará
        redisClient.RegisterChannel("orders:created")
        redisClient.RegisterChannel("orders:updated")
        redisClient.RegisterChannel("orders:cancelled")
    }

    return &OrderEventPublisher{
        redis: redisClient,
        log:   logger,
    }
}
```

### 3. Ejemplo Completo - Módulo de Caché de Configuración

```go
// services/modules/invoicing/bundle.go
package invoicing

import (
    "github.com/gin-gonic/gin"
    "github.com/secamc93/probability/back/central/shared/db"
    "github.com/secamc93/probability/back/central/shared/redis"
    "github.com/secamc93/probability/back/central/shared/log"
    "github.com/secamc93/probability/back/central/shared/env"
)

func New(
    router *gin.RouterGroup,
    database db.IDatabase,
    redisClient redis.IRedis,
    logger log.ILogger,
    config env.IConfig,
) {
    // 1. REGISTRAR PREFIJOS DE CACHÉ al inicio del módulo
    if redisClient != nil {
        // Registrar todos los patrones de keys que este módulo usará
        redisClient.RegisterCachePrefix("probability:invoicing:config:*")
        redisClient.RegisterCachePrefix("probability:invoicing:retry:*")
        redisClient.RegisterCachePrefix("probability:invoicing:stats:*")
    }

    // 2. Inicializar repositorio (que usará estos prefijos)
    repo := repository.New(database, redisClient, config, logger)

    // 3. Inicializar use cases
    useCase := usecases.New(repo, logger, config)

    // ... resto del bundle
}
```

### 4. Ejemplo - Repository con Caché

```go
// services/modules/invoicing/internal/infra/secondary/repository/config_cache.go
package repository

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
)

const (
    // Prefijo registrado en bundle.go
    ConfigCachePrefix = "probability:invoicing:config"
)

func (r *Repository) getConfigCache(ctx context.Context, integrationID uint) (*entities.InvoicingConfig, error) {
    if r.redis == nil {
        return nil, nil
    }

    // Usar el prefijo registrado
    key := fmt.Sprintf("%s:%d", ConfigCachePrefix, integrationID)

    data, err := r.redis.Get(ctx, key)
    if err != nil {
        return nil, nil // Cache miss
    }

    var config entities.InvoicingConfig
    if err := json.Unmarshal([]byte(data), &config); err != nil {
        return nil, err
    }

    return &config, nil
}

func (r *Repository) setConfigCache(ctx context.Context, config *entities.InvoicingConfig) error {
    if r.redis == nil {
        return nil
    }

    key := fmt.Sprintf("%s:%d", ConfigCachePrefix, config.IntegrationID)
    data, err := json.Marshal(config)
    if err != nil {
        return err
    }

    ttl := 1 * time.Hour
    return r.redis.Set(ctx, key, data, ttl)
}
```

## Output en Logs de Startup

Cuando el servidor inicia, verás algo como:

```
 🚀 Servidor HTTP iniciado correctamente
 📍 Disponible en: http://localhost:8000
 📖 Documentación: http://localhost:8000/docs/index.html

 🗄️  Conexión PostgreSQL: postgres://localhost:5433/probability

 🐰 RabbitMQ: amqp://localhost:5672/
    📥 Colas activas:
       • order.created
       • invoice.sync
       • invoice.retry

 🔴 Redis: redis://localhost:6379
    💾 Prefijos de caché:
       • probability:invoicing:config:*
       • probability:invoicing:retry:*
       • probability:orders:*
       • probability:sessions:*
    📡 Canales pub/sub:
       • orders:created
       • orders:updated
       • invoices:synced

 ☁️  S3 Storage: s3://probability-bucket (us-east-1)
```

## Beneficios

1. **Visibilidad** - Ver qué módulos usan Redis y para qué
2. **Debugging** - Identificar rápidamente patrones de keys en uso
3. **Documentación automática** - Los logs sirven como documentación viva
4. **Optimización** - Detectar prefijos redundantes o mal organizados

## Convenciones

### Nomenclatura de Prefijos

Seguir el patrón:
```
probability:{modulo}:{tipo}:*
```

Ejemplos:
- `probability:invoicing:config:*` - Configuraciones de facturación
- `probability:orders:cache:*` - Caché de órdenes
- `probability:sessions:user:*` - Sesiones de usuario
- `probability:analytics:daily:*` - Estadísticas diarias

### Nomenclatura de Canales

Seguir el patrón:
```
{entidad}:{evento}
```

Ejemplos:
- `orders:created` - Nueva orden creada
- `orders:updated` - Orden actualizada
- `invoices:synced` - Factura sincronizada
- `payments:completed` - Pago completado

## Cuándo Registrar

- **Al inicio del módulo** - En el `bundle.go` o `constructor.go`
- **Una sola vez** - No repetir registros en cada operación
- **Solo lo que se usa** - No registrar prefijos "por las dudas"
- **Antes de usar** - Registrar antes de hacer Set/Get con ese prefijo

## Cuándo NO Registrar

- ❌ Keys temporales o únicas (session tokens, CSRF tokens)
- ❌ Prefijos dinámicos que cambian constantemente
- ❌ Keys de testing o desarrollo
- ❌ Patrones que no seguirán usándose

---

**Última actualización:** 2026-02-10
**Propósito:** Documentar el uso del sistema de tracking de Redis
