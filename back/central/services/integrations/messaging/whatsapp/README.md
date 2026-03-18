# WhatsApp Business Cloud API Integration

Integración completa y bidireccional con WhatsApp Business Cloud API para gestión de conversaciones con clientes.

## Cuenta Meta

| Campo | Valor |
|-------|-------|
| Porfolio Empresarial | **Probabilityapp** |
| App | **ProbabilityIA** (ID: `2812884712240202`) |
| WABA Producción | `1302830408357767` |
| WABA Test | `946521194991666` |
| Phone Number ID (Producción) | `1077369948787698` (+57 300 5636160) |
| Phone Number ID (Test) | `921919641007826` (+1 555 172 9816) |
| System User | `cam-adm` (ID: `61579689993959`) |
| PIN 2FA del número | `170122` |

> **IMPORTANTE**: Los tokens del system user se generan desde **Probabilityapp** (Business Settings → System Users → cam-adm → Generar identificador). El system user debe tener asignada la app ProbabilityIA y las cuentas de WhatsApp como activos.

## Características

- 12 plantillas conversacionales (11 de negocio + 1 de prueba de conexión)
- Botones interactivos (Quick Reply) para respuestas predefinidas
- Webhooks bidireccionales para recibir respuestas de usuarios en tiempo real
- Persistencia de conversaciones en PostgreSQL con ventana de 24h
- Máquina de estados para flujos conversacionales automatizados
- Eventos de negocio publicados en RabbitMQ
- Credenciales de plataforma compartidas (todas las integraciones WhatsApp usan las mismas credenciales del integration_type)
- Webhook autenticado via verify_token + HMAC-SHA256 (sin JWT — Meta no envía JWT)
- Credenciales leídas desde Redis cache via core (no depende de Redis directamente)

## API de Meta — Referencia Rápida

Todas las llamadas usan `https://graph.facebook.com/v22.0/` con header `Authorization: Bearer {ACCESS_TOKEN}`.

### Listar números de teléfono de un WABA

```bash
curl -s "https://graph.facebook.com/v22.0/{WABA_ID}/phone_numbers?fields=verified_name,display_phone_number,quality_rating,status,messaging_limit_tier,platform_type" \
  -H "Authorization: Bearer {TOKEN}" | jq .
```

### Ver suscripciones de webhook de la app

Requiere App Token (`APP_ID|APP_SECRET`):

```bash
curl -s "https://graph.facebook.com/v22.0/{APP_ID}/subscriptions?access_token={APP_ID}|{APP_SECRET}" | jq .
```

### Suscribir campos de webhook via API

```bash
curl -X POST "https://graph.facebook.com/v22.0/{APP_ID}/subscriptions" \
  -d "object=whatsapp_business_account" \
  -d "callback_url=https://www.probabilityia.com.co/api/v1/integrations/whatsapp/webhook" \
  -d "verify_token={VERIFY_TOKEN}" \
  -d "fields=messages,message_template_status_update" \
  -d "access_token={APP_ID}|{APP_SECRET}"
```

### Suscribir app a un WABA (recibir eventos de esa cuenta)

```bash
curl -X POST "https://graph.facebook.com/v22.0/{WABA_ID}/subscribed_apps" \
  -H "Authorization: Bearer {TOKEN}"
```

### Ver health status de un número

```bash
curl -s "https://graph.facebook.com/v22.0/{PHONE_NUMBER_ID}?fields=health_status,quality_score" \
  -H "Authorization: Bearer {TOKEN}" | jq .
```

### Registrar número de teléfono (cambiar estado de Pendiente a Connected)

```bash
curl -X POST "https://graph.facebook.com/v22.0/{PHONE_NUMBER_ID}/register" \
  -H "Authorization: Bearer {TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "messaging_product": "whatsapp",
    "pin": "123456"
  }'
```

El `pin` es un PIN de 6 dígitos que eliges — es para la verificación en dos pasos. Guardarlo.

### Verificar estado de un número

```bash
curl -s "https://graph.facebook.com/v22.0/{PHONE_NUMBER_ID}?fields=verified_name,display_phone_number,quality_rating,platform_type,status,name_status" \
  -H "Authorization: Bearer {TOKEN}" | jq .
```

### Debug de token (ver permisos y expiración)

```bash
curl -s "https://graph.facebook.com/v22.0/debug_token?input_token={TOKEN}&access_token={TOKEN}" | jq .
```

