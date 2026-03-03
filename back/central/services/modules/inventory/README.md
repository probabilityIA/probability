# Inventory Module

Gestiona niveles de inventario, movimientos de stock y operaciones transaccionales vinculadas al ciclo de vida de las ordenes.

## Arquitectura

```
inventory/
├── bundle.go                           # Composicion del modulo
└── internal/
    ├── domain/
    │   ├── entities/                   # InventoryLevel, StockMovement, StockMovementType
    │   ├── dtos/                       # Params de consulta, DTOs transaccionales, DTOs de ordenes
    │   ├── ports/                      # IRepository, ISyncPublisher, IInventoryEventPublisher
    │   └── errors/                     # Errores de dominio
    ├── app/                            # Casos de uso
    │   ├── constructor.go
    │   ├── list_inventory.go           # Listar stock por bodega
    │   ├── get_inventory.go            # Stock de un producto en todas las bodegas
    │   ├── adjust_stock.go             # Ajuste manual
    │   ├── transfer_stock.go           # Transferencia entre bodegas
    │   ├── list_movements.go           # Historial de movimientos
    │   ├── reserve_stock.go            # Reservar stock por orden nueva
    │   ├── confirm_sale.go             # Confirmar venta (shipped/completed)
    │   ├── release_stock.go            # Liberar reserva (cancelled)
    │   ├── return_stock.go             # Devolver stock (refunded)
    │   ├── list_movement_types.go      # Listar tipos de movimiento
    │   ├── create_movement_type.go     # Crear tipo de movimiento
    │   ├── update_movement_type.go     # Actualizar tipo
    │   └── delete_movement_type.go     # Eliminar tipo
    └── infra/
        ├── primary/
        │   ├── handlers/               # HTTP handlers (Gin)
        │   │   ├── routes.go
        │   │   ├── request/            # Structs de request HTTP
        │   │   └── response/           # Structs de response HTTP
        │   └── queue/
        │       └── order_consumer.go   # Consumer RabbitMQ (orders.events.inventory)
        └── secondary/
            ├── repository/             # GORM + transacciones atomicas
            │   ├── inventory_level_queries.go
            │   ├── adjust_stock_tx.go
            │   ├── transfer_stock_tx.go
            │   ├── reserve_stock_tx.go
            │   ├── confirm_sale_tx.go
            │   ├── release_stock_tx.go
            │   ├── return_stock_tx.go
            │   ├── stock_movements.go
            │   ├── movement_type_queries.go
            │   └── product_integration_queries.go
            ├── redis/
            │   ├── inventory_cache.go  # Cache de niveles (TTL)
            │   └── event_publisher.go  # Pub/sub para SSE
            └── queue/
                └── sync_publisher.go   # Publicar sync a integraciones externas
```

## Endpoints HTTP

| Metodo | Ruta | Descripcion |
|--------|------|-------------|
| GET | `/inventory/product/:productId` | Stock de un producto en todas las bodegas |
| GET | `/inventory/warehouse/:warehouseId` | Niveles de stock de una bodega (paginado) |
| POST | `/inventory/adjust` | Ajustar stock manualmente |
| POST | `/inventory/transfer` | Transferir stock entre bodegas |
| GET | `/inventory/movements` | Listar movimientos de stock (paginado) |
| GET | `/inventory/movement-types` | Listar tipos de movimiento |
| POST | `/inventory/movement-types` | Crear tipo de movimiento |
| PUT | `/inventory/movement-types/:id` | Actualizar tipo |
| DELETE | `/inventory/movement-types/:id` | Eliminar tipo |

Todos los endpoints soportan super admin con `?business_id=X` query param.

### Query Params comunes

**GET /inventory/warehouse/:warehouseId:**
- `page`, `page_size` — Paginacion
- `search` — Buscar por nombre/SKU de producto
- `low_stock` — Filtrar productos con stock bajo (`true`)
- `business_id` — Requerido para super admin

**GET /inventory/movements:**
- `page`, `page_size` — Paginacion
- `warehouse_id` — Filtrar por bodega
- `product_id` — Filtrar por producto
- `movement_type_id` — Filtrar por tipo de movimiento
- `business_id` — Requerido para super admin

## Consumer RabbitMQ

Cola: `orders.events.inventory`

