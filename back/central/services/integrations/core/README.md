# Modulo: Integration Core

Hub central de gestion de integraciones del proyecto Probability. Administra el ciclo de vida completo de las integraciones con servicios externos (e-commerce, facturacion electronica, mensajeria, transporte, pagos) y actua como punto de registro y coordinacion para todos los proveedores.

---

## Proposito

Integration Core es el nucleo del sistema de integraciones. Sus responsabilidades principales son:

- **CRUD de integraciones**: Crear, leer, actualizar y eliminar integraciones configuradas por cada negocio.
- **CRUD de tipos de integracion**: Administrar el catalogo de tipos disponibles (Shopify, WhatsApp, Siigo, etc.) y sus categorias.
- **Encriptacion de credenciales**: Todas las credenciales se encriptan (AES) antes de persistirse y se desencriptan bajo demanda. Soporta credenciales de plataforma compartidas (`use_platform_token`).
- **Cache en Redis**: Metadata y credenciales desencriptadas se cachean (TTL 24h) con warm-up al iniciar el servidor para evitar cold start.
- **Registro de providers**: Patron registry thread-safe donde cada modulo de integracion registra su implementacion de `IIntegrationContract`.
- **Delegacion de operaciones**: Sync de ordenes, test de conexion y gestion de webhooks se delegan al provider registrado para cada tipo.
- **Observadores**: Permite a otros modulos suscribirse a eventos de creacion de integraciones mediante callbacks.

---

## Entidades principales

| Entidad | Descripcion |
|---------|-------------|
| `IntegrationCategory` | Categoria que agrupa tipos de integraciones (ecommerce, invoicing, messaging, transport, etc.). |
| `IntegrationType` | Tipo de integracion disponible en el catalogo (Shopify, Siigo, EnvioClick, etc.). Define schemas de config/credenciales, URLs base, icono e imagen. |
| `Integration` | Instancia concreta de una integracion configurada por un negocio. Contiene config (JSON), credenciales encriptadas, estado activo/inactivo, modo testing y referencia al tipo. |
| `PublicIntegration` | Vista publica de una integracion (sin credenciales, con config deserializada). Usada por modulos consumidores. |
| `IntegrationWithCredentials` | Integracion con credenciales desencriptadas. Solo uso interno, nunca se expone via HTTP. |
| `IIntegrationContract` | Interfaz que todo provider debe implementar: `TestConnection`, `SyncOrders`, webhooks, `UpdateInventory`. Los providers que no soporten una operacion embeben `BaseIntegration`, que retorna `ErrNotSupported` por defecto. |

---

## Tipos de integracion

Constantes definidas en `type_codes.go`, mapeadas al campo `integration_types.id` en la base de datos:

| ID | Codigo | Categoria |
|----|--------|-----------|
| 1 | shopify | E-commerce |
| 2 | whatsapp | Mensajeria |
| 3 | mercado_libre | E-commerce |
| 4 | woocommerce | E-commerce |
| 5 | softpymes | Facturacion |
| 6 | platform | Interno |
| 7 | factus | Facturacion |
| 8 | siigo | Facturacion |
| 9 | alegra | Facturacion |
| 10 | world_office | Facturacion |
| 11 | helisa | Facturacion |
| 12 | envioclick | Transporte |
| 13 | enviame | Transporte |
| 14 | tu | Transporte |
| 15 | mipaquete | Transporte |
| 16 | vtex | E-commerce |
| 17 | tiendanube | E-commerce |
| 18 | magento | E-commerce |
| 19 | amazon | Marketplace |
| 20 | falabella | Marketplace |
| 21 | exito | Marketplace |

La funcion canonica `IntegrationTypeCodeAsInt(code string) int` convierte codigos string a su ID numerico.

---

## Endpoints

Todos los endpoints requieren autenticacion JWT.

### Integraciones (`/integrations`)