### Listar templates

```bash
curl -s "https://graph.facebook.com/v22.0/{WABA_ID}/message_templates?fields=name,status,language,category&limit=50" \
  -H "Authorization: Bearer {TOKEN}" | jq .
```

### Crear template via API

Mucho más rápido que hacerlo desde la interfaz gráfica de Meta.

#### Template simple (sin variables ni botones)

```bash
curl -X POST "https://graph.facebook.com/v22.0/{WABA_ID}/message_templates" \
  -H "Authorization: Bearer {TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "mi_template",
    "language": "es",
    "category": "UTILITY",
    "components": [
      {
        "type": "BODY",
        "text": "Texto del mensaje sin variables."
      }
    ]
  }'
```

#### Template con variables

Las variables se definen como `{{1}}`, `{{2}}`, etc. Se **requiere** un `example` con datos de ejemplo.

```bash
curl -X POST "https://graph.facebook.com/v22.0/{WABA_ID}/message_templates" \
  -H "Authorization: Bearer {TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "mi_template_con_vars",
    "language": "es",
    "category": "UTILITY",
    "components": [
      {
        "type": "BODY",
        "text": "Hola {{1}}, tu pedido {{2}} está listo.",
        "example": {
          "body_text": [["Juan", "ORD-001"]]
        }
      }
    ]
  }'
```

#### Template con botones Quick Reply

```bash
curl -X POST "https://graph.facebook.com/v22.0/{WABA_ID}/message_templates" \
  -H "Authorization: Bearer {TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "mi_template_con_botones",
    "language": "es",
    "category": "UTILITY",
    "components": [
      {
        "type": "BODY",
        "text": "Hola {{1}}, ¿confirmas tu pedido {{2}}?",
        "example": {
          "body_text": [["Juan", "ORD-001"]]
        }
      },
      {
        "type": "BUTTONS",
        "buttons": [
          {"type": "QUICK_REPLY", "text": "Sí, confirmar"},
          {"type": "QUICK_REPLY", "text": "No, cancelar"}
        ]
      }
    ]
  }'
```

#### Reglas importantes de templates

- **Categorías**: `UTILITY` (transaccional), `MARKETING`, `AUTHENTICATION`
- **UTILITY** se aprueba rápido (minutos a horas)
- **Variables**: si el ratio variables/texto es muy alto, Meta rechaza. Agregar más texto si falla
- **Botones Quick Reply**: máximo 3 botones, texto máximo 20 caracteres cada uno
- **`hello_world`**: solo funciona con números de test de Meta, NO con números empresariales
- Los templates nuevos quedan en `PENDING` hasta que Meta los aprueba

### Eliminar template

```bash
curl -X DELETE "https://graph.facebook.com/v22.0/{WABA_ID}/message_templates?name=nombre_template" \
  -H "Authorization: Bearer {TOKEN}"
```

### Enviar mensaje con template

```bash
curl -X POST "https://graph.facebook.com/v22.0/{PHONE_NUMBER_ID}/messages" \
  -H "Authorization: Bearer {TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "messaging_product": "whatsapp",
    "to": "573001234567",
    "type": "template",
    "template": {
      "name": "confirmacion_pedido_contraentrega",
      "language": {"code": "es"},
      "components": [
        {
          "type": "body",
          "parameters": [
            {"type": "text", "text": "Juan"},
            {"type": "text", "text": "Mi Tienda"},
            {"type": "text", "text": "ORD-001"},
            {"type": "text", "text": "Calle 45 #23-67"},
            {"type": "text", "text": "1x Camiseta"}
          ]
        }
      ]
    }
  }'
```

## Templates del Sistema

Todos los templates están definidos en `internal/domain/entities/template.go` y creados en el WABA `1302830408357767`.

### Formato API para crear cada template

#### prueba_conexion (test de conexión)

```json
{
  "name": "prueba_conexion",
  "language": "es",
  "category": "UTILITY",
  "components": [
    {"type": "BODY", "text": "¡Hola! Este es un mensaje de prueba de Probability. Tu conexión de WhatsApp está funcionando correctamente. ✅"}
  ]
}
```

#### confirmacion_pedido_contraentrega

