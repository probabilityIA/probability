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

---

## Funcionalidades

### 1. Notificaciones automatizadas (bot de confirmación)

Cuando se crea un pedido con integración WhatsApp activa, el sistema envía automáticamente una plantilla de confirmación al cliente. El flujo soporta:

- Confirmación del pedido
- Cancelación con motivo libre
- Presentación de novedades (cambio de dirección, productos, medio de pago)
- Derivación a asesor humano (handoff)

Cada respuesta del cliente dispara una transición de estado y publica el evento de negocio correspondiente en RabbitMQ.

### 2. Chat bidireccional en tiempo real (dashboard)

El agente humano puede chatear directamente con el cliente desde el dashboard:

- **Lista de conversaciones** con filtros por estado y teléfono, paginada
- **Vista de mensajes** con historial completo y scroll automático
- **SSE (Server-Sent Events)**: los mensajes entrantes del cliente aparecen en tiempo real sin refrescar la página
- **Estados de mensajes en tiempo real**: ✓ enviado, ✓✓ entregado, ✓✓ azul leído — actualizados via SSE
- **Envío de respuestas manuales** con UI optimista (el mensaje aparece antes de la confirmación del servidor)
- **Alerta de ventana de 24h** cuando el cliente no ha respondido recientemente (Meta restringe mensajes libres)

### 3. Control IA Activa / Pausada

Permite controlar quién responde al cliente: la IA o el humano:

- **Toggle en el header del chat**: `🟢 IA Activa` → `🔴 IA Pausada`
- **Cuando IA está activa**: el compositor del chat está bloqueado — solo la IA responde
- **Cuando IA está pausada**: el compositor se habilita y el humano escribe libremente
- Al pausar la IA se activa automáticamente una `HumanSession` en Redis para ese teléfono
- El módulo `ai_sales` verifica el flag antes de responder y se detiene si hay una sesión humana activa
- TTL de 24h alineado con la ventana de servicio de WhatsApp

**Endpoints HTTP:**
```
POST /whatsapp/conversations/:id/pause-ai   # Pausa IA + activa HumanSession
POST /whatsapp/conversations/:id/resume-ai  # Reactiva IA
```

### 4. Agente AI Sales (integración)

Cuando un cliente escribe sin tener una conversación activa de bot en curso, el mensaje se reenvía al módulo `ai_sales`:

- El agente usa Amazon Bedrock (Nova) con tool use para buscar productos y crear órdenes
- Mantiene historial de sesión en Redis (últimos 20 mensajes)
- Respeta la pausa: si un humano tomó control (`whatsapp:ai_paused:{phone}`), el agente ignora el mensaje y retorna `nil`

### 5. Sesión Humana (HumanSession)

Mecanismo que enruta los mensajes entrantes del cliente al dashboard en lugar del bot:

- Clave Redis: `whatsapp:human_session:{phone}` con TTL 24h
- Se activa automáticamente al enviar una respuesta manual o al pausar la IA
- Cuando existe, los mensajes del cliente se persisten en BD y se publican via SSE al dashboard
- Si no existe (ni conversación activa), el mensaje se reenvía al agente AI Sales

---

## Arquitectura

### Redis — Claves de estado

| Clave | TTL | Propósito |
|-------|-----|-----------|
| `whatsapp:conv:{id}` | 25h | Datos de la conversación |
| `whatsapp:conv:idx:po:{phone}:{order}` | 25h | Índice phone+order → conv ID |
| `whatsapp:conv:idx:active:{phone}` | 25h | Índice teléfono activo → conv ID |
| `whatsapp:human_session:{phone}` | 24h | Sesión de atención humana activa |
| `whatsapp:ai_paused:{phone}` | 24h | Flag de IA pausada (humano en control) |

### Flujo de mensaje entrante

```
Cliente escribe
      │
      ▼
¿Conversación activa en Redis?
      ├─ SÍ → procesar flujo de bot → publicar SSE → persistir BD
      └─ NO
           │
           ▼
      ¿HumanSession activa?
           ├─ SÍ → persistir BD + publicar SSE al dashboard
           └─ NO
                │
                ▼
           ¿AIForwarder disponible?
                └─ SÍ → reenviar a ai_sales (que chequea IsAIPaused)
```

### Flujo de respuesta manual (humano)

```
Agente escribe en dashboard
      │
      ▼
SendManualReply → envía texto libre a WhatsApp API
      │
      ▼
ActivateHumanSession en Redis (TTL 24h)
      │
      ▼
Persiste en BD via RabbitMQ (persistence publisher)
```

### RabbitMQ — Eventos publicados

