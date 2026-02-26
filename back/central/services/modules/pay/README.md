# Módulo `pay` — Pagos con Pasarelas Externas

Módulo de pagos de primera clase para procesar transacciones reales a través de pasarelas externas (Nequi, y futuras). Gestiona el ciclo de vida completo de una transacción: creación, procesamiento asíncrono vía colas, reintentos y notificaciones en tiempo real (SSE).

> **Distinto de `wallet`:** `wallet` maneja recargas manuales con aprobación de admin. `pay` procesa pagos reales contra APIs externas.

---

## Flujo de una Transacción

```
Frontend
  │
  ▼
POST /pay/transactions
  │
  ▼
CreatePayment (UseCase)
  ├─► Crear PaymentTransaction en DB  (status: pending)
  ├─► Crear PaymentSyncLog           (status: processing)
  └─► Publicar a "pay.requests"
          │
          ▼
  integrations/pay/router
  (enruta según gateway_code → "pay.nequi.requests")
          │
          ▼
  integrations/pay/nequi
  (llama API Nequi → publica resultado a "pay.responses")
          │
          ▼
ResponseConsumer (este módulo)
  ├─► Actualizar PaymentTransaction  (status: completed | failed)
  ├─► Actualizar PaymentSyncLog
  └─► SSE vía Redis Pub/Sub → Frontend
```

---

## Estructura de Carpetas

```
pay/
├── bundle.go                            # Punto de entrada — ensambla todo
└── internal/
    ├── domain/
    │   ├── entities/
    │   │   ├── payment_transaction.go   # Entidad principal (sin tags)
    │   │   └── payment_sync_log.go      # Registro de cada intento
    │   ├── dtos/
    │   │   ├── create_payment_dto.go    # DTO de entrada
    │   │   └── payment_message.go       # Mensajes RabbitMQ (request/response)
    │   ├── ports/
    │   │   └── ports.go                 # IRepository, IUseCase, IRequestPublisher, ISSEPublisher
    │   ├── constants/
    │   │   └── constants.go             # Estados, gateways, nombres de colas
    │   └── errors/
    │       └── errors.go                # Errores tipados del dominio
    ├── app/
    │   ├── constructor.go               # Factory: New() → IUseCase
    │   ├── create_payment.go            # Iniciar pago
    │   ├── process_payment_response.go  # Procesar respuesta del gateway
    │   ├── retry_payment.go             # Reintentar pago fallido
    │   ├── get_payment.go               # Consultar transacción por ID
    │   └── list_payments.go             # Listar con paginación
    └── infra/
        ├── primary/
        │   ├── handlers/
        │   │   ├── constructor.go
        │   │   ├── routes.go
        │   │   ├── create_payment_handler.go
        │   │   ├── get_payment_handler.go
        │   │   ├── list_payments_handler.go
        │   │   ├── request/
        │   │   │   └── create_payment.go
        │   │   ├── response/
        │   │   │   └── payment.go
        │   │   └── mappers/
        │   │       └── mapper.go
        │   └── queue/consumer/
        │       ├── constructor.go
        │       ├── response_consumer.go  # Consume "pay.responses"
        │       └── retry_consumer.go     # Cron cada ~5 min
        └── secondary/
            ├── repository/
            │   ├── constructor.go
            │   └── repository.go         # GORM — usa migration/shared/models
            ├── queue/
            │   └── request_publisher.go  # Publica a "pay.requests"
            └── redis/
                └── sse_publisher.go      # Redis Pub/Sub para SSE
```

---

## Endpoints HTTP

| Método | Ruta | Descripción |
|--------|------|-------------|
| `POST` | `/pay/transactions` | Crear nueva transacción de pago |
| `GET` | `/pay/transactions` | Listar transacciones (paginado) |
| `GET` | `/pay/transactions/:id` | Obtener transacción por ID |

### POST /pay/transactions — Request

```json
{
  "amount": 50000,
  "currency": "COP",
  "gateway_code": "nequi",
  "payment_method": "qr_code",
  "description": "Recarga de saldo",
  "callback_url": "https://example.com/callback",
  "metadata": {}
}
```

### POST /pay/transactions — Response `201`

```json
{
  "id": 42,
  "business_id": 5,
  "reference": "a1b2c3d4e5f6...",
  "amount": 50000,
  "currency": "COP",
  "status": "pending",
  "gateway_code": "nequi",
  "payment_method": "qr_code",
  "created_at": "2026-02-26T10:00:00Z"
}
```

