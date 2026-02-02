# 🛍️ Workflow de Testing de Shopify

## 📊 Arquitectura del Sistema

```
┌─────────────────────────────────────────────────────────────┐
│                    FLUJO COMPLETO DE TESTING                 │
└─────────────────────────────────────────────────────────────┘

   ┌──────────────────────┐
   │  Integration Test    │  Puerto: N/A (no es servidor)
   │   (Simulador CLI)    │  Ubicación: /back/integrationTest
   │                      │
   │  - Menú interactivo  │
   │  - Simula webhooks   │
   │  - Envía a Central   │
   └──────┬───────────────┘
          │
          │ HTTP POST (webhooks)
          │ http://localhost:3050/api/v1/integrations/shopify/webhook
          │
          ↓
   ┌──────────────────────┐
   │   Central Backend    │  Puerto: 3050
   │   (Sistema Real)     │  Ubicación: /back/central
   │                      │
   │  - Recibe webhooks   │
   │  - ⚠️ NO PROCESA ⚠️  │ ← PROBLEMA ACTUAL
   │  - Solo responde OK  │
   └──────────────────────┘
```

## ✅ ESTADO ACTUAL - WEBHOOK RECONECTADO (2 Febrero 2026)

**El webhook handler ESTÁ procesando las órdenes correctamente:**

```go
// /back/central/services/integrations/ecommerce/shopify/internal/infra/primary/handlers/webhook.go

// Respond 200 OK as fast as possible as required by Shopify
c.JSON(http.StatusOK, response.WebhookResponse{
    Success: true,
    Message: "Recibido",
})

// ✅ Procesar el webhook de forma asíncrona
go h.processWebhookAsync(headers.Topic, headers.ShopDomain, bodyBytes)
```

**Resultado:**
- ✅ Webhook llega correctamente (200 OK)
- ✅ Payload tiene estructura correcta de Shopify
- ✅ Webhook SE PROCESA en goroutine (no bloquea respuesta)
- ✅ Orden SE PUBLICA a RabbitMQ (`probability.orders.canonical`)
- ✅ Order Consumer la procesa
- ✅ Orden APARECE en la base de datos

**Nota histórica:** El procesamiento estuvo desconectado desde el 26 de enero 2026 (commit 18c9e87 por Juan Sebastian Mendoza) y fue reconectado el 2 de febrero 2026.

## 🎯 Variables de Entorno Configuradas

### Integration Test (`.env`)

```env
# El simulador envía webhooks a esta URL
WEBHOOK_BASE_URL=http://localhost:3050

# Dominio de la tienda de prueba
SHOPIFY_SHOP_DOMAIN=tienda-test.myshopify.com

# Secret para firmar webhooks (hardcoded en el simulador por ahora)
# Actualmente usa: "test_secret_key_for_integration_tests"

# Versión de API de Shopify
SHOPIFY_API_VERSION=2024-10
```

### Central Backend (`.env`)

**⚠️ ACTUALMENTE NO CONFIGURADO:**

```env
# Secret para validar webhooks de Shopify
# Si no se configura, se omite la validación HMAC (solo para desarrollo)
# SHOPIFY_API_SECRET=test_secret_key_for_integration_tests

# Cola de RabbitMQ para órdenes canónicas
RABBITMQ_ORDERS_CANONICAL_QUEUE=probability.orders.canonical
```

**Estado actual en logs:**
```
🔐 Verificando HMAC has_secret=false secret_prefix=NO_SECRET
⚠️ SHOPIFY_API_SECRET no configurado - omitiendo validación HMAC
```

## 🚀 Cómo Usar

### 1. Iniciar Central Backend

```bash
cd /home/cam/Desktop/probability/back/central
go run cmd/main.go
```

**Verificar que esté corriendo:**
```bash
curl http://localhost:3050/health
# Debe responder: 200 OK
```

### 2. Iniciar Integration Test (Simulador)

```bash
cd /home/cam/Desktop/probability/back/integrationTest
go run cmd/main.go
```

**Menú que aparece:**
```
=== Simulador de Webhooks ===

📦 SHOPIFY:
1. orders/create (crear nueva orden aleatoria)
2. orders/paid (marcar orden como pagada)
3. orders/updated (actualizar orden existente)
4. orders/cancelled (cancelar orden)
5. orders/fulfilled (marcar orden como cumplida)
6. orders/partially_fulfilled (cumplimiento parcial)
7. Listar órdenes almacenadas

💬 WHATSAPP:
8. Simular respuesta de usuario (manual)
...

0. Salir
```