| Metodo | Ruta | Descripcion |
|--------|------|-------------|
| GET | `/integrations` | Listar integraciones (paginado, con filtros por tipo, categoria, business, estado, busqueda) |
| GET | `/integrations/simple` | Listar integraciones en formato simplificado |
| GET | `/integrations/:id` | Obtener integracion por ID (credenciales solo para super admin) |
| GET | `/integrations/type/:type` | Obtener integracion por codigo de tipo |
| POST | `/integrations` | Crear integracion (valida conexion con el provider antes de guardar) |
| PUT | `/integrations/:id` | Actualizar integracion |
| DELETE | `/integrations/:id` | Eliminar integracion |
| POST | `/integrations/test` | Test de conexion sin guardar (con config y credenciales raw) |
| POST | `/integrations/:id/test` | Test de conexion de una integracion existente |
| POST | `/integrations/:id/sync` | Sincronizar ordenes de una integracion |
| POST | `/integrations/sync-orders/business/:business_id` | Sincronizar ordenes de todas las integraciones activas de un negocio |
| PUT | `/integrations/:id/activate` | Activar integracion |
| PUT | `/integrations/:id/deactivate` | Desactivar integracion |
| PUT | `/integrations/:id/set-default` | Marcar integracion como default de su tipo |
| GET | `/integrations/:id/webhook` | Obtener URL del webhook |
| GET | `/integrations/:id/webhooks` | Listar webhooks registrados en la plataforma externa |
| GET | `/integrations/:id/webhooks/verify` | Verificar webhooks existentes que coincidan con nuestra URL |
| POST | `/integrations/:id/webhooks/create` | Crear webhook en la plataforma externa |
| DELETE | `/integrations/:id/webhooks/:webhook_id` | Eliminar webhook |

### Tipos de integracion (`/integration-types`)

| Metodo | Ruta | Descripcion |
|--------|------|-------------|
| GET | `/integration-types` | Listar todos los tipos (filtrable por categoria) |
| GET | `/integration-types/active` | Listar solo tipos activos |
| GET | `/integration-types/:id` | Obtener tipo por ID |
| GET | `/integration-types/:id/platform-credentials` | Obtener credenciales de plataforma desencriptadas |
| GET | `/integration-types/code/:code` | Obtener tipo por codigo |
| POST | `/integration-types` | Crear tipo de integracion (soporta subir imagen a S3) |
| PUT | `/integration-types/:id` | Actualizar tipo de integracion |
| DELETE | `/integration-types/:id` | Eliminar tipo (falla si tiene integraciones asociadas) |

### Categorias (`/integration-categories`)

| Metodo | Ruta | Descripcion |
|--------|------|-------------|
| GET | `/integration-categories` | Listar todas las categorias de integraciones |

---

## Integracion con otros modulos

Integration Core se inicializa primero en `integrations/bundle.go` y se pasa como dependencia (`IIntegrationCore`) a todos los sub-modulos:

```
integrations/bundle.go
  |
  +-- core.New()            --> IIntegrationCore (hub central)
  |
  +-- messaging.New(core)   --> WhatsApp registra su provider (type 2)
  +-- ecommerce.New(core)   --> Shopify, WooCommerce, MeLi, VTEX, Tiendanube,
  |                              Magento, Amazon, Falabella, Exito registran providers
  +-- invoicing.New(core)   --> Softpymes, Factus, Siigo, Alegra, WorldOffice,
  |                              Helisa registran providers
  +-- transport.New(core)   --> EnvioClick, Enviame, Tu, MiPaquete registran providers
  +-- pay.New()             --> Independiente (Nequi, etc. - no registra providers en core)
```

### Interfaces expuestas

- **`IIntegrationService`**: Interfaz liviana para modulos consumidores. Permite consultar integraciones por ID o external ID, desencriptar credenciales individuales, actualizar config y obtener credenciales de plataforma. Usada por los consumers de facturacion, e-commerce, etc.
- **`IIntegrationCore`**: Interfaz completa que extiende `IIntegrationService` con registro de providers, observadores, sync de ordenes, test de conexion y operaciones de webhooks. Solo usada por `integrations/bundle.go` y los bundles de providers.

### Flujo tipico de un modulo consumidor

1. En su `bundle.go`, el modulo recibe `core.IIntegrationCore`.
2. Registra su provider: `core.RegisterIntegration(core.IntegrationTypeShopify, shopifyProvider)`.
3. Opcionalmente se suscribe a eventos: `core.OnIntegrationCreated(tipo, callback)`.
4. Para operaciones en runtime, sus consumers usan `IIntegrationService` para obtener credenciales desencriptadas y configuracion.

### Re-exports publicos

Como `internal/domain` no es importable desde fuera, `bundle.go` re-exporta los tipos y constantes necesarios como type aliases:

```go
type IIntegrationContract = domain.IIntegrationContract
type BaseIntegration       = domain.BaseIntegration
type WebhookInfo           = domain.WebhookInfo
type PublicIntegration     = domain.PublicIntegration
// + constantes IntegrationTypeShopify ... IntegrationTypeExito
```

---

## Arquitectura (capas hexagonales)

```
core/
+-- bundle.go                              # Constructor, facade, re-exports publicos
+-- internal/
    +-- domain/
    |   +-- entities.go                    # IntegrationCategory, IntegrationType, Integration, PublicIntegration
    |   +-- ports.go                       # IRepository, IEncryptionService, IIntegrationCache, IIntegrationUseCase
    |   +-- dtos.go                        # CreateIntegrationDTO, UpdateIntegrationDTO, IntegrationFilters
    |   +-- errors.go                      # Errores de dominio
    |   +-- provider_contract.go           # IIntegrationContract, BaseIntegration, ErrNotSupported
    |   +-- type_codes.go                  # Constantes de tipos (IDs) y funcion IntegrationTypeCodeAsInt
    |
    +-- app/
    |   +-- usecaseintegrations/           # Caso de uso principal de integraciones
    |   |   +-- constructor.go             # New(), RegisterProvider, RegisterObserver
    |   |   +-- create-integration.go      # Crear con validacion de conexion, cache y webhooks automaticos
    |   |   +-- update-integration.go      # Actualizar con re-encriptacion e invalidacion de cache
    |   |   +-- test-integration.go        # Test de conexion via provider registrado
    |   |   +-- sync_orders.go             # Sync delegado al provider (por integracion o por negocio)
    |   |   +-- webhook_ops.go             # CRUD de webhooks delegado al provider
    |   |   +-- warm_cache.go              # Pre-carga de integraciones activas en Redis al startup
    |   |   +-- consumer_methods.go        # Metodos de conveniencia: GetByExternalID, UpdateConfig, OnCreated
    |   |   +-- provider_registry.go       # Registry thread-safe map[int]IIntegrationContract
    |   |
    |   +-- usecaseintegrationtype/        # Caso de uso de tipos de integracion y categorias
    |       +-- constructor.go
    |       +-- create-integration-type.go # Soporta upload de imagen a S3, encriptacion de platform credentials
    |       +-- get-platform-credentials.go
    |
    +-- infra/
        +-- primary/handlers/
        |   +-- handlerintegrations/       # Handlers HTTP de integraciones
        |   |   +-- router.go              # Registro de rutas /integrations
        |   |   +-- request/               # Structs de request HTTP
        |   |   +-- response/              # Structs de response HTTP
        |   |   +-- mapper/                # Mapeo domain <-> HTTP
        |   |
        |   +-- handlerintegrationtype/    # Handlers HTTP de tipos de integracion
        |       +-- router.go              # Registro de rutas /integration-types y /integration-categories
        |       +-- request/, response/, mapper/
        |
        +-- secondary/
            +-- repository/                # Persistencia PostgreSQL via GORM
            |   +-- constructor.go         # New(db, logger, encryptionService, cache)
            |   +-- integration_repository.go          # CRUD integraciones + encriptacion automatica
            |   +-- integration_type_repository.go     # CRUD tipos de integracion
            |   +-- integration_category_repository.go # Consultas de categorias
            |
            +-- cache/                     # Cache Redis
            |   +-- integration_cache.go   # Implementacion de IIntegrationCache (TTL 24h)
            |   +-- keys.go               # Formato: integration:meta:{id}, integration:creds:{id}, etc.
            |
            +-- encryption/               # Servicio de encriptacion AES-256
                +-- service.go             # EncryptCredentials, DecryptCredentials, EncryptValue, DecryptValue
```

### Flujo de datos

1. **Request HTTP** llega al handler.
2. El handler parsea el request y llama al **use case**.
3. El use case consulta primero el **cache** (Redis). Si hay hit, retorna directamente.
4. Si hay miss, consulta el **repositorio** (PostgreSQL via GORM).
5. Las credenciales se desencriptan/encriptan a traves del **servicio de encriptacion**.
6. Para operaciones como sync o test, el use case consulta el **provider registry** y delega al provider correspondiente.
7. Al crear/actualizar integraciones, se invalida y re-cachea la metadata y credenciales en Redis.

---

Ultima actualizacion: 2026-03-01
