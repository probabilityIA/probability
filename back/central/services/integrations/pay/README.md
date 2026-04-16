# Modulo Pay -- Integraciones de Pasarelas de Pago

Capa de integracion que conecta el sistema con las APIs externas de pasarelas de pago. Actua como broker asincrono: recibe solicitudes genericas de la cola `pay.requests`, las enruta a la pasarela correspondiente segun el campo `gateway_code`, ejecuta la operacion de pago y devuelve el resultado estandarizado a `pay.responses`.

---

## Proposito

- Procesar solicitudes de pago de forma asincrona a traves de RabbitMQ.
- Abstraer la comunicacion con multiples pasarelas de pago bajo una interfaz unificada.
- Enrutar dinamicamente cada solicitud al gateway correcto mediante un router centralizado.
- Publicar respuestas estandarizadas (exito o error) en una cola comun de salida.

---

## Gateways Soportados

| Gateway | Codigo | Cola | Estado | Operacion |
|---------|--------|------|--------|-----------|
| Nequi | `nequi` | `pay.nequi.requests` | Implementado | GenerateQR |
| Bold | `bold` | `pay.bold.requests` | Implementado | CreatePaymentLink |
| Wompi | `wompi` | `pay.wompi.requests` | Esqueleto (TODO) | CreateTransaction |
| Stripe | `stripe` | `pay.stripe.requests` | Esqueleto (TODO) | CreatePaymentIntent |
| PayU | `payu` | `pay.payu.requests` | Esqueleto (TODO) | CreateTransaction |
| ePayco | `epayco` | `pay.epayco.requests` | Esqueleto (TODO) | CreateCheckout |
| MercadoPago | `melipago` | `pay.melipago.requests` | Esqueleto (TODO) | CreatePreference |

---

## Arquitectura General

```
"pay.requests" (generico)
        |
        v
  +-----------+
  |  router/  |  <-- Lee gateway_code, reenvio al gateway correcto
  +-----------+
        |
        +---> "pay.nequi.requests"    ---> nequi/
        +---> "pay.bold.requests"     ---> bold/
        +---> "pay.wompi.requests"    ---> wompi/     (pendiente)
        +---> "pay.stripe.requests"   ---> stripe/    (pendiente)
        +---> "pay.payu.requests"     ---> payu/      (pendiente)
        +---> "pay.epayco.requests"   ---> epayco/    (pendiente)
        +---> "pay.melipago.requests" ---> melipago/  (pendiente)
                      |
                      v
              API externa del gateway
                      |
                      v
              "pay.responses" (resultado estandarizado)
                      |
                      v
              modules/pay ResponseConsumer
```

---

## Estructura de Carpetas

```
integrations/pay/
  bundle.go                 # Ensambla todos los gateways + router
  router/
    bundle.go               # Consume pay.requests, enruta por gateway_code
  bold/                     # Gateway Bold (implementado)
    bundle.go
    internal/
      domain/
        entities/bold_payment.go       # BoldPaymentResult
        errors/errors.go               # Errores de dominio
        ports/ports.go                 # IBoldClient, IIntegrationRepository, IResponsePublisher
      app/
        constructor.go                 # New() -> IUseCase
        process_payment.go             # Logica de procesamiento
      infra/
        primary/consumer/
          bold_consumer.go             # Consume pay.bold.requests
        secondary/
          client/bold_client.go        # HTTP -> API Bold
          queue/response_publisher.go  # Publica a pay.responses
          repository/integration_repository.go  # Credenciales desde DB
  nequi/                    # Gateway Nequi (implementado) -- misma estructura
  wompi/                    # Gateway Wompi (esqueleto)    -- misma estructura
  stripe/                   # Gateway Stripe (esqueleto)   -- misma estructura
  payu/                     # Gateway PayU (esqueleto)     -- misma estructura
  epayco/                   # Gateway ePayco (esqueleto)   -- misma estructura
  melipago/                 # Gateway MercadoPago (esqueleto) -- misma estructura
```

Todos los gateways siguen exactamente la misma estructura hexagonal interna.

---

## Colas de RabbitMQ

| Cola | Productor | Consumidor | Descripcion |
|------|-----------|------------|-------------|
| `pay.requests` | `modules/pay` | `router/` | Solicitud de pago generica |
| `pay.nequi.requests` | `router/` | `nequi/` | Solicitud para Nequi |
| `pay.bold.requests` | `router/` | `bold/` | Solicitud para Bold |
| `pay.wompi.requests` | `router/` | `wompi/` | Solicitud para Wompi |
| `pay.stripe.requests` | `router/` | `stripe/` | Solicitud para Stripe |
| `pay.payu.requests` | `router/` | `payu/` | Solicitud para PayU |
| `pay.epayco.requests` | `router/` | `epayco/` | Solicitud para ePayco |
| `pay.melipago.requests` | `router/` | `melipago/` | Solicitud para MercadoPago |
| `pay.responses` | Todos los gateways | `modules/pay` | Resultado estandarizado |

Las constantes de colas estan centralizadas en `shared/rabbitmq/queues.go`.

---

## Mensajes

### PaymentRequestMsg (entrada a cada gateway)

```json
{
  "payment_transaction_id": 42,
  "business_id": 1,
  "gateway_code": "bold",
  "amount": 50000,
  "currency": "COP",
  "reference": "ORD-2026-001",
  "payment_method": "CREDIT_CARD",
  "description": "Pago de pedido",
  "metadata": {},
  "correlation_id": "uuid-xxx",
  "timestamp": "2026-03-01T10:00:00Z"
}
```