### GET /pay/transactions — Query params

| Param | Default | Descripción |
|-------|---------|-------------|
| `page` | `1` | Página |
| `page_size` | `10` | Registros por página (max 100) |

---

## Colas RabbitMQ

| Cola | Dirección | Descripción |
|------|-----------|-------------|
| `pay.requests` | Publica | Solicitud de pago hacia el router de integraciones |
| `pay.responses` | Consume | Resultado del gateway → actualiza DB + SSE |

---

## Estados de una Transacción

```
pending ──► processing ──► completed
                │
                └──► failed (max 3 reintentos)
                       │
                       └──► cancelled (si se agotaron reintentos)
```

| Estado | Descripción |
|--------|-------------|
| `pending` | Recién creada, esperando procesamiento |
| `processing` | En cola / siendo procesada por el gateway |
| `completed` | Pago exitoso |
| `failed` | Error — puede reintentarse |
| `cancelled` | Agotó los reintentos máximos |

---

## Reintentos Automáticos

`RetryConsumer` corre cada ~5 minutos y busca `PaymentSyncLog` con:
- `status = 'failed'`
- `retry_count < 3`
- `next_retry_at <= now()`

Por cada registro encontrado:
1. Cancela sync logs pendientes anteriores
2. Crea un nuevo sync log con `retry_count + 1`
3. Re-publica el mensaje a `pay.requests`

Máximo **3 reintentos**. Al agotarse, la transacción queda en `cancelled`.

---

## SSE (Tiempo Real)

Cuando una transacción cambia de estado, el módulo publica en el canal Redis:

```
probability:pay:state:events
```

Payload del evento:

```json
{
  "event": "payment.completed",
  "transaction_id": 42,
  "reference": "a1b2c3d4...",
  "status": "completed",
  "gateway_code": "nequi",
  "external_id": "NQI-2026-XYZ",
  "timestamp": "2026-02-26T10:01:00Z"
}
```

Eventos posibles: `payment.completed`, `payment.failed`, `payment.processing`

---

## Modelos de Base de Datos

> Definidos en `back/migration/shared/models/pay.go`

### `payment_transaction`

| Columna | Tipo | Descripción |
|---------|------|-------------|
| `id` | bigint PK | Auto-incremental |
| `business_id` | bigint | FK → businesses (multi-tenant) |
| `amount` | numeric | Monto (ej: `50000` = $50.000 COP) |
| `currency` | varchar | `"COP"` |
| `status` | varchar | `pending\|processing\|completed\|failed\|cancelled` |
| `gateway_code` | varchar | `"nequi"` |
| `external_id` | varchar? | ID de la transacción en el gateway |
| `reference` | varchar UNIQUE | Referencia interna (UUID hex) |
| `payment_method` | varchar? | `"qr_code"\|"payment_link"` |
| `description` | varchar? | Descripción del pago |
| `callback_url` | varchar? | URL de retorno |
| `metadata` | jsonb? | Datos adicionales |
| `gateway_response` | jsonb? | Respuesta cruda del gateway |

### `payment_sync_log`

| Columna | Tipo | Descripción |
|---------|------|-------------|
| `payment_transaction_id` | bigint FK | Transacción asociada |
| `status` | varchar | Estado de este intento |
| `retry_count` | int | Número de intento (0 = primero) |
| `gateway_request` | jsonb? | Request enviado al gateway |
| `gateway_response` | jsonb? | Respuesta del gateway |
| `error_message` | varchar? | Mensaje de error legible |
| `next_retry_at` | timestamptz? | Cuándo reintentar |

---

## Gateways Soportados

| Gateway | `gateway_code` | Estado |
|---------|---------------|--------|
| Nequi | `nequi` | Implementado |
| Wompi | `wompi` | Pendiente |
| PayU | `payU` | Pendiente |
| Stripe | `stripe` | Pendiente |
| ePayco | `ePayco` | Pendiente |
| Bold | `bold` | Pendiente |
| MercadoPago | `meliPago` | Pendiente |

---

## Variables de Entorno

No requiere variables propias. Usa las compartidas del proyecto:

```env
# RabbitMQ (requerido para processing)
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672

# Redis (requerido para SSE)
REDIS_HOST=localhost
REDIS_PORT=6379

# Database
DB_HOST=localhost
DB_PORT=5433
```