```json
{
  "name": "confirmacion_pedido_contraentrega",
  "language": "es",
  "category": "UTILITY",
  "components": [
    {
      "type": "BODY",
      "text": "Hola {{1}}, tu pedido en {{2}} ha sido recibido.\n\nOrden: {{3}}\nDirección: {{4}}\nProductos: {{5}}\n\n¿Confirmas tu pedido?",
      "example": {"body_text": [["Juan", "Mi Tienda", "ORD-001", "Calle 45 #23-67 Bogotá", "1x Camiseta Negra"]]}
    },
    {
      "type": "BUTTONS",
      "buttons": [
        {"type": "QUICK_REPLY", "text": "Confirmar pedido"},
        {"type": "QUICK_REPLY", "text": "No confirmar"}
      ]
    }
  ]
}
```

#### pedido_confirmado_v2

```json
{
  "name": "pedido_confirmado_v2",
  "language": "es",
  "category": "UTILITY",
  "components": [
    {"type": "BODY", "text": "¡Hola {{1}}! 👋\nTu pedido ha sido *confirmado* exitosamente ✅\n\n📋 *Resumen de tu pedido:*\n🛒 *Pedido:* {{2}}\n🏪 *Tienda:* {{3}}\n📍 *Dirección de entrega:* {{4}}\n📦 *Productos:* {{5}}\n\nPronto estará en camino 🚚\n¡Gracias por tu compra! 🙏", "example": {"body_text": [["Juan Pérez", "ORD-001", "Mi Tienda", "Calle 45 #23-67 Bogotá", "1x Camiseta Negra, 2x Pantalón"]]}}
  ]
}
```

#### pedido_cancelado

```json
{
  "name": "pedido_cancelado",
  "language": "es",
  "category": "UTILITY",
  "components": [
    {"type": "BODY", "text": "Tu pedido {{1}} ha sido cancelado. Si tienes alguna duda, contáctanos.", "example": {"body_text": [["ORD-001"]]}}
  ]
}
```

#### menu_no_confirmacion

```json
{
  "name": "menu_no_confirmacion",
  "language": "es",
  "category": "UTILITY",
  "components": [
    {"type": "BODY", "text": "Entendemos que no deseas confirmar el pedido {{1}}. ¿Qué te gustaría hacer?", "example": {"body_text": [["ORD-001"]]}},
    {"type": "BUTTONS", "buttons": [{"type": "QUICK_REPLY", "text": "Presentar novedad"}, {"type": "QUICK_REPLY", "text": "Cancelar pedido"}, {"type": "QUICK_REPLY", "text": "Asesor"}]}
  ]
}
```

#### confirmar_cancelacion_pedido

```json
{
  "name": "confirmar_cancelacion_pedido",
  "language": "es",
  "category": "UTILITY",
  "components": [
    {"type": "BODY", "text": "¿Estás seguro de que deseas cancelar el pedido {{1}}? Esta acción no se puede deshacer.", "example": {"body_text": [["ORD-001"]]}},
    {"type": "BUTTONS", "buttons": [{"type": "QUICK_REPLY", "text": "Sí, cancelar"}, {"type": "QUICK_REPLY", "text": "No, volver"}]}
  ]
}
```

#### tipo_novedad_pedido

```json
{
  "name": "tipo_novedad_pedido",
  "language": "es",
  "category": "UTILITY",
  "components": [
    {"type": "BODY", "text": "Selecciona el tipo de novedad que deseas reportar para tu pedido:"},
    {"type": "BUTTONS", "buttons": [{"type": "QUICK_REPLY", "text": "Cambio de dirección"}, {"type": "QUICK_REPLY", "text": "Cambio de productos"}, {"type": "QUICK_REPLY", "text": "Cambio medio de pago"}]}
  ]
}
```

#### motivo_cancelacion_pedido

```json
{
  "name": "motivo_cancelacion_pedido",
  "language": "es",
  "category": "UTILITY",
  "components": [
    {"type": "BODY", "text": "Por favor, cuéntanos el motivo por el cual deseas cancelar tu pedido. Escribe tu respuesta a continuación:"}
  ]
}
```

#### handoff_asesor

```json
{
  "name": "handoff_asesor",
  "language": "es",
  "category": "UTILITY",
  "components": [
    {"type": "BODY", "text": "Te estamos conectando con un asesor. Por favor espera un momento, pronto te atenderemos."}
  ]
}
```

#### novedad_cambio_direccion

```json
{
  "name": "novedad_cambio_direccion",
  "language": "es",
  "category": "UTILITY",
  "components": [
    {"type": "BODY", "text": "Hemos recibido tu solicitud de cambio de dirección. Nuestro equipo la procesará a la brevedad."}
  ]
}
```

