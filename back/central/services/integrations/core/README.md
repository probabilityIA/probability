# Módulo Core de Integraciones

Sistema centralizado para gestionar todas las integraciones externas de Probability (e-commerce, facturación, mensajería, etc.).

## 📋 Índice

- [Descripción](#descripción)
- [Conceptos Clave](#conceptos-clave)
- [Arquitectura](#arquitectura)
- [Provider Registry](#provider-registry)
- [Tipos de Integración](#tipos-de-integración)
- [Categorías de Integraciones](#categorías-de-integraciones)
- [Flujo Completo](#flujo-completo)
- [API Endpoints](#api-endpoints)
- [Agregar un Nuevo Provider](#agregar-un-nuevo-provider)
- [Base de Datos](#base-de-datos)

---

## Descripción

Este módulo proporciona la infraestructura común para **todas las integraciones** de Probability:

- ✅ Catálogo unificado de tipos de integraciones
- ✅ Gestión de credenciales encriptadas
- ✅ Configuración por negocio (multi-tenant)
- ✅ Test de conexión para validar credenciales
- ✅ Webhooks y sincronización de órdenes
- ✅ Provider registry unificado (sin dependencias circulares)
- ✅ Sistema de categorías extensible

---

## Conceptos Clave

### 🎯 Integration Type (Tipo de Integración)

**Definición**: Representa un **tipo** de integración disponible en el "marketplace" de Probability.

**Analogía**: Es como una **app en el App Store** (disponible para instalar).

```sql
-- Catálogo de integraciones disponibles
integration_types:
id | code         | name         | category
---|--------------|--------------|----------
1  | shopify      | Shopify      | ecommerce
2  | whatsapp     | WhatsApp     | messaging
3  | mercadolibre | MercadoLibre | ecommerce
4  | woocommerce  | WooCommerce  | ecommerce
5  | softpymes    | Softpymes    | invoicing
7  | factus       | Factus       | invoicing
8  | siigo        | Siigo        | invoicing
```

---

### 🔌 Integration (Integración Configurada)

**Definición**: Representa una **instancia configurada** de un tipo de integración para un negocio específico.

**Analogía**: Es como una **app instalada** en tu teléfono con TUS configuraciones.

```sql
-- Instancias configuradas por negocio
integrations:
id | business_id | integration_type_id | name                      | credentials
---|-------------|---------------------|---------------------------|------------------
1  | 1           | 1 (shopify)         | Mi Tiendita - Shopify     | {api_key: "..."}
2  | 1           | 5 (softpymes)       | Mi Tiendita - Softpymes   | {api_key: "...", nit: "900..."}
3  | 2           | 1 (shopify)         | Tu Negocio - Shopify      | {api_key: "..."}
4  | 2           | 7 (factus)          | Tu Negocio - Factus       | {token: "..."}
```

---

### 📊 Category (Categoría)

**Definición**: Agrupa tipos de integraciones por su propósito.

**Categorías Actuales**:
- `ecommerce` - Plataformas de venta (Shopify, MeLi, WooCommerce)
- `invoicing` - Proveedores de facturación electrónica (Softpymes, Factus, Siigo)
- `messaging` - Canales de mensajería (WhatsApp)
- `payment` - Procesadores de pago *[futuro]*
- `shipping` - Operadores logísticos *[futuro]*

---

## Arquitectura

### Estructura de paquetes

```
core/
├── bundle.go                          ← Fachada pública (re-exports + IIntegrationCore)
└── internal/                          ← Todo privado — no importar desde fuera
    ├── domain/
    │   ├── entities.go                ← Integration, IntegrationType, PublicIntegration
    │   ├── dtos.go                    ← CreateIntegrationDTO, UpdateIntegrationDTO, etc.
    │   ├── ports.go                   ← IIntegrationUseCase, IRepository, etc.
    │   ├── provider_contract.go       ← IIntegrationContract, BaseIntegration, ErrNotSupported
    │   ├── type_codes.go              ← Constantes IntegrationTypeShopify=1…Siigo=8
    │   └── errors.go                  ← Errores de dominio
    ├── app/
    │   ├── usecaseintegrations/
    │   │   ├── constructor.go         ← New(repo, enc, cache, log, config)
    │   │   ├── provider_registry.go   ← Registro unificado int→IIntegrationContract
    │   │   ├── sync_orders.go         ← SyncOrdersByIntegrationID, SyncOrdersByBusiness
    │   │   ├── webhook_ops.go         ← GetWebhookURL, CreateWebhookForIntegration, etc.
    │   │   ├── consumer_methods.go    ← GetIntegrationByExternalID, OnIntegrationCreated, etc.
    │   │   └── *.go                   ← CRUD: create, update, delete, list, activate…
    │   └── usecaseintegrationtype/
    └── infra/
        ├── primary/handlers/
        └── secondary/
            ├── repository/
            ├── cache/
            └── encryption/
```

### Flujo de dependencias

```
                ┌──────────────────────────────────────────────┐
                │              bundle.go (IIntegrationCore)     │
                │   thin facade — delega todo al use case       │
                └────────────────────┬─────────────────────────┘
                                     │
                                     ▼
                ┌──────────────────────────────────────────────┐
                │     IntegrationUseCase (IIntegrationUseCase)  │
                │                                              │
                │  ┌──────────────────────────────────────┐   │
                │  │  providerReg map[int]IIntegrationContract │ │
                │  │  (Shopify, WhatsApp, Factus, Siigo…)  │   │
                │  └──────────────────────────────────────┘   │
                └──────────┬──────────────┬───────────────────┘
                           │              │
                           ▼              ▼
                      IRepository   IEncryptionService
```

**Antes del refactor** había una dependencia circular:
```
integrationCore → useCase → webhookCreator (= integrationCore)   ← CIRCULAR ❌
```

**Ahora** el use case contiene el registro de providers y no depende de `bundle.go`:
```
bundle.go → useCase (con providerReg) → providers   ← SIN CIRCULAR ✓
```

### Tipos re-exportados en `bundle.go`

El `domain` vive en `internal/` y no es importable directamente desde módulos externos. `bundle.go` re-exporta los tipos necesarios como type aliases:

```go
// En core/bundle.go
type IIntegrationContract  = domain.IIntegrationContract
type BaseIntegration        = domain.BaseIntegration
type WebhookInfo            = domain.WebhookInfo
type PublicIntegration      = domain.PublicIntegration
type IntegrationWithCredentials = domain.IntegrationWithCredentials

var ErrNotSupported = domain.ErrNotSupported

const (
    IntegrationTypeShopify      = domain.IntegrationTypeShopify      // 1
    IntegrationTypeWhatsApp     = domain.IntegrationTypeWhatsApp     // 2
    IntegrationTypeMercadoLibre = domain.IntegrationTypeMercadoLibre // 3
    IntegrationTypeWoocommerce  = domain.IntegrationTypeWoocommerce  // 4
    IntegrationTypeInvoicing    = domain.IntegrationTypeInvoicing    // 5 (Softpymes)
    IntegrationTypePlatform     = domain.IntegrationTypePlatform     // 6
    IntegrationTypeFactus       = domain.IntegrationTypeFactus       // 7
    IntegrationTypeSiigo        = domain.IntegrationTypeSiigo        // 8
)
```

Los módulos externos (shopify, factus, etc.) importan **solo el paquete `core`**, nunca `core/internal/domain`.

---

## Provider Registry

El use case mantiene un registro unificado `map[int]IIntegrationContract` donde cada provider registra sus capacidades.

### IIntegrationContract

```go
type IIntegrationContract interface {
    // Obligatorio — toda integración debe implementarlo
    TestConnection(ctx, config, credentials) error

    // Sincronización de órdenes (ej: Shopify)
    SyncOrdersByIntegrationID(ctx, integrationID string) error
    SyncOrdersByIntegrationIDWithParams(ctx, integrationID string, params interface{}) error

    // Webhooks — URL informativa
    GetWebhookURL(ctx, baseURL string, integrationID uint) (*WebhookInfo, error)

    // Webhooks — operaciones CRUD en plataformas externas
    ListWebhooks(ctx, integrationID string) ([]interface{}, error)
    DeleteWebhook(ctx, integrationID, webhookID string) error
    VerifyWebhooksByURL(ctx, integrationID, baseURL string) ([]interface{}, error)
    CreateWebhook(ctx, integrationID, baseURL string) (interface{}, error)
}
```

### BaseIntegration

Struct que implementa todos los métodos de `IIntegrationContract` retornando `ErrNotSupported`. Los providers embeben este struct y solo sobrescriben lo que soportan:

```go
type MyProvider struct {
    core.BaseIntegration   // todos los métodos no implementados → ErrNotSupported
}

// Solo sobrescribir lo que el provider soporta
func (p *MyProvider) TestConnection(ctx context.Context, config, creds map[string]interface{}) error {
    // validar credenciales contra la API externa
}

func (p *MyProvider) SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error {
    // sincronizar órdenes
}
```

### Registrar un provider

Los sub-módulos registran su implementación en `bundle.go` al inicializarse:

```go
// En shopify/bundle.go
coreIntegration.RegisterIntegration(core.IntegrationTypeShopify, shopifyCore)

// En factus/bundle.go
coreIntegration.RegisterIntegration(core.IntegrationTypeFactus, factusProvider)
```

### Observer de creación

Los sub-módulos pueden reaccionar cuando se crea una integración de su tipo:

```go
// En shopify/bundle.go
coreIntegration.OnIntegrationCreated(core.IntegrationTypeShopify,
    func(ctx context.Context, integration *core.PublicIntegration) {
        // ej: crear webhooks automáticamente
        useCase.CreateWebhook(ctx, fmt.Sprintf("%d", integration.ID), baseURL)
    },
)
```

---

## Tipos de Integración

Constantes canónicas definidas en `internal/domain/type_codes.go` y re-exportadas desde `bundle.go`:

| Constante | ID | Código BD | Provider | Estado |
|-----------|-----|-----------|----------|--------|
| `IntegrationTypeShopify` | 1 | `shopify` | shopify/ | ✅ Activo |
| `IntegrationTypeWhatsApp` | 2 | `whatsapp` | messaging/whatsapp/ | ✅ Activo |
| `IntegrationTypeMercadoLibre` | 3 | `mercadolibre` | — | 🔜 Próximamente |
| `IntegrationTypeWoocommerce` | 4 | `woocommerce` | — | 🔜 Próximamente |
| `IntegrationTypeInvoicing` | 5 | `softpymes` | invoicing/softpymes/ | ✅ Activo |
| `IntegrationTypePlatform` | 6 | `platform` | interno | — |
| `IntegrationTypeFactus` | 7 | `factus` | invoicing/factus/ | ✅ Activo |
| `IntegrationTypeSiigo` | 8 | `siigo` | invoicing/siigo/ | ✅ Activo |

Función de conversión código→ID (canónica):
```go
domain.IntegrationTypeCodeAsInt("shopify") // → 1
domain.IntegrationTypeCodeAsInt("factus")  // → 7
```

---

## Categorías de Integraciones

### 📦 E-commerce (Inbound)

**Propósito**: Recibir órdenes de plataformas de venta.

| Código | Nombre | Estado |
|--------|--------|--------|
| `shopify` | Shopify | ✅ Activo — webhooks + sync |
| `mercadolibre` | MercadoLibre | 🔜 Próximamente |
| `woocommerce` | WooCommerce | 🔜 Próximamente |

**Flujo**:
```
Cliente compra → Webhook → Probability → Crea orden → Notifica al negocio
```

---

### 💳 Facturación Electrónica (Outbound)

**Propósito**: Emitir facturas electrónicas ante autoridades fiscales (DIAN, etc.).

| Código | Nombre | País | Estado |
|--------|--------|------|--------|
| `softpymes` | Softpymes | 🇨🇴 Colombia | ✅ Activo |
| `factus` | Factus | 🇨🇴 Colombia | ✅ Activo |
| `siigo` | Siigo | 🇨🇴 Colombia | ✅ Activo |
| `alegra` | Alegra | 🇨🇴🇲🇽🇵🇪 Multi | 🔜 Próximamente |

**Flujo**:
```
Orden creada → Probability → Genera factura → API proveedor → DIAN → Cliente recibe factura
```

Se vincula una integración de e-commerce con una de facturación vía `invoicing_configs`.

---

### 📧 Mensajería (Bidirectional)

| Código | Nombre | Estado |
|--------|--------|--------|
| `whatsapp` | WhatsApp Business | ✅ Activo — Meta Business API |
| `telegram` | Telegram Bot | 🔜 Próximamente |

---

## Flujo Completo

### Ejemplo: Facturar Órdenes de Shopify con Factus

**Paso 1: Conectar Shopify**
```bash
POST /api/integrations
{
  "business_id": 1,
  "integration_type_id": 1,
  "name": "Mi Tiendita - Shopify",
  "credentials": {
    "api_key": "shpat_...",
    "api_secret": "shpss_...",
    "shop_domain": "mitiendita.myshopify.com"
  }
}
# Resultado: se crea webhook automáticamente en Shopify
```

**Paso 2: Conectar Factus**
```bash
POST /api/integrations
{
  "business_id": 1,
  "integration_type_id": 7,
  "name": "Mi Tiendita - Factus",
  "credentials": {
    "client_id": "...",
    "client_secret": "...",
    "username": "...",
    "password": "..."
  }
}
```

**Paso 3: Vincular vía invoicing_configs**
```bash
POST /api/invoicing/configs
{
  "source_integration_id": 1,      # Shopify
  "invoicing_integration_id": 2,   # Factus
  "enabled": true,
  "auto_invoice": true,
  "filters": { "min_amount": 50000, "only_paid": true }
}
```

---

## API Endpoints

### Integration Types

```http
GET    /api/integrations/types              # Listar tipos disponibles
GET    /api/integrations/types/:id          # Obtener tipo
POST   /api/integrations/types              # Crear tipo (admin)
PUT    /api/integrations/types/:id          # Actualizar tipo (admin)
DELETE /api/integrations/types/:id          # Eliminar tipo (admin)
```

### Integrations

```http
GET    /api/integrations                    # Listar integraciones del negocio
GET    /api/integrations/:id                # Obtener integración
POST   /api/integrations                    # Crear integración
PUT    /api/integrations/:id                # Actualizar integración
DELETE /api/integrations/:id                # Eliminar integración
POST   /api/integrations/:id/test           # Probar conexión
POST   /api/integrations/:id/activate       # Activar
POST   /api/integrations/:id/deactivate     # Desactivar
POST   /api/integrations/:id/sync           # Sincronizar órdenes (últimos 30 días)
```

### Webhooks

```http
GET    /api/integrations/:id/webhook                    # Obtener URL del webhook
GET    /api/integrations/:id/webhooks                   # Listar webhooks registrados
POST   /api/integrations/:id/webhooks/create            # Crear webhooks en plataforma externa
DELETE /api/integrations/:id/webhooks/:webhook_id       # Eliminar webhook
POST   /api/integrations/:id/webhooks/verify            # Verificar webhooks por URL
```

---

## Agregar un Nuevo Provider

### 1. Crear el paquete del provider

```
services/integrations/
└── mi-categoria/
    └── mi-provider/
        ├── bundle.go
        └── internal/
            ├── app/
            └── infra/
```

### 2. Implementar `IIntegrationContract`

```go
package miprovider

import (
    "context"
    "github.com/secamc93/probability/back/central/services/integrations/core"
)

type MiProvider struct {
    core.BaseIntegration   // default: ErrNotSupported para métodos no implementados
    client *MyAPIClient
}

func New(client *MyAPIClient) *MiProvider {
    return &MiProvider{client: client}
}

// Obligatorio
func (p *MiProvider) TestConnection(ctx context.Context, config, creds map[string]interface{}) error {
    // validar credenciales contra la API
    return p.client.Ping(creds["api_key"].(string))
}

// Opcional — solo si el provider soporta sync de órdenes
func (p *MiProvider) SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error {
    // ...
}
```

### 3. Registrar en `bundle.go`

```go
package miprovider

import (
    "github.com/secamc93/probability/back/central/services/integrations/core"
)

func New(router, logger, config, coreIntegration core.IIntegrationCore, ...) {
    provider := newMiProvider(...)

    // Registrar en el use case del core
    coreIntegration.RegisterIntegration(core.IntegrationTypeMiProvider, provider)

    // Opcional: reaccionar cuando se crea una integración de este tipo
    coreIntegration.OnIntegrationCreated(core.IntegrationTypeMiProvider,
        func(ctx context.Context, integration *core.PublicIntegration) {
            // setup automático
        },
    )
}
```

### 4. Agregar constante de tipo

En `core/internal/domain/type_codes.go`, agregar:
```go
const (
    // ... existentes ...
    IntegrationTypeMiProvider = 9 // Mi Provider
)
```

Y el case en `IntegrationTypeCodeAsInt`:
```go
case "mi_provider":
    return IntegrationTypeMiProvider
```

### 5. Registrar en `integrations/bundle.go`

```go
// En services/integrations/bundle.go
miprovider.New(subRouter, logger, config, coreIntegration, ...)
```

---

## Base de Datos

### Tabla: `integration_types`

```sql
CREATE TABLE integration_types (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,           -- 'shopify', 'factus', etc.
    name VARCHAR(100) NOT NULL,
    category_id INTEGER REFERENCES integration_categories(id),
    is_active BOOLEAN DEFAULT true,
    config_schema JSONB,                        -- Esquema JSON de configuración
    credentials_schema JSONB,                   -- Esquema JSON de credenciales (campos requeridos)
    image_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Tabla: `integrations`

```sql
CREATE TABLE integrations (
    id SERIAL PRIMARY KEY,
    business_id INTEGER NOT NULL,
    integration_type_id INTEGER NOT NULL REFERENCES integration_types(id),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50),
    store_id VARCHAR(100),                      -- ID externo (ej: shop domain de Shopify)
    is_active BOOLEAN DEFAULT true,
    is_default BOOLEAN DEFAULT false,
    config JSONB,                               -- Configuración específica del negocio
    credentials JSONB,                          -- Credenciales encriptadas (AES-256)
    description TEXT,
    last_sync_at TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,

    UNIQUE(business_id, integration_type_id, code)
);

CREATE INDEX idx_integrations_business ON integrations(business_id);
CREATE INDEX idx_integrations_type ON integrations(integration_type_id);
CREATE INDEX idx_integrations_active ON integrations(business_id, is_active);
```

---

## Seguridad

### Encriptación de Credenciales

Todas las credenciales se guardan **encriptadas** usando AES-256.

```go
// Al crear/actualizar — automático en el use case
encrypted, _ := encryption.Encrypt(rawCredentials)
integration.Credentials = encrypted

// Al usar integración — el use case desencripta al retornar IntegrationWithCredentials
rawCredentials := integration.DecryptedCredentials
```

**Variable de entorno requerida**: `ENCRYPTION_KEY`

---

## Relación con Otros Módulos

### Módulo `invoicing`

Consume `IIntegrationService` para obtener configuración del provider de facturación y desencriptar credenciales:

```go
type IIntegrationService interface {
    GetIntegrationByID(ctx, integrationID string) (*core.PublicIntegration, error)
    DecryptCredential(ctx, integrationID, fieldName string) (string, error)
    UpdateIntegrationConfig(ctx, integrationID string, newConfig map[string]interface{}) error
}
```

### Módulo `orders`

Las órdenes registran de qué integración provienen (`integration_id`).

---

## Roadmap

| Categoría | Prioridad | Estado |
|-----------|-----------|--------|
| `ecommerce/shopify` | Alta | ✅ Completo |
| `invoicing/softpymes` | Alta | ✅ Completo |
| `invoicing/factus` | Alta | ✅ Completo |
| `invoicing/siigo` | Alta | ✅ Completo |
| `messaging/whatsapp` | Media | ✅ Activo |
| `ecommerce/mercadolibre` | Alta | 🔜 Q2 2026 |
| `ecommerce/woocommerce` | Media | 🔜 Q2 2026 |
| `invoicing/alegra` | Media | 🔜 Q2 2026 |
| `payment/*` | Media | 🔜 Q3 2026 |
| `shipping/*` | Baja | 🔜 Q3 2026 |

---

**Última actualización**: 2026-02-22