### 3. Simular Creación de Orden (Opción 1)

**Ejemplo de uso:**

```
Opción: 1
```

**Lo que sucede internamente:**

1. **Genera orden aleatoria** con datos realistas:
   ```json
   {
     "id": 6386354755,
     "admin_graphql_api_id": "gid://shopify/Order/6386354755",
     "name": "#1001",
     "order_number": 1001,
     "email": "cliente@example.com",
     "currency": "COP",
     "total_price": "125000.00",
     "financial_status": "pending",
     "fulfillment_status": null,
     "line_items": [...],
     "customer": {...},
     "shipping_address": {...}
   }
   ```

2. **Calcula HMAC-SHA256** del payload
   ```go
   secret := "test_secret_key_for_integration_tests"
   mac := hmac.New(sha256.New, []byte(secret))
   mac.Write(payloadBytes)
   hmac := base64.StdEncoding.EncodeToString(mac.Sum(nil))
   ```

3. **Envía webhook** al central
   ```
   POST http://localhost:3050/api/v1/integrations/shopify/webhook
   Headers:
     Content-Type: application/json
     X-Shopify-Topic: orders/create
     X-Shopify-Shop-Domain: tienda-test.myshopify.com
     X-Shopify-Hmac-Sha256: <hmac>
     X-Shopify-API-Version: 2024-10
     X-Shopify-Webhook-Id: test-1738456789
   Body: <json_payload>
   ```

4. **Central responde** 200 OK
   ```
   ✅ Webhook recibido
   📦 Payload del webhook (7766 bytes)
   ⚠️ Webhook aceptado pero procesamiento asíncrono no implementado
   ```

5. **⚠️ Orden NO se procesa** - El webhook handler solo responde OK pero no hace nada más

### 4. Ver Logs Detallados

**Logs del Central:**
```bash
cd /home/cam/Desktop/probability/back/central
go run cmd/main.go

# Output:
02-02 01:14:21 INF Webhook recibido de Shopify function=WebhookHandler hmac=NDWPWN2/m/bJr/RSIiH5mMMGphF7L21CZXgwmL0uF0M= module=shopify shop_domain=tienda-test.myshopify.com topic=orders/create
02-02 01:14:21 INF 📦 Payload del webhook function=WebhookHandler module=shopify payload_preview="{\"id\":6386354755,\"admin_graphql_api_id\":\"gid://shopify/Order/6386354755\"..." payload_size=7766 shop_domain=tienda-test.myshopify.com topic=orders/create
02-02 01:14:21 INF 🔐 Verificando HMAC function=WebhookHandler has_secret=false module=shopify secret_prefix=NO_SECRET
02-02 01:14:21 WRN ⚠️ SHOPIFY_API_SECRET no configurado - omitiendo validación HMAC function=Next module=shopify
02-02 01:14:21 INF ⚠️ Webhook aceptado pero procesamiento asíncrono no implementado - payload NO se procesó function=WebhookHandler module=shopify shop_domain=tienda-test.myshopify.com topic=orders/create
 POST /api/v1/integrations/shopify/webhook 200 in 1ms
```

**Logs del Simulador:**
```
01:14:21 INF Orden aleatoria creada order_number=#1001
01:14:21 INF Enviando webhook shop_domain=tienda-test.myshopify.com topic=orders/create url=http://localhost:3050/api/v1/integrations/shopify/webhook
01:14:21 INF Webhook enviado exitosamente status_code=200 topic=orders/create
```

## 📋 Estructura de Datos Generados

### Orden Completa

El simulador genera órdenes con todos los campos requeridos por Shopify:

```json
{
  "id": 6386354755,
  "admin_graphql_api_id": "gid://shopify/Order/6386354755",
  "app_id": 134662,
  "browser_ip": "186.85.156.47",
  "buyer_accepts_marketing": true,
  "cancel_reason": null,
  "cancelled_at": null,
  "cart_token": "ccfada82dfa87a1",
  "checkout_id": 8789924299,
  "checkout_token": "ct_740d0c7610a8c542",
  "client_details": {
    "accept_language": "es",
    "browser_height": null,
    "browser_ip": "137.69.111.163",
    "browser_width": null,
    "session_hash": null,
    "user_agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X)"
  },
  "closed_at": null,
  "confirmed": true,
  "contact_email": "cliente@example.com",
  "created_at": "2026-02-02T00:50:00Z",
  "currency": "COP",
  "current_subtotal_price": "105042.02",
  "current_subtotal_price_set": {
    "shop_money": {
      "amount": "105042.02",
      "currency_code": "COP"
    },
    "presentment_money": {
      "amount": "105042.02",
      "currency_code": "COP"
    }
  },
  "current_total_price": "128165.00",
  "financial_status": "pending",
  "fulfillment_status": null,
  "line_items": [
    {
      "id": 123456789,
      "variant_id": 987654321,
      "title": "Producto Example",
      "quantity": 2,
      "price": "52521.01",
      "sku": "SKU-12345",
      "variant_title": "Talla M / Color Azul",
      "vendor": "Probability Store",
      "fulfillment_status": null,
      "requires_shipping": true,
      "taxable": true,
      "gift_card": false,
      "name": "Producto Example - Talla M / Color Azul",
      "properties": [],
      "product_exists": true
    }
  ],
  "customer": {
    "id": 123456789,
    "email": "cliente@example.com",
    "first_name": "Juan",
    "last_name": "Pérez",
    "phone": "+573001234567",
    "orders_count": 5,
    "total_spent": "450000.00",
    "default_address": {
      "first_name": "Juan",
      "last_name": "Pérez",
      "address1": "Calle 123 #45-67",
      "city": "Bogotá",
      "province": "Cundinamarca",
      "country": "Colombia",
      "zip": "110111",
      "phone": "+573001234567"
    }
  },
  "shipping_address": {
    "first_name": "Juan",
    "last_name": "Pérez",
    "address1": "Calle 123 #45-67",
    "address2": "Apto 301",
    "city": "Bogotá",
    "province": "Cundinamarca",
    "country": "Colombia",
    "zip": "110111",
    "phone": "+573001234567",
    "latitude": 4.6097,
    "longitude": -74.0817
  },
  "shipping_lines": [
    {
      "id": 123456789,
      "title": "Entrega Estándar CUNDINAMARCA (3 a 6 días hábiles municipios principales en Cundinamarca - 3 o más días a otros municipios)",
      "code": "standard_cundinamarca",
      "price": "3.15",
      "source": "shopify"
    }
  ],
  "tax_lines": [
    {
      "price": "19957.98",
      "rate": 0.19,
      "title": "IVA"
    }
  ]
}
```

### Campos Aleatorios Generados

| Campo | Valores Posibles |
|-------|------------------|
| `currency` | COP, USD, EUR |
| `source_name` | web, pos, mobile, api |
| `shipping_lines[0].title` | Entrega Estándar CUNDINAMARCA, Envío Express, Envío Gratis, Recogida en Tienda |
| `client_details.user_agent` | Windows, Mac, iPhone, Android |
| `financial_status` (create) | pending |
| `financial_status` (paid) | paid |
| `financial_status` (cancelled) | refunded |
| `fulfillment_status` (fulfilled) | fulfilled |
| `fulfillment_status` (partial) | partial |

## 🔐 Validación de Firma HMAC-SHA256

### Simulador (Generación)

```go
secret := "test_secret_key_for_integration_tests"
mac := hmac.New(sha256.New, []byte(secret))
mac.Write(payloadBytes)
signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

// Header enviado
X-Shopify-Hmac-Sha256: <signature>
```

### Central (Validación)

**⚠️ ACTUALMENTE DESHABILITADO:**

```go
// Si SHOPIFY_API_SECRET está vacío, se omite validación
if shopifySecret != "" {
    if !VerifyWebhookHMAC(bodyBytes, headers.Hmac, shopifySecret) {
        return "invalid signature"
    }
}
```

**Para habilitar validación:**

Agregar en `/back/central/.env`:
```env
SHOPIFY_API_SECRET=test_secret_key_for_integration_tests
```

## 🧪 Escenarios de Testing

### Escenario 1: Crear Nueva Orden

1. **Simulador opción 1:** `orders/create`
2. **Se genera orden** con `financial_status=pending`
3. **Webhook enviado** al central
4. **Central responde** 200 OK
5. **⚠️ Orden NO se procesa** - Solo se acepta pero no se guarda