Escucha eventos del ciclo de vida de ordenes y ejecuta operaciones de inventario automaticas:

| Evento | Accion |
|--------|--------|
| `order.created` | Reservar stock de cada item en la bodega por defecto |
| `order.cancelled` | Liberar stock reservado |
| `order.shipped` | Confirmar venta (reserved → sold) |
| `order.completed` | Confirmar venta |
| `order.refunded` | Devolver stock al inventario |
| `order.status_changed` | Enrutar segun keywords del estado |

### Flujo de stock por orden

```
Orden creada → quantity reservada (reserved_qty += N, available_qty -= N)
    │
    ├─ Orden cancelada → reserva liberada (reserved_qty -= N, available_qty += N)
    │
    └─ Orden enviada/completada → venta confirmada (reserved_qty -= N, quantity -= N)
         │
         └─ Orden reembolsada → stock devuelto (quantity += N, available_qty += N)
```

## Operaciones Transaccionales

Todas las operaciones que modifican stock usan `SELECT FOR UPDATE` + commit atomico para evitar race conditions:

- **AdjustStockTx** — Ajuste manual con movimiento de auditoria
- **TransferStockTx** — Descuenta de origen, suma en destino (2 movimientos)
- **ReserveStockTx** — Valida available_qty, reserva por cada item de la orden
- **ConfirmSaleTx** — Convierte reserva en venta confirmada
- **ReleaseStockTx** — Libera reserva al cancelar
- **ReturnStockTx** — Reingresa stock al inventario

## Cache (Redis)

- Los niveles de inventario se cachean por `product:{productID}:warehouse:{warehouseID}`
- TTL configurable via env
- Se invalida automaticamente al modificar stock
- Resiliente: si Redis no esta disponible, las queries van directo a PostgreSQL

## Sync Publisher

Despues de cada operacion de stock, si el producto tiene integraciones vinculadas, se publica un mensaje a RabbitMQ para sincronizar el inventario en los canales externos (Shopify, MercadoLibre, etc).

## Entidades

### InventoryLevel

```
ProductID     string  — UUID del producto
WarehouseID   uint    — Bodega
LocationID    *uint   — Ubicacion dentro de la bodega (opcional)
BusinessID    uint    — Negocio
Quantity      int     — Stock total
ReservedQty   int     — Reservado por ordenes pendientes
AvailableQty  int     — Disponible (quantity - reserved)
MinStock      *int    — Umbral minimo (alerta)
MaxStock      *int    — Capacidad maxima
ReorderPoint  *int    — Punto de reorden
```

### StockMovement

```
ProductID       string  — UUID del producto
WarehouseID     uint    — Bodega
MovementTypeID  uint    — Tipo de movimiento (FK → stock_movement_types)
Quantity        int     — Cantidad (positivo = entrada, negativo = salida)
PreviousQty     int     — Stock antes del movimiento
NewQty          int     — Stock despues del movimiento
Reason          string  — Razon del movimiento
ReferenceType   *string — Tipo de referencia (order, transfer, adjustment)
ReferenceID     *string — ID de la referencia
Notes           string  — Notas adicionales
```

### StockMovementType (catalogo)

| ID | Code | Nombre | Direccion |
|----|------|--------|-----------|
| 1 | inbound | Entrada de mercancia | in |
| 2 | outbound | Salida de mercancia | out |
| 3 | adjustment | Ajuste de inventario | neutral |
| 4 | transfer | Transferencia entre bodegas | neutral |
| 5 | return | Devolucion | in |
| 6 | sync | Sincronizacion desde canal | neutral |
| 7 | reserve | Reserva de stock | neutral |
| 8 | confirm_sale | Confirmacion de venta | out |
| 9 | release | Liberacion de reserva | neutral |

## Dependencias

```go
inventory.New(router, database, logger, environment, rabbitMQ, redisClient)
```

- `db.IDatabase` — PostgreSQL via GORM
- `log.ILogger` — Zerolog
- `env.IConfig` — Variables de entorno
- `rabbitmq.IQueue` — RabbitMQ (consumer + publisher)
- `redis.IRedis` — Redis (cache + event pub/sub, opcional)

## Tests

```bash
go test ./internal/app/...
go test ./internal/infra/primary/handlers/...
```

Tests unitarios con mocks para repository, logger, publisher y event publisher.