### PaymentResponseMsg (salida a pay.responses)

Exito:
```json
{
  "payment_transaction_id": 42,
  "gateway_code": "bold",
  "status": "success",
  "external_id": "LINK-123",
  "gateway_response": {
    "payment_link_id": "LINK-123",
    "checkout_url": "https://checkout.bold.co/...",
    "status": "ACTIVE"
  },
  "correlation_id": "uuid-xxx",
  "processing_time_ms": 850,
  "timestamp": "2026-03-01T10:00:01Z"
}
```

Error:
```json
{
  "payment_transaction_id": 42,
  "gateway_code": "bold",
  "status": "error",
  "error": "bold api error: status=401",
  "error_code": "api_error",
  "correlation_id": "uuid-xxx",
  "processing_time_ms": 320,
  "timestamp": "2026-03-01T10:00:01Z"
}
```

---

## Gateways Implementados

### Bold

- **API Base**: `https://integrations.api.bold.co`
- **Endpoint**: `POST /online/link/v1` (crear link de pago)
- **Autenticacion**: Header `Authorization: x-api-key <api_key>`
- **Moneda por defecto**: COP
- **Metodos de pago**: CREDIT_CARD, PSE, NEQUI, BOTON_BANCOLOMBIA
- **Credenciales en DB**: `integration_types` donde `code = 'bold_pay'`
- **Campos de credenciales**: `api_key`, `environment`

### Nequi

- **API Base sandbox**: `https://api.sandbox.nequi.com/payments/v2`
- **API Base produccion**: `https://api.nequi.com/payments/v2`
- **Endpoint**: `POST /-services-paymentservice-generatecodeqr` (generar QR de pago)
- **Autenticacion**: Header `x-api-key: <api_key>`
- **Credenciales en DB**: `integration_types` donde `code = 'nequi_pay'`
- **Campos de credenciales**: `api_key`, `environment`, `phone_code`

---

## Gateways Pendientes (Esqueleto)

Los siguientes gateways tienen la estructura hexagonal completa pero sus clientes HTTP retornan `"not yet implemented"`. Solo falta implementar la llamada real a la API:

| Gateway | Documentacion de referencia |
|---------|-----------------------------|
| Wompi | https://docs.wompi.co/ |
| Stripe | https://stripe.com/docs/api |
| PayU | https://developers.payulatam.com/ |
| ePayco | https://docs.epayco.co/ |
| MercadoPago | https://www.mercadopago.com.co/developers/es/reference |

---

## Credenciales y Seguridad

Las credenciales de cada gateway se almacenan cifradas en la tabla `integration_types`, campo `platform_credentials_encrypted`, usando AES-256-GCM. La clave de descifrado proviene de la variable de entorno `ENCRYPTION_KEY`.

Configuraciones por gateway:

| Gateway | Campos de credenciales |
|---------|------------------------|
| Bold | `api_key`, `environment` |
| Nequi | `api_key`, `environment`, `phone_code` |
| Wompi | `private_key`, `environment` |
| Stripe | `secret_key`, `environment` |
| PayU | `api_key`, `api_login`, `account_id`, `merchant_id`, `environment` |
| ePayco | `customer_id`, `key`, `environment` |
| MercadoPago | `access_token`, `environment` |

---

## Endpoints HTTP

Este modulo **no expone endpoints HTTP**. Toda la comunicacion es asincrona mediante colas de RabbitMQ.

---

## Integracion con Otros Modulos

- **modules/pay**: Productor de `pay.requests` y consumidor de `pay.responses`. Es el modulo de negocio que orquesta los pagos.
- **shared/rabbitmq**: Constantes de colas centralizadas en `shared/rabbitmq/queues.go`.
- **shared/db**: Acceso a la tabla `integration_types` para obtener credenciales encriptadas.
- **shared/env**: Lectura de `ENCRYPTION_KEY` y otras variables de configuracion.
- **core / invoicing / ecommerce / messaging / transport**: No hay dependencia directa. La comunicacion con otros modulos de integracion se realiza exclusivamente mediante colas RabbitMQ.

---

## Orden de Inicializacion

1. Se inicializan los 7 gateways (cada uno declara su cola y levanta su consumer en goroutine).
2. Se inicializa el router al final para que las colas de los gateways ya esten declaradas.
3. El router comienza a escuchar en `pay.requests` y distribuye mensajes.

---

## Variables de Entorno

| Variable | Descripcion |
|----------|-------------|
| `ENCRYPTION_KEY` | Clave AES-256-GCM para descifrar credenciales. Obligatoria en produccion. |
| `RABBITMQ_HOST` | Host de RabbitMQ |
| `RABBITMQ_PORT` | Puerto de RabbitMQ |
| `RABBITMQ_USER` | Usuario de RabbitMQ |
| `RABBITMQ_PASS` | Contrasena de RabbitMQ |

---

## Agregar un Nuevo Gateway

1. Crear carpeta `integrations/pay/<gateway>/` replicando la estructura de `bold/` o `nequi/`.
2. Implementar el cliente HTTP con la API real del gateway.
3. Agregar un `case "<gateway>"` en `router/bundle.go` dentro de `getGatewayQueue()`.
4. Agregar la constante de cola en `shared/rabbitmq/queues.go`.
5. Registrar `<gateway>.New(...)` en `integrations/pay/bundle.go`.
6. Crear el registro correspondiente en `integration_types` con las credenciales cifradas.