**Verificación esperada (cuando se implemente procesamiento):**
```sql
SELECT * FROM orders WHERE external_id = '6386354755';
-- Debe aparecer la orden
```

### Escenario 2: Marcar Orden Como Pagada

1. **Simulador opción 2:** `orders/paid`
2. **Se actualiza orden existente** con `financial_status=paid`
3. **Webhook enviado** con topic `orders/paid`
4. **Central responde** 200 OK
5. **⚠️ Orden NO se procesa**

### Escenario 3: Cancelar Orden

1. **Simulador opción 4:** `orders/cancelled`
2. **Se actualiza orden** con `financial_status=refunded`, `cancel_reason=customer`
3. **Webhook enviado** con topic `orders/cancelled`
4. **Central responde** 200 OK
5. **⚠️ Orden NO se procesa**

### Escenario 4: Orden Cumplida

1. **Simulador opción 5:** `orders/fulfilled`
2. **Se actualiza orden** con `fulfillment_status=fulfilled`
3. **Se generan fulfillments** con tracking number
4. **Webhook enviado** con topic `orders/fulfilled`
5. **Central responde** 200 OK
6. **⚠️ Orden NO se procesa**

## 🛠️ Próximos Pasos para Implementar Procesamiento

Para que las órdenes se procesen correctamente, se debe modificar el webhook handler:

```go
// webhook.go - IMPLEMENTACIÓN PENDIENTE

// Respond 200 OK as fast as possible
c.JSON(http.StatusOK, response.WebhookResponse{
    Success: true,
    Message: "Recibido",
})

// Procesar de manera asíncrona (goroutine)
go func() {
    ctx := context.Background()

    // Parsear payload a domain.ShopifyOrder
    var order domain.ShopifyOrder
    if err := json.Unmarshal(bodyBytes, &order); err != nil {
        h.logger.Error(ctx).Err(err).Msg("Error al parsear orden de Shopify")
        return
    }

    // Llamar al use case correspondiente según el topic
    switch headers.Topic {
    case "orders/create":
        if err := h.useCase.CreateOrder(ctx, headers.ShopDomain, &order, bodyBytes); err != nil {
            h.logger.Error(ctx).Err(err).Msg("Error al crear orden")
        }
    case "orders/paid":
        if err := h.useCase.ProcessOrderPaid(ctx, headers.ShopDomain, &order); err != nil {
            h.logger.Error(ctx).Err(err).Msg("Error al procesar orden pagada")
        }
    // ... otros casos
    }
}()
```

## 🐛 Debugging

### Ver payload completo del webhook

```bash
# En el central, los logs muestran los primeros 500 caracteres
# Para ver el payload completo, modificar webhook.go:

payloadPreview := string(bodyBytes)  // Quitar el límite de 500 chars
```

### Verificar que el simulador genera datos correctos

```bash
cd /home/cam/Desktop/probability/back/integrationTest
go run cmd/main.go

# Opción 1: Crear orden
# Opción 7: Listar órdenes almacenadas
```

### Problemas comunes

| Problema | Causa | Solución |
|----------|-------|----------|
| "HMAC inválido" | Secret no coincide | Configurar mismo secret en ambos .env |
| "connection refused" | Central no está corriendo | Iniciar el central en puerto 3050 |
| Webhook responde OK pero no se procesa | Handler no implementado | Implementar procesamiento asíncrono en webhook.go |
| Orden no aparece en BD | Use case no se llama | Modificar webhook handler para llamar a use cases |

## ✅ Checklist de Verificación

Antes de empezar testing:

- [x] Central corriendo en puerto 3050
- [x] Simulador compilado sin errores
- [x] Webhook endpoint `/api/v1/integrations/shopify/webhook` implementado
- [x] Logs detallados agregados al webhook handler
- [ ] Variables de entorno configuradas (SHOPIFY_API_SECRET opcional)
- [ ] Webhook handler procesa órdenes (PENDIENTE DE IMPLEMENTAR)
- [ ] Use cases se llaman correctamente (PENDIENTE)
- [ ] Órdenes se publican a RabbitMQ (PENDIENTE)
- [ ] Órdenes se guardan en BD (PENDIENTE)

---

**Última actualización:** 2026-02-02
**Estado actual:** Webhook handler solo acepta webhooks pero NO los procesa
**Próximo paso:** Implementar procesamiento asíncrono en webhook.go
