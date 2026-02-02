# Módulo de Simulación de WhatsApp

Módulo de pruebas de integración para el sistema de mensajería de WhatsApp, siguiendo arquitectura hexagonal.

## 🏗️ Arquitectura Hexagonal

```
whatsapp/
├── bundle.go                        # Punto de entrada del módulo
└── internal/
    ├── domain/                      # Capa de dominio (sin dependencias externas)
    │   ├── entities.go              # Entidades de dominio (Conversation, MessageLog, etc.)
    │   ├── ports.go                 # Interfaces (IWebhookClient)
    │   └── repository.go            # Repositorio en memoria
    ├── app/
    │   └── usecases/                # Casos de uso
    │       ├── constructor.go       # Constructor de use cases
    │       └── conversation_simulator.go  # Lógica de simulación
    └── infra/
        └── primary/
            └── client/
                └── webhook_client.go  # Cliente HTTP para enviar webhooks
```

## 🎯 Funcionalidades

### 1. Simulación Manual
Simular respuestas de usuario específicas para testing:

```go
whatsappIntegration.SimulateUserResponse("+573001234567", "Confirmar pedido")
```

### 2. Simulación Automática
Simular respuestas basadas en el template enviado:

```go
whatsappIntegration.SimulateAutoResponse("+573001234567", "confirmacion_pedido_contraentrega")
// Automáticamente responde: "Confirmar pedido"
```

### 3. Listar Conversaciones
Ver todas las conversaciones simuladas:

```go
conversations := whatsappIntegration.GetAllConversations()
```

### 4. Ver Mensajes
Ver todos los mensajes de una conversación:

```go
messages := whatsappIntegration.GetMessages(conversationID)
```

## 📋 Mapeo de Templates → Respuestas

| Template | Respuesta Automática |
|----------|---------------------|
| `confirmacion_pedido_contraentrega` | "Confirmar pedido" |
| `menu_no_confirmacion` | "Presentar novedad" |
| `solicitud_novedad` | "Otro" |
| `novedad_otro` | "El producto llegó dañado" |
| `solicitud_cancelacion` | "Sí, cancelar" |
| `motivo_cancelacion` | "Ya no lo necesito" |
| `confirmacion_cancelacion_pago` | "Sí, he cancelado el pago" |
| `pedido_cancelado` | (sin respuesta) |
| `pedido_confirmado` | (sin respuesta) |
| Cualquier otro | "Asesor humano" |

## ⚙️ Configuración

Variables en `.env`:

```env
# URL base del sistema real
WEBHOOK_BASE_URL=http://localhost:3050

# Secret para firmar webhooks (debe coincidir con el sistema real)
WHATSAPP_WEBHOOK_SECRET=test_webhook_secret

# Delay en segundos antes de enviar respuestas automáticas
WHATSAPP_AUTO_REPLY_DELAY=2
```

## 🚀 Uso desde el Menú Interactivo

```bash
cd /back/integrationTest
go run cmd/main.go
```

**Opciones disponibles:**

```
💬 WHATSAPP:
8. Simular respuesta de usuario (manual)
9. Simular respuesta automática (por template)
10. Listar conversaciones almacenadas
```

### Ejemplo de Flujo

1. **Iniciar el simulador interactivo**
2. **Seleccionar opción 8** (respuesta manual)
3. **Ingresar:** `+573001234567`
4. **Ingresar:** `Confirmar pedido`
5. **El simulador:**
   - Espera 2 segundos (configurable)
   - Envía webhook de estado "delivered"
   - Espera 500ms
   - Envía webhook de estado "read"
   - Espera 500ms
   - Envía webhook con respuesta del usuario

## 🔄 Flujo de Webhooks

```
1. Sistema Real envía template al usuario real
   ↓
2. Usuario usa el simulador (opción 8 o 9)
   ↓
3. Simulador espera WHATSAPP_AUTO_REPLY_DELAY (2s default)
   ↓
4. Simulador envía webhook: status "delivered"
   ↓
5. Simulador espera 500ms
   ↓
6. Simulador envía webhook: status "read"
   ↓
7. Simulador espera 500ms
   ↓
8. Simulador envía webhook: mensaje de usuario (button response)
   ↓
9. Sistema Real procesa la respuesta y actualiza el estado
```

