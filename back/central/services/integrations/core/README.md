# Módulo Core de Integraciones

Sistema centralizado para gestionar todas las integraciones externas de Probability (e-commerce, facturación, mensajería, etc.).

## 📋 Índice

- [Descripción](#descripción)
- [Conceptos Clave](#conceptos-clave)
- [Arquitectura](#arquitectura)
- [Categorías de Integraciones](#categorías-de-integraciones)
- [Flujo Completo](#flujo-completo)
- [API Endpoints](#api-endpoints)
- [Ejemplos de Uso](#ejemplos-de-uso)
- [Base de Datos](#base-de-datos)

---

## Descripción

Este módulo proporciona la infraestructura común para **todas las integraciones** de Probability:

- ✅ Catálogo unificado de tipos de integraciones
- ✅ Gestión de credenciales encriptadas
- ✅ Configuración por negocio (multi-tenant)
- ✅ Test de conexión para validar credenciales
- ✅ Webhooks y sincronización
- ✅ Sistema de categorías extensible

---

## Conceptos Clave

### 🎯 Integration Type (Tipo de Integración)

**Definición**: Representa un **tipo** de integración disponible en el "marketplace" de Probability.

**Ejemplos**:
- Shopify (e-commerce)
- Softpymes (facturación electrónica)
- WhatsApp (mensajería)
- MercadoLibre (marketplace)

**Analogía**: Es como una **app en el App Store** (disponible para instalar).

```sql
-- Catálogo de integraciones disponibles
integration_types:
id | code           | name           | category    | direction
---|----------------|----------------|-------------|-------------
1  | shopify        | Shopify        | ecommerce   | inbound
2  | mercadolibre   | MercadoLibre   | ecommerce   | inbound
3  | whatsapp       | WhatsApp       | messaging   | bidirectional
4  | softpymes      | Softpymes      | invoicing   | outbound
5  | alegra         | Alegra         | invoicing   | outbound
```

---

### 🔌 Integration (Integración Configurada)

**Definición**: Representa una **instancia configurada** de un tipo de integración para un negocio específico.

**Ejemplos**:
- "Mi Tiendita - Shopify" (business_id=1, type=shopify)
- "Mi Tiendita - Softpymes" (business_id=1, type=softpymes)
- "Tu Negocio - Alegra" (business_id=2, type=alegra)

**Analogía**: Es como una **app instalada** en tu teléfono con TUS configuraciones.

```sql
-- Instancias configuradas por negocio
integrations:
id | business_id | integration_type_id | name                      | credentials
---|-------------|---------------------|---------------------------|------------------
1  | 1           | 1 (shopify)         | Mi Tiendita - Shopify     | {api_key: "..."}
2  | 1           | 4 (softpymes)       | Mi Tiendita - Softpymes   | {api_key: "...", nit: "900..."}
3  | 2           | 1 (shopify)         | Tu Negocio - Shopify      | {api_key: "..."}
4  | 2           | 5 (alegra)          | Tu Negocio - Alegra       | {token: "..."}
```

---

### 📊 Category (Categoría)

**Definición**: Agrupa tipos de integraciones por su propósito.

**Categorías Actuales**:
- `ecommerce` - Plataformas de venta (Shopify, MeLi, Amazon)
- `invoicing` - Proveedores de facturación electrónica (Softpymes, Alegra, Siigo)
- `messaging` - Canales de mensajería (WhatsApp, Telegram)
- `payment` - Procesadores de pago (Stripe, PayPal) *[futuro]*
- `shipping` - Operadores logísticos (FedEx, DHL) *[futuro]*
- `accounting` - Software contable (QuickBooks, Xero) *[futuro]*

---

### 🔄 Direction (Dirección del Flujo)

**Definición**: Define la dirección del flujo de datos.

- `inbound` - Reciben datos en Probability (webhooks de Shopify, órdenes de MeLi)
- `outbound` - Envían datos desde Probability (facturas a Softpymes, notificaciones a WhatsApp)
- `bidirectional` - Ambas direcciones (WhatsApp recibe y envía mensajes)

---

## Arquitectura

```
┌─────────────────────────────────────────────────────────────────┐
│                    CATÁLOGO DE INTEGRACIONES                    │
│                     (integration_types)                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  📦 E-commerce        💳 Facturación      📧 Mensajería         │
│  ─────────────        ──────────────      ─────────────         │
│  • Shopify            • Softpymes         • WhatsApp            │
│  • MercadoLibre       • Alegra            • Telegram            │
│  • Amazon             • Siigo             • SMS                 │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
                            │
                            │ Cada negocio "instala" lo que necesita
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│              INSTANCIAS CONFIGURADAS (integrations)             │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Mi Tiendita (business_id=1):                                   │
│    ✓ Shopify (#1) - {api_key: "sk_prod_..."}                   │
│    ✓ Softpymes (#2) - {api_key: "...", nit: "900123456"}       │
│    ✓ WhatsApp (#5) - {phone_id: "...", token: "..."}           │
│                                                                 │
│  Tu Negocio (business_id=2):                                    │
│    ✓ Shopify (#3) - {api_key: "sk_prod_different..."}          │
│    ✓ Alegra (#4) - {token: "..."}                              │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## Categorías de Integraciones

### 📦 E-commerce (Inbound)

**Propósito**: Recibir órdenes de plataformas de venta.

| Código | Nombre | Estado | Descripción |
|--------|--------|--------|-------------|
| `shopify` | Shopify | ✅ Activo | Tienda online con webhooks |
| `mercadolibre` | MercadoLibre | ✅ Activo | Marketplace LATAM |
| `amazon` | Amazon | 🔜 Próximamente | Marketplace global |
| `woocommerce` | WooCommerce | 🔜 Próximamente | Plugin WordPress |

**Flujo**:
```
Cliente compra → Webhook → Probability → Crea orden → Notifica al negocio
```

---

### 💳 Facturación Electrónica (Outbound)

**Propósito**: Emitir facturas electrónicas ante autoridades fiscales (DIAN, SAT, etc.).

| Código | Nombre | País | Estado |
|--------|--------|------|--------|
| `softpymes` | Softpymes | 🇨🇴 Colombia | ✅ Activo |
| `alegra` | Alegra | 🇨🇴🇲🇽🇵🇪 Multi | 🔜 Próximamente |
| `siigo` | Siigo | 🇨🇴 Colombia | 🔜 Próximamente |
| `facturama` | Facturama | 🇲🇽 México | 🔜 Próximamente |

**Flujo**:
```
Orden creada → Probability → Genera factura → Softpymes API → DIAN → Cliente recibe factura
```

**Configuración Especial**:
- Ver módulo `invoicing` para configurar facturación automática
- Se vincula una integración de e-commerce con una de facturación vía `invoicing_configs`

---

### 📧 Mensajería (Bidirectional)

**Propósito**: Comunicación con clientes (notificaciones, soporte).

| Código | Nombre | Estado | Descripción |
|--------|--------|--------|-------------|
| `whatsapp` | WhatsApp Business | ✅ Activo | Meta Business API |
| `telegram` | Telegram Bot | 🔜 Próximamente | Bot API |
| `sms` | SMS Gateway | 🔜 Próximamente | Twilio/AWS SNS |

**Flujo**:
```
Outbound: Orden pagada → Probability → WhatsApp → Cliente recibe mensaje
Inbound:  Cliente pregunta → WhatsApp → Probability → Bot responde
```

---

### 🚚 Logística (Futuro)

| Código | Nombre | Estado |
|--------|--------|--------|
| `fedex` | FedEx | 🔜 Planeado |
| `dhl` | DHL Express | 🔜 Planeado |
| `coordinadora` | Coordinadora (CO) | 🔜 Planeado |

---

### 💰 Pagos (Futuro)

| Código | Nombre | Estado |
|--------|--------|--------|
| `stripe` | Stripe | 🔜 Planeado |
| `paypal` | PayPal | 🔜 Planeado |
| `wompi` | Wompi (CO) | 🔜 Planeado |

---

## Flujo Completo

### Ejemplo: Facturar Órdenes de Shopify con Softpymes

#### **Paso 1: Conectar Shopify**

```bash
POST /api/integrations
{
  "business_id": 1,
  "integration_type_id": 1,  # Shopify
  "name": "Mi Tiendita - Shopify",
  "credentials": {
    "api_key": "shpat_...",
    "api_secret": "shpss_...",
    "shop_domain": "mitiendita.myshopify.com"
  },
  "config": {
    "sync_products": true,
    "sync_orders": true
  }
}
```

**Respuesta**:
```json
{
  "id": 1,
  "integration_type": {
    "code": "shopify",
    "category": "ecommerce"
  },
  "is_active": true
}
```

---

#### **Paso 2: Conectar Softpymes**

```bash
POST /api/integrations
{
  "business_id": 1,
  "integration_type_id": 4,  # Softpymes
  "name": "Mi Tiendita - Softpymes",
  "credentials": {
    "api_key": "sk_live_...",
    "secret_key": "sk_secret_...",
    "company_nit": "900123456-7"
  },
  "config": {
    "max_retries": 3,
    "auto_send_email": true
  }
}
```

**Respuesta**:
```json
{
  "id": 2,
  "integration_type": {
    "code": "softpymes",
    "category": "invoicing"
  },
  "is_active": true
}
```

---

#### **Paso 3: Vincular Shopify con Softpymes**

```bash
POST /api/invoicing/configs
{
  "business_id": 1,
  "source_integration_id": 1,      # Shopify
  "invoicing_integration_id": 2,   # Softpymes
  "enabled": true,
  "auto_invoice": true,
  "filters": {
    "min_amount": 50000,
    "only_paid": true
  }
}
```

**Resultado**: Ahora las órdenes de Shopify se facturan automáticamente con Softpymes.

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

**Filtros**:
```
?category=ecommerce        # Filtrar por categoría
&direction=inbound         # Filtrar por dirección
&is_active=true            # Solo activos
```

---

### Integrations

```http
GET    /api/integrations                    # Listar integraciones del negocio
GET    /api/integrations/:id                # Obtener integración
POST   /api/integrations                    # Crear integración
PUT    /api/integrations/:id                # Actualizar integración
DELETE /api/integrations/:id                # Eliminar integración
POST   /api/integrations/:id/test           # Probar conexión
POST   /api/integrations/:id/activate       # Activar integración
POST   /api/integrations/:id/deactivate     # Desactivar integración
```

**Filtros**:
```
?business_id=1             # Filtrar por negocio
&category=invoicing        # Filtrar por categoría
&is_active=true            # Solo activas
&integration_type_id=4     # Filtrar por tipo
```

---

## Ejemplos de Uso

### 1. Listar Integraciones Disponibles

```bash
GET /api/integrations/types?category=invoicing

Response:
[
  {
    "id": 4,
    "code": "softpymes",
    "name": "Softpymes",
    "category": "invoicing",
    "direction": "outbound",
    "icon": "https://cdn.probability.com/integrations/softpymes.svg",
    "description": "Proveedor de facturación electrónica para Colombia (DIAN)",
    "is_active": true,
    "supported_countries": ["CO"]
  },
  {
    "id": 5,
    "code": "alegra",
    "name": "Alegra",
    "category": "invoicing",
    "direction": "outbound",
    "is_active": false
  }
]
```

---

### 2. Conectar Nueva Integración

```bash
POST /api/integrations
{
  "business_id": 1,
  "integration_type_id": 4,
  "name": "Softpymes - Producción",
  "credentials": {
    "api_key": "your_api_key",
    "company_nit": "900123456-7"
  }
}

Response:
{
  "id": 10,
  "business_id": 1,
  "integration_type": {
    "id": 4,
    "code": "softpymes",
    "name": "Softpymes",
    "category": "invoicing"
  },
  "name": "Softpymes - Producción",
  "is_active": true,
  "created_at": "2026-01-31T10:00:00Z"
}
```

---

### 3. Listar Mis Integraciones

```bash
GET /api/integrations?business_id=1

Response:
[
  {
    "id": 1,
    "name": "Mi Tiendita - Shopify",
    "integration_type": {
      "code": "shopify",
      "category": "ecommerce"
    },
    "is_active": true,
    "last_sync": "2026-01-31T09:45:00Z"
  },
  {
    "id": 2,
    "name": "Mi Tiendita - Softpymes",
    "integration_type": {
      "code": "softpymes",
      "category": "invoicing"
    },
    "is_active": true
  }
]
```

---

### 4. Probar Conexión

```bash
POST /api/integrations/2/test

Response:
{
  "success": true,
  "message": "Conexión exitosa con Softpymes",
  "details": {
    "company_name": "Mi Tiendita SAS",
    "nit": "900123456-7",
    "environment": "production"
  }
}
```

---

## Base de Datos

### Tabla: `integration_types`

```sql
CREATE TABLE integration_types (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,           -- 'shopify', 'softpymes', etc.
    name VARCHAR(100) NOT NULL,                 -- 'Shopify', 'Softpymes', etc.
    category VARCHAR(50) NOT NULL,              -- 'ecommerce', 'invoicing', etc.
    direction VARCHAR(20) NOT NULL,             -- 'inbound', 'outbound', 'bidirectional'
    description TEXT,
    icon VARCHAR(255),                          -- URL del icono
    image_url VARCHAR(255),                     -- URL de imagen de portada
    is_active BOOLEAN DEFAULT true,
    config_schema JSONB,                        -- Esquema JSON de configuración
    credentials_schema JSONB,                   -- Esquema JSON de credenciales
    api_base_url VARCHAR(255),                  -- URL base del API
    documentation_url VARCHAR(255),             -- URL de documentación
    supported_countries TEXT[],                 -- ['CO', 'MX', 'PE']
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_integration_types_category ON integration_types(category);
CREATE INDEX idx_integration_types_active ON integration_types(is_active);
```

---

### Tabla: `integrations`

```sql
CREATE TABLE integrations (
    id SERIAL PRIMARY KEY,
    business_id INTEGER NOT NULL,
    integration_type_id INTEGER NOT NULL REFERENCES integration_types(id),
    name VARCHAR(255) NOT NULL,                 -- "Mi Tiendita - Shopify"
    code VARCHAR(50),                           -- Código único opcional
    store_id VARCHAR(100),                      -- ID de la tienda externa
    is_active BOOLEAN DEFAULT true,
    is_default BOOLEAN DEFAULT false,           -- ¿Es la integración por defecto?
    config JSONB,                               -- Configuración específica
    credentials JSONB,                          -- Credenciales encriptadas
    description TEXT,
    last_sync_at TIMESTAMP,                     -- Última sincronización
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

## Relación con Otros Módulos

### Módulo `invoicing`

El módulo de facturación usa `integrations` para configurar proveedores de facturación.

**Tabla de vinculación**: `invoicing_configs`

```sql
CREATE TABLE invoicing_configs (
    id SERIAL PRIMARY KEY,
    business_id INTEGER NOT NULL,
    source_integration_id INTEGER NOT NULL,      -- FK a integrations (Shopify, MeLi)
    invoicing_integration_id INTEGER NOT NULL,   -- FK a integrations (Softpymes, Alegra)
    enabled BOOLEAN DEFAULT true,
    auto_invoice BOOLEAN DEFAULT false,
    filters JSONB,                               -- Filtros de facturación
    created_at TIMESTAMP DEFAULT NOW()
);
```

**Relación**:
- `source_integration_id` → Integración de e-commerce (category='ecommerce')
- `invoicing_integration_id` → Integración de facturación (category='invoicing')

---

### Módulo `orders`

Las órdenes guardan de qué integración provienen:

```sql
CREATE TABLE orders (
    id UUID PRIMARY KEY,
    integration_id INTEGER REFERENCES integrations(id),
    business_id INTEGER NOT NULL,
    -- ...
);
```

---

## Expansión Futura

### Próximas Categorías

Ver sección [Planificación de Categorías](#planificación-de-categorías) para detalles de expansión.

| Categoría | Prioridad | Estado |
|-----------|-----------|--------|
| `ecommerce` | Alta | ✅ Implementado |
| `invoicing` | Alta | ✅ Implementado |
| `messaging` | Media | ✅ Implementado |
| `payment` | Media | 🔜 Q2 2026 |
| `shipping` | Media | 🔜 Q2 2026 |
| `accounting` | Baja | 🔜 Q3 2026 |
| `analytics` | Baja | 🔜 Q4 2026 |

---

## Arquitectura Hexagonal

Este módulo sigue arquitectura hexagonal:

```
core/
├── bundle.go
└── internal/
    ├── domain/
    │   ├── entities.go       # IntegrationType, Integration
    │   ├── ports.go          # Interfaces
    │   ├── dtos.go           # DTOs
    │   └── enums.go          # Constantes
    ├── app/
    │   ├── usecaseintegrations/
    │   └── usecaseintegrationtype/
    └── infra/
        ├── primary/
        │   └── handlers/
        └── secondary/
            ├── repository/
            └── encryption/
```

---

## Seguridad

### Encriptación de Credenciales

Todas las credenciales se guardan **encriptadas** usando AES-256.

```go
// Al crear integración
credentials, _ := encryption.Encrypt(rawCredentials)
integration.Credentials = credentials

// Al usar integración
rawCredentials, _ := encryption.Decrypt(integration.Credentials)
```

**Variable de entorno requerida**: `ENCRYPTION_KEY`

---

## Testing

```bash
# Tests unitarios
go test ./internal/domain/...
go test ./internal/app/...

# Tests de integración
go test ./internal/infra/...

# Test end-to-end
go test ./... -tags=e2e
```

---

## Contribuir

Ver archivo `/.claude/rules/architecture.md` para reglas de arquitectura hexagonal.

---

**Última actualización**: 2026-01-31
