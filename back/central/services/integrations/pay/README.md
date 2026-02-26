# Integración `pay` — Pasarelas de Pago

Capa de integración que conecta el módulo `modules/pay` con las APIs externas de pasarelas de pago. Actúa como broker: recibe solicitudes genéricas de la cola `pay.requests`, las enruta a la pasarela correspondiente, ejecuta el pago y devuelve el resultado a `pay.responses`.

---

## Arquitectura General

```
"pay.requests" (genérico)
        │
        ▼
  ┌─────────────┐
  │   router/   │  ← Lee gateway_code → enruta
  └─────────────┘
        │
        ├──► "pay.nequi.requests"  ──► nequi/
        ├──► "pay.wompi.requests"  ──► wompi/   (pendiente)
        ├──► "pay.stripe.requests" ──► stripe/  (pendiente)
        └──► ...
                │
                ▼
        API externa (Nequi, etc.)
                │
                ▼
        "pay.responses" (resultado)
                │
                ▼
        modules/pay ResponseConsumer
```

---

## Estructura de Carpetas

```
integrations/pay/
├── bundle.go               # Ensambla nequi + router
├── nequi/                  # Pasarela Nequi (implementada)
│   ├── bundle.go
│   └── internal/
│       ├── domain/
│       │   ├── entities/
│       │   │   └── nequi_payment.go        # NequiPaymentResult (sin tags)
│       │   ├── ports/
│       │   │   └── ports.go                # INequiClient, IIntegrationRepository, IResponsePublisher
│       │   └── errors/
│       │       └── errors.go
│       ├── app/
│       │   ├── constructor.go              # New() → IUseCase
│       │   └── process_payment.go          # Lógica principal
│       └── infra/
│           ├── primary/consumer/
│           │   └── nequi_consumer.go       # Consume "pay.nequi.requests"
│           └── secondary/
│               ├── client/
│               │   └── nequi_client.go     # HTTP client resty → API Nequi
│               ├── queue/
│               │   └── response_publisher.go # Publica a "pay.responses"
│               └── repository/
│                   └── integration_repository.go # Lee credenciales de DB
├── router/
│   └── bundle.go           # Consume "pay.requests" → enruta al gateway correcto
├── wompi/                  # (pendiente)
├── stripe/                 # (pendiente)
├── payU/                   # (pendiente)
├── ePayco/                 # (pendiente)
├── bold/                   # (pendiente)
└── meliPago/               # (pendiente)
```

---

## Colas RabbitMQ

| Cola | Productor | Consumidor | Descripción |
|------|-----------|------------|-------------|
| `pay.requests` | `modules/pay` | `router/` | Solicitud de pago genérica |
| `pay.nequi.requests` | `router/` | `nequi/` | Solicitud específica para Nequi |
| `pay.responses` | `nequi/` | `modules/pay` | Resultado del gateway |

---

## Router (`router/`)

Consume la cola genérica `pay.requests` y reenvía el mensaje a la cola específica del gateway indicado en el campo `gateway_code`.

```
Mensaje de pay.requests:
{
  "gateway_code": "nequi",   ← campo de ruteo
  "payment_transaction_id": 42,
  "amount": 50000,
  "reference": "a1b2c3...",
  ...
}

→ publica a: "pay.nequi.requests"
```

**Diseño:** el router no conoce las implementaciones de los gateways, solo los nombres de las colas. Para agregar un nuevo gateway basta con añadir un `case` en el switch de ruteo.

---

## Nequi (`nequi/`)

Implementación completa de la integración con la API de Nequi Colombia.

### Flujo interno

```
"pay.nequi.requests"
        │
        ▼
nequi_consumer.go
        │
        ▼
process_payment.go (UseCase)
  ├─► GetNequiConfig()           → lee credenciales de integration_types
  ├─► nequiClient.GenerateQR()  → POST a API Nequi
  └─► responsePublisher.Publish() → "pay.responses" (éxito o error)
```

### Credenciales

Se leen de la tabla `integration_types` donde `code = 'nequi_pay'`, campo `platform_credentials_encrypted` (AES-256-GCM). La clave de cifrado viene de la variable de entorno `ENCRYPTION_KEY`.

