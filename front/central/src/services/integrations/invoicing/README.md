# Modulo de Integraciones de Facturacion Electronica

## Arquitectura General

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           FRONTEND (Next.js)                            │
│                                                                         │
│  ConfigForm / EditForm  →  Server Actions  →  Backend API               │
│  (credenciales, config,     (integrations/core)                         │
│   is_testing toggle)                                                    │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                        BACKEND (Go - Hexagonal)                         │
│                                                                         │
│  ┌──────────────────┐    ┌──────────────┐    ┌───────────────────────┐  │
│  │ modules/invoicing │───▶│   ROUTER     │───▶│ Provider Consumers    │  │
│  │   (core module)   │    │ (by provider │    │  ┌─────────────────┐  │  │
│  │                   │    │   field)     │    │  │   Softpymes     │  │  │
│  │ CreateInvoice()   │    └──────────────┘    │  │   Siigo         │  │  │
│  │ RetryInvoice()    │                        │  │   Factus        │  │  │
│  │ CreateJournal()   │    ┌──────────────┐    │  └─────────────────┘  │  │
│  │ CancelInvoice()   │◀───│  Response    │◀───│                       │  │
│  │ CompareInvoices() │    │  Consumer    │    │  Cada consumer llama  │  │
│  │ ListItems()       │    └──────────────┘    │  a su API externa y   │  │
│  └──────────────────┘                        │  publica response     │  │
│                                              └───────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Flujo Completo de una Factura

### 1. Publicacion (modules/invoicing)

```
CreateInvoice(orderID)
  │
  ├─ 1. GetOrderByID()
  ├─ 2. GetConfigByIntegration() → fallback GetEnabledConfigByBusiness()
  ├─ 3. Validar: enabled, invoiceable, filtros
  ├─ 4. resolveProvider(integrationID) → "softpymes" | "siigo" | "factus"
  ├─ 5. Crear invoice record (status: pending)
  ├─ 6. Crear invoice items
  ├─ 7. Crear sync log
  ├─ 8. Publicar a cola: invoicing.requests
  │     { provider: "siigo", operation: "create", invoice_data: {...} }
  │
  └─ 9. Si provider == "siigo" && enable_journal == true:
        └─ CreateJournal() (auto-trigger, non-blocking)
```

### 2. Routing (integrations/invoicing/router)

```
invoicing.requests
  │
  ├─ Lee campo "provider" del mensaje (sin transformar payload)
  │
  ├─ "softpymes" → invoicing.softpymes.requests
  ├─ "factus"    → invoicing.factus.requests
  ├─ "siigo"     → invoicing.siigo.requests
  ├─ "alegra"    → invoicing.alegra.requests
  ├─ "world_office" → invoicing.world_office.requests
  └─ "helisa"    → invoicing.helisa.requests
```

### 3. Consumer del Proveedor (ej: Siigo)

```
invoicing.siigo.requests
  │
  ├─ GetIntegrationByID() → obtiene IsTesting, BaseURL, BaseURLTest
  ├─ DecryptCredentials() → username, access_key, account_id, partner_id
  │
  ├─ Resolver URL efectiva:
  │   if integration.IsTesting && BaseURLTest != "" → usa BaseURLTest (mock)
  │   else → usa apiURL de credenciales o default
  │
  ├─ Switch por operacion:
  │   ├─ "create" / "retry"    → processCreateInvoice()
  │   ├─ "create_journal"      → processCreateJournal()  (solo Siigo)
  │   └─ default               → error
  │
  └─ Publicar response a: invoicing.responses
```

### 4. Response Consumer (modules/invoicing)

```
invoicing.responses
  │
  ├─ Discriminar por operation:
  │   ├─ "compare"    → handleCompareResponse() → SSE
  │   ├─ "list_items" → handleListItemsResponse() → SSE
  │   └─ otros        → handleSuccess/Error/PendingValidation
  │
  ├─ handleSuccess():
  │   ├─ Actualizar invoice (status: issued, invoice_number, CUFE, URLs)
  │   ├─ Actualizar sync log (status: success)
  │   ├─ Si operation != "create_journal": actualizar order invoice info
  │   ├─ Publicar evento RabbitMQ: invoice.created
  │   └─ Publicar evento SSE: invoice.created
  │
  ├─ handlePendingValidation():
  │   ├─ Mantener invoice en pending (DIAN validando)
  │   ├─ Programar check_status con backoff (15min, 30min, 1h, 2h, 4h)
  │   └─ Publicar SSE: invoice.pending_validation
  │
  └─ handleError():
      ├─ Marcar invoice como failed
      ├─ Programar retry con backoff (5min, 15min, 30min)
      └─ Publicar SSE: invoice.failed
```