## 📊 Estructura de Webhooks

### Webhook de Estado

```json
{
  "object": "whatsapp_business_account",
  "entry": [{
    "id": "123456789",
    "changes": [{
      "field": "messages",
      "value": {
        "messaging_product": "whatsapp",
        "metadata": {
          "display_phone_number": "+15551234567",
          "phone_number_id": "123456789012345"
        },
        "statuses": [{
          "id": "wamid.HBg...",
          "status": "delivered",
          "timestamp": "1738454567",
          "recipient_id": "+573001234567",
          "conversation": {
            "id": "conversation_...",
            "origin": { "type": "business_initiated" }
          }
        }]
      }
    }]
  }]
}
```

### Webhook de Mensaje de Usuario

```json
{
  "object": "whatsapp_business_account",
  "entry": [{
    "id": "123456789",
    "changes": [{
      "field": "messages",
      "value": {
        "messaging_product": "whatsapp",
        "metadata": {
          "display_phone_number": "+15551234567",
          "phone_number_id": "123456789012345"
        },
        "contacts": [{
          "profile": { "name": "Test User" },
          "wa_id": "+573001234567"
        }],
        "messages": [{
          "from": "+573001234567",
          "id": "wamid.HBg...",
          "timestamp": "1738454570",
          "type": "button",
          "button": {
            "payload": "Confirmar pedido",
            "text": "Confirmar pedido"
          }
        }]
      }
    }]
  }]
}
```

## 🔐 Firma HMAC-SHA256

Los webhooks se firman con HMAC-SHA256 usando el `WHATSAPP_WEBHOOK_SECRET`:

```
X-Hub-Signature-256: sha256=<hex_signature>
```

El sistema real debe validar esta firma para verificar que el webhook es legítimo.

## 🧪 Testing

El simulador permite probar:

- ✅ Flujos completos de conversación
- ✅ Diferentes respuestas de usuario
- ✅ Estados de mensaje (delivered, read)
- ✅ Firma HMAC de webhooks
- ✅ Templates sin respuesta esperada
- ✅ Respuestas automáticas vs. manuales

## 📝 Notas Importantes

1. **Repositorio en Memoria**: Las conversaciones se almacenan en memoria y se pierden al reiniciar
2. **No Valida Business Logic**: Solo simula webhooks, no valida reglas de negocio
3. **Delay Configurable**: Ajustar `WHATSAPP_AUTO_REPLY_DELAY` según necesidad
4. **Same Process**: El simulador corre en el mismo proceso que otros simuladores (Shopify, etc.)

## 🎯 Casos de Uso

### Testing de Confirmación de Pedido

```
1. Sistema envía template "confirmacion_pedido_contraentrega"
2. Simulador responde "Confirmar pedido" (automático o manual)
3. Sistema actualiza orden a CONFIRMED
4. Verificar en BD: estado de conversación y orden
```

### Testing de Cancelación

```
1. Sistema envía template "confirmacion_pedido_contraentrega"
2. Simulador responde "No confirmar"
3. Sistema envía "menu_no_confirmacion"
4. Simulador responde "Cancelar pedido"
5. Sistema procesa cancelación
```

### Testing de Novedad

```
1. Sistema envía template "confirmacion_pedido_contraentrega"
2. Simulador responde "No confirmar"
3. Sistema envía "menu_no_confirmacion"
4. Simulador responde "Presentar novedad"
5. Sistema abre caso de soporte
```

## 🔗 Integración con Sistema Real

El sistema real debe:

1. ✅ Tener endpoint: `POST /api/integrations/whatsapp/webhook`
2. ✅ Validar firma HMAC-SHA256
3. ✅ Procesar webhooks de estados (delivered, read)
4. ✅ Procesar webhooks de mensajes (button, text)
5. ✅ Actualizar estado de conversaciones
6. ✅ Actualizar estado de órdenes según respuestas

---

**Implementado siguiendo:**
- ✅ Arquitectura Hexagonal
- ✅ Domain sin dependencias externas
- ✅ Ports and Adapters pattern
- ✅ Repository pattern (in-memory)
- ✅ Use Cases para lógica de negocio
- ✅ Infraestructura separada del dominio