#### novedad_cambio_productos

```json
{
  "name": "novedad_cambio_productos",
  "language": "es",
  "category": "UTILITY",
  "components": [
    {"type": "BODY", "text": "Hemos recibido tu solicitud de cambio de productos. Nuestro equipo la revisará y te contactaremos pronto."}
  ]
}
```

#### novedad_cambio_medio_pago

```json
{
  "name": "novedad_cambio_medio_pago",
  "language": "es",
  "category": "UTILITY",
  "components": [
    {"type": "BODY", "text": "Hemos recibido tu solicitud de cambio de medio de pago. Nuestro equipo la procesará a la brevedad."}
  ]
}
```

#### alerta_servidor

```json
{
  "name": "alerta_servidor",
  "language": "es",
  "category": "UTILITY",
  "components": [
    {"type": "BODY", "text": "⚠️ Alerta del servidor - Tipo: {{1}}\n\nDetalle: {{2}}\n\nPor favor revisa el estado del sistema.", "example": {"body_text": [["RAM", "87.3% - supera umbral de 85%"]]}}
  ]
}
```

## Tabla resumen de templates

| Template | Variables | Botones | Uso |
|----------|-----------|---------|-----|
| `prueba_conexion` | — | — | Test de conexión |
| `confirmacion_pedido_contraentrega` | nombre, tienda, orden, dirección, productos | Confirmar / No confirmar | Inicio del flujo |
| `pedido_confirmado_v2` | nombre, #pedido, tienda, dirección, productos | — | Confirmación exitosa (con resumen e iconos) |
| `pedido_cancelado` | #pedido | — | Cancelación completada |
| `menu_no_confirmacion` | #pedido | Novedad / Cancelar / Asesor | Menú al no confirmar |
| `confirmar_cancelacion_pedido` | #pedido | Sí cancelar / No volver | Confirmar cancelación |
| `tipo_novedad_pedido` | — | Dirección / Productos / Pago | Tipo de novedad |
| `motivo_cancelacion_pedido` | — | — | Pide motivo (texto libre) |
| `handoff_asesor` | — | — | Derivar a humano |
| `novedad_cambio_direccion` | — | — | ACK cambio dirección |
| `novedad_cambio_productos` | — | — | ACK cambio productos |
| `novedad_cambio_medio_pago` | — | — | ACK cambio medio pago |
| `alerta_servidor` | tipo, detalle | — | Alerta de monitoreo |

## Flujo de Conversación

```
[INICIO] → confirmacion_pedido_contraentrega
            ├─ "Confirmar pedido" → pedido_confirmado [FIN]
            └─ "No confirmar" → menu_no_confirmacion
                                 ├─ "Presentar novedad" → tipo_novedad_pedido
                                 │                          ├─ "Cambio de dirección" → novedad_cambio_direccion [FIN]
                                 │                          ├─ "Cambio de productos" → novedad_cambio_productos [FIN]
                                 │                          └─ "Cambio medio de pago" → novedad_cambio_medio_pago [FIN]
                                 ├─ "Cancelar pedido" → confirmar_cancelacion_pedido
                                 │                       ├─ "Sí, cancelar" → motivo_cancelacion_pedido
                                 │                       │                    └─ [texto libre] → pedido_cancelado [FIN]
                                 │                       └─ "No, volver" → menu_no_confirmacion
                                 └─ "Asesor" → handoff_asesor [FIN - HUMANO]
```

## Arquitectura

### Credenciales de plataforma (compartidas)

Todas las integraciones WhatsApp (una por business) usan las mismas credenciales configuradas en el **integration_type** WhatsApp (ID: 2):
- `phone_number_id` — ID del número en Meta
- `access_token` — Token permanente del system user
- `verify_token` — Para verificación de webhook
- `webhook_secret` — Para validación HMAC-SHA256

Estas credenciales se cachean en Redis (`integration:platform_creds:2`) durante el warmup del servidor y se leen via `core.GetCachedPlatformCredentials()`.

### Webhook (sin JWT)

Las rutas del webhook (`GET/POST /api/v1/integrations/whatsapp/webhook`) **no usan JWT** porque Meta envía su propia autenticación:
- **GET** (verificación): Meta envía `hub.verify_token` → se compara con `verify_token` de platform_creds
- **POST** (eventos): Meta firma el payload con HMAC-SHA256 → se valida con `webhook_secret` de platform_creds