**Estructura de credenciales** (almacenadas cifradas):

```json
{
  "api_key":     "sk_live_...",
  "environment": "sandbox | production",
  "phone_code":  "NIT_1"
}
```

### API Nequi — GenerateQR

| Parámetro | Descripción |
|-----------|-------------|
| Endpoint sandbox | `https://sandbox.nequi.com.co/payments/v2` |
| Endpoint producción | `https://api.nequi.com.co/payments/v2` |
| Método | `POST /-services-paymentservice-generatecodeqr` |
| Autenticación | Header `x-api-key: <api_key>` |
| Content-Type | `application/json` |

**Request:**
```json
{
  "RequestMessage": {
    "RequestHeader": {
      "Channel": "PC001",
      "RequestDate": "2026-02-26T10:00:00-05:00",
      "MessageID": "<uuid>",
      "ClientID": "<phone_code>",
      "Destination": {
        "Name": "PaymentService",
        "Namespace": "http://www.namespace.co/types",
        "Operation": "generateCodeQR"
      }
    },
    "RequestBody": {
      "any": {
        "generateCodeQRRQ": {
          "Code":       "<phone_code>",
          "Value":      "50000",
          "Reference1": "<reference>",
          "Reference2": "pay",
          "Reference3": "<transaction_id>"
        }
      }
    }
  }
}
```

**Response exitosa:**
```json
{
  "ResponseMessage": {
    "ResponseHeader": {
      "ResponseCode": { "EntityCode": "0", "SystemCode": "200" }
    },
    "ResponseBody": {
      "any": {
        "generateCodeQRRS": {
          "QRValue": "<qr_string>",
          "transactionID": "NQI-2026-XYZ"
        }
      }
    }
  }
}
```

### Mensaje de respuesta publicado a `pay.responses`

**Éxito:**
```json
{
  "payment_transaction_id": 42,
  "gateway_code": "nequi",
  "status": "success",
  "external_id": "NQI-2026-XYZ",
  "gateway_response": { "qr_value": "...", "transaction_id": "..." },
  "correlation_id": "...",
  "timestamp": "2026-02-26T10:01:00Z",
  "processing_time_ms": 850
}
```

**Error:**
```json
{
  "payment_transaction_id": 42,
  "gateway_code": "nequi",
  "status": "error",
  "error": "Nequi API returned status 400",
  "error_code": "NEQUI_API_ERROR",
  "correlation_id": "...",
  "timestamp": "2026-02-26T10:01:00Z",
  "processing_time_ms": 320
}
```

---

## Configurar una Nueva Pasarela

Para agregar una nueva pasarela (ej: Wompi):

1. **Crear carpeta** `integrations/pay/wompi/` con la misma estructura que `nequi/`
2. **Implementar** `INequiClient` equivalente (`IWompiClient`)
3. **Crear consumer** que escuche `pay.wompi.requests`
4. **Registrar en router** (`router/bundle.go`) el nuevo `case "wompi"`
5. **Registrar en bundle** (`integrations/pay/bundle.go`) el nuevo `wompi.New(...)`
6. **Agregar constante** `GatewayWompi = "wompi"` en `modules/pay/internal/domain/constants/`
7. **Crear tipo de integración** `wompi_pay` en la tabla `integration_types` con las credenciales

---

## Configuración de Integración en DB

Antes de usar una pasarela, debe existir un registro en `integration_types`:

```sql
-- Verificar si existe
SELECT id, code, name FROM integration_types WHERE code = 'nequi_pay';

-- Crear si no existe (vía admin UI o migración)
INSERT INTO integration_types (code, name, description, created_at, updated_at)
VALUES ('nequi_pay', 'Nequi Pay', 'Pasarela de pagos Nequi Colombia', NOW(), NOW());
```

Las credenciales se configuran desde el admin → Tipos de Integración → Nequi Pay → Editar (se almacenan cifradas automáticamente).

---

## Variables de Entorno

```env
# Clave de descifrado de credenciales (32 bytes, mismo que el resto del sistema)
ENCRYPTION_KEY=your-32-byte-encryption-key

# RabbitMQ
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=admin
RABBITMQ_PASS=admin
```