---

## Colas RabbitMQ

| Cola | Publicador | Consumidor | Proposito |
|------|-----------|------------|-----------|
| `invoicing.requests` | modules/invoicing | router | Cola unificada de entrada |
| `invoicing.softpymes.requests` | router | Softpymes consumer | Requests para Softpymes |
| `invoicing.factus.requests` | router | Factus consumer | Requests para Factus |
| `invoicing.siigo.requests` | router | Siigo consumer | Requests para Siigo |
| `invoicing.alegra.requests` | router | Alegra consumer | Requests para Alegra (TODO) |
| `invoicing.world_office.requests` | router | World Office consumer | Requests para World Office (TODO) |
| `invoicing.helisa.requests` | router | Helisa consumer | Requests para Helisa (TODO) |
| `invoicing.responses` | Todos los providers | Response consumer | Respuestas unificadas |

---

## Operaciones por Proveedor

| Operacion | Softpymes | Siigo | Factus |
|-----------|:---------:|:-----:|:------:|
| `create` | Si | Si | Si |
| `retry` | Si | Si | Si |
| `cancel` | Si | No | No |
| `check_status` | Si | No | No |
| `compare` | Si | No | No |
| `list_items` | Si | No | No |
| `create_journal` | No | Si | No |

---

## Resolucion de URL (is_testing)

La logica de testing NO la maneja el modulo central de invoicing. Cada consumer de proveedor la resuelve independientemente:

```
1. Frontend: Toggle "Modo de Pruebas" → is_testing: true/false
2. Backend (integrations/core): Guarda is_testing en tabla integrations
3. Consumer del proveedor:
   integration = integrationCore.GetIntegrationByID(id)
   // ↑ Trae IsTesting, BaseURL, BaseURLTest del cache/DB

   effectiveURL = integration.BaseURL          // URL produccion
   if integration.IsTesting && integration.BaseURLTest != "" {
       effectiveURL = integration.BaseURLTest  // URL mock/sandbox
   }
4. El consumer pasa effectiveURL al cliente HTTP del proveedor
```

El modulo central de invoicing publica el mensaje con `provider` y `operation` — no envía la URL ni el flag `is_testing`. Cada consumer consulta el estado actual directamente via `integrationCore`.

### integration_types (tabla DB)

| type_id | code | base_url | base_url_test |
|---------|------|----------|---------------|
| 5 | softpymes | https://api.softpymes.com | http://back-testing:9090 |
| 7 | factus | https://api.factus.com.co | (vacio) |
| 8 | siigo | https://api.siigo.com | http://back-testing:9090 |

`base_url_test` apunta a un proyecto mock interno que simula las APIs de los proveedores para pruebas sin afectar DIAN.

---

## Cache de Configuraciones (Redis)

El modulo de invoicing cachea las configuraciones de facturacion en Redis para evitar consultas repetidas a la DB:

```
Claves:
  probability:invoicing:config:{integration_id}     → Config por integracion e-commerce
  probability:invoicing:config:business:{business_id} → Config activa del negocio

TTL: 1 hora

Invalidacion:
  - Al crear/actualizar/eliminar una config
  - Al habilitar/deshabilitar una config
  - Al cambiar auto_invoice

Fallback:
  - Si Redis no disponible → cache miss → consulta directa a DB
  - Si cache miss → consulta DB → cachea resultado
```

### Flujo con cache

```
CreateInvoice(orderID)
  └─ GetConfigByIntegration(integrationID)
      ├─ 1. Buscar en Redis: probability:invoicing:config:{integrationID}
      ├─ 2. Si hit → retornar config cacheada
      └─ 3. Si miss → consultar DB → cachear en Redis → retornar

  └─ Fallback: GetEnabledConfigByBusiness(businessID)
      ├─ 1. Buscar en Redis: probability:invoicing:config:business:{businessID}
      ├─ 2. Si hit → retornar
      └─ 3. Si miss → consultar DB → cachear → retornar
```

---

## Retry y Journal

### Retry de Facturas Fallidas

```
RetryInvoice(invoiceID)
  ├─ Validar status == "failed"
  ├─ Lock optimista: marcar como pending
  ├─ Crear nuevo sync log (retry_count + 1)
  ├─ Si invoice tiene metadata.original_operation == "create_journal":
  │   └─ Publicar con operation = "create_journal" (no "retry")
  └─ Si no: publicar con operation = "retry"
```

