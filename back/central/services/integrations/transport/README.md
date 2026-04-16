# Modulo Transport

Modulo de integracion con proveedores de transporte (carriers) para cotizacion, generacion de guias, rastreo y cancelacion de envios.

## Proposito

Este modulo centraliza la comunicacion con multiples proveedores de transporte mediante un patron de routing basado en mensajeria (RabbitMQ). El modulo `shipments` publica solicitudes a una cola unificada (`transport.requests`) y el router interno las distribuye al carrier correspondiente segun el `integration_type_id`.

Cada proveedor se implementa como un sub-modulo aislado con arquitectura hexagonal propia.

## Proveedores Soportados

| Proveedor  | type_id | Queue                         | Estado         |
|------------|---------|-------------------------------|----------------|
| EnvioClick | 12      | transport.envioclick.requests | Completo       |
| Enviame    | 13      | transport.enviame.requests    | Esqueleto (TODO) |
| Tu         | 14      | transport.tu.requests         | Solo queue declarada |
| MiPaquete  | 15      | transport.mipaquete.requests  | Esqueleto (TODO) |

## Operaciones

Todos los proveedores exponen (o expondran) las mismas cuatro operaciones:

- **quote** - Cotizar tarifas de envio dados origen, destino y paquetes.
- **generate** - Crear envio y generar guia de transporte (tracking number + etiqueta).
- **track** - Consultar el estado y eventos de rastreo de un envio.
- **cancel** - Cancelar un envio previamente generado.

## Entidades Principales

### Mensaje de Solicitud (`TransportRequestMessage`)

```
shipment_id        - ID del envio (opcional)
provider           - Nombre del carrier (envioclick, enviame, mipaquete)
operation          - Operacion solicitada (quote, generate, track, cancel)
correlation_id     - ID de correlacion para trazabilidad
business_id        - ID del negocio
integration_id     - ID de la integracion (para resolver credenciales)
base_url           - URL base del API (opcional, se resuelve automaticamente)
is_test            - Indica si es ambiente de pruebas
payload            - Datos especificos de la operacion (map generico)
```

### Mensaje de Respuesta (`TransportResponseMessage`)

```
shipment_id        - ID del envio
business_id        - ID del negocio
provider           - Nombre del carrier
operation          - Operacion ejecutada
status             - "success" o "error"
correlation_id     - ID de correlacion
is_test            - Indica si es ambiente de pruebas
data               - Datos de la respuesta (map generico)
error              - Mensaje de error (si aplica)
```

### Entidades de EnvioClick (proveedor completo)

- `QuoteRequest` / `QuoteResponse` - Cotizacion con paquetes, origen, destino y tarifas.
- `GenerateResponse` - Numero de tracking, URL de etiqueta, referencia.
- `TrackingResponse` - Estado actual, carrier, historial de eventos.
- `CancelResponse` - Estado y mensaje de la cancelacion.

## Colas RabbitMQ

```
transport.requests                  <- Cola de entrada unificada (publicada por shipments)
transport.envioclick.requests       <- Cola especifica para EnvioClick
transport.enviame.requests          <- Cola especifica para Enviame
transport.tu.requests               <- Cola especifica para Tu
transport.mipaquete.requests        <- Cola especifica para MiPaquete
transport.responses                 -> Cola de salida unificada (consumida por shipments)
```

### Flujo de mensajes

```
shipments --> [transport.requests] --> Router --> [transport.{carrier}.requests] --> Consumer del carrier
                                                                                         |
shipments <-- [transport.responses] <------- Response Publisher <-------------------------+
```

## Resolucion de Credenciales

Los consumers resuelven las credenciales API de cada carrier de forma dinamica:

- **EnvioClick**: Soporta credenciales por negocio (api_key desencriptada de la integracion) o token de plataforma compartido (`use_platform_token=true` en config, que lee `ENVIOCLICK_API_KEY` de `platform_credentials_encrypted` del integration type).
- **Enviame**: Desencripta `api_key` de la integracion. Fallback a variable de entorno `ENVIAME_API_KEY`.
- **MiPaquete**: Desencripta `api_key` de la integracion. Fallback a variable de entorno `MIPAQUETE_API_KEY`.

La interfaz `ICredentialResolver` (replicada localmente por aislamiento de modulos) es satisfecha por `core.IIntegrationCore`.

## Resolucion de Base URL (EnvioClick)

Prioridad de resolucion del URL base del API:

1. `base_url_test` del config de la integracion (si existe)
2. `base_url` del mensaje de solicitud
3. URL por defecto: `https://api.envioclickpro.com.co/api/v2`

## Endpoints HTTP

Este modulo **no expone endpoints HTTP**. Toda la comunicacion se realiza via colas RabbitMQ.

## Integracion con Otros Modulos

- **core** (`services/integrations/core`): Provee `IIntegrationCore` para resolver y desencriptar credenciales de integraciones. Tambien registra cada carrier como integracion conocida via `RegisterIntegration()`.
- **shipments** (`services/modules/shipments`): Modulo productor. Publica solicitudes de transporte a `transport.requests` y consume respuestas de `transport.responses`.
- **shared/rabbitmq**: Define las constantes de nombres de colas centralizadas.
- **shared/httpclient**: Cliente HTTP reutilizable con reintentos, timeouts y logging.
- **shared/log**: Logger estructurado (zerolog) con soporte de modulos.

## Arquitectura

```
transport/
|-- bundle.go                          # Inicializa todos los carriers y el router
|-- router/
|   +-- bundle.go                      # Router: consume transport.requests y reenvia al carrier
|-- envioclick/                        # Proveedor completo
|   |-- bundle.go
|   +-- internal/
|       |-- domain/
|       |   |-- entities.go            # Modelos del API (QuoteRequest, etc.)
|       |   +-- ports.go              # IEnvioClickClient
|       |-- app/
|       |   |-- constructor.go         # IUseCase
|       |   +-- operations.go          # Quote, Generate, Track, Cancel
|       +-- infra/
|           |-- primary/consumer/      # TransportRequestConsumer (RabbitMQ)
|           +-- secondary/
|               |-- client/            # HTTP client (quote, generate, track, cancel)
|               +-- queue/             # ResponsePublisher (transport.responses)
|-- enviame/                           # Proveedor esqueleto (misma estructura)
|   +-- internal/ ...
+-- mipaquete/                         # Proveedor esqueleto (misma estructura)
    +-- internal/ ...
```

Cada sub-modulo de carrier sigue la arquitectura hexagonal del proyecto:

- **domain**: Entidades, puertos (interfaces del cliente HTTP) y errores.
- **app**: Casos de uso que orquestan las operaciones.
- **infra/primary/consumer**: Consumer de RabbitMQ que recibe solicitudes ya ruteadas.
- **infra/secondary/client**: Cliente HTTP que se comunica con el API externo del carrier.
- **infra/secondary/queue**: Publisher que envia las respuestas a `transport.responses`.

## Estado de Implementacion

- **EnvioClick**: Completamente implementado. Soporta las 4 operaciones (quote, generate, track, cancel) con manejo de errores en espanol, parsing tolerante del tracker (string o numero), y mapeo de errores amigables del API.
- **Enviame**: Estructura completa pero todas las operaciones retornan "not yet implemented". Los modelos del dominio son placeholders genericos.
- **MiPaquete**: Mismo estado que Enviame. Estructura lista para implementar cuando se disponga del API.
- **Tu**: Solo existe la queue declarada en el router. No tiene sub-modulo implementado.