| Exchange / Cola | Evento | Cuándo |
|-----------------|--------|--------|
| `orders.confirmation.requested` | Enviar confirmación | Orden creada con config WhatsApp |
| `orders.whatsapp.confirmed` | Pedido confirmado | Usuario presiona "Confirmar pedido" |
| `orders.whatsapp.cancelled` | Pedido cancelado | Usuario completa flujo de cancelación |
| `orders.whatsapp.novelty` | Novedad solicitada | Usuario selecciona tipo de novedad |
| `customer.whatsapp.handoff` | Handoff a humano | Usuario presiona "Asesor" |
| SSE exchange | `whatsapp.message_received` | Mensaje entrante del cliente |
| SSE exchange | `whatsapp.conversation_started` | Nueva conversación iniciada |
| SSE exchange | `whatsapp.message_status_updated` | Estado de mensaje cambia (sent/delivered/read) |
| Persistence exchange | `conversation.created/updated/expired` | Estado de conversación |
| Persistence exchange | `message_log.created` | Nuevo mensaje (inbound o outbound) |
| Persistence exchange | `message_log.status_updated` | Estado actualizado de mensaje |

### Credenciales de plataforma (compartidas)

Todas las integraciones WhatsApp usan las mismas credenciales del **integration_type** WhatsApp (ID: 2):
- `phone_number_id` — ID del número en Meta
- `access_token` — Token permanente del system user
- `verify_token` — Para verificación de webhook
- `webhook_secret` — Para validación HMAC-SHA256

Cacheadas en Redis (`integration:platform_creds:2`) durante el warmup del servidor. Se leen via `core.GetCachedPlatformCredentials()`.

### Webhook (sin JWT)

Las rutas del webhook (`GET/POST /api/v1/integrations/whatsapp/webhook`) **no usan JWT** porque Meta envía su propia autenticación:
- **GET** (verificación): Meta envía `hub.verify_token` → se compara con `verify_token` de platform_creds
- **POST** (eventos): Meta firma el payload con HMAC-SHA256 → se valida con `webhook_secret` de platform_creds

---

## Estructura de Directorios

```
whatsApp/
├── bundle.go                                  # Ensamblaje de componentes (DI)
├── internal/
│   ├── domain/
│   │   ├── entities/
│   │   │   ├── conversation.go                # Entidades de conversación y estados
│   │   │   ├── message.go                     # Entidades de mensajes
│   │   │   └── template.go                    # Catálogo de plantillas
│   │   ├── dtos/
│   │   ├── ports/
│   │   │   └── ports.go                       # Interfaces (IConversationCache, ISSEEventPublisher, etc.)
│   │   └── errors/
│   ├── app/
│   │   ├── usecasemessaging/
│   │   │   ├── constructor.go                 # IUseCase interface + dependencias
│   │   │   ├── handle-webhook.go              # Procesamiento de webhooks + SSE + HumanSession
│   │   │   ├── send-template-message.go       # Envío de plantillas
│   │   │   ├── send_manual_reply.go           # Respuesta manual + ActivateHumanSession
│   │   │   ├── ai_control.go                  # PauseAI / ResumeAI
│   │   │   ├── conversation-manager.go        # Máquina de estados
│   │   │   └── utils.go                       # NormalizePhoneNumber, helpers
│   │   └── usecasetestconnection/
│   │       └── test-connection.go             # Test con template prueba_conexion
│   └── infra/
│       ├── primary/
│       │   ├── handlers/
│       │   │   ├── constructor.go             # IHandler + struct handler
│       │   │   ├── routes.go                  # Webhook SIN JWT, demás CON JWT
│       │   │   ├── webhook_handler.go         # Procesa mensajes y estados de Meta
│       │   │   ├── send_template_handler.go   # Envía plantilla manual
│       │   │   ├── manual_reply_handler.go    # POST /conversations/:id/reply
│       │   │   ├── ai_control_handler.go      # POST /conversations/:id/pause-ai y /resume-ai
│       │   │   └── request/
│       │   │       ├── manual_reply.go
│       │   │       └── ai_control.go
│       │   └── queue/
│       │       ├── consumerorder/             # Consume orders.confirmation.requested
│       │       └── consumeralert/             # Consume alertas de monitoreo
│       └── secondary/
│           ├── client/                        # Cliente HTTP WhatsApp API
│           ├── cache/
│           │   ├── conversation_cache.go      # Redis: conversaciones + HumanSession + AIPaused
│           │   └── credentials_cache.go       # Lee platform_creds desde Redis
│           └── queue/
│               ├── sse_publisher.go           # Publica eventos SSE via RabbitMQ
│               ├── event_publisher.go         # Publica eventos de negocio
│               └── persistence_publisher.go   # Publica eventos de persistencia
```

---

## Endpoints HTTP