### Estructura de Directorios

```
whatsApp/
├── bundle.go                           # Ensamblaje de componentes (DI)
├── internal/
│   ├── domain/
│   │   ├── entities/
│   │   │   ├── conversation.go         # Entidades de conversación y estados
│   │   │   ├── message.go             # Entidades de mensajes
│   │   │   └── template.go            # Catálogo de 12 plantillas
│   │   ├── dtos/
│   │   ├── ports/
│   │   │   └── ports.go               # Interfaces (incluye IPlatformCredentialsGetter)
│   │   └── errors/
│   ├── app/
│   │   ├── usecasemessaging/
│   │   │   ├── send-template-message.go    # Envío de plantillas
│   │   │   ├── handle-webhook.go           # Procesamiento de webhooks
│   │   │   └── conversation-manager.go     # Máquina de estados
│   │   └── usecasetestconnection/
│   │       └── test-connection.go      # Test con template prueba_conexion
│   └── infra/
│       ├── primary/
│       │   ├── handlers/
│       │   │   ├── routes.go           # Webhook SIN JWT, send-template CON JWT
│       │   │   ├── webhook_handler.go  # Lee verify_token/secret de Redis via core
│       │   │   └── request/
│       │   │       └── webhook_payload.go  # Structs alineados con formato real de Meta
│       │   └── queue/
│       │       ├── consumerorder/      # Consume orders.confirmation.requested
│       │       └── consumeralert/      # Consume alertas de monitoreo
│       └── secondary/
│           ├── client/                 # Cliente HTTP WhatsApp API
│           ├── cache/
│           │   └── credentials_cache.go  # Siempre usa platform_creds (no propias)
│           └── queue/                  # Publishers RabbitMQ
```

## Webhook Payload de Meta

El struct `request/webhook_payload.go` está alineado con el formato real que Meta envía:

```json
{
  "object": "whatsapp_business_account",
  "entry": [{
    "id": "WABA_ID",
    "changes": [{
      "value": {
        "messaging_product": "whatsapp",
        "metadata": {
          "display_phone_number": "15551729816",
          "phone_number_id": "921919641007826"
        },
        "statuses": [{
          "id": "wamid.xxx",
          "status": "sent",
          "timestamp": "1773621987",
          "recipient_id": "573023406789",
          "recipient_logical_id": "160181375791296",
          "conversation": {
            "id": "abc123",
            "expiration_timestamp": "1773621988",
            "origin": {"type": "utility"}
          },
          "pricing": {
            "billable": true,
            "pricing_model": "PMP",
            "category": "utility",
            "type": "regular"
          }
        }]
      },
      "field": "messages"
    }]
  }]
}
```

**Nota**: `origin` es un **objeto** `{"type": "..."}`, no un string. El mapper extrae `origin.type` al DTO de dominio.

## Eventos de RabbitMQ

| Cola | Evento | Cuándo |
|------|--------|--------|
| `orders.confirmation.requested` | Enviar confirmación | Orden creada con config WhatsApp |
| `orders.whatsapp.confirmed` | Pedido confirmado | Usuario presiona "Confirmar pedido" |
| `orders.whatsapp.cancelled` | Pedido cancelado | Usuario completa flujo de cancelación |
| `orders.whatsapp.novelty` | Novedad solicitada | Usuario selecciona tipo de novedad |
| `customer.whatsapp.handoff` | Handoff a humano | Usuario presiona "Asesor" |

## Troubleshooting

### Token sin permisos para un WABA

El token del system user solo tiene acceso a los WABAs asignados como activos. Si da error de permisos:
1. Business Settings → System Users → cam-adm → Activos asignados
2. Verificar que el WABA esté asignado con control total
3. Regenerar token

### Template rechazado por "too many variables"

Meta rechaza templates donde el ratio variables/texto es alto. Solución: agregar más texto al body.

### `hello_world` no funciona con número empresarial

Es normal — `hello_world` solo funciona con números de test de Meta. Usar `prueba_conexion` para testing.

### Webhook da 401

Las rutas del webhook NO deben pasar por JWT. Verificar que `routes.go` registre GET/POST `/webhook` sin `middleware.JWT()`.

### Credenciales no encontradas en cache

El warmup del servidor debe cachear `platform_creds` en Redis. Verificar que `warm_cache.go` ejecute `warmPlatformCredentials()`.