Esto garantiza que el retry de un journal llega al `case "create_journal"` del consumer Siigo, no al `case "retry"`.

### Journal (Comprobante Contable - Solo Siigo)

```
CreateJournal(orderID)
  ├─ Validar: provider == "siigo" && enable_journal == true
  ├─ Crear invoice record con metadata: {type: "journal", original_operation: "create_journal"}
  ├─ Crear items y sync log
  └─ Publicar con operation = "create_journal"

Response:
  ├─ handleSuccess() actualiza el invoice record del journal
  └─ NO actualiza order invoice info (skip para journals)
```

**Auto-trigger:** Al crear una factura exitosamente para un negocio Siigo con `enable_journal: true`, se dispara `CreateJournal()` automaticamente (non-blocking).

---

## Frontend - Estructura por Proveedor

```
invoicing/
├── softpymes/          (type_id: 5 - Completo)
│   ├── domain/types.ts
│   └── ui/components/
│       ├── SoftpymesConfigForm.tsx      Crear integracion
│       ├── SoftpymesEditForm.tsx        Editar integracion
│       └── SoftpymesIntegrationView.tsx Vista resumen (solo este provider)
│
├── factus/             (type_id: 7 - Completo)
│   ├── domain/types.ts
│   └── ui/components/
│       ├── FactusConfigForm.tsx
│       └── FactusEditForm.tsx
│
├── siigo/              (type_id: 8 - Completo)
│   ├── domain/types.ts
│   └── ui/components/
│       ├── SiigoConfigForm.tsx
│       └── SiigoEditForm.tsx
│
├── alegra/             (type_id: 9 - Esqueleto)
│   ├── domain/types.ts
│   └── ui/components/
│       ├── AlegraConfigForm.tsx
│       └── AlegraEditForm.tsx
│
├── world_office/       (type_id: 10 - Esqueleto)
│   ├── domain/types.ts
│   └── ui/components/
│       ├── WorldOfficeConfigForm.tsx
│       └── WorldOfficeEditForm.tsx
│
└── helisa/             (type_id: 11 - Esqueleto)
    ├── domain/types.ts
    └── ui/components/
        ├── HelisaConfigForm.tsx
        └── HelisaEditForm.tsx
```

### Patron comun de cada formulario

1. **Configuracion General**: Nombre + selector de negocio (super admin)
2. **Credenciales**: Campos especificos del proveedor + visibilidad de passwords
3. **Probar Conexion**: `testConnectionRawAction(providerCode, config, credentials)`
4. **Modo de Pruebas**: Toggle `is_testing` (Softpymes y Siigo) — muestra URL sandbox
5. **Acciones**: Crear / Actualizar + Cancelar

### Server Actions usadas

Todas desde `@/services/integrations/core/infra/actions`:

| Action | Uso |
|--------|-----|
| `createIntegrationAction()` | Crear nueva integracion |
| `updateIntegrationAction()` | Actualizar integracion existente |
| `testConnectionRawAction()` | Probar conexion sin crear |
| `getBusinessesSimpleAction()` | Listar negocios (super admin) |

---

## Credenciales por Proveedor

| Proveedor | Credenciales | Config Especifica |
|-----------|-------------|-------------------|
| Softpymes | api_key, api_secret | company_nit, company_name, referer, resolution_id, branch_code, seller_nit |
| Factus | client_id, client_secret, username, password | numbering_range_id, default_tax_rate, payment_form, payment_method_code |
| Siigo | username, access_key, account_id, partner_id | (via invoice_config JSONB en invoicing config) |
| Alegra | email, token | (pendiente) |
| World Office | username, password, company_code | (pendiente) |
| Helisa | username, password, company_id | (pendiente) |

---

## Estado de Implementacion

| Proveedor | Frontend | Backend Consumer | API Client | Operaciones |
|-----------|:--------:|:---------------:|:----------:|:-----------:|
| Softpymes | Completo | Completo | Completo | create, retry, cancel, check_status, compare, list_items |
| Factus | Completo | Completo | Completo | create, retry |
| Siigo | Completo | Completo | Completo | create, retry, create_journal |
| Alegra | Esqueleto | Esqueleto | Pendiente | - |
| World Office | Esqueleto | Esqueleto | Pendiente | - |
| Helisa | Esqueleto | Esqueleto | Pendiente | - |