| Método | Ruta | Auth | Descripción |
|--------|------|------|-------------|
| `GET` | `/whatsapp/webhook` | — | Verificación de webhook Meta |
| `POST` | `/whatsapp/webhook` | — | Eventos entrantes de Meta |
| `POST` | `/whatsapp/send-template` | JWT | Envío manual de plantilla |
| `POST` | `/whatsapp/conversations/:id/reply` | JWT | Respuesta manual del agente |
| `POST` | `/whatsapp/conversations/:id/pause-ai` | JWT | Pausa IA, activa sesión humana |
| `POST` | `/whatsapp/conversations/:id/resume-ai` | JWT | Reactiva IA |

---

## Flujo de Conversación (bot automatizado)

```
[INICIO] → confirmacion_pedido_contraentrega
            ├─ "Confirmar pedido" → pedido_confirmado_v2 [FIN: event confirmed]
            └─ "No confirmar" → menu_no_confirmacion
                                 ├─ "Presentar novedad" → tipo_novedad_pedido
                                 │                          ├─ "Cambio de dirección" → novedad_cambio_direccion [FIN: event novelty]
                                 │                          ├─ "Cambio de productos" → novedad_cambio_productos [FIN: event novelty]
                                 │                          └─ "Cambio medio de pago" → novedad_cambio_medio_pago [FIN: event novelty]
                                 ├─ "Cancelar pedido" → confirmar_cancelacion_pedido
                                 │                       ├─ "Sí, cancelar" → motivo_cancelacion_pedido
                                 │                       │                    └─ [texto libre] → pedido_cancelado [FIN: event cancelled]
                                 │                       └─ "No, volver" → menu_no_confirmacion
                                 └─ "Asesor" → handoff_asesor [FIN: event handoff → HUMANO]
```

---

## Templates del Sistema

| Template | Variables | Botones | Uso |
|----------|-----------|---------|-----|
| `prueba_conexion` | — | — | Test de conexión |
| `confirmacion_pedido_contraentrega` | nombre, tienda, orden, dirección, productos | Confirmar / No confirmar | Inicio del flujo |
| `pedido_confirmado_v2` | nombre, #pedido, tienda, dirección, productos | — | Confirmación exitosa |
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

---

## API de Meta — Referencia Rápida

Todas las llamadas usan `https://graph.facebook.com/v22.0/` con header `Authorization: Bearer {ACCESS_TOKEN}`.

### Listar números de teléfono de un WABA

```bash
curl -s "https://graph.facebook.com/v22.0/{WABA_ID}/phone_numbers?fields=verified_name,display_phone_number,quality_rating,status,messaging_limit_tier,platform_type" \
  -H "Authorization: Bearer {TOKEN}" | jq .
```

### Ver suscripciones de webhook de la app

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

### Enviar mensaje de texto libre (dentro de ventana de 24h)

```bash
curl -X POST "https://graph.facebook.com/v22.0/{PHONE_NUMBER_ID}/messages" \
  -H "Authorization: Bearer {TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "messaging_product": "whatsapp",
    "to": "573001234567",
    "type": "text",
    "text": {"body": "Hola, ¿en qué te puedo ayudar?"}
  }'
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

### Debug de token

```bash
curl -s "https://graph.facebook.com/v22.0/debug_token?input_token={TOKEN}&access_token={TOKEN}" | jq .
```

### Listar templates

```bash
curl -s "https://graph.facebook.com/v22.0/{WABA_ID}/message_templates?fields=name,status,language,category&limit=50" \
  -H "Authorization: Bearer {TOKEN}" | jq .
```

---

## Troubleshooting

### Webhook da 401
Las rutas del webhook NO deben pasar por JWT. Verificar que `routes.go` registre `GET/POST /webhook` sin `middleware.JWT()`.

### Template rechazado por "too many variables"
Meta rechaza templates donde el ratio variables/texto es alto. Solución: agregar más texto al body.

### `hello_world` no funciona con número empresarial
Es normal — `hello_world` solo funciona con números de test de Meta. Usar `prueba_conexion` para testing.

### Mensajes del cliente no aparecen en el dashboard
1. Verificar que exista `HumanSession` en Redis: `whatsapp:human_session:{phone}`
2. Verificar que el SSE publisher esté conectado a RabbitMQ
3. Verificar que el frontend tenga el hook `useSSE` suscrito a `whatsapp.message_received`

### IA sigue respondiendo aunque esté pausada
Verificar que existe la clave `whatsapp:ai_paused:{phone}` en Redis. El módulo `ai_sales` la lee con `IsAIPaused()` al inicio de cada mensaje entrante.

### Credenciales no encontradas en cache
El warmup del servidor debe cachear `platform_creds` en Redis. Verificar que `warm_cache.go` ejecute `warmPlatformCredentials()`.

### Token sin permisos para un WABA
1. Business Settings → System Users → cam-adm → Activos asignados
2. Verificar que el WABA esté asignado con control total
3. Regenerar token
